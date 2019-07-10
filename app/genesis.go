package app

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cybercongress/cyberd/types/coin"
	"github.com/cybercongress/cyberd/util"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
	"io/ioutil"
	"time"
)

const (
	defaultUnbondingTime = 60 * 60 * 24 * 3 * 7 * time.Second // 3 weeks
)

// State to Unmarshal
type GenesisState struct {
	Accounts     []GenesisAccount      `json:"accounts"`
	AuthData     auth.GenesisState     `json:"auth"`
	BankData     bank.GenesisState     `json:"bank"`
	DistrData    distr.GenesisState    `json:"distr"`
	MintData     mint.GenesisState     `json:"mint"`
	StakingData  staking.GenesisState  `json:"staking"`
	SlashingData slashing.GenesisState `json:"slashing"`
	GovData      gov.GenesisState      `json:"gov"`
	GenTxs       []json.RawMessage     `json:"gentxs"`
}

func (gs *GenesisState) GetAddresses() []sdk.AccAddress {
	addresses := make([]sdk.AccAddress, 0, len(gs.Accounts))
	for _, acc := range gs.Accounts {
		addresses = append(addresses, acc.Address)
	}
	return addresses
}

func NewGenesisState(
	accounts []GenesisAccount, authData auth.GenesisState,
	stakingData staking.GenesisState, mintData mint.GenesisState,
	distrData distr.GenesisState, govData gov.GenesisState,
	slashingData slashing.GenesisState,
) GenesisState {

	return GenesisState{
		Accounts:     accounts,
		AuthData:     authData,
		StakingData:  stakingData,
		MintData:     mintData,
		DistrData:    distrData,
		SlashingData: slashingData,
		GovData:      govData,
	}
}

type GenesisAccount struct {
	Address   sdk.AccAddress `json:"addr"`
	Amount    int64          `json:"amt"`
	AccNumber uint64         `json:"nmb"`
}

func NewGenesisAccount(acc auth.Account) GenesisAccount {
	return GenesisAccount{
		Address:   acc.GetAddress(),
		Amount:    acc.GetCoins().AmountOf(coin.CYB).Int64(),
		AccNumber: acc.GetAccountNumber(),
	}
}

// convert GenesisAccount to auth.BaseAccount
func (ga *GenesisAccount) ToAccount() (acc *auth.BaseAccount) {
	return &auth.BaseAccount{
		Address:       ga.Address,
		Coins:         sdk.Coins{sdk.NewInt64Coin(coin.CYB, ga.Amount)},
		AccountNumber: ga.AccNumber,
	}
}

// NewDefaultGenesisState generates the default state for cyberd.
func NewDefaultGenesisState() GenesisState {
	return GenesisState{
		Accounts: nil,
		AuthData: auth.GenesisState{
			Params: auth.Params{
				MaxMemoCharacters: 256,
			},
		},
		BankData: bank.GenesisState{
			SendEnabled: true,
		},
		MintData: mint.GenesisState{
			Minter: mint.Minter{
				Inflation: sdk.NewDecWithPrec(13, 2),
				AnnualProvisions: sdk.NewDec(0),
			},
			Params: mint.Params{
				MintDenom:           coin.CYB,
				InflationRateChange: sdk.NewDecWithPrec(13, 2),
				InflationMax:        sdk.NewDecWithPrec(20, 2),
				InflationMin:        sdk.NewDecWithPrec(7, 2),
				GoalBonded:          sdk.NewDecWithPrec(67, 2),
				BlocksPerYear:       uint64(60 * 60 * 8766 / 5), // assuming 5 second block times
			},
		},
		StakingData: staking.GenesisState{
			Pool: staking.InitialPool(),
			Params: types.Params{
				UnbondingTime: defaultUnbondingTime,
				MaxValidators: 146,
				MaxEntries:    7,
				BondDenom:     coin.CYB,
			},
		},
		SlashingData: slashing.GenesisState{
			Params: slashing.Params{
				MaxEvidenceAge:          defaultUnbondingTime,
				SignedBlocksWindow:      60 * 30, // ~30min
				DowntimeJailDuration:    0,
				MinSignedPerWindow:      sdk.NewDecWithPrec(70, 2),           // 70%
				SlashFractionDoubleSign: sdk.NewDecWithPrec(20, 2),           // 20%
				SlashFractionDowntime:   sdk.NewDec(1).Quo(sdk.NewDec(1000)), // 0.1%
			},
		},
		DistrData: distr.GenesisState{
			FeePool:             distr.InitialFeePool(),
			CommunityTax:        sdk.NewDecWithPrec(0, 2), // 0%
			BaseProposerReward:  sdk.NewDecWithPrec(1, 2), // 1%
			BonusProposerReward: sdk.NewDecWithPrec(4, 2), // 4%
			WithdrawAddrEnabled: true,
			PreviousProposer:    nil,
		},
		GovData: gov.GenesisState{
			StartingProposalID: 1,
			DepositParams: gov.DepositParams{
				MinDeposit:       sdk.Coins{coin.NewCybCoin(500 * coin.Giga)}, //top 50 of current network
				MaxDepositPeriod: gov.DefaultPeriod,
			},
			VotingParams: gov.VotingParams{
				VotingPeriod: gov.DefaultPeriod,
			},
			TallyParams: gov.TallyParams{
				Quorum:    sdk.NewDecWithPrec(334, 3),
				Threshold: sdk.NewDecWithPrec(5, 1),
				Veto:      sdk.NewDecWithPrec(334, 3),
			},
		},
		GenTxs: []json.RawMessage{},
	}
}

