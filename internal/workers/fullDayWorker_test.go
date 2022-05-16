package workers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFullDayWorkerOneHourUp(t *testing.T) {
	counter := 0
	var TestWorkerFuncError = func() error {
		counter = 1
		return nil
	}
	currentTime := time.Now().Add(time.Hour)
	go FullDayWorker(nil, TestWorkerFuncError, currentTime.Hour())
	time.Sleep(time.Millisecond * 50)
	assert.Equal(t, counter, 0)
}

func TestFullDayWorkerOneHourDown(t *testing.T) {
	counter := 0
	var TestWorkerFuncError = func() error {
		counter = 1
		return nil
	}
	currentTime := time.Now().Add(-time.Hour)
	fmt.Println(currentTime.Hour())
	go FullDayWorker(nil, TestWorkerFuncError, currentTime.Hour())
	time.Sleep(time.Millisecond * 50)
	assert.Equal(t, counter, 0)
}

func TestFullDayWorkerWithWakeUp(t *testing.T) {
	counter := 0
	var TestWorkerFuncError = func() error {
		counter = 1
		return nil
	}
	currentTime := time.Now()
	go FullDayWorker(nil, TestWorkerFuncError, currentTime.Hour())
	time.Sleep(time.Millisecond * 50)
	assert.Equal(t, counter, 1)
}

func TestFullDayWorkerWeekend(t *testing.T) {
	counter := 0
	var TestWorkerFuncError = func() error {
		counter = 1
		return nil
	}
	currentTime := time.Now()
	go FullDayWorker(nil, TestWorkerFuncError, currentTime.Hour())
	time.Sleep(time.Millisecond * 50)
	assert.Equal(t, counter, 1)
}
