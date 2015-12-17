package veritrans_test

import (
	vt "github.com/apseria/veritrans-golang"

	"testing"
	//"net/http/httptest"
	//"fmt"
	//"io/ioutil"
	"encoding/json"
	"net/http"
	//"bytes"
)

var (
	cert       = "./cacert.pem"
	production = false
	serverKey  = ""
	v          = vt.New(serverKey, production, cert)
)

func TestCtor(t *testing.T) {
	if v == nil {
		t.Fatalf("Expected non nil, happened: %v", v)
	}
}

var urlTests = []struct {
	in  bool
	out string
	exp string
}{
	{
		in:  false,
		exp: vt.SANDBOX_BASE_URL,
	},
	{
		in:  true,
		exp: vt.PRODUCTION_BASE_URL,
	},
}

func TestGetBaseUrl(t *testing.T) {
	for i, test := range urlTests {
		v := vt.New(serverKey, test.in, cert)

		var err error
		test.out, err = v.GetBaseUrl()
		if err != nil {
			t.Fatalf("#%d Expected %v got %v", i, test.exp, test.out)
		}
	}
}

var transaction = make(map[string]interface{})

var transaction_details = map[string]string{
	"order_id":     "0900090976",
	"gross_amount": "10000",
}

func TestVtWebCharge(t *testing.T) {
	transaction["transaction_details"] = transaction_details
	transaction["payment_type"] = "vtweb"

	resp, err := v.VtWebCharge(transaction)
	if err != nil {
		t.Fatalf("Expected nil got: %v", err)
	}

	var body vt.VTWebChargeResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Expected nil got: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected response status %v, got %v", http.StatusOK, resp.StatusCode)
	}

	if body.StatusCode != "201" {
		t.Fatalf("Expected response body status 201, got %v", body.StatusCode)
	}

	if body.RedirectUrl == "" {
		t.Fatalf("Expected redirect url, got %v", body.RedirectUrl)
	}
}

// todo: write more tests

