package txlist

import (
	"github.com/Futuremine-chain/futuremine/common/db/base"
	"github.com/Futuremine-chain/futuremine/types"
)

const (
	path   = "txlist"
	bucket = "txlist"
)

type TxListDB struct {
	base *base.Base
}

func NewTxListDB(dir string) (*TxListDB, error) {
	base, err := base.Open(dir + "/" + path)
	if err != nil {
		return nil, err
	}
	return &TxListDB{
		base: base,
	}, nil
}

func (t *TxListDB) Read() []types.ITransaction {
	return nil
}

func (t *TxListDB) Save(tx types.ITransaction) {
	key := base.Key(bucket, tx.Hash().Bytes())
	t.base.Update(key, tx.ToRlp().Bytes())
}

func (t *TxListDB) Delete(tx types.ITransaction) {
	key := base.Key(bucket, tx.Hash().Bytes())
	t.base.Delete(key)
}

func (t *TxListDB) Clear() {
	t.base.Clear(bucket)
}

func (t *TxListDB) Close() error {
	return t.base.Close()
}
