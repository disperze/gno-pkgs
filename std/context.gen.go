package std

import (
	"fmt"
	"runtime"
	"strings"
)

type contextI interface {
	GetChainID() string
	GetHeight() int64
	GetTimestamp() int64
	GetOrigCaller() Address
	GetOrigPkgAddr() Address
	GetOrigSend() Coins
	GetOrigSendSpent() *Coins
	GetBanker() Banker
	GetRealmPath() string
	GetCallerAt(n int) Address
	IsOriginCall() bool
}

var df = &defContext{
	chainId:     "dev",
	height:      1,
	timestamp:   1671041345,
	callerAt:    make(map[int]Address),
	callerPkgAt: make(map[string]Address),
}

var ctx contextI = df

type defContext struct {
	chainId       string
	height        int64
	timestamp     int64
	origCaller    Address
	origPkgAddr   Address
	origSend      Coins
	origSendSpent *Coins
	banker        Banker
	realmPath     string
	callerAt      map[int]Address
	callerPkgAt   map[string]Address
	originCall    bool
}

// GetBanker implements contextI
func (d defContext) GetBanker() Banker {
	return d.banker
}

// GetChainID implements contextI
func (d defContext) GetChainID() string {
	return d.chainId
}

// GetHeight implements contextI
func (d defContext) GetHeight() int64 {
	return d.height
}

// GetOrigCaller implements contextI
func (d defContext) GetOrigCaller() Address {
	if d.origCaller == "" {
		panic(shimWarn)
	}
	return d.origCaller
}

// GetOrigPkgAddr implements contextI
func (d defContext) GetOrigPkgAddr() Address {
	if d.origPkgAddr == "" {
		panic(shimWarn)
	}
	return d.origPkgAddr
}

// GetOrigSend implements contextI
func (d defContext) GetOrigSend() Coins {
	return d.origSend
}

// GetOrigSendSpent implements contextI
func (d defContext) GetOrigSendSpent() *Coins {
	return d.origSendSpent
}

// GetTimestamp implements contextI
func (d defContext) GetTimestamp() int64 {
	return d.timestamp
}

// GetRealmPath implements contextI
func (d defContext) GetRealmPath() string {
	if d.realmPath == "" {
		panic(shimWarn)
	}
	return d.realmPath
}

// GetCallerAt implements contextI
func (d defContext) GetCallerAt(n int) Address {
	if n < 0 {
		panic("negative caller index")
	}

	if caller, ok := d.callerAt[n]; ok {
		return caller
	}

	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	origPkg := d.getPkgName(funcName)
	if caller, ok := d.callerPkgAt[fmt.Sprintf("%s-%d", origPkg, n)]; ok {
		return caller
	}

	panic(shimWarn)
}

// IsOriginCall implements contextI
func (d defContext) IsOriginCall() bool {
	return d.originCall
}

func (defContext) getPkgName(funcName string) string {
	lastSlash := strings.LastIndexByte(funcName, '/')
	if lastSlash < 0 {
		lastSlash = 0
	}
	lastDot := strings.LastIndexByte(funcName[lastSlash:], '.') + lastSlash
	return funcName[:lastDot]
}
