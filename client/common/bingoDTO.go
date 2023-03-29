package common

import (
	"github.com/pkg/errors"
	"strconv"
)

type BingoDTO struct {
	Name string `json:"name"`
	Document int `json:"document"`
	BornDate string `json:"born_date"`
	Number int `json:"number"`
	Surname string `json:"surname"`
}

type BingoCommunication struct {
	Action string `json:"action"`
	Data []BingoDTO `json:"data,omitempty"`
}

type BingoResponse struct {
	AmountProcessed int `json:"amount_processed"`
	Status string `json:"status"`
	Winners []string `json:"winners,omitempty"`
}

func (r *BingoResponse) adder(message string, wholeMessage string, position int) (int, bool, error) {
	posToAdd := 1
	interString := ""
	var err error
	if message == "\"status\":" {
		r.Status, posToAdd, _ = findNextSep(wholeMessage, position + 1)
		return posToAdd, true, nil
	}
	if message == "\"document\":" {
		interString, posToAdd, _ = findNextSep(wholeMessage, position + 1)
		r.AmountProcessed, err = strconv.Atoi(interString)
		if err != nil {
			return 0, false, errors.Errorf("document isn't number, can't translate class, err: %v", err)
		}
		return posToAdd, true, err
	}
	if message == "\"winners\":[" {
		winners := make([]string, 0)
		for s, pos, finish := "", 0, false; !finish; {
			s, pos, finish = findNextSep(wholeMessage, position + posToAdd)
			posToAdd += pos
			if s == " " || s == "" {
				posToAdd += 1
				continue
			}
			winners = append(winners, s)
		}
		r.Winners = winners
		return posToAdd, true, nil
	}
	return posToAdd, false, err
}

type SendableMessage interface {
	ToByteArray() ([]byte, error)
	ToObject([]byte) error
}

type adderable interface {
	adder(message string, wholeMessage string, position int) (int, bool, error)
}

func (b *BingoDTO) ToByteArray() ([]byte, error) {
	sep := "\","
	message := "{" + "\"name\":\"" + b.Name + sep + "\"document\":\"" + strconv.Itoa(b.Document) + sep + "\"born_date\":\"" + b.BornDate + sep + "\"number\":\"" + strconv.Itoa(b.Number) + sep + "\"surname\":\"" + b.Surname + "\"}"
	return []byte(message), nil
}

func findNextSep(message string, position int) (string, int, bool) {
	untilSeparator := ""
	for i := position; i < len(message); i += 1 {
		char := string(message[i])
		if char == "," || char == "}" || char == "]" {
			return untilSeparator, i - position, char == "]"
		}
		if char == "\"" {
			continue
		}
		untilSeparator += char
	}
	return "", 0, false
}

func (b *BingoDTO) adder(message string, wholeMessage string, position int) (int, bool, error) {
	posToAdd := 0
	interString := ""
	var err error
	if message == "\"name\":" {
		b.Name, posToAdd, _ = findNextSep(wholeMessage, position + 1)
		return posToAdd, true, nil
	}
	if message == "\"document\":" {
		interString, posToAdd, _ = findNextSep(wholeMessage, position + 1)
		b.Document, err = strconv.Atoi(interString)
		if err != nil {
			return 0, false, errors.Errorf("document isn't number, can't translate class, err: %v", err)
		}
		return posToAdd, true, err
	}
	if message == "\"born_date\":" {
		b.BornDate, posToAdd, _ = findNextSep(wholeMessage, position + 1)
		return posToAdd, true, nil
	}
	if message == "\"number\":" {
		interString, posToAdd, _ = findNextSep(wholeMessage, position + 1)
		b.Number, err = strconv.Atoi(interString)
		if err != nil {
			return 0, false, errors.Errorf("Number isn't number, can't translate class, err: %v", err)
		}
		return posToAdd, true, err
	}
	if message == "\"surname\":" {
		b.Surname, posToAdd, _ = findNextSep(wholeMessage, position + 1)
		return posToAdd, true, nil
	}
	return 0, false, nil
}

func (b *BingoDTO) ToObject(bytes []byte) error {
	return processBytes(bytes, b)
}

func processBytes(bytes []byte, b adderable) error {
	receivedInfo := string(bytes)
	message := ""
	for i := 0; i < len(receivedInfo); i += 1 {
		char := string(receivedInfo[i])
		if char == "{" || char == "," {
			message = ""
			continue
		}
		if char == "}" {
			break
		}
		message += char
		if finalPosition, added, err := b.adder(message, receivedInfo, i); added {
			i += finalPosition
			message = ""
		} else if err != nil {
			return err
		}
	}
	return nil
}

func (r *BingoResponse) ToObject(bytes []byte) error  {
	return processBytes(bytes, r)
}

func (b *BingoCommunication) ToByteArray() ([]byte, error) {
	message := "{" + "\"action\":\"" + b.Action + "\""
	if b.Data != nil {
		message += "\"data\":["
		for i := 0; i < len(b.Data); i+=1 {
			data, _ := b.Data[i].ToByteArray()
			message += string(data)
			if i != len(b.Data) -1 {
				message += ","
			}
		}
		message += "]"
	}
	message += "}"
	return []byte(message), nil
}
