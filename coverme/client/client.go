// +build !change

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/slon/shad-go/coverme/models"
)

type Client struct {
	addr string
}

func New(addr string) *Client {
	return &Client{addr: addr}
}

func (c *Client) Add(r *models.AddRequest) (*models.Todo, error) {
	data, _ := json.Marshal(r)

	resp, err := http.Post(c.addr+"/todo/create", "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var todo *models.Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	return todo, err
}

func (c *Client) Get(id models.ID) (*models.Todo, error) {
	resp, err := http.Get(c.addr + fmt.Sprintf("/todo/%d", id))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var todo *models.Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	return todo, err
}

func (c *Client) List() ([]*models.Todo, error) {
	resp, err := http.Get(c.addr + "/todo")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var todos []*models.Todo
	err = json.NewDecoder(resp.Body).Decode(&todos)
	return todos, err
}
