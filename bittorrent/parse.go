package bittorrent

import (
	"bytes"
	"crypto/sha1"
	"io"

	"github.com/jackpal/bencode-go"
)

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

// GetInfoHash returns the info hash of the info
func (info bencodeInfo) GetInfoHash() ([20]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := bencode.Marshal(buffer, info)
	if err != nil {
		return [20]byte{}, err
	}
	return sha1.Sum(buffer.Bytes()), nil
}

// GetPieces returns the pieces in the info
func (info bencodeInfo) GetPieces() ([][20]byte, error) {
	pieces := make([][20]byte, info.Length/info.PieceLength)
	for i := range pieces {
		copy(pieces[i][:], info.Pieces[i*20:(i+1)*20])
	}
	return pieces, nil
}

// BencodeTorrent is a torrent file that has been unmarshalled from bencode
type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

// TorrentFile is the flattened version of BencodeTorrent
type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	Pieces      [][20]byte
	PieceLength int
	Length      int
	Name        string
}

// GetTorrentFile returns a TorrentFile from a BencodeTorrent
func (bto *bencodeTorrent) GetTorrentFile() (TorrentFile, error) {
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
	bto := bencodeTorrent{}
	err := bencode.Unmarshal(r, &bto)
	if err != nil {
		return TorrentFile{}, err
	}
	return bto.GetTorrentFile()
}

// AnnounceResponse is the response to an announce request
type AnnounceResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

// UnmarshalAnnounceResponse unmarshals an announce response
func UnmarshalAnnounceResponse(r io.Reader) (AnnounceResponse, error) {
	ar := AnnounceResponse{}
	err := bencode.Unmarshal(r, &ar)
	if err != nil {
		return AnnounceResponse{}, err
	}
	return ar, nil
}
