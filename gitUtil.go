package main

import (
	"bytes"
	"os/exec"
	"strings"
)

func getStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func getStagedFiles() (string, error) {
	cmd := exec.Command("git", "diff", "--staged", "--name-status")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func commitChanges(message string) (string, error) {
	cmd := exec.Command("git", "commit", "-m", message)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()

	if err != nil {
		return "", err
	}

	// Extract commit hash from the output
	output := outb.String()
	if strings.Contains(output, "nothing to commit") {
		return "", nil // Or an error if you prefer
	}

	commitHashLine := strings.SplitN(output, "\n", 2)[0]
	commitHash := strings.SplitN(commitHashLine, " ", 2)[1]
	return commitHash, nil
}
