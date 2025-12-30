# devmux

devmux is a small Go-based CLI tool that creates a fully configured tmux development environment from a single YAML file.

It provides a one-command, reproducible development setup: tmux sessions, windows, panes, editors, and long-running commands are all created automatically.
README.md

---

## Features

- Declarative tmux session setup using YAML
- One tmux window per git repository
- Automatically opens editor in each repository
- Optional command windows (Docker, logs, watchers, etc.)
- Multi-pane layouts per window
- Session reuse or forced recreation
- Configurable log levels
- Single static Go binary

---

## Installation

### Build from source

Clone the repository and build the binary:

```bash
git clone https://github.com/yourusername/devmux.git
cd devmux
make
```

This will produce a `devmux` binary in the project directory.

### Install into $GOBIN

```bash
make install
```

---

## Usage

```bash
devmux [options] <config.yml>
```

### Examples

```bash
devmux my-project.yml
devmux --force my-project.yml
devmux --log-level debug my-project.yml
devmux --version
```

---

## Command-line options

- `--log-level`  
  Set log level: debug, info, warn, error (default: info)

- `--force`  
  Kill and recreate the tmux session if it already exists

- `--version`  
  Print version information and exit

- `--help`  
  Show usage information

---

## Configuration file

Example `my-project.yml`:

```yaml
session: my-project

repos:
  - name: project-a
    path: ~/src/repos/project-a
    editor: nvim

  - name: my-project
    path: ~/src/repos/my-project
    editor: nvim

commands:
  - name: docker-project-a
    path: ~/src/repos/project-a
    panes:
      - cmd: docker compose up
      - cmd: docker compose logs -f

  - name: docker-project-b
    path: ~/src/repos/project-b
    cmd: docker compose up
```

---

## How it works

1. A tmux session is created using the configured session name
2. Each repository gets its own tmux window
   - The window title is set to the repo name
   - The working directory is set
   - Editor (e.g neovim) is launched automatically
3. Command windows can:
   - Run a single command
   - Or create multiple panes, each running its own command
4. If the session already exists:
   - It is reused by default
   - Or recreated when the `--force` flag is supplied

---

## Pane layouts

When multiple panes are defined:

```yaml
panes:
  - cmd: docker compose up
  - cmd: docker compose logs -f
```

tmux panes are created automatically and an even horizontal layout is applied.

---

## Development

### Build

```bash
make
```

### Clean

```bash
make clean
```

### Versioning

The version string is injected at build time using ldflags:

```bash
make build
devmux --version
```

---

## Requirements

- tmux
- Go 1.20 or newer
- Neovim (or another editor of your choice)

---

## License

MIT

---

## Contributing

Issues and pull requests are welcome.

Keep it simple, explicit, and predictable — just like tmux.
