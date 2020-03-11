// +build !solution

package retryupdate

import "gitlab.com/slon/shad-go/retryupdate/kvapi"

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	panic("implement me")
}
