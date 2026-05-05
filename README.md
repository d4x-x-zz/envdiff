# envdiff

> CLI tool to diff `.env` files across environments and flag missing or mismatched keys.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff
go build -o envdiff .
```

---

## Usage

```bash
envdiff [flags] <base-env> <compare-env> [additional-envs...]
```

### Example

```bash
envdiff .env.example .env.production
```

**Sample output:**

```
MISSING in .env.production:
  - DATABASE_URL
  - REDIS_HOST

MISMATCHED keys (present in both, different values):
  - APP_ENV   [.env.example: "development"]  [.env.production: "production"]

✔ All other keys match.
```

### Flags

| Flag | Description |
|------|-------------|
| `--strict` | Exit with non-zero status if any differences are found |
| `--ignore KEY` | Ignore a specific key during comparison |
| `--json` | Output results in JSON format |

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)