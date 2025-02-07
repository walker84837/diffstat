package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func main() {
	if len(os.Args) != 3 {
		color.Red("Usage: diffstat <main branch> <feature branch>")
		os.Exit(1)
	}

	branch1 := os.Args[1]
	branch2 := os.Args[2]

	// We now count the total lines in the target branch (branch2) rather than the working directory.
	totalLines, err := getTotalLines(branch2)
	if err != nil {
		color.Red("Error calculating total lines: %v", err)
		os.Exit(1)
	}

	changedLines, err := getChangedLines(branch1, branch2)
	if err != nil {
		color.Red("Error calculating changed lines: %v", err)
		os.Exit(1)
	}

	if changedLines == 0 {
		color.Green("No changes between the branches.")
		return
	}

	percentageChange := float64(changedLines) / float64(totalLines) * 100

	fmt.Printf("Total lines in branch %s: %s\n", branch2, color.CyanString("%d", totalLines))
	fmt.Printf("Lines changed between %s and %s: %s\n", branch1, branch2, color.YellowString("%d", changedLines))
	fmt.Printf("Percentage of change: %s\n", color.GreenString("%.2f%%", percentageChange))
}

// getTotalLines returns the total number of lines in all files of the given branch.
func getTotalLines(branch string) (int, error) {
	// List all files tracked in the given branch.
	out, err := exec.Command("git", "ls-tree", "-r", "--name-only", branch).Output()
	if err != nil {
		return 0, err
	}
	files := strings.Split(strings.TrimSpace(string(out)), "\n")
	totalLines := 0

	for _, file := range files {
		// Get file content from the branch
		content, err := exec.Command("git", "show", branch+":"+file).Output()
		if err != nil {
			// If file is binary or cannot be retrieved, skip it
			continue
		}
		// Count the number of lines
		lines := strings.Split(string(content), "\n")
		totalLines += len(lines)
	}

	return totalLines, nil
}

// getChangedLines calculates the number of changed lines between branch1 and branch2.
func getChangedLines(branch1, branch2 string) (int, error) {
	out, err := exec.Command("git", "diff", "--numstat", branch1, branch2).Output()
	if err != nil {
		return 0, err
	}

	totalChanges := 0
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		addedStr, deletedStr := fields[0], fields[1]

		if addedStr == "-" || deletedStr == "-" {
			// Handle binary or new/deleted files
			totalChanges += estimateBinaryOrNewFileChange(branch1, branch2, fields[2])
		} else {
			// Convert to integers
			a, err := strconv.Atoi(addedStr)
			if err != nil {
				return 0, fmt.Errorf("failed to parse added lines: %s", addedStr)
			}
			d, err := strconv.Atoi(deletedStr)
			if err != nil {
				return 0, fmt.Errorf("failed to parse deleted lines: %s", deletedStr)
			}
			totalChanges += a + d
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return totalChanges, nil
}

// estimateBinaryOrNewFileChange estimates line changes for binary files or newly created/deleted files.
func estimateBinaryOrNewFileChange(branch1, branch2, file string) int {
	// Get file size in bytes from both branches
	size1 := getFileSize(branch1, file)
	size2 := getFileSize(branch2, file)

	sizeChange := abs(size2 - size1)

	// Estimate number of lines based on size change
	estimatedLines := sizeChange / 100 // Rough estimate: 100 bytes per line
	if estimatedLines == 0 {
		estimatedLines = 1 // Minimum change of 1 line
	}

	return estimatedLines
}

// getFileSize returns the file size in bytes for a given branch and file.
func getFileSize(branch, file string) int {
	out, err := exec.Command("git", "cat-file", "-s", branch+":"+file).Output()
	if err != nil {
		return 0 // Assume deleted or new file
	}
	size, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0
	}
	return size
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
