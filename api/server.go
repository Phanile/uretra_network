package api

import (
	"encoding/gob"
	"encoding/hex"
	"github.com/go-kit/log"
	"github.com/labstack/echo/v4"
	"net/http"
	"uretra-network/core"
	"uretra-network/types"
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

type GetBalanceResponse struct {
	Balance uint64 `json:"balance"`
	Error   string `json:"error"`
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
	e.GET("/getBalance/:address", s.handleGetBalance)

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

func (s *Server) handleGetBalance(c echo.Context) error {
	hexAddr := c.Param("address")
	addrBytes, err := hex.DecodeString(hexAddr)

	resp := GetBalanceResponse{}

	if err != nil {
		resp.Error = err.Error()
		return c.JSON(http.StatusBadRequest, resp)
	}

	balance, errBalance := s.bc.GetAccounts().GetBalance(types.AddressFromBytes(addrBytes))

	if errBalance != nil {
		resp.Error = errBalance.Error()
		return c.JSON(http.StatusInternalServerError, resp)
	}

	resp.Balance = balance
	return c.JSON(http.StatusOK, resp)
}
