package etc

import (
	"fmt"
	"testing"
)

func TestBuffer(t *testing.T) {

	var tstfloats = []float32{
		6.022140857,
		12312315.34,
		3245365.342342,
	}
	var testdata = []uint64{
		6022140857,
		0xDEADBEEF,
		8982343454,
		0000000000,
		9999999999,
	}

	e := NewBuffer()

	for _, v := range testdata {
		e.Write_LEB128_Uint(v)
	}

	for _, v := range tstfloats {
		e.WriteFloat32(v)
	}

	u, _ := ParseUUID("123e4567-e89b-12d3-a456-426655440000")
	e.WriteUUID(u)
	fmt.Println(e.Len())

	for i := 0; i < len(testdata); i++ {
		ca := e.Read_LEB128_Uint()
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

	uid := e.ReadUUID()
	if !uid.Cmp(u) {
		t.Fatal("UUID mismatch")
	}
}
