package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDSN(t *testing.T) {
	tests := []struct {
		Name     string
		Input    string
		Expected *DSN
		WantErr  bool
	}{
		{
			Name:  "URI",
			Input: "postgres://postgres:r5$$gkt5$$4f4R43@infra-postgres18/aaaa_bbbb",
			Expected: &DSN{
				DatabaseName:             "aaaa_bbbb",
				PostgresConnectionString: "postgres://postgres:r5$$gkt5$$4f4R43@infra-postgres18",
			},
		},
		{
			Name:  "URI with query parameters",
			Input: "postgresql://postgres:secret@postgres:5432/service?sslmode=require",
			Expected: &DSN{
				DatabaseName:             "service",
				PostgresConnectionString: "postgresql://postgres:secret@postgres:5432?sslmode=require",
			},
		},
		{
			Name:  "keyword value DSN with dbname in the middle",
			Input: "host=postgres user=admin dbname=service password=secret sslmode=disable",
			Expected: &DSN{
				DatabaseName:             "service",
				PostgresConnectionString: "host=postgres user=admin password=secret sslmode=disable",
			},
		},
		{
			Name:  "keyword value DSN with dbname at the end",
			Input: "host=postgres password=secret=part dbname=service",
			Expected: &DSN{
				DatabaseName:             "service",
				PostgresConnectionString: "host=postgres password=secret=part",
			},
		},
		{
			Name:    "URI without database name",
			Input:   "postgres://postgres:secret@postgres:5432",
			WantErr: true,
		},
		{
			Name:    "keyword value DSN without database name",
			Input:   "host=postgres user=admin password=secret",
			WantErr: true,
		},
		{
			Name:    "empty DSN",
			Input:   "",
			WantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got, err := ParseDSN(test.Input)
			if test.WantErr {
				require.Error(t, err)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, test.Expected, got)
		})
	}
}
