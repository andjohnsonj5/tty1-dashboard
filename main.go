package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func readUptime() string {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "unknown"
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return "unknown"
	}
	return fields[0]
}

func main() {
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	// Hide cursor while running.
	fmt.Fprint(writer, "\x1b[?25l")
	writer.Flush()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	render := func() {
		fmt.Fprint(writer, "\x1b[2J\x1b[H")
		fmt.Fprintln(writer, "TTY1 dashboard (demo)")
		fmt.Fprintln(writer, "---------------------")
		fmt.Fprintf(writer, "Time   : %s\n", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Fprintf(writer, "Uptime : %s seconds\n", readUptime())
		fmt.Fprintln(writer, "")
		fmt.Fprintln(writer, "Press Ctrl+C to exit.")
		writer.Flush()
	}

	render()

	for {
		select {
		case <-ticker.C:
			render()
		case <-sigCh:
			fmt.Fprint(writer, "\x1b[?25h")
			writer.Flush()
			return
		}
	}
}
