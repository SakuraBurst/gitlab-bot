package workers

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOneMinuteWorker(t *testing.T) {
	counter := 0
	var TestWorkerFuncError = func() error {
		counter = 1
		return nil
	}
	go OneMinuteWorker(nil, TestWorkerFuncError)
	time.Sleep(time.Millisecond * 50)
	assert.Equal(t, counter, 1)
}
