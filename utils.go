package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func unknownWorkspaceError(name string) string {
	return "The '" + name + "' workspace is unkown."
}

func unknownRepositoryError(name string) string {
	return "The '" + name + "' repository is unkown."
}

// ExitWithError print an error and exit the programm
func ExitWithError(err string) {
	fmt.Println("[Fatal error]", err)
	os.Exit(-1)
}

// CheckError if there is an error, print it and exit the program
func CheckError(err error) {
	if err != nil {
		ExitWithError(err.Error())
	}
}

// GetWorkspacePath get the path to the given workspace as found in the configuration
func (configuration *PidWsConfiguration) GetWorkspacePath(name string) (string, error) {
	var pathToWs = ""
	for _, ws := range configuration.Workspaces {
		if ws.Name == name {
			pathToWs = ws.Path
		}
	}
	if pathToWs == "" {
		return pathToWs, errors.New(unknownWorkspaceError(name))
	}
	return pathToWs, nil
}

// GetDefaultWorkspace get the default workspace. Return name if not empty and the one from the configuration otherwise
func (configuration *PidWsConfiguration) GetDefaultWorkspace(name string) (string, error) {
	var ws = configuration.Default
	if name != "" {
		ws = name
	}
	_, err := configuration.GetWorkspacePath(ws)
	return ws, err
}

// DoesWorkspaceExist check if the given workspace exists in the configuration
func (configuration *PidWsConfiguration) DoesWorkspaceExist(name string) bool {
	for _, ws := range configuration.Workspaces {
		if ws.Name == name {
			return true
		}
	}
	return false
}

// SetDefaultWorkspace set the default workspace in the configuration
func (configuration *PidWsConfiguration) SetDefaultWorkspace(name string) error {
	if configuration.DoesWorkspaceExist(name) {
		configuration.Default = name
		configuration.Save()
		return nil
	}
	return errors.New(unknownWorkspaceError(name))
}

// RegisterWorkspace add a new worksapce in the configuration
func (configuration *PidWsConfiguration) RegisterWorkspace(name string, directory string) error {
	if !configuration.DoesWorkspaceExist(name) {
		var newWs WorkspacePath
		newWs.Name = name
		newWs.Path = directory
		configuration.Workspaces = append(configuration.Workspaces, newWs)
		configuration.Save()
		return nil
	}
	return errors.New("The '" + name + "' workspace already exists.")
}

// AddWorkspace add a new worksapce in the configuration
func (configuration *PidWsConfiguration) AddWorkspace(name string, directory string) error {
	currentPath, _ := os.Getwd()
	return configuration.RegisterWorkspace(name, currentPath+"/"+directory)
}

// RemoveWorkspace remove the workspace from the configuration
func (configuration *PidWsConfiguration) RemoveWorkspace(name string) error {
	for idx, ws := range configuration.Workspaces {
		if ws.Name == name {
			configuration.Workspaces = append(configuration.Workspaces[:idx], configuration.Workspaces[idx+1:]...)
			configuration.Save()
			return nil
		}
	}
	return errors.New(unknownWorkspaceError(name))
}

// ConfigureWorkspace configure the workspace by calling cmake from its pid folder
func (configuration *PidWsConfiguration) ConfigureWorkspace(name string) error {
	path, _ := configuration.GetWorkspacePath(name)
	os.Chdir(path + "/pid")
	cmd := exec.Command("cmake", []string{".."}...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// OpenShellInWorkspace Start a new shell from the workspace root directory
func (configuration *PidWsConfiguration) OpenShellInWorkspace(name string) error {
	pathToWs, err := configuration.GetWorkspacePath(name)
	if err != nil {
		return err
	}
	shell := os.Getenv("SHELL")
	err = syscall.Exec(shell, []string{shell, "-c", "cd " + pathToWs + " && exec \"" + shell + "\""}, os.Environ())
	return err
}

// ExecuteCommandInWorkspace execute the given command inside the workspace by passing the arguments to make
func (configuration *PidWsConfiguration) ExecuteCommandInWorkspace(name string, args []string) error {
	pathToWs, err := configuration.GetWorkspacePath(name)
	if err != nil {
		return err
	}

	os.Chdir(pathToWs + "/pid")
	cmd := exec.Command("make", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (configuration *PidWsConfiguration) DoesRepositoryExists(name string) bool {
	for _, repo := range configuration.Repositories {
		if repo.Name == name {
			return true
		}
	}
	return false
}

func (configuration *PidWsConfiguration) RegisterRepository(name string, address string) error {
	if configuration.DoesRepositoryExists(name) {
		return errors.New("The '" + name + "' already exists.")
	}
	var newRepo Repository
	newRepo.Name = name
	newRepo.Address = address
	configuration.Repositories = append(configuration.Repositories, newRepo)
	configuration.Save()
	return nil
}

// RemoveWorkspace remove the workspace from the configuration
func (configuration *PidWsConfiguration) RemoveRepository(name string) error {
	for idx, ws := range configuration.Repositories {
		if ws.Name == name {
			configuration.Repositories = append(configuration.Repositories[:idx], configuration.Repositories[idx+1:]...)
			configuration.Save()
			return nil
		}
	}
	return errors.New(unknownRepositoryError(name))
}

// FindRepo find a repository in the configuration
func (configuration *PidWsConfiguration) FindRepository(repository string) (Repository, error) {
	for _, repo := range configuration.Repositories {
		if repo.Name == repository {
			return repo, nil
		}
	}
	return Repository{}, errors.New("The '" + repository + "' repository is unkown.")
}

// CloneRepo clone a repository inside the given directory
func CloneRepository(repo Repository, directory string) error {
	cmd := exec.Command("git", []string{"clone", repo.Address, directory}...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// PrintWorkspaceSuccessMessage helper function to print messages on success from the main function
func (configuration *PidWsConfiguration) PrintWorkspaceSuccessMessage(name string, message string) {
	path, _ := configuration.GetWorkspacePath(name)
	fmt.Println("The '" + name + "' workspace (" + path + ") " + message + ".")
}

// PrintRepositorySuccessMessage helper function to print messages on success from the main function
func (configuration *PidWsConfiguration) PrintRepositorySuccessMessage(name string, message string) {
	repo, _ := configuration.FindRepository(name)
	fmt.Println("The '" + repo.Name + "' repository (" + repo.Address + ") " + message + ".")
}
