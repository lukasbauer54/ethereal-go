package ethereal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Network endpoints for different chains
var endpoints = map[int]string{
	1:      "https://api.etherscan.io",
	3:      "https://api-ropsten.etherscan.io",
	4:      "https://api-rinkeby.etherscan.io",
	5:      "https://api-goerli.etherscan.io",
	42:     "https://api-kovan.etherscan.io",
	137:    "https://api.polygonscan.com",
	43114:  "https://api.avax.network",
	250:    "https://api.ftmscan.com",
	42161:  "https://api.arbiscan.io",
	10:     "https://api-optimistic.etherscan.io",
}

const ethStartTimestamp = 1438214400 // July 30, 2015 UTC

// EtherscanError represents an error from the Etherscan API
type EtherscanError struct {
	Message string
}

func (e *EtherscanError) Error() string {
	return fmt.Sprintf("etherscan error: %s", e.Message)
}

// EtherscanNetworkConfig represents configuration for a specific network
type EtherscanNetworkConfig struct {
	Key string `json:"key"`
}

// EtherscanConfig represents the complete Etherscan configuration
type EtherscanConfig struct {
	ChainID    int                             `json:"chain_id"`
	Timeout    int                             `json:"timeout"`
	Networks   map[string]EtherscanNetworkConfig `json:"networks"`
}

// Etherscan provides access to Etherscan API functionality
type Etherscan struct {
	config   EtherscanConfig
	cache    *Cache
	chainID  int
	client   *http.Client
}

// NewEtherscan creates a new Etherscan instance
func NewEtherscan(config EtherscanConfig, cache *Cache) *Etherscan {
	return &Etherscan{
		config:  config,
		cache:   cache,
		chainID: config.ChainID,
		client:  &http.Client{Timeout: time.Duration(config.Timeout) * time.Second},
	}
}

// GetBlockByTimestamp gets the block number for a given timestamp
func (e *Etherscan) GetBlockByTimestamp(timestamp int64, closest string) (int64, error) {
	if closest == "" {
		closest = "after"
	}
	
	cacheKey := fmt.Sprintf("etherscan:block:%d:%s", timestamp, closest)
	if cached, err := e.cache.Get(cacheKey); err == nil {
		return cached.(int64), nil
	}

	params := map[string]string{
		"module":    "block",
		"action":    "getblocknobytime",
		"timestamp": strconv.FormatInt(timestamp, 10),
		"closest":   closest,
	}

	result, err := e.fetch(params)
	if err != nil {
		return 0, err
	}

	blockNum, err := strconv.ParseInt(result.(string), 10, 64)
	if err != nil {
		return 0, err
	}

	e.cache.Set(cacheKey, blockNum)
	return blockNum, nil
}

// ToBlock converts a timestamp to a block number
func (e *Etherscan) ToBlock(ts interface{}) (int64, error) {
	switch v := ts.(type) {
	case int64:
		if v < ethStartTimestamp {
			return v, nil
		}
		return e.GetBlockByTimestamp(v, "after")
	case time.Time:
		return e.GetBlockByTimestamp(v.Unix(), "after")
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return 0, err
		}
		return e.GetBlockByTimestamp(t.Unix(), "after")
	default:
		return 0, fmt.Errorf("unsupported timestamp type")
	}
}

// GetABI gets the ABI for a given address
func (e *Etherscan) GetABI(address string) (string, error) {
	cacheKey := fmt.Sprintf("etherscan:abi:%s", address)
	if cached, err := e.cache.Get(cacheKey); err == nil {
		return cached.(string), nil
	}

	params := map[string]string{
		"module":  "contract",
		"action":  "getabi",
		"address": address,
	}

	result, err := e.fetch(params)
	if err != nil {
		return "", err
	}

	abi := result.(string)
	e.cache.Set(cacheKey, abi)
	return abi, nil
}

type etherscanResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func (e *Etherscan) fetch(params map[string]string) (interface{}, error) {
	endpoint := endpoints[e.chainID]
	if endpoint == "" {
		return nil, fmt.Errorf("unsupported chain ID: %d", e.chainID)
	}

	url := fmt.Sprintf("%s/api?", endpoint)
	for k, v := range params {
		url += fmt.Sprintf("%s=%s&", k, v)
	}
	
	network := getNetwork(e.chainID)
	apiKey := e.config.Networks[network].Key
	url += fmt.Sprintf("apiKey=%s", apiKey)

	resp, err := e.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result etherscanResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Status != "1" {
		return nil, &EtherscanError{Message: result.Message}
	}

	return result.Result, nil
}

func getNetwork(chainID int) string {
	switch chainID {
	case 1:
		return "mainnet"
	case 137:
		return "polygon"
	case 43114:
		return "avalanche"
	case 250:
		return "ftm"
	case 42161:
		return "arbitrum"
	case 10:
		return "optimism"
	default:
		return "mainnet"
	}
} 