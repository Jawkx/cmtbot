package main

import (
	"bytes"
	"os/exec"
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

	output := outb.String()

	return output, nil
}
