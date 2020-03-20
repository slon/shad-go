// +build !change

package testequal

// T is an interface wrapper around *testing.T.
type T interface {
	Errorf(format string, args ...interface{})
	Helper()
	FailNow()
}
