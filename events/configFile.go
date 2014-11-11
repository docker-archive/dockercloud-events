package events

import (
	"encoding/json"
	"log"
	"os"
)

var conf *Configuration

type Configuration struct {
	CertCommonName  string
	DockerBinaryURL string
	DockerHost      string
	TutumHost       string
	TutumToken      string
	TutumUUID       string
}

func parseConfigFile(file string) (*Configuration, error) {
	var conf Configuration
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	//read and decode json format config file
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func GetConf(configFile string) (*Configuration){
    if conf == nil {
	    configuration, err := parseConfigFile(configFile)
	    if err != nil {
		    SendError(err)
		    log.Fatalf("Cannot parse configuration file(%s):%s\n", configFile, err.Error())
    	}
        conf = configuration    
    }
    return conf
}
