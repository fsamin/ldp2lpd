package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/url"
	"time"

	"golang.org/x/net/websocket"
)

var pingCodec = websocket.Codec{Marshal: ping, Unmarshal: nil}

func ping(v interface{}) (msg []byte, payloadType byte, err error) {
	return nil, websocket.PingFrame, nil
}

var timeUnit time.Duration
var verbose bool

func main() {
	c := getConf()
	u, err := url.Parse(c.Address)
	if err != nil {
		log.Fatal(err)
	}

	var unit time.Duration
	var factor = time.Duration(c.TimeFactor)
	switch c.TimeUnit {
	case "millisecond":
		unit = time.Millisecond
	case "second":
		unit = time.Second
	case "minute":
		unit = time.Minute
	case "hour":
		unit = time.Hour
	default:
		log.Fatal("unsupported time unit")
	}
	timeUnit = factor * unit

	verbose = c.Verbose

	for {
		// Try to connect
		log.Printf("Connecting to %s...\n", u.Host)

		ws, err := websocket.Dial(u.String(), "", "http://mySelf")
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Connected!")

		go display(context.Background(), globalCounters)
		go rolling(context.Background(), globalCounters)

		var msg []byte
		for {
			ws.SetReadDeadline(time.Now().Add(5 * time.Second))
			err := websocket.Message.Receive(ws, &msg)
			if t, ok := err.(net.Error); ok && t.Timeout() {
				// Timeout, send a Pong && continue
				pingCodec.Send(ws, nil)
				continue
			}

			if err != nil {
				log.Printf("Error while reading from %q: %q. Will try to reconnect after 1s...\n", u.Host, err.Error())
				time.Sleep(1 * time.Second)
				break
			}

			// Extract Message
			var logMessage struct {
				Message string `json:"message"`
			}
			json.Unmarshal(msg, &logMessage)

			// Extract infos
			var message message
			json.Unmarshal([]byte(logMessage.Message), &message)

			count(message)
		}
	}
}
