package grpc_pool

import (
	"errors"
	"file-server-gateway/model"
	"fmt"
	"sort"

	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/config"
	fs "smart.gitlab.biomind.com.cn/intelligent-system/biogo/file_server"
)

var (
	ConnMap       map[string]Pool
	NodeMap       map[string]*fs.ServerNode
	leastNodePool Pool
	current       int
)

func init() {
	ConnMap = make(map[string]Pool)
	NodeMap = make(map[string]*fs.ServerNode)
}

func LoadLeastNode(m map[string]*fs.ServerNode) (err error) {
	nodes := getNodes(m)
	if len(nodes) == 0 {
		return errors.New("server node null")
	}

	pool, err := getSpeifiedConn(nodes[0].NodeName)
	if err != nil {
		return
	}
	leastNodePool = pool
	return
}

func getNodes(m map[string]*fs.ServerNode) []*fs.ServerNode {
	nodes := make([]*fs.ServerNode, len(m))
	index := 0
	for _, v := range m {
		nodes[index] = v
		index++
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].DirSize < nodes[j].DirSize
	})
	return nodes
}

func getSpeifiedConn(nodeName string) (Pool, error) {
	key := fmt.Sprintf("/%s/%s/%s",
		config.GlobalConfig.Namespace,
		model.FileServerNodePrefix,
		nodeName,
	)
	pool, ok := ConnMap[fmt.Sprintf(key)]
	if !ok {
		return nil, fmt.Errorf("not match key %s for pool", key)
	}
	return pool, nil
}

func GetLeastNodePool() Pool {
	return leastNodePool
}

func GetNodeConn() (Pool, error) {
	nodes := getNodes(NodeMap)
	if len(nodes) == 0 {
		return nil, errors.New("server node null")
	}
	if current >= len(nodes) {
		current = 0
	}
	index := (current + 1) % len(nodes)
	node := nodes[index]
	current++
	pool, err := getSpeifiedConn(node.NodeName)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
