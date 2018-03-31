package etc

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestLeb128(t *testing.T) {
	cases := [][]uint64{
		{5000, 69},
	}

	for _, v := range cases {
		t1 := Leb128Encode(v)
		fmt.Println(t1)
		_, t2 := Leb128Decode(len(v), t1)

		fmt.Println(spew.Sdump(t2))
	}

	var finTest uint64 = 0x4245345245
	b := NewBuffer()
	b.Write_LEB128_Uint(finTest)

	dat := b.Read_LEB128_Uint()
	if dat != finTest {
		t.Fatal("Mismatch int read")
	}

	fmt.Println("TEST:", b.Bytes())
}
