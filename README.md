# go-speedmeter

Internet speed meter terminal buat Linux — lightweight, real-time, zero dependency.

Dibikin buat STB (Set-Top Box) ARM / perangkat low-resource, tapi jalan di mana aja.

```
╔══════════════════════════════════════════╗
║  🌐 SpeedMeter — eth0                   ║
╠══════════════════════════════════════════╣
║                                          ║
║  ⬇ DOWNLOAD:    2.45 Mbps               ║
║    ██████████████░░░░░░░░░░░░░░░░░░░░    ║
║                                          ║
║  ⬆ UPLOAD:     512.00 Kbps              ║
║    ██████░░░░░░░░░░░░░░░░░░░░░░░░░░░░    ║
╠══════════════════════════════════════════╣
║  📊 Total RX: 3.57 GB                   ║
║  📊 Total TX: 2.76 GB                   ║
║  🚀 Peak DL:  8.92 Mbps                 ║
║  🚀 Peak UL:  1.63 Mbps                 ║
║  ⏱  Runtime:  5m30s                     ║
╠══════════════════════════════════════════╣
║  [Ctrl+C] exit                           ║
╚══════════════════════════════════════════╝
```

## Fitur

- **Real-time**: Update tiap 1 detik, cuma angka yang refresh (gak flicker)
- **Zero dependency**: 100% Go standard library, gak pake package eksternal
- **Static binary**: ~1.6MB, tinggal copas ke server mana aja
- **Terminal detection**: Auto detect TUI vs plain-text mode
- **Peak tracking**: Catat peak download & upload selama session
- **Multi interface**: Support milih interface lewat argumen

## Install

### Pre-built binary

Download dari [Releases](https://github.com/chermiel/go-speedmeter/releases).

### Build dari source

```bash
git clone https://github.com/chermiel/go-speedmeter.git
cd go-speedmeter
go build -ldflags="-s -w" -o speedmeter .
```

## Usage

```bash
# Default: eth0
./speedmeter

# Interface lain
./speedmeter wlan0
./speedmeter tailscale0

# Exit
Ctrl+C
```

## Cara kerja

Baca `/proc/net/dev` tiap detik, hitung delta bytes → speed (bps), render ke TUI pake ANSI escape codes. Layout statis digambar sekali, tiap tick cuma update angka-angkanya (no flicker).

## License

MIT
