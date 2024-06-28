package bittorrent

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
)

func bitTorrentHandshake(infoHash string, peerID string) []byte {
	return []byte("\x13BitTorrent protocol\x00\x00\x00\x00\x00\x10\x00\x00\x00\x00\x00\x00\x00\x00" + infoHash + peerID)
}

func readHandshake(r io.Reader) (string, string, error) {
	buf := make([]byte, 68)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return "", "", err
	}
	if !bytes.Equal(buf[:20], []byte("\x13BitTorrent protocol")) {
		return "", "", errors.New("invalid handshake")
	}
	infoHash := hex.EncodeToString(buf[20:40])
	peerID := hex.EncodeToString(buf[40:60])
	return infoHash, peerID, nil
}
