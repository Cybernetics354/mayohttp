# Mayohttp

A TUI HTTP client written in Go with Pipelines on mind.

![preview](./readme/preview.png "TUI Preview")

![piping](./readme/pipe_preview.png "Piping Preview")

![command pallete](./readme/command_pallete.png "Command Pallete")

![select method](./readme/select_method.png "Select Method")

## Main features

- Terminal pipelines
- Filtering response before pipe-ing (eg. req header, req body, res header, res body)
- Sessions
- Environment file
- You can use env variables literally on anything (url, pipe, response, header) with $NAME syntax
- Open / Edit anything with your favorite editor (set the $EDITOR on your environment)

## 🚀 Automatic Install (Recommended)
> [!IMPORTANT]  
> ✅ Supported OS: Linux and MacOS

```bash
$ /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Cybernetics354/mayohttp/main/.scripts/install-linux.sh)"
```

### ⚙️ Installation parameters
| Variable      | Parameter       | Example                                                     | Remarks                                               |
|---------------|-----------------|-------------------------------------------------------------|-------------------------------------------------------|
| `VERSION`     | `--version`     | `VERSION=v0.0.2` or<br>`--version v0.0.2`                   | Install specific version (default: latest release)    |
| `DESTINATION` | `--destination` | `DESTINATION=/usr/bin` or<br>`--destination /usr/local/bin` | Install to custom directory (default: `~/.local/bin`) |

Example:
```bash
$ VERSION=v0.0.2 DESTINATION=/usr/bin /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Cybernetics354/mayohttp/main/.scripts/install-linux.sh)"
```

## ✋ Manual Install (Alternative)
If you prefer, you can install manually by hand:

1. Go to the [GitHub releases page](https://github.com/Cybernetics354/mayohttp/releases)

2. Download the asset matching your OS and architecture, e.g.:
```
mayohttp_0.0.2_linux_amd64.tar.gz
```

3. Download the checksum file:
```
mayohttp_0.0.2_checksums.txt
```

4. Verify the checksum:
```bash
$ grep " mayohttp_0.0.2_linux_amd64.tar.gz$" mayohttp_0.0.2_checksums.txt | sha256sum -c -
```

5. Extract the archive:
```bash
$ tar -xzf mayohttp_0.0.2_linux_amd64.tar.gz
```

6. Move the executable into your preferred directory:
```bash
$ chmod +x mayohttp
$ mv mayohttp ~/.local/bin/
```

7. If `~/.local/bin` is not in your `PATH`, add it in your `~/.bashrc`, `~/.zshrc` or `~/.profile`:
```bash
$ export PATH="$HOME/.local/bin:$PATH"
```

8. Reload the shell profile:
```bash
$ source ~/.profile
```
