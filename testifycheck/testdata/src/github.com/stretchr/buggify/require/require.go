package require

type TestingT interface {
	Errorf(format string, args ...interface{})
}

func NotNil(t TestingT, object interface{}, msgAndArgs ...interface{}) {
	panic("not implemented")
}

func NotNilf(t TestingT, object interface{}, msg string, args ...interface{}) {
	panic("not implemented")
}

func Nil(t TestingT, object interface{}, msgAndArgs ...interface{}) {
	panic("not implemented")
}

func Nilf(t TestingT, object interface{}, msg string, args ...interface{}) {
	panic("not implemented")
}

func NoError(t TestingT, object interface{}, msgAndArgs ...interface{}) {
	panic("not implemented")
}

func NoErrorf(t TestingT, err error, msg string, args ...interface{}) {
	panic("not implemented")
}

func Error(t TestingT, object interface{}, msgAndArgs ...interface{}) {
	panic("not implemented")
}

func Errorf(t TestingT, err error, msg string, args ...interface{}) {
	panic("not implemented")
}

type Assertions struct{}

func New(t TestingT) *Assertions {
	return nil
}

func (*Assertions) NotNil(object interface{}, msgAndArgs ...interface{}) {
	panic("not implemented")
}

func (*Assertions) NotNilf(object interface{}, msg string, args ...interface{}) {
	panic("not implemented")
}

func (*Assertions) Nil(object interface{}, msgAndArgs ...interface{}) {
	panic("not implemented")
}

func (*Assertions) Nilf(object interface{}, msg string, args ...interface{}) {
	panic("not implemented")
}

func (*Assertions) NoError(object interface{}, msgAndArgs ...interface{}) {
	panic("not implemented")
}

func (*Assertions) NoErrorf(object interface{}, msg string, args ...interface{}) {
	panic("not implemented")
}

func (*Assertions) Error(object interface{}, msgAndArgs ...interface{}) {
	panic("not implemented")
}

func (*Assertions) Errorf(object interface{}, msg string, args ...interface{}) {
	panic("not implemented")
}
