package canopusgo

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/alterraacademy/canopus-gopkg/config"
)

// Result extract Canopus Response (json key result)
type Result struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// CommonMessage canopus request/response
type CommonMessage struct {
	Response struct {
		Result Result                 `json:"result"`
		Data   map[string]interface{} `json:"data"`
	} `json:"response,omitempty"`
	Request struct {
		Result Result                 `json:"result"`
		Data   map[string]interface{} `json:"data"`
	} `json:"request,omitempty"`
	Signature string `json:"signature"`
}

func ValidateResponse(resp []byte) (CommonMessage, error) {
	var response CommonMessage

	err := json.Unmarshal(resp, &response)
	if err != nil {
		return CommonMessage{}, err
	}

	if response.Response.Result.Status == "ERROR" {
		switch response.Response.Result.Code {
		case "00000002":
			err := errors.New(config.Err002 + response.Response.Result.Message)
			return CommonMessage{}, err
		case "00000004":
			err := errors.New(config.Err004 + response.Response.Result.Message)
			return CommonMessage{}, err
		case "00000007":
			err := errors.New(config.Err007 + response.Response.Result.Message)
			return CommonMessage{}, err
		case "00000008":
			err := errors.New(config.Err008 + response.Response.Result.Message)
			return CommonMessage{}, err
		case "00000019":
			err := errors.New(config.Err019 + response.Response.Result.Message)
			return CommonMessage{}, err
		case "00000020":
			err := errors.New(config.Err020 + response.Response.Result.Message)
			return CommonMessage{}, err
		default:
			err := errors.New(config.ErrUnknown + response.Response.Result.Message)
			return CommonMessage{}, err
		}
	}

	return response, nil
}

// Callback canopus body when payment settlement, expired and canopus notification
type Callback struct {
	Request struct {
		Data struct {
			Amount          int    `json:"amount"`
			Bank            string `json:"bank"`
			MerchantID      string `json:"merchantID"`
			MerchantOrderID string `json:"merchantOrderId"`
			Number          string `json:"number"`
			OrderID         string `json:"orderID"`
			Status          string `json:"status"`
			Time            struct {
				Created time.Time `json:"created"`
				Updated time.Time `json:"updated"`
			} `json:"time"`
			TransactionID string `json:"transactionID"`
			Type          string `json:"type"`
		} `json:"data"`
		Result Result `json:"result"`
	} `json:"request"`
	Signature string `json:"signature"`
}

// PaymentMethod payment method available
type PaymentMethod struct {
	Key         string      `json:"key" validate:"required"`
	Name        string      `json:"name"`
	Type        string      `json:"type" validate:"required"`
	Logo        string      `json:"logo"`
	Instruction interface{} `json:"instruction"`
}

// CartResponse response after generate cart
type CartResponse struct {
	CartID string `json:"cartID"`
	PayTo  string `json:"payto"`
	Amount string `json:"amount"`
	Bank   string `json:"bank"`
	Number string `json:"number"`
}
