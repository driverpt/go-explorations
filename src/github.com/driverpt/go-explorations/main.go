package main 

import (
    "os"
    "net"
    "net/http"
    "fmt"
    "github.com/driverpt/go-explorations/balancer"
)

func main() {
    port := ""
    if len(os.Args) > 1 { 
        port = os.Args[1]
    } else { 
        port = "8000"
    }

    endpoints := []balancer.Endpoint{ balancer.CreateEndpointFromIP(net.ParseIP("127.0.0.1"), "10080" ),
                                      balancer.CreateEndpointFromIP(net.ParseIP("127.0.0.1"), "10081" ),
                                      }
    
    rr := balancer.NewRoundRobinLoadBalancer(endpoints)
    rr.Init()
    
    fmt.Println("Starting Server @ port 8000");
    server := &http.Server {
	    Addr: ":" + port,
	    Handler: http.HandlerFunc(rr.ServeHTTP),
	}
    server.ListenAndServe();
}

