package etc

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
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

func (g UUID) UnmarshalJSON(b []byte) error {
	bt := MkBuffer(b)
	bt.ReadByte()
	uidStr := string(bt.ReadBytes(36))
	pu, err := ParseUUID(uidStr)
	if err != nil {
		return err
	}

	g = pu
	return nil
}

func (e *Buffer) WriteUUID(g UUID) {
	e.EncodeUnsignedVarint(g.Big())
}

func (e *Buffer) ReadUUID() UUID {
	ud := e.DecodeUnsignedVarint(20)
	de := ud.Bytes()
	if len(de) != 16 {
		return UUID_NULL
	}

	g := UUID{
		binary.BigEndian.Uint64(de[0:8]),
		binary.BigEndian.Uint64(de[8:16]),
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

// TODO: use bit mask and shift for encoding/decoding to string format.
func (g UUID) String() string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		g.SliceString(0, 4),
		g.SliceString(4, 6),
		g.SliceString(6, 8),
		g.SliceString(8, 10),
		g.SliceString(10, 16))
}

func (g UUID) SliceString(x, y int) string {
	return hex.EncodeToString(g.Bytes()[x:y])
}

func p64(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
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
	part1 := data.ReadBigUint64()

	data = NewBuffer()
	data.Write(clock_seq_hi_and_reserved_clock_seq_low)
	data.Write(node)
	part2 := data.ReadBigUint64()

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

	pr, _ := ParseUUID(string(str))
	return pr
}
