package portergosdk

import (
	"bytes"
	"slices"
)

func getConnectLength(s Session, props []Prop) int {
	return 0
}

func CreateConnectPacket(s Session, length int) ([]byte, error) {
	return nil, nil
}

func buildConnect(
	cid string,
	keepAlive uint16,
	creds *credential,
) ([]byte, error) {
	// make connect packet
	var (
		msg,
		vhBuff,
		lenBuff,
		idBuff,
		usrBuff,
		pwdBuff bytes.Buffer
	)

	if err := msg.WriteByte(ConnectCMD); err != nil {
		return nil, err
	}

	n, err := vhBuff.Write([]byte{
		0x00, 0x04, 0x4d, 0x51, 0x54, 0x54, 0x05,
	})
	if err != nil {
		return nil, err
	}

	props := make([]Prop, 0, 7)

	var flag uint8 = 0

	if creds != nil {
		authProp, err := NewProperty(
			EncString,
			0x15,
			creds.authMethod,
		)
		if err != nil {
			return nil, err
		}
		props = append(props, authProp)

		if creds.usr != nil {
			flag ^= 0x80
			if err := writeUTFString(&usrBuff, *creds.usr); err != nil {
				return nil, err
			}
			n += usrBuff.Len()
		}

		if creds.pwd != nil {
			flag ^= 0x40
			if err := writeUTFString(&pwdBuff, *creds.pwd); err != nil {
				return nil, err
			}
			n += pwdBuff.Len()
		}
	}

	// Var header flag
	if err := vhBuff.WriteByte(flag); err != nil {
		return nil, err
	}
	n++

	if err := writeUint16(&vhBuff, keepAlive); err != nil {
		return nil, err
	}
	n += 2

	var propLenBuff bytes.Buffer
	var propBuff bytes.Buffer

	for _, prop := range props {
		if err := propBuff.WriteByte(prop.key); err != nil {
			return nil, err
		}
		n++

		nv, err := propBuff.Write(prop.value)
		if err != nil {
			return nil, err
		}
		n += nv
	}

	if err := encodeVarInt(&propLenBuff, propBuff.Len()); err != nil {
		return nil, err
	}

	n += propBuff.Len()

	// id
	if err := writeUTFString(&idBuff, cid); err != nil {
		return nil, err
	}
	n += idBuff.Len()

	// usr & pwd
	if ul := usrBuff.Len(); ul > 0 {
		n += ul
	}

	if pl := pwdBuff.Len(); pl > 0 {
		n += pl
	}

	// Encode packet remaining length
	if err := encodeVarInt(&lenBuff, n); err != nil {
		return nil, err
	}

	whole := slices.Concat(
		lenBuff.Bytes(),
		vhBuff.Bytes(),
		propLenBuff.Bytes(),
		propBuff.Bytes(),
		idBuff.Bytes(),
		usrBuff.Bytes(),
		pwdBuff.Bytes(),
	)

	if _, err := msg.Write(whole); err != nil {
		return nil, err
	}

	return msg.Bytes(), nil
}
