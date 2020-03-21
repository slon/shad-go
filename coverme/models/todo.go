// +build !change

package models

type ID int

type AddRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Todo struct {
	ID       ID     `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Finished bool   `json:"finished"`
}

func (t *Todo) MarkFinished() {
	t.Finished = true
}

func (t *Todo) MarkUnfinished() {
	t.Finished = false
}
