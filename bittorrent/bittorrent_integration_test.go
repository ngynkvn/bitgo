package bittorrent

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/jackpal/bencode-go"
	"github.com/stretchr/testify/assert"
)

func TestDebianTorrent(t *testing.T) {
	f, err := os.Open("testdata/debian.torrent")
	assert.NoError(t, err)
	torrent, err := Open(f)
	assert.NoError(t, err)

	assert.Equal(t, torrent.InfoHash, [20]uint8{0x2b, 0x66, 0x98, 0x0, 0x93, 0xbc, 0x11, 0x80, 0x6f, 0xab, 0x50, 0xcb, 0x3c, 0xb4, 0x18, 0x35, 0xb9, 0x5a, 0x3, 0x62})
	assert.Equal(t, torrent.Name, "debian-12.5.0-amd64-netinst.iso")
	assert.Equal(t, torrent.PieceLength, 262144)
	assert.Equal(t, torrent.Length, 659554304)
	assert.Equal(t, torrent.Announce, "http://bttracker.debian.org:6969/announce")
}

func TestDebianTorrentRequest(t *testing.T) {
	f, err := os.Open("testdata/debian.torrent")
	assert.NoError(t, err)
	torrent, err := Open(f)
	assert.NoError(t, err)

	base, err := url.Parse(torrent.Announce)
	assert.NoError(t, err)

	params := url.Values{}
	params.Add("info_hash", string(torrent.InfoHash[:]))
	params.Add("peer_id", "kevinkevinkevinkevin")
	params.Add("port", "6881")
	params.Add("uploaded", "0")
	params.Add("downloaded", "0")
	params.Add("compact", "1")
	params.Add("left", strconv.Itoa(torrent.Length))

	base.RawQuery = params.Encode()

	resp, err := http.Get(base.String())
	assert.NoError(t, err)
	out, err := bencode.Decode(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, out, "interval")
	assert.Contains(t, out, "peers")
}
