package packet

import (
	"errors"
	"io"
)

type SubscribeControlPacket struct {
	// Bits 3,2,1 and 0 of the fixed header of the SUBSCRIBE Control Packet are reserved and MUST be set to 0,0,1 and 0 respectively. The Server MUST treat any other value as malformed and close the Network Connection [MQTT-3.8.1-1].
	// TODO fail packet deserializing when this is not the case
	FixedHeader    FixedHeader
	VariableHeader SubscribeVariableHeader // 2 Bytes
	Payload        SubscribePayload
}

type SubscribeVariableHeader struct {
	PacketID int // int16
}

type SubscribePayload struct {
	Subscriptions []Subscription
}

type Subscription struct {
	Topic string
	QoS   QosLevel
}

func readSubscribeVariableHeader(r io.Reader) (n int, vh SubscribeVariableHeader, err error) {
	packetID, err := readUint16(r)
	if err != nil {
		return 0, SubscribeVariableHeader{}, err
	}

	return 2, SubscribeVariableHeader{PacketID: packetID}, nil
}

func readSubscribePayload(r io.Reader, remainingLength int) (n int, payload SubscribePayload, err error) {
	for n < remainingLength {
		topicLength, err := readUint16(r)
		n += 2 // TODO get this info from readUint16, in case of errors it's maybe not exactly 2
		if err != nil {
			return n, SubscribePayload{}, err
		}

		topic := make([]byte, topicLength)
		bytesRead, err := io.ReadFull(r, topic)
		n += bytesRead
		if err != nil {
			return n, SubscribePayload{}, err
		}

		qos := make([]byte, 1)
		bytesRead, err = io.ReadFull(r, qos)
		n += bytesRead
		if err != nil {
			return n, SubscribePayload{}, err
		}

		sub := Subscription{}
		sub.Topic = string(topic)

		if qos[0]&252 > 0 {
			return n, SubscribePayload{}, errors.New("Invalid Subscribe payload. Reserved bits of QoS are non-zero")
		}

		if qos[0]&1 > 0 && qos[0]&2 > 0 {
			return n, SubscribePayload{}, errors.New("Invalid QoS level in payload. It is not allowed to set both bits")
		}

		if qos[0]&1 > 0 {
			sub.QoS = QoSLevelAtLeastOnce
		} else if qos[0]&2 > 0 {
			sub.QoS = QoSLevelExactlyOnce
		} else {
			sub.QoS = QoSLevelNone
		}
		payload.Subscriptions = append(payload.Subscriptions, sub)
	}
	return
}
