# Introduction

With pidws you can easily manage your PID workspaces and execute workspace commands from anywhere.

The main features are the:
 * Creation/registration/removal of PID workspaces
 * Registration/removal of PID workspace Git repositories
 * Setting of the default workspace to forward commands to
 * Change directory to one of your workspaces

# Installation

## From source

You can compile it yourself by running (Golang is required):

```
go get github.com/BenjaminNavarro/pidws
```

Then you can move or link the `pidws` executable to a suitable location, e.g `mv pidws /usr/local/bin`.

## From binaries

You can get the latest's release Linux binaries from [here](https://github.com/BenjaminNavarro/pidws/releases). Extract the archive's content, e.g with `tar -xf pidws.tar.gz`, and move or link the `pidws` executable to a suitable location, e.g `mv pidws /usr/local/bin`.


# Usage

Here are the flags you can pass to pidws (you can either use the `-` or the `--` prefix):
```
  -cd string
    	Change the current directory to the specified workspace (a new shell will be created).
  -create string
    	Create a new workspace in the current directory. To specify its name, see -name
  -default string
    	Set the workspace that will receive workspace commands when none is specified (see -use).
  -name string
    	To be used together with -create-ws, -register-workspace or -register-repository to specify the name of the new workspace or repository.
  -register-repository string
    	Add a new repository in the configuration file.
  -register-workspace string
    	Add an already existing workspace in the configuration file.
  -remove-repository string
    	Remove a repository from the list of the known ones.
  -remove-workspace string
    	Remove a workspace from the list of the known ones. Files won't be deleted.
  -use string
    	If set, it will override the default workspace for the current commands.
```
All other commands will be passed to the active PID workspace.

All the settings are stored in your home folder under the `.pidws.yaml` file. You can edit it manually but be careful of what you do.
