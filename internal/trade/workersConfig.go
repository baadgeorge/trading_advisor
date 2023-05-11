package trade

import (
	"final/internal/strategy"
	"github.com/google/uuid"
)

// worker configuration struct which send by channel for Run
type WorkerConfig struct {
	workerId uint32
	figi     string
	strategy strategy.Strategy
}

// TODO
func NewWorkerConfig(figi string, strategy strategy.Strategy) *WorkerConfig {
	return &WorkerConfig{
		workerId: uuid.New().ID(),
		figi:     figi,
		strategy: strategy,
	}
}

func (wc *WorkerConfig) GetWorkerId() uint32 {
	return wc.workerId
}
