//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package main

import (
	"fmt"
	"net"

	"gitlab-odx.oracledx.com/wcai/speedle/api/pms"
	"gitlab-odx.oracledx.com/wcai/speedle/pkg/cmd/flags"
	"gitlab-odx.oracledx.com/wcai/speedle/pkg/errors"
	"gitlab-odx.oracledx.com/wcai/speedle/pkg/logging"
	"gitlab-odx.oracledx.com/wcai/speedle/pkg/store"
	"gitlab-odx.oracledx.com/wcai/speedle/pkg/svcs/pmsgrpc"
	"gitlab-odx.oracledx.com/wcai/speedle/pkg/svcs/pmsgrpc/pb"
	"gitlab-odx.oracledx.com/wcai/speedle/pkg/svcs/pmsrest"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var gitCommit string
var productVersion string
var goVersion string

func printVersionInfo() {
	fmt.Printf("speedle-pms:\n")
	fmt.Printf(" Version:       %s\n", productVersion)
	fmt.Printf(" Go Version:    %s\n", goVersion)
	fmt.Printf(" Git commit:    %s\n", gitCommit)
}

func main() {

	storeParamsMap := store.GetAllStoreParams()

	var params flags.Parameters
	params.ParseFlags(flags.DefaultPolicyMgmtEndPoint, printVersionInfo, storeParamsMap)
	params.ValidateFlags()

	conf, _ := params.Param2Config(storeParamsMap)

	// Initialize the logging
	if conf.LogConfig != nil {
		err := logging.InitLog(conf.LogConfig)
		if err != nil {
			log.Errorf("Policy_mgmt failed to initialize the log module, err: %v.", err)
		}
	} else {
		log.Error("No any log configurations for pmsserver.")
	}

	// Initialize the Audit logging
	if conf.AuditLogConfig != nil {
		err := logging.InitAuditLog(conf.AuditLogConfig)
		if err != nil {
			log.Errorf("Policy_mgmt failed to initialize the audit log module, err: %v.", err)
		}
	} else {
		log.Error("No any audit log configurations for Policy_mgmt.")
	}

	ps, err := store.NewStore(conf.StoreConfig.StoreType, conf.StoreConfig.StoreProps)
	if err != nil {
		log.Panic(err)
	}

	server, err := newGRPCServer(ps)
	if err != nil {
		log.Panic(err)
	}
	log.Info("Starting the gRPC server for pmsserver...")
	go listenGRPCServer(server)

	log.Info("Starting the REST server for pmsserver...")
	panic(listenAndServe(&params, ps))

}

func newGRPCServer(ps pms.PolicyStoreManager) (*grpc.Server, error) {
	server := grpc.NewServer()
	pb.RegisterPolicyManagerServer(server, pmsgrpc.NewServiceImpl(ps))
	reflection.Register(server)
	return server, nil
}

func listenGRPCServer(server *grpc.Server) error {
	// Register reflection service on gRPC server.
	lis, err := net.Listen("tcp", ":50001")
	if err != nil {
		return errors.Wrap(err, errors.ServerError, "failed to listen on endpoint :50001")
	}
	if err := server.Serve(lis); err != nil {
		return errors.Wrap(err, errors.ServerError, "failed to serve for endpoint :50001")
	}
	return nil
}

func listenAndServe(params *flags.Parameters, ps pms.PolicyStoreManager) error {
	routers, err := pmsrest.NewRouter(ps)
	if err != nil {
		log.Error("Fail to create handler...")
		return err
	}
	params.ListenAndServe(routers)
	return nil
}
