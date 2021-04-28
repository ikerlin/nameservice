package keeper

import (
	"github.com/ikerlin/nameservice/x/nameservice/types"
)

var _ types.QueryServer = Keeper{}
