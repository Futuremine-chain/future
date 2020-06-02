package txlist

import (
	"container/heap"
	"github.com/Futuremine-chain/futuremine/common/validator"
	"github.com/Futuremine-chain/futuremine/types"
)

type Sorted struct {
	txs   map[string]types.ITransaction
	cache map[string]types.ITransaction
	index *txInfos
	db    ITxListDB
}

func NewSorted(db ITxListDB) *Sorted {
	return &Sorted{
		txs:   make(map[string]types.ITransaction),
		cache: make(map[string]types.ITransaction),
		index: new(txInfos),
		db:    db,
	}
}

func (t *Sorted) Put(tx types.ITransaction) {
	t.txs[tx.From().String()] = tx
	t.cache[tx.From().String()] = tx
	heap.Push(t.index, &txInfo{
		address: tx.From().String(),
		txHash:  tx.Hash().String(),
		fees:    tx.Fee(),
		nonce:   tx.Nonce(),
		time:    tx.Time(),
	})
	t.db.Save(tx)
}

func (t *Sorted) All() types.ITransactions {
	var all types.ITransactions
	for _, tx := range t.cache {
		all.Add(tx)
	}
	return all
}

func (t *Sorted) NeedPackaged(count int) types.ITransactions {
	var txs types.ITransactions
	rIndex := t.index.CopySelf()

	for rIndex.Len() > 0 && count > 0 {
		ti := heap.Pop(rIndex).(*txInfo)
		tx := t.txs[ti.address]
		txs.Add(tx)
		count--
	}
	return txs
}

func (t *Sorted) GetByAddress(addr string) types.ITransaction {
	return t.txs[addr]
}

// If the transaction pool is full, delete the transaction with a small fee
func (t *Sorted) PopMin(fees uint64) types.ITransaction {
	if t.Len() > 0 {
		if (*t.index)[0].fees <= fees {
			ti := heap.Remove(t.index, 0).(*txInfo)
			tx := t.txs[ti.address]
			delete(t.txs, ti.address)
			delete(t.cache, ti.address)
			t.db.Delete(tx)
			return tx
		}
	}
	return nil
}

func (t *Sorted) Len() int { return len(t.txs) }

func (t *Sorted) Exist(from string, txHash string) bool {
	tx, ok := t.cache[from]
	if ok {
		return tx.Hash().String() == txHash
	}
	return false
}

func (t *Sorted) Remove(tx types.ITransaction) {
	for i, ti := range *(t.index) {
		if ti.txHash == tx.Hash().String() {
			heap.Remove(t.index, i)
			delete(t.txs, tx.From().String())
			delete(t.cache, tx.From().String())
			t.db.Delete(tx)
			return
		}
	}
}

// Delete already packed transactions
func (t *Sorted) RemoveExecuted(v validator.IValidator) {
	for _, tx := range t.cache {
		if err := v.Check(tx); err != nil {
			t.Remove(tx)
		}
	}
}

// Delete expired transactions
func (t *Sorted) RemoveExpiredTx(timeThreshold int64) {
	for _, tx := range t.cache {
		if tx.Time() <= timeThreshold {
			t.Remove(tx)
		}
	}
}

type txInfos []*txInfo

type txInfo struct {
	address string
	txHash  string
	fees    uint64
	nonce   uint64
	time    int64
}

func (t txInfos) Len() int           { return len(t) }
func (t txInfos) Less(i, j int) bool { return t[i].fees > t[j].fees }
func (t txInfos) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

func (t *txInfos) Push(x interface{}) {
	*t = append(*t, x.(*txInfo))
}

func (t *txInfos) Pop() interface{} {
	old := *t
	n := len(old)
	x := old[n-1]
	*t = old[0 : n-1]
	return x
}

func (t *txInfos) CopySelf() *txInfos {
	reReelList := new(txInfos)
	for _, nonce := range *t {
		*reReelList = append(*reReelList, nonce)
	}
	return reReelList
}

func (t *txInfos) FindIndex(addr string, nonce uint64) int {
	for index, ti := range *t {
		if ti.address == addr && ti.nonce <= nonce {
			return index
		}
	}
	return -1
}
