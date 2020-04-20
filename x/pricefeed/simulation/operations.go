package simulation

import (
	"math/rand"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/kava-labs/kava/app/helpers"
	appparams "github.com/kava-labs/kava/app/params"
	"github.com/kava-labs/kava/x/pricefeed/keeper"
	"github.com/kava-labs/kava/x/pricefeed/types"
)

var (
	noOpMsg   = simulation.NoOpMsg(types.ModuleName)
	btcPrices = []sdk.Dec{}
	bnbPrices = []sdk.Dec{}
	xrpPrices = []sdk.Dec{}
	genPrices sync.Once
)

// Simulation operation weights constants
const (
	OpWeightMsgUpdatePrices = "op_weight_msg_update_prices"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simulation.AppParams, cdc *codec.Codec, ak auth.AccountKeeper,
	k keeper.Keeper, wContents []simulation.WeightedProposalContent,
) simulation.WeightedOperations {
	var weightMsgUpdatePrices int

	appParams.GetOrGenerate(cdc, OpWeightMsgUpdatePrices, &weightMsgUpdatePrices, nil,
		func(_ *rand.Rand) {
			weightMsgUpdatePrices = appparams.DefaultWeightMsgUpdatePrices
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgUpdatePrices,
			SimulateMsgUpdatePrices(ak, k, blocks),
		),
	}
}

// SimulateMsgUpdatePrices updates the prices of various assets by randomly varying them based on current price
func SimulateMsgUpdatePrices(ak auth.AccountKeeper, keeper.Keeper, blocks int) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		genPrices.Do(func() {
			// generate a random walk for each asset exactly once, with observations equal to the number of blocks in the sim
			for _, m := range keeper.GetMarkets(ctx) {
				startPrice := getStartPrice(m.MarketID)
				// allow prices to fluctuate from 10x GAINZ to 100x REKT
				maxPrice := sdk.MustNewDecFromStr("10.0").Mul(startPrice)
				minPrice := sdk.MustNewDecFromStr("0.01").Mul(startPrice)
				previousPrice := startPrice
				for i := 0; i < blocks; i++ {
					increment := getIncrement(m.MarketID)
					// note calling r instead of rand here breaks determinism
					upDown := rand.Intn(2)
					if upDown == 0 {
						if previousPrice.Add(increment).GT(maxPrice) {
							previousPrice = maxPrice
						} else {
							previousPrice = previousPrice.Add(increment)
						}
					} else {
						if previousPrice.Sub(increment).LT(minPrice) {
							previousPrice = minPrice
						} else {
							previousPrice = previousPrice.Sub(increment)
						}
					}
					setPrice(m.MarketID, previousPrice)
				}
			}
		})

		randomMarket := pickRandomAsset(ctx, keeper, r)
		marketID := randomMarket.MarketID
		address := getRandomOracle(r, randomMarket)
		acc := ak.GetAccount(ctx, address)
		if acc == nil {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		price := pickNewRandomPrice(marketID, int(ctx.BlockHeight()))

		// get the expiry time based off the current time
		expiry := getExpiryTime(ctx)

		// now create the msg to post price
		msg := types.NewMsgPostPrice(address, marketID, price, expiry)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{acc.GetAccountNumber()},
			[]uint64{acc.GetSequence()},
			bidder.PrivKey,
		)

		_, result, err := app.Deliver(tx)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}
		return simulation.NewOperationMsg(msg, true, result.Log), nil, nil
	}
}

func getStartPrice(marketID string) (startPrice sdk.Dec) {
	switch marketID {
	case "btc:usd":
		return sdk.MustNewDecFromStr("7000")
	case "bnb:usd":
		return sdk.MustNewDecFromStr("15")
	case "xrp:usd":
		return sdk.MustNewDecFromStr("0.25")
	}
	return sdk.MustNewDecFromStr("100")
}

func getIncrement(marketID string) (increment sdk.Dec) {
	startPrice := getStartPrice(marketID)
	divisor := sdk.MustNewDecFromStr("20")
	increment = startPrice.Quo(divisor)
	return increment
}

func setPrice(marketID string, price sdk.Dec) {
	switch marketID {
	case "btc:usd":
		btcPrices = append(btcPrices, price)
		return
	case "bnb:usd":
		bnbPrices = append(bnbPrices, price)
		return
	case "xrp:usd":
		xrpPrices = append(xrpPrices, price)
	}
	return
}

func pickNewRandomPrice(marketID string, blockHeight int) (newPrice sdk.Dec) {
	switch marketID {
	case "btc:usd":
		return btcPrices[blockHeight-1]
	case "bnb:usd":
		return bnbPrices[blockHeight-1]
	case "xrp:usd":
		return xrpPrices[blockHeight-1]
	}
	panic("invalid price request")
}

// getRandomOracle picks a random oracle from the list of oracles
func getRandomOracle(r *rand.Rand, market types.Market) sdk.AccAddress {
	randomIndex := simulation.RandIntBetween(r, 0, len(market.Oracles))
	return market.Oracles[randomIndex]
}

// pickRandomAsset picks a random asset out of the assets with equal probability
// it returns the Market which includes the base asset as one of its fields
func pickRandomAsset(ctx sdk.Context, keeper keeper.Keeper, r *rand.Rand) (market types.Market) {
	// get the params
	params := keeper.GetParams(ctx)
	// now pick a random asset
	randomIndex := simulation.RandIntBetween(r, 0, len(params.Markets))
	return params.Markets[randomIndex]
}

// getExpiryTime gets a price expiry time by taking the current time and adding a delta to it
func getExpiryTime(ctx sdk.Context) (t time.Time) {
	// need to use the blocktime from the context as the context generates random start time when running simulations
	return ctx.BlockTime().Add(time.Second * 1000000)
}