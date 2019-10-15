# systemd-auto-config
## Program to configure systemd services interactively 

### Presets available:

1: Simple - A single command is executed and is considered to be the main process/service

2: Forking - A Service which expects the main process/service to fork and launch a child process

3: Oneshot - A Service which is expected to start and finish its job before other units are launched

### Compiled binaries:

> [Linux amd64 / x86_64](https://cdn.sidsun.com/systemd-auto-config/systemd-auto-config_linux-amd64)

### Debian Packages:

> [amd64](https://cdn.sidsun.com/systemd-auto-config/systemd-auto-config.deb)

### Use YAPPA ( Yet Another Personal Package Archive ):

```bash
curl -s --compressed "https://sid-sun.github.io/yappa/KEY.gpg" | sudo apt-key add -
curl -s --compressed "https://sid-sun.github.io/yappa/yappa.list" | sudo tee /etc/apt/sources.list.d/yappa.list
sudo apt update
sudo apt install systemd-auto-config
```

:)