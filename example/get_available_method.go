package examples

import (
	"fmt"
	"io/ioutil"

	"github.com/alterraacademy/canopus-gopkg"
)

func GetAvailableMethod() {
	privKey, err := ioutil.ReadFile("/your/path/to/M-00001.key")
	if err != nil {
		fmt.Println(err)
	}
	privPem, err := ioutil.ReadFile("/your/path/to/M-00001.pem")
	if err != nil {
		fmt.Println(err)
	}

	canopusClient := canopus.NewAPICLient(&canopus.ConfigOptions{
		MerchantKey: privKey,
		MerchantPem: privPem,
		MerchantID:  "M-0001",
		Secret:      "yoursecret",
	})

	res, err := canopusClient.GetAvailableMethod(10)
	fmt.Printf("API Result : %+v\nError : %+v", res, err)
}
