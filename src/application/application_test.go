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
	"github.com/stretchr/testify/suite"
)

const (
	ExpectedJsonForGetAlbums = "[{\"id\":1,\"title\":\"Waffle Irons\",\"artist\":{\"id\":42,\"name\":\"James\"},\"tracks\":[{\"id\":1,\"title\":\"Track 1.1\",\"album\":1,\"rating\":5},{\"id\":2,\"title\":\"Track 1.2\",\"album\":1,\"rating\":5},{\"id\":3,\"title\":\"Track 1.3\",\"album\":1,\"rating\":5}],\"published\":false,\"rating\":5},{\"id\":2,\"title\":\"Something New\",\"artist\":{\"id\":42,\"name\":\"James\"},\"tracks\":[{\"id\":4,\"title\":\"Track 2.1\",\"album\":2,\"rating\":5},{\"id\":5,\"title\":\"Track 2.2\",\"album\":2,\"rating\":5},{\"id\":6,\"title\":\"Track 2.3\",\"album\":2,\"rating\":5}],\"published\":false,\"rating\":3},{\"id\":3,\"title\":\"Something New (Deluxe)\",\"artist\":{\"id\":42,\"name\":\"James\"},\"tracks\":[{\"id\":7,\"title\":\"Track 3.1\",\"album\":3,\"rating\":5},{\"id\":8,\"title\":\"Track 3.2\",\"album\":3,\"rating\":5},{\"id\":9,\"title\":\"Track 3.3\",\"album\":3,\"rating\":5}],\"published\":true,\"rating\":5}]"
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

type AppSuite struct {
	suite.Suite
	ctrl *gomock.Controller
	cfg  config.Config
}

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppSuite))
}

func (suite *AppSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.cfg = config.Config{
		ServerHost: "",
		ServerPort: 8080,
	}
}

func (suite *AppSuite) TestCreateAlbum() {
	defer suite.ctrl.Finish()

	album := model.Album{
		Title: "Something Wicked This Way Comes",
		Artist: model.Artist{
			Id:   1,
			Name: "James",
		},
		Published: false,
		Rating:    0,
	}

	mockAlbumDao := mock.NewMockAlbumDao(suite.ctrl)
	mockAlbumDao.EXPECT().
		Save(gomock.Eq(album)).
		Return(int64(1), nil).
		Times(1)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(album)
	suite.Nil(err)

	buffer := bytes.NewBuffer(body)
	resp, err := http.Post("http://localhost:8080/api/v1/album", "application/json", buffer)
	suite.Nil(err)
	defer resp.Body.Close()

	album.Id = 1
	body, err = json.Marshal(album)
	suite.Nil(err)

	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal(string(body), string(retBody))
}

