package scan

import (
	"context"
	"fmt"
	"net"
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
	NoPing bool
	NoLogs bool
	Lock   *semaphore.Weighted
}

// Run runs the portScanner.
func (ps *Scanner) Run(min, max int, timeout time.Duration) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	if !ps.NoPing && !ps.NoLogs && !ps.pingIcmpEchoRequest(timeout) {
		log.Warnf("%s does not respond to ICMP requets. If you think it drops ICMP packets, retry with flag -np.", ps.IP)
		return
	}

	for port := min; port <= max; port++ {
		wg.Add(1)
		ps.Lock.Acquire(context.TODO(), 1)
		go func(port int) {
			defer ps.Lock.Release(1)
			defer wg.Done()
			scanPort(ps.IP, ps.Prot, port, timeout)
		}(port)
	}
}

func scanPort(ip, prot string, port int, timeout time.Duration) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout(prot, target, timeout)
	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			scanPort(ip, prot, port, timeout)
		}
		return
	}

	conn.Close()
	fmt.Printf("%d/%s  \topen\n", port, prot)
}
