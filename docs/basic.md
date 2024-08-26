## Basic Usage

```
A simple tool to update DNS dynamically if you want to create DNS records
for hosts that IP address might changes. Support multiple DNS providers.

Usage:
  micro-ddns [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  run         Start micro-ddns server.
  version     Print version information about micro-ddns.

Flags:
  -h, --help          help for micro-ddns
  -v, --verbose int   log level, available options are -4 (DEBUG), 0 (INFO), 4 (WARN) and 8 (ERROR)

Use "micro-ddns [command] --help" for more information about a command.
```

Basically when you want to run this program, simply type

```bash
micro-ddns run -c /path/to/config.yaml
```

If you want to use Docker or other containerized solutions to run this program,
 you can use the Docker image of this program by using this command (assumes you are using Docker)

```bash
docker run --name ddns -v /host/path/to/config.yaml:/etc/micro-ddns/config.yaml masteryyh/micro-ddns:latest
```
