//go:build gno
// +build gno

package std

import (
	"time"

	"github.com/gnolang/gno/pkgs/bech32"
	gno "github.com/gnolang/gno/pkgs/gnolang"
)

const shimWarn = "stdshim cannot be used to run code"

func AssertOriginCall() {
	if !ctx.IsOriginCall() {
		panic("invalid non-origin call")
	}
}

func IsOriginCall() (isOrigin bool) {
	return ctx.IsOriginCall()
}

func Hash(bz []byte) (hash [20]byte) {
	return gno.HashBytes(bz)
}

func CurrentRealmPath() string {
	return ctx.GetRealmPath()
}

func GetChainID() string {
	return ctx.GetChainID()
}

func GetHeight() int64 {
	return ctx.GetHeight()
}

func GetOrigSend() Coins {
	return ctx.GetOrigSend()
}

func GetOrigCaller() Address {
	return ctx.GetOrigCaller()
}

func GetOrigPkgAddr() Address {
	return ctx.GetOrigPkgAddr()
}

func GetCallerAt(n int) Address {
	return ctx.GetCallerAt(n)
}

func GetBanker(bankerType BankerType) Banker {
	banker := ctx.GetBanker()
	switch bankerType {
	case BankerTypeReadonly:
		banker = NewReadonlyBanker(banker)
	case BankerTypeOrigSend:
		banker = NewOrigSendBanker(banker, ctx.GetOrigPkgAddr(), ctx.GetOrigSend(), ctx.GetOrigSendSpent())
	case BankerTypeRealmSend:
		banker = NewRealmSendBanker(banker, ctx.GetOrigPkgAddr())
	case BankerTypeRealmIssue:
		banker = banker
	default:
		panic("should not happen") // defensive
	}

	return banker
}

func GetTimestamp() Time {
	return Time(ctx.GetTimestamp())
}

func FormatTimestamp(timestamp Time, format string) string {
	t := time.Unix(int64(timestamp), 0).Round(0).UTC()
	return t.Format(format)
}

func EncodeBech32(prefix string, bytes [20]byte) (addr Address) {
	b32, err := bech32.ConvertAndEncode(prefix, bytes[:])
	if err != nil {
		panic(err)
	}
	return Address(b32)
}

func DecodeBech32(addr Address) (prefix string, bytes [20]byte, ok bool) {
	prefix, bz, err := bech32.Decode(addr.String())
	if err != nil {
		return "", bytes, false
	}

	var b [20]byte
	copy(b[:], bz)
	return prefix, b, true
}

func DerivePkgAddr(pkgPath string) (addr Address) {
	b32 := gno.DerivePkgAddr(pkgPath).Bech32()

	return Address(b32)
}
