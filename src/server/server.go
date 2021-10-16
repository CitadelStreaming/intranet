package server

import (
    "fmt"
    "net/http"

    "citadel_intranet/src/config"

    "github.com/kataras/muxie"
)

type Server struct {
    server http.Server
    Mux *muxie.Mux
}

func NewServer(cfg config.Config) Server {
    mux := muxie.NewMux()

    mux.Handle("/*file", http.FileServer(http.Dir(cfg.ServerFilePath)))

    server := http.Server{
        Addr: getAddressString(cfg),
        Handler: mux,
    }

    go server.ListenAndServe()

    return Server{
        server: server,
        Mux: mux,
    }
}

func (this Server) Close() {
    this.server.Close()
}

func getAddressString(cfg config.Config) string {
    return fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
}
