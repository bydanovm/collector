package redis

import (
	"math"
	"sync"
)

type Process struct {
	id      uint64
	protect sync.Mutex
}

// Получение ID процесса использующего транзакцию
func (p *Process) GetID() (id uint64) {
	p.protect.Lock()
	defer p.protect.Unlock()
	id = 1
	if p.id > 0 && p.id < math.MaxUint32 {
		id = p.id
	} else {
		p.id = id
	}
	p.id++
	return id
}
