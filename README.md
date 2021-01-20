# posc

High performance port scanner written in Golang.

## Usage

Most simple usage:

```
$ posc <ip or address>
```

In order to enable ping requests, you must launch it as root:

```
$ sudo posc 127.0.0.1
```

If `posc` can't reach the target with ICMP, it will warn you and stop the scan. You can ask it to scan even if the target doesn't responds to ICMP ping requests with the flag `-np`:

```
$ sudo posc -np 127.0.0.1
```

The complete list of the options is available with the flag `-h`:

```
$ posc -h
Usage of ./posc: [OPTIONS] target

OPTIONS

-h              Display this help
-limit int      Number of files that can be opened (default 1024)
-np             Disable ping
-p string       Protocol to use (default "tcp")
-q              Enable quiet mode (no logs)
```

## Demo

[![Demo](https://asciinema.org/a/pXWO6QoLBlqufMwoIhcILyvF7.svg)](https://asciinema.org/a/pXWO6QoLBlqufMwoIhcILyvF7)

## License

[MIT](https://choosealicense.com/licenses/mit/)