package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func main() {
	if len(os.Args) != 3 {
		color.Red("Usage: %s <branch1> <branch2>", os.Args[0])
		os.Exit(1)
	}

	branch1 := os.Args[1]
	branch2 := os.Args[2]

	totalLines, err := getTotalLines()
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

	fmt.Printf("Total lines in repository: %s\n", color.CyanString("%d", totalLines))
	fmt.Printf("Lines changed between %s and %s: %s\n", branch1, branch2, color.YellowString("%d", changedLines))
	fmt.Printf("Percentage of change: %s\n", color.GreenString("%.2f%%", percentageChange))
}

// getTotalLines returns the total number of lines in the repository
func getTotalLines() (int, error) {
	out, err := exec.Command("git", "ls-files").Output()
	if err != nil {
		return 0, err
	}

	files := strings.Split(strings.TrimSpace(string(out)), "\n")
	totalLines := 0

	for _, file := range files {
		out, err := exec.Command("wc", "-l", file).Output()
		if err != nil {
			return 0, err
		}
		lines, err := strconv.Atoi(strings.Fields(string(out))[0])
		if err != nil {
			return 0, err
		}
		totalLines += lines
	}

	return totalLines, nil
}

// getChangedLines returns the number of lines changed between two branches
func getChangedLines(branch1, branch2 string) (int, error) {
	out, err := exec.Command("git", "diff", "--stat", branch1, branch2).Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	lastLine := lines[len(lines)-1]

	// Extract the number of changed lines from the last line of git diff --stat
	fields := strings.Fields(lastLine)
	if len(fields) < 1 {
		return 0, nil
	}

	changedLines, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, err
	}

	return changedLines, nil
}
