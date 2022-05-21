package basa_dannih

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasaDannihWriteToBd(t *testing.T) {
	bd := BasaDannihMySQLPostgresMongoPgAdmin777{}
	bd.WriteToBD(1)
	assert.True(t, bd[1])
	assert.Equal(t, 1, len(bd))
}

func TestReadFromBdFalse(t *testing.T) {
	bd := BasaDannihMySQLPostgresMongoPgAdmin777{}
	assert.False(t, bd.ReadFromBd(1))
	assert.Equal(t, 0, len(bd))
}

func TestReadFromBdTrue(t *testing.T) {
	bd := BasaDannihMySQLPostgresMongoPgAdmin777{}
	bd[1] = true
	assert.True(t, bd[1])
}
