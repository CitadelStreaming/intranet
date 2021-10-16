package application_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"citadel_intranet/src/application"
	"citadel_intranet/src/config"
	"citadel_intranet/src/db"
	"citadel_intranet/src/db/dao/mock"
	"citadel_intranet/src/db/model"
	"citadel_intranet/src/server"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	ExpectedJsonForGetAlbums = "[{\"id\":1,\"title\":\"Waffle Irons\",\"artist\":{\"id\":42,\"name\":\"James\"},\"tracks\":[{\"id\":1,\"title\":\"Track 1.1\",\"rating\":5},{\"id\":2,\"title\":\"Track 1.2\",\"rating\":5},{\"id\":3,\"title\":\"Track 1.3\",\"rating\":5}],\"published\":false,\"rating\":5},{\"id\":2,\"title\":\"Something New\",\"artist\":{\"id\":42,\"name\":\"James\"},\"tracks\":[{\"id\":4,\"title\":\"Track 2.1\",\"rating\":5},{\"id\":5,\"title\":\"Track 2.2\",\"rating\":5},{\"id\":6,\"title\":\"Track 2.3\",\"rating\":5}],\"published\":false,\"rating\":3},{\"id\":3,\"title\":\"Something New (Deluxe)\",\"artist\":{\"id\":42,\"name\":\"James\"},\"tracks\":[{\"id\":7,\"title\":\"Track 3.1\",\"rating\":5},{\"id\":8,\"title\":\"Track 3.2\",\"rating\":5},{\"id\":9,\"title\":\"Track 3.3\",\"rating\":5}],\"published\":true,\"rating\":5}]"
)

func TestGetAlbums(t *testing.T) {
	assert := assert.New(t)

	wd, err := os.Getwd()
	assert.Nil(err)
	assert.NotEqual("", wd)

	cfg := config.Config{
		ServerHost:     "",
		ServerPort:     8080,
		ServerFilePath: wd + "/testcontents",
	}
	logrus.Info("Setting web path to: ", cfg.ServerFilePath)

	server := server.NewServer(cfg)

	mockDb, mock, err := sqlmock.New()
	assert.Nil(err)

	mockAlbums := sqlmock.NewRows([]string{"id", "title", "artist", "published", "rating"}).
		AddRow(1, "Waffle Irons", 42, false, 5).
		AddRow(2, "Something New", 42, false, 3).
		AddRow(3, "Something New (Deluxe)", 42, true, 5)
	mock.ExpectQuery(`
        SELECT
            \*
        FROM album
    `).
		WillReturnRows(mockAlbums)
	mock.ExpectQuery(`
        SELECT
            \*
        FROM artist
        WHERE id = ?
    `).
		WithArgs(42).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(42, "James"))
	mockTracks1 := sqlmock.NewRows([]string{"id", "title", "album", "rating"}).
		AddRow(1, "Track 1.1", 1, 5).
		AddRow(2, "Track 1.2", 1, 5).
		AddRow(3, "Track 1.3", 1, 5)
	mock.ExpectQuery(`
        SELECT
            \*
        FROM track
        WHERE album = ?
    `).
		WithArgs(1).
		WillReturnRows(mockTracks1)
	mock.ExpectQuery(`
        SELECT
            \*
        FROM artist
        WHERE id = ?
    `).
		WithArgs(42).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(42, "James"))
	mockTracks2 := sqlmock.NewRows([]string{"id", "title", "album", "rating"}).
		AddRow(4, "Track 2.1", 2, 5).
		AddRow(5, "Track 2.2", 2, 5).
		AddRow(6, "Track 2.3", 2, 5)
	mock.ExpectQuery(`
        SELECT
            \*
        FROM track
        WHERE album = ?
    `).
		WithArgs(2).
		WillReturnRows(mockTracks2)
	mock.ExpectQuery(`
        SELECT
            \*
        FROM artist
        WHERE id = ?
    `).
		WithArgs(42).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(42, "James"))
	mockTracks3 := sqlmock.NewRows([]string{"id", "title", "album", "rating"}).
		AddRow(7, "Track 3.1", 3, 5).
		AddRow(8, "Track 3.2", 3, 5).
		AddRow(9, "Track 3.3", 3, 5)
	mock.ExpectQuery(`
        SELECT
            \*
        FROM track
        WHERE album = ?
    `).
		WithArgs(3).
		WillReturnRows(mockTracks3)

	dbClient := db.NewDatabaseClientFromConnection(mockDb)

	app := application.NewApp(dbClient, server)
	assert.NotNil(app)
	defer app.Close()
	app.Run()

	resp, err := http.Get("http://localhost:8080/api/v1/album")
	assert.Nil(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)
	assert.Equal(ExpectedJsonForGetAlbums, string(body))
}

