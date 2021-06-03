package keeper

import (
	"fmt"

	crud "github.com/iov-one/cosmos-sdk-crud"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/iov-one/starnamed/x/escrow/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keeper defines the escrow keeper
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.Marshaler
	paramSpace    paramstypes.Subspace
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	//TODO : store it into persistent storage !!
	lastBlockTime uint64
	storeHolders  map[types.TypeID]types.StoreHolder
	blockedAddrs  map[string]bool
}

// NewKeeper creates a new escrow Keeper instance
func NewKeeper(
	cdc codec.Marshaler,
	key sdk.StoreKey,
	paramSpace paramstypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	blockedAddrs map[string]bool,
	storeHolders map[types.TypeID]types.StoreHolder,
) Keeper {
	// ensure the escrow module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramSpace:    paramSpace,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		storeHolders:  storeHolders,
		blockedAddrs:  blockedAddrs,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("irismod/%s", types.ModuleName))
}

// GetEscrowAccount returns the escrow module account
func (k Keeper) GetEscrowAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

func (k Keeper) FetchNextId() tmbytes.HexBytes {
	//TODO: not implemented
	return nil
}

func (k Keeper) NextId() {
	//TODO: not implemented
}

func (k Keeper) getStoreForID(id types.TypeID) (crud.Store, error) {
	store, ok := k.storeHolders[id]
	if !ok {
		return nil, types.ErrUnknownTypeID
	}
	return store.GetCRUDStore(), nil
}

func (k Keeper) SetLastBlockTime(date uint64) {
	k.lastBlockTime = date
}

func (k Keeper) GetLastBlockTime() uint64 {
	return k.lastBlockTime
}

func (k Keeper) isBlockedAddr(address string) bool {
	return k.blockedAddrs[address]
}
