//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package flags

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/oracle/speedle/pkg/assertion"
	"github.com/oracle/speedle/pkg/cfg"
	"github.com/oracle/speedle/pkg/errors"
	"github.com/oracle/speedle/pkg/logging"

	"strconv"

	"github.com/natefinch/lumberjack"
	"github.com/spf13/pflag"
)

// Parameters is the parameters for Speedle
type Parameters struct {
	////////// Common flags //////////
	Version bool

	//////////Server config//////////////
	ConfigFile      StrParamDetail
	Endpoint        StrParamDetail
	Insecure        StrParamDetail
	EnableAuthz     StrParamDetail
	KeyPath         StrParamDetail
	CertPath        StrParamDetail
	ClientCertPath  StrParamDetail
	ForceClientCert StrParamDetail
	/////////Store config////////////////
	StoreType         StrParamDetail
	StoreWatchEnabled StrParamDetail

	////////Log config/////////////////////
	LogConf      LogParameters // normal log configuration
	AuditLogConf LogParameters // audit log configuration

	// AsserterParameters asserter webhook configuration
	AsserterConf AsserterParameters
}

// LogParameters is the parameters for log configuration
type LogParameters struct {
	LogLevel        StrParamDetail
	LogFormatter    StrParamDetail
	LogReportCaller StrParamDetail
	LogFileName     StrParamDetail
	LogMaxSize      StrParamDetail
	LogMaxAge       StrParamDetail
	LogMaxBackups   StrParamDetail
	LogLocalTime    StrParamDetail
	LogCompress     StrParamDetail
}

// AsserterParameters asserter webhook configurations
type AsserterParameters struct {
	AsserterEndpoint       StrParamDetail
	AsserterClientKeyPath  StrParamDetail
	AsserterCaPath         StrParamDetail
	AsserterClientCertPath StrParamDetail
	AsserterClientTimeout  StrParamDetail
}

const (
	DefaultPolicyMgmtEndPoint = "0.0.0.0:6733"
	DefaultAuthzCheckEndPoint = "0.0.0.0:6734"
	DefaultInsecure           = true
	DefaultEnableAuthz        = false

	DefaultStoreType = cfg.StorageTypeFile //file

	DefaultStoreWatchEnabled = true

	EnvVarPrefix = "SPDL"

	DefaultAuditLogLevel      = "info"
	DefaultAuditLogFormatter  = logging.FORMATTER_JSON
	DefaultAuditLogFilename   = ""   // Write audit log to os.Stderr by default
	DefaultAuditLogMaxSize    = "10" // MB
	DefaultAuditLogMaxAge     = "0"  // Never remove old files based on age
	DefaultAuditLogMaxBackups = "5"  // maximum number of old log files to retain

	DefaultAsserterClientTimeout = "5"
)

type StrParamDetail struct {
	Name         string
	ShortName    string
	DefaultValue string
	Usage        string
	Value        string
}

func (k *Parameters) NewHTTPServer(handler http.Handler) (*http.Server, error) {
	insecure, _ := strconv.ParseBool(k.Insecure.Value)
	if insecure {
		server := http.Server{
			Addr:    k.Endpoint.Value,
			Handler: handler,
		}
		return &server, nil
	}
	return k.newTLSServer(handler)
}

func (k *Parameters) newTLSServer(handler http.Handler) (*http.Server, error) {
	// Set HTTPS client
	tlsConfig := &tls.Config{}

	if k.ClientCertPath.Value != "" {
		caCert, err := ioutil.ReadFile(k.ClientCertPath.Value)
		if err != nil {
			return nil, errors.Wrapf(err, errors.ConfigError, "unable to read client CA certification from file %s", k.ClientCertPath.Value)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, errors.Wrap(err, errors.ConfigError, "failed to append certificates to pool")
		}

		tlsConfig.ClientCAs = caCertPool
		forceClientCert, _ := strconv.ParseBool(k.ForceClientCert.Value)
		if forceClientCert {
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		} else {
			tlsConfig.ClientAuth = tls.VerifyClientCertIfGiven
		}
	} else {
		tlsConfig.ClientAuth = tls.NoClientCert
	}

	tlsConfig.BuildNameToCertificate()

	server := http.Server{
		Addr:      k.Endpoint.Value,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}
	return &server, nil
}

