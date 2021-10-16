package application

import (
	"encoding/json"
	"errors"
	"net/http"

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

}

func (this App) Close() {
	this.server.Close()
	this.db.Close()
}

func writeBack(out *http.ResponseWriter, value interface{}) error {
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
		written, err := (*out).Write(buffer[offset:])
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

func (this App) createAlbum(out http.ResponseWriter, req *http.Request) {
	var err error
	album := model.Album{}
	muxie.JSON.Bind(req, &album)

	// Ensure that we didn't get an ID sent to us, we are creating a new Album,
	// _not_ updating.
	album.Id = 0

	if album.Artist.Id == 0 {
		if album.Artist.Name != "" {
			logrus.Info("Saving off artist with name=", album.Artist.Name, " to ", this.db.Artist)
			// We need to insert the artist, which is apparently new.
			album.Artist.Id, err = this.db.Artist.Save(album.Artist)
			if err != nil {
				out.WriteHeader(http.StatusInternalServerError)
				writeBack(&out, err)
				return
			}
		} else {
			// Someone sent us an invalid request.
			out.WriteHeader(http.StatusBadRequest)
			writeBack(&out, errors.New("Invalid album artist provided. Name cannot be empty when inserting an artist."))
			return
		}
	}

	album.Id, err = this.db.Album.Save(album)
	if err != nil {
		out.WriteHeader(http.StatusInternalServerError)
		writeBack(&out, err)
		return
	}

	out.WriteHeader(http.StatusCreated)
	muxie.JSON.Dispatch(out, &album)
}
