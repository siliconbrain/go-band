package band

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBand_DisbandWithoutMembers(t *testing.T) {
	b := New()
	require.True(t, open(b.Disbanding()))
	require.True(t, open(b.Disbanded()))
	b.Disband(42)
	require.False(t, open(b.Disbanding()))
	require.False(t, open(b.Disbanded()))
	require.Equal(t, 42, b.Follow())
}

func TestBand_DisbandWithMembers(t *testing.T) {
	b := New()
	m1 := b.Join()
	m2 := b.Join()
	require.True(t, open(b.Disbanding()))
	require.True(t, open(b.Disbanded()))
	b.Disband(42)
	require.False(t, open(b.Disbanding()))
	require.True(t, open(b.Disbanded()))
	m1.Leave()
	require.False(t, open(b.Disbanding()))
	require.True(t, open(b.Disbanded()))
	m2.Leave()
	require.False(t, open(b.Disbanding()))
	require.False(t, open(b.Disbanded()))
}

func TestBand_Collab(t *testing.T) {
	b := New()
	require.Equal(t, 0, b.members)
	b.Collab(func() {
		require.Equal(t, 1, b.members)
	})
	require.Equal(t, 0, b.members)
}
