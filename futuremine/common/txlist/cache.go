package txlist

import (
	"fmt"
	"github.com/Futuremine-chain/futuremine/types"
	"strconv"
)

type Cache struct {
	txs      map[string]types.ITransaction
	nonceTxs map[string]string
	db       ITxListDB
}

func NewCache(db ITxListDB) *Cache {
	return &Cache{
		txs:      make(map[string]types.ITransaction),
		nonceTxs: make(map[string]string),
		db:       db,
	}
}

func (c *Cache) Put(tx types.ITransaction) error {
	if c.Exist(tx.Hash().String()) {
		return fmt.Errorf("transation hash %s exsit", tx.Hash())
	}
	nonceKey := nonceKey(tx)
	if oldTxHash := c.getHash(nonceKey); oldTxHash != "" {
		oldTx := c.txs[oldTxHash]
		if oldTx.Fee() > tx.Fee() {
			return fmt.Errorf("transation nonce %d exist, the fees must biger than before %d", tx.Nonce(), oldTx.Fee())
		}
		c.Remove(oldTx)
	}
	c.txs[tx.Hash().String()] = tx
	c.nonceTxs[nonceKey] = tx.Hash().String()
	c.db.Save(tx)
	return nil
}

func (c *Cache) Remove(tx types.ITransaction) {
	delete(c.txs, tx.Hash().String())
	delete(c.nonceTxs, nonceKey(tx))
	c.db.Delete(tx)
}

func (c *Cache) Exist(txHash string) bool {
	_, ok := c.txs[txHash]
	return ok
}

func (c *Cache) Len() int {
	return len(c.txs)
}

func (c *Cache) All() types.ITransactions {
	var all types.ITransactions
	for _, tx := range c.txs {
		all.Add(tx)
	}
	return all
}

func (c *Cache) getHash(nonceKey string) string {
	return c.nonceTxs[nonceKey]
}

func nonceKey(tx types.ITransaction) string {
	return tx.From().String() + "_" + strconv.FormatUint(tx.Nonce(), 10)
}
