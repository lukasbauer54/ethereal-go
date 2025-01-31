package ethereal

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Etherscan interface defines the methods required from an Etherscan client
type EtherscanClient interface {
	GetContractABI(address string) (string, error)
	GetContractSource(address string) (string, error)
}

// CacheClient interface defines the methods required from a cache implementation
type CacheClient interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
}

// EventFilter defines parameters for filtering contract events
type EventFilter struct {
	FromBlock uint64
	ToBlock   uint64
	Address   string
	Topics    []string
}

// ContractEvent represents a contract event definition
type ContractEvent struct {
	Name      string
	Anonymous bool
	Inputs    []EventInput
}

// EventInput represents an event parameter
type EventInput struct {
	Name       string
	Type       string
	Indexed    bool
	Components []EventInput
}

// ProxyInfo contains information about a proxy contract
type ProxyInfo struct {
	IsProxy        bool
	Implementation string
	ProxyType      string // "EIP1967", "EIP897", "Custom"
}

// FunctionCall represents a contract function call
type FunctionCall struct {
	Name   string
	Inputs []FunctionParam
}

// FunctionParam represents a function parameter
type FunctionParam struct {
	Name  string
	Type  string
	Value interface{}
}

// Contracts handles Ethereum contract operations
type Contracts struct {
	etherscan EtherscanClient
	cache     CacheClient
}

// NewContracts creates a new Contracts instance
func NewContracts(etherscan EtherscanClient, cache CacheClient) *Contracts {
	return &Contracts{
		etherscan: etherscan,
		cache:     cache,
	}
}

// GetABI retrieves and parses the ABI for a contract
func (c *Contracts) GetABI(address string, resolveProxy bool) ([]map[string]interface{}, error) {
	if address == "" {
		return nil, errors.New("address cannot be empty")
	}

	cacheKey := fmt.Sprintf("abi_%s_%v", address, resolveProxy)

	// Try to get from cache first
	if cached, err := c.cache.Get(cacheKey); err == nil {
		return cached.([]map[string]interface{}), nil
	}

	// Get ABI from Etherscan
	abiString, err := c.etherscan.GetContractABI(address)
	if err != nil {
		return nil, fmt.Errorf("failed to get ABI from Etherscan: %w", err)
	}

	// Parse ABI
	var abiArray []map[string]interface{}
	if err := json.Unmarshal([]byte(abiString), &abiArray); err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Cache the result
	if err := c.cache.Set(cacheKey, abiArray); err != nil {
		return nil, fmt.Errorf("failed to cache ABI: %w", err)
	}

	return abiArray, nil
}

// ListEvents returns all event names defined in the contract
func (c *Contracts) ListEvents(address string, resolveProxy bool) ([]string, error) {
	abiMap, err := c.GetABI(address, resolveProxy)
	if err != nil {
		return nil, err
	}

	var events []string
	for name, item := range abiMap {
		if itemMap, ok := item.(map[string]interface{}); ok {
			if itemMap["type"] == "event" {
				events = append(events, name)
			}
		}
	}

	return events, nil
}

// GetEvents retrieves events from a contract
func (c *Contracts) GetEvents(address string, event string, filter EventFilter, resolveProxy bool) ([]interface{}, error) {
	if address == "" {
		return nil, errors.New("address cannot be empty")
	}
	if event == "" {
		return nil, errors.New("event name cannot be empty")
	}

	// Validate filter parameters
	if filter.FromBlock > filter.ToBlock && filter.ToBlock != 0 {
		return nil, errors.New("fromBlock cannot be greater than toBlock")
	}

	// Get ABI to validate event exists
	abi, err := c.GetABI(address, resolveProxy)
	if err != nil {
		return nil, fmt.Errorf("failed to get ABI: %w", err)
	}

	// Verify event exists in ABI
	eventExists := false
	for _, item := range abi {
		if itemMap, ok := item.(map[string]interface{}); ok {
			if itemMap["type"] == "event" && itemMap["name"] == event {
				eventExists = true
				break
			}
		}
	}

	if !eventExists {
		return nil, fmt.Errorf("event %s does not exist in contract ABI", event)
	}

	// Implementation would need to use etherscan API or direct blockchain connection
	// to fetch the actual events. This is a placeholder that returns an empty slice.
	return []interface{}{}, nil
}

