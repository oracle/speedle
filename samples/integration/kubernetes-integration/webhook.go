package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"speedle/api/authz"
	"speedle/rest/authz/client"
	"strings"

	"k8s.io/api/authorization/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type SubjectAccessReviewResponse struct {
	metav1.TypeMeta `json:",inline"`
	// Status is filled in by the server and indicates whether the request is allowed or not
	// +optional
	Status v1beta1.SubjectAccessReviewStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type Parameters struct {
	ClusterName    string
	KeyPath        string
	CertPath       string
	ClientCertPath string
	ListenPort     string
	KubeConfigFile string
	SpeedleHost    string
}

func usage() {
	flag.Usage()
	os.Exit(1)
}

func main() {
	var params Parameters
	params.parseFlags()
	validateParameters(&params)
	handler, err := New(&params)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", handler.atzHandler)
	log.Fatal(twoWaysTLS(&params))
}

func validateParameters(params *Parameters) {
	if len(params.ListenPort) == 0 {
		fmt.Fprintln(os.Stderr, "No listen port found.")
		usage()
	}
	if len(params.KeyPath) == 0 {
		fmt.Fprintln(os.Stderr, "No webhook key file found.")
		usage()
	}
	if len(params.CertPath) == 0 {
		fmt.Fprintln(os.Stderr, "No webhook certificate file found.")
		usage()
	}
	if len(params.ClientCertPath) == 0 {
		fmt.Fprintln(os.Stderr, "No k8s certificate file found.")
		usage()
	}
	if len(params.ClusterName) == 0 {
		fmt.Fprintln(os.Stderr, "No cluster name found.")
		usage()
	}
	if len(params.KubeConfigFile) == 0 {
		fmt.Fprintln(os.Stderr, "No kubernetes configure file found.")
		usage()
	}
	if len(params.SpeedleHost) == 0 {
		fmt.Fprintln(os.Stderr, "No Speedle host found.")
		usage()
	}
}

func twoWaysTLS(params *Parameters) error {
	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	caCert, err := ioutil.ReadFile(params.ClientCertPath)
	if err != nil {
		return err
	}
	caCertPool := x509.NewCertPool()
	// Set HTTPS client
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAnyClientCert,
	}
	tlsConfig.BuildNameToCertificate()
	server := &http.Server{
		Addr:      hostName + ":" + params.ListenPort,
		TLSConfig: tlsConfig,
	}
	caCertPool.AppendCertsFromPEM(caCert)
	return server.ListenAndServeTLS(params.CertPath, params.KeyPath)
}

func (k *Parameters) parseFlags() {
	flag.StringVar(&k.KeyPath, "key", "", "Server key file path.")
	flag.StringVar(&k.CertPath, "cert", "", "Server certifice file path.")
	flag.StringVar(&k.ClientCertPath, "client-cert", "", "Client certifice file path.")
	flag.StringVar(&k.ListenPort, "port", "8443", "Listen port.")
	flag.StringVar(&k.ClusterName, "cluster-name", "", "Cluster name.")
	home, _ := os.LookupEnv("HOME")
	defaultConfigPath := strings.Join([]string{home, ".kube", "config"}, string(os.PathSeparator))
	flag.StringVar(&k.KubeConfigFile, "kubeconfig", defaultConfigPath, "K8s configure file.")
	flag.StringVar(&k.SpeedleHost, "speedle-host", "", "Speedle host.")
	flag.Parse()
}

type handlerImpl struct {
	speedleClient client.ADSClient
	clusterName   string
	kubeClientset *kubernetes.Clientset
}

func New(params *Parameters) (*handlerImpl, error) {
	config, err := clientcmd.BuildConfigFromFlags("", params.KubeConfigFile)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	connProperties := map[string]string{
		authz.HOST_PROP:      params.SpeedleHost,
		authz.IS_SECURE_PROP: "false",
	}

	speedleClient, err := client.New(connProperties)
	if err != nil {
		return nil, err
	}

	handler := handlerImpl{
		speedleClient: speedleClient,
		clusterName:   params.ClusterName,
		kubeClientset: clientset,
	}
	return &handler, nil
}

func (impl *handlerImpl) atzHandler(w http.ResponseWriter, r *http.Request) {
	if method := r.Method; method != "POST" {
		http.Error(w, "Method "+method+" is not allowed.", http.StatusMethodNotAllowed)
		return
	}
	fmt.Printf("Method, %q\n", r.Method)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4*1024))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	fmt.Println("-------------------------------")
	fmt.Print(string(body))
	var ar v1beta1.SubjectAccessReview
	if err := json.Unmarshal(body, &ar); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	allowed, err := impl.evaluatePolicies(impl.clusterName, &ar)
	err = nil
	fmt.Println("-------------------------------")

	var ret SubjectAccessReviewResponse
	ret.APIVersion = ar.APIVersion
	ret.Kind = ar.Kind
	ret.Status.Allowed = allowed
	if err != nil {
		ret.Status.Reason = err.Error()
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
}

