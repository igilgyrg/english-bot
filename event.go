package english

const (
	Unknown Type = iota
	Message
)

type Type int

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
