package application

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"citadel_intranet/src/db"
	"citadel_intranet/src/db/model"
	"citadel_intranet/src/server"

	"github.com/kataras/muxie"
	"github.com/sirupsen/logrus"
)

type App struct {
	server server.Server
	db     db.DatabaseClient
}

type AppErr struct {
	Error string `json:"error"`
}

func NewApp(db db.DatabaseClient, server server.Server) Application {
	return App{
		server: server,
		db:     db,
	}
}

func (this App) Run() {
	this.server.Mux.Handle("/api/v1/album", muxie.Methods().
		HandleFunc(http.MethodGet, this.retrieveAllAlbums).
		HandleFunc(http.MethodPost, this.createAlbum))

	this.server.Mux.Handle("/api/v1/album/:id", muxie.Methods().
		HandleFunc(http.MethodGet, this.retrieveAlbum).
		HandleFunc(http.MethodPut, this.updateAlbum).
		HandleFunc(http.MethodDelete, this.removeAlbum))

	this.server.Mux.Handle("/api/v1/artist", muxie.Methods().
		HandleFunc(http.MethodGet, this.retrieveArtists).
		HandleFunc(http.MethodPost, this.createArtist))
}

func (this App) Close() {
	this.server.Close()
	this.db.Close()
}

func writeBack(out http.ResponseWriter, value interface{}) error {
	if v, ok := value.(error); ok {
		value = AppErr{v.Error()}
	}

	buffer, err := json.Marshal(value)
	if err != nil {
		return err
	}

	bufferSize := len(buffer)
	offset := 0

	for offset < bufferSize {
		written, err := out.Write(buffer[offset:])
		if err != nil {
			return err
		}

		offset += written
	}

	return nil
}

func (this App) retrieveAllAlbums(out http.ResponseWriter, req *http.Request) {
	albums := this.db.Album.LoadAll()
	muxie.JSON.Dispatch(out, albums)
}

func (this App) upsertAlbum(out http.ResponseWriter, req *http.Request, albumId int64) {
	var err error
	album := model.Album{}
	muxie.JSON.Bind(req, &album)

	album.Id = albumId

	if album.Artist.Id == 0 {
		if album.Artist.Name != "" {
			logrus.Info("Saving off artist with name=", album.Artist.Name, " to ", this.db.Artist)
			// We need to insert the artist, which is apparently new.
			album.Artist.Id, err = this.db.Artist.Save(album.Artist)
			if err != nil {
				out.WriteHeader(http.StatusInternalServerError)
				writeBack(out, err)
				return
			}
		} else {
			// Someone sent us an invalid request.
			out.WriteHeader(http.StatusBadRequest)
			writeBack(out, errors.New("Invalid album artist provided. Name cannot be empty when inserting an artist."))
			return
		}
	}

	album.Id, err = this.db.Album.Save(album)
	if err != nil {
		out.WriteHeader(http.StatusInternalServerError)
		writeBack(out, err)
		return
	}

	if albumId == 0 {
		out.WriteHeader(http.StatusCreated)
		muxie.JSON.Dispatch(out, &album)
	} else {
		out.WriteHeader(http.StatusOK)
	}
}

func (this App) createAlbum(out http.ResponseWriter, req *http.Request) {
	// Ensure that we didn't get an ID sent to us, we are creating a new Album,
	// _not_ updating.
	this.upsertAlbum(out, req, 0)
}

func parseIdFromUrl(out http.ResponseWriter) int64 {
	albumIdStr := muxie.GetParam(out, "id")
	albumId, err := strconv.ParseInt(albumIdStr, 10, 64)
	if err != nil {
		out.WriteHeader(http.StatusBadRequest)
		writeBack(out, errors.New("Invalid ID provided. Must be an integer."))
		return 0
	}

	return albumId
}

func (this App) retrieveAlbum(out http.ResponseWriter, req *http.Request) {
	var albumId int64
	if albumId = parseIdFromUrl(out); albumId == 0 {
		return
	}

	album := this.db.Album.Load(albumId)
	if album == nil {
		out.WriteHeader(http.StatusNotFound)
		writeBack(out, errors.New("Album not found."))
		return
	}
	muxie.JSON.Dispatch(out, album)
}

func (this App) updateAlbum(out http.ResponseWriter, req *http.Request) {
	var albumId int64
	if albumId = parseIdFromUrl(out); albumId == 0 {
		return
	}

	this.upsertAlbum(out, req, albumId)
}

func (this App) removeAlbum(out http.ResponseWriter, req *http.Request) {
	var albumId int64
	if albumId = parseIdFromUrl(out); albumId == 0 {
		return
	}

	_, err := this.db.Album.Delete(model.Album{Id: albumId})
	if err != nil {
		out.WriteHeader(http.StatusInternalServerError)
		writeBack(out, err)
	}
}

func (this App) retrieveArtists(out http.ResponseWriter, req *http.Request) {
	artists := this.db.Artist.LoadAll()
	muxie.JSON.Dispatch(out, artists)
}

func (this App) createArtist(out http.ResponseWriter, req *http.Request) {
	var err error
	artist := model.Artist{}
	muxie.JSON.Bind(req, &artist)

	artist.Id = 0

	if artist.Name == "" {
		// Someone sent us an invalid request.
		out.WriteHeader(http.StatusBadRequest)
		writeBack(out, errors.New("Artists must be named."))
		return
	}

	logrus.Info("Saving off artist with name=", artist.Name, " to ", this.db.Artist)
	// We need to insert the artist, which is apparently new.
	artist.Id, err = this.db.Artist.Save(artist)
	if err != nil {
		out.WriteHeader(http.StatusInternalServerError)
		writeBack(out, err)
		return
	}

	out.WriteHeader(http.StatusCreated)
	muxie.JSON.Dispatch(out, &artist)
}
