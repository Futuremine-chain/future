package types

type MessageType int

type Message_Header struct {
	Type MessageType
}

func (t *Message_Header) Check() error {
	if err := t.CheckType(); err != nil {
		return err
	}

	/*	if err := t.verifyTxHash(); err != nil {
			return err
		}

		if err := t.verifyTxFrom(); err != nil {
			return err
		}

		if err := t.verifyTxFees(); err != nil {
			return err
		}

		if err := t.verifyTxSinger(); err != nil {
			return err
		}*/
	return nil
}

func (t *Message_Header) CheckType() error {
	/*switch t.TxType {
	case NormalMessage:
		return nil
	case ContractMessage:
		return nil
	case VoteToCandidate:
		return nil
	case LoginCandidate:
		return nil
	case LogoutCandidate:
		return nil
	}
	return ErrTxType*/
	return nil
}
