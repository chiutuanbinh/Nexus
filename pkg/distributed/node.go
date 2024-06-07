package distributed

import "nexus/pkg/common"

type NodeConfig struct {
	LeaderAddr string
}

type NodeIface interface {
	GetID() string
	Put(entry common.Tuple) error
	Get(key string) (common.Tuple, error)
	GetLeader() NodeIface
	GetFollowers() []NodeIface
}
