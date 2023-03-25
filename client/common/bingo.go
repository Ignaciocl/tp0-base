package common

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type BingoService struct {
	ClientInfo BingoDTO
	Id int
	AmountToSend int
}

func (b BingoService) ProcessInformation(c *Client) error {
	if err := c.OpenConnection(); err != nil {
		return err
	}
	users, _ := b.readCsv()
	amountOfUsers := len(users)
	for i := 0; i <= amountOfUsers; i += b.AmountToSend {
		amount := i + b.AmountToSend
		lastBatch := false
		if amount >= amountOfUsers {
			amount = amountOfUsers
			lastBatch = true
		}
		msg := BingoCommunication{
			Action: "sendingBatch",
			Data:   users[i:amount],
		}
		res, err := SendAndReceive(c, msg, lastBatch)
		if err != nil {
			return err
		}
		if amount-i < res.AmountProcessed {
			log.Errorf("processed less than it should, amount expected: %d, amount processed: %d", amount-i, res.AmountProcessed)
			return fmt.Errorf("processed less than it should, amount expected: %d, amount processed: %d", amount-i, res.AmountProcessed)
		}
	}
	c.CloseConnection()
	log.Infof("action: apuestas_enviadas | result: success | cantidad: %d | status: %s", amountOfUsers, "ok")
	for {
		if err := c.OpenConnection(); err != nil {
			return err
		}
		msg := BingoCommunication{
			Action: "findMeMyOgre",
			Data:   nil,
		}
		if res, err := SendAndReceive(c, msg, true); err != nil {
			log.Errorf("error while receiving winners: %v", err)
			return err
		} else if res.Status == "foundOgre" {
			log.Infof("winners are: %v", res.Winners)
			break
		}
		c.CloseConnection()

	}
	return nil
}

func SendAndReceive(c *Client, msg BingoCommunication, lastBatch bool) (BingoResponse, error) {
	info, _ := json.Marshal(msg)
	if err := c.SendData(info, lastBatch); err != nil {
		return BingoResponse{}, err
	}
	return getResponse(c)
}

func getResponse(c *Client) (BingoResponse, error) {
	data, err := c.ReceiveData()
	if data == nil || err != nil {
		log.Errorf("could not read data for bingo")
		return BingoResponse{}, err
	}
	var res BingoResponse
	if err := json.Unmarshal(data, &res); err != nil {
		log.Errorf("could not understand response from otherside, %v, message received was: %v", err, string(data))
		return BingoResponse{}, err
	}
	return res, nil
}

func (b BingoService) readCsv() ([]BingoDTO, error) {
	f, err := os.Open(fmt.Sprintf("/dataset/agency-%d.csv", b.Id))
	if err != nil {
		log.Errorf("could not open file %s", fmt.Sprintf("/dataset/agency-%d", b.Id))
		return nil, err
	}
	defer f.Close()

	records, _ := csv.NewReader(f).ReadAll()
	data := make([]BingoDTO, 0)
	for i, row := range records {
		name := row[0]
		surname := row[1]
		document, errDoc := strconv.Atoi(row[2])
		if errDoc != nil {
			log.Errorf("invalid register with index %d, document not a number %v", i, errDoc)
			return nil, errDoc
		}
		bornDate := row[3]
		number, errNumber := strconv.Atoi(row[4])
		if errDoc != nil {
			log.Errorf("invalid register with index %d, Number not a number %v", i, errNumber)
			return nil, errDoc
		}
		r := BingoDTO {
			Name: name,
			Document: document,
			BornDate: bornDate,
			Number: number,
			Surname: surname,
		}
		data = append(data, r)
	}
	return data, nil
}
