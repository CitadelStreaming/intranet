package db

import (
	"database/sql"
	"fmt"

	"citadel_intranet/src/config"
	"citadel_intranet/src/db/dao"
	"citadel_intranet/src/db/dao/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type DatabaseClient struct {
	db     *sql.DB
	Artist dao.ArtistDao
	Album  dao.AlbumDao
	Track  dao.TrackDao
}

func NewDatabaseClientFromConnection(db *sql.DB) DatabaseClient {
	client := DatabaseClient{
		db:     db,
		Artist: mysql.NewArtistDao(db),
		Album:  nil,
		Track:  mysql.NewTrackDao(db),
	}
	client.Album = mysql.NewAlbumDao(client.db, client.Artist, client.Track)

	return client
}

func NewDatabaseClient(cfg config.Config) DatabaseClient {
	db, err := sql.Open("mysql", getConnectionString(cfg))
	if err != nil {
		logrus.Panic("Unable to connect to database: ", err.Error())
	}

	return NewDatabaseClientFromConnection(db)
}

func (this DatabaseClient) Migrate(migrationsPath string) {
	migrate(this.db, migrationsPath)
}

func (this DatabaseClient) Close() {
	if this.Artist != nil {
		this.Artist.Close()
	}

	if this.Album != nil {
		this.Album.Close()
	}

	if this.Track != nil {
		this.Track.Close()
	}

	if this.db != nil {
		this.db.Close()
	}
}

func getConnectionString(cfg config.Config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		cfg.DbUser,
		cfg.DbPass,
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbName,
	)
}
