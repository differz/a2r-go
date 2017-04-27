package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// Configuration file name
const (
	configFile string = "../config.json"
)

// A Configuration structure
type Configuration struct {
	AsteriskServer   string
	AsteriskUser     string
	AsteriskPassword string
	AsteriskDatabase string

	HelpDeskServer   string
	HelpDeskUser     string
	HelpDeskPassword string
	HelpDeskDatabase string

	RedmineServer string
	RedmineKey    string

	InputMethod  string
	OutputMethod string
}

// Singleton object pointer
var (
	config *Configuration
	once   sync.Once
)

// New return pointer to configuration object
// May be read from file or create new
func New() *Configuration {

	once.Do(func() {
		config = &Configuration{}
	})

	return config
}

// Init read config from file, if not exist create a new one
func (conf *Configuration) Init() error {

	err := conf.Read()
	if err != nil {
		err = conf.Create()
	}

	return err
}

// Read function reads configuration from file and decode it
func (conf *Configuration) Read() error {

	file, err := os.Open(configFile)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return err
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(conf)
	if err != nil {
		fmt.Println("Error parsing config file: ", err)
		log.Fatal(err)
		// can't get here too
	}

	return err
}

// Create function return empty config file
func (conf *Configuration) Create() error {

	// test for camelCase, MixedCase, lowercase
	var jsonBlob = []byte(`{
		"asteriskServer"  : "192.168.0.4",
		"AsteriskUser"    : ""
		"AsteriskPassword": ""
		"AsteriskDatabase": ""
		"HelpDeskServer"  : "192.168.0.1",
		"HelpDeskUser"    : "",
		"HelpDeskPassword": "",
		"HelpDeskDatabase": "",
		"redmineserver"   : "192.168.0.22",
		"RedmineKey"      : "70d0affd9bc31f5ea60d6cf15483318ca9a2a8db",
		"InputMethod"     : "CDR",
		"OutputMethod"    : "API"
	}`)
	optimize := true

	err := json.Unmarshal(jsonBlob, conf)
	if err == nil {
		if optimize {
			jsonBlob, _ = json.Marshal(conf)
		}
		err = ioutil.WriteFile(configFile, jsonBlob, 0644)
	}

	return err
}

////////////////////////////////////////
// fmt.Printf("%+v", conf)
