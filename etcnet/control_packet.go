package etcnet

import (
	"time"
)

const (
	// Constants
	CONN_INIT   uint64 = 1
	CONN_ACCEPT uint64 = 2
	CONN_REJECT uint64 = 3
	CONN_DATA   uint64 = 4

	// Control flags
	CONTROL_ACQUIRE_TOKEN     uint8 = 1
	CONTROL_PING              uint8 = 1 << 1
	CONTROL_PAYLOAD           uint8 = 1 << 2
	CONTROL_INITIALIZE_STREAM uint8 = 1 << 3
	CONTROL_HANGUP_STREAM     uint8 = 1 << 4
	CONTROL_ACK               uint8 = 1 << 5

	REJECT_SIGNING    uint8 = 0
	REJECT_SESSION    uint8 = 1
	REJECT_CLOCK      uint8 = 2
	REJECT_NO_SESSION uint8 = 3
)

type ControlPacket struct {
	MessageID uint64
	Flags     uint8
	Time      time.Time

	// Flags & CONTROL_INITIALIZE_STREAM
	InitStreamKey string
	InitStreamID  uint64

	// Flags & CONTROL_INITIALIZE_STREAM
	HangupStreamID uint64

	// Flags & CONTROL_ACK
	Acking uint64

	// Flags & CONTROL_PAYLOAD
	PayloadStreamReference uint64
	PayloadData            []byte
}

func (c ControlPacket) Flagged(f uint8) bool {
	return c.Flags&f != 0
}

func (cp ControlPacket) Parse(b []byte) {
	i := etc.FromBytes(b)

	cp.MessageID = i.ReadUint()
	cp.Flags = i.ReadByte()
	cp.Time = i.ReadTime()

	if cp.Flagged(CONTROL_INITIALIZE_STREAM) {
		cp.InitStreamKey = i.ReadUTF8()
		cp.InitStreamID = i.ReadUint()
	}

	if cp.Flagged(CONTROL_HANGUP_STREAM) {
		cp.HangupStreamID = i.ReadUint()
	}

	if cp.Flagged(CONTROL_ACK) {
		cp.Acking = i.ReadUint()
	}

	if cp.Flagged(CONTROL_PAYLOAD) {
		cp.PayloadStreamReference = i.ReadUint()
		cp.PayloadData = i.ReadRemainder()
	}

	i.Flush()
	i = nil
}

func (cp ControlPacket) Encode() []byte {
	o := etc.NewBuffer()
	o.WriteUint(cp.MessageID)
	o.WriteByte(cp.Flags)
	o.WriteTime(cp.Time)

	if cp.Flagged(CONTROL_INITIALIZE_STREAM) {
		o.WriteUTF8(cp.InitStreamKey)
		o.WriteUint(cp.InitStreamID)
	}

	if cp.Flagged(CONTROL_HANGUP_STREAM) {
		o.WriteUint(cp.HangupStreamID)
	}

	if cp.Flagged(CONTROL_ACK) {
		o.WriteUint(cp.Acking)
	}

	if cp.Flagged(CONTROL_PAYLOAD) {
		o.WriteUint(cp.PayloadStreamReference)
		o.Write(cp.PayloadData)
	}

	return o.Bytes()
}



