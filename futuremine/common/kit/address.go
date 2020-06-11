package kit

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/crypto/ecc/secp256k1"
)

func GenerateAddress(net string, key *secp256k1.PublicKey) (arry.Address, error) {
	return arry.Address{}, nil
}

func CheckAddress(net string, address arry.Address) bool {
	return false
}
