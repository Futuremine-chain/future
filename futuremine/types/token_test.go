package types

import (
	"fmt"
	"testing"
)

func TestTokenRecord_Check(t1 *testing.T) {
	amount1 := 4699999999999999999
	amount2 := 4699999999999999999
	amount3 := Amount(amount1 + amount2)
	fmt.Println(amount3)
	fmt.Println(amount3.ToCoin())
}
