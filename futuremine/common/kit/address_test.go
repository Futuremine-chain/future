package kit

import (
	"fmt"
	"github.com/Futuremine-chain/futuremine/common/param"
	"testing"
)

func TestGenerateAddress(t *testing.T) {
	e, _ := Entropy()
	m, _ := Mnemonic(e)
	key, _ := MnemonicToEc(m)
	addr, _ := GenerateAddress(param.MainNet, key.PubKey().SerializeCompressedString())
	fmt.Println(addr)
	if !CheckAddress(param.MainNet, addr) {
		t.Fatalf("failed")
	}
}
