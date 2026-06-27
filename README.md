# itunesScript

A shell project to control **iTunes** from the console on macOS.

It wraps AppleScript behind a small command dispatcher so playback can be driven
from the terminal — handy for hotkeys, remote shells, or scripting your music
without leaving the keyboard.

## Commands

Built-in commands live under `builtin/`:

| Command    | Action                          |
| ---------- | ------------------------------- |
| `play` / `pause` / `stop` | Transport controls.  |
| `next` / `previous`       | Skip between tracks. |
| `mute` / `unmute` / `vol` | Volume controls.     |
| `shuffle`                 | Toggle shuffle.      |
| `status`                  | Show what's playing. |
| `quit`                    | Quit iTunes.         |

## Install

```sh
git clone https://github.com/helmedeiros/itunesScript.git
chmod +x itunesScript.sh
```

Add the clone directory to your `PATH` so the command is available everywhere.

## License

[MIT](LICENSE)
