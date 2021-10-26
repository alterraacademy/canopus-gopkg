package examples

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/alterraacademy/canopus-gopkg"
)

func CreateCart() {
	privKey, err := ioutil.ReadFile("/your/path/to/M-00001.key")
	if err != nil {
		fmt.Println(err)
	}
	privPem, err := ioutil.ReadFile("/your/path/to/M-0001.pem")
	if err != nil {
		fmt.Println(err)
	}

	canopusClient := canopus.NewAPICLient(&canopus.ConfigOptions{
		MerchantKey: privKey,
		MerchantPem: privPem,
		MerchantID:  "M-0001",
		Secret:      "yoursecret",
		Timeout:     20,
	})

	var payload canopus.CartPayload
	var paymentMethod canopus.PaymentMethod

	exp := time.Now().Local().Add(time.Hour * time.Duration(1)).Format(time.RFC3339) // 1 hour

	paymentMethod.Type = "op"
	paymentMethod.Key = "RFSP"

	arrItemDetail := []canopus.CartPayloadItemDetail{}
	itemDetail := canopus.CartPayloadItemDetail{}
	itemDetail.SKU = "PULSA20"
	itemDetail.AdditionalInfo.NoHandphone = "085721777738"
	itemDetail.Desc = "Pulsa 50 ribu ke 0811xxxxxx"
	itemDetail.Name = "Pulsa 50 ribu"
	itemDetail.Price = 50000
	itemDetail.Quantity = 1
	arrItemDetail = append(arrItemDetail, itemDetail)

	payload.ItemDetails = arrItemDetail

	payload.CustomerDetails.Email = "mega.octavia@gmail.com"
	payload.CustomerDetails.FirstName = "mega"
	payload.CustomerDetails.LastName = "octavia"
	payload.CustomerDetails.Phone = "08112342321"

	payload.CartDetails.ID = "TESTMOP-201911121424"
	payload.CartDetails.Title = "TEST CANO ORDER"
	payload.CartDetails.Currency = "IDR"
	payload.CartDetails.Amount = 50000
	payload.CartDetails.ExpiredAt = exp
	payload.CartDetails.Payment.Key = paymentMethod.Key
	payload.CartDetails.Payment.Type = paymentMethod.Type

	payload.URL.CancelURL = "https://yourdomain.com/cancel"
	payload.URL.ReturnURL = "https://yourdomain.com/return"
	payload.URL.NotificationURL = "https://yourdomain.com/notification"

	payTo, err := canopusClient.GenerateCart(payload, paymentMethod)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(payTo)
}
