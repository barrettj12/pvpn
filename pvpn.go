package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	// For the "(re)connect" commands, block the terminal and disconnect when
	// an interrupt is received.
	// Other commands should be passed verbatim to protonvpn-cli.
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "c", "r":
			runBlocking()
			return
		}
	}

	runNonBlocking()
}

func runNonBlocking() {
	pvpnCmd(os.Args[1:]...).Run()
}

func runBlocking() {
	cmd := pvpnCmd(os.Args[1:]...)
	// TODO: cancel command context when signal is received
	cmd.Start()

	signals := make(chan os.Signal)
	signal.Notify(signals,
		syscall.SIGHUP, // doesn't seem to work
		os.Interrupt,
	)

	<-signals
	pvpnCmd("d").Run()
}

func pvpnCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("protonvpn-cli", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
