package api

import (
	"bitgo/bittorrent"
	"bitgo/cmd/server/messages"
	"os"

	"github.com/jackpal/bencode-go"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (api *API) AddTorrent(p messages.ParamsAddTorrent) error {
	file, err := os.Open(p.File)
	if err != nil {
		return err
	}

	torrentInfo := bittorrent.BencodeTorrent{}
	err = bencode.Unmarshal(file, &torrentInfo)
	if err != nil {
		return err
	}

	_, err = api.db.Exec(`INSERT INTO torrents(file, path, torrentinfo, progress) VALUES (?, ?, ?, ?)`, p.File, p.OutputPath, torrentInfo, 0.0)
	return err
}

type Torrent struct {
	File     string
	Path     string
	Progress float32
}

func (api *API) GetTorrents() ([]Torrent, error) {
	result := []Torrent{}
	err := api.db.Select(&result, `SELECT file, path, progress from torrents`)
	return result, err
}

// TODO(refactor)
func StartAppDB() (*sqlx.DB, func()) {
	db := sqlx.MustOpen("sqlite3", DBFilePath)
	log.Info().Msg("DB connected")
	dbInitialization(db)
	return db, func() { os.Remove(DBFilePath); db.Close() }
}

func dbInitialization(db *sqlx.DB) {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS torrents(
        file text UNIQUE,
        path text,
        progress float
    );`)
	if err != nil {
		panic(err.Error())
	}
}

// TODO(config)
const SocketFilePath = "/tmp/bitgo.sock"

// TODO(config): probably not the best place to persist application state
const DBFilePath = "/tmp/bitgo.db"
