package etc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
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
		panic(err)
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

	bodyTest := FromString(`hello
		world how are you doing
		`)

	str, err := bodyTest.ReadUntilToken("how")
	if err != nil {
		t.Fatal(err)
	}

	if str != `hello
		world ` {
		t.Fatal("Invalid read", str)
	}

	str, err = bodyTest.ReadUntilToken("you")
	if err != nil {
		t.Fatal(err)
	}

	if str != ` are ` {
		t.Fatal("invalid read x2 \"" + str + "\"")
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
		16428,
		148,
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
		'ß',
		'Ö',
		'イ',
		'ษ',
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
	e.WriteInvertedString(4, "x86")

	for _, v := range strs {
		e.WriteUTF8(v)
	}

	u, _ := ParseUUID("123e4567-e89b-12d3-a456-426655440000")
	e.WriteUUID(u)
	fmt.Println(e.Len())

	tfs := "testing some freaking strings"
	e.Write([]byte(tfs))

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
		sz := runeSize(runes[i])
		run, detected, err := e.ReadRune()
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}

		if detected != sz {
			t.Fatal(string(runes[i]), "expected", sz, "got", detected)
		}

		if run != runes[i] {
			t.Fatal(name, string(runes[i]), "mismatch with", runes[i], "(got ", run, ")")
		}
	}

	enGB := e.ReadInvertedString(4)
	wow := e.ReadInvertedString(4)
	arch := e.ReadBytes(4)

	if enGB != "enGB" && wow != "WoW" {
		t.Fatal("could not handle bliz-style inverted strings (got " + enGB + " " + wow + ")")
	}

	if !bytes.Equal(arch, []byte{'6', '8', 'x', 0}) {
		t.Fatal("inverted string failed")
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
		fmt.Println(spew.Sdump([]byte(tfs)))
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

func TestReflection(t *testing.T) {
	for _, v := range []struct {
		Name  string
		Value interface{}
	}{
		{
			"Float64 buffer",
			struct {
				Key   string
				Value []float64
			}{
				"Some coordinates",
				[]float64{
					33.677216,
					-106.476059,
				},
			},
		},

		{
			"Map encoding",
			map[string]uint64{
				"Sixty Three":                   64,
				"Forty Nine":                    49,
				"Neunzehnhundertneunundsiebzig": 1979,
			},
		},
	} {

		b1 := time.Now()
		fmt.Println("encoding ", v.Name)
		bytes, err := Marshal(v.Value)
		if err != nil {
			t.Fatal(err)
		}
		bench1 := time.Since(b1)
		fmt.Println(spew.Sdump(bytes))

		b2 := time.Now()
		jBytes, err := json.Marshal(v.Value)
		if err != nil {
			t.Fatal(err)
		}

		bench2 := time.Since(b2)

		if bench2 < bench1 {
			fmt.Println("JSON encoding faster for", v.Name)
		} else {
			fmt.Println("etc encoding faster for", v.Name)
		}

		fmt.Println(string(jBytes))
	}

	type TestPtr struct {
		Data *int64
	}

	type TestRecord struct {
		Key      string
		Coords   [2]float32
		Comments []string
		XY       struct {
			X struct {
				Y int64
			}
		}

		Ptr *TestPtr

		Map map[string]int64
	}

	tr := TestRecord{
		Key: "Trinity",
		Coords: [2]float32{
			33.677216,
			-106.476059,
		},
		Comments: []string{
			"hello",
			"world",
		},
	}

	tr.XY.X.Y = 50000
	tr.Ptr = &TestPtr{}
	i := int64(420)
	tr.Ptr.Data = &i
	tr.Map = make(map[string]int64)
	tr.Map["Sixty Four"] = 64

	b, err := Marshal(tr)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("UNMARSHAL")

	var tr2 TestRecord
	err = Unmarshal(b, &tr2)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(spew.Sdump(tr2))

	fixed, err := Marshal(struct {
		X FixedInt64
		Y FixedInt64BE
	}{
		550,
		550,
	})

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("fixed =>", spew.Sdump(fixed))
}

func TestWindowsPlatform(t *testing.T) {
	check := func(t *testing.T, path string, shouldEqual Path) {
		prs := parseWinPath([]rune(path))
		if !reflect.DeepEqual(prs, shouldEqual) {
			t.Fatal("error parsing", path, spew.Sdump([]string(prs)))
		}
	}

	// MSYS-style paths must be parsed correctly
	check(t, "/c/Windows/System32/", Path{"...", "C", "Windows", "System32"})
	check(t, "D:/WeirdPath/Includes Spaces and Forward Slash/", Path{"...", "D", "WeirdPath", "Includes Spaces and Forward Slash"})
}

// func TestDir(t *testing.T) {
// 	fmt.Println(spew.Sdump([]string(Gopath())))
// }
