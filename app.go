package main

import (
	"fmt"
	"log"
	"strings"
)

type Web3 struct {
	ChainID int
}

type EtherealFacade struct{}

type AppContainer struct {
	Config AppConfig
}

type AppConfig struct {
	Logging    LoggingConfig
	Etherscan  EtherscanConfig
}

type LoggingConfig struct {
	Loggers map[string]*Logger
}

type Logger struct {
	Level string
}

type EtherscanConfig struct {
	ChainID int
}

type Ethereal struct {
	w3       *Web3
	logLevel *string
	facade   *EtherealFacade
}

func NewEthereal(w3 *Web3, logLevel *string) *Ethereal {
	return &Ethereal{
		w3:       w3,
		logLevel: logLevel,
	}
}

func (e *Ethereal) GetAttribute(name string) interface{} {
	if name != "e" {
		return e.w3
	}

	if e.facade == nil {
		app, err := e.initApp()
		if err != nil {
			log.Fatalf("Error initializing Ethereal: %v", err)
		}
		e.facade = app.EtherealFacade(e.w3)
	}

	return e.facade
}

func (e *Ethereal) initApp() (*AppContainer, error) {
	var logLevel string
	if e.logLevel != nil {
		logLevel = strings.ToUpper(*e.logLevel)
	}

	chainID := e.w3.ChainID

	app := &AppContainer{
		Config: AppConfig{
			Logging: LoggingConfig{
				Loggers: map[string]*Logger{
					"root": {},
				},
			},
			Etherscan: EtherscanConfig{},
		},
	}

	if logLevel != "" {
		app.Config.Logging.Loggers["root"].Level = logLevel
	}

	app.Config.Etherscan.ChainID = chainID

	if err := app.InitResources(); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *AppContainer) EtherealFacade(w3 *Web3) *EtherealFacade {
	return &EtherealFacade{}
}

func (app *AppContainer) InitResources() error {
	fmt.Println("Initializing resources...")
	return nil
}

func main() {
	w3 := &Web3{ChainID: 1}
	logLevel := "INFO"
	ethereal := NewEthereal(w3, &logLevel)

	facade := ethereal.GetAttribute("e")
	fmt.Printf("Facade: %v\n", facade)
}
