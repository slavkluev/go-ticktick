# Contributing to go-ticktick

We welcome contributions! Whether it's a bug report, feature suggestion, or pull request — all input is appreciated.

## Reporting Issues

- Search [existing issues](https://github.com/slavkluev/go-ticktick/issues) before opening a new one.
- Include Go version, OS, and a minimal code example to reproduce the problem.
- **Security vulnerabilities** must be reported by email to <hi@kliuev.dev>, not as public issues.

## Submitting a Pull Request

1. Fork the repository and create a branch from `main`.
2. Add or update tests for your changes.
3. Make sure everything passes:
   ```bash
   make test
   make lint
   ```
4. Keep commits focused — one logical change per commit.
5. Open a pull request with a clear description of what changed and why.

## Development

Requires Go 1.25+ and Docker (for linting).

```bash
# Run all tests
make test

# Run a single test
go test -run TestName ./...

# Lint (runs golangci-lint via Docker)
make lint
```

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).