package ethereal

import (
	"testing"
	"strings"
	"github.com/stretchr/testify/assert"
)

func TestNewAccounts(t *testing.T) {
	a := NewAccounts()
	assert.NotNil(t, a, "NewAccounts should return a non-nil Accounts instance")
}

func TestGenerateSeedPhrase(t *testing.T) {
	tests := []struct {
		strength    int
		expectsError bool
		name         string
	}{
		{128, false, "Valid strength (128 bits)"},
		{256, false, "Valid strength (256 bits)"},
		{64, true, "Invalid strength (too low)"},
		{384, true, "Invalid strength (too high)"},
	}
	a := NewAccounts()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seedPhrase, err := a.GenerateSeedPhrase(tt.strength)

			if tt.expectsError {
				assert.Error(t, err, "Expected an error for strength %d", tt.strength)
			} else {
				assert.NoError(t, err, "Did not expect an error for strength %d", tt.strength)
				assert.NotEmpty(t, seedPhrase, "Seed phrase should not be empty")
				assert.True(t, strings.Count(seedPhrase, " ") >= 11, "Seed phrase should contain at least 12 words")
			}
		})
	}
}

func TestDeriveAccount(t *testing.T) {
	validSeedPhrase := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	invalidSeedPhrase := "invalid seed phrase"
	validPassphrase := "myPassphrase"
	tests := []struct {
		seedPhrase   string
		index        int
		passphrase   string
		expectsError bool
		name         string
	}{
		{validSeedPhrase, 0, "", false, "Valid seed phrase without passphrase"},
		{validSeedPhrase, 1, validPassphrase, false, "Valid seed phrase with passphrase"},
		{validSeedPhrase, -1, "", true, "Negative index"},
		{invalidSeedPhrase, 0, "", true, "Invalid seed phrase"},
	}
	a := NewAccounts()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := a.DeriveAccount(tt.seedPhrase, tt.index, tt.passphrase)

			if tt.expectsError {
				assert.Error(t, err, "Expected an error for seedPhrase %s and index %d", tt.seedPhrase, tt.index)
				assert.Nil(t, account, "Account should be nil when error occurs")
			} else {
				assert.NoError(t, err, "Did not expect an error for seedPhrase %s and index %d", tt.seedPhrase, tt.index)
				assert.NotNil(t, account, "Account should not be nil for valid inputs")
			}
		})
	}
}
