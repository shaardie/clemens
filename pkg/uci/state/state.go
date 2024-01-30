package state

import "sync/atomic"

type stateTypes int64

const (
	IDLE stateTypes = iota
	POSITION_SET
	RUNNING
)

type State interface {
	Get() stateTypes
	Set(v stateTypes)
	private()
}

type stateImpl atomic.Int64

func New() *stateImpl {
	return new(stateImpl)
}

func (s *stateImpl) Get() stateTypes {
	return stateTypes((*atomic.Int64)(s).Load())
}

func (s *stateImpl) Set(v stateTypes) {
	(*atomic.Int64)(s).Store(int64(v))
}

func (s *stateImpl) private() {}
