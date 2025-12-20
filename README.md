# TTY1 Dashboard Demo

This is a minimal Go program that runs on tty1 and redraws the screen every
two seconds. It is intended to be launched by systemd and replaces getty on
tty1.

## What it does

- Clears the screen and renders a simple status view.
- Refreshes every 2 seconds.
- Handles SIGINT/SIGTERM and restores the cursor.

## Files

- `main.go` - dashboard program
- `go.mod` - module metadata
- `/usr/local/bin/tty1-dashboard` - compiled binary
- `/etc/systemd/system/tty1-dashboard.service` - systemd unit

## Build

```bash
go build -o /usr/local/bin/tty1-dashboard
```

## Install (systemd)

1) Disable getty on tty1:

```bash
systemctl disable --now getty@tty1.service
```

2) Install the unit:

```ini
[Unit]
Description=TTY1 Dashboard Demo
After=systemd-user-sessions.service
Conflicts=getty@tty1.service

[Service]
ExecStart=/usr/local/bin/tty1-dashboard
Restart=always
RestartSec=1
StandardInput=tty
StandardOutput=tty
TTYPath=/dev/tty1
TTYReset=yes
TTYVHangup=yes
TTYVTDisallocate=yes

[Install]
WantedBy=multi-user.target
```

3) Enable and start:

```bash
systemctl daemon-reload
systemctl enable --now tty1-dashboard.service
```

## SELinux note

If SELinux blocks execution (status shows 203/EXEC), relabel the binary:

```bash
restorecon -v /usr/local/bin/tty1-dashboard
```

## Restore getty

```bash
systemctl disable --now tty1-dashboard.service
systemctl enable --now getty@tty1.service
```
