package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

type HeartbeatClient struct {
	l        *zap.Logger
	endpoint string
}

func NewHeartbeatClient(l *zap.Logger, endpoint string) *HeartbeatClient {
	return &HeartbeatClient{l: l, endpoint: endpoint}
}

func (c *HeartbeatClient) Heartbeat(ctx context.Context, req *HeartbeatRequest) (*HeartbeatResponse, error) {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+"/heartbeat", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}

	httpRsp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode != http.StatusOK {
		errorMsg, _ := ioutil.ReadAll(httpRsp.Body)
		return nil, fmt.Errorf("heartbeat failed: %s", errorMsg)
	}

	rspJSON, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}

	rsp := &HeartbeatResponse{}
	if err = json.Unmarshal(rspJSON, rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}
