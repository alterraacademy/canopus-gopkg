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

var logLevel = map[LogType]int{
	Error: 1,
	Info:  2,
	Debug: 3,
}

func (e LogType) CheckLevel() bool {
	for k, v := range logLevel {
		if k == e && v <= DefaultLogLevel {
			return true
		}
	}
	return false
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
	DefaultLogLevel    = 2
	DefaultLogFormat   = "{\"log_type\":\"%s\",\"timestamp\":\"%s\",\"event\":\"%s\",\"detail\":\"%s\"}"

	DefaultHttpClient = &http.Client{Timeout: DefaultHttpTimeout}
	newLog            = log.New(os.Stdout, "", 0)
	Logger            = func(logtype LogType, event string, detail interface{}) {
		timestamp := time.Now().UTC().Format("2006-01-02 15:04:05")
		var jsonFormat []byte
		if logtype == Error {
			jsonFormat = []byte(detail.(error).Error())
		} else {
			jsonFormat, _ = json.Marshal(detail)
		}

		if logtype.CheckLevel() {
			newLog.Printf(DefaultLogFormat, logtype.ToString(), timestamp, event, jsonFormat)
		}
	}
)
