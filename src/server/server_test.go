package server_test

import (
    "io/ioutil"
    "net/http"
	"os"
	"testing"

	"citadel_intranet/src/config"
	"citadel_intranet/src/server"

	"github.com/stretchr/testify/assert"
	"github.com/sirupsen/logrus"
)

func TestServerSetup(t *testing.T) {
	assert := assert.New(t)

    wd, err := os.Getwd()
    assert.Nil(err)
    assert.NotEqual("", wd)

    cfg := config.Config{
        ServerHost: "",
        ServerPort: 8080,
        ServerFilePath: wd + "/testcontents",
    }
    logrus.Info("Setting web path to: ", cfg.ServerFilePath)

    server := server.NewServer(cfg)

    resp, err := http.Get("http://localhost:8080/index.html")
    assert.Nil(err)
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    assert.Nil(err)
    assert.Equal(`<!DOCTYPE html>
<html>
    <body>
        Hello Tests!
    </body>
</html>
`, string(body))

    server.Close()
}
