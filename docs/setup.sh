#!/bin/bash

# Setup script for lanup documentation

set -e

echo "Setting up lanup documentation..."

# Check if Hugo is installed
if ! command -v hugo &> /dev/null; then
    echo "❌ Hugo is not installed"
    echo ""
    echo "Please install Hugo first:"
    echo "  macOS:   brew install hugo"
    echo "  Linux:   sudo apt-get install hugo"
    echo "  Windows: choco install hugo-extended"
    echo ""
    echo "Or visit: https://gohugo.io/installation/"
    exit 1
fi

echo "✓ Hugo is installed ($(hugo version))"

# Check if theme exists
if [ ! -d "themes/book" ]; then
    echo "Installing Hugo Book theme..."
    
    # Try git submodule first
    if git rev-parse --git-dir > /dev/null 2>&1; then
        git submodule add https://github.com/alex-shpak/hugo-book themes/book || true
        git submodule update --init --recursive
    else
        # Clone directly if not in a git repo
        mkdir -p themes
        git clone https://github.com/alex-shpak/hugo-book themes/book
    fi
    
    echo "✓ Theme installed"
else
    echo "✓ Theme already installed"
fi

echo ""
echo "✅ Setup complete!"
echo ""
echo "To start the development server:"
echo "  cd docs"
echo "  hugo server"
echo ""
echo "Then visit: http://localhost:1313"
