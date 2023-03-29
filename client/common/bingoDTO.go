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

type SendableMessage interface {
	ToByteArray() ([]byte, error)
	ToObject([]byte) error
}

func (b *BingoDTO) ToByteArray() ([]byte, error) {
	sep := "\","
	message := "{" + "\"name\":\"" + b.Name + sep + "\"document\":\"" + strconv.Itoa(b.Document) + sep + "\"born_date\":\"" + b.BornDate + sep + "\"number\":\"" + strconv.Itoa(b.Number) + sep + "\"surname\":\"" + b.Surname + "\"}"
	return []byte(message), nil
}

func findNextSep(message string, position int) (string, int) {
	untilSeparator := ""
	for i := position; i <= len(message); i += 1 {
		char := string(message[i])
		if char == "," || char == "}" {
			return untilSeparator, i - position
		}
		if char == "\"" {
			continue
		}
		untilSeparator += char
	}
	return "", 0
}

func (b *BingoDTO) adder(message string, wholeMessage string, position int) (int, bool, error) {
	posToAdd := 0
	interString := ""
	var err error
	if message == "\"name\":" {
		b.Name, posToAdd = findNextSep(wholeMessage, position + 1)
		return posToAdd, true, nil
	}
	if message == "\"document\":" {
		interString, posToAdd = findNextSep(wholeMessage, position + 1)
		b.Document, err = strconv.Atoi(interString)
		if err != nil {
			return 0, false, errors.Errorf("document isn't number, can't translate class, err: %v", err)
		}
		return posToAdd, true, err
	}
	if message == "\"born_date\":" {
		b.BornDate, posToAdd = findNextSep(wholeMessage, position + 1)
		return posToAdd, true, nil
	}
	if message == "\"number\":" {
		interString, posToAdd = findNextSep(wholeMessage, position + 1)
		b.Number, err = strconv.Atoi(interString)
		if err != nil {
			return 0, false, errors.Errorf("Number isn't number, can't translate class, err: %v", err)
		}
		return posToAdd, true, err
	}
	if message == "\"surname\":" {
		b.Surname, posToAdd = findNextSep(wholeMessage, position + 1)
		return posToAdd, true, nil
	}
	return 0, false, nil
}

func (b *BingoDTO) ToObject(bytes []byte) error {
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
