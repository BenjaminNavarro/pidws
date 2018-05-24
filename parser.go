package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	yaml "gopkg.in/yaml.v2"
)

func getConfigurationFilePath() string {
	currentUser, err := user.Current()
	CheckError(err)
	path := currentUser.HomeDir + "/.pidws.yaml"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Create(path)
	}
	return path
}

// ReadPidWsConfiguration parses a YAML configuration file and populates the configuration struct
func (configuration *PidWsConfiguration) Read() {
	source, err := ioutil.ReadFile(getConfigurationFilePath())
	CheckError(err)

	err = yaml.Unmarshal(source, configuration)
	CheckError(err)

	if len(configuration.Workspaces) == 0 {
		fmt.Println("[Warning] The 'workspaces' field in the configuration file is empty. You have to create a workspace before using workspace commands")
	}
	if configuration.Default == "" {
		fmt.Println("[Warning] The 'default' field in the configuration file is not set. You have to set the active workspace before using workspace commands")
	}
	if len(configuration.Repositories) == 0 {
		fmt.Println("[Warning] The 'repositories' field in the configuration file is empty. You have to set at least a repository before creating a workspace")
	}
}

// SavePidWsConfiguration saves the given configuration as a YAML file
func (configuration *PidWsConfiguration) Save() {
	out, err := yaml.Marshal(configuration)
	CheckError(err)

	err = ioutil.WriteFile(getConfigurationFilePath(), out, 'w')
	CheckError(err)
}
