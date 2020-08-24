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
	addr, _ := GenerateAddress(param.TestNet, key.PubKey().SerializeCompressedString())
	fmt.Println(addr)
	if !CheckAddress(param.TestNet, addr) {
		t.Fatalf("failed")
	}
	fmt.Println(GenerateTokenAddress(param.TestNet, addr, "SA"))
}