func (impl *handlerImpl) fillAttributes(context *authz.RequestContext, spec *v1beta1.SubjectAccessReviewSpec) {
	if spec == nil || spec.ResourceAttributes == nil {
		return
	}
	namespace := spec.ResourceAttributes.Namespace
	resource := spec.ResourceAttributes.Resource
	resourceName := spec.ResourceAttributes.Name

	if len(resource) == 0 || len(resourceName) == 0 {
		return
	}

	impl.fillAttributesWithResource(context, namespace, resource, resourceName)
}

func (impl *handlerImpl) getResourceByName(namespace string, resource string, resourceName string) (*metav1.ObjectMeta, []metav1.OwnerReference, error) {
	var (
		resMeta *metav1.ObjectMeta
		owners  []metav1.OwnerReference
	)

	switch resource {
	case "pods":
		pod, err := impl.kubeClientset.CoreV1().Pods(namespace).Get(resourceName, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when get pod details due to error %v.", err)
			return nil, owners, err
		}
		resMeta = &pod.ObjectMeta
		owners = pod.OwnerReferences
		return resMeta, owners, nil
	case "daemonsets":
		daemonset, err := impl.kubeClientset.Extensions().DaemonSets(namespace).Get(resourceName, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when get daemon set details due to error %v.", err)
			return nil, owners, err
		}
		resMeta = &daemonset.ObjectMeta
		owners = daemonset.OwnerReferences
		return resMeta, owners, nil
	case "deployments":
		deployment, err := impl.kubeClientset.Extensions().Deployments(namespace).Get(resourceName, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when get deployment details due to error %v.", err)
			return nil, owners, err
		}
		resMeta = &deployment.ObjectMeta
		owners = deployment.OwnerReferences
		return resMeta, owners, nil
	case "replicasets":
		repSet, err := impl.kubeClientset.Extensions().ReplicaSets(namespace).Get(resourceName, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when get replication set details due to error %v.", err)
			return nil, owners, err
		}
		resMeta = &repSet.ObjectMeta
		owners = repSet.OwnerReferences
		return resMeta, owners, nil
	case "replicationcontrollers":
		rc, err := impl.kubeClientset.Core().ReplicationControllers(namespace).Get(resourceName, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when get replication controller details due to error %v.", err)
			return nil, owners, err
		}
		resMeta = &rc.ObjectMeta
		owners = rc.OwnerReferences
		return resMeta, owners, nil
	case "statefulsets":
		sf, err := impl.kubeClientset.Apps().StatefulSets(namespace).Get(resourceName, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when get stateful set details due to error %v.", err)
			return nil, owners, err
		}
		resMeta = &sf.ObjectMeta
		owners = sf.OwnerReferences
		return resMeta, owners, nil
	case "nodes":
		node, err := impl.kubeClientset.Core().Nodes().Get(resourceName, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when get node details due to error %v.", err)
			return nil, owners, err
		}
		resMeta = &node.ObjectMeta
		owners = node.OwnerReferences
		return resMeta, owners, nil
	case "namespaces":
		ns, err := impl.kubeClientset.Core().Namespaces().Get(resourceName, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when get node details due to error %v.", err)
			return nil, owners, err
		}
		resMeta = &ns.ObjectMeta
		owners = ns.OwnerReferences
		return resMeta, owners, nil
	}
	return nil, owners, fmt.Errorf("Not support resource %s", resource)
}

var kindResourceMap = map[string]string{
	"Pod":                   "pods",
	"DaemonSet":             "daemonsets",
	"Deployment":            "deployments",
	"ReplicaSet":            "replicasets",
	"ReplicationController": "replicationcontrollers",
	"StatefulSet":           "statefulsets",
	"Node":                  "nodes",
	"Namespace":             "namespaces",
}

func (impl *handlerImpl) fillAttributesWithOwner(context *authz.RequestContext, namespace string, resource string, resourceName string) {
	fmt.Printf("[Retrieve owner] Resource %s, resource name: %s.\n", resource, resourceName)
	resMeta, owners, err := impl.getResourceByName(namespace, resource, resourceName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get resource details due to error %v.\n", err)
		return
	}

	for _, owner := range owners {
		context.Attributes["owner_"+owner.Kind] = owner.Name
		resource, ok := kindResourceMap[owner.Kind]
		if !ok {
			fmt.Printf("Unsupport resource kind %s.\n", owner.Kind)
			continue
		}
		impl.fillAttributesWithOwner(context, resMeta.Namespace, resource, owner.Name)
	}
}

