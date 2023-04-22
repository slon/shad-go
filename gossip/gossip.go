//go:build !solution

package gossip

import (
	"time"

	"gitlab.com/slon/shad-go/gossip/meshpb"
	"google.golang.org/grpc"
)

type PeerConfig struct {
	SelfEndpoint string
	PingPeriod   time.Duration
}

type Peer struct {
	config PeerConfig
}

func (p *Peer) UpdateMeta(meta *meshpb.PeerMeta) {
	panic("implement me")
}

func (p *Peer) AddSeed(seed string) {
	panic("implement me")
}

func (p *Peer) Addr() string {
	return p.config.SelfEndpoint
}

func (p *Peer) GetMembers() map[string]*meshpb.PeerMeta {
	panic("implement me")
}

func (p *Peer) Run() {
	panic("implement me")
}

func (p *Peer) Stop() {
	panic("implement me")
}

func NewPeer(config PeerConfig) *Peer {
	panic("implement me")
}
