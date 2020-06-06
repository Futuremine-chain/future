package private

import (
	"fmt"
	"github.com/Futuremine-chain/futuremine/futuremine/common/arry"
	"github.com/Futuremine-chain/futuremine/futuremine/common/keystore"
	"github.com/Futuremine-chain/futuremine/tools/crypto/ecc/secp256k1"
	"github.com/Futuremine-chain/futuremine/tools/utils"
)

const defaultKey = "fmc"

type Private struct {
	address  arry.Address
	priKey   *secp256k1.PrivateKey
	mnemonic string
}

func NewPrivate() *Private {
	return &Private{}
}

func (p *Private) PrivateKey() *secp256k1.PrivateKey {
	return p.priKey
}

func (p *Private) Load(path string) {

}

func LoadNodePrivate(file string, key string) (*Private, error) {
	if !utils.Exist(file) {
		return nil, fmt.Errorf("%s is not exsists", file)
	}
	j, err := keystore.ReadJson(file)
	if err != nil {
		return nil, fmt.Errorf("read json file %s failed! %s", file, err.Error())
	}
	privJson, err := keystore.DecryptPrivate([]byte(key), j)
	if err != nil {
		return nil, fmt.Errorf("decrypt priavte failed! %s", err.Error())
	}
	privKey, err := secp256k1.ParseStringToPrivate(privJson.Private)
	if err != nil {
		return nil, fmt.Errorf("parse priavte failed! %s", err.Error())
	}
	return &Private{arry.StringToAddress(j.Address), privKey, privJson.Mnemonic}, nil
}
