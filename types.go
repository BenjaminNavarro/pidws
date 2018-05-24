package main

import "fmt"

// WorkspacePath stores the path of a specific workspace
type WorkspacePath struct {
	Name string
	Path string
}

// Repository stores the path of a specific workspace
type Repository struct {
	Name    string
	Address string
}

// PidWsConfiguration configuration parameters for the access and creation of PID workspaces.
// Fields must start with an uppercase letter
type PidWsConfiguration struct {
	Workspaces   []WorkspacePath
	Default      string
	Repositories []Repository
}

// Print outputs all the fields to the standard output
func (configuration *PidWsConfiguration) Print() {
	fmt.Println("Workspaces:", configuration.Workspaces)
	fmt.Println("Active:", configuration.Default)
	fmt.Println("Repositories:", configuration.Repositories)
}
