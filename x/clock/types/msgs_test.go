package types

import (
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgsTestSuite struct {
	suite.Suite
	govModule string
}

func TestMsgsTestSuite(t *testing.T) {
	suite.Run(t, new(MsgsTestSuite))
}

func (suite *MsgsTestSuite) SetupTest() {
	suite.govModule = "juno10d07y265gmmuvt4z0w9aw880jnsr700jvss730"
}

func (suite *MsgsTestSuite) TestMsgUpdateParams() {
	var limit uint64 = 100_000

	p := MsgUpdateParams{
		Authority: suite.govModule,
		Params: Params{
			ContractGasLimit: limit,
		},
	}

	acc, _ := sdk.AccAddressFromBech32(p.Authority)

	msg := NewMsgUpdateParams(acc, limit)

	suite.Require().Equal(RouterKey, msg.Route())
	suite.Require().Equal(TypeMsgUpdateParams, msg.Type())
	suite.Require().NotNil(msg.GetSigners())
}
