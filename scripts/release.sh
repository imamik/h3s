#!/bin/bash

# Exit on error
set -e

# Add all changes
git add .

# Get commit message from user
echo "Enter commit message:"
read commit_msg

# Commit changes
git commit -m "$commit_msg"

# Get the latest tag
latest_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

# Extract version numbers
IFS='.' read -r major minor patch <<< "$(echo "$latest_tag" | sed 's/v//')"

# Increment patch version
new_patch=$((patch + 1))
new_tag="v$major.$minor.$new_patch"

# Create new tag
git tag "$new_tag"

# Push changes and tags
git push origin main --tags

echo "Successfully:"
echo "- Committed changes with message: $commit_msg"
echo "- Created new tag: $new_tag"
echo "- Pushed changes and tags to remote" 