package main

type Config struct {
	Session  string        `yaml:"session"`
	Repos    []RepoConfig  `yaml:"repos"`
	Commands []CmdConfig   `yaml:"commands"`
}

type RepoConfig struct {
	Name   string `yaml:"name"`
	Path   string `yaml:"path"`
	Editor string `yaml:"editor"`
}

type CmdConfig struct {
	Name  string       `yaml:"name"`
	Path  string       `yaml:"path"`
	Cmd   string       `yaml:"cmd"`
	Panes []PaneConfig `yaml:"panes"`
}

type PaneConfig struct {
	Cmd string `yaml:"cmd"`
}
