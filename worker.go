package main 

import (
	"fmt"
	"os"
)
var (
	serverIpPort string
)



func main() {

	err := ParseArguments()
	if err != nil {
		panic(err)
	}
	fmt.Println("serverIpPort:", serverIpPort)
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