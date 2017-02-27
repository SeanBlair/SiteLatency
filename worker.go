package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"sort"
	"strconv"
	"strings"
	// "io/ioutil"
	"net/http"
	"time"
)

var (
	portForWorkerRPC  string
	portForPingServer string
	serverIpPort      string
)

// //Messages
// const (
// 	JOIN = iota
// 	// GETWORKERS
// )

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

	join()

	fmt.Println("Successfully joined. Ports: Server:", portForWorkerRPC, "PingServer:", portForPingServer)

	// listen on own ip, specific port so server knows how to access
	listen(":" + portForWorkerRPC)
}

func (p *WorkerServer) PingSite(req MWebsiteReq, resp *LatencyStats) error {
	fmt.Println("received call to PingSite")
	*resp = pingSite(req)
	return nil
}

func (p *WorkerServer) PingServer(samples int, resp *LatencyStats) error {
	fmt.Println("received call to PingServer")
	// TODO
	// pingServer(samples)
	*resp = LatencyStats{7, 7, 7}
	return nil
}

func pingSite(req MWebsiteReq) (stats LatencyStats) {
	var latencyList []int

	for i := 0; i < req.SamplesPerWorker; i++ {
		latency := pingSiteOnce(req.URI)
		latencyList = append(latencyList, latency)
	}

	fmt.Println("latencyList before sorting:", latencyList)

	sort.Ints(latencyList)

	fmt.Println("latencyList after sorting:", latencyList)
	min := latencyList[0]
	max := latencyList[len(latencyList)-1]
	median := getMedian(latencyList)
	stats = LatencyStats{min, median, max}
	return
}

func pingSiteOnce(uri string) (l int) {
	start := time.Now()
	// res, err := http.Get(uri)
	_, err := http.Get(uri)
	elapsed := time.Since(start)

	checkError("Error in pingSiteOnce(), http.Get():", err, true)
	// html, err := ioutil.ReadAll(res.Body)
	// res.Body.Close()
	// checkError("Error in pingSiteOnce(), ioutil.ReadAll():", err, true)
	// fmt.Printf("%s", html)

	l = int(elapsed / time.Millisecond)

	return l
}

// list is sorted
func getMedian(list []int) (m int) {
	length := len(list)
	var middle int = length / 2
	if (length % 2) == 1 {
		return list[middle]
	} else {
		a := list[middle-1]
		b := list[middle]
		return (a + b) / 2
	}
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

	// TODO make more elegant than space delimiter...
	port, err := bufio.NewReader(conn).ReadString(' ')
	checkError("Error in join(), bufio.NewReader(conn).ReadString()", err, true)
	fmt.Println("Message from server: ", port)

	portForWorkerRPC = strings.Trim(port, " ")
	fmt.Println("My portForWorkerRPC is:", portForWorkerRPC)

	portValue, err := strconv.Atoi(portForWorkerRPC)
	checkError("Error in join(), strconv.Atoi()", err, true)

	portForPingServer = strconv.Itoa(portValue + 1)
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

// Prints msg + err to console and exits program if exit == true
func checkError(msg string, err error, exit bool) {
	if err != nil {
		log.Println(msg, err)
		if exit {
			os.Exit(-1)
		}
	}
}
