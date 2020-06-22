package kit

import (
	"bytes"
	"errors"
	"github.com/Futuremine-chain/futuremine/common/param"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/crypto/base58"
	"github.com/Futuremine-chain/futuremine/tools/crypto/ecc/secp256k1"
	"github.com/Futuremine-chain/futuremine/tools/crypto/hash"
)

const addressLength = 35
const addressBytesLength = 26

func GenerateAddress(net string, key *secp256k1.PublicKey) (arry.Address, error) {
	ver := []byte{}
	switch net {
	case param.MainNet:
		ver = append(ver, param.MainNetParam.PubKeyHashAddrID[0:]...)
	case param.TestNet:
		ver = append(ver, param.TestNetParam.PubKeyHashAddrID[0:]...)
	default:
		return arry.Address{}, errors.New("wrong network")
	}
	hashed256 := hash.Hash(key.SerializeCompressed())
	hashed160, err := hash.Hash160(hashed256.Bytes())
	if err != nil {
		return arry.Address{}, err
	}

	addNet := append(ver, hashed160...)
	hashed1 := hash.Hash(addNet)
	hashed2 := hash.Hash(hashed1.Bytes())
	checkSum := hashed2[0:4]
	hashedCheck1 := append(addNet, checkSum...)
	return arry.StringToAddress(base58.Encode(hashedCheck1)), nil
}

func CheckAddress(net string, address arry.Address) bool {
	ver := []byte{}
	switch net {
	case param.MainNet:
		ver = append(ver, param.MainNetParam.PubKeyHashAddrID[0:]...)
	case param.TestNet:
		ver = append(ver, param.TestNetParam.PubKeyHashAddrID[0:]...)
	default:
		return false
	}
	addr := address.String()
	if len(addr) != addressLength {
		return false
	}
	addrBytes := base58.Decode(addr)
	if len(addrBytes) != addressBytesLength {
		return false
	}
	checkSum := addrBytes[len(addrBytes)-4:]
	checkBytes := addrBytes[0 : len(addrBytes)-4]
	checkBytesHashed1 := hash.Hash(checkBytes)
	checkBytesHashed2 := hash.Hash(checkBytesHashed1.Bytes())
	netBytes := checkBytes[0:2]
	if bytes.Compare(ver, netBytes) != 0 {
		return false
	}
	return bytes.Compare(checkSum, checkBytesHashed2[0:4]) == 0
}
