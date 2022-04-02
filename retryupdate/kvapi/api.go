//go:build !change

package kvapi

import "github.com/gofrs/uuid"

type (
	Client interface {
		// Get the value of key.
		Get(req *GetRequest) (*GetResponse, error)

		// Set key to hold given value.
		Set(rsp *SetRequest) (*SetResponse, error)
	}

	GetRequest struct {
		Key string
	}

	GetResponse struct {
		Value   string
		Version uuid.UUID
	}

	SetRequest struct {
		Key, Value string

		// OldVersion field is checked before updating value.
		//
		// When updating old value, OldVersion must specify uuid of currently stored value.
		// When creating new value, OldVersion must hold zero value of UUID.
		OldVersion, NewVersion uuid.UUID
	}

	SetResponse struct{}
)
