//go:build !solution
// +build !solution

package fileleak

type testingT interface {
	Errorf(msg string, args ...interface{})
	Cleanup(func())
}

func VerifyNone(t testingT) {
	panic("implement me")
}
