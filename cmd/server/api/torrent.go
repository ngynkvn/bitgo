package api

import (
	"bitgo/cmd/server/messages"

	"github.com/sanity-io/litter"
)

func (api *API) AddTorrent(p messages.ParamsAddTorrent) error {
	litter.Dump(p)
	_, err := api.db.DB.Exec(`INSERT INTO torrents(file, path, progress) VALUES (?, ?, ?)`, p.File, p.OutputPath, 0.0)
	return err
}
