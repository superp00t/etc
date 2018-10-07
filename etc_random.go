package etc

import (
	"crypto/rand"
	"math/big"
	"reflect"
)

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

//Random random intergers
func RandomBigInt(min, max *big.Int) *big.Int {
	range_ := new(big.Int).Sub(max, min)
	bi, err := rand.Int(rand.Reader, range_)
	if err != nil {
		panic(err)
	}
	return new(big.Int).Add(min, bi)
}

func RandomIndex(v interface{}) int {
	if reflect.TypeOf(v).Kind() == reflect.Slice {
		s := reflect.ValueOf(v)
		bi := RandomBigInt(big.NewInt(0), big.NewInt(int64(s.Len())))
		return int(bi.Uint64())
	}

	return 0
}

func GenerateMiniGUID() uint64 {
	return 0
}
