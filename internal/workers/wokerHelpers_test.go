package workers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConstants(t *testing.T) {
	assert.Equal(t, fullDay, 24)
	assert.Equal(t, maxPossibleErrors, 100)
}

func TestErrorCheckerNoErrors(t *testing.T) {
	stopChan := make(chan error, 1)

	go errorChecker(nil, stopChan)
	time.Sleep(time.Millisecond * 50)
	assert.Len(t, stopChan, 0)
	assert.Equal(t, errorCounter, 0)
	errorCounter = 0
}

func TestErrorChecker99ErrorsAndOneNil(t *testing.T) {
	stopChan := make(chan error, 1)
	errorCounter = 99
	go errorChecker(nil, stopChan)
	time.Sleep(time.Millisecond * 50)
	assert.Len(t, stopChan, 0)
	assert.Equal(t, errorCounter, 0)
}

func TestErrorCheckerFinalError(t *testing.T) {
	stopChan := make(chan error, 1)
	errorCounter = 99
	go errorChecker(errors.New("test"), stopChan)
	time.Sleep(time.Millisecond * 50)
	assert.Len(t, stopChan, 1)
	assert.NotNil(t, <-stopChan)
	assert.Equal(t, errorCounter, 100)
}
