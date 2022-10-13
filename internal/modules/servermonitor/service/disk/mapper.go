package disk

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/unknwon/com"
)

type MountpointIDMapper struct {
	filePath string
}

func GetMountpointIDMapper(dataFolder string) (*MountpointIDMapper, error) {
	mFilePath := filepath.Join(dataFolder, "mountpointid")
	if !com.IsFile(mFilePath) {
		if _, err := os.Create(mFilePath); err != nil {
			return nil, err
		}
	}
	mapper := MountpointIDMapper{mFilePath}

	return &mapper, nil
}

func (mp *MountpointIDMapper) GetMountpointID(mountpoint string) (int, error) {
	file, err := os.OpenFile(mp.filePath, os.O_RDWR, 0666)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	mMap := make(map[string]int)
	decoder := json.NewDecoder(file)
	decoder.Decode(&mMap)

	if id, ok := mMap[mountpoint]; ok {
		return id, nil
	}

	maxId := 0
	for _, value := range mMap {
		if value > maxId {
			maxId = value
		}
	}

	mMap[mountpoint] = maxId + 1
	if err = file.Truncate(0); err != nil {
		return 0, err
	}
	if _, err = file.Seek(0, 0); err != nil {
		return 0, err
	}
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&mMap); err != nil {
		return 0, err
	}

	return mMap[mountpoint], nil
}
