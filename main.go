package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mdp/qrterminal/v3"
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
	secondsFloat, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return "unknown"
	}
	totalSeconds := int64(secondsFloat)
	days := totalSeconds / 86400
	hours := (totalSeconds % 86400) / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
}

func readIPBrief() string {
	output, err := exec.Command("ip", "-br", "a").Output()
	if err != nil {
		return "unavailable"
	}
	return strings.TrimSpace(string(output))
}

func readMachineID() string {
	data, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(data))
}

func renderMachineIDQR(writer *bufio.Writer, machineID string) {
	if machineID == "" || machineID == "unknown" {
		fmt.Fprintln(writer, "unavailable")
		return
	}
	config := qrterminal.Config{
		Level:      qrterminal.L,
		Writer:     writer,
		HalfBlocks: true,
		QuietZone:  1,
	}
	qrterminal.GenerateWithConfig(machineID, config)
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
		beijing, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			beijing = time.FixedZone("UTC+8", 8*60*60)
		}
		ipBrief := readIPBrief()
		machineID := readMachineID()

		fmt.Fprint(writer, "\x1b[2J\x1b[H")
		fmt.Fprintln(writer, "TTY1 dashboard")
		fmt.Fprintln(writer, "---------------------")
		fmt.Fprintf(writer, "Time   : %s (UTC+8 BeiJing)\n", time.Now().In(beijing).Format("2006-01-02 15:04:05"))
		fmt.Fprintf(writer, "Uptime : %s seconds\n", readUptime())
		fmt.Fprintf(writer, "IP     : %s\n", ipBrief)
		renderMachineIDQR(writer, machineID)
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
