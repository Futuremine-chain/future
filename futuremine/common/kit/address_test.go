package kit

import (
	"fmt"
	"github.com/Futuremine-chain/futuremine/futuremine/common/param"
	"testing"
)

func TestGenerateAddress(t *testing.T) {
	e, _ := Entropy()
	m, _ := Mnemonic(e)
	key, _ := MnemonicToEc(m)
	addr, _ := GenerateAddress(param.TestNet, key.PubKey())
	fmt.Println(addr.String())
	if !CheckAddress(param.TestNet, addr) {
		t.Fatalf("failed")
	}
}
