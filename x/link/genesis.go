package link

import (
	"bufio"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/io"
	"github.com/cybercongress/cyberd/x/link/keeper"
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

const (
	LinksFileName = "config/links"
)

func InitGenesis(
	ctx sdk.Context, cidNumKeeper keeper.CidNumberKeeper, linkIndexedKeeper *keeper.LinkIndexedKeeper, logger log.Logger,
) (err error) {

	linksFilePath := io.RootifyPath(LinksFileName)
	linksFile, err := os.Open(linksFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info("File with links not found. Empty set will be used")
			return nil
		}
		return
	}
	reader := bufio.NewReader(linksFile) // 4096 bytes chunk size

	// initialize slices to read data
	err = cidNumKeeper.LoadFromReader(ctx, reader)
	if err != nil {
		return
	}

	// Read all links
	err = linkIndexedKeeper.LoadFromReader(ctx, reader)
	if err != nil {
		return
	}

	return
}

func WriteGenesis(
	ctx sdk.Context, cidNumKeeper keeper.CidNumberKeeper, linkIndexedKeeper *keeper.LinkIndexedKeeper, logger log.Logger,
) (err error) {

	linksFilePath := io.RootifyPath(LinksFileName)
	linksFile, err := os.Create(linksFilePath)
	if err != nil {
		return
	}

	writer := bufio.NewWriter(linksFile) // 4096 byte chunk
	err = cidNumKeeper.WriteCids(ctx, writer)
	if err != nil {
		return
	}
	err = linkIndexedKeeper.WriteLinks(ctx, writer)
	if err != nil {
		return
	}

	err = writer.Flush()
	if err != nil {
		return
	}
	err = linksFile.Close()

	logger.Info("Cids and links exported. File created.", "path", linksFilePath)
	return
}
