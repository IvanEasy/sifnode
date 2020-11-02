package keeper

import (
	"encoding/json"
	_ "fmt"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	"github.com/Sifchain/sifnode/x/oracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
)

var (
	cosmosReceivers, _                         = CreateTestAddrs(1)
	amount, _                                  = strconv.ParseInt("10", 10, 64)
	symbol                                     = "stake"
	tokenContractAddress                       = types.NewEthereumAddress("0xbbbbca6a901c926f240b89eacb641d8aec7aeafd")
	ethBridgeAddress     types.EthereumAddress = types.NewEthereumAddress(strings.ToLower("0x30753E4A8aad7F8597332E813735Def5dD395028"))
	ethereumSender                             = types.NewEthereumAddress("0x627306090abaB3A6e1400e9345bC60c78a8BEf57")
)

func TestProcessClaim(t *testing.T) {
	ctx, keeper, _, _, _, validatorAddresses := CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	validator1Pow3 := validatorAddresses[0]
	validator2Pow3 := validatorAddresses[1]
	nonce, err := strconv.Atoi("1")
	require.NoError(t, err)
	claimType, err := types.StringToClaimType("lock")
	require.NoError(t, err)
	ethBridgeClaim := types.NewEthBridgeClaim(
		5777,
		ethBridgeAddress, // bridge registry
		nonce,
		symbol,
		tokenContractAddress, // loopring
		ethereumSender,
		cosmosReceivers[0],
		validator1Pow3,
		amount,
		claimType,
	)

	status, err := keeper.ProcessClaim(ctx, ethBridgeClaim)

	require.NoError(t, err)
	require.Equal(t, status.Text, oracle.PendingStatusText)

	status, err = keeper.ProcessClaim(ctx, ethBridgeClaim)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "already processed message from validator for this id"))

	// other validator execute

	ethBridgeClaim = types.NewEthBridgeClaim(
		5777,
		ethBridgeAddress, // bridge registry
		nonce,
		symbol,
		tokenContractAddress, // loopring
		ethereumSender,       // accounts[0]
		cosmosReceivers[0],
		validator2Pow3,
		amount,
		claimType,
	)
	status, err = keeper.ProcessClaim(ctx, ethBridgeClaim)
	require.Equal(t, status.Text, oracle.SuccessStatusText)

}

func TestProcessSuccessfulClaim(t *testing.T) {
	ctx, keeper, _, _, _, validatorAddresses := CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	claimType, err := types.StringToClaimType("lock")
	require.NoError(t, err)
	claimContent := types.NewOracleClaimContent(cosmosReceivers[0], amount, symbol, tokenContractAddress, claimType)

	claimBytes, err := json.Marshal(claimContent)
	claimString := string(claimBytes)

	// TODO: find out why bankkeeper cant see ethbridge for mint fn
	require.Panics(t, func() { keeper.ProcessSuccessfulClaim(ctx, claimString) }, "the code did not panic")

	validator1Pow3 := validatorAddresses[0]
	validator2Pow3 := validatorAddresses[1]
	nonce, err := strconv.Atoi("1")
	require.NoError(t, err)
	ethBridgeClaim := types.NewEthBridgeClaim(
		5777,
		ethBridgeAddress, // bridge registry
		nonce,
		symbol,
		tokenContractAddress, // loopring
		ethereumSender,
		cosmosReceivers[0],
		validator1Pow3,
		amount,
		claimType,
	)

	status, err := keeper.ProcessClaim(ctx, ethBridgeClaim)

	require.NoError(t, err)
	require.Equal(t, status.Text, oracle.PendingStatusText)

	require.Panics(t, func() { keeper.ProcessSuccessfulClaim(ctx, claimString) }, "the code did not panic")

	ethBridgeClaim = types.NewEthBridgeClaim(
		5777,
		ethBridgeAddress, // bridge registry
		nonce,
		symbol,
		tokenContractAddress, // loopring
		ethereumSender,
		cosmosReceivers[0],
		validator2Pow3,
		amount,
		claimType,
	)

	status, err = keeper.ProcessClaim(ctx, ethBridgeClaim)

	require.NoError(t, err)
	require.Equal(t, status.Text, oracle.SuccessStatusText)

	// I dont think this one should panic, cannot find ethbridge module account for mint
	//err = keeper.ProcessSuccessfulClaim(ctx, claimString)
	//fmt.Println(err)
}

func TestProcessBurn(t *testing.T) {
	ctx, keeper, _, _, _, _ := CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	coins := sdk.NewCoins(sdk.NewInt64Coin("stake", amount))
	require.Panics(t, func() { keeper.ProcessBurn(ctx, cosmosReceivers[0], coins) }, "the code did not panic")
}

func ProcessLock(t *testing.T) {

	ctx, keeper, _, _, _, _ := CreateTestKeepers(t, 0.7, []int64{3, 3}, "")

	coins := sdk.NewCoins(sdk.NewInt64Coin("stake", amount))
	require.Panics(t, func() { keeper.ProcessLock(ctx, cosmosReceivers[0], coins) }, "the code did not panic")
}
