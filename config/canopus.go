package config

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type EnvironmentType uint8
type APIType uint8
type LogType uint8

const (
	_ EnvironmentType = iota
	Staging
	Production

	_ APIType = iota
	API
	SNAP

	_ LogType = iota
	Error
	Info
	Debug

	Err002     = "required parameter is missing, "
	Err004     = "parameter has illegal value, "
	Err007     = "signature is invalid, "
	Err008     = "payement method is not available, "
	Err019     = "payment method error, "
	Err020     = "duplicate order id, "
	ErrUnknown = "unknown error, "
)

var envURL = map[EnvironmentType]string{
	Staging:    "https://canopus-auth.sumpahpalapa.com",
	Production: "https://canopus-auth.sepulsa.id",
}

// BaseUrl To get Midtrans Base URL
func (e EnvironmentType) BaseUrl() string {
	for k, v := range envURL {
		if k == e {
			return v
		}
	}
	return "undefined"
}

var canopusType = map[APIType]string{
	API:  "api",
	SNAP: "snap",
}

func (e APIType) ToString() string {
	for k, v := range canopusType {
		if k == e {
			return v
		}
	}
	return "undefined"
}

func (e LogType) CheckLevel() bool {
	return e <= LogType(DefaultLogLevel)
}

var logType = map[LogType]string{
	Error: "error",
	Info:  "info",
	Debug: "debug",
}

func (e LogType) ToString() string {
	for k, v := range logType {
		if k == e {
			return v
		}
	}
	return "undefined"
}

var (
	Environment        = Staging
	DefaultAPIType     = SNAP
	DefaultHttpTimeout = time.Second * time.Duration(10)
	DefaultLogLevel    = Info
	DefaultLogFormat   = "{\"timestamp\":\"%s\",\"log_type\":\"%s\",\"event\":\"%s\",\"detail\":\"%s\"}"

	DefaultHttpClient = &http.Client{Timeout: DefaultHttpTimeout}
	newLog            = log.New(os.Stdout, "", 0)
	Logger            = func(logtype LogType, event string, detail interface{}) {
		timestamp := time.Now().UTC().Format("2006-01-02 15:04:05")
		var jsonFormat []byte
		switch v := detail.(type) {
		case error:
			jsonFormat = []byte(v.Error())
		case string:
			jsonFormat = []byte(v)
		case map[string]interface{}:
			jsonFormat, _ = json.Marshal(detail)
		default:
			jsonFormat = []byte("Someting error.")
		}
		if logtype.CheckLevel() {
			newLog.Printf(DefaultLogFormat, timestamp, logtype.ToString(), event, jsonFormat)
		}
	}
)
