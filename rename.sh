#!/bin/bash

# Treat unset variables as an error
set -u
# Exit on error
set -e 

# Obtain repository owner and name from argument
owner=$(echo "$1" | cut -d'/' -f1)
name=$(echo "$1" | cut -d'/' -f2)

# Check if owner or name are empty
if [ -z "$owner" ] || [ -z "$name" ]; then
    echo "Invalid repository name '$1'"
    exit 1
fi

# Lowercase and uppercase
name_low=$(echo "$name" | tr '[:upper:]' '[:lower:]')
name_upp=$(echo "$name" | tr '[:lower:]' '[:upper:]')

# Rename files and folders
mv cmd/goobar "cmd/$name_low"
mv goobar.go "$name_low.go"

# Replace file contents
find . -type f -exec sed -i "s/goobar/$name_low/g" {} \;
find . -type f -exec sed -i "s/GOOBAR/$name_upp/g" {} \;
find . -type f -exec sed -i "s/igolaizola/$owner/g" {} \;

# Override README.md
echo "# $name" > README.md

# Remove this script folder
rm -rf rename.sh

# Remove .git folder
rm -rf .git
