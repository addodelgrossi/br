package brasilapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCEP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/cep/v1/01001000" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"cep":"01001000","state":"SP","city":"Sao Paulo","neighborhood":"Se","street":"Praca da Se","service":"mock"}`))
	}))
	defer server.Close()

	client := NewClient(server.Client())
	client.BaseURL = server.URL + "/api"

	got, err := client.CEP("01001-000")
	if err != nil {
		t.Fatalf("CEP() error = %v", err)
	}
	if got.City != "Sao Paulo" || got.State != "SP" {
		t.Fatalf("CEP() = %+v", got)
	}
}

func TestBankByCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/banks/v1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"ispb":"60701190","name":"ITAU","code":341,"fullName":"ITAU UNIBANCO S.A."}]`))
	}))
	defer server.Close()

	client := NewClient(server.Client())
	client.BaseURL = server.URL + "/api"

	got, err := client.BankByCode("341")
	if err != nil {
		t.Fatalf("BankByCode() error = %v", err)
	}
	if got.FullName != "ITAU UNIBANCO S.A." {
		t.Fatalf("BankByCode() = %+v", got)
	}
}
