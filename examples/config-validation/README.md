# Config Validation Example

This example demonstrates how **goStartyUpy**'s `configcheck` package
validates configuration loaded by **goConfy**.

## What it shows

- Loading a YAML configuration file with goConfy.
- Running `configcheck.RunStartupCheck` to validate required fields.
- Clear error diagnostics when required configuration keys are missing.

## Files

| File                  | Purpose                                   |
|-----------------------|-------------------------------------------|
| `config.yml`          | Valid configuration with all required keys |
| `config_invalid.yml`  | Broken configuration with missing keys    |
| `main.go`             | Runnable example application              |
| `config_validation_test.go` | Tests for both valid and invalid configs |

## Running the example

From the repository root:

```bash
go run ./examples/config-validation/
```

The application loads `config.yml`, validates it, and prints the startup
banner followed by the loaded values.

## Running the tests

```bash
go test ./examples/config-validation/ -v
```

## What happens with invalid configuration

When required fields are missing, `configcheck` produces a diagnostic like:

```
Config validation failed:

Missing configuration keys:
  - database.port
  - database.password
  - redis.address
```

The error message lists each missing key using its dot-separated YAML path,
making it easy to find and fix the problem.
