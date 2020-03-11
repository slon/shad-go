package retryupdate_test

//go:generate mockgen -destination mock_test.go -package retryupdate_test gitlab.com/slon/shad-go/retryupdate/kvapi Client

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"gitlab.com/slon/shad-go/retryupdate"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

var (
	K0 = "K0"
	V0 = "V0"
	V1 = "V1"
	V2 = "V2"
	V3 = "V3"

	UUID0 = uuid.Must(uuid.NewV4())
	UUID1 = uuid.Must(uuid.NewV4())
	UUID2 = uuid.Must(uuid.NewV4())

	errUpdate = errors.New("update error")

	errGetAuth = &kvapi.APIError{Method: "get", Err: &kvapi.AuthError{Msg: "token expired"}}
	errSetAuth = &kvapi.APIError{Method: "set", Err: &kvapi.AuthError{Msg: "token expired"}}

	errGetNoKey = &kvapi.APIError{Method: "get", Err: kvapi.ErrKeyNotFound}
	errSetNoKey = &kvapi.APIError{Method: "set", Err: kvapi.ErrKeyNotFound}

	errGetTemporary = &kvapi.APIError{Method: "get", Err: errors.New("unavailable")}
	errSetTemporary = &kvapi.APIError{Method: "set", Err: errors.New("unavailable")}
)

type setMatcher struct {
	kvapi.SetRequest

	save *uuid.UUID
}

func (m setMatcher) Matches(x interface{}) bool {
	if arg, ok := x.(*kvapi.SetRequest); ok {
		if m.save != nil {
			*m.save = arg.NewVersion
		}

		return arg.Key == m.Key && arg.Value == m.Value && arg.OldVersion == m.OldVersion
	}

	return false
}

func (m setMatcher) String() string {
	return fmt.Sprintf("%v", m.SetRequest)
}

func SetRequest(k, v string, oldVersion uuid.UUID, saveUUID ...*uuid.UUID) gomock.Matcher {
	m := setMatcher{
		SetRequest: kvapi.SetRequest{
			Key:        k,
			Value:      v,
			OldVersion: oldVersion,
		},
	}

	if len(saveUUID) > 1 {
		panic("error")
	}

	if len(saveUUID) == 1 {
		m.save = saveUUID[0]
	}

	return m
}

func updateFn(oldValue *string) (string, error) {
	switch {
	case oldValue == nil:
		return V0, nil
	case *oldValue == V0:
		return V1, nil
	case *oldValue == V1:
		return V2, nil
	case *oldValue == V2:
		return V3, nil
	default:
		return "", errUpdate
	}
}

func TestSimpleUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := NewMockClient(ctrl)
	gomock.InOrder(
		c.EXPECT().
			Get(&kvapi.GetRequest{Key: K0}).
			Return(&kvapi.GetResponse{Value: V0, Version: UUID0}, nil),

		c.EXPECT().
			Set(SetRequest(K0, V1, UUID0)).
			Return(&kvapi.SetResponse{}, nil),
	)

	require.NoError(t, retryupdate.UpdateValue(c, K0, updateFn))
}

func TestUpdateFnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := NewMockClient(ctrl)
	gomock.InOrder(
		c.EXPECT().
			Get(&kvapi.GetRequest{Key: K0}).
			Return(&kvapi.GetResponse{Value: V3, Version: UUID0}, nil),
	)

	require.Equal(t, errUpdate, retryupdate.UpdateValue(c, K0, updateFn))
}
func TestCreateKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := NewMockClient(ctrl)
	gomock.InOrder(
		c.EXPECT().
			Get(&kvapi.GetRequest{Key: K0}).
			Return(nil, errGetNoKey),

		c.EXPECT().
			Set(SetRequest(K0, V0, uuid.UUID{})).
			Return(&kvapi.SetResponse{}, nil),
	)

	require.NoError(t, retryupdate.UpdateValue(c, K0, updateFn))
}

func TestKeyVanished(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := NewMockClient(ctrl)
	gomock.InOrder(
		c.EXPECT().
			Get(&kvapi.GetRequest{Key: K0}).
			Return(&kvapi.GetResponse{Value: V2, Version: UUID0}, nil),

		c.EXPECT().
			Set(SetRequest(K0, V3, UUID0)).
			Return(nil, errSetNoKey),

		c.EXPECT().
			Set(SetRequest(K0, V0, uuid.UUID{})).
			Return(&kvapi.SetResponse{}, nil),
	)

	require.NoError(t, retryupdate.UpdateValue(c, K0, updateFn))
}

func TestFailOnAuthErrorInGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := NewMockClient(ctrl)
	gomock.InOrder(
		c.EXPECT().
			Get(&kvapi.GetRequest{Key: K0}).
			Return(nil, errGetAuth),
	)

	require.Equal(t, errGetAuth, retryupdate.UpdateValue(c, K0, updateFn))
}

func TestFailOnAuthErrorInSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := NewMockClient(ctrl)
	gomock.InOrder(
		c.EXPECT().
			Get(&kvapi.GetRequest{Key: K0}).
			Return(nil, errGetNoKey),

		c.EXPECT().
			Set(SetRequest(K0, V0, uuid.UUID{})).
			Return(nil, errSetAuth),
	)

	require.Equal(t, errSetAuth, retryupdate.UpdateValue(c, K0, updateFn))
}

func TestRetryGetError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := NewMockClient(ctrl)
	gomock.InOrder(
		c.EXPECT().
			Get(&kvapi.GetRequest{Key: K0}).
			Return(nil, errGetTemporary),

		c.EXPECT().
			Get(&kvapi.GetRequest{Key: K0}).
			Return(nil, errGetTemporary),

		c.EXPECT().
			Get(&kvapi.GetRequest{Key: K0}).
			Return(&kvapi.GetResponse{Value: V0, Version: UUID0}, nil),

		c.EXPECT().
			Set(SetRequest(K0, V1, UUID0)).
			Return(&kvapi.SetResponse{}, nil),
	)

	require.NoError(t, retryupdate.UpdateValue(c, K0, updateFn))
}

func TestRetrySetError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := NewMockClient(ctrl)
	gomock.InOrder(
		c.EXPECT().
			Get(&kvapi.GetRequest{Key: K0}).
			Return(&kvapi.GetResponse{Value: V0, Version: UUID0}, nil),

		c.EXPECT().
			Set(SetRequest(K0, V1, UUID0)).
			Return(nil, errSetTemporary),

		c.EXPECT().
			Set(SetRequest(K0, V1, UUID0)).
			Return(&kvapi.SetResponse{}, nil),
	)

	require.NoError(t, retryupdate.UpdateValue(c, K0, updateFn))
}

func TestRetrySetConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := NewMockClient(ctrl)
	gomock.InOrder(
		c.EXPECT().
			Get(&kvapi.GetRequest{Key: K0}).
			Return(&kvapi.GetResponse{Value: V0, Version: UUID0}, nil),

		c.EXPECT().
			Set(SetRequest(K0, V1, UUID0)).
			Return(nil, errSetTemporary),

		c.EXPECT().
			Set(SetRequest(K0, V1, UUID0)).
			Return(nil, &kvapi.APIError{Method: "set", Err: &kvapi.ConflictError{ExpectedVersion: UUID1, ProvidedVersion: UUID0}}),

		c.EXPECT().
			Get(&kvapi.GetRequest{Key: K0}).
			Return(&kvapi.GetResponse{Value: V2, Version: UUID1}, nil),

		c.EXPECT().
			Set(SetRequest(K0, V3, UUID1)).
			Return(&kvapi.SetResponse{}, nil),
	)

	require.NoError(t, retryupdate.UpdateValue(c, K0, updateFn))
}

func TestRetrySetFalseConflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conflictErr := &kvapi.ConflictError{ProvidedVersion: UUID0}

	c := NewMockClient(ctrl)
	gomock.InOrder(
		c.EXPECT().
			Get(&kvapi.GetRequest{Key: K0}).
			Return(&kvapi.GetResponse{Value: V0, Version: UUID0}, nil),

		// first Set updates key state, but returns an error.
		c.EXPECT().
			Set(SetRequest(K0, V1, UUID0, &conflictErr.ExpectedVersion)).
			Return(nil, errSetTemporary),

		// second Set returns conflict with ExpectedVersion == OldVersion from previous request.
		c.EXPECT().
			Set(SetRequest(K0, V1, UUID0)).
			Return(nil, &kvapi.APIError{Method: "set", Err: conflictErr}),
	)

	require.NoError(t, retryupdate.UpdateValue(c, K0, updateFn))
}
