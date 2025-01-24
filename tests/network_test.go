package main

import (
	"testing"
	"time"
)

func TestLoadProviderFromURI(t *testing.T) {
	tests := []struct {
		name         string
		uriString    string
		timeout      time.Duration
		expectedType interface{}
		expectError  bool
	}{
		{
			name:         "Valid file URI",
			uriString:    "file:///tmp/ipc.sock",
			timeout:      0,
			expectedType: IPCProvider{},
			expectError:  false,
		},
		{
			name:         "Valid HTTP URI",
			uriString:    "http://example.com",
			timeout:      5 * time.Second,
			expectedType: HTTPProvider{},
			expectError:  false,
		},
		{
			name:         "Valid WebSocket URI",
			uriString:    "ws://example.com",
			timeout:      5 * time.Second,
			expectedType: WebsocketProvider{},
			expectError:  false,
		},
		{
			name:         "Unsupported scheme",
			uriString:    "ftp://example.com",
			timeout:      0,
			expectedType: nil,
			expectError:  true,
		},
		{
			name:         "Invalid URI",
			uriString:    "http://:invalid_uri",
			timeout:      0,
			expectedType: nil,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := LoadProviderFromURI(tt.uriString, tt.timeout)
			if (err != nil) != tt.expectError {
				t.Errorf("unexpected error status: got %v, want %v", err != nil, tt.expectError)
			}
			if err == nil && tt.expectedType != nil {
				switch tt.expectedType.(type) {
				case IPCProvider:
					if _, ok := provider.(IPCProvider); !ok {
						t.Errorf("expected IPCProvider, got %T", provider)
					}
				case HTTPProvider:
					if _, ok := provider.(HTTPProvider); !ok {
						t.Errorf("expected HTTPProvider, got %T", provider)
					}
				case WebsocketProvider:
					if _, ok := provider.(WebsocketProvider); !ok {
						t.Errorf("expected WebsocketProvider, got %T", provider)
					}
				default:
					t.Errorf("unexpected provider type: %T", provider)
				}
			}
		})
	}
}

func TestGetChainID(t *testing.T) {
	tests := []struct {
		name        string
		network     string
		expectedID  int
		expectError bool
	}{
		{"Valid network - mainnet", "mainnet", 1, false},
		{"Valid network - rinkeby", "rinkeby", 4, false},
		{"Case insensitive network", "POLYGON", 137, false},
		{"Invalid network", "invalidnetwork", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chainID, err := GetChainID(tt.network)
			if (err != nil) != tt.expectError {
				t.Errorf("unexpected error status: got %v, want %v", err != nil, tt.expectError)
			}
			if chainID != tt.expectedID {
				t.Errorf("unexpected chain ID: got %d, want %d", chainID, tt.expectedID)
			}
		})
	}
}

func TestGetNetwork(t *testing.T) {
	tests := []struct {
		name          string
		chainID       int
		expectedName  string
		expectError   bool
	}{
		{"Valid chain ID - mainnet", 1, "mainnet", false},
		{"Valid chain ID - polygon", 137, "polygon", false},
		{"Invalid chain ID", 9999, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			network, err := GetNetwork(tt.chainID)
			if (err != nil) != tt.expectError {
				t.Errorf("unexpected error status: got %v, want %v", err != nil, tt.expectError)
			}
			if network != tt.expectedName {
				t.Errorf("unexpected network name: got %q, want %q", network, tt.expectedName)
			}
		})
	}
}
