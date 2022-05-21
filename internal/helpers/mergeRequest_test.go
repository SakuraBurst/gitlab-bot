package helpers

import (
	"github.com/SakuraBurst/gitlab-bot/pkg/basa_dannih"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWriteMrsToBd(t *testing.T) {
	bd := basa_dannih.BasaDannihMySQLPostgresMongoPgAdmin777{}
	testMr := models.MergeRequest{
		Iid: 1337,
	}
	WriteMrsToBd(bd, testMr)
	assert.Len(t, bd, 1)
	assert.True(t, bd.ReadFromBd(1337))
}

func TestOnlyNewMrsLength0(t *testing.T) {
	bd := basa_dannih.BasaDannihMySQLPostgresMongoPgAdmin777{}
	mrs, writtenMrs, ok := OnlyNewMrs(nil, bd)
	assert.Nil(t, mrs)
	assert.Zero(t, writtenMrs)
	assert.False(t, ok)
}

func TestOnlyNewMrsWithWrittenMrs(t *testing.T) {
	bd := basa_dannih.BasaDannihMySQLPostgresMongoPgAdmin777{}
	bd.WriteToBD(1337)
	mrs, writtenMrs, ok := OnlyNewMrs([]models.MergeRequest{{
		Iid: 1337,
	}}, bd)
	assert.Equal(t, 1, cap(mrs))
	assert.Zero(t, writtenMrs)
	assert.False(t, ok)
}

func TestOnlyNewMrsWithUnWrittenMrs(t *testing.T) {
	bd := basa_dannih.BasaDannihMySQLPostgresMongoPgAdmin777{}
	mrs, writtenMrs, ok := OnlyNewMrs([]models.MergeRequest{{
		Iid: 1337,
	}}, bd)
	assert.Len(t, mrs, 1)
	assert.Equal(t, 1, cap(mrs))
	assert.Equal(t, writtenMrs, 1)
	assert.Equal(t, 1337, mrs[0].Iid)
	assert.True(t, ok)
}
