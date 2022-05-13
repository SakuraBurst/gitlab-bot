package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SakuraBurst/gitlab-bot/api/clients"
	"github.com/SakuraBurst/gitlab-bot/pkg/telegram"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

type LoggerInfoMessage struct {
	File     string    `json:"file"`
	Function string    `json:"func"`
	Level    string    `json:"level"`
	Message  string    `json:"msg"`
	Time     time.Time `json:"time"`
}

func TestInit(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString("[\n")
	currentTime := time.Now().Truncate(time.Second)
	Init(log.TraceLevel, buffer)
	log.Info(log.InfoLevel.String())
	log.Warn(log.WarnLevel.String())
	log.Debug(log.DebugLevel.String())
	log.Trace(log.TraceLevel.String())
	log.Error(log.ErrorLevel.String())
	buffer.Truncate(buffer.Len() - 2)
	buffer.WriteString("\n]")
	assert.NotEmpty(t, buffer.Bytes(), "Writer не должен быть пустым")
	assert.Equal(t, buffer.Len(), 1344, "Writer должен быть больше или равен установленной длинны")
	logs := make([]LoggerInfoMessage, 6)
	err := json.Unmarshal(buffer.Bytes(), &logs)
	require.Nil(t, err)
	require.Len(t, logs, 6)
	levels := []string{log.InfoLevel.String(), log.WarnLevel.String(), log.DebugLevel.String(), log.TraceLevel.String(), log.ErrorLevel.String()}
	for i := 1; i < 6; i++ {
		assert.Equal(t, logs[i].Time, currentTime)
		assert.Equal(t, logs[i].Message, levels[i-1])
		assert.Equal(t, logs[i].Level, levels[i-1])
	}
}

// 150966050
func TestFatalReminderHookError(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	Init(log.TraceLevel, buffer)
	clients.EnableMock()
	fr := &FatalNotifier{
		Bot:     telegram.Bot{},
		LogFile: nil,
	}
	AddHook(fr)
	err := errors.New("error")
	clients.Mocks.AddMock("https://api.telegram.org/bot/sendMessage", clients.Mock{
		Response: nil,
		Err:      err,
	})

	assert.PanicsWithValue(t, err, func() {
		log.Fatal("govno")
	}, "Должна произойти паника")

	clients.Mocks.ClearMocks()

	clients.Mocks["https://api.telegram.org/bot/sendMessage"] = clients.Mock{
		Response: nil,
		Err:      nil,
	}

}
