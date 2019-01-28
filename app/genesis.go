package app

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cybercongress/cyberd/types/coin"
	"github.com/pkg/errors"
	tmtypes "github.com/tendermint/tendermint/types"
	"time"
)

// State to Unmarshal
type GenesisState struct {
	Accounts     []GenesisAccount      `json:"accounts"`
	AuthData     auth.GenesisState     `json:"auth"`
	DistrData    distr.GenesisState    `json:"distr"`
	MintData     mint.GenesisState     `json:"mint"`
	StakingData  staking.GenesisState  `json:"staking"`
	SlashingData slashing.GenesisState `json:"slashing"`
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
	distrData distr.GenesisState,
	slashingData slashing.GenesisState,
) GenesisState {

	return GenesisState{
		Accounts:     accounts,
		AuthData:     authData,
		StakingData:  stakingData,
		MintData:     mintData,
		DistrData:    distrData,
		SlashingData: slashingData,
	}
}

type GenesisAccount struct {
	Address sdk.AccAddress `json:"addr"`
	Amount  int64          `json:"amt"`
}

func NewGenesisAccount(acc *auth.BaseAccount) GenesisAccount {
	return GenesisAccount{
		Address: acc.Address,
		Amount:  acc.Coins.AmountOf(coin.CYB).Int64(),
	}
}

func NewGenesisAccountI(acc auth.Account) GenesisAccount {
	return GenesisAccount{
		Address: acc.GetAddress(),
		Amount:  acc.GetCoins().AmountOf(coin.CYB).Int64(),
	}
}

// convert GenesisAccount to auth.BaseAccount
func (ga *GenesisAccount) ToAccount() (acc *auth.BaseAccount) {
	return &auth.BaseAccount{
		Address: ga.Address,
		Coins:   sdk.Coins{sdk.NewInt64Coin(coin.CYB, ga.Amount)},
	}
}

const (
	// defaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	defaultUnbondingTime = 60 * 60 * 24 * 3 * time.Second
)

// todo set params for each module
// NewDefaultGenesisState generates the default state for cyberd.
func NewDefaultGenesisState() GenesisState {
	return GenesisState{
		Accounts: nil,
		AuthData: auth.GenesisState{
			Params: auth.Params{
				MaxMemoCharacters: 256,
			},
		},
		MintData: mint.GenesisState{
			Params: mint.Params{
				MintDenom:           coin.CYB,
				InflationRateChange: sdk.NewDecWithPrec(0, 2),
				InflationMax:        sdk.NewDecWithPrec(200, 2),
				InflationMin:        sdk.NewDecWithPrec(200, 2),
				GoalBonded:          sdk.NewDecWithPrec(99, 2),
				BlocksPerYear:       uint64(60 * 60 * 24 * 365), // assuming 1 second block times
			},
		},
		StakingData: staking.GenesisState{
			Pool: staking.InitialPool(),
			Params: types.Params{
				UnbondingTime: defaultUnbondingTime,
				MaxValidators: 146,
				BondDenom:     coin.CYB,
			},
		},
		SlashingData: slashing.DefaultGenesisState(),
		DistrData:    distr.DefaultGenesisState(),
		GenTxs:       nil,
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

	stakeData := genesisState.StakingData
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

	for _, acc := range genesisState.Accounts {
		stakeData.Pool.NotBondedTokens = stakeData.Pool.NotBondedTokens.Add(sdk.NewInt(acc.Amount))
	}
	genesisState.StakingData = stakeData
	genesisState.GenTxs = appGenTxs
	return genesisState, nil
}

// CyberdAppGenState but with JSON
func CyberdAppGenStateJSON(cdc *codec.Codec, genDoc tmtypes.GenesisDoc, appGenTxs []json.RawMessage) (
	appState json.RawMessage, err error) {
	// create the final app state
	genesisState, err := CyberdAppGenState(cdc, genDoc, appGenTxs)
	if err != nil {
		return nil, err
	}
	return codec.MarshalJSONIndent(cdc, genesisState)
}

// CyberdValidateGenesisState ensures that the genesis state obeys the expected invariants
func CyberdValidateGenesisState(genesisState GenesisState) (err error) {

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
