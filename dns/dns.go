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
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(wire))
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 4096) //max response length is 4096
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, err
	}

	resp := new(dns.Msg)
	if err := resp.Unpack(buffer[:n]); err != nil {
		return nil, err
	}

	return resp, nil
}
