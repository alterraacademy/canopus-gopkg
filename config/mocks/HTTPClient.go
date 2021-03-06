// Code generated by mockery 2.7.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// HTTPClient is an autogenerated mock type for the HTTPClient type
type HTTPClient struct {
	mock.Mock
}

// RequestDo provides a mock function with given fields: method, url, bodyJson, token, signature
func (_m *HTTPClient) RequestDo(method string, url string, bodyJson []byte, token string, signature string) ([]byte, error) {
	ret := _m.Called(method, url, bodyJson, token, signature)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, string, []byte, string, string) []byte); ok {
		r0 = rf(method, url, bodyJson, token, signature)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, []byte, string, string) error); ok {
		r1 = rf(method, url, bodyJson, token, signature)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
