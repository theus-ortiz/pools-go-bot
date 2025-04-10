package pools

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const dataDir = "data/pools"

func ensureDataDir() error {
	return os.MkdirAll(dataDir, os.ModePerm)
}

func getUserFilePath(userID string) string {
	return filepath.Join(dataDir, fmt.Sprintf("%s.json", userID))
}

func LoadUserPools(userID string) (*UserPools, error) {
	if err := ensureDataDir(); err != nil {
		return nil, err
	}

	path := getUserFilePath(userID)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return &UserPools{Owner: userID, Positions: []Position{}}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var userPools UserPools
	if err := json.Unmarshal(data, &userPools); err != nil {
		return nil, err
	}

	return &userPools, nil
}

func SaveUserPools(pools *UserPools) error {
	if err := ensureDataDir(); err != nil {
		return err
	}

	path := getUserFilePath(pools.Owner)
	data, err := json.MarshalIndent(pools, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func AddPool(userID, address, network string) error {
	userPools, err := LoadUserPools(userID)
	if err != nil {
		return err
	}

	// Evita duplicidade
	for _, pos := range userPools.Positions {
		if pos.Address == address && pos.Network == network {
			return fmt.Errorf("esta pool já foi adicionada")
		}
	}

	newPosition := Position{
		ID:      fmt.Sprintf("%d", len(userPools.Positions)+1),
		Address: address,
		Network: network,
	}

	userPools.Positions = append(userPools.Positions, newPosition)

	return SaveUserPools(userPools)
}