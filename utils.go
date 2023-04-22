package phylum

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func ExtractRemediation(issue IssuesListItem) (string, error) {
	var result string
	desc := issue.Description

	recPat := regexp.MustCompile(`### Recommendation`)

	if len(desc) < 1 {
		// no description
		return "", errors.New("No description found")
	}

	start := recPat.FindStringIndex(desc) //[0]
	if len(start) > 0 {
		offset := start[0]
		result = desc[offset:]
	}
	//result = desc[start:]

	return result, nil
}

func GetTokenFromCLI() (string, error) {
	var stdErrBytes bytes.Buffer
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user homedir")
	}
	_ = home
	pathThing := filepath.Join(home, ".config/phylum/fb-settings.yaml")
	var phylumTokenArgs = []string{"-c", pathThing, "auth", "token"}
	phylumTokenCmd := exec.Command("phylum", phylumTokenArgs...)
	phylumTokenCmd.Stderr = &stdErrBytes
	output, err := phylumTokenCmd.Output()
	stdErrString := stdErrBytes.String()
	if err != nil {
		fmt.Printf("error running phylum auth token: %v\n", err)
		fmt.Printf("stderr: %v\n", stdErrString)
	}
	return strings.TrimSuffix(string(output), "\n"), nil
}

func GetAuthUser() (string, error) {
	var stdErrBytes bytes.Buffer
	//var phylumTokenArgs = []string{"-c", "$HOME/.config/phylum/fb-settings.yaml", "auth", "status"}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user homedir")
	}
	_ = home
	pathThing := filepath.Join(home, ".config/phylum/fb-settings.yaml")
	phylumTokenCmd := exec.Command("phylum", "-c", pathThing, "auth", "status")
	phylumTokenCmd.Stderr = &stdErrBytes
	output, err := phylumTokenCmd.Output()
	stdErrString := stdErrBytes.String()
	if err != nil {
		fmt.Printf("error running phylum auth status: %v\n", err)
		fmt.Printf("stderr: %v\n", stdErrString)
	}
	return strings.TrimSuffix(string(output), "\n"), nil
}
