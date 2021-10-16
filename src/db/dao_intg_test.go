// +build integration

package db_test

import (
	"os"
	"testing"

	"citadel_intranet/src/config"
	"citadel_intranet/src/db"
	"citadel_intranet/src/db/model"

	"github.com/stretchr/testify/assert"
)

func TestDaoCalls(t *testing.T) {
	assert := assert.New(t)

    wd, err := os.Getwd()
    assert.Nil(err)

    cfg := config.Config{
        DbHost: "localhost",
        DbPort: 3306,
        DbUser: "root",
        DbPass: "pass",
        DbName: "testbed",

        MigrationsPath: wd + "/../../migrations/",
    }

    db := db.NewDatabaseClient(cfg)
    assert.NotNil(db)
    db.Migrate(cfg.MigrationsPath)
    defer db.Close()

    artist := model.Artist{
        Name: "James",
    }

    artistId, err := db.Artist.Save(artist)
    assert.Nil(err)
    artist.Id = artistId

    artists := db.Artist.LoadAll()
    assert.Len(artists, 1)
    assert.Equal(artist, artists[0])

    album := model.Album{
        Title: "Something Awesome",
        Artist: artist,
        Tracks: []model.Track{},
        Published: false,
        Rating: 0,
    }

    album.Id, err = db.Album.Save(album)
    assert.Nil(err)

    albums := db.Album.LoadAll()
    assert.Len(albums, 1)
    assert.Equal(album, albums[0])

    track := model.Track{
        Title: "Track 1",
        AlbumId: album.Id,
        Rating: 0,
    }

    track.Id, err = db.Track.Save(track)
    assert.Nil(err)

    tracks := db.Track.LoadAll()
    assert.Len(tracks, 1)
    assert.Equal(track, tracks[0])

    tracksForAlbum := db.Track.LoadForAlbum(album.Id)
    assert.Len(tracksForAlbum, 1)
    assert.Equal(track, tracksForAlbum[0])

    rows, err := db.Album.Delete(album)
    assert.Nil(err)
    assert.Equal(int64(1), rows)

    retrievedArtist := db.Artist.Load(artist.Id)
    assert.Equal(artist, *retrievedArtist)

    tracksForAlbum = db.Track.LoadForAlbum(album.Id)
    assert.Len(tracksForAlbum, 0)
}
