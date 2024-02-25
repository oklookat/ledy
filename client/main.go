package client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oklookat/ledy/effect"
)

func New() *Client {
	return &Client{}
}

type Client struct {
	conn *websocket.Conn
}

func (c *Client) Connect() error {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	finderCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer func() {
		cancel()
	}()
	data, err := findService(finderCtx)
	if err != nil {
		return err
	}
	if data == nil {
		return errors.New("server not found")
	}

	connURL := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", data.IP, data.Port), Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(connURL.String(), nil)
	if err != nil {
		return err
	}
	c.conn = conn

	return nil
}

func (c *Client) SetColors(leds effect.LEDS) error {
	if c.conn == nil {
		return nil
	}
	cmd := newCommandSetColors(leds)
	return c.conn.WriteMessage(websocket.BinaryMessage, cmd[:])
}

func (c *Client) SetColorCorrection(v ColorCorrection) error {
	if c.conn == nil {
		return nil
	}
	cmd := newCommandSetColorCorrection(v)
	return c.conn.WriteMessage(websocket.BinaryMessage, cmd[:])
}

func (c *Client) SetColorTemperature(v ColorTemperature) error {
	if c.conn == nil {
		return nil
	}
	cmd := newCommandSetColorTemperature(v)
	return c.conn.WriteMessage(websocket.BinaryMessage, cmd[:])
}
