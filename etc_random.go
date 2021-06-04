package etc

import (
	"crypto/rand"
	"io"
	"math/big"
	"reflect"
	"time"
)

func GenerateRandomUUID() UUID {
	var uuid UUID
	_, err := io.ReadFull(rand.Reader, uuid[:])
	if err != nil {
		return NullUUID
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return uuid
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
