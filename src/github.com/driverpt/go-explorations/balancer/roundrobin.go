package balancer

import (
	"net"
	"net/http"
	"net/url"
	"sync"
	"fmt"
)

func NewRoundRobinLoadBalancer(endpoints []Endpoint) *RoundRobinLoadBalancer {
    return &RoundRobinLoadBalancer{
        		config: LoadBalancerConfig{ hosts: endpoints },
        		transport: &http.Transport{DisableKeepAlives: false, DisableCompression: false}, 
        	}
}

type RoundRobinLoadBalancer struct {
    config    LoadBalancerConfig
    state     LoadBalancerState
    lock      sync.RWMutex
    transport *http.Transport
}

func (rrlb *RoundRobinLoadBalancer) Init() error {
    rrlb.state = LoadBalancerState{ currentEndpoints: rrlb.config.hosts, currentIndex: 0 }
    return nil
}

func (rrlb *RoundRobinLoadBalancer) NextEndpoint() (int, Endpoint, error) {
	rrlb.lock.Lock()
	// This will be called when the function returns
	defer rrlb.lock.Unlock()
	
	rrlb.state.currentIndex = (rrlb.state.currentIndex + 1) % len(rrlb.state.currentEndpoints)
	
	return rrlb.state.currentIndex, rrlb.state.currentEndpoints[rrlb.state.currentIndex], nil	
}

func (rrlb *RoundRobinLoadBalancer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
    _, endpoint, _ := rrlb.NextEndpoint()

    fmt.Println("Next Endpoint => " + net.JoinHostPort(endpoint.host.String(), endpoint.port))
    
    request.Header.Add("X-FORWARDED-FOR", request.RemoteAddr)
    request.Header.Add("Server", "Go Load Balancer")
	fmt.Println(request.Proto)
	client := &http.Client{}
	
	newRequest := new(http.Request)
	*newRequest = *request
	
	fmt.Println("Request URI => " + request.RequestURI)
	
	uri, _ := url.ParseRequestURI(request.RequestURI)
	
	if len(uri.Scheme) == 0 {
	    uri.Scheme = "http"
	}
	
	newRequest.URL = uri
	
	newRequest.URL.Host = net.JoinHostPort(endpoint.host.String(), endpoint.port)
	newRequest.URL.User = request.URL.User
	
	newRequest.RequestURI = ""
	
    clientResponse, _ := client.Do(newRequest)
    
    clientResponse.Write(response)
}
