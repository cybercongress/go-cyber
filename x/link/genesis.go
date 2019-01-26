package link

import (
	"bytes"
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/cyberd/store"
	"github.com/cybercongress/cyberd/x/link/keeper"
	. "github.com/cybercongress/cyberd/x/link/types"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	"io/ioutil"
	"os"
	"path/filepath"
)

const LinksFileName = "links"

func GetLinksFilePath() string {
	rootDir := viper.GetString(cli.HomeFlag)
	return filepath.Join(rootDir, LinksFileName)
}

func InitGenesis(
	ctx sdk.Context, mainKeeper store.MainKeeper, cidNumKeeper keeper.CidNumberKeeper,
	linkIndexedKeeper *keeper.LinkIndexedKeeper, logger log.Logger,
) {
	// initialize links

	_, err := os.Stat(GetLinksFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info("File with links not found. Empty links set will be used")
		}
	} else {
		data, err := ioutil.ReadFile(GetLinksFilePath())
		if err != nil {
			panic(err)
		}

		// import cids
		cidCount := binary.LittleEndian.Uint64(data[0:8])

		cursor := uint64(8)
		// read all CIDs with their numbers
		for ; cursor < 8+(cidCount*54); cursor = cursor + 54 {
			cid := Cid(data[cursor : cursor+46])
			cidNumber := CidNumber(binary.LittleEndian.Uint64(data[cursor+46 : cursor+54]))
			cidNumKeeper.PutCid(ctx, cid, cidNumber)
		}
		lastCidIndex := cidCount - 1
		lastCidIndexBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(lastCidIndexBytes, lastCidIndex)
		mainKeeper.SetLastCidIndex(ctx, lastCidIndexBytes)

		for ; cursor < uint64(len(data)); cursor = cursor + 24 {
			linkIndexedKeeper.PutLink(ctx, UnmarshalBinaryLink(data[cursor:cursor+24]))
		}
	}
}

func ExportLinks(
	ctx sdk.Context, cidNumKeeper keeper.CidNumberKeeper, linkIndexedKeeper *keeper.LinkIndexedKeeper,
) (err error) {

	linksFile, err := os.Create(GetLinksFilePath())
	if err != nil {
		return
	}

	buf := new(bytes.Buffer)
	dataAsBytes := make([]byte, 8) //common bytes array to convert uints

	// export CIDs in binary format: <n><cid1><cid1_number><cid2><cid2_number>....<cidn><cidn_number>
	cidsCount := cidNumKeeper.GetCidsCount(ctx)
	binary.LittleEndian.PutUint64(dataAsBytes, cidsCount)
	buf.Write(dataAsBytes)

	cidNumKeeper.Iterate(ctx, func(cid Cid, number CidNumber) {
		buf.Write([]byte(cid))
		binary.LittleEndian.PutUint64(dataAsBytes, uint64(number))
		buf.Write(dataAsBytes)
	})

	// export links in binary format: <cid_number_from><cid_number_to><acc_number>
	linkIndexedKeeper.Iterate(ctx, func(link CompactLink) { buf.Write(link.MarshalBinary()) })

	_, err = linksFile.Write(buf.Bytes())
	if err != nil {
		return
	}
	err = linksFile.Close()
	return
}
