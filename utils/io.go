package utils

import (
	"errors"
	"io"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
)

func RootifyPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	rootDir := viper.GetString(cli.HomeFlag)
	return filepath.Join(rootDir, path)
}

func ReadExactlyNBytes(reader io.Reader, n uint64) ([]byte, error) {
	data := make([]byte, n)
	currentlyReadedBytes := uint64(0)
	for currentlyReadedBytes < n {
		readedBytes, err := reader.Read(data[currentlyReadedBytes:n])
		if err != nil {
			return nil, err
		}
		if readedBytes == 0 {
			return nil, errors.New("not enough bytes tor read")
		}
		currentlyReadedBytes += uint64(readedBytes)
	}

	return data, nil
}
