package msglist

import (
	"github.com/Futuremine-chain/futuremine/common/db/base"
	"github.com/Futuremine-chain/futuremine/types"
)

const (
	path   = "msglist"
	bucket = "msglist"
)

type MsgListDB struct {
	base *base.Base
}

func Open(path string) (*MsgListDB, error) {
	var err error
	baseDB, err := base.Open(path)
	if err != nil {
		return nil, err
	}
	return &MsgListDB{base: baseDB}, nil
}

func (t *MsgListDB) Read() []types.IMessage {
	return nil
}

func (t *MsgListDB) Save(msg types.IMessage) {
	key := base.Key(bucket, msg.Hash().Bytes())
	t.base.Put(key, msg.ToRlp().Bytes())
}

func (t *MsgListDB) Delete(msg types.IMessage) {
	key := base.Key(bucket, msg.Hash().Bytes())
	t.base.Delete(key)
}

func (t *MsgListDB) Clear() {
	t.base.Clear(bucket)
}

func (t *MsgListDB) Close() error {
	return t.base.Close()
}
