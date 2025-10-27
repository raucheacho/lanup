#!/bin/bash

# Setup script for lanup documentation

set -e

echo "🚀 Configuration de la documentation lanup..."
echo ""

# Check if Hugo is installed
if ! command -v hugo &> /dev/null; then
    echo "❌ Hugo n'est pas installé"
    echo ""
    echo "Installez Hugo d'abord :"
    echo "  macOS:   brew install hugo"
    echo "  Linux:   sudo apt-get install hugo"
    echo "  Windows: choco install hugo-extended"
    echo ""
    echo "Ou visitez: https://gohugo.io/installation/"
    exit 1
fi

echo "✓ Hugo est installé ($(hugo version | head -n 1))"

# Check if theme exists
if [ ! -d "themes/book" ]; then
    echo "📦 Installation du thème Hugo Book..."
    
    # Create themes directory
    mkdir -p themes
    
    # Clone the theme
    git clone https://github.com/alex-shpak/hugo-book themes/book
    
    echo "✓ Thème installé"
else
    echo "✓ Thème déjà installé"
fi

echo ""
echo "✅ Configuration terminée !"
echo ""
echo "Pour démarrer le serveur de développement :"
echo "  hugo server"
echo ""
echo "Puis visitez: http://localhost:1313"
