package patch

type Patch interface {
	ToKey() string
	GetLineNumber() int
	GetData() string
	GetFilePath() string
	GetOperand() string
	AppendData(data string)
}

const (
	TypeOperandReplace = "<<<"
	TypeOperandAppend  = ">>>"
)
