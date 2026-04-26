package internal

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type DSN struct {
	DatabaseName string

	PostgresConnectionString string
}

var parsers = []func(string) (*DSN, error){
	parseDSNGo,
	parseDSNUri,
}

func ParseDSN(dsnVal string) (*DSN, error) {
	errs := make([]error, 0)

	for _, parser := range parsers {
		dsn, err := parser(dsnVal)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		return dsn, nil
	}

	return nil, errors.Join(errs...)
}

func parseDSNUri(dsnVal string) (*DSN, error) {
	dsnURI, err := url.Parse(dsnVal)
	if err != nil {
		return nil, fmt.Errorf("parse dsn as url: %w", err)
	}

	name := dsnURI.Path
	dsnURI.Path = ""

	return &DSN{
		DatabaseName:             name,
		PostgresConnectionString: dsnURI.String(),
	}, nil
}

func parseDSNGo(dsnVal string) (*DSN, error) {
	parts := strings.Split(dsnVal, " ")
	dsn := DSN{}

	for i, part := range parts {
		partKV := strings.SplitN(part, "=", 2)
		if len(partKV) != 2 {
			return nil, fmt.Errorf("invalid value in %d segment", i)
		}

		key := partKV[0]
		value := partKV[1]
		if key == "dbname" {
			dsn.DatabaseName = value
			continue
		}

		dsn.PostgresConnectionString += part
		if i < len(parts)-1 {
			dsn.PostgresConnectionString += " "
		}
	}

	if dsn.DatabaseName == "" {
		return nil, errors.New("database name not found")
	}

	return &dsn, nil
}
