package act_status

import (
	"errors"
	"github.com/Futuremine-chain/futuremine/common/account"
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/futuremine/db/status/act_db"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/utils"
	"github.com/Futuremine-chain/futuremine/types"
	"sync"
)

const account_db = "account_db"

type ActStatus struct {
	db        IActDB
	mutex     sync.RWMutex
	confirmed uint64
}

func NewActStatus() (*ActStatus, error) {
	db, err := act_db.Open(config.App.Setting().Data + "/" + account_db)
	if err != nil {
		return nil, err
	}
	return &ActStatus{db: db}, nil
}

// Initialize account balance root hash
func (a *ActStatus) SetTrieRoot(stateRoot arry.Hash) error {
	return a.db.SetRoot(stateRoot)
}

// Get account status, if the account status needs to be updated
// according to the effective block height, it will be updated,
// but not stored.
func (a *ActStatus) Account(address arry.Address) account.IAccount {
	a.mutex.RLock()
	account := a.db.Account(address)
	a.mutex.RUnlock()

	if account.NeedUpdate() {
		account = a.updateLocked(address)
	}
	return account
}

func (a *ActStatus) Nonce(address arry.Address) uint64 {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	return a.db.Nonce(address)
}

// Update sender account status based on transaction information
func (a *ActStatus) FromTransaction(tx types.ITransaction, height uint64) error {
	if tx.IsCoinBase() {
		return nil
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()

	fromAct := a.db.Account(tx.From())
	err := fromAct.UpdateLocked(a.confirmed)
	if err != nil {
		return err
	}

	err = fromAct.FromTransaction(tx, height)
	if err != nil {
		return err
	}

	a.setAccount(fromAct)
	return nil
}

// Update the receiver's account status based on transaction information
func (a *ActStatus) ToTransaction(tx types.ITransaction, height uint64) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	var toAct account.IAccount

	toAct = a.db.Account(tx.To())
	err := toAct.UpdateLocked(a.confirmed)
	if err != nil {
		return err
	}
	err = toAct.ToTransaction(tx, height)
	if err != nil {
		return err
	}

	a.setAccount(toAct)
	return nil
}

func (a *ActStatus) SetConfirmed(height uint64) {
	a.confirmed = height
}

// Verify the status of the trading account
func (a *ActStatus) Check(tx types.ITransaction) error {
	if tx.Time() > utils.NowUnix() {
		return errors.New("incorrect transaction time")
	}

	account := a.Account(tx.From())
	return account.Check(tx)
}

func (a *ActStatus) Commit() (arry.Hash, error) {
	return a.db.Commit()
}

func (a *ActStatus) Root() arry.Hash {
	return a.db.Root()
}

func (a *ActStatus) Close() error {
	return a.db.Close()
}

func (a *ActStatus) setAccount(account account.IAccount) {
	a.db.SetAccount(account)
}

// Update the locked balance of an account
func (a *ActStatus) updateLocked(address arry.Address) account.IAccount {
	act := a.db.Account(address)
	act.UpdateLocked(a.confirmed)
	return act
}
