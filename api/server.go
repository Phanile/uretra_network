package api

import (
	"encoding/gob"
	"github.com/go-kit/log"
	"github.com/labstack/echo/v4"
	"net/http"
	"uretra-network/core"
)

type ServerConfig struct {
	ListenAddr string
	Logger     log.Logger
}

type Server struct {
	txChan chan *core.Transaction
	ServerConfig
	bc *core.Blockchain
}

func NewServer(config ServerConfig, bc *core.Blockchain, txChan chan *core.Transaction) *Server {
	return &Server{
		ServerConfig: config,
		bc:           bc,
		txChan:       txChan,
	}
}

func (s *Server) Start() error {
	e := echo.New()

	e.POST("/tx", s.handlePostTransaction)

	return e.Start(s.ListenAddr)
}

func (s *Server) handlePostTransaction(c echo.Context) error {
	tx := &core.Transaction{}

	err := gob.NewDecoder(c.Request().Body).Decode(tx)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	s.txChan <- tx

	return nil
}
