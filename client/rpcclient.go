package main

import (
	"fmt"
	"golang-raft/server"
	"net/rpc"
	"strconv"
)

func main() {
	var (
		addr    = "127.0.0.1:" + strconv.Itoa(9121)
		request = &server.CommandArgs{
			CommandName: "get",
			Params:      []string{"foo"},
		}
		response = new(server.Reply)
	)

	// Establish the connection to the adddress of the
	// RPC server
	client, _ := rpc.Dial("tcp", addr)
	defer client.Close()

	// Perform a procedure call (core.HandlerName == Handler.Execute)
	// with the Request as specified and a pointer to a response
	// to have our response back.
	_ = client.Call("CommandHandler.Handle", request, response)
	fmt.Println(response.Value)
}
