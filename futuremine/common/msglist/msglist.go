package msglist

import (
	"fmt"
	"github.com/Futuremine-chain/futuremine/common/account"
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/common/db/msglist"
	"github.com/Futuremine-chain/futuremine/common/validator"
	"github.com/Futuremine-chain/futuremine/types"
	"sync"
)

const msgList_db = "msg_List_db"
const maxPoolTx = 100000

type MsgManagement struct {
	cache     *Cache
	ready     *Sorted
	validator validator.IValidator
	actStatus account.IActStatus
	mutex     sync.RWMutex
	msgDB     ITxListDB
}

func NewMsgManagement(validator validator.IValidator, actStatus account.IActStatus) (*MsgManagement, error) {
	msgDB, err := msglist.Open(config.Param.Data + "/" + msgList_db)
	if err != nil {
		return nil, err
	}
	return &MsgManagement{
		cache:     NewCache(msgDB),
		ready:     NewSorted(msgDB),
		validator: validator,
		actStatus: actStatus,
		msgDB:     msgDB,
	}, nil
}

func (t *MsgManagement) Read() error {
	msgs := t.msgDB.Read()
	if msgs != nil {
		for _, msg := range msgs {
			t.Put(msg)
		}
	}
	return nil
}

func (t *MsgManagement) Close() error {
	t.msgDB.Close()
	return nil
}

func (t *MsgManagement) Count() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.cache.Len() + t.ready.Len()
}

func (t *MsgManagement) Put(msg types.IMessage) error {
	if t.Exist(msg) {
		return fmt.Errorf("the message %s already exists", msg.Hash().String())
	}
	if err := t.validator.CheckMsg(msg, false); err != nil {
		return err
	}

	if t.cache.Len() >= maxPoolTx {
		t.DeleteEnd(msg)
	}
	return t.put(msg)
}

func (t *MsgManagement) put(msg types.IMessage) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	from := msg.From().String()
	nonce := t.actStatus.Nonce(msg.From())
	if nonce == msg.Nonce()-1 {
		oldTx := t.ready.GetByAddress(from)
		if oldTx != nil {
			if oldTx.Nonce() == msg.Nonce() && oldTx.Fee() < msg.Fee() {
				t.ready.Remove(oldTx)
			} else if oldTx.Nonce() < msg.Nonce() {
				t.ready.Remove(oldTx)
			} else if oldTx.Nonce() == msg.Nonce() {
				return fmt.Errorf("the same nonce %d message already exists, so if you want to replace the nonce message, add a fee", msg.Nonce())
			} else {
				return fmt.Errorf("the nonce value %d is repeated, increase the nonce value", msg.Nonce())
			}
		}
		t.ready.Put(msg)
	} else if nonce >= msg.Nonce() {
		return fmt.Errorf("the nonce value %d is repeated, increase the nonce value", msg.Nonce())
	} else {
		return t.cache.Put(msg)
	}
	return nil
}

func (t *MsgManagement) DeleteAndUpdate(messages []types.IMessage) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for _, msg := range messages {
		t.Remove(msg)
	}
	t.update()
}

func (t *MsgManagement) DeleteEnd(newTx types.IMessage) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.ready.PopMin(newTx.Fee())
}

func (t *MsgManagement) NeedPackaged(count int) []types.IMessage {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.ready.NeedPackaged(count)
}

func (t *MsgManagement) GetAll() ([]types.IMessage, []types.IMessage) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	readyTxs := t.ready.All()
	cacheTxs := t.cache.All()
	return readyTxs, cacheTxs
}

func (t *MsgManagement) Exist(msg types.IMessage) bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if !t.ready.Exist(msg.From().String(), msg.Hash().String()) {
		return t.cache.Exist(msg.From().String())
	}
	return true
}

func (t *MsgManagement) Update() {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	t.update()
}

func (t *MsgManagement) update() {
	t.ready.RemoveExecuted(t.validator)
	for _, msg := range t.cache.msgs {
		nonce := t.actStatus.Nonce(msg.From())
		if nonce < msg.Nonce()-1 {
			continue
		}
		if nonce == msg.Nonce()-1 {
			t.ready.Put(msg)
		}
		t.cache.Remove(msg)
	}
}

func (t *MsgManagement) DeleteExpired(timeThreshold int64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.ready.RemoveExpiredTx(timeThreshold)

	for _, msg := range t.cache.msgs {
		if msg.Time() <= timeThreshold {
			t.cache.Remove(msg)
		}
	}
}

func (t *MsgManagement) Remove(msg types.IMessage) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.cache.Remove(msg)
	t.ready.Remove(msg)
}
