package graphql

import (
	"os"
	"path/filepath"
)

func ReadQuery(filename string) (string, error) {
	path := filepath.Join("internal", "graphql", "queries", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
