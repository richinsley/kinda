package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-git/go-git/v5"
	kinda "github.com/richinsley/kinda/pkg"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		// Handle the error if there is one
		fmt.Printf("Error getting current working directory: %v\n", err)
		return
	}

	// Specify the binary folder to place micromamba in
	rootDirectory := filepath.Join(cwd, "micromamba")
	fmt.Println("Creating Kinda repo at: ", rootDirectory)
	version := "3.10"
	env, err := kinda.CreateEnvironment("myenv"+version, rootDirectory, version, "conda-forge", kinda.ShowVerbose)
	if err != nil {
		fmt.Printf("Error creating environment: %v\n", err)
		return
	}
	fmt.Printf("Created environment: %s\n", env.Name)

	// Clone the given repository to the given directory
	fmt.Printf("git clone https://github.com/go-git/go-git")

	comfyFolder := filepath.Join(cwd, "comfyui")
	repo, err := git.PlainClone(comfyFolder, false, &git.CloneOptions{
		URL:      "https://github.com/comfyanonymous/ComfyUI.git",
		Progress: os.Stdout,
	})

	if err != nil && err.Error() != "repository already exists" {
		fmt.Printf("Error cloning: %v\n", err)
		return
	}

	if repo != nil {
		if runtime.GOOS == "windows" {
			// pre-install some required packages for cuda
			packages := []string{
				"torch",
				"torchvision",
				"torchaudio",
			}
			err = env.PipInstallPackages(packages, "", "https://download.pytorch.org/whl/cu121", false)
			if err != nil {
				fmt.Printf("Error installing requirements: %v\n", err)
				return
			}
		}
		// install the pip requirements
		requirementsPath := filepath.Join(comfyFolder, "requirements.txt")
		err = env.PipInstallRequirmements(requirementsPath, kinda.ShowVerbose)
		if err != nil {
			fmt.Printf("Error installing requirements: %v\n", err)
			return
		}
	}

	// run the python script
	// python main.py --highvram --listen
	env.RunPythonScriptFromFile(filepath.Join(comfyFolder, "main.py"), "--highvram", "--listen")
}
