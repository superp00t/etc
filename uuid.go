package etc

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"
)

type UUID [2]uint64

var UUID_NULL UUID = UUID{
	0x0000000000000000,
	0x0000000000000000,
}

func (g UUID) Big() *big.Int {
	return big.NewInt(0).SetBytes(g.Bytes())
}

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
	bt := MkBuffer(b)
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
	e.WriteUint64(g[0])
	e.WriteUint64(g[1])
}

func (e *Buffer) ReadUUID() UUID {
	g := UUID{
		e.ReadUint64(),
		e.ReadUint64(),
	}

	return g
}

func (g UUID) Bytes() []byte {
	return append(
		p64(g[0]),
		p64(g[1])...,
	)
}

func (g UUID) Cmp(u UUID) bool {
	return g[0] == u[0] && g[1] == u[1]
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

func (g UUID) TimeStitch() uint64 {
	return (uint64(g.Get_time_low()) | uint64(g.Get_time_mid())>>32 | uint64(g.Get_time_hi_and_version())&0xFF>>34)
}

func (g UUID) Time() time.Time {
	return time.Unix(0, int64(g.TimeStitch())*int64(time.Millisecond))
}

func (g UUID) Get_time_low() uint32 {
	return uint32(g[0] & 0xFFFFFFFF)
}

func (g UUID) Get_time_mid() uint16 {
	return uint16((g[0] >> 32) & 0xFFFF)
}

func (g UUID) Get_time_hi_and_version() uint16 {
	return uint16((g[0] >> 48) & 0xFFFF)
}

func (g UUID) Get_clock_seq_hi_and_res_clock_seq_low() uint16 {
	return uint16(g[1] & 0xFFFF)
}

func (g UUID) Get_node() uint64 {
	return uint64(g[1] >> 16)
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
		return UUID_NULL, fmt.Errorf("etc: invalid UUID '%s'", s)
	}

	chrs := strings.Split(s, "-")
	time_low := l64(chrs[0])[:4]
	time_mid := l64(chrs[1])[:2]
	time_hi_and_version := l64(chrs[2])[:2]
	clock_seq_hi_and_reserved_clock_seq_low := l64(chrs[3])[:2]
	node := l64(chrs[4])[:6]

	data := NewBuffer()
	data.Write(time_low)
	data.Write(time_mid)
	data.Write(time_hi_and_version)
	part1 := data.ReadUint64()

	data = NewBuffer()
	data.Write(clock_seq_hi_and_reserved_clock_seq_low)
	data.Write(node)
	part2 := data.ReadUint64()

	return UUID{part1, part2}, nil
}

func GenerateRandomUUID() UUID {
	by := NewBuffer()
	by.WriteRandom(16)

	part1 := by.ReadBigUint64()
	part2 := by.ReadBigUint64()

	gu := UUID{part1, part2}
	str := []rune(gu.String())

	str[14] = '4'
	str[19] = '8'

	pr, err := ParseUUID(string(str))
	if err != nil {
		panic(err)
	}
	return pr
}
