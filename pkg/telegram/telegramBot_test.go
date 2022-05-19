package telegram

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBot(t *testing.T) {
	testString := "test"
	bot := NewBot(testString, testString)
	assert.Equal(t, bot.token, testString)
	assert.Equal(t, bot.mainChannel, testString)
}
