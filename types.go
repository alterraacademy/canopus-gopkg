package canopusgo

import (
	"encoding/json"
	"time"
)

// RawMessage partly get json for validate request/response and signature
type RawMessage struct {
	Response  json.RawMessage `json:"response,omitempty"`
	Request   json.RawMessage `json:"request,omitempty"`
	Signature string          `json:"signature"`
}

// CommonMessage canopus request/response
type CommonMessage struct {
	Response  Message `json:"response,omitempty"`
	Request   Message `json:"request,omitempty"`
	Signature string  `json:"signature"`
}

// Message canopus subparameter
type Message struct {
	Result Result                 `json:"result"`
	Data   map[string]interface{} `json:"data"`
}

// Result extract Canopus Response (json key result)
type Result struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// NotifData extract data when notification callback
type NotifData struct {
	Amount float64 `json:"amount"`
	Bank   string  `json:"bank"`
	Number string  `json:"number"`
}

// PaymentMethod payment method available
type PaymentMethod struct {
	Key         string      `json:"key" validate:"required"`
	Name        string      `json:"name"`
	Type        string      `json:"type" validate:"required"`
	Logo        string      `json:"logo"`
	Instruction interface{} `json:"instruction"`
}

// CartPayload payload to create cart
type CartPayload struct {
	CartDetails struct {
		ID      string `json:"id" validate:"required"`
		Payment struct {
			Key  string `json:"key"`
			Type string `json:"type"`
		} `json:"payment"`
		Amount    float64 `json:"amount" validate:"required"`
		Title     string  `json:"title" validate:"required"`
		Currency  string  `json:"currency" validate:"required"`
		ExpiredAt string  `json:"expiredAt" validate:"required"`
	} `json:"cartDetails" validate:"required"`
	ItemDetails     []CartPayloadItemDetail `json:"itemDetails" validate:"required"`
	CustomerDetails struct {
		FirstName      string `json:"firstName" validate:"required"`
		LastName       string `json:"lastName"`
		Email          string `json:"email" validate:"required"`
		Phone          string `json:"phone" validate:"required"`
		BillingAddress struct {
			FirstName  string `json:"firstName"`
			LastName   string `json:"lastName"`
			Phone      string `json:"phone"`
			Address    string `json:"address"`
			City       string `json:"city"`
			PostalCode string `json:"postalCode"`
		} `json:"billingAddress"`
		ShippingAddress struct {
			FirstName  string `json:"firstName"`
			LastName   string `json:"lastName"`
			Phone      string `json:"phone"`
			Address    string `json:"address"`
			City       string `json:"city"`
			PostalCode string `json:"postalCode"`
		} `json:"shippingAddress"`
	} `json:"customerDetails" validate:"required"`
	Environment struct {
		Agent   string `json:"agent"`
		Mode    string `json:"mode"`
		Os      string `json:"os"`
		Version string `json:"version"`
	} `json:"environment"`
	URL struct {
		ReturnURL       string `json:"returnURL" validate:"required"`
		CancelURL       string `json:"cancelURL" validate:"required"`
		NotificationURL string `json:"notificationURL" validate:"required"`
	} `json:"url" validate:"required"`
	ExtendInfo struct {
		AdditionalPrefix string `json:"additionalPrefix"`
	} `json:"extendInfo"`
}

// CartPayloadItemDetail item cart detail
type CartPayloadItemDetail struct {
	Name           string  `json:"name" validate:"required"`
	Desc           string  `json:"desc"`
	Price          float64 `json:"price" validate:"required"`
	Quantity       int     `json:"quantity" validate:"required"`
	SKU            string  `json:"SKU" validate:"required"`
	AdditionalInfo struct {
		NoHandphone string `json:"No Handphone"`
	} `json:"additionalInfo"`
}

// CartResponse response after generate cart
type CartResponse struct {
	CartID string `json:"cartID"`
	PayTo  string `json:"payto"`
	Amount string `json:"amount"`
	Bank   string `json:"bank"`
	Number string `json:"number"`
}

// DataCallback data callback from canopus
type DataCallback struct {
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
}

// Callback canopus body when payment settlement, expired and canopus notification
type Callback struct {
	Request struct {
		Data   DataCallback `json:"data"`
		Result Result       `json:"result"`
	} `json:"request"`
	Signature string `json:"signature"`
}

// InitService to create canoput service
type InitService struct {
	Type        string // snap | api
	MerchantKey []byte `validate:"required"`
	MerchantPem []byte `validate:"required"`
	MerchantID  string `validate:"required"`
	Secret      string `validate:"required"`
	TimeOut     time.Duration
}
