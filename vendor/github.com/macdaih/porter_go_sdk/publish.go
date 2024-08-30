package portergosdk

import (
	"fmt"
)

type ContentType string

const (
	Json = "json"
	Text = "text"
)

type AppMessage struct {
	TopicName   string
	Format      bool
	ContentType ContentType
	Correlation string
	SubID       string
	Payload     []byte
}

func readPublish(b []byte) (AppMessage, error) {
	var msg AppMessage
	cursor := 1
	// TODO publis cmd with retain, dup and qos
	// Remaining length
	length, err := decodeVarint(b[1:])
	if err != nil {
		return msg, err
	}

	if len(b) < int(length) {
		return msg, fmt.Errorf("malformed packet : invalid length")
	}

	cursor += evalBytes(length)
	// Read topic
	topic, err := readUTFString(b[cursor:])
	if err != nil {
		return msg, err
	}
	msg.TopicName = topic

	cursor += len(topic) + 2
	// No packet ID for now

	//read props
	propsLen, err := decodeVarint(b[cursor:])
	if err != nil {
		return msg, err
	}

	ceil := cursor + int(propsLen)
	for cursor < ceil {
		if cursor > int(length) {
			return msg, fmt.Errorf("malformed packet : cursor exceeded length")
		}
		switch b[cursor] {
		case 0x01:
			cursor++
			// payload format indicator
			indicator := b[cursor]
			msg.Format = indicator >= 1
			cursor++
		case 0x03:
			cursor++
			// content type
			content, err := readUTFString(b[cursor:])
			if err != nil {
				return msg, err
			}
			msg.ContentType = ContentType(content)
		default:
			cursor++
			continue
		}
	}

	if msg.Format {
		payload, err := readUTFString(b[cursor+1:])
		//msg.Payload = b[cursor+1 : length-1]
		if err != nil {
			return msg, err
		}
		msg.Payload = []byte(payload)
	}

	return msg, nil
}
