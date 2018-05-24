package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// Get the path to the configuration file and parse it
	var configuration PidWsConfiguration
	configuration.Read()

	/* Set the flags accepted by the command line.
	 * All arguments not starting with a '-' will be passed to the make command of the selected workspace.
	 */
	cdWsPtr := flag.String("cd", "", "Change the current directory to the specified workspace (a new shell will be created).")
	createWsPtr := flag.String("create", "", "Create a new workspace in the current directory. To specify its name, see -name")
	defaultWsPtr := flag.String("default", "", "Set the workspace that will receive workspace commands when none is specified (see -use).")
	namePtr := flag.String("name", "", "To be used together with -create-ws, -register-workspace or -register-repository to specify the name of the new workspace or repository.")
	useWsPtr := flag.String("use", "", "If set, it will override the default workspace for the current commands.")
	removeWsPtr := flag.String("remove-workspace", "", "Remove a workspace from the list of the known ones. Files won't be deleted.")
	removeRepoPtr := flag.String("remove-repository", "", "Remove a repository from the list of the known ones.")
	registerWsPtr := flag.String("register-workspace", "", "Add an already existing workspace in the configuration file.")
	registerRepoPtr := flag.String("register-repository", "", "Add a new repository in the configuration file.")

	flag.Parse()

	// If no argument is passed, print usage and exit
	if len(os.Args) == 1 {
		defaultWsPath, _ := configuration.GetWorkspacePath(configuration.Default)
		fmt.Println("pidws is an utility to help you manage and interact with your PID workspaces.")
		fmt.Println("If you want to find out more about PID, go and check http://pid.lirmm.net/pid-framework")
		if defaultWsPath != "" {
			fmt.Println("\nThe current default workspace is '" + configuration.Default + " '(" + defaultWsPath + ").\n")
		}
		fmt.Println("usage:")
		flag.PrintDefaults()
		fmt.Println("\nAll other commands will be passed to the active PID workspace.")
		os.Exit(0)
	}

	// Handle the flags if specified by the user
	if *removeWsPtr != "" {
		path, err := configuration.GetWorkspacePath(*removeWsPtr)
		CheckError(err)

		configuration.RemoveWorkspace(*removeWsPtr)

		fmt.Println("The '" + *removeWsPtr + "' workspace (" + path + ") has been successfully removed from the known workspaces.")
	}
	if *removeRepoPtr != "" {
		repo, err := configuration.FindRepository(*removeRepoPtr)
		CheckError(err)

		configuration.RemoveRepository(*removeRepoPtr)

		fmt.Println("The '" + repo.Name + "' repository (" + repo.Address + ") has been successfully removed from the known repositories.")
	}
	if *registerWsPtr != "" {
		if *namePtr == "" {
			ExitWithError("You must provide a -name argument together with -register-workspace")
		}
		err := configuration.RegisterWorkspace(*namePtr, *registerWsPtr)
		CheckError(err)

		configuration.PrintWorkspaceSuccessMessage(*namePtr, "has been registered successfully")
	}
	if *registerRepoPtr != "" {
		if *namePtr == "" {
			ExitWithError("You must provide a -name argument together with -register-repository")
		}
		err := configuration.RegisterRepository(*namePtr, *registerRepoPtr)
		CheckError(err)

		fmt.Println("The '" + *namePtr + "' repository (" + *registerRepoPtr + ") has been registered successfully.")
	}
	if *createWsPtr != "" {
		if *namePtr == "" {
			*namePtr = "pid"
		}
		if configuration.DoesWorkspaceExist(*namePtr) {
			log.Fatal("You cannot create a workspace with the name '" + *namePtr + "' because it already exists")
		}

		repo, err := configuration.FindRepository(*createWsPtr)
		CheckError(err)

		var dirName = *namePtr + "-workspace"
		err = CloneRepository(repo, dirName)
		CheckError(err)

		err = configuration.AddWorkspace(*namePtr, dirName)
		CheckError(err)

		err = configuration.ConfigureWorkspace(*namePtr)
		CheckError(err)

		configuration.PrintWorkspaceSuccessMessage(*namePtr, "has been successfully created using the '"+repo.Name+"' repository")
	}
	if *defaultWsPtr != "" {
		err := configuration.SetDefaultWorkspace(*defaultWsPtr)
		CheckError(err)
		configuration.PrintWorkspaceSuccessMessage(*defaultWsPtr, "is now the default workspace")
	}
	if *cdWsPtr != "" {
		err := configuration.OpenShellInWorkspace(*cdWsPtr)
		CheckError(err)
	}

	// if some arguments where not parsed by 'flag' forward them to the selected workspace
	if len(flag.Args()) > 0 {
		// Set the workspace to be used, being either the default one or the one specified with the 'use-ws' flag
		wsInUse, err := configuration.GetDefaultWorkspace(*useWsPtr)
		CheckError(err)

		err = configuration.ExecuteCommandInWorkspace(wsInUse, flag.Args())
		CheckError(err)
	}
}
