package txlist

import (
	"fmt"
	"github.com/Futuremine-chain/futuremine/common/account"
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/common/db/txlist"
	"github.com/Futuremine-chain/futuremine/common/validator"
	"github.com/Futuremine-chain/futuremine/types"
	"sync"
)

const txList_db = "tx_List_db"
const maxPoolTx = 100000

type TxManagement struct {
	cache     *Cache
	ready     *Sorted
	validator validator.IValidator
	actStatus account.IActStatus
	mutex     sync.RWMutex
	txDB      ITxListDB
}

func NewTxManagement(validator validator.IValidator, actStatus account.IActStatus) (*TxManagement, error) {
	txDB, err := txlist.Open(config.App.Setting().Data + "/" + txList_db)
	if err != nil {
		return nil, err
	}
	return &TxManagement{
		cache:     NewCache(txDB),
		ready:     NewSorted(txDB),
		validator: validator,
		actStatus: actStatus,
		txDB:      txDB,
	}, nil
}

func (t *TxManagement) Read() error {
	txs := t.txDB.Read()
	if txs != nil {
		for _, tx := range txs {
			t.Put(tx)
		}
	}
	return nil
}

func (t *TxManagement) Close() error {
	t.txDB.Close()
	return nil
}

func (t *TxManagement) Count() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.cache.Len() + t.ready.Len()
}

func (t *TxManagement) Put(tx types.ITransaction) error {
	if t.Exist(tx) {
		return fmt.Errorf("the transaction %s already exists", tx.Hash().String())
	}
	if err := t.validator.Check(tx); err != nil {
		return err
	}

	if t.cache.Len() >= maxPoolTx {
		t.DeleteEnd(tx)
	}
	return t.put(tx)
}

func (t *TxManagement) put(tx types.ITransaction) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	from := tx.From().String()
	nonce := t.actStatus.Nonce(tx.From())
	if nonce == tx.Nonce()-1 {
		oldTx := t.ready.GetByAddress(from)
		if oldTx != nil {
			if oldTx.Nonce() == tx.Nonce() && oldTx.Fee() < tx.Fee() {
				t.ready.Remove(oldTx)
			} else if oldTx.Nonce() < tx.Nonce() {
				t.ready.Remove(oldTx)
			} else if oldTx.Nonce() == tx.Nonce() {
				return fmt.Errorf("the same nonce %d transaction already exists, so if you want to replace the nonce transaction, add a fee", tx.Nonce())
			} else {
				return fmt.Errorf("the nonce value %d is repeated, increase the nonce value", tx.Nonce())
			}
		}
		t.ready.Put(tx)
	} else if nonce >= tx.Nonce() {
		return fmt.Errorf("the nonce value %d is repeated, increase the nonce value", tx.Nonce())
	} else {
		return t.cache.Put(tx)
	}
	return nil
}

func (t *TxManagement) DeleteAndUpdate(transactions types.ITransactions) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for _, tx := range transactions.Txs() {
		t.Remove(tx)
	}
	t.update()
}

func (t *TxManagement) DeleteEnd(newTx types.ITransaction) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.ready.PopMin(newTx.Fee())
}

func (t *TxManagement) NeedPackaged(count int) types.ITransactions {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.ready.NeedPackaged(count)
}

func (t *TxManagement) GetAll() (types.ITransactions, types.ITransactions) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	readyTxs := t.ready.All()
	cacheTxs := t.cache.All()
	return readyTxs, cacheTxs
}

func (t *TxManagement) Exist(tx types.ITransaction) bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if !t.ready.Exist(tx.From().String(), tx.Hash().String()) {
		return t.cache.Exist(tx.From().String())
	}
	return true
}

func (t *TxManagement) Update() {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	t.update()
}

func (t *TxManagement) update() {
	t.ready.RemoveExecuted(t.validator)
	for _, tx := range t.cache.txs {
		nonce := t.actStatus.Nonce(tx.From())
		if nonce < tx.Nonce()-1 {
			continue
		}
		if nonce == tx.Nonce()-1 {
			t.ready.Put(tx)
		}
		t.cache.Remove(tx)
	}
}

func (t *TxManagement) DeleteExpired(timeThreshold int64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.ready.RemoveExpiredTx(timeThreshold)

	for _, tx := range t.cache.txs {
		if tx.Time() <= timeThreshold {
			t.cache.Remove(tx)
		}
	}
}

func (t *TxManagement) Remove(tx types.ITransaction) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.cache.Remove(tx)
	t.ready.Remove(tx)
}
