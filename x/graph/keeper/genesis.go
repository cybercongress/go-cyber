package keeper

import (
	"bufio"
	"os"
	"path/filepath"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/util"
)

const (
	LinksFileName       = "config/graph"
	LinksExportFileName = "export/graph"
)

func InitGenesis(
	ctx sdk.Context, gk GraphKeeper, ik *IndexKeeper,
) (err error) {
	linksFilePath := util.RootifyPath(LinksFileName)
	linksFile, err := os.Open(linksFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			gk.Logger(ctx).Info("File with cyberlinks not found. Empty set will be used")
			return nil
		}
		return err
	}
	reader := bufio.NewReader(linksFile) // 4096 bytes chunk size

	// initialize slices to read data
	err = gk.LoadFromReader(ctx, reader)
	if err != nil {
		return err
	}

	// Read all links
	err = ik.LoadFromReader(ctx, reader)
	if err != nil {
		return err
	}
	return
}

func WriteGenesis(
	ctx sdk.Context, gk GraphKeeper, ik *IndexKeeper,
) (err error) {

	linksFilePath := util.RootifyPath(LinksExportFileName)
	dirName := filepath.Dir(linksFilePath)
	if _, err := os.Stat(dirName); err != nil {
		if err = os.MkdirAll(dirName, os.ModePerm); err != nil {
			return err
		}
	}

	linksFile, err := os.Create(linksFilePath)
	if err != nil {
		return
	}

	writer := bufio.NewWriter(linksFile) // 4096 byte chunk
	err = gk.WriteCids(ctx, writer)
	if err != nil {
		return
	}
	err = ik.WriteLinks(ctx, writer)
	if err != nil {
		return
	}

	err = writer.Flush()
	if err != nil {
		return
	}
	err = linksFile.Close()

	gk.Logger(ctx).Info("CIDs and cyberlinks exported. File created.", "path", linksFilePath)
	return
}
