package etherscan_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/yourproject/etherscan" // Update this import path to match your project
)

func TestNewClient(t *testing.T) {
	apiKey := "test_api_key"
	client := etherscan.NewClient(apiKey)

	if client.GetAPIKey() != apiKey { // Assuming you have a getter method
		t.Errorf("Expected apiKey to be %s, got %s", apiKey, client.GetAPIKey())
	}

	if client.GetBaseURL() != "https://api.etherscan.io/api" { // Assuming you have a getter method
		t.Errorf("Expected baseURL to be https://api.etherscan.io/api, got %s", client.GetBaseURL())
	}
}

func TestGetAccountBalance(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("module") != "account" {
			t.Error("Expected module=account")
		}
		if r.URL.Query().Get("action") != "balance" {
			t.Error("Expected action=balance")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status": "1",
			"message": "OK",
			"result": "1000000000000000000"
		}`))
	}))
	defer server.Close()

	client := etherscan.NewClient("test_api_key")
	client.SetBaseURL(server.URL) // Assuming you have a setter method

	address := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	balance, err := client.GetAccountBalance(address)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedBalance := "1000000000000000000"
	if balance != expectedBalance {
		t.Errorf("Expected balance %s, got %s", expectedBalance, balance)
	}
}

func TestGetAccountBalanceError(t *testing.T) {
	client := etherscan.NewClient("test_api_key")
	_, err := client.GetAccountBalance("invalid_address")

	if err == nil {
		t.Error("Expected error for invalid address, got nil")
	}
}

func TestGetTransactionHistory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("module") != "account" {
			t.Error("Expected module=account")
		}
		if r.URL.Query().Get("action") != "txlist" {
			t.Error("Expected action=txlist")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status": "1",
			"message": "OK",
			"result": [
				{
					"blockNumber": "12345",
					"timeStamp": "1609459200",
					"hash": "0x123...",
					"from": "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
					"to": "0x742d35Cc6634C0532925a3b844Bc454e4438f44f",
					"value": "1000000000000000000"
				}
			]
		}`))
	}))
	defer server.Close()

	client := etherscan.NewClient("test_api_key")
	client.SetBaseURL(server.URL)

	address := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
	transactions, err := client.GetTransactionHistory(address)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(transactions))
	}

	tx := transactions[0]
	if tx.BlockNumber != "12345" {
		t.Errorf("Expected block number 12345, got %s", tx.BlockNumber)
	}
}

func TestRateLimiting(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount > 5 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "1", "message": "OK", "result": "1000000000000000000"}`))
	}))
	defer server.Close()

	client := etherscan.NewClient("test_api_key")
	client.SetBaseURL(server.URL)

	address := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"

	for i := 0; i < 6; i++ {
		_, err := client.GetAccountBalance(address)
		if i < 5 && err != nil {
			t.Errorf("Expected no error for request %d, got %v", i+1, err)
		}
		if i == 5 && err == nil {
			t.Error("Expected rate limit error for 6th request, got nil")
		}
	}
}
