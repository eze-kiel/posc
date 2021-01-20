package scan

import (
	"net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// This function comes from github.com/projectdiscovery/naabu.
// I added RTT log.
func (ps *Scanner) pingIcmpEchoRequest(timeout time.Duration) bool {
	destAddr := &net.IPAddr{IP: net.ParseIP(ps.IP)}
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return false
	}
	defer c.Close()

	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Data: []byte(""),
		},
	}

	data, err := m.Marshal(nil)
	if err != nil {
		return false
	}

	start := time.Now()
	_, err = c.WriteTo(data, destAddr)
	if err != nil {
		return false
	}

	reply := make([]byte, 1500)
	err = c.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return false
	}
	n, SourceIP, err := c.ReadFrom(reply)
	// timeout
	if err != nil {
		return false
	}
	// if anything is read from the connection it means that the host is alive
	if destAddr.String() == SourceIP.String() && n > 0 {
		if !ps.NoLogs {
			rtt := time.Since(start)
			log.Infof("ICMP RTT: %s", rtt)
		}
		return true
	}
	return false
}
