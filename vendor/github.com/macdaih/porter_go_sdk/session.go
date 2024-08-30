package portergosdk

import (
	"bytes"
	"fmt"
)

const (
	PasswordMethod string = "Password"
)

type PropType string

const (
	Varint    PropType = "varint"
	Uint32    PropType = "uint32"
	Uint16    PropType = "uint16"
	Byte      PropType = "uint8"
	EncString PropType = "string"
)

type Prop struct {
	key   byte
	value []byte
}

func parseVarintProp(prop *Prop, value uint32) error {
	buff := new(bytes.Buffer)
	if err := encodeVarInt(buff, int(value)); err != nil {
		return err
	}

	prop.value = buff.Bytes()
	return nil
}

func parseU32Prop(prop *Prop, value uint32) error {
	buff := new(bytes.Buffer)
	if err := writeUint32(buff, value); err != nil {
		return err
	}

	prop.value = buff.Bytes()
	return nil
}

func parseU16Prop(prop *Prop, value uint16) error {
	buff := new(bytes.Buffer)
	if err := writeUint16(buff, value); err != nil {
		return err
	}

	prop.value = buff.Bytes()
	return nil
}

func parseByteProp(prop *Prop, value byte) {
	prop.value = []byte{value}
}

func parseStringProp(prop *Prop, value string) error {
	buff := new(bytes.Buffer)
	if err := writeUTFString(buff, value); err != nil {
		return err
	}
	prop.value = buff.Bytes()
	return nil
}

func NewProperty(pt PropType, key byte, value any) (Prop, error) {
	prop := Prop{key: key}

	switch pt {
	case Varint:
		v, ok := value.(uint32)
		if !ok {
			return prop, fmt.Errorf("failed to parse variable integer property : wrong type provided")
		}
		err := parseVarintProp(&prop, v)
		return prop, err
	case Uint32:
		v, ok := value.(uint32)
		if !ok {
			return prop, fmt.Errorf("failed to parse unsigned integer 32 property : wrong type provided")
		}
		err := parseU32Prop(&prop, v)
		return prop, err
	case Uint16:
		v, ok := value.(uint16)
		if !ok {
			return prop, fmt.Errorf("failed to parse unsigned integer 16 property : wrong type provided")
		}
		err := parseU16Prop(&prop, v)
		return prop, err
	case Byte:
		v, ok := value.(byte)
		if !ok {
			return prop, fmt.Errorf("failed to parse byte property : wrong type provided")
		}
		parseByteProp(&prop, v)
		return prop, nil
	case EncString:
		v, ok := value.(string)
		if !ok {
			return prop, fmt.Errorf("failed to parse string property : wrong type provided")
		}
		err := parseStringProp(&prop, v)
		return prop, err
	default:
		return prop, fmt.Errorf(
			"unknown property type provided %s",
			string(pt),
		)
	}
	return prop, nil
}

type Session struct {
	clientID   string
	keepAlive  uint16
	authMethod string

	cleanStart bool
	will       bool
	retain     bool
	pwdFlag    bool
	usrFlag    bool
	qos        uint8

	usr *string
	pwd *string
}
