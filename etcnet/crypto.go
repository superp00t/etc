package etcnet

import (
	"github.com/superp00t/etc"
	"golang.org/x/crypto/curve25519"
)

func _GENERATE_TOKEN() uint64 {

}

func _GEN_DH_KEYPAIR() (*[32]byte, *[32]byte) {
	pub := new([32]byte)
	sec := new([32]byte)

	e := etc.NewBuffer()
	e.WriteRandom(32)
	copy(sec[:], e.Bytes())

	curve25519.ScalarBaseMult(pub, sec)
	return pub, sec
}

func _ENCRYPT_AND_AUTH(key, data []byte) (*[24]byte, []byte) {
	bx := box.Seal()
}

func _DH_AGREEMENT(theirpub, mypriv *[32]byte) *[32]byte {
	out := new([32]byte)

	curve25519.ScalarMult(out, mypriv, theirpub)

	return out
}



