# Contributing to NYXORA

First off, thank you for considering contributing to NYXORA! We welcome contributions from everyone.

## Code of Conduct

This project adheres to a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before submitting a bug report:
- Check the [issues](https://github.com/nyxora-network/nyxora/issues) to see if it's already reported
- Collect information: OS version, Go version, steps to reproduce, error output

**Submit a bug report** by opening a [new issue](https://github.com/nyxora-network/nyxora/issues/new?template=bug_report.md).

### Suggesting Features

Open a [feature request](https://github.com/nyxora-network/nyxora/issues/new?template=feature_request.md) describing:
- The problem you're solving
- How you envision the solution
- Any alternatives considered

### Adding a New Transport

1. Create `internal/transport/<name>.go` implementing the `Transport` interface
2. Register it in `internal/transport/registry.go`
3. Create `tunnels/<name>/` with install scripts and manifest
4. Add scoring weights in `internal/transport/scoring.go`
5. Write tests and run `make test`

### Improving the TUI

The interactive TUI is in `internal/interactive/` and uses:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — styling
- [Bubbles](https://github.com/charmbracelet/bubbles) — components (textinput, spinner, progress)

Theme colors use Catppuccin TrueColor hex values in `theme.go`.

## Pull Request Process

1. Fork the repo and create your branch from `main`
2. Run `make test` and `make vet` — all must pass
3. Add tests for new functionality
4. Update documentation if needed
5. Make sure your code follows existing conventions
6. Submit the PR with a clear description

### Commit Style

Use conventional commit messages:
- `feat:` — new feature
- `fix:` — bug fix
- `refactor:` — code change without fix/feature
- `docs:` — documentation only
- `test:` — test additions/fixes
- `style:` — formatting, style changes
- `chore:` — maintenance, dependencies

## Development Setup

```bash
# Fork and clone
git clone https://github.com/YOUR_USERNAME/nyxora.git
cd nyxora

# Add upstream remote
git remote add upstream https://github.com/nyxora-network/nyxora.git

# Create feature branch
git checkout -b feat/your-feature

# Make changes, then:
make test
make vet
make build

# Commit and push
git commit -m "feat: add your feature"
git push origin feat/your-feature
```

## Questions?

Open a [Discussion](https://github.com/nyxora-network/nyxora/discussions) or join our [Telegram](https://t.me/NyxoraCore).
