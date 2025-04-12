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
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	path := getUserFilePath(userID)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &UserPools{Owner: userID}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read user file: %w", err)
	}

	var userPools UserPools
	if err := json.Unmarshal(data, &userPools); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user pools: %w", err)
	}

	return &userPools, nil
}

func SaveUserPools(pools *UserPools) error {
	if err := ensureDataDir(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(pools, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(getUserFilePath(pools.Owner), data, 0644)
}

func AddPosition(userID, address, network string) error {
	userPools, err := LoadUserPools(userID)
	if err != nil {
		return err
	}

	for _, pos := range userPools.Positions {
		if pos.Address == address && pos.Network == network {
			return fmt.Errorf("position already exists")
		}
	}

	userPools.Positions = append(userPools.Positions, Position{
		Address: address,
		Network: network,
	})

	return SaveUserPools(userPools)
}
