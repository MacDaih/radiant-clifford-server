package portergosdk

import (
	"bytes"
	"slices"
)

func buildSubscribe(
	topics []string,
	pktID uint16,
) ([]byte, error) {
	var (
		msg,
		idBuff,
		propBuff,
		payloadBuff bytes.Buffer
	)

	if err := msg.WriteByte(SubscribeCMD); err != nil {
		return nil, err
	}

	if err := writeUint16(&idBuff, pktID); err != nil {
		return nil, err
	}

	// TODO deal with props
	if err := encodeVarInt(&propBuff, 0); err != nil {
		return nil, err
	}

	// payload
	for _, topic := range topics {
		if err := writeUTFString(&payloadBuff, topic); err != nil {
			return nil, err
		}

		// TODO handle subscription options
		if err := payloadBuff.WriteByte(0); err != nil {
			return nil, err
		}
	}

	if err := encodeVarInt(
		&msg,
		(idBuff.Len() + propBuff.Len() + payloadBuff.Len()),
	); err != nil {
		return nil, err
	}

	if _, err := msg.Write(
		slices.Concat(
			idBuff.Bytes(),
			propBuff.Bytes(),
			payloadBuff.Bytes(),
		),
	); err != nil {
		return nil, err
	}

	return msg.Bytes(), nil
}
