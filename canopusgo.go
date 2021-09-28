package canopusgo

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	validator "github.com/go-playground/validator/v10"
)

type Canopus struct {
	Type        string // snap | api
	MerchantKey []byte
	MerchantPem []byte
	MerchantID  string
	Secret      string
	Client      *http.Client
	Token       string
	Validator   *validator.Validate
}

// TODO: create validation in parameter for each function
func CreateService(init InitService) (*Canopus, error) {
	validate := validator.New()

	err := validate.Struct(init)
	if err != nil {
		return &Canopus{}, err
	}

	canType := DefaultType
	timeout := DefaultDuration

	if init.Type != "" {
		canType = init.Type
	}
	if init.TimeOut != 0 {
		timeout = init.TimeOut
	}

	client := &http.Client{Timeout: timeout}
	return &Canopus{
		Type:        canType,
		MerchantKey: init.MerchantKey,
		MerchantPem: init.MerchantPem,
		MerchantID:  init.MerchantID,
		Secret:      init.Secret,
		Client:      client,
		Validator:   validate,
	}, nil
}

// GenerateSignature ...
func (cano *Canopus) GenerateSignature(payload []byte) (string, error) {
	var signature string
	var err error

	block, _ := pem.Decode(cano.MerchantKey)
	if block == nil {
		err = errors.New("failed to decode key file")
		return signature, err
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return signature, err
	}

	hashed := sha256.Sum256(payload)

	s, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed[:])
	if err != nil {
		return signature, err
	}

	signature = base64.StdEncoding.EncodeToString(s)

	return signature, err
}

// GetToken ...
func (cano *Canopus) GetToken() (string, error) {
	url := fmt.Sprintf("%v/api/v1/merchants/%v/token", BaseURL, cano.MerchantID)
	body := make(map[string]string)
	body["secret"] = cano.Secret
	bodyJson, _ := json.Marshal(body)
	signature, err := cano.GenerateSignature(bodyJson)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))

	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Signature", signature)

	resp, err := cano.Client.Do(req)
	if err != nil {
		return "", err
	}
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errMsg := fmt.Sprintf("error http code: %v", resp.StatusCode)
		return "", errors.New(errMsg)
	}

	respResult, err := ValidateResponse(response)
	if err != nil {
		return "", err
	}

	token := respResult.Response.Data["token"].(string)
	if token == "invalid" || len(token) <= 0 {
		return "", errors.New("invalid token")
	}

	// set new token to struct
	cano.Token = token
	return token, nil
}

// GetAvailableMethod Get all available method by merhantID
func (cano *Canopus) GetAvailableMethod(amount float64) ([]PaymentMethod, error) {
	url := fmt.Sprintf("%v/api/v1/merchants/%v/method", BaseURL, cano.MerchantID)
	body := make(map[string]interface{})
	body["MerchantID"] = cano.MerchantID
	body["amount"] = amount
	bodyJson, _ := json.Marshal(body)
	signature, err := cano.GenerateSignature(bodyJson)
	if err != nil {
		return []PaymentMethod{}, err
	}

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(bodyJson))

	if err != nil {
		return []PaymentMethod{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Signature", signature)

	resp, err := cano.Client.Do(req)
	if err != nil {
		return []PaymentMethod{}, err
	}
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []PaymentMethod{}, err
	}

	defer resp.Body.Close()

	respResult, err := ValidateResponse(response)
	if err != nil {
		return []PaymentMethod{}, err
	}

	// _, err = VerifySignature(cano.MerchantPem, response, respResult.Signature)
	// if err != nil {
	// 	return []PaymentMethod{}, err
	// }

	var result []PaymentMethod
	for _, payment := range respResult.Response.Data["method"].([]interface{}) {
		var tmp PaymentMethod
		paymentJson, err := json.Marshal(payment)
		if err != nil {
			return []PaymentMethod{}, err
		}
		err = json.Unmarshal(paymentJson, &tmp)
		if err != nil {
			return []PaymentMethod{}, err
		}
		result = append(result, tmp)
	}

	return result, nil
}

func (cano *Canopus) GenerateCart(payload CartPayload, paymentMethod PaymentMethod) (CartResponse, error) {
	var url string
	var response []byte

	err := cano.Validator.Struct(payload)
	if err != nil {
		return CartResponse{}, err
	}

	err = cano.Validator.Struct(paymentMethod)
	if err != nil {
		return CartResponse{}, err
	}

	if cano.Type == "snap" {
		url = fmt.Sprintf("%v/api/v1/merchants/%v/%v", BaseURL, cano.MerchantID, "snap/cart")
	} else {
		url = fmt.Sprintf("%v/api/v1/merchants/%v/%v", BaseURL, cano.MerchantID, "cart")
	}

	bodyJson, _ := json.Marshal(payload)
	signature, err := cano.GenerateSignature(bodyJson)
	if err != nil {
		return CartResponse{}, err
	}

	resp, err := callCanopus(url, bodyJson, cano, signature)
	if err != nil {
		return CartResponse{}, err
	}

	// token expired
	if resp.StatusCode == 403 {
		_, err := cano.GetToken()

		if err != nil {
			return CartResponse{}, err
		}
		resp, err = callCanopus(url, bodyJson, cano, signature)
		if err != nil {
			return CartResponse{}, err
		}
	} else if resp.StatusCode != 200 {
		errMsg := fmt.Sprintf("error http code: %v", resp.StatusCode)
		return CartResponse{}, errors.New(errMsg)
	}
	defer resp.Body.Close()

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return CartResponse{}, err
	}
	respResult, err := ValidateResponse(response)
	if err != nil {
		return CartResponse{}, err
	}

	// TODO: cannot varifysignature, return false
	// _, err = VerifySignature(cano.MerchantPem, response, respResult.Signature)
	// if err != nil {
	// 	return CartResponse{}, err
	// }

	var result CartResponse
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

func callCanopus(url string, bodyJson []byte, cano *Canopus, signature string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))

	if err != nil {
		return nil, err
	}

	if cano.Token == "" {
		_, err := cano.GetToken()
		if err != nil {
			return nil, err
		}
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Signature", signature)
	req.Header.Add("authorization", "Bearer "+cano.Token)

	resp, err := cano.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
