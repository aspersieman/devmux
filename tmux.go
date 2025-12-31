package main

import (
	"fmt"
	"os"
	"os/exec"
)

func tmux(args ...string) error {
	Debug(fmt.Sprintf("tmux %s", args))
	cmd := exec.Command("tmux", args...)
	return cmd.Run()
}

func tmuxSend(target, command string) error {
	return tmux(
		"send-keys",
		"-t", target,
		command,
		"C-m",
	)
}

func tmuxSplit(target string, vertical bool) error {

	flag := "-h"
	if vertical {
		flag = "-v"
	}
	return tmux("split-window", flag, "-t", target)
}

func tmuxAttach(session string) {
	if os.Getenv("TMUX") != "" {
		// Already inside tmux
		cmd := exec.Command("tmux", "switch-client", "-t", session)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
		return
	}

	// Not inside tmux → attach normally
	cmd := exec.Command("tmux", "attach", "-t", session)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		Err(fmt.Sprintf("failed to attach to tmux: %v", err))
		os.Exit(1)
	}
}

func tmuxHasSession(session string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", session)
	return cmd.Run() == nil
}
