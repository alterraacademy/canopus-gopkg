package helper

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func GenerateSignature(merchantKey []byte, payload []byte) (string, error) {
	var signature string
	var err error

	block, _ := pem.Decode(merchantKey)
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
