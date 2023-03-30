package common

import (
	"bytes"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
	ClosingMessage string
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Fatalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met

func (c *Client) OpenConnection() error {
	if err := c.createClientSocket(); err != nil {
		log.Errorf("error while openning connection, %v", err)
		return err
	}
	return nil
}

func (c *Client) CloseConnection() {
	c.conn.Close()
}

func (c *Client) SendData(bytes []byte) error {
	bytesToSend := append(bytes, c.config.ClosingMessage...)
	eightKB := 8 * 1024
	size := len(bytesToSend)
	for i := 0; i <= len(bytesToSend); i += eightKB {
		var sending []byte
		if size < i + eightKB {
			sending = bytesToSend[i : size]
		} else {
			sending = bytesToSend[i: i + eightKB]
		}
		amountSent, err := c.conn.Write(sending)
		if err != nil {
			log.Errorf("weird error happened, not stopping but something should be checked: %v", err)
		}
		if dif := len(sending) - amountSent; dif > 0 { // Avoiding short write
			i -= dif
		}
	}
	return nil
}

func (c *Client) ReceiveData() ([]byte, error) {
	eightKB := 8 * 1024
	received := make([]byte, eightKB)
	checkedValue := []byte(c.config.ClosingMessage)
	total := make([]byte, 0)
	for {
		if i, err := c.conn.Read(received); err != nil {
			log.Errorf("error while receiving message, ending receiver: %v", err)
			return nil, err
		} else {
			total = append(total, received[0:i]...)
		}
		if bytes.HasSuffix(total, checkedValue) {
			break
		}
	}
	finalData := total[0:len(total) - len(c.config.ClosingMessage)]
	return finalData, nil
}
