package canopusgo

import (
	"errors"
	"time"
)

var (
	BaseURL         = "https://canopus-auth.sumpahpalapa.com"
	DefaultDuration = time.Second * time.Duration(10)
	DefaultType     = "api"
	Err002          = errors.New("required parameter is missing")
	Err004          = errors.New("parameter has illegal value")
	Err007          = errors.New("signature is invalid")
	Err008          = errors.New("payement method is not available")
)
