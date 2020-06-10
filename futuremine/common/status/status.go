package status

import (
	"github.com/Futuremine-chain/futuremine/common/account"
	"github.com/Futuremine-chain/futuremine/common/dpos"
	"github.com/Futuremine-chain/futuremine/common/token"
)

type FMCStatus struct {
	actStatus   account.IActStatus
	dPosStatus  dpos.IDPosStatus
	tokenStatus token.ITokenStatus
}

func NewFMCStatus(actStatus account.IActStatus, dPosStatus dpos.IDPosStatus, tokenStatus token.ITokenStatus) *FMCStatus {
	return &FMCStatus{
		actStatus:   actStatus,
		dPosStatus:  dPosStatus,
		tokenStatus: tokenStatus,
	}
}
