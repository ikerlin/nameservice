package keeper

import (
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ikerlin/nameservice/x/nameservice/types"
	"strconv"
)

// GetWhoisCount get the total number of whois
func (k Keeper) GetWhoisCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WhoisCountKey))
	byteKey := types.KeyPrefix(types.WhoisCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to iint64
		panic("cannot decode count")
	}

	return count
}

// SetWhoisCount set the total number of whois
func (k Keeper) SetWhoisCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WhoisCountKey))
	byteKey := types.KeyPrefix(types.WhoisCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendWhois appends a whois in the store with a new id and update the count
func (k Keeper) AppendWhois(
	ctx sdk.Context,
	creator string,
	value string,
	price string,
) uint64 {
	// Create the whois
	count := k.GetWhoisCount(ctx)
	var whois = types.Whois{
		Creator: creator,
		Id:      count,
		Value:   value,
		Price:   price,
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WhoisKey))
	store.Set(GetWhoisIDBytes(whois.Id), k.cdc.MustMarshalBinaryBare(&whois))

	// Update whois count
	k.SetWhoisCount(ctx, count+1)

	return count
}

// SetWhois set a specific whois in the store
func (k Keeper) SetWhois(ctx sdk.Context, name string, whois types.Whois) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WhoisKey))
	b := k.cdc.MustMarshalBinaryBare(&whois)
	store.Set([]byte(name), b)
}

// GetWhois returns a whois from its id
func (k Keeper) GetWhois(ctx sdk.Context, name string) types.Whois {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WhoisKey))
	var whois types.Whois
	k.cdc.MustUnmarshalBinaryBare(store.Get([]byte(name)), &whois)
	return whois
}

// HasWhois checks if the whois exists in the store
func (k Keeper) HasWhois(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WhoisKey))
	return store.Has(GetWhoisIDBytes(id))
}

// GetWhoisOwner returns the creator of the whois
func (k Keeper) GetWhoisOwner(ctx sdk.Context, name string) string {
	return k.GetWhois(ctx, name).Creator
}

// RemoveWhois removes a whois from the store
func (k Keeper) RemoveWhois(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WhoisKey))
	store.Delete(GetWhoisIDBytes(id))
}

// GetAllWhois returns all whois
func (k Keeper) GetAllWhois(ctx sdk.Context) (list []types.Whois) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WhoisKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Whois
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// SetNameValue set name's value
func (k Keeper) SetNameValue(ctx sdk.Context, name string, value string) {
	whois := k.GetWhois(ctx, name)
	whois.Value = value
	k.SetWhois(ctx, name, whois)
}

// GetWhoisIDBytes returns the byte representation of the ID
func GetWhoisIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetWhoisIDFromBytes returns ID in uint64 format from a byte array
func GetWhoisIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