func (k *Parameters) ListenAndServe(s *http.Server) error {
	insecure, _ := strconv.ParseBool(k.Insecure.Value)
	if insecure {
		return s.ListenAndServe()
	}
	return k.listenAndServeTLS(s)
}

func (k *Parameters) listenAndServeTLS(s *http.Server) error {
	return s.ListenAndServeTLS(k.CertPath.Value, k.KeyPath.Value)
}

// ParseFlags parses command line arguments
func (k *Parameters) ParseFlags(defaultEndpoint string, printVersionInfoFun func(), storeParamsMap map[string]string) {
	var params []*StrParamDetail
	k.ConfigFile = StrParamDetail{Name: "config-file", ShortName: "k", Usage: "Configuration file."}
	params = append(params, &k.ConfigFile)
	k.Endpoint = StrParamDetail{Name: "endpoint", DefaultValue: defaultEndpoint, Usage: "Server config: Endpoint the server listen and serve."}
	params = append(params, &k.Endpoint)
	k.Insecure = StrParamDetail{Name: "insecure", DefaultValue: strconv.FormatBool(DefaultInsecure), Usage: "Server config: Disable transport security."}
	params = append(params, &k.Insecure)
	k.EnableAuthz = StrParamDetail{Name: "enable-authz", DefaultValue: strconv.FormatBool(DefaultEnableAuthz), Usage: "Server config: Enable authorization check."}
	params = append(params, &k.EnableAuthz)
	k.CertPath = StrParamDetail{Name: "cert", Usage: "Server config: Server certifice file path."}
	params = append(params, &k.CertPath)
	k.KeyPath = StrParamDetail{Name: "key", Usage: "Server config: Server key file path."}
	params = append(params, &k.KeyPath)
	k.ClientCertPath = StrParamDetail{Name: "client-cert", ShortName: "c", Usage: "Server config: Client certifice file path."}
	params = append(params, &k.ClientCertPath)
	k.ForceClientCert = StrParamDetail{Name: "force-client-cert", ShortName: "f", Usage: "Server config: Force Client certification."}
	params = append(params, &k.ForceClientCert)

	k.StoreType = StrParamDetail{Name: "store-type", DefaultValue: DefaultStoreType, Usage: "Store config: Policy store type, etcd or file."}
	params = append(params, &k.StoreType)
	k.StoreWatchEnabled = StrParamDetail{Name: "enable-watch", DefaultValue: strconv.FormatBool(DefaultStoreWatchEnabled), Usage: "Evaluator config: Whether enable watch store changes."}
	params = append(params, &k.StoreWatchEnabled)

	// Log configurations
	k.LogConf.LogLevel = StrParamDetail{Name: "log-level", Usage: "Log config: log level, available levels are panic, fatal, error, warn, info and debug."}
	params = append(params, &k.LogConf.LogLevel)
	k.LogConf.LogFormatter = StrParamDetail{Name: "log-formatter", Usage: "Log config: log formatter, available values are text and json."}
	params = append(params, &k.LogConf.LogFormatter)
	k.LogConf.LogReportCaller = StrParamDetail{Name: "log-reportcaller", DefaultValue: strconv.FormatBool(false), Usage: "Log config: if the caller(file, line and function) is included in the log entry."}
	params = append(params, &k.LogConf.LogReportCaller)
	k.LogConf.LogFileName = StrParamDetail{Name: "log-filename", Usage: "Log config: log file name."}
	params = append(params, &k.LogConf.LogFileName)
	k.LogConf.LogMaxSize = StrParamDetail{Name: "log-maxsize", Usage: "Log config: maximum size in megabytes of the log file before it gets rotated."}
	params = append(params, &k.LogConf.LogMaxSize)
	k.LogConf.LogCompress = StrParamDetail{Name: "log-compress", Usage: "Log config: if the rotated log files should be compressed."}
	params = append(params, &k.LogConf.LogCompress)
	k.LogConf.LogMaxBackups = StrParamDetail{Name: "log-maxbackups", Usage: "Log config: maximum number of old log files to retain."}
	params = append(params, &k.LogConf.LogMaxBackups)
	k.LogConf.LogMaxAge = StrParamDetail{Name: "log-maxage", Usage: "Log config: maximum number of days to retain old log files."}
	params = append(params, &k.LogConf.LogMaxAge)
	k.LogConf.LogLocalTime = StrParamDetail{Name: "log-localtime", Usage: "Log config: if local time is used for formatting the timestamps in backup files."}
	params = append(params, &k.LogConf.LogLocalTime)

	// Audit Log configurations
	k.AuditLogConf.LogLevel = StrParamDetail{Name: "auditlog-level", DefaultValue: DefaultAuditLogLevel, Usage: "Audit Log config: log level, available levels are panic, fatal, error, warn, info and debug."}
	params = append(params, &k.AuditLogConf.LogLevel)
	k.AuditLogConf.LogFormatter = StrParamDetail{Name: "auditlog-formatter", DefaultValue: DefaultAuditLogFormatter, Usage: "Audit Log config: log formatter, available values are text and json."}
	params = append(params, &k.AuditLogConf.LogFormatter)
	k.AuditLogConf.LogReportCaller = StrParamDetail{Name: "auditlog-reportcaller", DefaultValue: strconv.FormatBool(false), Usage: "Audit Log config: if the caller(file, line and function) is included in the log entry."}
	params = append(params, &k.AuditLogConf.LogReportCaller)
	k.AuditLogConf.LogFileName = StrParamDetail{Name: "auditlog-filename", DefaultValue: DefaultAuditLogFilename, Usage: "Audit Log config: log file name."}
	params = append(params, &k.AuditLogConf.LogFileName)
	k.AuditLogConf.LogMaxSize = StrParamDetail{Name: "auditlog-maxsize", DefaultValue: DefaultAuditLogMaxSize, Usage: "Audit Log config: maximum size in megabytes of the log file before it gets rotated."}
	params = append(params, &k.AuditLogConf.LogMaxSize)
	k.AuditLogConf.LogCompress = StrParamDetail{Name: "auditlog-compress", DefaultValue: "false", Usage: "Audit Log config: if the rotated log files should be compressed."}
	params = append(params, &k.AuditLogConf.LogCompress)
	k.AuditLogConf.LogMaxBackups = StrParamDetail{Name: "auditlog-maxbackups", DefaultValue: DefaultAuditLogMaxBackups, Usage: "Audit Log config: maximum number of old log files to retain."}
	params = append(params, &k.AuditLogConf.LogMaxBackups)
	k.AuditLogConf.LogMaxAge = StrParamDetail{Name: "auditlog-maxage", DefaultValue: DefaultAuditLogMaxAge, Usage: "Audit Log config: maximum number of days to retain old log files."}
	params = append(params, &k.AuditLogConf.LogMaxAge)
	k.AuditLogConf.LogLocalTime = StrParamDetail{Name: "auditlog-localtime", DefaultValue: "false", Usage: "Audit Log config: if local time is used for formatting the timestamps in backup files."}
	params = append(params, &k.AuditLogConf.LogLocalTime)

	k.AsserterConf.AsserterEndpoint = StrParamDetail{Name: "asserter-endpoint", Usage: "Assertion server endpoint."}
	params = append(params, &k.AsserterConf.AsserterEndpoint)
	k.AsserterConf.AsserterClientKeyPath = StrParamDetail{Name: "asserter-client-key", Usage: "Assertion service client key file."}
	params = append(params, &k.AsserterConf.AsserterClientKeyPath)
	k.AsserterConf.AsserterClientCertPath = StrParamDetail{Name: "asserter-client-cert", Usage: "Assertion service client cert file."}
	params = append(params, &k.AsserterConf.AsserterClientCertPath)
	k.AsserterConf.AsserterCaPath = StrParamDetail{Name: "asserter-ca-cert", Usage: "Assertion service CA cert file."}
	params = append(params, &k.AsserterConf.AsserterCaPath)
	k.AsserterConf.AsserterClientTimeout = StrParamDetail{Name: "asserter-client-timeout", DefaultValue: DefaultAsserterClientTimeout, Usage: "Assertion service client http timeout value."}
	params = append(params, &k.AsserterConf.AsserterClientTimeout)

	pflag.BoolVarP(&k.Version, "version", "", false, "print version information")

	for _, paramDetail := range params {
		pflag.StringVarP(&(paramDetail.Value), paramDetail.Name, paramDetail.ShortName, paramDetail.DefaultValue, paramDetail.Usage)
	}
	pflag.Parse()

	if k.Version {
		printVersionInfoFun()
		os.Exit(0)
	}

	if len(k.ConfigFile.Value) == 0 {
		envVarName := FlagToEnv(k.ConfigFile.Name)
		val := os.Getenv(envVarName)
		if len(val) != 0 {
			k.ConfigFile.Value = val
		}
	}

	var conf *cfg.Config
	if k.ConfigFile.Value != "" {
		var err error
		conf, err = cfg.ReadConfig(k.ConfigFile.Value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fail to parse config file %s, error is %v. \n", k.ConfigFile.Value, err)
			k.usage()
		}
	} else {
		conf = nil
	}

	pflag.VisitAll(func(f *pflag.Flag) {
		key := FlagToEnv(f.Name)
		if !f.Changed {
			//if not set from command line, search it from environment variable
			val := os.Getenv(key)
			if val != "" {
				f.Value.Set(val)
			} else {
				//if not set from environment variable, search it from config file
				switch f.Name {
				case k.Endpoint.Name:
					if conf != nil && conf.ServerConfig != nil && len(conf.ServerConfig.Endpoint) != 0 {
						f.Value.Set(conf.ServerConfig.Endpoint)
					}
				case k.Insecure.Name:
					if conf != nil && conf.ServerConfig != nil && len(conf.ServerConfig.Insecure) != 0 {
						f.Value.Set(conf.ServerConfig.Insecure)
					}
				case k.EnableAuthz.Name:
					if conf != nil && conf.ServerConfig != nil && len(conf.ServerConfig.EnableAuthz) != 0 {
						f.Value.Set(conf.ServerConfig.EnableAuthz)
					}
				case k.KeyPath.Name:
					if conf != nil && conf.ServerConfig != nil && len(conf.ServerConfig.KeyPath) != 0 {
						f.Value.Set(conf.ServerConfig.KeyPath)
					}
				case k.CertPath.Name:
					if conf != nil && conf.ServerConfig != nil && len(conf.ServerConfig.CertPath) != 0 {
						f.Value.Set(conf.ServerConfig.CertPath)
					}
				case k.ForceClientCert.Name:
					if conf != nil && conf.ServerConfig != nil {
						f.Value.Set(strconv.FormatBool(conf.ServerConfig.ForceClientCert))
					}
				case k.ClientCertPath.Name:
					if conf != nil && conf.ServerConfig != nil && len(conf.ServerConfig.ClientCertPath) != 0 {
						f.Value.Set(conf.ServerConfig.ClientCertPath)
					}
				case k.StoreType.Name:
					if conf != nil && conf.StoreConfig != nil && len(conf.StoreConfig.StoreType) != 0 {
						f.Value.Set(conf.StoreConfig.StoreType)
					}
				case k.StoreWatchEnabled.Name:
					if conf != nil {
						f.Value.Set(strconv.FormatBool(conf.EnableWatch))
					}
				// Log configurations
				case k.LogConf.LogLevel.Name:
					if conf != nil && conf.LogConfig != nil {
						f.Value.Set(conf.LogConfig.Level)
					}
				case k.LogConf.LogFormatter.Name:
					if conf != nil && conf.LogConfig != nil {
						f.Value.Set(conf.LogConfig.Formatter)
					}
				case k.LogConf.LogReportCaller.Name:
					if conf != nil && conf.LogConfig != nil {
						f.Value.Set(strconv.FormatBool(conf.LogConfig.SetReportCaller))
					}
				case k.LogConf.LogFileName.Name:
					if conf != nil && conf.LogConfig != nil && conf.LogConfig.RotationConfig != nil {
						f.Value.Set(conf.LogConfig.RotationConfig.Filename)
					}
				case k.LogConf.LogMaxSize.Name:
					if conf != nil && conf.LogConfig != nil && conf.LogConfig.RotationConfig != nil {
						f.Value.Set(strconv.Itoa(conf.LogConfig.RotationConfig.MaxSize))
					}
				case k.LogConf.LogMaxAge.Name:
					if conf != nil && conf.LogConfig != nil && conf.LogConfig.RotationConfig != nil {
						f.Value.Set(strconv.Itoa(conf.LogConfig.RotationConfig.MaxAge))
					}
				case k.LogConf.LogMaxBackups.Name:
					if conf != nil && conf.LogConfig != nil && conf.LogConfig.RotationConfig != nil {
						f.Value.Set(strconv.Itoa(conf.LogConfig.RotationConfig.MaxBackups))
					}
				case k.LogConf.LogCompress.Name:
					if conf != nil && conf.LogConfig != nil && conf.LogConfig.RotationConfig != nil {
						f.Value.Set(strconv.FormatBool(conf.LogConfig.RotationConfig.Compress))
					}
				case k.LogConf.LogLocalTime.Name:
					if conf != nil && conf.LogConfig != nil && conf.LogConfig.RotationConfig != nil {
						f.Value.Set(strconv.FormatBool(conf.LogConfig.RotationConfig.LocalTime))
					}
				// Audit Log configurations
				case k.AuditLogConf.LogLevel.Name:
					if conf != nil && conf.AuditLogConfig != nil {
						f.Value.Set(conf.AuditLogConfig.Level)
					}
				case k.AuditLogConf.LogFormatter.Name:
					if conf != nil && conf.AuditLogConfig != nil {
						f.Value.Set(conf.AuditLogConfig.Formatter)
					}
				case k.AuditLogConf.LogReportCaller.Name:
					if conf != nil && conf.AuditLogConfig != nil {
						f.Value.Set(strconv.FormatBool(conf.AuditLogConfig.SetReportCaller))
					}
				case k.AuditLogConf.LogFileName.Name:
					if conf != nil && conf.AuditLogConfig != nil && conf.AuditLogConfig.RotationConfig != nil {
						f.Value.Set(conf.AuditLogConfig.RotationConfig.Filename)
					}
				case k.AuditLogConf.LogMaxSize.Name:
					if conf != nil && conf.AuditLogConfig != nil && conf.AuditLogConfig.RotationConfig != nil {
						f.Value.Set(strconv.Itoa(conf.AuditLogConfig.RotationConfig.MaxSize))
					}
				case k.AuditLogConf.LogMaxAge.Name:
					if conf != nil && conf.AuditLogConfig != nil && conf.AuditLogConfig.RotationConfig != nil {
						f.Value.Set(strconv.Itoa(conf.AuditLogConfig.RotationConfig.MaxAge))
					}
				case k.AuditLogConf.LogMaxBackups.Name:
					if conf != nil && conf.AuditLogConfig != nil && conf.AuditLogConfig.RotationConfig != nil {
						f.Value.Set(strconv.Itoa(conf.AuditLogConfig.RotationConfig.MaxBackups))
					}
				case k.AuditLogConf.LogCompress.Name:
					if conf != nil && conf.AuditLogConfig != nil && conf.AuditLogConfig.RotationConfig != nil {
						f.Value.Set(strconv.FormatBool(conf.AuditLogConfig.RotationConfig.Compress))
					}
				case k.AuditLogConf.LogLocalTime.Name:
					if conf != nil && conf.AuditLogConfig != nil && conf.AuditLogConfig.RotationConfig != nil {
						f.Value.Set(strconv.FormatBool(conf.AuditLogConfig.RotationConfig.LocalTime))
					}
					// Asserter webhook configurations
				case k.AsserterConf.AsserterEndpoint.Name:
					if conf != nil && conf.AsserterWebhookConfig != nil {
						f.Value.Set(conf.AsserterWebhookConfig.Endpoint)
					}
				case k.AsserterConf.AsserterCaPath.Name:
					if conf != nil && conf.AsserterWebhookConfig != nil {
						f.Value.Set(conf.AsserterWebhookConfig.CACert)
					}
				case k.AsserterConf.AsserterClientCertPath.Name:
					if conf != nil && conf.AsserterWebhookConfig != nil {
						f.Value.Set(conf.AsserterWebhookConfig.ClientCert)
					}
				case k.AsserterConf.AsserterClientKeyPath.Name:
					if conf != nil && conf.AsserterWebhookConfig != nil {
						f.Value.Set(conf.AsserterWebhookConfig.ClientKey)
					}
				case k.AsserterConf.AsserterClientTimeout.Name:
					if conf != nil && conf.AsserterWebhookConfig != nil {
						f.Value.Set(string(conf.AsserterWebhookConfig.HTTPTimeout))
					}
				default:
					//
				}

				key, ok := storeParamsMap[f.Name]
				if ok {
					if conf != nil && conf.StoreConfig != nil && conf.StoreConfig.StoreProps != nil {
						if value, ok := conf.StoreConfig.StoreProps[key]; ok {
							switch x := value.(type) {
							case bool:
								f.Value.Set(strconv.FormatBool(value.(bool)))
							case int:
								f.Value.Set(strconv.Itoa(value.(int)))
							case string:
								f.Value.Set(value.(string))
							default:
								fmt.Printf("Unsupported type: %T\n", x)
							}
						}
					}
				}
			}

		}
	})

	fmt.Printf("parameters:%v\n", k)
}

