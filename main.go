// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"go.etcd.io/etcd/raft/raftpb"
	"strings"

	"golang-raft/server"
)

func main() {
	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	id := flag.Int("id", 1, "node ID")
	rpcPort := flag.Int("port", 9121, "rpc server port")
	join := flag.Bool("join", false, "join an existing cluster")
	flag.Parse()

	proposeC := make(chan string)
	defer close(proposeC)
	confChangeC := make(chan raftpb.ConfChange)
	defer close(confChangeC)

	// raft provides a commit stream for the proposals from the http api
	var kvs *server.Kvstore
	getSnapshot := func() ([]byte, error) { return kvs.GetSnapshot() }
	commitC, errorC, snapshotterReady := server.NewRaftNode(*id, strings.Split(*cluster, ","), *join, getSnapshot, proposeC, confChangeC)

	kvs = server.NewKVStore(<-snapshotterReady, proposeC, commitC, errorC)

	server.StartRPCServer(kvs, *rpcPort, confChangeC, errorC)
}