func (suite *AppSuite) TestCreateAlbumBadArtist() {
	defer suite.ctrl.Finish()

	album := model.Album{
		Title:     "Something Wicked This Way Comes",
		Artist:    model.Artist{},
		Published: false,
		Rating:    0,
	}

	mockAlbumDao := mock.NewMockAlbumDao(suite.ctrl)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(album)
	suite.Nil(err)

	buffer := bytes.NewBuffer(body)
	resp, err := http.Post("http://localhost:8080/api/v1/album", "application/json", buffer)
	suite.Nil(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal("{\"error\":\"Invalid album artist provided. Name cannot be empty when inserting an artist.\"}", string(retBody))
}

func (suite *AppSuite) TestCreateAlbumNewArtist() {
	defer suite.ctrl.Finish()

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

	mockArtistDao := mock.NewMockArtistDao(suite.ctrl)
	mockArtistDao.EXPECT().
		Save(gomock.Eq(album.Artist)).
		Return(int64(42), nil).
		Times(1)
	mockArtistDao.EXPECT().Close().Times(1)

	mockAlbumDao := mock.NewMockAlbumDao(suite.ctrl)
	mockAlbumDao.EXPECT().
		Save(gomock.Eq(albumPostEdit)).
		Return(int64(1), nil).
		Times(1)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Album:  mockAlbumDao,
		Artist: mockArtistDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(album)
	suite.Nil(err)

	buffer := bytes.NewBuffer(body)
	resp, err := http.Post("http://localhost:8080/api/v1/album", "application/json", buffer)
	suite.Nil(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)
	defer resp.Body.Close()

	albumPostEdit.Id = 1
	body, err = json.Marshal(albumPostEdit)
	suite.Nil(err)

	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal(string(body), string(retBody))
}

func (suite *AppSuite) TestCreateAlbumNewArtistError() {
	defer suite.ctrl.Finish()

	album := model.Album{
		Title: "Something Wicked This Way Comes",
		Artist: model.Artist{
			Name: "James",
		},
		Published: false,
		Rating:    0,
	}

	mockArtistDao := mock.NewMockArtistDao(suite.ctrl)
	mockArtistDao.EXPECT().
		Save(gomock.Eq(album.Artist)).
		Return(int64(0), errors.New("Something went bad with the artist")).
		Times(1)
	mockArtistDao.EXPECT().Close().Times(1)

	mockAlbumDao := mock.NewMockAlbumDao(suite.ctrl)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Album:  mockAlbumDao,
		Artist: mockArtistDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(album)
	suite.Nil(err)

	buffer := bytes.NewBuffer(body)
	resp, err := http.Post("http://localhost:8080/api/v1/album", "application/json", buffer)
	suite.Nil(err)
	suite.Equal(http.StatusInternalServerError, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal("{\"error\":\"Something went bad with the artist\"}", string(retBody))
}

func (suite *AppSuite) TestCreateAlbumError() {
	defer suite.ctrl.Finish()

	album := model.Album{
		Title: "Something Wicked This Way Comes",
		Artist: model.Artist{
			Id:   42,
			Name: "James",
		},
		Published: false,
		Rating:    0,
	}

	mockArtistDao := mock.NewMockArtistDao(suite.ctrl)
	mockArtistDao.EXPECT().Close().Times(1)

	mockAlbumDao := mock.NewMockAlbumDao(suite.ctrl)
	mockAlbumDao.EXPECT().
		Save(gomock.Eq(album)).
		Return(int64(0), errors.New("Waffles")).
		Times(1)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Album:  mockAlbumDao,
		Artist: mockArtistDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(album)
	suite.Nil(err)

	buffer := bytes.NewBuffer(body)
	resp, err := http.Post("http://localhost:8080/api/v1/album", "application/json", buffer)
	suite.Nil(err)
	suite.Equal(http.StatusInternalServerError, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal("{\"error\":\"Waffles\"}", string(retBody))
}

func (suite *AppSuite) TestRetrieveAlbum() {
	defer suite.ctrl.Finish()

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

	mockAlbumDao := mock.NewMockAlbumDao(suite.ctrl)
	mockAlbumDao.EXPECT().
		Load(gomock.Eq(album.Id)).
		Return(&album).
		Times(1)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	resp, err := http.Get("http://localhost:8080/api/v1/album/123")
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	retAlbum := model.Album{}
	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Nil(json.Unmarshal(retBody, &retAlbum))
	suite.Equal(album, retAlbum)
}

func (suite *AppSuite) TestRetrieveAlbumNotFound() {
	defer suite.ctrl.Finish()

	mockAlbumDao := mock.NewMockAlbumDao(suite.ctrl)
	mockAlbumDao.EXPECT().
		Load(gomock.Eq(int64(13))).
		Return(nil).
		Times(1)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	resp, err := http.Get("http://localhost:8080/api/v1/album/13")
	suite.Nil(err)
	suite.Equal(http.StatusNotFound, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal("{\"error\":\"Album not found.\"}", string(retBody))
}

func (suite *AppSuite) TestRetrieveAlbumInvalidId() {
	defer suite.ctrl.Finish()

	mockAlbumDao := mock.NewMockAlbumDao(suite.ctrl)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	resp, err := http.Get("http://localhost:8080/api/v1/album/cats")
	suite.Nil(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal("{\"error\":\"Invalid ID provided. Must be an integer.\"}", string(retBody))
}

func (suite *AppSuite) TestUpdateAlbumInvalidId() {
	defer suite.ctrl.Finish()

	mockAlbumDao := mock.NewMockAlbumDao(suite.ctrl)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/album/cats", nil)
	suite.Nil(err)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	suite.Nil(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal("{\"error\":\"Invalid ID provided. Must be an integer.\"}", string(retBody))
}

func (suite *AppSuite) TestUpdateAlbum() {
	defer suite.ctrl.Finish()

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

	mockAlbumDao := mock.NewMockAlbumDao(suite.ctrl)
	mockAlbumDao.EXPECT().
		Save(gomock.Eq(album)).
		Return(int64(456), nil).
		Times(1)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(album)
	suite.Nil(err)

	buffer := bytes.NewBuffer(body)

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/album/456", buffer)
	suite.Nil(err)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
}

func (suite *AppSuite) TestRemoveAlbumInvalidId() {
	defer suite.ctrl.Finish()

	mockAlbumDao := mock.NewMockAlbumDao(suite.ctrl)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/album/cats", nil)
	suite.Nil(err)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	suite.Nil(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal("{\"error\":\"Invalid ID provided. Must be an integer.\"}", string(retBody))
}

func (suite *AppSuite) TestRemoveAlbum() {
	defer suite.ctrl.Finish()

	mockAlbumDao := mock.NewMockAlbumDao(suite.ctrl)
	mockAlbumDao.EXPECT().
		Delete(gomock.Eq(model.Album{Id: 456})).
		Return(int64(1), nil).
		Times(1)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/album/456", nil)
	suite.Nil(err)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
}

func (suite *AppSuite) TestRemoveAlbumError() {
	defer suite.ctrl.Finish()

	mockAlbumDao := mock.NewMockAlbumDao(suite.ctrl)
	mockAlbumDao.EXPECT().
		Delete(gomock.Eq(model.Album{Id: 456})).
		Return(int64(0), errors.New("Unable to delete album")).
		Times(1)
	mockAlbumDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Album: mockAlbumDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/album/456", nil)
	suite.Nil(err)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	suite.Nil(err)
	suite.Equal(http.StatusInternalServerError, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal("{\"error\":\"Unable to delete album\"}", string(retBody))
}

func (suite *AppSuite) TestGetArtists() {
	defer suite.ctrl.Finish()

	artists := []model.Artist{
		{Id: 1, Name: "James"},
		{Id: 2, Name: "Bobby"},
		{Id: 3, Name: "Jayne"},
	}

	mockArtistDao := mock.NewMockArtistDao(suite.ctrl)
	mockArtistDao.EXPECT().
		LoadAll().
		Return(artists).
		Times(1)
	mockArtistDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Artist: mockArtistDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	resp, err := http.Get("http://localhost:8080/api/v1/artist")
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	retArtists := []model.Artist{}
	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Nil(json.Unmarshal(retBody, &retArtists))
	suite.Equal(artists, retArtists)
}

func (suite *AppSuite) TestCreateArtist() {
	defer suite.ctrl.Finish()

	artist := model.Artist{
		Name: "James",
	}

	mockArtistDao := mock.NewMockArtistDao(suite.ctrl)
	mockArtistDao.EXPECT().
		Save(gomock.Eq(artist)).
		Return(int64(1), nil).
		Times(1)
	mockArtistDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Artist: mockArtistDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(artist)
	suite.Nil(err)

	buffer := bytes.NewBuffer(body)
	resp, err := http.Post("http://localhost:8080/api/v1/artist", "application/json", buffer)
	suite.Nil(err)
	defer resp.Body.Close()

	artist.Id = 1
	body, err = json.Marshal(artist)
	suite.Nil(err)

	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal(string(body), string(retBody))
}

func (suite *AppSuite) TestCreateArtistEmptyName() {
	defer suite.ctrl.Finish()

	artist := model.Artist{}

	mockArtistDao := mock.NewMockArtistDao(suite.ctrl)
	mockArtistDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Artist: mockArtistDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(artist)
	suite.Nil(err)

	buffer := bytes.NewBuffer(body)
	resp, err := http.Post("http://localhost:8080/api/v1/artist", "application/json", buffer)
	suite.Nil(err)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal("{\"error\":\"Artists must be named.\"}", string(retBody))
}

func (suite *AppSuite) TestCreateArtistError() {
	defer suite.ctrl.Finish()

	artist := model.Artist{
		Name: "James",
	}

	mockArtistDao := mock.NewMockArtistDao(suite.ctrl)
	mockArtistDao.EXPECT().
		Save(gomock.Eq(artist)).
		Return(int64(0), errors.New("Wat")).
		Times(1)
	mockArtistDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Artist: mockArtistDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(artist)
	suite.Nil(err)

	buffer := bytes.NewBuffer(body)
	resp, err := http.Post("http://localhost:8080/api/v1/artist", "application/json", buffer)
	suite.Nil(err)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal("{\"error\":\"Wat\"}", string(retBody))
}

func (suite *AppSuite) TestUpdateTrackInvalidId() {
	defer suite.ctrl.Finish()

	mockTrackDao := mock.NewMockTrackDao(suite.ctrl)
	mockTrackDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Track: mockTrackDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/track/cats", nil)
	suite.Nil(err)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	suite.Nil(err)
	suite.Equal(http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal("{\"error\":\"Invalid ID provided. Must be an integer.\"}", string(retBody))
}

func (suite *AppSuite) TestUpdateTrack() {
	defer suite.ctrl.Finish()

	track := model.Track{
		Id:     456,
		Title:  "Something Wicked This Way Comes",
		Rating: 0,
	}

	mockTrackDao := mock.NewMockTrackDao(suite.ctrl)
	mockTrackDao.EXPECT().
		Save(gomock.Eq(track)).
		Return(int64(456), nil).
		Times(1)
	mockTrackDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Track: mockTrackDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(track)
	suite.Nil(err)

	buffer := bytes.NewBuffer(body)

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/track/456", buffer)
	suite.Nil(err)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
}

func (suite *AppSuite) TestUpdateTrackError() {
	defer suite.ctrl.Finish()

	track := model.Track{
		Id:     456,
		Title:  "Something Wicked This Way Comes",
		Rating: 0,
	}

	mockTrackDao := mock.NewMockTrackDao(suite.ctrl)
	mockTrackDao.EXPECT().
		Save(gomock.Eq(track)).
		Return(int64(0), errors.New("Bad day")).
		Times(1)
	mockTrackDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Track: mockTrackDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(track)
	suite.Nil(err)

	buffer := bytes.NewBuffer(body)

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/track/456", buffer)
	suite.Nil(err)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	suite.Nil(err)
	suite.Equal(http.StatusInternalServerError, resp.StatusCode)
	defer resp.Body.Close()

	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal("{\"error\":\"Bad day\"}", string(retBody))
}

func (suite *AppSuite) TestCreateTrack() {
	defer suite.ctrl.Finish()

	track := model.Track{
		Title:  "Something Wicked This Way Comes",
		Rating: 0,
	}

	mockTrackDao := mock.NewMockTrackDao(suite.ctrl)
	mockTrackDao.EXPECT().
		Save(gomock.Eq(track)).
		Return(int64(111), nil).
		Times(1)
	mockTrackDao.EXPECT().Close().Times(1)

	server := server.NewServer(suite.cfg)
	dbClient := db.DatabaseClient{
		Track: mockTrackDao,
	}

	app := application.NewApp(dbClient, server)
	suite.NotNil(app)
	defer app.Close()
	app.Run()

	body, err := json.Marshal(track)
	suite.Nil(err)

	buffer := bytes.NewBuffer(body)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/track", buffer)
	suite.Nil(err)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	suite.Nil(err)
	suite.Equal(http.StatusCreated, resp.StatusCode)
	defer resp.Body.Close()

	track.Id = 111
	body, err = json.Marshal(track)
	suite.Nil(err)

	retBody, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal(string(body), string(retBody))
}
