# Config Validation Example

This example demonstrates how **goStartyUpy**'s `configcheck` package
validates configuration loaded by **goConfy**.

## What it shows

- Loading a YAML configuration file with goConfy.
- Running `configcheck.RunStartupCheck` to validate required fields.
- Clear error diagnostics when required configuration keys are missing.
- A **full reference configuration** covering every supported option.

## Files

| File                           | Purpose                                                         |
|--------------------------------|-----------------------------------------------------------------|
| `config.full.example.yml`     | **Full reference** – every supported configuration option       |
| `config.minimal.yml`          | Smallest valid configuration (required fields only)             |
| `config.yml`                  | Typical valid configuration                                     |
| `config_invalid.yml`          | Intentionally broken configuration for negative testing         |
| `main.go`                     | Runnable example application                                    |
| `config_validation_test.go`   | Tests for all configuration variants                            |

## Full example configuration

The file **`config.full.example.yml`** is the canonical reference for all
configuration options supported by goStartyUpy. It contains every
user-relevant setting with explanatory comments and realistic values.

Use it as a starting point when integrating goStartyUpy into your project:

```bash
cp examples/config-validation/config.full.example.yml config.yml
```

Edit the values to match your environment and remove any sections you do
not need.

## Running the example

From the repository root:

```bash
go run ./examples/config-validation/
```

The application loads `config.full.example.yml`, validates it, and prints
the startup banner followed by the loaded values.

## Running the tests

```bash
go test ./examples/config-validation/ -v
```

The test suite validates:

- **Minimal config** passes validation (only required fields).
- **Full example config** passes validation and has all key fields populated.
- **Valid config** (`config.yml`) passes validation.
- **Invalid config** fails with the expected missing keys.

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
