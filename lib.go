package canopusgo

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
)

func ValidateResponse(resp []byte) (CommonMessage, error) {
	var response CommonMessage

	err := json.Unmarshal(resp, &response)
	if err != nil {
		return CommonMessage{}, err
	}

	switch response.Response.Result.Code {
	case "00000002":
		err := errors.New(Err002.Error() + response.Response.Result.Message)
		return CommonMessage{}, err
	case "00000004":
		err := errors.New(Err004.Error() + response.Response.Result.Message)
		return CommonMessage{}, err
	case "00000007":
		err := errors.New(Err007.Error() + response.Response.Result.Message)
		return CommonMessage{}, err
	case "00000008":
		err := errors.New(Err008.Error() + response.Response.Result.Message)
		return CommonMessage{}, err
	case "00000020":
		err := errors.New(Err020.Error() + response.Response.Result.Message)
		return CommonMessage{}, err
	}

	return response, nil
}

func VerifySignature(merchantPem []byte, response []byte, signature string) (bool, error) {
	block, _ := pem.Decode(merchantPem)
	if block == nil {
		return false, errors.New("failed to decode merchantPem")
	}

	pubkey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	sDec, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}

	hashed := sha256.Sum256(response)

	err = rsa.VerifyPKCS1v15(pubkey, crypto.SHA256, hashed[:], sDec)
	if err != nil {
		return false, err
	}
	return true, nil
}
