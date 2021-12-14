package keeper

import (
	"bufio"
	"encoding/csv"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cyberbankkeeper "github.com/cybercongress/go-cyber/x/cyberbank/keeper"
	"github.com/cybercongress/go-cyber/x/graph/types"
	"os"
	"path/filepath"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/utils"
)

const (
	LinksFileName       = "config/graph"
	LinksExportFileName = "export/graph"
)

// TODO make refactoring of initial rank calculation on graph load from binary
func InitGenesis(
	ctx sdk.Context, gk GraphKeeper, ik *IndexKeeper,
) (err error) {
	linksFilePath := utils.RootifyPath(LinksFileName)
	linksFile, err := os.Open(linksFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			gk.Logger(ctx).Info("File with cyberlinks and particles not found. Empty set will be used")
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

	linksFilePath := utils.RootifyPath(LinksExportFileName)
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

	gk.Logger(ctx).Info("Particles and cyberlinks exported. File created.", "path", linksFilePath)
	return
}

func WriteGraphCSV(
	ctx sdk.Context,
	gk GraphKeeper,
	ik *IndexKeeper,
	ak authkeeper.AccountKeeper,
	bk *cyberbankkeeper.IndexedKeeper,
) (err error) {
	linksFilePath := utils.RootifyPath(LinksExportFileName)
	dirName := filepath.Dir(linksFilePath)
	if _, err := os.Stat(dirName); err != nil {
		if err = os.MkdirAll(dirName, os.ModePerm); err != nil {
			return err
		}
	}

	linksFile, err := os.Create(linksFilePath)
	if err != nil {
		return err
	}

	var addressMap = make(map[uint64]sdk.AccAddress, 0)
	ak.IterateAccounts(ctx, func(account authtypes.AccountI) bool {
		addressMap[account.GetAccountNumber()] = account.GetAddress()
		return false
	})

	w := csv.NewWriter(linksFile)
	defer w.Flush()

	gk.IterateLinks(ctx, func(link types.CompactLink){
		var row []string
		row = append(row, string(ik.GetCid(ctx, types.CidNumber(link.From))))
		row = append(row, string(ik.GetCid(ctx, types.CidNumber(link.To))))
		row = append(row, addressMap[link.Account].String())
		row = append(row, strconv.FormatUint(link.From, 10))
		row = append(row, strconv.FormatUint(link.To, 10))
		row = append(row, strconv.FormatUint(link.Account, 10))
		row = append(row, strconv.FormatUint(uint64(bk.GetAccountTotalStakeAmper(ctx, addressMap[link.Account])), 10))
		if err := w.Write(row); err != nil {
			return
		}
	})
	return nil
}