package main

import (
	"fmt"
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
