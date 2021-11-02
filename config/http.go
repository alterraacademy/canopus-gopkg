package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

type HTTPClient interface {
	RequestDo(method string, url string, bodyJson []byte, token string, signature string) ([]byte, error)
}

type HttpClientImp struct {
	HttpClient *http.Client
}

func GetHTTPClient() HTTPClient {
	return &HttpClientImp{
		HttpClient: DefaultHttpClient,
	}
}

func (c *HttpClientImp) RequestDo(method string, url string, bodyJson []byte, token string, signature string) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		Logger(Error, "create new http client request", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Signature", signature)
	req.Header.Add("Authorization", "Bearer "+token)
	if Environment == Staging {
		req.Header.Add("x-mock-ip", "127.0.0.1")
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		Logger(Error, "do client request", err)
		return nil, err
	}
	defer resp.Body.Close()

	dumpReq, _ := httputil.DumpRequest(req, true)
	dumpRes, _ := httputil.DumpResponse(resp, true)
	log := map[string]interface{}{
		"request":  string(dumpReq) + "\nPayload: " + string(bodyJson),
		"response": string(dumpRes),
	}
	Logger(Info, "do http request", log)

	// token expired
	if resp.StatusCode == 403 {
		return nil, errors.New("token expired or invalid")
	} else if resp.StatusCode != 200 {
		errMsg := fmt.Sprintf("error http code: %v", resp.StatusCode)
		return nil, errors.New(errMsg)
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger(Error, "read http client response", err)
		return nil, err
	}

	return response, nil
}
