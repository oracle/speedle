//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package cfg

import (
	"encoding/json"
	"os"
	"testing"
)

func WriteConfig(mConfig Config, fileName string) error {
	jsonFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	psB, err := json.MarshalIndent(mConfig, "", "    ")
	if err != nil {
		return err
	}
	_, err = jsonFile.Write(psB)
	return err
}

func TestReadConfig(t *testing.T) {
	fileConfig, err := ReadConfig("./config_file.json")
	if err != nil {
		t.Error("Fail to read file store config")
	}
	if fileConfig.StoreConfig.StoreType != StorageTypeFile {
		t.Error("Read config error")
	}

	etcdConfig, err := ReadConfig("./config_etcd.json")
	if err != nil {
		t.Error("Fail to read etcd store config")
	}
	if etcdConfig.StoreConfig.StoreType != "etcd" {
		t.Error("Read config error")
	}
}

func TestWriteConfig(t *testing.T) {
	storeConfig := StoreConfig{
		StoreType: StorageTypeFile,
		StoreProps: map[string]interface{}{
			"FileLocation": "./ps.json",
		},
	}
	mConfig := Config{
		StoreConfig: &storeConfig,
	}
	WriteConfig(mConfig, "./myconfig.json")
}
