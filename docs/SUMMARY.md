# Documentation Summary

This Hugo documentation site provides comprehensive documentation for lanup.

## What's Included

### Pages

1. **Home** (`content/_index.md`)
   - Project overview
   - Quick start
   - Feature highlights

2. **Installation** (`content/docs/installation.md`)
   - Installation methods
   - Platform-specific instructions
   - Verification steps

3. **Getting Started** (`content/docs/getting-started.md`)
   - Step-by-step tutorial
   - Configuration basics
   - First run guide

4. **Commands Reference** (`content/docs/commands.md`)
   - Complete command documentation
   - All flags and options
   - Usage examples

5. **Configuration** (`content/docs/configuration.md`)
   - Project configuration
   - Global configuration
   - Environment file format
   - Configuration examples

6. **Use Cases** (`content/docs/use-cases.md`)
   - Mobile app development
   - Supabase integration
   - Docker development
   - Team collaboration
   - Multi-device testing

7. **Troubleshooting** (`content/docs/troubleshooting.md`)
   - Common problems and solutions
   - Network issues
   - Permission errors
   - Service accessibility

8. **Development Guide** (`content/docs/development.md`)
   - Contributing guidelines
   - Project structure
   - Building and testing
   - Release process

## Setup Instructions

### Quick Setup

```bash
cd docs
./setup.sh
hugo server
```

### Manual Setup

1. Install Hugo (v0.100.0+)
2. Install theme:
   ```bash
   cd docs
   git submodule add https://github.com/alex-shpak/hugo-book themes/book
   ```
3. Run server:
   ```bash
   hugo server
   ```

## Deployment Options

### GitHub Pages

The site is configured to deploy automatically to GitHub Pages via GitHub Actions.

Workflow file: `.github/workflows/deploy.yml`

### Netlify

Configuration file: `netlify.toml`

Simply connect your repository to Netlify.

### Manual Build

```bash
cd docs
hugo --minify
```

Output will be in `docs/public/`

## Theme

Uses [Hugo Book](https://github.com/alex-shpak/hugo-book) theme for clean, documentation-focused design.

## Customization

- **Site config:** `config.toml`
- **Content:** `content/` directory
- **Theme overrides:** `layouts/` directory
- **Static assets:** `static/` directory

## Maintenance

### Adding New Pages

1. Create markdown file in `content/docs/`
2. Add frontmatter with title and weight
3. Write content
4. Test with `hugo server`

### Updating Content

1. Edit markdown files in `content/docs/`
2. Test locally
3. Commit and push (auto-deploys if configured)

### Theme Updates

```bash
cd docs/themes/book
git pull origin main
```

## Features

- ✅ Clean, modern design
- ✅ Mobile responsive
- ✅ Search functionality
- ✅ Syntax highlighting
- ✅ Dark mode support
- ✅ Fast static site
- ✅ SEO optimized
- ✅ Easy to maintain
