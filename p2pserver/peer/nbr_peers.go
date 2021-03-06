package peer

import (
	"fmt"
	"sync"

	"dan-road-vbft/p2pserver/common"
	"dan-road-vbft/p2pserver/message/types"
)

//NbrPeers: The neigbor list
type NbrPeers struct {
	sync.RWMutex
	List map[uint64]*Peer
}

//Broadcast tranfer msg buffer to all establish peer
func (this *NbrPeers) Broadcast(msg types.Message, isConsensus bool) {
	this.RLock()
	defer this.RUnlock()
	for _, node := range this.List {
		if node.syncState == common.ESTABLISH && node.GetRelay() == true {
			node.Send(msg, isConsensus)
		}
	}
}

//NodeExisted return when peer in nbr list
func (this *NbrPeers) NodeExisted(uid uint64) bool {
	_, ok := this.List[uid]
	return ok
}

//GetPeer return peer according to id
func (this *NbrPeers) GetPeer(id uint64) *Peer {
	this.Lock()
	defer this.Unlock()
	n, ok := this.List[id]
	if ok == false {
		return nil
	}
	return n
}

//AddNbrNode add peer to nbr list
func (this *NbrPeers) AddNbrNode(p *Peer) {
	this.Lock()
	defer this.Unlock()

	if this.NodeExisted(p.GetID()) {
		fmt.Printf("[p2p]insert an existed node\n")
	} else {
		this.List[p.GetID()] = p
	}
}

//DelNbrNode delete peer from nbr list
func (this *NbrPeers) DelNbrNode(id uint64) (*Peer, bool) {
	this.Lock()
	defer this.Unlock()

	n, ok := this.List[id]
	if ok == false {
		return nil, false
	}
	delete(this.List, id)
	return n, true
}

//initialize nbr list
func (this *NbrPeers) Init() {
	this.List = make(map[uint64]*Peer)
}

//NodeEstablished whether peer established according to id
func (this *NbrPeers) NodeEstablished(id uint64) bool {
	this.RLock()
	defer this.RUnlock()

	n, ok := this.List[id]
	if ok == false {
		return false
	}

	if n.syncState != common.ESTABLISH {
		return false
	}

	return true
}

//GetNeighborAddrs return all establish peer address
func (this *NbrPeers) GetNeighborAddrs() []common.PeerAddr {
	this.RLock()
	defer this.RUnlock()

	var addrs []common.PeerAddr
	for _, p := range this.List {
		if p.GetSyncState() != common.ESTABLISH {
			continue
		}
		var addr common.PeerAddr
		addr.IpAddr, _ = p.GetAddr16()
		addr.Time = p.GetTimeStamp()
		addr.Services = p.GetServices()
		addr.Port = p.GetSyncPort()
		addr.ID = p.GetID()
		addrs = append(addrs, addr)
	}

	return addrs
}

//GetNeighborHeights return the id-height map of nbr peers
func (this *NbrPeers) GetNeighborHeights() map[uint64]uint64 {
	this.RLock()
	defer this.RUnlock()

	hm := make(map[uint64]uint64)
	for _, n := range this.List {
		if n.GetSyncState() == common.ESTABLISH {
			hm[n.GetID()] = n.GetHeight()
		}
	}
	return hm
}

//GetNeighbors return all establish peers in nbr list
func (this *NbrPeers) GetNeighbors() []*Peer {
	this.RLock()
	defer this.RUnlock()
	peers := []*Peer{}
	for _, n := range this.List {
		if n.GetSyncState() == common.ESTABLISH {
			node := n
			peers = append(peers, node)
		}
	}
	return peers
}

//GetNbrNodeCnt return count of establish peers in nbrlist
func (this *NbrPeers) GetNbrNodeCnt() uint32 {
	this.RLock()
	defer this.RUnlock()
	var count uint32
	for _, n := range this.List {
		if n.GetSyncState() == common.ESTABLISH {
			count++
		}
	}
	return count
}
