package std

// defContext is not exposed,
// use native injections std.GetChainID(), std.GetHeight() etc instead.

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

var ctx contextI = defContext{}

type defContext struct {
}

// GetBanker implements contextI
func (defContext) GetBanker() Banker {
	panic(shimWarn)
}

// GetChainID implements contextI
func (defContext) GetChainID() string {
	panic(shimWarn)
}

// GetHeight implements contextI
func (defContext) GetHeight() int64 {
	panic(shimWarn)
}

// GetOrigCaller implements contextI
func (defContext) GetOrigCaller() Address {
	panic(shimWarn)
}

// GetOrigPkgAddr implements contextI
func (defContext) GetOrigPkgAddr() Address {
	panic(shimWarn)
}

// GetOrigSend implements contextI
func (defContext) GetOrigSend() Coins {
	panic(shimWarn)
}

// GetOrigSendSpent implements contextI
func (defContext) GetOrigSendSpent() *Coins {
	panic(shimWarn)
}

// GetTimestamp implements contextI
func (defContext) GetTimestamp() int64 {
	panic(shimWarn)
}

// GetRealmPath implements contextI
func (defContext) GetRealmPath() string {
	panic(shimWarn)
}

// GetCallerAt implements contextI
func (defContext) GetCallerAt(n int) Address {
	panic(shimWarn)
}

// IsOriginCall implements contextI
func (defContext) IsOriginCall() bool {
	panic(shimWarn)
}
