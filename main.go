package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	GitCommit  string
	BuildTime  string
	GitBranch  string
	GitVersion string
)

// fileExists checks if a file or directory exists at the given path.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

// getCodeBin returns the name of the VSCode binary based on the operating system.
// It checks for the standard VSCode binary and, if not found, checks for the Insiders version.
func getCodeBin() string {
	var codeBin string
	if runtime.GOOS == "windows" {
		codeBin = "Code.exe"
		if !fileExists(codeBin) {
			codeBin = "Code - Insiders.exe"
		}
	} else {
		codeBin = "code"
		if !fileExists(codeBin) {
			codeBin = "code-insiders"
		}
	}
	return codeBin
}

func main() {
	var userDataDirDefaultValue = "private-data"
	var extensionsDirDefaultValue = "private-extensions"
	var codeBin = getCodeBin()
	wd, _ := os.Getwd()
	log.Println("Working Directory:", wd)

	codePath := filepath.Join(wd, codeBin)
	if !fileExists(codePath) {
		log.Fatalf("Not found %s", codePath)
		return
	}
	var launchArgs []string
	launchArgs = append(launchArgs, fmt.Sprintf("--user-data-dir=%s", userDataDirDefaultValue))
	launchArgs = append(launchArgs, fmt.Sprintf("--extensions-dir=%s", extensionsDirDefaultValue))
	if runtime.GOOS == "windows" && len(os.Args) == 1 {
		launchVSCode(wd, codePath, launchArgs...)
		return
	}

	var userDataDir string
	var extensionsDir string
	var rootCmd = &cobra.Command{
		Use:   "vscode_launcher [projectDir]",
		Short: "Start VSCode with optional user-data-dir and extensions-dir",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			launchArgs = launchArgs[:0]
			if len(args) > 0 {
				launchArgs = append(launchArgs, args[0])
			}
			if userDataDir != "" {
				launchArgs = append(launchArgs, fmt.Sprintf("--user-data-dir=%s", userDataDir))
			}
			if extensionsDir != "" {
				launchArgs = append(launchArgs, fmt.Sprintf("--extensions-dir=%s", extensionsDir))
			}
			return launchVSCode(wd, codePath, launchArgs...)
		},
	}

	rootCmd.Flags().StringVarP(&userDataDir, "user-data-dir", "u", userDataDirDefaultValue, "VSCode user data directory")
	rootCmd.Flags().StringVarP(&extensionsDir, "extensions-dir", "e", extensionsDirDefaultValue, "VSCode extensions directory")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func launchVSCode(wd string, codePath string, args ...string) error {
	log.Println("Launching VSCode:", codePath)
	if runtime.GOOS == "windows" {
		allArgs := append([]string{"/C", "start", "", codePath}, args...)
		log.Println("Command Arguments:", allArgs)
		cmd := exec.Command("cmd", allArgs...)
		cmd.Dir = wd
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	} else {
		allArgs := append([]string{codePath}, args...)
		log.Println("Command Arguments:", allArgs)
		cmd := exec.Command("sh", "-c", fmt.Sprintf("%s", filepath.Join(allArgs...)))
		cmd.Dir = wd
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
}
