package app

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cybercongress/cyberd/types/coin"
	"github.com/cybercongress/cyberd/x/mint"
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

const (
	// defaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	defaultUnbondingTime = 60 * 60 * 24 * 3 * 7 * time.Second
)

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
				TokensPerBlock: 634195840,
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
