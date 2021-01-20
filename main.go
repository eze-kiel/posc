package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/eze-kiel/posc/scan"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

type options struct {
	limit    int64
	protocol string
	noping   bool
	stfu     bool
	help     bool
}

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		log.Fatalf("%s\n", err)
	}
}

func run(args []string, stdout io.Writer) error {
	var opts options

	flag.BoolVar(&opts.help, "h", false, "Display this help")
	flag.Int64Var(&opts.limit, "limit", 2048, "Number of files that can be opened")
	flag.BoolVar(&opts.noping, "np", false, "Disable ping")
	flag.StringVar(&opts.protocol, "p", "tcp", "Protocol to use")
	flag.BoolVar(&opts.stfu, "q", false, "Enable quiet mode (no logs)")
	flag.Parse()

	// Display help is asked and exit.
	if opts.help {
		usage(os.Args[0])
		os.Exit(0)
	}

	// Check if a target has been provided.
	if len(flag.Args()) != 1 {
		log.Error("no target provided")
		usage(os.Args[0])
		log.Errorf("exiting with status 1")
		os.Exit(1)
	}

	ips, err := net.LookupHost(flag.Arg(0))
	if err != nil {
		log.Fatalf("can not resolve %s", flag.Arg(0))
	}

	// Some recap about parameters.
	if !opts.stfu {
		log.Infof("max open files: %d", opts.limit)
		log.Infof("target: %s (%s)", flag.Arg(0), ips[0])
		log.Infof("protocol: %s", opts.protocol)
	}

	// If the program is not launched as root, disable ping requests and log.
	if os.Getenv("SUDO_USER") == "" && !opts.stfu {
		log.Warn("not running as root, ping has been disabled.")
		opts.noping = true
	}

	ps := &scan.Scanner{
		IP:     flag.Arg(0),
		Prot:   opts.protocol,
		NoPing: opts.noping,
		NoLogs: opts.stfu,
		Lock:   semaphore.NewWeighted(opts.limit),
	}

	start := time.Now()
	ps.Run(1, 65535, 500*time.Millisecond)
	elapsed := time.Since(start)

	if !opts.stfu {
		log.Info("scan ended")
		log.Infof("time elapsed: %s", elapsed)
	}
	return nil
}

func usage(name string) {
	fmt.Printf("Usage of %s: [OPTIONS] target", name)
	fmt.Print(`

OPTIONS

-h		Display this help
-limit int	Number of files that can be opened (default 1024)
-np		Disable ping
-p string	Protocol to use (default "tcp")
-q    		Enable quiet mode (no logs)
`)
}
