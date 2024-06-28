package bittorrent

import (
	"encoding/binary"
	"errors"
	"net"
)

// Peer is a peer in a peer list
// It is made up of an IP and a port number
type Peer struct {
	IP   net.IP
	Port uint16
}

// IntoPeer converts a byte slice into a peer
func IntoPeer(bytes []byte) (Peer, error) {
	if len(bytes) != 6 {
		return Peer{}, errors.New("invalid peer")
	}
	return Peer{
		IP:   net.IP(bytes[0:4]),
		Port: binary.BigEndian.Uint16(bytes[4:6]),
	}, nil
}

// GetPeers returns the peers in the response
func (ar AnnounceResponse) GetPeers() ([]Peer, error) {
	nPeers := len(ar.Peers) / 6
	peers := make([]Peer, nPeers)
	for i := range peers {
		section := ar.Peers[i*6 : (i+1)*6]
		peer, err := IntoPeer([]byte(section))
		if err != nil {
			return nil, err
		}
		peers[i] = peer
	}
	return peers, nil
}
