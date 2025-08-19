package keeper_test

import (
	_ "embed"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/v6/x/clock/types"
)

// Test register clock contract.
func (s *IntegrationTestSuite) TestRegisterClockContract() {
	_, _, addr := testdata.KeyTestPubAddr()
	_, _, addr2 := testdata.KeyTestPubAddr()
	_ = s.FundAccount(s.ctx, addr, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1_000_000))))

	// Store code
	s.StoreCode()
	contractAddress := s.InstantiateContract(addr.String(), "")
	contractAddressWithAdmin := s.InstantiateContract(addr.String(), addr2.String())

	for _, tc := range []struct {
		desc     string
		sender   string
		contract string
		isJailed bool
		success  bool
	}{
		{
			desc:     "Success - Register Contract",
			sender:   addr.String(),
			contract: contractAddress,
			success:  true,
		},
		{
			desc:     "Success - Register Contract With Admin",
			sender:   addr2.String(),
			contract: contractAddressWithAdmin,
			success:  true,
		},
		{
			desc:     "Fail - Register Contract With Admin, But With Creator Addr",
			sender:   addr.String(),
			contract: contractAddressWithAdmin,
			success:  false,
		},
		{
			desc:     "Error - Invalid Sender",
			sender:   addr2.String(),
			contract: contractAddress,
			success:  false,
		},
		{
			desc:     "Fail - Invalid Contract Address",
			sender:   addr.String(),
			contract: "Invalid",
			success:  false,
		},
		{
			desc:     "Fail - Invalid Sender Address",
			sender:   "Invalid",
			contract: contractAddress,
			success:  false,
		},
		{
			desc:     "Fail - Contract Already Jailed",
			sender:   addr.String(),
			contract: contractAddress,
			isJailed: true,
			success:  false,
		},
	} {
		tc := tc
		s.Run(tc.desc, func() {
			// Set params
			params := types.DefaultParams()
			err := s.app.AppKeepers.ClockKeeper.SetParams(s.ctx, params)
			s.Require().NoError(err)

			// Jail contract if needed
			if tc.isJailed {
				s.RegisterClockContract(tc.sender, tc.contract)
				err := s.app.AppKeepers.ClockKeeper.SetJailStatus(s.ctx, tc.contract, true)
				s.Require().NoError(err)
			}

			// Try to register contract
			res, err := s.clockMsgServer.RegisterClockContract(s.ctx, &types.MsgRegisterClockContract{
				SenderAddress:   tc.sender,
				ContractAddress: tc.contract,
			})

			if !tc.success {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(res, &types.MsgRegisterClockContractResponse{})
			}

			// Ensure contract is unregistered
			s.app.AppKeepers.ClockKeeper.RemoveContract(s.ctx, contractAddress)
			s.app.AppKeepers.ClockKeeper.RemoveContract(s.ctx, contractAddressWithAdmin)
		})
	}
}

// Test standard unregistration of clock contracts.
func (s *IntegrationTestSuite) TestUnregisterClockContract() {
	_, _, addr := testdata.KeyTestPubAddr()
	_, _, addr2 := testdata.KeyTestPubAddr()
	_ = s.FundAccount(s.ctx, addr, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1_000_000))))

	s.StoreCode()
	contractAddress := s.InstantiateContract(addr.String(), "")
	contractAddressWithAdmin := s.InstantiateContract(addr.String(), addr2.String())

	for _, tc := range []struct {
		desc     string
		sender   string
		contract string
		success  bool
	}{
		{
			desc:     "Success - Unregister Contract",
			sender:   addr.String(),
			contract: contractAddress,
			success:  true,
		},
		{
			desc:     "Success - Unregister Contract With Admin",
			sender:   addr2.String(),
			contract: contractAddressWithAdmin,
			success:  true,
		},
		{
			desc:     "Fail - Unregister Contract With Admin, But With Creator Addr",
			sender:   addr.String(),
			contract: contractAddressWithAdmin,
			success:  false,
		},
		{
			desc:     "Error - Invalid Sender",
			sender:   addr2.String(),
			contract: contractAddress,
			success:  false,
		},
		{
			desc:     "Fail - Invalid Contract Address",
			sender:   addr.String(),
			contract: "Invalid",
			success:  false,
		},
		{
			desc:     "Fail - Invalid Sender Address",
			sender:   "Invalid",
			contract: contractAddress,
			success:  false,
		},
	} {
		tc := tc
		s.Run(tc.desc, func() {
			s.RegisterClockContract(addr.String(), contractAddress)
			s.RegisterClockContract(addr2.String(), contractAddressWithAdmin)

			// Set params
			params := types.DefaultParams()
			err := s.app.AppKeepers.ClockKeeper.SetParams(s.ctx, params)
			s.Require().NoError(err)

			// Try to register all contracts
			res, err := s.clockMsgServer.UnregisterClockContract(s.ctx, &types.MsgUnregisterClockContract{
				SenderAddress:   tc.sender,
				ContractAddress: tc.contract,
			})

			if !tc.success {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(res, &types.MsgUnregisterClockContractResponse{})
			}

			// Ensure contract is unregistered
			s.app.AppKeepers.ClockKeeper.RemoveContract(s.ctx, contractAddress)
			s.app.AppKeepers.ClockKeeper.RemoveContract(s.ctx, contractAddressWithAdmin)
		})
	}
}

