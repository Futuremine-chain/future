package types

import (
	"github.com/Futuremine-chain/futuremine/futuremine/common/arry"
	"github.com/Futuremine-chain/futuremine/tools/crypto/ecc/secp256k1"
)

func GenerateAddress(net string, key *secp256k1.PublicKey) (arry.Address, error) {
	return arry.Address{}, nil
}
