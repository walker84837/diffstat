#!/bin/bash

# Get the total number of lines in the repository (including all files)
total_lines=$(git ls-files | xargs wc -l | tail -n 1 | awk '{print $1}')

# Get the number of lines changed between the two branches
changed_lines=$(git diff --stat <branch1> <branch2> | awk '{s+=$1} END {print s}')

# Calculate the percentage of change
if [ -z "$changed_lines" ]; then
  echo "No changes between the branches."
  exit 0
fi

percentage_change=$(echo "scale=2; ($changed_lines / $total_lines) * 100" | bc)

echo "Percentage of change: $percentage_change%"
