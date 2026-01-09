package store

import (
	"encoding/json"
	"errors"
	"go-lab/app/go_dag/model"
)

type Store interface {
	SaveWorkflow(wf *model.Workflow) error
	LoadWorkflow(id string) (*model.Workflow, error)
}

type MemStore struct {
	ss map[string]string
}

func NewMemStore() *MemStore {
	return &MemStore{
		ss: make(map[string]string, 0),
	}
}

func (ms *MemStore) SaveWorkflow(wf *model.Workflow) error {
	bs, err := json.Marshal(wf)
	if err != nil {
		return err
	}
	ms.ss[wf.ID] = string(bs)
	return nil
}

func (ms *MemStore) LoadWorkflow(id string) (*model.Workflow, error) {

	if v, ok := ms.ss[id]; ok {

		var m = &model.Workflow{}
		err := json.Unmarshal([]byte(v), m)
		if err != nil {
			return nil, errors.New("unmarshal error")
		}
		return m, nil

	}

	return nil, errors.New("the workflow does not exist")
}
