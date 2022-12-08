package std

import "fmt"

// Realm functions can call std.GetBanker(options) to get
// a banker instance. Banker objects cannot be persisted,
// but can be passed onto other functions to be transacted
// on. A banker instance can be passed onto other realm
// functions; this allows other realms to spend coins on
// behalf of the first realm.
//
// Banker panics on errors instead of returning errors.
// This also helps simplify the interface and prevent
// hidden bugs (e.g. ignoring errors)
//
// NOTE: this Gno interface is satisfied by a native go
// type, and those can't return non-primitive objects
// (without confusion).
type Banker interface {
	GetCoins(addr Address) (dst Coins)
	SendCoins(from, to Address, amt Coins)
	TotalCoin(denom string) int64
	IssueCoin(addr Address, denom string, amount int64)
	RemoveCoin(addr Address, denom string, amount int64)
}

// Also available natively in stdlibs/context.go
type BankerType uint8

// Also available natively in stdlibs/context.go
const (
	// Can only read state.
	BankerTypeReadonly BankerType = iota
	// Can only send from tx send.
	BankerTypeOrigSend
	// Can send from all realm coins.
	BankerTypeRealmSend
	// Can issue and remove realm coins.
	BankerTypeRealmIssue
)

//----------------------------------------
// ReadonlyBanker

type ReadonlyBanker struct {
	banker Banker
}

func NewReadonlyBanker(banker Banker) ReadonlyBanker {
	return ReadonlyBanker{banker}
}

func (rb ReadonlyBanker) GetCoins(addr Address) (dst Coins) {
	return rb.banker.GetCoins(addr)
}

func (rb ReadonlyBanker) SendCoins(from, to Address, amt Coins) {
	panic("ReadonlyBanker cannot send coins")
}

func (rb ReadonlyBanker) TotalCoin(denom string) int64 {
	return rb.banker.TotalCoin(denom)
}

func (rb ReadonlyBanker) IssueCoin(addr Address, denom string, amount int64) {
	panic("ReadonlyBanker cannot issue coins")
}

func (rb ReadonlyBanker) RemoveCoin(addr Address, denom string, amount int64) {
	panic("ReadonlyBanker cannot remove coins")
}

//----------------------------------------
// OrigSendBanker

type OrigSendBanker struct {
	banker        Banker
	pkgAddr       Address
	origSend      Coins
	origSendSpent *Coins
}

func NewOrigSendBanker(banker Banker, pkgAddr Address, origSend Coins, origSendSpent *Coins) OrigSendBanker {
	if origSendSpent == nil {
		panic("origSendSpent cannot be nil")
	}
	return OrigSendBanker{
		banker:        banker,
		pkgAddr:       pkgAddr,
		origSend:      origSend,
		origSendSpent: origSendSpent,
	}
}

func (osb OrigSendBanker) GetCoins(addr Address) (dst Coins) {
	return osb.banker.GetCoins(addr)
}

func (osb OrigSendBanker) SendCoins(from, to Address, amt Coins) {
	if from != osb.pkgAddr {
		panic(fmt.Sprintf(
			"OrigSendBanker can only send from the realm package address %q, but got %q",
			osb.pkgAddr, from))
	}
	spent := (*osb.origSendSpent).Add(amt)
	if !osb.isAllGTE(osb.origSend, spent) {
		panic(fmt.Sprintf(
			`cannot send "%v", limit "%v" exceeded with "%v" already spent`,
			amt, osb.origSend, *osb.origSendSpent))
	}
	osb.banker.SendCoins(from, to, amt)
	*osb.origSendSpent = spent
}

func (osb OrigSendBanker) TotalCoin(denom string) int64 {
	return osb.banker.TotalCoin(denom)
}

func (osb OrigSendBanker) IssueCoin(addr Address, denom string, amount int64) {
	panic("OrigSendBanker cannot issue coins")
}

func (osb OrigSendBanker) RemoveCoin(addr Address, denom string, amount int64) {
	panic("OrigSendBanker cannot remove coins")
}

func (osb OrigSendBanker) isAllGTE(coins, coinsB Coins) bool {
	if len(coinsB) == 0 {
		return true
	}

	if len(coins) == 0 {
		return false
	}

	for _, coinB := range coinsB {
		if coinB.Amount > coins.AmountOf(coinB.Denom) {
			return false
		}
	}

	return true
}

//----------------------------------------
// RealmSendBanker

type RealmSendBanker struct {
	banker  Banker
	pkgAddr Address
}

func NewRealmSendBanker(banker Banker, pkgAddr Address) RealmSendBanker {
	return RealmSendBanker{
		banker:  banker,
		pkgAddr: pkgAddr,
	}
}

func (rsb RealmSendBanker) GetCoins(addr Address) (dst Coins) {
	return rsb.banker.GetCoins(addr)
}

func (rsb RealmSendBanker) SendCoins(from, to Address, amt Coins) {
	if from != rsb.pkgAddr {
		panic(fmt.Sprintf(
			"RealmSendBanker can only send from the realm package address %q, but got %q",
			rsb.pkgAddr, from))
	}
	rsb.banker.SendCoins(from, to, amt)
}

func (rsb RealmSendBanker) TotalCoin(denom string) int64 {
	return rsb.banker.TotalCoin(denom)
}

func (rsb RealmSendBanker) IssueCoin(addr Address, denom string, amount int64) {
	panic("RealmSendBanker cannot issue coins")
}

func (rsb RealmSendBanker) RemoveCoin(addr Address, denom string, amount int64) {
	panic("RealmSendBanker cannot remove coins")
}
