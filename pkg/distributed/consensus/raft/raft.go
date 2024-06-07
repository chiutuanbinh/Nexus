package raft

type StateType uint8

const (
	Follower StateType = iota
	Candidate
	Leader
)

type Peer struct {
	State  StateType
	Leader *Peer
	Peers  []*Peer
}

func CreatePeer() *Peer {
	return &Peer{
		State:  Follower,
		Leader: nil,
		Peers:  make([]*Peer, 0),
	}
}

func (p *Peer) GetLeader() *Peer {
	return p.Leader
}

func (p *Peer) Run() {

}
