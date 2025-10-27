# lanup Documentation

This directory contains the Hugo-based documentation for lanup.

## Prerequisites

- [Hugo](https://gohugo.io/) v0.100.0 or higher

## Installation

### Install Hugo

**macOS:**
```bash
brew install hugo
```

**Linux:**
```bash
# Debian/Ubuntu
sudo apt-get install hugo

# Arch Linux
sudo pacman -S hugo
```

**Windows:**
```bash
choco install hugo-extended
```

Or download from [Hugo Releases](https://github.com/gohugoio/hugo/releases).

### Install Theme

```bash
cd docs
git submodule add https://github.com/alex-shpak/hugo-book themes/book
```

Or if you prefer to download manually:
```bash
cd docs/themes
git clone https://github.com/alex-shpak/hugo-book book
```

## Development

### Run Local Server

```bash
cd docs
hugo server
```

Visit `http://localhost:1313` to view the documentation.

### Build for Production

```bash
cd docs
hugo
```

The generated site will be in the `public/` directory.

## Structure

```
docs/
├── config.toml           # Hugo configuration
├── content/              # Documentation content
│   ├── _index.md        # Homepage
│   └── docs/            # Documentation pages
│       ├── _index.md
│       ├── installation.md
│       ├── getting-started.md
│       ├── commands.md
│       ├── configuration.md
│       ├── use-cases.md
│       ├── troubleshooting.md
│       └── development.md
├── themes/              # Hugo themes
│   └── book/           # Book theme
└── public/             # Generated site (gitignored)
```

## Adding New Pages

1. Create a new markdown file in `content/docs/`:
   ```bash
   hugo new docs/my-page.md
   ```

2. Add frontmatter:
   ```yaml
   ---
   title: "My Page"
   weight: 10
   ---
   ```

3. Write content in markdown

4. Test locally:
   ```bash
   hugo server
   ```

## Deployment

### GitHub Pages

1. Build the site:
   ```bash
   hugo
   ```

2. Deploy the `public/` directory to GitHub Pages

### Netlify

1. Connect your repository to Netlify
2. Set build command: `cd docs && hugo`
3. Set publish directory: `docs/public`

### Vercel

1. Connect your repository to Vercel
2. Set build command: `cd docs && hugo`
3. Set output directory: `docs/public`

## Theme Customization

The documentation uses the [Hugo Book](https://github.com/alex-shpak/hugo-book) theme.

To customize:
- Edit `config.toml` for site-wide settings
- Override theme templates in `layouts/`
- Add custom CSS in `static/css/`

## Contributing

When adding documentation:
- Use clear, concise language
- Include code examples
- Add screenshots where helpful
- Test all commands and examples
- Update the navigation if needed
