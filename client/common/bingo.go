package common

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type BingoService struct {
	ClientInfo BingoDTO
}

func (b BingoService) ProcessInformation(c *Client) error {
	if err := c.OpenConnection(); err != nil {
		return err
	}
	defer c.CloseConnection()
	info, _ := json.Marshal(b.ClientInfo)
	if err := c.SendData(info); err != nil {
		return err
	}
	data, err := c.ReceiveData()
	if data == nil || err != nil {
		log.Errorf("could not read data for bingo")
		return err
	}
	var res BingoDTO
	if err := json.Unmarshal(data, &res); err != nil {
		log.Errorf("could not understand response from otherside, %v", err)
		return err
	}
	log.Infof("action: apuesta_enviada | result: success | dni: %d | numero: %d", res.Document, res.Number)
	return nil
}

