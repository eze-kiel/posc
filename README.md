# posc

(Another) Small and fast port scanner written in Golang.

This is more a PoC I wrote for [scan-exporter](https://github.com/devops-works/scan-exporter) than a real program.

## Installation

You can either:

1. Clone the repo and build it (Golang is required):

```
$ git clone https://github.com/eze-kiel/posc.git
$ cd posc
$ go build .
```

2. Download the latest build [from the releases](https://github.com/eze-kiel/posc/releases)

## Usage

Most simple usage:

```
$ posc target port-range
```

For example: 

```
$ posc 195.66.45.12 reserved
```

In order to enable ping requests, you must launch it as root:

```
$ sudo posc target port-range
```

If `posc` can't reach the target with ICMP, it will warn you and stop the scan. You can ask it to scan even if the target doesn't responds to ICMP ping requests with the flag `-np`:

```
$ sudo posc -np target port-range
```

The complete list of the options is available with the flag `-h`:

```
$ posc -h
Usage of ./posc: [OPTIONS] target

Target can be an IP address or an URL
Ports can be: "all, "reserved", "22", "22-443", "1-1023,1337-4242"...

OPTIONS

-h		Display this help
-limit int	Number of files that can be opened (default 1024)
-np		Disable ping
-p string	Protocol to use (default "tcp")
-q    		Enable quiet mode (no logs)
```

## Demo

[![Demo](https://asciinema.org/a/8f5fPT9Ou3VemY7kLtLwNKTDw.svg)](https://asciinema.org/a/8f5fPT9Ou3VemY7kLtLwNKTDw)

## License

[MIT](https://choosealicense.com/licenses/mit/)
