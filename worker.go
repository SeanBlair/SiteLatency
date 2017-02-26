package main 

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"strconv"
	"net/rpc"
)
var (
	outboundIp string
	portForMServerRPC string
	portForGetSite string
	portForPingServer string
	serverIpPort string
)

type MServer int



func main() {

	err := ParseArguments()
	if err != nil {
		panic(err)
	}
	fmt.Println("serverIpPort:", serverIpPort)

	outboundIp = GetOutboundIP()

	fmt.Println("this workers outboundIp is:", outboundIp)

	join()

	fmt.Println("Successfully joined. Ports: Server:", portForMServerRPC, "PingServer:", portForPingServer, "GetSite:", portForGetSite)

	// TODO listen rpc outboundIp:portForMServerRPC
}

func join() {
	var joinResp int
	// joinReq := JoinRequest{myIpPort}
	client, err := rpc.Dial("tcp", serverIpPort)
	checkError("rpc.Dial in join()", err, false)
	err = client.Call("MServer.Join", outboundIp, &joinResp)
	checkError("client.Call(MServer.Join: ", err, false)
	err = client.Close()
	checkError("client.Close() in join call: ", err, false)

	portForMServerRPC = strconv.Itoa(joinResp)
	portForPingServer = strconv.Itoa(joinResp + 1)
	portForGetSite = strconv.Itoa(joinResp + 2)
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
func GetOutboundIP() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().String()
    idx := strings.LastIndex(localAddr, ":")

    return localAddr[0:idx]
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