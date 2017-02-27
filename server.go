/*
Implements the server in assignment 4 for UBC CS 416 2016 W2.

Usage:

go run server.go [worker-incoming ip:port] [client-incoming ip:port]

Example:

go run server.go 127.0.0.1:1111 127.0.0.1:2222
*/

package main 

import (
	"fmt"
	"log"
	"net"
	"os"
	"net/rpc"
	"strings"
	"strconv"
)
var (
	workerIncomingIpPort string
	clientIncomingIpPort string
	nextWorkerRPCPort int = 20000
	Workers []Worker
)

type Worker struct {
	Ip string
	Port int
}


// A stats struct that summarizes a set of latency measurements to an
// internet host.
type LatencyStats struct {
	Min    int // min measured latency in milliseconds to host
	Median int // median measured latency in milliseconds to host
	Max    int // max measured latency in milliseconds to host
}

/////////////// RPC structs

// Resource server type.
type MServer int

// Request that client sends in RPC call to MServer.MeasureWebsite
type MWebsiteReq struct {
	URI              string // URI of the website to measure
	SamplesPerWorker int    // Number of samples, >= 1
}

// Response to:
// MServer.MeasureWebsite:
//   - latency stats per worker to a *URI*
//   - (optional) Diff map
// MServer.GetWorkers
//   - latency stats per worker to the *server*
type MRes struct {
	Stats map[string]LatencyStats    // map: workerIP -> LatencyStats
	Diff  map[string]map[string]bool // map: [workerIP x workerIP] -> True/False
}

// Request that client sends in RPC call to MServer.GetWorkers
type MWorkersReq struct {
	SamplesPerWorker int // Number of samples, >= 1
}

/////////////// /RPC structs



func main() {

	err := ParseArguments()
	if err != nil {
		panic(err)
	}
	fmt.Println("workerIncomingIpPort:", workerIncomingIpPort, "clientIncomingIpPort:", clientIncomingIpPort)

	go listenClient()
	listenWorkers()
}

func (p *MServer) MeasureWebsite(mSiteReq MWebsiteReq, mRes *MRes) error {
	*mRes = measureWebsite(mSiteReq)
	return nil
}

func (p *MServer) GetWorkers(workerReq MWorkersReq, wRes *MRes) error {
	*wRes = getWorkers(workerReq.SamplesPerWorker)
	return nil
}

// func (p *MServer) Join(workerIP string, port *int) error {
// 	*port = nextWorkerPort
// 	Workers = append(Workers, Worker{workerIP, nextWorkerPort})
// 	nextWorkerPort += 10
// 	fmt.Println(Workers)
// 	return nil
// }

func listenWorkers() {
	ln, err := net.Listen("tcp", workerIncomingIpPort)
	checkError("Error in listenWorkers(), net.Listen():", err, true)
	for {
		conn, err := ln.Accept()
		checkError("Error in listenWorkers(), ln.Accept():", err, true)
		// go joinWorker(conn)
		joinWorker(conn)
		fmt.Println("Worker joined. Workers:", Workers)
	}
}

func joinWorker(conn net.Conn) {
	workerIpPort := conn.RemoteAddr().String()
	fmt.Println("joining Workers ip:", workerIpPort)

	workerIp := workerIpPort[:strings.Index(workerIpPort, ":")]

	Workers = append(Workers, Worker{workerIp, nextWorkerRPCPort})
	// send to socket
	// TODO change to not require space delimiter
    fmt.Fprintf(conn, strconv.Itoa(nextWorkerRPCPort) + " ")
    nextWorkerRPCPort += 10
}

func listenClient() {
	mServer := rpc.NewServer()
	m := new(MServer)
	mServer.Register(m)
	l, err := net.Listen("tcp", clientIncomingIpPort)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go mServer.ServeConn(conn)
	}
}


func getWorkers(samples int) (res MRes) {
	fmt.Println("GetWorkers called with samples:", samples)

	res.Stats = make(map[string]LatencyStats)

	for _, worker := range Workers {
		stats := pingServer(worker, samples)
		res.Stats[worker.Ip] = stats
		res.Diff = nil
	}

	// res = MRes{map[string]LatencyStats{
	// 	"hardcodedWorkerIp" : LatencyStats{3,2,1},
	// 	},
	// 	nil}
		return
}

func pingServer(w Worker, samples int) (st LatencyStats) {
	wIpPort := getWorkerIpPort(w)
	client, err := rpc.Dial("tcp", wIpPort)
	checkError("rpc.Dial in pingServer()", err, true)
	err = client.Call("WorkerServer.PingServer", samples, &st)
	checkError("client.Call(WorkerServer.PingServer: ", err, true)
	err = client.Close()
	checkError("client.Close() in pingServer call: ", err, true)
	return
}

func measureWebsite(mSite MWebsiteReq) (res MRes) {
	fmt.Println("Website to measure:", mSite.URI, "SamplesPerWorker:", mSite.SamplesPerWorker)

	res.Stats = make(map[string]LatencyStats)
	

	for _, worker := range Workers {
		stats := pingSite(worker, mSite)
		res.Stats[worker.Ip] = stats
		res.Diff = nil
	}

	// res = MRes{map[string]LatencyStats{
	// 	"hardcodedWorkerIp" : LatencyStats{1,2,3},
	// 	},
	// 	nil}

	return
}

func pingSite(w Worker, req MWebsiteReq) (st LatencyStats) {
	wIpPort := getWorkerIpPort(w)
	client, err := rpc.Dial("tcp", wIpPort)
	checkError("rpc.Dial in pingSite()", err, true)
	err = client.Call("WorkerServer.PingSite", req, &st)
	checkError("client.Call(WorkerServer.PingSite: ", err, true)
	err = client.Close()
	checkError("client.Close() in pingSite call: ", err, true)
	return
}

func getWorkerIpPort(w Worker) (s string) {
	s = w.Ip + ":" + strconv.Itoa(w.Port)
	return
}

func ParseArguments() (err error) {
	arguments := os.Args[1:]
	if len(arguments) == 2 {
		workerIncomingIpPort = arguments[0]
		clientIncomingIpPort = arguments[1]
		} else {
			err = fmt.Errorf("Usage: {go run server.go [worker-incoming ip:port] [client-incoming ip:port]}")
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