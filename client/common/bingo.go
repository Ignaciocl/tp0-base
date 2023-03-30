package common

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type BingoService struct {
	ClientInfo BingoDTO
}

func (b BingoService) ProcessInformation(c *Client) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	sigtermWasCalled := false
	go func() {
		<- sigs // If we receive something, sigsterm was called
		sigtermWasCalled = true
	}()
	if err := c.OpenConnection(); err != nil {
		return err
	}
	defer c.CloseConnection()
	if sigtermWasCalled {
		return nil
	}
	info, _ := b.ClientInfo.ToByteArray()
	if err := c.SendData(info); err != nil {
		return err
	}
	if sigtermWasCalled {
		return nil
	}
	data, err := c.ReceiveData()
	if data == nil || err != nil {
		log.Errorf("could not read data for bingo")
		return err
	}
	var res BingoDTO
	if err := res.ToObject(data); err != nil {
		log.Errorf("could not understand response from otherside, %v", err)
		return err
	}
	log.Infof("action: apuesta_enviada | result: success | dni: %d | numero: %d", res.Document, res.Number)
	return nil
}

