package private

import (
	"encoding/json"
	"fmt"
	"github.com/Futuremine-chain/futuremine/futuremine/common/keystore"
	"github.com/Futuremine-chain/futuremine/futuremine/common/kit"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/crypto/ecc/secp256k1"
	"github.com/Futuremine-chain/futuremine/tools/crypto/mnemonic"
	"github.com/Futuremine-chain/futuremine/tools/utils"
	"io/ioutil"
)

var defaultPrivateFile = "key.json"
var defaultPrivatePass = "fmc"

type Private struct {
	address  arry.Address
	priKey   *secp256k1.PrivateKey
	mnemonic string
}

func NewPrivate() *Private {
	return &Private{}
}

func CreatePrivate(network string) (*Private, error) {
	entropy, err := mnemonic.Entropy()
	if err != nil {
		return nil, fmt.Errorf("failed to create entropy, %s", err.Error())
	}
	mnemonicStr, err := mnemonic.Mnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("failed to create mnemonic, %s", err.Error())
	}
	key, err := mnemonic.MnemonicToEc(mnemonicStr)
	if err != nil {
		return nil, fmt.Errorf("failed to mnemonic to ec, %s", err.Error())
	}
	address, err := kit.GenerateAddress(network, key.PubKey())
	if err != nil {
		return nil, err
	}
	return &Private{address, key, mnemonicStr}, nil
}

func (p *Private) PrivateKey() *secp256k1.PrivateKey {
	return p.priKey
}

func (p *Private) Load(file string, key string) error {
	if !utils.Exist(file) {
		return fmt.Errorf("%s is not exsists", file)
	}
	j, err := keystore.ReadJson(file)
	if err != nil {
		return fmt.Errorf("read json file %s failed! %s", file, err.Error())
	}
	privJson, err := keystore.DecryptPrivate([]byte(key), j)
	if err != nil {
		return fmt.Errorf("decrypt priavte failed! %s", err.Error())
	}
	privKey, err := secp256k1.ParseStringToPrivate(privJson.Private)
	if err != nil {
		return fmt.Errorf("parse priavte failed! %s", err.Error())
	}
	p.address = arry.StringToAddress(j.Address)
	p.priKey = privKey
	p.mnemonic = privJson.Mnemonic
	return nil
}

func (p *Private) Create(net, file, key string) error {
	private, err := CreatePrivate(net)
	if err != nil {
		return err
	}
	j, err := keystore.PrivateToJson(net, private.priKey, private.mnemonic, []byte(key))
	if err != nil {
		return fmt.Errorf("key json creation failed! %s", err.Error())
	}
	bytes, _ := json.Marshal(j)
	if err = ioutil.WriteFile(file, bytes, 0644); err != nil {
		return fmt.Errorf("write json file %s failed! %s", file, err.Error())
	}
	p.address = private.address
	p.priKey = private.priKey
	p.mnemonic = private.mnemonic
	return nil
}

func (p *Private) Serialize() []byte {
	// TODO
	return nil
}
