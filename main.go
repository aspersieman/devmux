package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

var version = "dev"

func parseLogLevel(s string) (Level, error) {
	switch s {
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "warn":
		return LevelWarn, nil
	case "error":
		return LevelError, nil
	default:
		return LevelInfo, fmt.Errorf("invalid log level: %s", s)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "devmux %s\n\n", version)
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  devmux [options] <config.yml>\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func main() {
	logLevelFlag := flag.String(
		"log-level",
		"info",
		"log level: debug, info, warn, error",
	)

	forceFlag := flag.Bool(
		"force",
		false,
		"kill and recreate tmux session if it exists",
	)

	versionFlag := flag.Bool(
		"version",
		false,
		"print version and exit",
	)

	flag.Usage = usage
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		return
	}

	level, err := parseLogLevel(*logLevelFlag)
	if err != nil {
		Err(err.Error())
		os.Exit(1)
	}
	LogLevel = level

	if flag.NArg() < 1 {
		usage()
		os.Exit(1)
	}

	cfgPath := flag.Arg(0)
	cfg := loadConfig(cfgPath)

	session := cfg.Session
	if session == "" {
		Err("session name is required")
		usage()
		os.Exit(1)
	}

	if tmuxHasSession(session) {
		if *forceFlag {
			logger("killing existing tmux session", LevelWarn)
			_ = tmux("kill-session", "-t", session)
		} else {
			_ = tmux("attach", "-t", session)
			return
		}
	}

	// Session does not exist, create it
	if err := tmux("new-session", "-d", "-s", session); err != nil {
		Err(fmt.Sprintf("failed to create session: %v\r", err))
		os.Exit(1)
	}

	// Repos, editor windows
	for i, repo := range cfg.Repos {
		target := fmt.Sprintf("%s:%d", session, i+1)

		tmux("new-window",
			"-t", session,
			"-n", repo.Name,
		)

		path := expand(repo.Path)
		cmd := fmt.Sprintf("cd %s && %s .", path, repo.Editor)
		tmuxSend(target, cmd)
	}

	// Command windows
	for _, c := range cfg.Commands {
		tmux(
			"new-window",
			"-t", session,
			"-n", c.Name,
		)

		target := fmt.Sprintf("%s:%s", session, c.Name)
		baseCmd := fmt.Sprintf("cd %s", expand(c.Path))

		// Pane-based layout
		if len(c.Panes) > 0 {
			// First pane
			tmuxSend(target, fmt.Sprintf("%s && %s", baseCmd, c.Panes[0].Cmd))

			// Additional panes
			for i := 1; i < len(c.Panes); i++ {
				tmuxSplit(target, true) // horizontal split
				paneTarget := fmt.Sprintf("%s.%d", target, i)
				tmuxSend(paneTarget, fmt.Sprintf("%s && %s", baseCmd, c.Panes[i].Cmd))
			}

			// Nice layout
			tmux("select-layout", "-t", target, "even-vertical")
			continue
		}

		// Single command fallback
		if c.Cmd != "" {
			tmuxSend(target, fmt.Sprintf("%s && %s", baseCmd, c.Cmd))
		}
	}

	// Focus first window and attach
	tmux("select-window", "-t", session+":1")
	tmux("attach", "-t", session)
}

func loadConfig(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		Err(err.Error())
		os.Exit(1)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}

func tmuxHasSession(session string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", session)
	return cmd.Run() == nil
}

func expand(path string) string {
	if path[:1] == "~" {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[1:])
	}
	return path
}
