package helpers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestFileStatMust_Panic(t *testing.T) {
	assert.Panics(t, func() {
		FileStatMust(nil)
	})
}

func TestFileStatMust(t *testing.T) {
	file, err := os.Create("test.txt")
	require.Nil(t, err)
	fileStats := FileStatMust(file)
	assert.NotNil(t, fileStats)
	assert.Equal(t, int64(0), fileStats.Size())
	assert.Equal(t, "test.txt", fileStats.Name())
	err = file.Close()
	require.Nil(t, err)
	err = os.Remove("test.txt")
	require.Nil(t, err)
}
