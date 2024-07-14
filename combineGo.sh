#!/bin/bash

# Check if a path argument is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <path>"
    exit 1
fi

# Starting path from command line argument
start_path="$1"

# Check if the provided path exists
if [ ! -d "$start_path" ]; then
    echo "Error: The specified path does not exist or is not a directory."
    exit 1
fi

# Output file
output_file="combined_golang_files.txt"

# Clear the output file if it already exists
> "$output_file"

# Function to process files
process_files() {
    for file in "$1"/*; do
        if [ -d "$file" ]; then
            # If it's a directory, recurse into it
            process_files "$file"
        elif [ -f "$file" ] && [[ "$file" == *.go ]]; then
            # If it's a .go file, add its content to the output file
            echo -e "\n============== $file ==============" >> "$output_file"
            cat "$file" >> "$output_file"
        fi
    done
}

# Start processing from the specified directory
process_files "$start_path"

echo "All Golang files from $start_path have been combined into $output_file"