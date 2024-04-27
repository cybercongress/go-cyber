package keeper_test

//type KeeperTestSuite struct {
//	apptesting.KeeperTestHelper
//
//	queryClient     types.QueryClient
//	bankQueryClient banktypes.QueryClient
//	msgServer       types.MsgServer
//	// defaultDenom is on the suite, as it depends on the creator test address.
//	defaultDenom string
//}
//
//func TestKeeperTestSuite(t *testing.T) {
//	suite.Run(t, new(KeeperTestSuite))
//}
//
//func (suite *KeeperTestSuite) SetupTest() {
//	suite.Setup()
//
//	// Fund every TestAcc with two denoms, one of which is the denom creation fee
//	fundAccsAmount := sdk.NewCoins(sdk.NewCoin(types.DefaultParams().DenomCreationFee[0].Denom, types.DefaultParams().DenomCreationFee[0].Amount.MulRaw(100)), sdk.NewCoin(apptesting.SecondaryDenom, apptesting.SecondaryAmount))
//	for _, acc := range suite.TestAccs {
//		suite.FundAcc(acc, fundAccsAmount)
//	}
//
//	suite.queryClient = types.NewQueryClient(suite.QueryHelper)
//	suite.bankQueryClient = banktypes.NewQueryClient(suite.QueryHelper)
//	suite.msgServer = keeper.NewMsgServerImpl(suite.App.AppKeepers.TokenFactoryKeeper)
//}
//
//func (suite *KeeperTestSuite) CreateDefaultDenom() {
//	res, _ := suite.msgServer.CreateDenom(sdk.WrapSDKContext(suite.Ctx), types.NewMsgCreateDenom(suite.TestAccs[0].String(), "bitcoin"))
//	suite.defaultDenom = res.GetNewTokenDenom()
//}
//
//func (suite *KeeperTestSuite) TestCreateModuleAccount() {
//	app := suite.App
//
//	// remove module account
//	tokenfactoryModuleAccount := app.AppKeepers.AccountKeeper.GetAccount(suite.Ctx, app.AppKeepers.AccountKeeper.GetModuleAddress(types.ModuleName))
//	app.AppKeepers.AccountKeeper.RemoveAccount(suite.Ctx, tokenfactoryModuleAccount)
//
//	// ensure module account was removed
//	suite.Ctx = app.BaseApp.NewContext(false, tmproto.Header{ChainID: "testing"})
//	tokenfactoryModuleAccount = app.AppKeepers.AccountKeeper.GetAccount(suite.Ctx, app.AppKeepers.AccountKeeper.GetModuleAddress(types.ModuleName))
//	suite.Require().Nil(tokenfactoryModuleAccount)
//
//	// create module account
//	app.AppKeepers.TokenFactoryKeeper.CreateModuleAccount(suite.Ctx)
//
//	// check that the module account is now initialized
//	tokenfactoryModuleAccount = app.AppKeepers.AccountKeeper.GetAccount(suite.Ctx, app.AppKeepers.AccountKeeper.GetModuleAddress(types.ModuleName))
//	suite.Require().NotNil(tokenfactoryModuleAccount)
//}