func (impl *handlerImpl) fillAttributesWithResource(context *authz.RequestContext, namespace string, resource string, resourceName string) {
	fmt.Printf("Resource %s, resource name: %s.\n", resource, resourceName)

	var (
		resMeta *metav1.ObjectMeta
		owners  []metav1.OwnerReference
	)

	switch resource {
	case "pods":
		pod, err := impl.kubeClientset.CoreV1().Pods(namespace).Get(resourceName, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when get pod details due to error %v.", err)
			return
		}
		resMeta = &pod.ObjectMeta
		owners = pod.OwnerReferences
		break
	case "daemonsets":
		daemonset, err := impl.kubeClientset.Extensions().DaemonSets(namespace).Get(resourceName, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when get daemon set details due to error %v.", err)
			return
		}
		resMeta = &daemonset.ObjectMeta
		owners = daemonset.OwnerReferences
		break
	case "deployments":
		deployment, err := impl.kubeClientset.Extensions().Deployments(namespace).Get(resourceName, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when get deployment details due to error %v.", err)
			return
		}
		resMeta = &deployment.ObjectMeta
		owners = deployment.OwnerReferences
		break
	case "replicasets":
		repSet, err := impl.kubeClientset.Extensions().ReplicaSets(namespace).Get(resourceName, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when get replication set details due to error %v.", err)
			return
		}
		resMeta = &repSet.ObjectMeta
		owners = repSet.OwnerReferences
		break
	case "replicationcontrollers":
		rc, err := impl.kubeClientset.Core().ReplicationControllers(namespace).Get(resourceName, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when get replication controller details due to error %v.", err)
			return
		}
		resMeta = &rc.ObjectMeta
		owners = rc.OwnerReferences
		break
	case "statefulsets":
		sf, err := impl.kubeClientset.Apps().StatefulSets(namespace).Get(resourceName, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when get stateful set details due to error %v.", err)
			return
		}
		resMeta = &sf.ObjectMeta
		owners = sf.OwnerReferences
		break
	case "nodes":
		node, err := impl.kubeClientset.Core().Nodes().Get(resourceName, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when get node details due to error %v.", err)
			return
		}
		resMeta = &node.ObjectMeta
		owners = node.OwnerReferences
		break
	case "namespaces":
		ns, err := impl.kubeClientset.Core().Namespaces().Get(resourceName, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when get node details due to error %v.", err)
			return
		}
		resMeta = &ns.ObjectMeta
		owners = ns.OwnerReferences
		break
	default:
		return
	}

	if context.Attributes == nil {
		context.Attributes = make(map[string]interface{})
	}
	context.Attributes["name"] = resMeta.GetName()
	context.Attributes["namespace"] = resMeta.GetName()
	labels := resMeta.GetLabels()
	if labels != nil {
		prefix := "labels_"
		for labelName, labelValue := range labels {
			context.Attributes[prefix+labelName] = labelValue
		}
	}

	for _, owner := range owners {
		context.Attributes["owner_"+owner.Kind] = owner.Name
		resource, ok := kindResourceMap[owner.Kind]
		if !ok {
			fmt.Printf("Unsupport resource kind %s.\n", owner.Kind)
			continue
		}
		impl.fillAttributesWithOwner(context, resMeta.Namespace, resource, owner.Name)
	}

	fmt.Printf("context: %v.\n", context)
}

func fillResourceAndAction(spec *v1beta1.SubjectAccessReviewSpec) (string, string) {
	if spec.NonResourceAttributes != nil {
		return fmt.Sprintf("non-res:%s", spec.NonResourceAttributes.Path), spec.NonResourceAttributes.Verb
	}
	resource := "/" + spec.ResourceAttributes.Resource + "/"
	if len(spec.ResourceAttributes.Name) != 0 {
		resource = resource + spec.ResourceAttributes.Name + "/"
	}
	if len(spec.ResourceAttributes.Subresource) != 0 {
		resource = resource + spec.ResourceAttributes.Subresource
	}
	return resource, spec.ResourceAttributes.Verb
}

func (impl *handlerImpl) evaluatePolicies(clusterName string, req *v1beta1.SubjectAccessReview) (bool, error) {
	resource, action := fillResourceAndAction(&req.Spec)

	subject := authz.Subject{
		Principals: []*authz.Principal{
			{
				Type: "user",
				Name: req.Spec.User,
			},
		},
	}

	for _, v := range req.Spec.Groups {
		subject.Principals = append(subject.Principals, &authz.Principal{Type: "group", Name: v})
	}

	context := authz.RequestContext{
		Subject:     &subject,
		ServiceName: clusterName,
		Resource:    resource,
		Action:      action,
		Attributes:  nil,
	}
	impl.fillAttributes(&context, &req.Spec)
	fmt.Printf("--- Service Name: %v\n", clusterName)
	fmt.Printf("--- Resource Name: %v\n", resource)
	fmt.Printf("--- Action: %v\n", action)
	for _, v := range subject.Principals {
		fmt.Printf("--- %v: %v\n", v.Type, v.Name)
	}

	allowed, err := impl.speedleClient.IsAllowed(context)
	fmt.Printf("--- Result: %v\n", allowed)
	fmt.Printf("--- Error: %v\n", err)
	return allowed, err
}
