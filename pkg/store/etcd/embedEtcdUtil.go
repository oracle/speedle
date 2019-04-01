//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package etcd

import (
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/coreos/etcd/embed"
	"github.com/oracle/speedle/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var embededStarted = false
var isStartedByOtherProcess = false

const embeddedEtcdPort = 2379

//StartEmbeddedEtcd start a embed etcd which use a clean tmp directory to store data
func StartEmbeddedEtcd(dataDir string) (etcd *embed.Etcd, etcdDir string, err error) {
	if embededStarted {
		// Already started
		return nil, "", nil
	}
	if isEtcdPortOccupied() {
		//we assume the embeded etcd is already started by other process, and we use that etcd directly.
		//This is to support starting both mgmt server and atz server, and use the same embeded etcd in dev or test env.
		return nil, "", nil
	}
	etcdDir = dataDir
	if etcdDir == "" {
		etcdDir, err = ioutil.TempDir(os.TempDir(), "etcd.tmp")
		log.Infof("The embedded etcd store dir is %q", etcdDir)
		if err != nil {
			log.Error(err)
			return etcd, etcdDir, errors.Wrapf(err, errors.StoreError, "failed to create etcd dir")
		}
	}

	cfg := embed.NewConfig()
	cfg.Debug = true
	cfg.Dir = etcdDir
	etcd, err = embed.StartEtcd(cfg)
	if err != nil {
		log.Error(err)
		return etcd, etcdDir, errors.Wrapf(err, errors.StoreError, "failed to start embedded etcd server")
	}

	embededStarted = true
	select {
	case <-etcd.Server.ReadyNotify():
		log.Info("Etcd Server is ready!")
	case <-time.After(60 * time.Second):
		etcd.Server.Stop() // trigger a shutdown
		err = errors.New(errors.StoreError, "etcd Server took too long to start")
	}
	return etcd, etcdDir, err
}

//CleanEmbedEtcd free the resource of embed etcd, and remove the tmp directory which is used to store data
func CleanEmbeddedEtcd(etcd *embed.Etcd, etcdDir string) {
	if embededStarted {
		etcd.Close()
		os.RemoveAll(etcdDir)
		embededStarted = false
	}
}

func isEtcdPortOccupied() bool {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(embeddedEtcdPort))

	if err != nil {
		return true
	}
	ln.Close()
	return false
}
