package workers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConstants(t *testing.T) {
	assert.Equal(t, 24, fullDay)
	assert.Equal(t, 100, maxPossibleErrors)
}

func TestErrorCheckerNoErrors(t *testing.T) {
	stopChan := make(chan error, 1)

	go errorChecker(nil, stopChan)
	time.Sleep(time.Millisecond * 50)
	assert.Len(t, stopChan, 0)
	assert.Equal(t, 0, errorCounter)
	errorCounter = 0
}

func TestErrorChecker99ErrorsAndOneNil(t *testing.T) {
	stopChan := make(chan error, 1)
	errorCounter = 99
	go errorChecker(nil, stopChan)
	time.Sleep(time.Millisecond * 50)
	assert.Len(t, stopChan, 0)
	assert.Equal(t, 0, errorCounter)
}

func TestErrorCheckerFinalError(t *testing.T) {
	stopChan := make(chan error, 1)
	errorCounter = 99
	go errorChecker(errors.New("test"), stopChan)
	time.Sleep(time.Millisecond * 50)
	assert.Len(t, stopChan, 1)
	assert.NotNil(t, <-stopChan)
	assert.Equal(t, 100, errorCounter)
}
