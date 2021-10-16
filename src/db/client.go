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

type databaseClient struct {
	db     *sql.DB
	Artist dao.ArtistDao
	Album  dao.AlbumDao
	Track  dao.TrackDao
}

type DatabaseClient interface {
	/*
	   Shut down database connection and all associated resources.
	*/
	Close()
}

func NewDatabaseClient(cfg config.Config) databaseClient {
	db, err := sql.Open("mysql", getConnectionString(cfg))
	if err != nil {
		logrus.Panic("Unable to connect to database: ", err.Error())
	}

	client := databaseClient{
		db:     db,
		Artist: mysql.NewArtistDao(db),
		Album:  nil,
		Track:  mysql.NewTrackDao(db),
	}
	client.Album = mysql.NewAlbumDao(client.db, client.Artist, client.Track)

	migrate(db, cfg.MigrationsPath)

	return client
}

func (this databaseClient) Close() {
	this.Artist.Close()
	this.Album.Close()
	this.Track.Close()
	this.db.Close()
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
