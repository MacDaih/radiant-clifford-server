package portergosdk

import (
	"bytes"
	"fmt"
)

func writeUint16(buff *bytes.Buffer, in uint16) error {
	if err := buff.WriteByte(
		byte((in & 0xff00) >> 8),
	); err != nil {
		return err
	}

	return buff.WriteByte(
		byte((in & 0x00ff)),
	)
}

func readUint16(input []byte) (uint16, error) {
	if len(input) < 2 {
		return 0, fmt.Errorf("failed to read unsigned 16 bytes integer : invalid length")
	}

	return uint16(input[0]<<8 ^ input[1]), nil
}

func writeUint32(buff *bytes.Buffer, in uint32) error {

	if err := buff.WriteByte(byte((in & 0xff000000) >> 24)); err != nil {
		return err
	}

	if err := buff.WriteByte(byte((in & 0x00ff0000) >> 16)); err != nil {
		return err
	}

	if err := buff.WriteByte(byte((in & 0x0000ff00) >> 8)); err != nil {
		return err
	}

	return buff.WriteByte(byte(in & 0x000000ff))
}

func writeUTFString(buff *bytes.Buffer, str string) error {
	length := len(str)

	if err := writeUint16(buff, uint16(length)); err != nil {
		return err
	}

	_, err := buff.Write([]byte(str))
	return err
}

func readUTFString(str []byte) (string, error) {
	strlen, err := readUint16(str)
	if err != nil {
		return "", err
	}

	if int(strlen) >= len(str) {
		return "", fmt.Errorf("failed to read string : prefix length greater than actual length")
	}

	return string(str[2 : strlen+2]), nil
}

func encodeVarInt(buff *bytes.Buffer, input int) error {
	if input == 0 {
		return buff.WriteByte(0)
	}

	var enc byte
	count := 0
	for input > 0 && count < 5 {
		enc = byte(input % 128)
		input /= 128
		if input > 0 {
			enc |= 0x80
		}
		if err := buff.WriteByte(enc); err != nil {
			return err
		}
		count++
	}

	if count >= 5 {
		return fmt.Errorf("malformed byte")
	}

	return nil
}

func decodeVarint(input []byte) (uint32, error) {
	var remainingMult byte = 1
	var (
		lword byte
		b     byte
		lb    int
	)

	for i := 0; i < 4; i++ {
		lb++
		b = input[i]
		lword += (b & 127) * remainingMult
		remainingMult *= 128
		if (b & 128) == 0 {
			if lb > 1 && b == 0 {
				return 0, fmt.Errorf("malformed packet")
			}
			return uint32(lword), nil
		}
	}

	return 0, fmt.Errorf("malformed packet")
}

func evalBytes(value uint32) int {
	if value < 128 {
		return 1
	} else if value < 16384 {
		return 2
	} else if value < 2097152 {
		return 3
	} else if value < 268435456 {
		return 4
	} else {
		return 5
	}
}
