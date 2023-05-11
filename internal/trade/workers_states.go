package trade

type State int

const (
	Info_type State = iota
	Err_type
	Signal_type
	Cancel_type
)
