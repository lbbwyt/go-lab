package engine

import "sync"

type DefinitionStore interface {
	Get(id string) (*ProcessDef, bool)
	Save(def *ProcessDef)
}

type InstanceStore interface {
	SaveInst(inst *ProcessInst) error
	GetInst(id string) (*ProcessInst, bool)
	SaveToken(tok *Token) error
	GetToken(id string) (*Token, bool)
	DeleteToken(id string) error
	Tokens() []*Token
}

// 内存实现
type MemStore struct {
	defs   map[string]*ProcessDef
	insts  map[string]*ProcessInst
	tokens map[string]*Token

	defLock   sync.RWMutex
	instLock  sync.RWMutex
	tokenLock sync.RWMutex
}

func NewMemStore() *MemStore {
	return &MemStore{
		defs:   make(map[string]*ProcessDef),
		insts:  make(map[string]*ProcessInst),
		tokens: make(map[string]*Token),
	}
}

func (m *MemStore) Get(id string) (*ProcessDef, bool) {
	m.defLock.RLock()
	defer m.defLock.RUnlock()
	d, ok := m.defs[id]
	return d, ok
}
func (m *MemStore) Save(def *ProcessDef) {
	m.defLock.Lock()
	defer m.defLock.Unlock()
	m.defs[def.ID] = def
}

func (m *MemStore) SaveInst(inst *ProcessInst) error {
	m.instLock.Lock()
	defer m.instLock.Unlock()
	m.insts[inst.ID] = inst
	return nil
}
func (m *MemStore) GetInst(id string) (*ProcessInst, bool) {
	m.instLock.RLock()
	defer m.instLock.RUnlock()
	i, ok := m.insts[id]
	return i, ok
}

func (m *MemStore) SaveToken(tok *Token) error {
	m.tokenLock.Lock()
	defer m.tokenLock.Unlock()
	m.tokens[tok.ID] = tok
	return nil
}
func (m *MemStore) GetToken(id string) (*Token, bool) {
	m.tokenLock.RLock()
	defer m.tokenLock.RUnlock()
	t, ok := m.tokens[id]
	return t, ok
}
func (m *MemStore) DeleteToken(id string) error {
	m.tokenLock.Lock()
	defer m.tokenLock.Unlock()
	delete(m.tokens, id)
	return nil
}

// 返回当前所有 token 的快照（只读）
func (m *MemStore) Tokens() []*Token {
	m.tokenLock.RLock()
	defer m.tokenLock.RUnlock()

	list := make([]*Token, 0, len(m.tokens))
	for _, t := range m.tokens {
		// 复制一份避免外部并发修改
		cpy := *t
		list = append(list, &cpy)
	}
	return list
}
