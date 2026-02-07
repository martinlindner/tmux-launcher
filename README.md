# tmux-launcher

A smart tmux wrapper that automatically manages session attachment. Instead of manually deciding whether to run `tmux new` or `tmux attach`, tmux-launcher inspects existing sessions and does the right thing:

- **Single detached session** — auto-attaches to it
- **No sessions, or multiple detached sessions** — shows a TUI picker to select a session or create a new one
- **Already inside tmux** — warns and exits (configurable)

> [!IMPORTANT]  
> Disclaimer: This project is essentially a quick test drive of Claude Opus 4.6 via Claude Code. The codebase is 90% vibe-coded.

## Installation

Requires Go 1.21+ and tmux.

```bash
git clone https://github.com/martinlindner/tmux-launcher.git
cd tmux-launcher
make install  # builds and copies to ~/bin/
```

## Usage

```bash
tmux-launcher
```

### Flags

```
--allow-nested     Allow running inside an existing tmux session
--no-auto-attach   Always show the TUI picker instead of auto-attaching
```

### Config file

Optional config at `~/.config/tmux-launcher/config.yaml`:

```yaml
allow_nested: false
auto_attach: true
```

CLI flags override config file values.

## WSL2 setup

**As Windows Terminal default command** (in your WSL profile settings):

```json
{ "commandline": "wsl.exe ~ -e /home/<user>/bin/tmux-launcher" }
```

**From shell config** (`~/.bashrc` or `~/.zshrc`):

```bash
if [ -z "$TMUX" ]; then
    exec ~/bin/tmux-launcher
fi
```
