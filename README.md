# IP Range Extractor by Country (IPv4 / IPv6)

This Go program extracts IP address ranges for a specific country from the **IP2LOCATION-LITE-DB1** database and writes individual IPs to a text file.

Supports both **IPv4** and **IPv6** via command-line flags.

---

## ✅ Features

- Parses CSV database from IP2Location (Lite version).
- Supports IPv4 and IPv6 formats.
- Filters IP ranges by 2-letter country code.
- Skips overly large ranges (>10 million IPs).
- Logs all skipped/invalid ranges.
- High performance via Go routines and buffering.

---

## 📦 Requirements

- Go 1.2+
- IP2Location CSV files:
  - [`IP2LOCATION-LITE-DB1.CSV`](https://lite.ip2location.com/database/db1/ip-country)
  - [`IP2LOCATION-LITE-DB1.IPV6.CSV`](https://lite.ip2location.com/database/db1/ipv6)

Place the `.CSV` files in a `./db/` folder.

---

## 🚀 Usage

```bash
go run main.go <COUNTRY_CODE> [--v4 | --v6]
```

### Arguments:

- `<COUNTRY_CODE>`: Two-letter ISO 3166-1 alpha-2 country code (e.g. `US`, `UA`, `DE`).
- `--v4`: Process IPv4 data from `IP2LOCATION-LITE-DB1.CSV`.
- `--v6`: Process IPv6 data from `IP2LOCATION-LITE-DB1.IPV6.CSV`.

> ❗ You must specify either `--v4` or `--v6`, not both.

---

## 📁 Output

- Output file:
  - `UA.txt` for IPv4
  - `UA_ipv6.txt` for IPv6
- Log file:
  - `skipped_UA.log` for skipped or invalid ranges

---

## 🧪 Examples

### Extract IPv4 for Ukraine:
```bash
go run main.go UA --v4
```

### Extract IPv6 for Germany:
```bash
go run main.go DE --v6
```

---

## ⚠ Notes

- Large IP ranges (>10 million addresses) are automatically skipped and logged.
- The program uses goroutines and buffered channels to speed up processing.

---

## 📜 License

This project is provided as-is for educational and utility purposes.

IP data © [IP2Location LITE](https://lite.ip2location.com/) under Creative Commons Attribution-ShareAlike 4.0.
