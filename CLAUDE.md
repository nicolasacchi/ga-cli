# CLAUDE.md — ga-cli

Go CLI for Google Analytics 4. Single binary, JSON output, service account auth.

**APIs**: GA4 Data API v1beta (reports, realtime) + Admin API v1beta (accounts, properties, metadata).

## Authentication

1. `GOOGLE_APPLICATION_CREDENTIALS` env var — path to service account JSON key file
2. Falls back to Application Default Credentials (`gcloud auth application-default login`)

Service account needs **Viewer** role on the GA4 property.

## Commands

### accounts

```bash
ga-cli accounts list        # List all GA4 accounts and property summaries
```

### properties

```bash
ga-cli properties get 257850630    # Get property details (timezone, currency, type)
```

### report

```bash
ga-cli report run 257850630 --dimensions date --metrics sessions,activeUsers
ga-cli report run 257850630 --dimensions date,country --metrics sessions --start-date 2026-02-16 --end-date 2026-02-22
ga-cli report run 257850630 --dimensions pagePath --metrics screenPageViews --limit 20 --offset 0
```

| Flag | Default | Description |
|------|---------|-------------|
| `--start-date` | 7 days ago | YYYY-MM-DD |
| `--end-date` | yesterday | YYYY-MM-DD |
| `--dimensions` | (none) | Comma-separated (e.g., `date,country,pagePath`) |
| `--metrics` | `sessions,activeUsers` | Comma-separated |
| `--limit` | 100 | Max rows |
| `--offset` | 0 | Row offset for pagination |

### realtime

```bash
ga-cli realtime 257850630                                          # Active users now
ga-cli realtime 257850630 --dimensions country --metrics activeUsers --limit 10
```

| Flag | Default | Description |
|------|---------|-------------|
| `--dimensions` | (none) | Comma-separated |
| `--metrics` | `activeUsers` | Comma-separated |
| `--limit` | 100 | Max rows |

### metadata

```bash
ga-cli metadata 257850630    # List all available dimensions and metrics for a property
```

Returns JSON with `dimensions` and `metrics` arrays, each containing `apiName`, `uiName`, `description`, `category`.

## Output Format

All commands output indented JSON to stdout. Reports and realtime return:

```json
{
  "row_count": 7,
  "rows": [
    {"date": "20260222", "sessions": "1234", "activeUsers": "987"},
    {"date": "20260221", "sessions": "1100", "activeUsers": "900"}
  ]
}
```

Rows are flattened objects — dimension and metric names become keys.

Accounts list returns an array of account objects with nested `properties` arrays.

Metadata returns `{"dimensions": [...], "metrics": [...]}`.

## Build

```bash
make install                    # Install to $GOPATH/bin/ga-cli
make build                      # Build to ./bin/ga-cli
go install ./cmd/ga-cli         # Direct Go install
make test                       # Run tests
```

Requires Go 1.25+.

## Project Structure

```
cmd/ga-cli/main.go                    # Entry point
internal/auth/credentials.go          # Google Cloud credential handling
internal/commands/root.go             # Root cobra command
internal/commands/accounts.go         # accounts list
internal/commands/properties.go       # properties get
internal/commands/report.go           # report run
internal/commands/realtime.go         # realtime
internal/commands/metadata.go         # metadata
internal/output/json.go               # JSON output formatter
```
