package dns

import (
	"fmt"
	"log"
	"net"
	"nsproxy/config"
)

type Server struct {
	repo DNSRepository
}

func (s *Server) StartServer(host string, port string) error {
	// Resolve the proxy IP and port
	proxyAddr, err := net.ResolveUDPAddr("udp", host+":"+port)
	if err != nil {
		log.Fatal("Failed to resolve proxy address:", err)
	}

	// Create a UDP socket and bind it to the proxy IP and port
	proxySocket, err := net.ListenUDP("udp", proxyAddr)
	if err != nil {
		log.Fatal("Failed to start DNS proxy:", err)
	}
	defer proxySocket.Close()

	log.Printf("DNS proxy server is listening on %s:%s...", host, port)

	// Buffer for incoming DNS requests
	buffer := make([]byte, 4096)

	for {
		// Receive DNS requests from clients
		n, addr, err := proxySocket.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Failed to read DNS request:", err)
			continue
		}

		log.Printf("Received DNS request from %s:%d", addr.IP, addr.Port)

		// Handle the DNS request
		go s.requestHandler(buffer[:n], addr, proxySocket)
	}
}

func (s *Server) requestHandler(data []byte, address *net.UDPAddr, proxySocket *net.UDPConn) {
	// Create a UDP socket and send the DNS request to the DNS server
	// read from cache
	// read add ip addresses
	// set in cache

	dnsServers := config.GetConfig().ExternalDnsServers

	var dnsServerAddr *net.UDPAddr
	var err error
	for i := 0; i < len(dnsServers); i++ {
		dnsServerAddr, err = net.ResolveUDPAddr("udp", dnsServers[i])
		if err != nil {
			log.Println(fmt.Sprintf("Failed to resolve DNS server address with server %s:", dnsServers[i]), err)
			continue
		}
		break
	}

	if dnsServerAddr == nil {
		log.Println("DNS servers not responding !")
	}

	dnsConn, err := net.DialUDP("udp", nil, dnsServerAddr)
	if err != nil {
		log.Println("Failed to connect to DNS server:", err)
		return
	}
	defer dnsConn.Close()

	_, err = dnsConn.Write(data)
	if err != nil {
		log.Println("Failed to send DNS request:", err)
		return
	}

	// Receive the DNS response from the DNS server
	response := make([]byte, 4096)
	_, _, err = dnsConn.ReadFromUDP(response)
	if err != nil {
		log.Println("Failed to receive DNS response:", err)
		return
	}

	// Send the DNS response back to the client
	_, err = proxySocket.WriteToUDP(response, address)
	if err != nil {
		log.Println("Failed to send DNS response to client:", err)
		return
	}
}

func NewServer(repo DNSRepository) *Server {
	return &Server{repo: repo}
}
