// nolint:lll
// Generates the grpcadapter adapter's resource yaml. It contains the adapter's configuration, name,
// supported template names (authorization in this case), and whether it is session or no-session based.
//go:generate $GOPATH/src/istio.io/istio/bin/mixer_codegen.sh -a mixer/adapter/speedle/config/config.proto -x "-s=false -n speedlegrpcadapter -t authorization -o speedlegrpcadapter.yaml"

package speedlegrpcadapter

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"
	adptModel "istio.io/api/mixer/adapter/model/v1beta1"
	"istio.io/api/policy/v1beta1"
	"istio.io/istio/mixer/adapter/speedle"
	"istio.io/istio/mixer/adapter/speedle/config"
	"istio.io/istio/mixer/pkg/adapter"
	"istio.io/istio/mixer/pkg/pool"
	"istio.io/istio/mixer/pkg/runtime/handler"
	"istio.io/istio/mixer/template/authorization"
)

type (

	// Server is basic server interface
	Server interface {
		Addr() string
		Close() error
		Run(shutdown chan error)
	}

	// Speedle grpc server supports authorization template.
	GrpcServer struct {
		listener net.Listener
		shutdown chan error
		server   *grpc.Server

		rawcfg      []byte
		builder     adapter.HandlerBuilder
		env         adapter.Env
		builderLock sync.RWMutex

		h authorization.Handler
	}
)

var _ authorization.HandleAuthorizationServiceServer = &GrpcServer{}

///////////////// Configuration Methods ///////////////

func (s *GrpcServer) getHandler(rawcfg []byte) (authorization.Handler, error) {
	s.builderLock.RLock()
	if 0 == bytes.Compare(rawcfg, s.rawcfg) {
		h := s.h
		s.builderLock.RUnlock()
		return h, nil
	}
	s.builderLock.RUnlock()

	cfg := &config.Params{}

	if err := cfg.Unmarshal(rawcfg); err != nil {
		return nil, err
	}

	s.builderLock.Lock()
	defer s.builderLock.Unlock()

	if 0 == bytes.Compare(rawcfg, s.rawcfg) {
		return s.h, nil
	}

	s.env.Logger().Infof("Loaded handler with: %v", cfg)

	s.builder.SetAdapterConfig(cfg)
	if ce := s.builder.Validate(); ce != nil {
		return nil, ce
	}

	h, err := s.builder.Build(context.Background(), s.env)
	if err != nil {
		s.env.Logger().Errorf("could not build: %v", err)
		return nil, err
	}
	s.rawcfg = rawcfg
	s.h = h.(authorization.Handler)
	return s.h, err
}

func instance(in *authorization.InstanceMsg) *authorization.Instance {
	out := &authorization.Instance{
		Name:    in.Name,
		Subject: decodeSubject(in.Subject),
		Action:  decodeAction(in.Action),
	}
	return out
}

func decodeAction(action *authorization.ActionMsg) *authorization.Action {
	ac := &authorization.Action{
		Namespace:  action.Namespace,
		Service:    action.Service,
		Method:     action.Method,
		Path:       action.Path,
		Properties: decodeProperties(action.Properties),
	}
	return ac
}

func decodeSubject(subject *authorization.SubjectMsg) *authorization.Subject {
	su := &authorization.Subject{
		User:       subject.User,
		Groups:     subject.Groups,
		Properties: decodeProperties(subject.Properties),
	}
	return su
}

func decodeProperties(in map[string]*v1beta1.Value) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = decodeValue(v.GetValue())
	}
	return out
}

func decodeValue(in interface{}) interface{} {
	switch t := in.(type) {
	case *v1beta1.Value_StringValue:
		return t.StringValue
	case *v1beta1.Value_Int64Value:
		return t.Int64Value
	case *v1beta1.Value_DoubleValue:
		return t.DoubleValue
	default:
		return fmt.Sprintf("%v", in)
	}
}

func (s *GrpcServer) HandleAuthorization(ctx context.Context, r *authorization.HandleAuthorizationRequest) (*adptModel.CheckResult, error) {

	h, err := s.getHandler(r.AdapterConfig.Value)
	if err != nil {
		return nil, err
	}

	if result, err := h.HandleAuthorization(ctx, instance(r.Instance)); err != nil {
		s.env.Logger().Errorf("Could not process: %v", err)
		return nil, err
	} else {
		return &adptModel.CheckResult{Status: result.Status}, nil
	}
}

// Addr returns the listening address of the server
func (s *GrpcServer) Addr() string {
	return s.listener.Addr().String()
}

// Run starts the server run
func (s *GrpcServer) Run() {
	s.shutdown = make(chan error, 1)
	go func() {
		err := s.server.Serve(s.listener)

		s.shutdown <- err
	}()
}

// Wait waits for server to stop
func (s *GrpcServer) Wait() error {
	if s.shutdown == nil {
		return fmt.Errorf("server not running")
	}

	err := <-s.shutdown
	s.shutdown = nil
	return err
}

// Close gracefully shuts down the server
func (s *GrpcServer) Close() error {
	if s.shutdown != nil {
		s.server.GracefulStop()
		_ = s.Wait()
	}

	if s.listener != nil {
		_ = s.listener.Close()
	}

	return nil
}

// NewGrpcServer creates a new speedle grpc server from given args.
func NewGrpcServer(addr string) (*GrpcServer, error) {
	if addr == "" {
		addr = "0"
	}

	gp := pool.NewGoroutinePool(1, true)
	inf := speedle.GetInfo()
	s := &GrpcServer{builder: inf.NewBuilder(),
		env:    handler.NewEnv(0, "speedlegrpcadapter", gp),
		rawcfg: []byte{0xff, 0xff},
	}
	var err error
	if s.listener, err = net.Listen("tcp", addr); err != nil {
		_ = s.Close()
		return nil, fmt.Errorf("unable to listen on socket: %v", err)
	}

	fmt.Printf("listening on :%v\n", s.listener.Addr())
	s.server = grpc.NewServer()
	authorization.RegisterHandleAuthorizationServiceServer(s.server, s)

	return s, nil
}
