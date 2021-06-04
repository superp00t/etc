package etc

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

type UUID [16]byte

var NullUUID UUID

func (g UUID) MarshalJSON() ([]byte, error) {
	b := NewBuffer()
	fmt.Fprintf(b, "\"")
	fmt.Fprintf(b, g.String())
	fmt.Fprintf(b, "\"")

	return b.Bytes(), nil
}

func (g UUID) MarshalYAML() (interface{}, error) {
	return g.String(), nil
}

func (g UUID) UnmarshalYAML(b interface{}) error {
	var err error

	g, err = ParseUUID(b.(string))
	return err
}

func (g UUID) UnmarshalJSON(b []byte) error {
	bt := FromBytes(b)
	bt.ReadByte()
	uidStr := string(bt.ReadBytes(36))
	pu, err := ParseUUID(uidStr)
	if err != nil {
		return err
	}

	bt.ReadByte()
	g = pu
	return nil
}

func (e *Buffer) WriteUUID(g UUID) {
	e.Write(g[:])
}

func (e *Buffer) ReadUUID() UUID {
	var g UUID
	e.Read(g[:])
	return g
}

func (g UUID) Bytes() []byte {
	return []byte(g[:])
}

func (g UUID) Cmp(u UUID) bool {
	return bytes.Equal(g[:], u[:])
}

func (g UUID) String() string {
	data := g.Bytes()
	d := hex.EncodeToString(data)
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		d[0:8],
		d[8:12],
		d[12:16],
		d[16:20],
		d[20:32])
}

func p64(i uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, i)
	return b
}

func l64(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func ParseUUID(s string) (UUID, error) {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	ok := r.MatchString(s)
	if !ok {
		return NullUUID, fmt.Errorf("etc: invalid UUID '%s'", s)
	}

	bytes := []byte(strings.ReplaceAll(s, "-", ""))

	var u UUID
	if _, err := hex.Decode(u[:], bytes); err != nil {
		return u, nil
	}

	return u, nil
}
