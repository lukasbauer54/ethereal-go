package ethereal

import (
	"encoding/json"
	"errors"
	"fmt"
)

// EventFilter defines parameters for filtering contract events
type EventFilter struct {
	FromBlock uint64
	ToBlock   uint64
	Address   string
	Topics    []string
}

// Contracts handles Ethereum contract operations
type Contracts struct {
	etherscan *Etherscan
	cache     *Cache
}

// NewContracts creates a new Contracts instance
func NewContracts(etherscan *Etherscan, cache *Cache) *Contracts {
	return &Contracts{
		etherscan: etherscan,
		cache:     cache,
	}
}

// GetABI retrieves and parses the ABI for a contract
func (c *Contracts) GetABI(address string, resolveProxy bool) (map[string]interface{}, error) {
	if address == "" {
		return nil, errors.New("address cannot be empty")
	}

	cacheKey := fmt.Sprintf("abi_%s_%v", address, resolveProxy)

	// Try to get from cache first
	if cached, err := c.cache.Get(cacheKey); err == nil {
		return cached.(map[string]interface{}), nil
	}

	// Get ABI from Etherscan
	abiString, err := c.etherscan.GetContractABI(address)
	if err != nil {
		return nil, fmt.Errorf("failed to get ABI from Etherscan: %w", err)
	}

	// Parse ABI
	var abiMap map[string]interface{}
	if err := json.Unmarshal([]byte(abiString), &abiMap); err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Cache the result
	if err := c.cache.Set(cacheKey, abiMap); err != nil {
		return nil, fmt.Errorf("failed to cache ABI: %w", err)
	}

	return abiMap, nil
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
