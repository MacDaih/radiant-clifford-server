package portergosdk

const (
	MQTT = "MQTT"

	Version5 = 5
)

type CodeString string

const (
	CodeConnack    CodeString = "connack"
	CodePublish    CodeString = "publish"
	CodeSubAck     CodeString = "suback"
	CodePubAck     CodeString = "puback"
	CodeDisconnect CodeString = "disconnect"
	CodeUnknown    CodeString = "unknown"
)

const (
	ConnectCMD    byte = 0x10
	ConnackCMD    byte = 0x20
	PublishCMD    byte = 0x30
	SubackCMD     byte = 0x80
	DisconnectCMD byte = 0xe0
	PingReqCMD    byte = 0xC0
	SubscribeCMD  byte = 0x80
)

func parseCode(in byte) CodeString {
	switch in {
	case ConnackCMD:
		return CodeConnack
	case PublishCMD:
		return CodePublish
	case SubackCMD:
		return CodeSubAck
	case DisconnectCMD:
		return CodeDisconnect
	default:
		return CodeUnknown
	}
}
