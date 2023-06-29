package dns

import (
	"errors"
	"fmt"
	"log"
	"net"
	"nsproxy/config"

	"github.com/redis/go-redis/v9"
)

type Server struct {
	repo DNSRepository
}

var ErrUnspportedType = errors.New("Unsupported Query Type")
var ErrBadRequest = errors.New("Bad Request")
var ErrServerError = errors.New("there is problem with dns proxy or dns server")

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
	var err error
	pm, err := parseRequest(data)

	print(pm.NsName)
	print(pm.NsType)

	if err != nil {
		s.WriteErr(address, proxySocket, ErrBadRequest)
	}

	if pm.NsType != AAAA && pm.NsType != A {
		s.WriteErr(address, proxySocket, ErrUnspportedType)
		return
	}

	nsResponse, count, err := s.repo.GetDNSFromCache(pm.NsName)

	if err != nil && err != redis.Nil {
		s.WriteErr(address, proxySocket, ErrServerError)
		return
	}

	if nsResponse != nil {
		resp := s.generateCachedResponse(data, count, nsResponse)
		s.SendResponse(address, proxySocket, resp)

	}

	response, err := s.GetFromDNSServer(data)
	if err != nil {
		s.WriteErr(address, proxySocket, ErrServerError)
		return
	}

	cpResponse := make([]byte, len(response))
	copy(cpResponse, response)

	parsedResponse, err := parseResponse(response)
	if err != nil {
		s.WriteErr(address, proxySocket, ErrServerError)
		return
	}

	fmt.Printf("%v", parsedResponse.Answers)

	err = s.repo.SetDNSInCache(pm.NsName, parsedResponse.Answers, parsedResponse.AnswerCount)
	if err != nil {
		s.WriteErr(address, proxySocket, ErrServerError)
		return
	}
	s.SendResponse(address, proxySocket, cpResponse)
}

func (s *Server) SendResponse(address *net.UDPAddr, proxySocket *net.UDPConn, response []byte) {
	// Send the DNS response back to the client
	fmt.Printf("%v", response)
	var err error
	_, err = proxySocket.WriteToUDP(response, address)
	if err != nil {
		log.Println("Failed to send DNS response to client:", err)
		return
	}
}

func (s *Server) GetFromDNSServer(data []byte) ([]byte, error) {
	var err error
	dnsServers := config.GetConfig().ExternalDnsServers
	var dnsServerAddr *net.UDPAddr
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
		return nil, err
	}
	defer dnsConn.Close()

	_, err = dnsConn.Write(data)
	if err != nil {
		log.Println("Failed to send DNS request:", err)
		return nil, err
	}

	// Receive the DNS response from the DNS server
	response := make([]byte, 4096)
	_, _, err = dnsConn.ReadFromUDP(response)
	if err != nil {
		log.Println("Failed to receive DNS response:", err)
		return nil, err
	}

	return response, nil
}

func (s *Server) WriteErr(address *net.UDPAddr, proxySocket *net.UDPConn, err error) {
	s.SendResponse(address, proxySocket, []byte(err.Error()))
}

func (s *Server) generateCachedResponse(req []byte, count []byte, resp []byte) []byte {
	fresp := append(req[0:32], resp...)
	// update count
	fresp[6] = count[0]
	fresp[7] = count[1]
	// update headers
	fresp[2] = 129
	fresp[3] = 128
	return fresp
}
func NewServer(repo DNSRepository) *Server {
	return &Server{repo: repo}
}
