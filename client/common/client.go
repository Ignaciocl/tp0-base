package common

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
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
	ClosingBatch string
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
func (c *Client) StartClientLoop() {
	// autoincremental msgID to identify every message sent
	msgID := 1
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

loop:
	// Send messages if the loopLapse threshold has not been surpassed
	for timeout := time.After(c.config.LoopLapse); ; {
		select {
		case <-timeout:
			log.Infof("action: timeout_detected | result: success | client_id: %v",
				c.config.ID,
			)
			break loop
		case <-sigs:
			log.Infof("sigterm was called, finished execution")
			break loop
		default:
		}

		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		// TODO: Modify the send to avoid short-write
		fmt.Fprintf(
			c.conn,
			"[CLIENT %v] Message N°%v\n",
			c.config.ID,
			msgID,
		)
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		msgID++
		c.conn.Close()

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}
		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
			c.config.ID,
			msg,
		)

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) OpenConnection() error {
	for {
		if err := c.createClientSocket(); err != nil {
			log.Errorf("error while openning connection, %v", err)
			return err
		}
		msgReceived := make([]byte, 4)
		if receivedBytes, err := c.conn.Read(msgReceived); err == nil {
			msg := string(msgReceived[0:receivedBytes])
			if receivedBytes == 3 && string(msg) == "ack" {
				log.Infof("connection successful")
				break
			}
			if msg == "nack" {
				log.Info("nack received, waiting a second to ask again")
			} else {
				log.Infof("message: '%s' received, did not understand, waiting a second to try again", msg)
			}
		} else {
			log.Errorf("error while waiting for receive new connection, closing previous and waiting for new: %v", err)
		}
		_ = c.conn.Close()
		c.conn = nil
		time.Sleep(1)
	}
	return nil
}

func (c *Client) CloseConnection() {
	c.conn.Close()
}

func (c *Client) SendData(bytes []byte, lastBatch bool) error {
	closingMessage := c.config.ClosingMessage
	if lastBatch {
		closingMessage = c.config.ClosingBatch
	}
	bytesToSend := append(bytes, closingMessage...)
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
			log.Errorf("weird error happened, stopping but something should be checked: %v", err)
			return err
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
