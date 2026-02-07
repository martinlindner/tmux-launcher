# tmux-launcher

A smart tmux wrapper that automatically manages session attachment. Instead of manually deciding whether to run `tmux new` or `tmux attach`, tmux-launcher inspects existing sessions and does the right thing:

- **No sessions** — creates a new session automatically
- **Single detached session** — auto-attaches to it
- **Multiple sessions** — shows a TUI picker to select a session or create a new one
- **Already inside tmux** — warns and exits (configurable)

> [!IMPORTANT]  
> Disclaimer: This project is essentially a quick test drive of Claude Opus 4.6 via Claude Code. The codebase is 90% vibe-coded.

## Installation

Requires Go 1.25+ and tmux.

```bash
git clone https://github.com/martinlindner/tmux-launcher.git
cd tmux-launcher
make install  # builds and copies to ~/.local/bin/
```

## Usage

```bash
tmux-launcher
```

### Flags

```
--allow-nested         Allow running inside an existing tmux session
--no-auto-attach       Always show the TUI picker instead of auto-attaching
--no-auto-new-session  Show the TUI picker even when no sessions exist
```

### Config file

Optional config at `~/.config/tmux-launcher/config.yaml`:

```yaml
allow_nested: false
auto_attach: true
auto_new_session: true
```

CLI flags override config file values.

## WSL2 setup

**As Windows Terminal default command** (in your WSL profile settings):

```json
{ "commandline": "wsl.exe ~ -e /home/<user>/.local/bin/tmux-launcher" }
```

**From shell config** (`~/.bashrc` or `~/.zshrc`):

```bash
if [ -z "$TMUX" ]; then
    exec ~/.local/bin/tmux-launcher
fi
```
