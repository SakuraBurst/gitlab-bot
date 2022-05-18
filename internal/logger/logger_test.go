package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SakuraBurst/gitlab-bot/api/clients"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	"github.com/SakuraBurst/gitlab-bot/pkg/telegram"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"os/exec"
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
func TestFatalNotifierHookErrorRequest(t *testing.T) {
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
	})

	clients.Mocks.ClearMocks()
}

func TestFatalNotifierHookTelegramError(t *testing.T) {
	telegramUnauthorizedMock := models.TelegramError{
		Ok:          false,
		ErrorCode:   401,
		Description: "Unauthorized",
	}
	tgMockBytes, err := json.Marshal(telegramUnauthorizedMock)
	require.Nil(t, err)
	reader := bytes.NewReader(tgMockBytes)
	readCloser := io.NopCloser(reader)
	clients.EnableMock()
	clients.Mocks.AddMock("https://api.telegram.org/bot/sendMessage", clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusUnauthorized),
			StatusCode: http.StatusUnauthorized,
			Body:       readCloser,
		},
		Err: nil,
	})
	fr := &FatalNotifier{
		Bot:     telegram.Bot{},
		LogFile: nil,
	}
	AddHook(fr)
	assert.PanicsWithValue(t, telegramUnauthorizedMock, func() {
		log.Fatal("govno")
	})

	clients.Mocks.ClearMocks()
}

func TestFatalNotifierHook(t *testing.T) {
	if os.Getenv("FATAL") == "1" {
		telegramResponseMock := models.TelegramResponse{
			Ok:     true,
			Result: models.TelegramResult{},
		}
		tgMockBytes, err := json.Marshal(telegramResponseMock)
		require.Nil(t, err)
		reader := bytes.NewReader(tgMockBytes)
		readCloser := io.NopCloser(reader)
		clients.EnableMock()
		clients.Mocks.AddMock("https://api.telegram.org/bot/sendMessage", clients.Mock{
			Response: &http.Response{
				Status:     http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Body:       readCloser,
			},
			Err: nil,
		})
		fr := &FatalNotifier{
			Bot:     telegram.Bot{},
			LogFile: nil,
		}
		AddHook(fr)
		log.Fatal("govno")
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalNotifierHook")
	cmd.Env = append(os.Environ(), "FATAL=1")
	err := cmd.Run()
	e, ok := err.(*exec.ExitError)
	require.True(t, ok)
	assert.Equal(t, e.ExitCode(), 1, "Процесс должен завершиться с кодом 1")
}
