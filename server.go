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
	"net"
	"os"
	"net/rpc"
)
var (
	workerIncomingIpPort string
	clientIncomingIpPort string
)


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

	listenClient()
}

func (p *MServer) MeasureWebsite(mSiteReq MWebsiteReq, mRes *MRes) error {
	*mRes = measureWebsite(mSiteReq)
	return nil
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

func measureWebsite(mSite MWebsiteReq) (res MRes) {
	fmt.Println("Website to measure:", mSite.URI, "SamplesPerWorker:", mSite.SamplesPerWorker)

	res = MRes{map[string]LatencyStats{
		"hardcodedWorkerIp" : LatencyStats{1,2,3},
		},
		nil}
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