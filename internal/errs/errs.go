package errs

type Temporary struct {
	Err string
}

func (t *Temporary) Error() string {
	return t.Err
}

type Argument struct {
	Err string
}

func (t *Argument) Error() string {
	return t.Err
}
