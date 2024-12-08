package ethereal

import (
	"testing"
	"time"
)

func TestGetBlockByTimestamp(t *testing.T) {
	// Mock dependencies would be initialized here
	facade := NewEtherealFacade(
		&Etherscan{}, // mock
		nil,          // mock ethclient
		&Accounts{},  // mock
		&Cache{},     // mock
	)

	timestamp := time.Now().Unix()
	block, err := facade.GetBlockByTimestamp(timestamp)
	
	if err != nil {
		t.Errorf("GetBlockByTimestamp failed: %v", err)
	}
	
	if block <= 0 {
		t.Error("Expected block number to be greater than 0")
	}
}

func TestGenerateSeedPhrase(t *testing.T) {
	facade := NewEtherealFacade(
		&Etherscan{}, // mock
		nil,          // mock ethclient
		&Accounts{},  // mock
		&Cache{},     // mock
	)

	phrase, err := facade.GenerateSeedPhrase(128)
	
	if err != nil {
		t.Errorf("GenerateSeedPhrase failed: %v", err)
	}
	
	if len(phrase) == 0 {
		t.Error("Expected non-empty seed phrase")
	}
} 