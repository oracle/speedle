//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package file

import (
	"io/ioutil"
	"os"

	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/pkg/errors"
	"github.com/oracle/speedle/pkg/store"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

const (
	StoreType = "file"

	//following are keys of file store properties
	FileLocationKey = "FileLocation"

	FileLocationFlagName = "filestore-loc"

	DefaultFileStoreLocation = "/tmp/speedle-test-file-store.json"
)

type FileStoreBuilder struct{}

func (fs FileStoreBuilder) NewStore(config map[string]interface{}) (pms.PolicyStoreManager, error) {
	fileLocation, ok := config[FileLocationKey].(string)
	if !ok {
		return nil, errors.New(errors.ConfigError, "configure item FileLocation is not found")
	}
	if _, err := os.Stat(fileLocation); os.IsNotExist(err) {
		log.Info("policy store file does not exist, create one...")
		if err1 := ioutil.WriteFile(fileLocation, []byte("{}"), 0644); err1 != nil {
			log.Errorf("error creating policy store file: %v\n", err1)
			return nil, err1
		}
	}
	return &Store{FileLocation: fileLocation}, nil
}

func (fs FileStoreBuilder) GetStoreParams() map[string]string {
	return map[string]string{
		FileLocationFlagName: FileLocationKey,
	}

}

func init() {
	pflag.String(FileLocationFlagName, DefaultFileStoreLocation, "Store config: File location of file store.")

	store.Register(StoreType, FileStoreBuilder{})
}
