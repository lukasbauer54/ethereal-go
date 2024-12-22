package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Provider interface{}

type HTTPProvider struct {
	URI     string
	Timeout time.Duration
}

type IPCProvider struct {
	Path string
}

type WebsocketProvider struct {
	URI     string
	Timeout time.Duration
}

var ChainIDsIndex = map[int]string{
	1:     "mainnet",
	3:     "ropsten",
	4:     "rinkeby",
	5:     "goerli",
	42:    "kovan",
	137:   "polygon",
	43114: "avalanche",
	250:   "fantom",
	42161: "arbitrum",
	10:    "optimism",
}

var NetworksIndex = func() map[string]int {
	index := make(map[string]int)
	for id, name := range ChainIDsIndex {
		index[strings.ToLower(name)] = id
	}
	return index
}()

func LoadProviderFromURI(uriString string, timeout time.Duration) (Provider, error) {
	parsedURI, err := url.Parse(uriString)
	if err != nil {
		return nil, fmt.Errorf("invalid URI: %v", err)
	}

	switch parsedURI.Scheme {
	case "file":
		return IPCProvider{Path: parsedURI.Path}, nil
	case "http", "https":
		return HTTPProvider{URI: uriString, Timeout: timeout}, nil
	case "ws", "wss":
		return WebsocketProvider{URI: uriString, Timeout: timeout}, nil
	default:
		return nil, fmt.Errorf("Web3 does not know how to connect to scheme %q in %q", parsedURI.Scheme, uriString)
	}
}

func GetChainID(network string) (int, error) {
	chainID, exists := NetworksIndex[strings.ToLower(network)]
	if !exists {
		return 0, errors.New("network not found")
	}
	return chainID, nil
}

func GetNetwork(chainID int) (string, error) {
	network, exists := ChainIDsIndex[chainID]
	if !exists {
		return "", errors.New("chain ID not found")
	}
	return network, nil
}
