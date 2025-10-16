package main

import (
	"fmt"
	"os/exec"
)

// getGitDiff executes git diff command and returns the output
func getGitDiff(filePath string) (string, error) {
	var cmd *exec.Cmd

	if filePath == "." {
		// Get diff for all files
		cmd = exec.Command("git", "diff")
	} else {
		// Get diff for specific file
		cmd = exec.Command("git", "diff", filePath)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git diff failed: %v\nOutput: %s", err, string(output))
	}

	return string(output), nil
}