func TestCreateAlbum(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		ServerHost: "",
		ServerPort: 8080,
	}

	album := model.Album{
		Title: "Something Wicked This Way Comes",
		Artist: model.Artist{
			Id:   1,
			Name: "James",
		},
		Published: false,
		Rating:    0,
	}

	mockAlbumDao := mock.NewMockAlbumDao(ctrl)
	mockAlbumDao.EXPECT().
		Save(gomock.Eq(album)).
		Return(int64(1), nil).
		Times(1)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	assert.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(album)
	assert.Nil(err)

	buffer := bytes.NewBuffer(body)
	resp, err := http.Post("http://localhost:8080/api/v1/album", "application/json", buffer)
	assert.Nil(err)
	defer resp.Body.Close()

	album.Id = 1
	body, err = json.Marshal(album)
	assert.Nil(err)

	retBody, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)
	assert.Equal(string(body), string(retBody))
}

func TestCreateAlbumBadArtist(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		ServerHost: "",
		ServerPort: 8080,
	}

	album := model.Album{
		Title:     "Something Wicked This Way Comes",
		Artist:    model.Artist{},
		Published: false,
		Rating:    0,
	}

	mockAlbumDao := mock.NewMockAlbumDao(ctrl)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	assert.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(album)
	assert.Nil(err)

	buffer := bytes.NewBuffer(body)
	resp, err := http.Post("http://localhost:8080/api/v1/album", "application/json", buffer)
	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)
	assert.Equal("{\"error\":\"Invalid album artist provided. Name cannot be empty when inserting an artist.\"}", string(retBody))
}

func TestCreateAlbumNewArtist(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		ServerHost: "",
		ServerPort: 8080,
	}

	album := model.Album{
		Title: "Something Wicked This Way Comes",
		Artist: model.Artist{
			Name: "James",
		},
		Published: false,
		Rating:    0,
	}
	albumPostEdit := model.Album{
		Title: "Something Wicked This Way Comes",
		Artist: model.Artist{
			Id:   42,
			Name: "James",
		},
		Published: false,
		Rating:    0,
	}

	mockArtistDao := mock.NewMockArtistDao(ctrl)
	mockArtistDao.EXPECT().
		Save(gomock.Eq(album.Artist)).
		Return(int64(42), nil).
		Times(1)
	mockArtistDao.EXPECT().Close().Times(1)

	mockAlbumDao := mock.NewMockAlbumDao(ctrl)
	mockAlbumDao.EXPECT().
		Save(gomock.Eq(albumPostEdit)).
		Return(int64(1), nil).
		Times(1)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(cfg)
	dbClient := db.DatabaseClient{
		Album:  mockAlbumDao,
		Artist: mockArtistDao,
	}

	app := application.NewApp(dbClient, server)
	assert.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(album)
	assert.Nil(err)

	buffer := bytes.NewBuffer(body)
	resp, err := http.Post("http://localhost:8080/api/v1/album", "application/json", buffer)
	assert.Nil(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	defer resp.Body.Close()

	albumPostEdit.Id = 1
	body, err = json.Marshal(albumPostEdit)
	assert.Nil(err)

	retBody, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)
	assert.Equal(string(body), string(retBody))
}

func TestCreateAlbumNewArtistError(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		ServerHost: "",
		ServerPort: 8080,
	}

	album := model.Album{
		Title: "Something Wicked This Way Comes",
		Artist: model.Artist{
			Name: "James",
		},
		Published: false,
		Rating:    0,
	}

	mockArtistDao := mock.NewMockArtistDao(ctrl)
	mockArtistDao.EXPECT().
		Save(gomock.Eq(album.Artist)).
		Return(int64(0), errors.New("Something went bad with the artist")).
		Times(1)
	mockArtistDao.EXPECT().Close().Times(1)

	mockAlbumDao := mock.NewMockAlbumDao(ctrl)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(cfg)
	dbClient := db.DatabaseClient{
		Album:  mockAlbumDao,
		Artist: mockArtistDao,
	}

	app := application.NewApp(dbClient, server)
	assert.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(album)
	assert.Nil(err)

	buffer := bytes.NewBuffer(body)
	resp, err := http.Post("http://localhost:8080/api/v1/album", "application/json", buffer)
	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)
	assert.Equal("{\"error\":\"Something went bad with the artist\"}", string(retBody))
}

func TestCreateAlbumError(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		ServerHost: "",
		ServerPort: 8080,
	}

	album := model.Album{
		Title: "Something Wicked This Way Comes",
		Artist: model.Artist{
			Id:   42,
			Name: "James",
		},
		Published: false,
		Rating:    0,
	}

	mockArtistDao := mock.NewMockArtistDao(ctrl)
	mockArtistDao.EXPECT().Close().Times(1)

	mockAlbumDao := mock.NewMockAlbumDao(ctrl)
	mockAlbumDao.EXPECT().
		Save(gomock.Eq(album)).
		Return(int64(0), errors.New("Waffles")).
		Times(1)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(cfg)
	dbClient := db.DatabaseClient{
		Album:  mockAlbumDao,
		Artist: mockArtistDao,
	}

	app := application.NewApp(dbClient, server)
	assert.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(album)
	assert.Nil(err)

	buffer := bytes.NewBuffer(body)
	resp, err := http.Post("http://localhost:8080/api/v1/album", "application/json", buffer)
	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)
	assert.Equal("{\"error\":\"Waffles\"}", string(retBody))
}

