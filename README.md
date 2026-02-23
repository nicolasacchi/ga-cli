# ga-cli

Command-line tool for Google Analytics 4. Single binary, JSON output, service account auth.

Uses the [GA4 Data API](https://developers.google.com/analytics/devguides/reporting/data/v1) and [Admin API](https://developers.google.com/analytics/devguides/config/admin/v1).

## Install

### From source

```bash
go install github.com/nicolasacchi/ga-cli/cmd/ga-cli@latest
```

### From release

Download the binary for your platform from [Releases](https://github.com/nicolasacchi/ga-cli/releases).

```bash
curl -L https://github.com/nicolasacchi/ga-cli/releases/latest/download/ga-cli_Linux_x86_64.tar.gz | tar xz
mv ga-cli ~/.local/bin/
```

### From source (local)

```bash
git clone https://github.com/nicolasacchi/ga-cli.git
cd ga-cli
make install
```

## Authentication

Uses the standard Google Cloud credential chain:

1. `GOOGLE_APPLICATION_CREDENTIALS` environment variable (service account JSON key)
2. Application Default Credentials (`gcloud auth application-default login`)

The service account needs **Viewer** role on the GA4 property.

## Usage

### List accounts and properties

```bash
ga-cli accounts list
```

### Get property details

```bash
ga-cli properties get 257850630
```

### Run a report

```bash
# Sessions and active users by date (last 7 days)
ga-cli report run 257850630 --dimensions date --metrics sessions,activeUsers

# Page views by country for a specific date range
ga-cli report run 257850630 \
  --dimensions country \
  --metrics screenPageViews \
  --start-date 2026-02-16 \
  --end-date 2026-02-22

# Top pages with pagination
ga-cli report run 257850630 \
  --dimensions pagePath \
  --metrics sessions,activeUsers \
  --limit 50 --offset 0
```

### Realtime report

```bash
# Current active users
ga-cli realtime 257850630

# Active users by country
ga-cli realtime 257850630 --dimensions country --metrics activeUsers
```

### List available dimensions and metrics

```bash
ga-cli metadata 257850630
```

## Output

All commands output JSON to stdout. Example report output:

```json
{
  "row_count": 7,
  "rows": [
    {"date": "20260222", "sessions": "1234", "activeUsers": "987"},
    {"date": "20260221", "sessions": "1100", "activeUsers": "900"}
  ]
}
```

## License

MIT
