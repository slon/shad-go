// +build !change

package testequal

// T is an interface wrapper for *testing.T
// that contains only a small subset of methods.
type T interface {
	Errorf(format string, args ...interface{})
	Helper()
	FailNow()
}