// Test duplicate register/unregister clock contracts.
func (s *IntegrationTestSuite) TestDuplicateRegistrationChecks() {
	_, _, addr := testdata.KeyTestPubAddr()
	_ = s.FundAccount(s.ctx, addr, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1_000_000))))

	s.StoreCode()
	contractAddress := s.InstantiateContract(addr.String(), "")

	// Test double register, first succeed, second fail
	_, err := s.clockMsgServer.RegisterClockContract(s.ctx, &types.MsgRegisterClockContract{
		SenderAddress:   addr.String(),
		ContractAddress: contractAddress,
	})
	s.Require().NoError(err)

	_, err = s.clockMsgServer.RegisterClockContract(s.ctx, &types.MsgRegisterClockContract{
		SenderAddress:   addr.String(),
		ContractAddress: contractAddress,
	})
	s.Require().Error(err)

	// Test double unregister, first succeed, second fail
	_, err = s.clockMsgServer.UnregisterClockContract(s.ctx, &types.MsgUnregisterClockContract{
		SenderAddress:   addr.String(),
		ContractAddress: contractAddress,
	})
	s.Require().NoError(err)

	_, err = s.clockMsgServer.UnregisterClockContract(s.ctx, &types.MsgUnregisterClockContract{
		SenderAddress:   addr.String(),
		ContractAddress: contractAddress,
	})
	s.Require().Error(err)
}

// Test unjailing clock contracts.
func (s *IntegrationTestSuite) TestUnjailClockContract() {
	_, _, addr := testdata.KeyTestPubAddr()
	_, _, addr2 := testdata.KeyTestPubAddr()
	_ = s.FundAccount(s.ctx, addr, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1_000_000))))

	s.StoreCode()
	contractAddress := s.InstantiateContract(addr.String(), "")
	contractAddressWithAdmin := s.InstantiateContract(addr.String(), addr2.String())

	for _, tc := range []struct {
		desc     string
		sender   string
		contract string
		unjail   bool
		success  bool
	}{
		{
			desc:     "Success - Unjail Contract",
			sender:   addr.String(),
			contract: contractAddress,
			success:  true,
		},
		{
			desc:     "Success - Unjail Contract With Admin",
			sender:   addr2.String(),
			contract: contractAddressWithAdmin,
			success:  true,
		},
		{
			desc:     "Fail - Unjail Contract With Admin, But With Creator Addr",
			sender:   addr.String(),
			contract: contractAddressWithAdmin,
			success:  false,
		},
		{
			desc:     "Error - Invalid Sender",
			sender:   addr2.String(),
			contract: contractAddress,
			success:  false,
		},
		{
			desc:     "Fail - Invalid Contract Address",
			sender:   addr.String(),
			contract: "Invalid",
			success:  false,
		},
		{
			desc:     "Fail - Invalid Sender Address",
			sender:   "Invalid",
			contract: contractAddress,
			success:  false,
		},
		{
			desc:     "Fail - Contract Not Jailed",
			sender:   addr.String(),
			contract: contractAddress,
			unjail:   true,
			success:  false,
		},
	} {
		tc := tc
		s.Run(tc.desc, func() {
			s.RegisterClockContract(addr.String(), contractAddress)
			s.JailClockContract(contractAddress)
			s.RegisterClockContract(addr2.String(), contractAddressWithAdmin)
			s.JailClockContract(contractAddressWithAdmin)

			// Unjail contract if needed
			if tc.unjail {
				s.UnjailClockContract(addr.String(), contractAddress)
				s.UnjailClockContract(addr2.String(), contractAddressWithAdmin)
			}

			// Set params
			params := types.DefaultParams()
			err := s.app.AppKeepers.ClockKeeper.SetParams(s.ctx, params)
			s.Require().NoError(err)

			// Try to register all contracts
			res, err := s.clockMsgServer.UnjailClockContract(s.ctx, &types.MsgUnjailClockContract{
				SenderAddress:   tc.sender,
				ContractAddress: tc.contract,
			})

			if !tc.success {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(res, &types.MsgUnjailClockContractResponse{})
			}

			// Ensure contract is unregistered
			s.app.AppKeepers.ClockKeeper.RemoveContract(s.ctx, contractAddress)
			s.app.AppKeepers.ClockKeeper.RemoveContract(s.ctx, contractAddressWithAdmin)
		})
	}
}
