//go:build test
// +build test

package std

func TestSetContext(c contextI) {
	ctx = c
}

func TestSetCallerAt(n int, addr Address) {
	df.callerAt[n] = addr
}

func TestSetChainID(id string) {
	df.chainId = id
}

func TestSetHeight(h int64) {
	df.height = h
}

func TestSetTimestamp(t int64) {
	df.timestamp = t
}

func TestSetOrigCaller(addr Address) {
	df.origCaller = addr
}

func TestSetOrigPkgAddr(addr Address) {
	df.origPkgAddr = addr
}

func TestSetOrigSend(c Coins) {
	df.origSend = c
}

func TestSetOrigSendSpent(c *Coins) {
	df.origSendSpent = c
}

func TestSetBanker(b Banker) {
	df.banker = b
}

func TestSetRealmPath(p string) {
	df.realmPath = p
}

func TestSetOriginCall(b bool) {
	df.originCall = b
}
