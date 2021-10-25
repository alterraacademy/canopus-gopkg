# Canopus-GO

![GitHub](https://img.shields.io/github/license/alterraacademy/canopus-gopkg)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/alterraacademy/canopus-gopkg)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/alterraacademy/canopus-gopkg)
![GitHub all releases](https://img.shields.io/github/downloads/alterraacademy/canopus-gopkg/total)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/alterraacademy/canopus-gopkg?sort=semver)

## Overview

Canopus-GO is implementation of canopus service using Golang.

## Repository overview

Provide an overview of the directory structure and files, for example:
```bash
├── README.md
├── LICENSE
├── client.go
├── endpoint.go
├── request.go
├── response.go
├── go.mod
├── go.sum
├── config
│   ├── mocks
│   │   └── HTTPClient.go
│   ├── canopus.go
│   └── http.go
├── helper
│   └── signature.go
├── mocks
│   └── ClientMethod.go
└── example
    ├── create_cart.go
    └── get_available_method.go
```
## Usage

Create service first:
```go
key, err := ioutil.ReadFile("/your/path/to/M-00001.key")
if err != nil {
  fmt.Println(err)
}
pem, err := ioutil.ReadFile("/your/path/to/M-0001.pem")
if err != nil {
  fmt.Println(err)
}
canopusClient := canopus.NewAPICLient(&canopus.ConfigOptions{
		MerchantKey: key,
		MerchantPem: pem,
		MerchantID:  "M-0001",
		Secret:      "yoursupersecret",
	})
```

Then you can user this service to create payment

```go
GetAvailableMethod(amount float64) ([]PaymentMethod, error)
GenerateCart(payload CartPayload, paymentMethod PaymentMethod) (CartResponse, error)
```

Check some examples [canopusgo](https://github.com/alterraacademy/canopus-gopkg/tree/main/example)

### Additional Config

```go
// canopus.go
Environment        = Staging
DefaultAPIType     = SNAP
DefaultHttpTimeout = time.Second * time.Duration(10)
DefaultLogLevel    = 2
DefaultLogFormat   = "{\"log_type\":\"%s\",\"timestamp\":\"%s\",\"event\":\"%s\",\"detail\":\"%s\"}"

//example to change the default value
config.Environment = config.Production
```