package BasaDannih

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasaDannihMySQLPostgresMongoPgAdmin777_ReadAndWriteFromBd(t *testing.T) {
	bd := BasaDannihMySQLPostgresMongoPgAdmin777{}
	bd.WriteToBD(1)
	assert.True(t, bd.ReadFromBd(1), "После записи в бд значение должно быть тру")
	assert.Equal(t, len(bd), 1, "После записи в бд размер бд должен быть 1")
}
