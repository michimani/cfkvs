package libs

import (
	"fmt"
	"os"

	"github.com/michimani/cfkvs/types"
)

func GetKeyValueStoreDataFromFile(path string) (*types.KeyValueStoreData, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	bodyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	kvsData := types.KeyValueStoreData{}
	if err := kvsData.FromBytes(bodyBytes); err != nil {
		return nil, err
	}

	return &kvsData, nil
}
