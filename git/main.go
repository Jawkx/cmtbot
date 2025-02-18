package git

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/Jawkx/cmtbot/model"
)

func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

type StagedFile struct {
	Path   string
	Status model.FileStatus
}

func GetStagedFiles() ([]StagedFile, error) {
	var result []StagedFile
	cmd := exec.Command("git", "diff", "--staged", "--name-status")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return result, err
	}

	fileStrings := out.String()

	if fileStrings == "" {
		return result, nil
	}

	fileStringSlice := strings.Split(fileStrings, "\n")

	for _, fileString := range fileStringSlice {

		parts := strings.SplitN(fileString, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		statusString, filePath := parts[0], parts[1]

		var status model.FileStatus

		switch statusString {
		case model.StatusAdded.String():
			status = model.StatusAdded
		case model.StatusModified.String():
			status = model.StatusModified
		case model.StatusDeleted.String():
			status = model.StatusDeleted
		case model.StatusRenamed.String():
			status = model.StatusRenamed
		case model.StatusCopied.String():
			status = model.StatusCopied
		case model.StatusUnmodified.String():
			status = model.StatusUnmodified
		default:
			status = model.StatusUnknown
		}

		result = append(result, StagedFile{
			Path:   filePath,
			Status: status,
		})
	}

	return result, err
}

func CommitChanges(message string) (string, error) {
	cmd := exec.Command("git", "commit", "-m", message)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()

	if err != nil {
		return "", err
	}

	output := strings.TrimSpace(outb.String())

	return output, nil
}
