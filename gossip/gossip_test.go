package gossip_test

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/slon/shad-go/gossip"
	"gitlab.com/slon/shad-go/gossip/meshpb"
	"go.uber.org/goleak"
	"google.golang.org/grpc"
)

const (
	pingPeriod = time.Millisecond * 50
	waitPeriod = pingPeriod * 10
)

type env struct {
	newPeer func() (*gossip.Peer, func())
}

func newEnv(t *testing.T) *env {
	t.Cleanup(func() {
		goleak.VerifyNone(t)
	})

	return &env{
		newPeer: func() (*gossip.Peer, func()) {
			lsn, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)

			peer := gossip.NewPeer(gossip.PeerConfig{
				SelfEndpoint: lsn.Addr().String(),
				PingPeriod:   pingPeriod,
			})

			server := grpc.NewServer()
			meshpb.RegisterGossipServiceServer(server, peer)

			go func() { _ = server.Serve(lsn) }()
			go peer.Run()

			stop := func() {
				server.Stop()
				peer.Stop()
			}
			t.Cleanup(stop)

			return peer, stop
		},
	}
}

func TestGossip_SinglePeer(t *testing.T) {
	env := newEnv(t)

	peer0, _ := env.newPeer()

	members := peer0.GetMembers()
	require.Contains(t, members, peer0.Addr())

	peer0.UpdateMeta(&meshpb.PeerMeta{Name: "peer0"})

	members = peer0.GetMembers()
	require.Contains(t, members, peer0.Addr())
	require.Equal(t, "peer0", members[peer0.Addr()].Name)

	deadPeer, stopDeadPeer := env.newPeer()
	stopDeadPeer()

	peer0.AddSeed(deadPeer.Addr())
	time.Sleep(waitPeriod)

	require.Len(t, peer0.GetMembers(), 1)
}

func TestGossip_TwoPeers(t *testing.T) {
	env := newEnv(t)

	peer0, _ := env.newPeer()
	peer1, stop1 := env.newPeer()

	peer0.AddSeed(peer1.Addr())
	time.Sleep(waitPeriod)

	members0 := peer0.GetMembers()
	require.Contains(t, members0, peer0.Addr())
	require.Contains(t, members0, peer1.Addr())

	members1 := peer1.GetMembers()
	require.Contains(t, members1, peer0.Addr())
	require.Contains(t, members1, peer1.Addr())

	peer1.UpdateMeta(&meshpb.PeerMeta{Name: "bob"})
	time.Sleep(waitPeriod)
	require.Equal(t, "bob", peer0.GetMembers()[peer1.Addr()].Name)

	peer1.UpdateMeta(&meshpb.PeerMeta{Name: "sam"})
	time.Sleep(waitPeriod)
	require.Equal(t, "sam", peer0.GetMembers()[peer1.Addr()].Name)

	stop1()
	time.Sleep(waitPeriod)

	members0 = peer0.GetMembers()
	require.NotContains(t, members0, peer1.Addr())
}

func TestGossip_ManyPeers(t *testing.T) {
	env := newEnv(t)

	seed, stopSeed := env.newPeer()

	var peers []*gossip.Peer
	names := map[string]string{}

	for i := 0; i < 10; i++ {
		peer, _ := env.newPeer()
		peer.AddSeed(seed.Addr())
		peer.UpdateMeta(&meshpb.PeerMeta{Name: fmt.Sprint(i)})
		names[peer.Addr()] = fmt.Sprint(i)
		peers = append(peers, peer)
	}

	time.Sleep(waitPeriod)
	stopSeed()
	time.Sleep(waitPeriod)

	for _, peer := range peers {
		members := peer.GetMembers()
		require.NotContains(t, members, seed.Addr())
		for addr, name := range names {
			require.Contains(t, members, addr)
			require.Equal(t, members[addr].Name, name)
		}
	}

	peers[0].UpdateMeta(&meshpb.PeerMeta{Name: "leader"})
	time.Sleep(waitPeriod)

	for _, peer := range peers {
		members := peer.GetMembers()

		require.Contains(t, members, peers[0].Addr())
		require.Equal(t, members[peers[0].Addr()].Name, "leader")
	}
}

func TestGossip_Groups(t *testing.T) {
	env := newEnv(t)

	aSize, bSize := 1, 1
	seedA, _ := env.newPeer()
	seedB, _ := env.newPeer()

	for i := 0; i < 10; i++ {
		peer, _ := env.newPeer()

		if rand.Int()%2 == 0 {
			peer.AddSeed(seedA.Addr())
			aSize++
		} else {
			peer.AddSeed(seedB.Addr())
			bSize++
		}
	}

	time.Sleep(waitPeriod)

	require.Len(t, seedA.GetMembers(), aSize)
	require.Len(t, seedB.GetMembers(), bSize)

	seedA.AddSeed(seedB.Addr())
	time.Sleep(waitPeriod)

	require.Len(t, seedA.GetMembers(), aSize+bSize)
	require.Len(t, seedB.GetMembers(), aSize+bSize)
}
