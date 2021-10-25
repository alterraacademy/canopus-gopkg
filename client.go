package canopus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/alterraacademy/canopus-gopkg/config"
	"github.com/alterraacademy/canopus-gopkg/helper"

	"github.com/go-playground/validator/v10"
)

type ClientMethod interface {
	GetAvailableMethod(amount float64) ([]PaymentMethod, error)
	GenerateCart(payload CartPayload, paymentMethod PaymentMethod) (CartResponse, error)
}

type ConfigOptions struct {
	MerchantKey []byte `validate:"required"`
	MerchantPem []byte `validate:"required"`
	MerchantID  string `validate:"required"`
	Secret      string `validate:"required"`
}

type Client struct {
	Canopus    *ConfigOptions
	HTTPClient config.HTTPClient
	Token      string
	Validator  *validator.Validate
}

func NewAPICLient(cano *ConfigOptions) ClientMethod {
	validate := validator.New()

	err := validate.Struct(cano)
	if err != nil {
		return nil
	}

	client := &Client{
		Canopus:    cano,
		HTTPClient: config.GetHTTPClient(),
		Validator:  validate,
	}
	client.getToken()

	return client
}

func (c *Client) callClient(method, url string, bodyJSON []byte) (CommonMessage, error) {

	signature, err := helper.GenerateSignature(c.Canopus.MerchantKey, bodyJSON)
	if err != nil {
		config.Logger(config.Error, "call client - generate signature", err.Error())
		return CommonMessage{}, err
	}

	if c.Token == "create token" {
		c.Token = ""
	}

	fullURL := fmt.Sprintf(url, config.Environment.BaseUrl(), c.Canopus.MerchantID)
	result, err := c.HTTPClient.RequestDo(method, fullURL, bodyJSON, c.Token, signature)
	if err != nil {
		if strings.Contains(err.Error(), "token") {
			c.getToken()
		}
		config.Logger(config.Error, "call client - request http", err)
		return CommonMessage{}, err
	}

	respResult, err := ValidateResponse(result)
	if err != nil {
		config.Logger(config.Error, "call client - validate response", err)
		return CommonMessage{}, err
	}

	// never pass the verify.
	// _, err = helper.VerifySignature(c.Canopus.MerchantPem, result, respResult.Signature)
	// if err != nil {
	// 	config.Logger(config.Error, "call client - verify signature", err)
	// 	return CommonMessage{}, err
	// }

	return respResult, nil

}

func (c *Client) getToken() {

	c.Token = "create token"

	body := make(map[string]string)
	body["secret"] = c.Canopus.Secret

	bodyJson, _ := json.Marshal(body)
	respResult, err := c.callClient(http.MethodPost, PostTokenEndpoint, bodyJson)
	if err != nil {
		config.Logger(config.Error, "create token - call client", err)
	}

	token := respResult.Response.Data["token"].(string)
	if token == "invalid" || len(token) <= 0 {
		config.Logger(config.Error, "create token - invalid token", err)
	}

	// set new token to struct
	c.Token = token
}

// GetAvailableMethod Get all available method by merhantID
func (c *Client) GetAvailableMethod(amount float64) ([]PaymentMethod, error) {

	body := make(map[string]interface{})
	body["amount"] = amount

	bodyJson, _ := json.Marshal(body)
	respResult, err := c.callClient(http.MethodGet, GetPaymentMethod, bodyJson)
	if err != nil {
		config.Logger(config.Error, "available method - call client", err)
	}

	var result []PaymentMethod
	for _, payment := range respResult.Response.Data["method"].([]interface{}) {
		var tmp PaymentMethod
		paymentJson, _ := json.Marshal(payment)
		err = json.Unmarshal(paymentJson, &tmp)
		if err != nil {
			return []PaymentMethod{}, err
		}
		result = append(result, tmp)
	}

	return result, nil
}

func (c *Client) GenerateCart(payload CartPayload, paymentMethod PaymentMethod) (CartResponse, error) {

	err := c.Validator.Struct(payload)
	if err != nil {
		config.Logger(config.Error, "generate cart - validate payload", err)
		return CartResponse{}, err
	}

	err = c.Validator.Struct(paymentMethod)
	if err != nil {
		config.Logger(config.Error, "generate cart - validate payment method", err)
		return CartResponse{}, err
	}

	var url string = PostAPICartGenerate
	if config.DefaultAPIType == config.SNAP {
		url = PostSnapCartGenerate
	}

	bodyJson, _ := json.Marshal(payload)
	respResult, err := c.callClient(http.MethodPost, url, bodyJson)
	if err != nil {
		config.Logger(config.Error, "generate cart - call client", err)
	}

	var result CartResponse
	if config.DefaultAPIType == config.SNAP {
		result.CartID = respResult.Response.Data["cartID"].(string)
		result.PayTo = respResult.Response.Data["checkoutURL"].(string)
		return result, nil
	}

	switch paymentMethod.Type {
	case "op":
		result.CartID = respResult.Response.Data["cartID"].(string)
		result.PayTo = respResult.Response.Data["checkoutURL"].(string)
	case "va":
		result.Bank = respResult.Response.Data["payment"].(string)
		result.Number = respResult.Response.Data["number"].(string)
		result.Amount = respResult.Response.Data["amount"].(string)
	case "bt":
		result.Bank = respResult.Response.Data["bank"].(string)
		result.Number = respResult.Response.Data["account"].(string)
		result.Amount = respResult.Response.Data["amountTotal"].(string)
	case "cc":
		result.CartID = respResult.Response.Data["cartID"].(string)
		result.PayTo = respResult.Response.Data["checkoutURL"].(string)
	}

	return result, nil
}
