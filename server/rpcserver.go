package server

import (
	"errors"
	"fmt"
	"go.etcd.io/etcd/raft/raftpb"
	"log"
	"net"
	"net/rpc"
	"strconv"
	"strings"
)

type CommandArgs struct {
	CommandName string
	Params      []string
}

type Reply struct {
	Value interface{}
}

const ReplyStatusOk = "ok"

type CommandHandler struct {
	store       *Kvstore
	confChangeC chan<- raftpb.ConfChange
}

func errorParams(message string) error {
	log.Printf(message)
	return errors.New(message)
}

func (c *CommandHandler) Handle(args *CommandArgs, reply *Reply) error {
	switch {
	case strings.ToLower(args.CommandName) == "set":
		errMessage := "No sufficient args to put in kv store\n"

		if len(args.Params) != 2 {
			return errorParams(errMessage)
		}

		c.store.Propose(args.Params[0], args.Params[1])
		reply.Value = ReplyStatusOk
	case strings.ToLower(args.CommandName) == "get":
		errMessage := "No sufficient args to set in kv store\n"

		if len(args.Params) != 1 {
			return errorParams(errMessage)
		}
		if v, ok := c.store.Lookup(args.Params[0]); ok {
			reply.Value = v
		}
	case strings.ToLower(args.CommandName) == "register":
		errMessage := "No sufficient args to register new node\n"

		if len(args.Params) != 2 {
			return errorParams(errMessage)
		}

		nodeId, err := strconv.ParseUint(args.Params[0], 0, 64)
		if err != nil {
			errMessage = fmt.Sprintf("Failed to convert ID for conf change (%v)\n", err)
			return errorParams(errMessage)
		}

		cc := raftpb.ConfChange{
			Type:    raftpb.ConfChangeAddNode,
			NodeID:  nodeId,
			Context: []byte("http://127.0.0.1:" + args.Params[1]),
		}
		c.confChangeC <- cc
		reply.Value = ReplyStatusOk
	case strings.ToLower(args.CommandName) == "kill":
		errMessage := "No sufficient args to kill a node\n"

		if len(args.Params) != 1 {
			return errorParams(errMessage)
		}

		nodeId, err := strconv.ParseUint(args.Params[0], 0, 64)
		if err != nil {
			errMessage = fmt.Sprintf("Failed to convert ID for conf change (%v)\n", err)
			return errorParams(errMessage)
		}

		cc := raftpb.ConfChange{
			Type:   raftpb.ConfChangeRemoveNode,
			NodeID: nodeId,
		}
		c.confChangeC <- cc
		reply.Value = ReplyStatusOk
	}
	return nil
}

func StartRPCServer(kv *Kvstore, port int, confChangeC chan<- raftpb.ConfChange, errorC <-chan error) {
	err := rpc.Register(&CommandHandler{
		store:       kv,
		confChangeC: confChangeC,
	})

	if err != nil {
		log.Printf("Failed to register rpc server: %+v", err)
	}

	go func() {
		listener, _ := net.Listen("tcp", ":"+strconv.Itoa(port))
		defer listener.Close()

		// Wait for incoming connections
		rpc.Accept(listener)
	}()

	// exit when raft goes down
	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}
