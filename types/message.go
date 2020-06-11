package types

type IMessage interface {
	IMessageHeader
	IMessageBody

	ToRlp() IRlpMessage
	Check() error
}
