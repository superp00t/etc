package etc

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func rndm() string {
	b := NewBuffer()
	b.WriteRandom(20)
	return hex.EncodeToString(b.Bytes())
}

func getRandomFile() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("APPDATA") + "\\" + rndm()
	}

	return "/tmp/" + rndm()
}

func TestBuffer(t *testing.T) {
	t1, err := FileController(getRandomFile())
	if err != nil {
		t.Fatal(err)
	}

	t2 := NewBuffer()
	buffertest("file", t1, t)
	buffertest("memory", t2, t)

	t3 := NewBuffer()
	t3.WriteFixedString(4, "test")
	t3.WriteUint(12345678)

	fmt.Println(t3.Base64())
	fmt.Println(spew.Sdump(t3.Bytes()))

	buf := FromString("a whisper goes around the world...")

	var ot [35]byte
	i, err := buf.Read(ot[:])
	if err != io.EOF || i != 34 {
		t.Fatal(err, i, "could not read correctly", spew.Sdump(ot[:]))
	}

	g := NewBuffer()
	g.WriteUint(10)
	g.WriteUint(50)
	g.WriteUint(14)
	g.WriteUint(16428)
	g.WriteUint(148)

	fmt.Println(g.Rpos())
	ten := g.ReadUint()
	fmt.Println(g.Rpos())
	fifty := g.ReadUint()
	fmt.Println(g.Rpos())
	fourteen := g.ReadUint()

	if ten != 10 {
		t.Fatal("ten error")
	}

	if fifty != 50 {
		t.Fatal("fifty error")
	}

	if fourteen != 14 {
		t.Fatal("fourteen error", fourteen, g.Bytes())
	}

	if c := g.ReadUint(); c != 16428 {
		t.Fatal("leb err", c)
	}

	if g.ReadUint() != 148 {
		t.Fatal("leb err 2")
	}

}

func buffertest(name string, e *Buffer, t *testing.T) {
	bnch := time.Now()
	var tstfloats = []float32{
		6.022140857,
		12312315.34,
		0,
		3245365.342342,
	}
	var testdata = []uint64{
		0,
		123123,
		69420,
		6022140857,
		0xDEADBEEF,
		8982343454,
		0,
		18446744073709551557,
		0000000000,
		9999999999,
		18446744073709551615,
		10,
		50,
		14,
		0,
		16428,
	}
	var signtest = []int64{
		9223372036854775807,
		123123,
		-9223372036854775808,
		-50000,
		6459,
		1025484282308132903,
		1233333333345353,
		1844673709551557,
		100000000000,
		-53564644234,
	}
	var runes = []rune{
		'™',
		'🖕',
		'𝕬',
	}
	var strs = []string{
		"fffff",
		"Quoth the raven, nevermore",
		"Memes",
		"שני בוטנים צעדו ברחוב. אחד הותקף.",
		"Два арахиса шли по улице. Один из них подвергся нападению.",
		"一つの妖怪がヨーロッパにあらわれている、――共産主義の妖怪が。",
	}

	e.WriteFixedString(10, "testing")

	for _, v := range testdata {
		e.WriteUint(v)
	}

	for _, v := range tstfloats {
		e.WriteFloat32(v)
	}

	for _, v := range signtest {
		e.WriteInt(v)
	}

	for _, v := range runes {
		e.WriteRune(v)
	}

	e.WriteInvertedString(4, "enGB")
	e.WriteInvertedString(4, "WoW")

	for _, v := range strs {
		e.WriteUTF8(v)
	}

	u, _ := ParseUUID("123e4567-e89b-12d3-a456-426655440000")
	e.WriteUUID(u)
	fmt.Println(e.Len())

	tfs := "testityeah"
	e.WriteFixedString(10, tfs)

	if stt := e.ReadFixedString(10); stt != "testing" {
		t.Fatal(name, "invalid fixed string: "+stt, spew.Sdump([]byte(stt)))
	}

	fmt.Println(e.Rpos(), e.backend.Size(), spew.Sdump(e.Bytes()))
	for i := 0; i < len(testdata); i++ {
		ca := e.ReadUint()
		if testdata[i] != ca {
			t.Fatal(name, "mismatch with", testdata[i], "(got ", ca, ")")
		}
	}

	for i := 0; i < len(tstfloats); i++ {
		ca := e.ReadFloat32()
		if tstfloats[i] != ca {
			t.Fatal(name, "mismatch with", tstfloats, "(got ", ca, ")")
		}
	}

	for i := 0; i < len(signtest); i++ {
		ca := e.ReadInt()
		fmt.Println(ca)
		if signtest[i] != ca {
			t.Fatal(name, "mismatch with", signtest[i], "(got ", ca, ")")
		}
	}

	for i := 0; i < len(runes); i++ {
		run, _, _ := e.ReadRune()
		if run != runes[i] {
			t.Fatal(name, "mismatch with", runes[i], "(got ", run, ")")
		}
	}

	enGB := e.ReadInvertedString(4)
	wow := e.ReadInvertedString(4)

	if enGB != "enGB" && wow != "WoW" {
		t.Fatal("could not handle bliz-style inverted strings (got " + enGB + " " + wow + ")")
	}

	for i := 0; i < len(strs); i++ {
		s := e.ReadUTF8()
		if s != strs[i] {
			t.Fatal(name, "mismatch with", strs[i], "(got ", s, ")")
		}
	}

	uid := e.ReadUUID()
	if !uid.Cmp(u) {
		t.Fatal(name, "UUID mismatch")
	}

	str := e.ReadRemainder()
	cap := strings.TrimRight(string(str), "\x00")
	if cap != tfs {
		t.Fatal(name, "Invalid string '"+cap+"'", spew.Sdump(str))
	}

	ff := NewBuffer()
	ff.Write([]byte{
		49,
		50,
	})

	ff.ReadBytes(2)

	if ff.Available() != 0 {
		t.Fatal(name, "Invalid available", ff.Available())
	}

	fmt.Println(time.Since(bnch))
}
