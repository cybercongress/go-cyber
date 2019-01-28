package util

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"io"
	"path/filepath"
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
	bytesReaded, err := reader.Read(data)
	if err != nil {
		return nil, err
	}
	if uint64(bytesReaded) != n {
		return nil, errors.New("not enough bytes tor read")
	}
	return data, nil
}
