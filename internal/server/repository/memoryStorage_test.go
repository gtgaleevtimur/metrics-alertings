package repository

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewServerMemStorage(t *testing.T) {
	ms := NewServerMemStorage()
	require.NotNil(t, ms)
}

func TestList(t *testing.T) {
	ms := NewServerMemStorage()
	res := ms.List()
	require.NotNil(t, res)
}

func TestGet(t *testing.T) {
	ms := NewServerMemStorage()
	res, err := ms.Get("test")
	require.Equal(t, err.Error(), "not found")
	require.Nil(t, res)

}

func TestUpdate(t *testing.T) {
	ms := NewServerMemStorage()
	err := ms.Update("gauge", "test", 1.1)
	require.NoError(t, err)
}
