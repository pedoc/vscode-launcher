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

func main() {
	var userDataDirDefaultValue = "private-data"
	var extensionsDirDefaultValue = "private-extensions"

	wd, _ := os.Getwd()
	log.Println("Working Directory:", wd)
	codeBin := "Code.exe"
	if runtime.GOOS != "windows" {
		codeBin = "code"
	}
	codePath := filepath.Join(wd, codeBin)
	if _, err := os.Stat(codePath); os.IsNotExist(err) {
		log.Fatalf("Not found %s", codeBin)
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
		allArgs := append([]string{"/C", "start", codePath}, args...)
		cmd := exec.Command("cmd", allArgs...)
		cmd.Dir = wd
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	} else {
		allArgs := append([]string{codePath}, args...)
		cmd := exec.Command("sh", "-c", fmt.Sprintf("%s", filepath.Join(allArgs...)))
		cmd.Dir = wd
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
}
