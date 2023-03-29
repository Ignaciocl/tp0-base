package common

import (
	"encoding/csv"
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
	defer c.CloseConnection()
	users, _ := b.readCsv()
	amountOfUsers := len(users)
	for i := 0; i <= amountOfUsers; i += b.AmountToSend {
		amount := i + b.AmountToSend
		lastBatch := false
		if amount >= amountOfUsers {
			amount = amountOfUsers
			lastBatch = true
		}
		info := make([]byte, 0)
		info = append(info, '[')
		for j := i; j < amount; j +=1 {
			if bUser, err := users[j].ToByteArray(); err != nil {
				return err
			}else {
				info = append(info, bUser...)
			}
		}
		info = append(info, ']')
		if err := c.SendData(info, lastBatch); err != nil {
			return err
		}
		data, err := c.ReceiveData()
		if data == nil || err != nil {
			log.Errorf("could not read data for bingo")
			return err
		}
		var res BingoResponse
		if err := res.ToObject(data); err != nil {
			log.Errorf("could not understand response from otherside, %v, message received was: %v", err, string(data))
			return err
		}
		if amount - i < res.AmountProcessed {
			log.Errorf("processed less than it should, amount expected: %d, amount processed: %d", amount - i, res.AmountProcessed)
			return fmt.Errorf("processed less than it should, amount expected: %d, amount processed: %d", amount - i, res.AmountProcessed)
		}
	}
	log.Infof("action: apuestas_enviadas | result: success | cantidad: %d | status: %s", amountOfUsers, "ok")
	return nil
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
