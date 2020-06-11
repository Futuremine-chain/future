package types

type IMessage interface {
	IMessageHeader
	MsgBody() IMessageBody
	ToRlp() IRlpMessage
	Check() error
}
