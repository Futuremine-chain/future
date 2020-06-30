package kit

import (
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"testing"
)

func TestGenerateTokenAddress(t *testing.T) {
	token, err := GenerateTokenAddress("mainnet", arry.StringToAddress("vQqgdUAFeBXdEfd4FMWbg65sXuWnims8hrt"), "ABC")
	if err != nil {
		t.Fatal(err)
	}
	if !CheckTokenAddress("mainnet", token, "ABC") {
		t.Fatal("check failed")
	}
}
