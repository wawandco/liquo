package liquo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLFunc(t *testing.T) {
	r := require.New(t)
	c := ChangeSet{}
	c.SQL = []string{
		"SELECT 1;",
		"SELECT 2;",
	}

	r.Equal(c.sql(), "SELECT 1;\nSELECT 2;")
}
