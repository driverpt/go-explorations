package balancer

import (
	"net"
	"net/http"
)

type LoadBalancer interface {
    Init() error
    ServeHTTP(response http.ResponseWriter, request *http.Request)
    NextEndpoint() (int, Endpoint, error)
}

type LoadBalancerState struct {
    currentEndpoints []Endpoint
    currentIndex int
}

type Endpoint struct {
    host net.Addr
    port string
}

type LoadBalancerConfig struct {
    hosts []Endpoint
}

func CreateEndpointFromIP(ip net.IP, port string) Endpoint {
    return Endpoint{ host: &net.IPAddr{ IP: ip }, port: port } 
}