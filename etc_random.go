package etc

import (
	"crypto/rand"
	"math/big"
	"reflect"
	"time"
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

func RandomInt(min, max int) int {
	return int(RandomInt64(int64(min), int64(max)))
}

func RandomInt64(min, max int64) int64 {
	bi := RandomBigInt(big.NewInt(min), big.NewInt(max))
	return bi.Int64()
}

func RandomDuration(min, max time.Duration) time.Duration {
	return time.Duration(RandomInt64(int64(min), int64(max)))
}

func RandomIndex(v interface{}) int {
	if reflect.TypeOf(v).Kind() == reflect.Slice {
		s := reflect.ValueOf(v)
		return RandomInt(0, s.Len())
	}

	return 0
}