// FlagToEnv converts flag string to upper-case environment variable key string.
func FlagToEnv(name string) string {
	return EnvVarPrefix + "_" + strings.ToUpper(strings.Replace(name, "-", "_", -1))
}

func (k *Parameters) ValidateFlags() {
	insecure, err := strconv.ParseBool(k.Insecure.Value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid value for 'insecure' parameter: %s", k.Insecure.Value)
		k.usage()
	}

	if len(k.EnableAuthz.Value) != 0 {
		_, err = strconv.ParseBool(k.EnableAuthz.Value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid value for 'enableAuthz' parameter: %s", k.EnableAuthz.Value)
			k.usage()
		}
	}

	if !insecure {
		if k.CertPath.Value == "" || k.KeyPath.Value == "" {
			fmt.Fprintln(os.Stderr, "In secure mode, "+k.KeyPath.Name+", "+k.CertPath.Name+" should be passed.")
			k.usage()
		}

		if k.ForceClientCert.Value != "" {
			forceClientCert, err := strconv.ParseBool(k.ForceClientCert.Value)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid value for 'ForceClientCert' parameter: %s", k.ForceClientCert.Value)
				k.usage()
			}
			if forceClientCert && k.ClientCertPath.Value == "" {
				fmt.Fprintln(os.Stderr, "In secure mode and force client certification is enabled, "+k.ClientCertPath.Name+" should be passed.")
				k.usage()
			}
		}
	}
}

