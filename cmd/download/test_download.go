package main

import (
	"bitgo/bittorrent"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/sanity-io/litter"
)

func main() {
	path := os.Args[1]
	f, err := os.Open(path)
	tfile, err := bittorrent.Open(f)
	if err != nil {
		panic(err)
	}

	// Make announce request

	params := url.Values{}
	params.Add("info_hash", string(tfile.InfoHash[:]))
	params.Add("peer_id", "kevinkevinkevinkevin")
	params.Add("port", "6881")
	params.Add("uploaded", "0")
	params.Add("downloaded", "0")
	params.Add("compact", "1")
	params.Add("left", strconv.Itoa(tfile.Length))

	base, err := url.Parse(tfile.Announce)
	if err != nil {
		panic(err)
	}
	base.RawQuery = params.Encode()

	resp, err := http.Get(base.String())
	if err != nil {
		panic(err)
	}

	ar, err := bittorrent.UnmarshalAnnounceResponse(resp.Body)
	if err != nil {
		panic(err)
	}
	litter.Dump(ar)

	peers, err := ar.GetPeers()
	if err != nil {
		panic(err)
	}
	for _, peer := range peers[:5] {
		litter.Dump(peer)
	}
}
