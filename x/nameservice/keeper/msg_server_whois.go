package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ikerlin/nameservice/x/nameservice/types"
)

func (k msgServer) SetName(goCtx context.Context, msg *types.MsgSetName) (*types.MsgSetNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.HasWhois(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(types.ErrNameDoesNotExist, fmt.Sprintf("Name %s doesn't exist", msg.Name))
	}

	if msg.Creator != k.GetWhoisOwner(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}
	k.SetNameValue(ctx, msg.Name, msg.Value)
	return &types.MsgSetNameResponse{}, nil
}

func (k msgServer) BuyName(goCtx context.Context, msg *types.MsgBuyName) (*types.MsgBuyNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var whois types.Whois
	if k.HasWhois(ctx, msg.Name) {
		whois = k.GetWhois(ctx, msg.Name)
		price, _ := sdk.ParseCoinsNormalized(whois.Price)
		bidPrice, _ := sdk.ParseCoinsNormalized(msg.Bid)

		if price.IsAllGT(bidPrice) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid not high enough") // If not, throw an error
		}

		creator, err := sdk.AccAddressFromBech32(whois.Creator)
		if err != nil {
			panic(err)
		}
		buyer, err := sdk.AccAddressFromBech32(msg.Buyer)
		if err != nil {
			return nil, err
		}

		err = k.bankKeeper.SendCoins(ctx, buyer, creator, bidPrice)
		if err != nil {
			return nil, err
		}

		whois.Creator = msg.Buyer
		whois.Price = msg.Bid

	} else {
		bidPrice, _ := sdk.ParseCoinsNormalized(msg.Bid)
		buyer, err := sdk.AccAddressFromBech32(msg.Buyer)
		if err != nil {
			return nil, err
		}

		if k.MinPrice(ctx).IsAllGT(bidPrice) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid not high enough")
		}

		err = k.bankKeeper.SubtractCoins(ctx, buyer, bidPrice)
		if err != nil {
			return nil, err
		}

		whois = types.Whois{
			Creator: msg.Buyer,
			Price:   msg.Bid,
		}
	}

	k.SetWhois(ctx, msg.Name, whois)

	return &types.MsgBuyNameResponse{}, nil
}

func (k msgServer) DeleteName(goCtx context.Context, msg *types.MsgDeleteName) (*types.MsgDeleteNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgDeleteNameResponse{}, nil
}
