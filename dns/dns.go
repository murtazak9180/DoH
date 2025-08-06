package dns

import (
	"net"

	"github.com/miekg/dns"
)

func UpstreamDNS(msg *dns.Msg, upstreamAddr string) (*dns.Msg, error) {
	wire, err := msg.Pack()
	if err != nil {
		return nil, err
	}

	//wire to be sent on the upstream is now set.
	serverAddr, err := net.ResolveUDPAddr("udp", upstreamAddr)
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(wire))
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 512) //max response length is 512
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, err
	}

	resp := new(dns.Msg)
	//now fill resp and pack it and return it. fill till the nth byte.

}
