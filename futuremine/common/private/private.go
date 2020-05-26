package private

import "github.com/Futuremine-chain/futuremine/tools/crypto/ecc/secp256k1"

type Private struct{
	Address    hasharry.Address
	PrivateKey *secp256k1.PrivateKey
	Mnemonic   string
}
