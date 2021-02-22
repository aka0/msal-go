package msal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAcquireTokenForClient(t *testing.T) {

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	testApp, err := NewClientApplication("testTenantID", "testClientID", "testClientSecret", "https%3A%2F%2Fgraph.microsoft.com%2F.default", nil)
	if err != nil {
		t.Errorf(err.Error())
	}

	testApp.BaseURL = server.URL

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"token_type": "Bearer","expires_in": 3599,"access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik1uQ19WWmNBVGZNNXBP", "ext_expires_in": 1234}`)
	})

	got, err := testApp.AcquireTokenForClient()
	if err != nil {
		t.Errorf(err.Error())
	}

	want := &Token{AccessToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik1uQ19WWmNBVGZNNXBP", TokenType: "Bearer", ExpiresIn: 3599, ExtExpiresIn: 1234}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Want: %s, %s, %v, %v", want.TokenType, want.AccessToken, want.ExpiresIn, want.ExtExpiresIn)
		t.Errorf("Got: %s, %s, %v, %v", got.TokenType, got.AccessToken, got.ExpiresIn, got.ExtExpiresIn)
	}
}
