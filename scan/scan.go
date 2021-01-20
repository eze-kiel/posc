package scan

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

// Scanner is a scanner instance. It scans a given ip with a given protocol.
type Scanner struct {
	IP     string
	Prot   string
	Range  string
	NoPing bool
	NoLogs bool
	Lock   *semaphore.Weighted
}

// Run runs the portScanner.
func (ps *Scanner) Run(min, max int, timeout time.Duration) error {
	wg := sync.WaitGroup{}

	if !ps.NoPing && !ps.NoLogs && !ps.pingIcmpEchoRequest(timeout) {
		log.Warnf("%s does not respond to ICMP requets. If you think it drops ICMP packets, retry with flag -np.", ps.IP)
		return nil
	}
	ports, err := readPortsRange(ps.Range)
	if err != nil {
		return err
	}

	for _, p := range ports {
		wg.Add(1)
		ps.Lock.Acquire(context.TODO(), 1)
		go func(port int) {
			defer ps.Lock.Release(1)
			defer wg.Done()
			ps.scanPort(port, timeout)
		}(p)
	}
	wg.Wait()
	return nil
}

func (ps *Scanner) scanPort(port int, timeout time.Duration) {
	target := fmt.Sprintf("%s:%d", ps.IP, port)
	conn, err := net.DialTimeout(ps.Prot, target, timeout)
	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			ps.scanPort(port, timeout)
		}
		return
	}

	conn.Close()
	fmt.Printf("%d/%s  \topen\n", port, ps.Prot)
}

// readPortsRange transforms a range of ports given in conf to an array of
// effective ports
func readPortsRange(ranges string) ([]int, error) {
	ports := []int{}

	parts := strings.Split(ranges, ",")

	for _, spec := range parts {
		if spec == "" {
			continue
		}
		switch spec {
		case "all":
			for port := 1; port <= 65535; port++ {
				ports = append(ports, port)
			}
		case "reserved":
			for port := 1; port < 1024; port++ {
				ports = append(ports, port)
			}
		default:
			var decomposedRange []string

			if !strings.Contains(spec, "-") {
				decomposedRange = []string{spec, spec}
			} else {
				decomposedRange = strings.Split(spec, "-")
			}

			min, err := strconv.Atoi(decomposedRange[0])
			if err != nil {
				return nil, err
			}
			max, err := strconv.Atoi(decomposedRange[len(decomposedRange)-1])
			if err != nil {
				return nil, err
			}

			if min > max {
				return nil, fmt.Errorf("lower port %d is higher than high port %d", min, max)
			}
			if max > 65535 {
				return nil, fmt.Errorf("port %d is higher than max port", max)
			}
			for i := min; i <= max; i++ {
				ports = append(ports, i)
			}
		}
	}

	return ports, nil
}