func (k *Parameters) Param2Config(storeParamsMap map[string]string) (*cfg.Config, error) {

	conf := cfg.Config{}

	var storeConf cfg.StoreConfig
	storeConf.StoreType = k.StoreType.Value
	storeConf.StoreProps = make(map[string]interface{})

	for k, v := range storeParamsMap {
		f := pflag.Lookup(k)
		if f != nil {
			storeConf.StoreProps[v] = f.Value.String()
		}
	}

	conf.StoreConfig = &storeConf

	watchEnabled, _ := strconv.ParseBool(k.StoreWatchEnabled.Value)
	conf.EnableWatch = watchEnabled

	// Log Configuration
	if len(k.LogConf.LogLevel.Value) != 0 ||
		len(k.LogConf.LogFormatter.Value) != 0 ||
		len(k.LogConf.LogReportCaller.Value) != 0 ||
		len(k.LogConf.LogFileName.Value) != 0 ||
		len(k.LogConf.LogCompress.Value) != 0 ||
		len(k.LogConf.LogMaxSize.Value) != 0 ||
		len(k.LogConf.LogMaxAge.Value) != 0 ||
		len(k.LogConf.LogMaxBackups.Value) != 0 ||
		len(k.LogConf.LogLocalTime.Value) != 0 {
		var logConf logging.LogConfig
		if len(k.LogConf.LogLevel.Value) != 0 {
			logConf.Level = k.LogConf.LogLevel.Value
		}
		if len(k.LogConf.LogFormatter.Value) != 0 {
			logConf.Formatter = k.LogConf.LogFormatter.Value
		}
		if len(k.LogConf.LogReportCaller.Value) != 0 {
			value, _ := strconv.ParseBool(k.LogConf.LogReportCaller.Value)
			logConf.SetReportCaller = value
		}
		if len(k.LogConf.LogFileName.Value) != 0 {
			var rotateConfig lumberjack.Logger
			if len(k.LogConf.LogFileName.Value) != 0 {
				rotateConfig.Filename = k.LogConf.LogFileName.Value
			}
			if len(k.LogConf.LogMaxSize.Value) != 0 {
				value, _ := strconv.Atoi(k.LogConf.LogMaxSize.Value)
				rotateConfig.MaxSize = value
			}
			if len(k.LogConf.LogMaxAge.Value) != 0 {
				value, _ := strconv.Atoi(k.LogConf.LogMaxAge.Value)
				rotateConfig.MaxAge = value
			}
			if len(k.LogConf.LogMaxBackups.Value) != 0 {
				value, _ := strconv.Atoi(k.LogConf.LogMaxBackups.Value)
				rotateConfig.MaxBackups = value
			}
			if len(k.LogConf.LogCompress.Value) != 0 {
				value, _ := strconv.ParseBool(k.LogConf.LogCompress.Value)
				rotateConfig.Compress = value
			}
			if len(k.LogConf.LogLocalTime.Value) != 0 {
				value, _ := strconv.ParseBool(k.LogConf.LogLocalTime.Value)
				rotateConfig.Compress = value
			}
			logConf.RotationConfig = &rotateConfig
		}
		conf.LogConfig = &logConf
	}

	// Audit Log Configuration
	if len(k.AuditLogConf.LogLevel.Value) != 0 ||
		len(k.AuditLogConf.LogFormatter.Value) != 0 ||
		len(k.AuditLogConf.LogReportCaller.Value) != 0 ||
		len(k.AuditLogConf.LogFileName.Value) != 0 ||
		len(k.AuditLogConf.LogCompress.Value) != 0 ||
		len(k.AuditLogConf.LogMaxSize.Value) != 0 ||
		len(k.AuditLogConf.LogMaxAge.Value) != 0 ||
		len(k.AuditLogConf.LogMaxBackups.Value) != 0 ||
		len(k.AuditLogConf.LogLocalTime.Value) != 0 {
		var auditLogConf logging.LogConfig
		if len(k.AuditLogConf.LogLevel.Value) != 0 {
			auditLogConf.Level = k.AuditLogConf.LogLevel.Value
		}
		if len(k.AuditLogConf.LogFormatter.Value) != 0 {
			auditLogConf.Formatter = k.AuditLogConf.LogFormatter.Value
		}
		if len(k.AuditLogConf.LogReportCaller.Value) != 0 {
			value, _ := strconv.ParseBool(k.AuditLogConf.LogReportCaller.Value)
			auditLogConf.SetReportCaller = value
		}
		if len(k.AuditLogConf.LogFileName.Value) != 0 {
			var rotateConfig lumberjack.Logger
			if len(k.AuditLogConf.LogFileName.Value) != 0 {
				rotateConfig.Filename = k.AuditLogConf.LogFileName.Value
			}
			if len(k.AuditLogConf.LogMaxSize.Value) != 0 {
				value, _ := strconv.Atoi(k.AuditLogConf.LogMaxSize.Value)
				rotateConfig.MaxSize = value
			}
			if len(k.AuditLogConf.LogMaxAge.Value) != 0 {
				value, _ := strconv.Atoi(k.AuditLogConf.LogMaxAge.Value)
				rotateConfig.MaxAge = value
			}
			if len(k.AuditLogConf.LogMaxBackups.Value) != 0 {
				value, _ := strconv.Atoi(k.AuditLogConf.LogMaxBackups.Value)
				rotateConfig.MaxBackups = value
			}
			if len(k.AuditLogConf.LogCompress.Value) != 0 {
				value, _ := strconv.ParseBool(k.AuditLogConf.LogCompress.Value)
				rotateConfig.Compress = value
			}
			if len(k.AuditLogConf.LogLocalTime.Value) != 0 {
				value, _ := strconv.ParseBool(k.AuditLogConf.LogLocalTime.Value)
				rotateConfig.Compress = value
			}
			auditLogConf.RotationConfig = &rotateConfig
		}
		conf.AuditLogConfig = &auditLogConf
	}

	// Asserter webhook Configuration
	if len(k.AsserterConf.AsserterEndpoint.Value) != 0 {
		asserterConf := assertion.AsserterConfig{}
		asserterConf.Endpoint = k.AsserterConf.AsserterEndpoint.Value

		if len(k.AsserterConf.AsserterCaPath.Value) != 0 {
			asserterConf.CACert = k.AsserterConf.AsserterCaPath.Value
		}
		if len(k.AsserterConf.AsserterClientCertPath.Value) != 0 {
			asserterConf.ClientCert = k.AsserterConf.AsserterClientCertPath.Value
		}
		if len(k.AsserterConf.AsserterClientKeyPath.Value) != 0 {
			asserterConf.ClientKey = k.AsserterConf.AsserterClientKeyPath.Value
		}
		if len(k.AsserterConf.AsserterClientTimeout.Value) != 0 {
			timeout, err := strconv.Atoi(k.AsserterConf.AsserterClientTimeout.Value)
			if err == nil {
				asserterConf.HTTPTimeout = timeout
			}
		}
		conf.AsserterWebhookConfig = &asserterConf
	}

	fmt.Printf("%v\n", conf.AsserterWebhookConfig)

	return &conf, nil
}

func (k *Parameters) usage() {
	pflag.Usage()
	os.Exit(1)
}
