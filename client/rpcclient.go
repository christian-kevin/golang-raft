package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gookit/color"
	"golang-raft/server"

	"net/rpc"
	"os"
	"strconv"
	"strings"
)

func getNodePort(nodeId int) int {
	portStr := "1300" + strconv.Itoa(nodeId)
	port, _ := strconv.Atoi(portStr)
	return port
}

func main() {
	nodeId := flag.Int("id", 2, "node id")
	flag.Parse()

	var (
		addr    = "127.0.0.1:" + strconv.Itoa(getNodePort(*nodeId))
		request = &server.CommandArgs{
			CommandName: "get",
			Params:      []string{"foo"},
		}
	)

	// Establish the connection to the address of the
	// RPC server
	srv, err := rpc.Dial("tcp", addr)

	if err != nil {
		fmt.Printf("Server is down : %+v", err)
		return
	}

	defer srv.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		response := new(server.Reply)

		fmt.Println("Enter command: ")
		args, _ := reader.ReadString('\n')
		args = args[:len(args)-1]
		commands := strings.Split(args, " ")

		if len(commands) > 0 && commands[0] == "exit" {
			break
		}

		if len(commands) > 0 {
			params := make([]string, 0, 0)

			request.CommandName = commands[0]
			request.Params = append(params, commands[1:]...)

			err := srv.Call("CommandHandler.Handle", request, response)
			if err != nil {
				fmt.Printf("Server is down : %+v", err)
				break
			}

			if response.Value == nil {
				fmt.Println(response.Value)
				continue
			}

			color.Red.Println(response.Value.(string))
		}

	}
}