// GetFunctionSignature returns the function signature for a given function name
func (c *Contracts) GetFunctionSignature(address string, functionName string, resolveProxy bool) (string, error) {
	if address == "" {
		return "", errors.New("address cannot be empty")
	}
	if functionName == "" {
		return "", errors.New("function name cannot be empty")
	}

	abi, err := c.GetABI(address, resolveProxy)
	if err != nil {
		return "", fmt.Errorf("failed to get ABI: %w", err)
	}

	for _, item := range abi {
		if itemMap, ok := item.(map[string]interface{}); ok {
			if itemMap["type"] == "function" && itemMap["name"] == functionName {
				inputs, ok := itemMap["inputs"].([]interface{})
				if !ok {
					return "", errors.New("invalid ABI format: inputs not found")
				}

				signature := functionName + "("
				for i, input := range inputs {
					inputMap, ok := input.(map[string]interface{})
					if !ok {
						return "", errors.New("invalid ABI format: input type not found")
					}
					if i > 0 {
						signature += ","
					}
					signature += inputMap["type"].(string)
				}
				signature += ")"
				return signature, nil
			}
		}
	}

	return "", fmt.Errorf("function %s not found in contract ABI", functionName)
}

// IsContract checks if the given address is a contract
func (c *Contracts) IsContract(address string) (bool, error) {
	if address == "" {
		return false, errors.New("address cannot be empty")
	}

	cacheKey := fmt.Sprintf("is_contract_%s", address)

	// Try to get from cache first
	if cached, err := c.cache.Get(cacheKey); err == nil {
		return cached.(bool), nil
	}

	// Get contract source code from Etherscan
	// If it returns source code, it's a contract
	abi, err := c.etherscan.GetContractABI(address)
	isContract := err == nil && len(abi) > 0

	// Cache the result
	if err := c.cache.Set(cacheKey, isContract); err != nil {
		return false, fmt.Errorf("failed to cache contract check: %w", err)
	}

	return isContract, nil
}

// GetProxyInfo checks if a contract is a proxy and returns its implementation
func (c *Contracts) GetProxyInfo(address string) (*ProxyInfo, error) {
	if address == "" {
		return nil, errors.New("address cannot be empty")
	}

	cacheKey := fmt.Sprintf("proxy_info_%s", address)

	// Try to get from cache first
	if cached, err := c.cache.Get(cacheKey); err == nil {
		return cached.(*ProxyInfo), nil
	}

	// Get contract source code to check for proxy patterns
	source, err := c.etherscan.GetContractSource(address)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract source: %w", err)
	}

	info := &ProxyInfo{
		IsProxy: false,
	}

	// Check for EIP-1967 proxy
	if containsEIP1967Pattern(source) {
		info.IsProxy = true
		info.ProxyType = "EIP1967"
		// Implementation would need to read the implementation slot
		// 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc
	}

	// Check for EIP-897 proxy
	if containsEIP897Pattern(source) {
		info.IsProxy = true
		info.ProxyType = "EIP897"
		// Implementation would need to call implementation() function
	}

	// Cache the result
	if err := c.cache.Set(cacheKey, info); err != nil {
		return nil, fmt.Errorf("failed to cache proxy info: %w", err)
	}

	return info, nil
}

// EncodeFunctionCall encodes a function call into calldata
func (c *Contracts) EncodeFunctionCall(address string, call FunctionCall) (string, error) {
	if address == "" {
		return "", errors.New("address cannot be empty")
	}
	if call.Name == "" {
		return "", errors.New("function name cannot be empty")
	}

	// Get function signature
	signature, err := c.GetFunctionSignature(address, call.Name, true)
	if err != nil {
		return "", fmt.Errorf("failed to get function signature: %w", err)
	}

	// Implementation would need to:
	// 1. Calculate function selector (first 4 bytes of keccak256 of signature)
	// 2. Encode parameters according to their types
	// 3. Concatenate selector and encoded parameters

	return "", errors.New("not implemented")
}

// Helper functions for proxy detection
func containsEIP1967Pattern(source string) bool {
	// Implementation would check for EIP-1967 storage slots or patterns
	return false
}

func containsEIP897Pattern(source string) bool {
	// Implementation would check for EIP-897 implementation() function
	return false
}

// GetImplementationAddress returns the implementation address for a proxy contract
func (c *Contracts) GetImplementationAddress(address string) (string, error) {
	info, err := c.GetProxyInfo(address)
	if err != nil {
		return "", fmt.Errorf("failed to get proxy info: %w", err)
	}

	if !info.IsProxy {
		return "", fmt.Errorf("address %s is not a proxy contract", address)
	}

	if info.Implementation == "" {
		return "", fmt.Errorf("implementation address not found for proxy %s", address)
	}

	return info.Implementation, nil
}
