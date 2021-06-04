package etc

import (
	"fmt"
	"testing"
)

func TestUUID(t *testing.T) {
	ex := "123e4567-e89b-12d3-a456-426655440000"
	g, err := ParseUUID(ex)
	if err != nil {
		t.Fatal(err)
	}

	if g.String() != ex {
		t.Fatal("UUID string mismatch", g.String(), ex)
	}

	fmt.Println(g.String())

	fmt.Println(GenerateRandomUUID())
}
