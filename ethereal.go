package ethereal

import (
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EtherealFacade is the main class containing Ethereal's functionality
type EtherealFacade struct {
	etherscan *Etherscan
	web3      *ethclient.Client
	accounts  *Accounts
	cache     *Cache
}

// NewEtherealFacade creates a new instance of EtherealFacade
func NewEtherealFacade(etherscan *Etherscan, web3 *ethclient.Client, accounts *Accounts, cache *Cache) *EtherealFacade {
	return &EtherealFacade{
		etherscan: etherscan,
		web3:      web3,
		accounts:  accounts,
		cache:     cache,
	}
}

// GetBlockByTimestamp gets the block number for a given timestamp
func (e *EtherealFacade) GetBlockByTimestamp(timestamp int64) (int64, error) {
	return e.etherscan.GetBlockByTimestamp(timestamp)
}

// GetABI gets the ABI for a given address
func (e *EtherealFacade) GetABI(address string, resolveProxy bool) (map[string]interface{}, error) {
	contracts := e.contracts()
	return contracts.GetABI(address, resolveProxy)
}

// ListEvents gets a list of events for a given address
func (e *EtherealFacade) ListEvents(address string, resolveProxy bool) ([]string, error) {
	contracts := e.contracts()
	return contracts.ListEvents(address, resolveProxy)
}

// GetContract gets a contract for a given address
func (e *EtherealFacade) GetContract(address string, resolveProxy bool) (*bind.BoundContract, error) {
	contracts := e.contracts()
	return contracts.GetContract(address, resolveProxy)
}

// Account represents a derived Ethereum account
type Account struct {
	PublicKey  common.Address
	PrivateKey string
}

// DeriveAccount derives public and private key from a seed phrase
func (e *EtherealFacade) DeriveAccount(seedPhrase string, index int, passphrase string) (*Account, error) {
	return e.accounts.DeriveAccount(seedPhrase, index, passphrase)
}

// GenerateSeedPhrase generates a mnemonic
func (e *EtherealFacade) GenerateSeedPhrase(strength int) (string, error) {
	return e.accounts.GenerateSeedPhrase(strength)
}

// EventFilter represents filters for event queries
type EventFilter struct {
	FromTime        time.Time
	ToTime          time.Time
	ArgumentFilters map[string]interface{}
}

// GetEvents gets events for a given address
func (e *EtherealFacade) GetEvents(address string, event string, filter EventFilter, resolveProxy bool) ([]interface{}, error) {
	contracts := e.contracts()
	return contracts.GetEvents(address, event, filter, resolveProxy)
}

func (e *EtherealFacade) contracts() *Contracts {
	return NewContracts(e.etherscan, e.web3, e.cache)
}
