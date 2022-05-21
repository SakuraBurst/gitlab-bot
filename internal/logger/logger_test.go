package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SakuraBurst/gitlab-bot/api/clients"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	"github.com/SakuraBurst/gitlab-bot/pkg/services/telegram"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
	Init(log.TraceLevel, buffer)
	log.Info(log.InfoLevel.String())
	log.Warn(log.WarnLevel.String())
	log.Debug(log.DebugLevel.String())
	log.Trace(log.TraceLevel.String())
	log.Error(log.ErrorLevel.String())
	buffer.Truncate(buffer.Len() - 2)
	buffer.WriteString("\n]")
	assert.NotEmpty(t, buffer.Bytes(), "Writer не должен быть пустым")
	assert.GreaterOrEqual(t, 1344, buffer.Len(), "Writer должен быть больше или равен установленной длинны")
	logs := make([]LoggerInfoMessage, 6)
	err := json.Unmarshal(buffer.Bytes(), &logs)
	require.Nil(t, err)
	require.Len(t, logs, 6)
	levels := []string{log.InfoLevel.String(), log.WarnLevel.String(), log.DebugLevel.String(), log.TraceLevel.String(), log.ErrorLevel.String()}
	for i := 1; i < 6; i++ {
		assert.Equal(t, levels[i-1], logs[i].Message)
		assert.Equal(t, levels[i-1], logs[i].Level)
	}
}

// 150966050
func TestFatalNotifierHook_Fire_ErrorRequest(t *testing.T) {
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

	assert.PanicsWithError(t, err.Error(), func() {
		log.Fatal("govno")
	})

	clients.DisableMock()
}

func TestFatalNotifierHook_Fire_TelegramError(t *testing.T) {
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
	assert.PanicsWithError(t, telegramUnauthorizedMock.Error(), func() {
		log.Fatal("govno")
	})

	clients.DisableMock()
}

func TestFatalNotifierHook_Fire_OpenFileError(t *testing.T) {
	clients.EnableMock()
	telegramResponseMock := models.TelegramResponse{
		Ok:     true,
		Result: models.TelegramResult{},
	}
	tgMockBytes, err := json.Marshal(telegramResponseMock)
	require.Nil(t, err)
	reader := bytes.NewReader(tgMockBytes)
	readCloser := io.NopCloser(reader)

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
	//absPath, err := filepath.Abs(".")
	//require.Nil(t, err)
	//cantOpenFileErrorStringMock := "open " + filepath.FromSlash(absPath+"/logger.json: The system cannot find the file specified.")

	assert.Panics(t, func() {
		log.Fatal("govno")
	})

	clients.DisableMock()
}

func TestFatalNotifierHook_Fire_SendFileError(t *testing.T) {
	clients.EnableMock()
	reader := bytes.NewReader(nil)
	clients.Mocks.AddMock("https://api.telegram.org/bot/sendMessage", clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(reader),
		},
		Err: nil,
	})
	errorMock := errors.New("test")
	clients.Mocks.AddMock("https://api.telegram.org/bot/sendDocument", clients.Mock{
		Response: nil,
		Err:      errorMock,
	})
	fr := &FatalNotifier{
		Bot:     telegram.Bot{},
		LogFile: nil,
	}
	AddHook(fr)
	absPath, err := filepath.Abs("logger.json")
	require.Nil(t, err)
	file, err := os.Create(absPath)
	require.Nil(t, err)
	assert.PanicsWithError(t, errorMock.Error(), func() {
		log.Fatal("govno")
	})
	err = file.Close()
	require.Nil(t, err)
	err = os.Remove(absPath)
	require.Nil(t, err)
	clients.DisableMock()
}

func TestFatalNotifierHook_Fire(t *testing.T) {
	absPath, err := filepath.Abs("logger.json")
	if os.Getenv("FATAL") == "1" {
		clients.EnableMock()
		reader := bytes.NewReader(nil)
		okResponse := &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(reader),
		}
		clients.Mocks.AddMock("https://api.telegram.org/bot/sendMessage", clients.Mock{
			Response: okResponse,
			Err:      nil,
		})
		clients.Mocks.AddMock("https://api.telegram.org/bot/sendDocument", clients.Mock{
			Response: okResponse,
			Err:      nil,
		})
		fr := &FatalNotifier{
			Bot:     telegram.Bot{},
			LogFile: nil,
		}
		require.Nil(t, err)
		_, err := os.Create(absPath)
		require.Nil(t, err)
		AddHook(fr)
		log.Fatal("govno")
		clients.DisableMock()
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalNotifierHook")
	cmd.Env = append(os.Environ(), "FATAL=1")
	err = cmd.Run()
	e, ok := err.(*exec.ExitError)
	require.True(t, ok)
	assert.Equal(t, 1, e.ExitCode(), "Процесс должен завершиться с кодом 1")
	err = os.Remove(absPath)
	require.Nil(t, err)
}

// goland generated test
func TestGetLogLevel(t *testing.T) {
	type args struct {
		isProduction bool
	}
	tests := []struct {
		name string
		args args
		want log.Level
	}{
		{"isProductionFalse", args{false}, log.InfoLevel},
		{"isProductionTrue", args{true}, log.ErrorLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetLogLevel(tt.args.isProduction), "GetLogLevel(%v)", tt.args.isProduction)
		})
	}
}
