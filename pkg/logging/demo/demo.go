//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	logging "github.com/oracle/speedle/pkg/logging"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	LogConfig *logging.LogConfig
}

func main() {
	err := initLog()
	if err != nil {
		panic(err)
	}

	// Basic usages
	log.Debug("This is a debug log entry")
	log.Info("This is a info log entry")
	log.Warn("This is a warning log entry")
	log.Warning("This is a warning log entry again")
	log.Error("This is a error log entry")
	//log.Fatal("This is a fatal log entry")
	//log.Panic("This is a panic log entry")

	// Include some context values
	log.Debugf("This is a debug log entry, %s", "debug value")
	log.Infof("This is a info log entry, %s", "info value")
	log.Warnf("This is a warning log entry, %s", "warning value")
	log.Warningf("This is a warning log entry again, %s", "warning value")
	log.Errorf("This is a error log entry, %s", "error value")
	//log.Fatalf("This is a fatal log entry, %s", "fatal value")
	//log.Panicf("This is a panic log entry, %s", "panic value")

	// Another option to include some context values
	log.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	str := "This is the value of a variable"
	log.WithFields(log.Fields{
		"omg":    true,
		"number": 12234,
		"test":   str,
	}).Warn("The group's number increased tremendously!")

	log.WithFields(log.Fields{
		"omg":    true,
		"number": 100,
	}).Fatal("The ice breaks!")

	// A common pattern is to re-use fields between logging statements by re-using
	// the logrus.Entry returned from WithFields()
	contextLogger := log.WithFields(log.Fields{
		"common": "this is a common field",
		"other":  "I also should be logged always",
	})

	contextLogger.Info("I'll be logged with common and other field")
	contextLogger.Info("Me too")
}

func initLog() error {
	raw, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return fmt.Errorf("Failed to load the config file, err: %v. \n", err)
	}

	var config Config
	err = json.Unmarshal(raw, &config)
	if err != nil {
		return fmt.Errorf("Failed to parse the config file, err: %v. \n", err)
	}

	if config.LogConfig != nil {
		err = logging.InitLog(config.LogConfig)
		if err != nil {
			return fmt.Errorf("Failed to initialize the log library, err: %v. \n", err)
		}
	} else {
		return fmt.Errorf("No any log configurations.\n")
	}

	return nil
}
