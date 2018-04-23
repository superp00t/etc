package etc

import (
	"fmt"
	"testing"
)

func TestBuffer(t *testing.T) {

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
		0000000000,
		9999999999,
	}
	var signtest = []int64{
		123123,
		-50000,
		6459,
		100000000000,
		-53564644234,
	}

	e := NewBuffer()

	for _, v := range testdata {
		e.WriteUint(v)
	}

	for _, v := range tstfloats {
		e.WriteFloat32(v)
	}

	for _, v := range signtest {
		e.WriteInt(v)
	}

	u, _ := ParseUUID("123e4567-e89b-12d3-a456-426655440000")
	e.WriteUUID(u)
	fmt.Println(e.Len())

	tfs := "testityeah"
	e.WriteFixedString(10, tfs)

	for i := 0; i < len(testdata); i++ {
		ca := e.ReadUint()
		if testdata[i] != ca {
			t.Fatal("mismatch with", testdata[i], "(got ", ca, ")")
		}
	}

	for i := 0; i < len(tstfloats); i++ {
		ca := e.ReadFloat32()
		if tstfloats[i] != ca {
			t.Fatal("mismatch with", tstfloats, "(got ", ca, ")")
		}
	}

	for i := 0; i < len(signtest); i++ {
		ca := e.ReadInt()
		fmt.Println(ca)
		if signtest[i] != ca {
			t.Fatal("mismatch with", signtest[i], "(got ", ca, ")")
		}
	}

	uid := e.ReadUUID()
	if !uid.Cmp(u) {
		t.Fatal("UUID mismatch")
	}

	str := e.ReadRemainder()
	if string(str) != tfs {
		t.Fatal("Invalid string " + string(str))
	}

	ff := NewBuffer()
	ff.Write([]byte{
		49,
		50,
	})

	ff.ReadBytes(2)

	if ff.Available() != 0 {
		t.Fatal("Invalid available", ff.Available())
	}
}
