package bittorrent

import (
	"crypto/sha1"
	"io"

	"github.com/jackpal/bencode-go"
)

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece_length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type BencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	Pieces      [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func (info bencodeInfo) GetInfoHash() ([20]byte, error) {
	sha := sha1.New()
	err := bencode.Marshal(sha, info)
	if err != nil {
		return [20]byte{}, err
	}
	return sha1.Sum(nil), nil
}

func (info bencodeInfo) GetPieces() ([][20]byte, error) {
	pieces := make([][20]byte, info.PieceLength)
	for i := 0; i < len(pieces); i++ {
		copy(pieces[i][:], info.Pieces[i*20:(i+1)*20])
	}
	return pieces, nil
}

func (bto *BencodeTorrent) GetTorrentFile() (TorrentFile, error) {
	info, err := bto.Info.GetInfoHash()
	if err != nil {
		return TorrentFile{}, err
	}

	pieces, err := bto.Info.GetPieces()
	if err != nil {
		return TorrentFile{}, err
	}

	return TorrentFile{
		Announce:    bto.Announce,
		InfoHash:    info,
		Pieces:      pieces,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}, nil
}

// Open parses a torrent file
func Open(r io.Reader) (TorrentFile, error) {
	bto := BencodeTorrent{}
	err := bencode.Unmarshal(r, &bto)
	if err != nil {
		return TorrentFile{}, err
	}
	return bto.GetTorrentFile()
}
