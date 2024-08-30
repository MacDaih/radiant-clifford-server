package domain

type ErrNotFound struct {
	Msg string
}

func (nf ErrNotFound) Error() string {
	return nf.Msg
}
