package dexT1

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/KyberNetwork/ethrpc"
	"github.com/KyberNetwork/kyberswap-dex-lib/pkg/entity"
	"github.com/KyberNetwork/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestPoolListUpdater(t *testing.T) {
	_ = logger.SetLogLevel("debug")

	if os.Getenv("CI") != "" {
		t.Skip()
	}

	var (
		pools            []entity.Pool
		metadataBytes, _ = json.Marshal(map[string]interface{}{})
		err              error

		config = Config{
			DexReservesResolver: "0x90bFebd5Ac2d6787028DC8A58d0b5EE07b16E06F",
		}
	)

	// Setup RPC server
	rpcClient := ethrpc.New("https://ethereum.kyberengineering.io")
	rpcClient.SetMulticallContract(common.HexToAddress("0x5ba1e12693dc8f9c48aad8770482f4739beed696"))

	pu := NewPoolsListUpdater(&config, rpcClient)
	require.NotNil(t, pu)

	pools, _, err = pu.GetNewPools(context.Background(), metadataBytes)
	require.NoError(t, err)
	require.True(t, len(pools) >= 1)

	staticExtraBytes, _ := json.Marshal(&StaticExtra{
		DexReservesResolver: config.DexReservesResolver,
	})

	expectedPool0 := entity.Pool{
		Address:  "0x25F0A3B25cBC0Ca0417770f686209628323fF901",
		Exchange: "fluid-dex-t1",
		Type:     "fluid-dex-t1",
		Reserves: pools[0].Reserves,
		Tokens: []*entity.PoolToken{
			{
				Address:   "0x7f39C581F595B53c5cb19bD0b3f8dA6c935E2Ca0",
				Weight:    1,
				Swappable: true,
				Decimals:  18,
			},
			{
				Address:   "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
				Weight:    1,
				Swappable: true,
				Decimals:  18,
			},
		},
		SwapFee: 0.01,

		Extra:       pools[0].Extra,
		StaticExtra: string(staticExtraBytes),
	}

	require.Equal(t, expectedPool0, pools[0])

	var extra PoolExtra
	err = json.Unmarshal([]byte(pools[0].Extra), &extra)
	require.NoError(t, err)

	require.NotEqual(t, "0", pools[0].Reserves[0], "Reserve should not be zero")
	require.NotEqual(t, "0", pools[0].Reserves[1], "Reserve should not be zero")

	require.True(t, extra.CollateralReserves.Token0RealReserves.Cmp(big.NewInt(0)) > 0)
	require.True(t, extra.CollateralReserves.Token1RealReserves.Cmp(big.NewInt(0)) > 0)
	require.True(t, extra.CollateralReserves.Token0ImaginaryReserves.Cmp(big.NewInt(0)) > 0)
	require.True(t, extra.CollateralReserves.Token1ImaginaryReserves.Cmp(big.NewInt(0)) > 0)
	require.True(t, extra.DebtReserves.Token0Debt.Cmp(big.NewInt(0)) > 0)
	require.True(t, extra.DebtReserves.Token1Debt.Cmp(big.NewInt(0)) > 0)
	require.True(t, extra.DebtReserves.Token0RealReserves.Cmp(big.NewInt(0)) > 0)
	require.True(t, extra.DebtReserves.Token1RealReserves.Cmp(big.NewInt(0)) > 0)
	require.True(t, extra.DebtReserves.Token0ImaginaryReserves.Cmp(big.NewInt(0)) > 0)
	require.True(t, extra.DebtReserves.Token1ImaginaryReserves.Cmp(big.NewInt(0)) > 0)

	// Log all pools
	// for i, pool := range pools {
	// 	jsonEncoded, _ := json.MarshalIndent(pool, "", "  ")
	// 	t.Logf("Pool %d: %s\n", i, string(jsonEncoded))
	// }

}
