package main 

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	// "strings"
	// "strconv"
	"net/rpc"
)
var (
	outboundIp string
	portForWorkerRPC string
	portForGetSite string
	portForPingServer string
	serverIpPort string
)

//Messages
const (
	JOIN = iota
	// GETWORKERS
)

type MServer int

type WorkerServer int

// A stats struct that summarizes a set of latency measurements to an
// internet host.
type LatencyStats struct {
	Min    int // min measured latency in milliseconds to host
	Median int // median measured latency in milliseconds to host
	Max    int // max measured latency in milliseconds to host
}


// Request that client sends in RPC call to MServer.MeasureWebsite
type MWebsiteReq struct {
	URI              string // URI of the website to measure
	SamplesPerWorker int    // Number of samples, >= 1
}



func main() {

	err := ParseArguments()
	if err != nil {
		panic(err)
	}
	fmt.Println("serverIpPort:", serverIpPort)

	// outboundIp = GetOutboundIP()

	fmt.Println("this workers outboundIp is:", outboundIp)

	join()

	fmt.Println("Successfully joined. Ports: Server:", portForWorkerRPC, "PingServer:", portForPingServer, "GetSite:", portForGetSite)

	// TODO listen rpc outboundIp:portForMServerRPC
	// listen(outboundIp + ":" + portForWorkerRPC)
}

func (p *WorkerServer) PingSite(req MWebsiteReq, resp *LatencyStats) error {
	fmt.Println("received call to PingSite")
	// TODO
	// pingSite(req)
	*resp = LatencyStats{5,5,5}
	return nil
}

func (p *WorkerServer) PingServer(samples int, resp *LatencyStats) error {
	fmt.Println("received call to PingServer")
	// TODO
	// pingServer(samples)
	*resp = LatencyStats{7,7,7}
	return nil
}

func listen(ipPort string) {
	wServer := rpc.NewServer()
	w := new(WorkerServer)
	wServer.Register(w)
	l, err := net.Listen("tcp", ipPort)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go wServer.ServeConn(conn)
	}
}

func join() {

	conn, err := net.Dial("tcp", serverIpPort)
	checkError("Error in join(), net.Dial()", err, true)

	fmt.Println("dialed server")

	port, err := bufio.NewReader(conn).ReadString(' ')
	checkError("Error in join(), bufio.NewReader(conn).ReadString()", err, true)
    fmt.Println("Message from server: ", port)

	// var joinResp int
	// client, err := rpc.Dial("tcp", serverIpPort)
	// checkError("rpc.Dial in join()", err, true)
	// err = client.Call("MServer.Join", outboundIp, &joinResp)
	// checkError("client.Call(MServer.Join: ", err, true)
	// err = client.Close()
	// checkError("client.Close() in join call: ", err, true)

	// portForWorkerRPC = strconv.Itoa(joinResp)
	// portForPingServer = strconv.Itoa(joinResp + 1)
	// portForGetSite = strconv.Itoa(joinResp + 2)
}

func ParseArguments() (err error) {
	arguments := os.Args[1:]
	if len(arguments) == 1 {
		serverIpPort = arguments[0]
		} else {
			err = fmt.Errorf("Usage: {go run server.go [server ip:port]}")
			return
		}
	return
}


// From http://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
// Get preferred outbound ip of this machine
// func GetOutboundIP() string {
//     conn, err := net.Dial("udp", "8.8.8.8:80")
//     if err != nil {
//         log.Fatal(err)
//     }
//     defer conn.Close()

//     localAddr := conn.LocalAddr().String()
//     idx := strings.LastIndex(localAddr, ":")

//     return localAddr[0:idx]

//     return "localhost"
// }

// Prints msg + err to console and exits program if exit == true
func checkError(msg string, err error, exit bool) {
	if err != nil {
		log.Println(msg, err)
		if exit {
			os.Exit(-1)
		}
	}
}