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

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		log.Fatalf("%s\n", err)
	}
}

func run(args []string, stdout io.Writer) error {
	var ulimit int64
	var prot string
	var noping, stfu, h bool

	flag.BoolVar(&h, "h", false, "Display this help")
	flag.Int64Var(&ulimit, "limit", 2048, "Number of files that can be opened")
	flag.BoolVar(&noping, "np", false, "Disable ping")
	flag.StringVar(&prot, "p", "tcp", "Protocol to use")
	flag.BoolVar(&stfu, "q", false, "Enable quiet mode (no logs)")
	flag.Parse()

	// Display help is asked and exit.
	if h {
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
	if !stfu {
		log.Infof("max open files: %d", ulimit)
		log.Infof("target: %s (%s)", flag.Arg(0), ips[0])
		log.Infof("protocol: %s", prot)
	}

	// If the program is not launched as root, disable ping requests and log.
	if os.Getenv("SUDO_USER") == "" && !stfu {
		log.Warn("not running as root, ping has been disabled.")
		noping = true
	}

	ps := &scan.Scanner{
		IP:     flag.Arg(0),
		Prot:   prot,
		NoPing: noping,
		NoLogs: stfu,
		Lock:   semaphore.NewWeighted(ulimit),
	}

	start := time.Now()
	ps.Run(1, 65535, 500*time.Millisecond)
	elapsed := time.Since(start)

	if !stfu {
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
