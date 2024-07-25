package keeper_test

import (
	"crypto/sha256"
	"testing"
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/stretchr/testify/suite"

	_ "embed"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/cybercongress/go-cyber/v4/app"
	"github.com/cybercongress/go-cyber/v4/x/clock/keeper"
	"github.com/cybercongress/go-cyber/v4/x/clock/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx            sdk.Context
	app            *app.App
	bankKeeper     bankkeeper.Keeper
	queryClient    types.QueryClient
	clockMsgServer types.MsgServer
}

func (s *IntegrationTestSuite) SetupTest() {
	isCheckTx := false
	s.app = app.Setup(s.T())

	s.ctx = s.app.BaseApp.NewContext(isCheckTx, tmproto.Header{
		ChainID: "testing",
		Height:  1,
		Time:    time.Now().UTC(),
	})

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, s.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.NewQuerier(s.app.AppKeepers.ClockKeeper))

	s.queryClient = types.NewQueryClient(queryHelper)
	s.bankKeeper = s.app.AppKeepers.BankKeeper
	s.clockMsgServer = keeper.NewMsgServerImpl(s.app.AppKeepers.ClockKeeper)
}

func (s *IntegrationTestSuite) FundAccount(ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := s.bankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}

	return s.bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amounts)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

//go:embed testdata/clock_example.wasm
var wasmContract []byte

func (s *IntegrationTestSuite) StoreCode() {
	_, _, sender := testdata.KeyTestPubAddr()
	msg := wasmtypes.MsgStoreCodeFixture(func(m *wasmtypes.MsgStoreCode) {
		m.WASMByteCode = wasmContract
		m.Sender = sender.String()
	})
	rsp, err := s.app.MsgServiceRouter().Handler(msg)(s.ctx, msg)
	s.Require().NoError(err)
	var result wasmtypes.MsgStoreCodeResponse
	s.Require().NoError(s.app.AppCodec().Unmarshal(rsp.Data, &result))
	s.Require().Equal(uint64(1), result.CodeID)
	expHash := sha256.Sum256(wasmContract)
	s.Require().Equal(expHash[:], result.Checksum)
	// and
	info := s.app.AppKeepers.WasmKeeper.GetCodeInfo(s.ctx, 1)
	s.Require().NotNil(info)
	s.Require().Equal(expHash[:], info.CodeHash)
	s.Require().Equal(sender.String(), info.Creator)
	s.Require().Equal(wasmtypes.DefaultParams().InstantiateDefaultPermission.With(sender), info.InstantiateConfig)
}

func (s *IntegrationTestSuite) InstantiateContract(sender string, admin string) string {
	msgStoreCode := wasmtypes.MsgStoreCodeFixture(func(m *wasmtypes.MsgStoreCode) {
		m.WASMByteCode = wasmContract
		m.Sender = sender
	})
	_, err := s.app.MsgServiceRouter().Handler(msgStoreCode)(s.ctx, msgStoreCode)
	s.Require().NoError(err)

	msgInstantiate := wasmtypes.MsgInstantiateContractFixture(func(m *wasmtypes.MsgInstantiateContract) {
		m.Sender = sender
		m.Admin = admin
		m.Msg = []byte(`{}`)
	})
	resp, err := s.app.MsgServiceRouter().Handler(msgInstantiate)(s.ctx, msgInstantiate)
	s.Require().NoError(err)
	var result wasmtypes.MsgInstantiateContractResponse
	s.Require().NoError(s.app.AppCodec().Unmarshal(resp.Data, &result))
	contractInfo := s.app.AppKeepers.WasmKeeper.GetContractInfo(s.ctx, sdk.MustAccAddressFromBech32(result.Address))
	s.Require().Equal(contractInfo.CodeID, uint64(1))
	s.Require().Equal(contractInfo.Admin, admin)
	s.Require().Equal(contractInfo.Creator, sender)

	return result.Address
}

// Helper method for quickly registering a clock contract
func (s *IntegrationTestSuite) RegisterClockContract(senderAddress string, contractAddress string) {
	err := s.app.AppKeepers.ClockKeeper.RegisterContract(s.ctx, senderAddress, contractAddress)
	s.Require().NoError(err)
}

// Helper method for quickly unregistering a clock contract
func (s *IntegrationTestSuite) UnregisterClockContract(senderAddress string, contractAddress string) {
	err := s.app.AppKeepers.ClockKeeper.UnregisterContract(s.ctx, senderAddress, contractAddress)
	s.Require().NoError(err)
}

// Helper method for quickly jailing a clock contract
func (s *IntegrationTestSuite) JailClockContract(contractAddress string) {
	err := s.app.AppKeepers.ClockKeeper.SetJailStatus(s.ctx, contractAddress, true)
	s.Require().NoError(err)
}

// Helper method for quickly unjailing a clock contract
func (s *IntegrationTestSuite) UnjailClockContract(senderAddress string, contractAddress string) {
	err := s.app.AppKeepers.ClockKeeper.SetJailStatusBySender(s.ctx, senderAddress, contractAddress, false)
	s.Require().NoError(err)
}