// Create the core parameters for genesis initialization for cyberd
// note that the pubkey input is this machines pubkey
func CyberdAppGenState(cdc *codec.Codec, genDoc tmtypes.GenesisDoc, appGenTxs []json.RawMessage) (
	genesisState GenesisState, err error) {

	if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
		return genesisState, err
	}

	// if there are no gen txs to be processed, return the default empty state
	if len(appGenTxs) == 0 {
		return genesisState, errors.New("there must be at least one genesis tx")
	}

	for i, genTx := range appGenTxs {
		var tx auth.StdTx
		if err := cdc.UnmarshalJSON(genTx, &tx); err != nil {
			return genesisState, err
		}
		msgs := tx.GetMsgs()
		if len(msgs) != 1 {
			return genesisState, errors.New(
				"must provide genesis StdTx with exactly 1 CreateValidator message")
		}
		if _, ok := msgs[0].(staking.MsgCreateValidator); !ok {
			return genesisState, fmt.Errorf(
				"genesis transaction %v does not contain a MsgCreateValidator", i)
		}
	}

	genesisState.GenTxs = appGenTxs
	return genesisState, nil
}

//todo should be here?
func CyberdAppGenStateJSON(cdc *codec.Codec, genDoc tmtypes.GenesisDoc, appGenTxs []json.RawMessage) (
	appState json.RawMessage, err error) {
	// create the final app state
	genesisState, err := CyberdAppGenState(cdc, genDoc, appGenTxs)
	if err != nil {
		return nil, err
	}
	return codec.MarshalJSONIndent(cdc, genesisState)
}

// validateGenesisState ensures that the genesis state obeys the expected invariants
func validateGenesisState(genesisState GenesisState) (err error) {

	err = validateGenesisStateAccounts(genesisState.Accounts)
	if err != nil {
		return
	}

	if err := staking.ValidateGenesis(genesisState.StakingData); err != nil {
		return err
	}
	if err := mint.ValidateGenesis(genesisState.MintData); err != nil {
		return err
	}
	if err := distr.ValidateGenesis(genesisState.DistrData); err != nil {
		return err
	}
	if err := gov.ValidateGenesis(genesisState.GovData); err != nil {
		return err
	}
	if err := bank.ValidateGenesis(genesisState.BankData); err != nil {
		return err
	}
	if err := mint.ValidateGenesis(genesisState.MintData); err != nil {
		return err
	}

	return staking.ValidateGenesis(genesisState.StakingData)
}

// Ensures that there are no duplicate accounts in the genesis state,
func validateGenesisStateAccounts(accs []GenesisAccount) (err error) {
	addrMap := make(map[string]bool, len(accs))
	for i := 0; i < len(accs); i++ {
		acc := accs[i]
		strAddr := string(acc.Address)
		if _, ok := addrMap[strAddr]; ok {
			return fmt.Errorf("duplicate account in genesis state: Address %v", acc.Address.String())
		}
		addrMap[strAddr] = true
	}
	return
}

func LoadGenesisState(
	ctx *server.Context, cdc *codec.Codec,
) (genDoc tmtypes.GenesisDoc, state GenesisState, err error) {

	config := ctx.Config
	config.SetRoot(viper.GetString(cli.HomeFlag))

	genFile := config.GenesisFile()
	if !common.FileExists(genFile) {
		err = fmt.Errorf("%s does not exist, run `cyberd init` first", genFile)
		return
	}
	genDoc, err = LoadGenesisDoc(cdc, genFile)
	if err != nil {
		return
	}

	err = cdc.UnmarshalJSON(genDoc.AppState, &state)
	return
}

func SaveGenesisState(ctx *server.Context, cdc *codec.Codec, oldDoc tmtypes.GenesisDoc, state GenesisState) error {

	config := ctx.Config
	config.SetRoot(viper.GetString(cli.HomeFlag))
	genFile := config.GenesisFile()

	appStateJSON, err := cdc.MarshalJSON(&state)
	if err != nil {
		return err
	}

	return util.ExportGenesisFile(genFile, oldDoc.ChainID, oldDoc.Validators, appStateJSON)
}

func LoadGenesisDoc(cdc *amino.Codec, genFile string) (genDoc tmtypes.GenesisDoc, err error) {
	genContents, err := ioutil.ReadFile(genFile)
	if err != nil {
		return genDoc, err
	}

	if err := cdc.UnmarshalJSON(genContents, &genDoc); err != nil {
		return genDoc, err
	}

	return genDoc, err
}
