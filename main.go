package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	branch1   string
	branch2   string
	format    string
	separator string
	noColor   bool
)

func init() {
	flag.StringVar(&branch1, "branch1", "", "First branch to compare (required)")
	flag.StringVar(&branch2, "branch2", "", "Second branch to compare (required)")
	flag.StringVar(&format, "format", "text", "Output format: text, json, custom")
	flag.StringVar(&separator, "separator", "\n", "Separator for custom output format")
	flag.BoolVar(&noColor, "no-color", false, "Disable color output")
}

func main() {
	flag.Parse()

	if branch1 == "" || branch2 == "" {
		outputError("Both --branch1 and --branch2 are required")
	}

	// Configure color output
	if format != "text" {
		color.NoColor = true
	} else {
		color.NoColor = noColor
	}

	totalLines, err := getTotalLines()
	if err != nil {
		outputError(fmt.Sprintf("Error calculating total lines: %v", err))
	}

	changedLines, err := getChangedLines(branch1, branch2)
	if err != nil {
		outputError(fmt.Sprintf("Error calculating changed lines: %v", err))
	}

	if changedLines == 0 && format == "text" {
		color.Green("No changes between the branches.")
		return
	}

	var percentageChange float64
	if totalLines > 0 {
		percentageChange = float64(changedLines) / float64(totalLines) * 100
	}

	switch format {
	case "text":
		fmt.Printf("Total lines in repository: %s\n", color.CyanString("%d", totalLines))
		fmt.Printf("Lines changed between %s and %s: %s\n", branch1, branch2, color.YellowString("%d", changedLines))
		fmt.Printf("Percentage of change: %s\n", color.GreenString("%.2f%%", percentageChange))
	case "json":
		out := struct {
			TotalLines   int     `json:"totalLines"`
			ChangedLines int     `json:"changedLines"`
			Percentage   float64 `json:"percentage"`
			Branch1      string  `json:"branch1"`
			Branch2      string  `json:"branch2"`
		}{
			TotalLines:   totalLines,
			ChangedLines: changedLines,
			Percentage:   percentageChange,
			Branch1:      branch1,
			Branch2:      branch2,
		}
		jsonData, err := json.Marshal(out)
		if err != nil {
			outputError(fmt.Sprintf("Error generating JSON: %v", err))
		}
		fmt.Println(string(jsonData))
	case "custom":
		parts := []string{
			strconv.Itoa(totalLines),
			strconv.Itoa(changedLines),
			fmt.Sprintf("%.2f", percentageChange),
			branch1,
			branch2,
		}
		fmt.Println(strings.Join(parts, separator))
	default:
		outputError("Invalid output format specified")
	}
}

func outputError(message string) {
	switch format {
	case "json":
		errJSON := struct {
			Error string `json:"error"`
		}{
			Error: message,
		}
		jsonData, _ := json.Marshal(errJSON)
		fmt.Println(string(jsonData))
	default:
		color.Red(message)
	}
	os.Exit(1)
}

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

func getChangedLines(branch1, branch2 string) (int, error) {
	out, err := exec.Command("git", "diff", "--stat", branch1, branch2).Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 {
		return 0, nil
	}
	lastLine := lines[len(lines)-1]

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