func TestRetrieveAlbum(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		ServerHost: "",
		ServerPort: 8080,
	}

	album := model.Album{
		Id:    123,
		Title: "Something Wicked This Way Comes",
		Artist: model.Artist{
			Id:   42,
			Name: "James",
		},
		Published: false,
		Rating:    0,
	}

	mockAlbumDao := mock.NewMockAlbumDao(ctrl)
	mockAlbumDao.EXPECT().
		Load(gomock.Eq(album.Id)).
		Return(&album).
		Times(1)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	assert.NotNil(app)
	defer app.Close()
	app.Run()

	resp, err := http.Get("http://localhost:8080/api/v1/album/123")
	assert.Nil(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	retAlbum := model.Album{}
	retBody, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)
	assert.Nil(json.Unmarshal(retBody, &retAlbum))
	assert.Equal(album, retAlbum)
}

func TestRetrieveAlbumNotFound(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		ServerHost: "",
		ServerPort: 8080,
	}

	mockAlbumDao := mock.NewMockAlbumDao(ctrl)
	mockAlbumDao.EXPECT().
		Load(gomock.Eq(int64(13))).
		Return(nil).
		Times(1)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	assert.NotNil(app)
	defer app.Close()
	app.Run()

	resp, err := http.Get("http://localhost:8080/api/v1/album/13")
	assert.Nil(err)
	assert.Equal(http.StatusNotFound, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)
	assert.Equal("{\"error\":\"Album not found.\"}", string(retBody))
}

func TestRetrieveAlbumInvalidId(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		ServerHost: "",
		ServerPort: 8080,
	}

	mockAlbumDao := mock.NewMockAlbumDao(ctrl)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	assert.NotNil(app)
	defer app.Close()
	app.Run()

	resp, err := http.Get("http://localhost:8080/api/v1/album/cats")
	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)
	assert.Equal("{\"error\":\"Invalid ID provided. Must be an integer.\"}", string(retBody))
}

func TestUpdateAlbumInvalidId(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		ServerHost: "",
		ServerPort: 8080,
	}

	mockAlbumDao := mock.NewMockAlbumDao(ctrl)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	assert.NotNil(app)
	defer app.Close()
	app.Run()

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/album/cats", nil)
	assert.Nil(err)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)
	assert.Equal("{\"error\":\"Invalid ID provided. Must be an integer.\"}", string(retBody))
}

func TestUpdateAlbum(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		ServerHost: "",
		ServerPort: 8080,
	}

	album := model.Album{
		Id:    456,
		Title: "Something Wicked This Way Comes",
		Artist: model.Artist{
			Id:   42,
			Name: "James",
		},
		Published: false,
		Rating:    0,
	}

	mockAlbumDao := mock.NewMockAlbumDao(ctrl)
	mockAlbumDao.EXPECT().
		Save(gomock.Eq(album)).
		Return(int64(456), nil).
		Times(1)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	assert.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(album)
	assert.Nil(err)

	buffer := bytes.NewBuffer(body)

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/album/456", buffer)
	assert.Nil(err)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	assert.Nil(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
}

func TestRemoveAlbumInvalidId(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		ServerHost: "",
		ServerPort: 8080,
	}

	mockAlbumDao := mock.NewMockAlbumDao(ctrl)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	assert.NotNil(app)
	defer app.Close()
	app.Run()

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/album/cats", nil)
	assert.Nil(err)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)
	assert.Equal("{\"error\":\"Invalid ID provided. Must be an integer.\"}", string(retBody))
}

func TestRemoveAlbum(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		ServerHost: "",
		ServerPort: 8080,
	}

	mockAlbumDao := mock.NewMockAlbumDao(ctrl)
	mockAlbumDao.EXPECT().
		Delete(gomock.Eq(model.Album{Id: 456})).
		Return(int64(1), nil).
		Times(1)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	assert.NotNil(app)
	defer app.Close()
	app.Run()

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/album/456", nil)
	assert.Nil(err)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	assert.Nil(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
}

func TestRemoveAlbumError(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := config.Config{
		ServerHost: "",
		ServerPort: 8080,
	}

	mockAlbumDao := mock.NewMockAlbumDao(ctrl)
	mockAlbumDao.EXPECT().
		Delete(gomock.Eq(model.Album{Id: 456})).
		Return(int64(0), errors.New("Unable to delete album")).
		Times(1)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	assert.NotNil(app)
	defer app.Close()
	app.Run()

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/album/456", nil)
	assert.Nil(err)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	assert.Nil(err)
	assert.Equal(http.StatusInternalServerError, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)
	assert.Equal("{\"error\":\"Unable to delete album\"}", string(retBody))
}
