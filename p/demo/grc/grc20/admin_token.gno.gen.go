// Code generated by github.com/gnolang/gno. DO NOT EDIT.




package grc20

import (
	"github.com/gnolang/gno/stdlibs/stdshim"

	"gno.tools/p/demo/avl"
	"gno.tools/p/demo/ufmt"
)

// AdminToken implements a token factory with admin helpers.
//
// Warning: you should not expose this struct to enduser directly.
//
// It allows token administrators to call privileged helpers
// like Mint, Burn, or any Transfer helpers by passing custom owners.
//
// You should initialize your token, then call AdminToken.SafeGRC20() to
// expose a safe instance to the endusers.
type AdminToken struct {
	name        string
	symbol      string
	decimals    uint
	totalSupply uint64
	balances    *avl.MutTree // std.Address(owner) -> uint64
	allowances  *avl.MutTree // string(owner+":"+spender) -> uint64
}

// safeToken implements the IGRC20 interface.
//
// It is generated by AdminToken.SafeGRC20().
// It can safely be explosed publicly.
type safeToken struct {
	IGRC20 // implements the GRC20 interface.

	factory *AdminToken
}

func NewAdminToken(name, symbol string, decimals uint) *AdminToken {
	// FIXME: check for limits

	return &AdminToken{
		name:     name,
		symbol:   symbol,
		decimals: decimals,

		balances:   avl.NewMutTree(),
		allowances: avl.NewMutTree(),
	}
}

func (t *AdminToken) GetName() string     { return t.name }
func (t *AdminToken) GetSymbol() string   { return t.symbol }
func (t *AdminToken) GetDecimals() uint   { return t.decimals }
func (t *AdminToken) TotalSupply() uint64 { return t.totalSupply }

func (t *AdminToken) BalanceOf(owner std.Address) (uint64, error) {
	return t.balanceOf(owner)
}

func (t *AdminToken) Transfer(owner, to std.Address, amount uint64) error {
	return t.transfer(owner, to, amount)
}

func (t *AdminToken) Allowance(owner, spender std.Address) (uint64, error) {
	return t.allowance(owner, spender)
}

func (t *AdminToken) Approve(owner, spender std.Address, amount uint64) error {
	return t.approve(owner, spender, amount)
}

func (t *AdminToken) TransferFrom(spender, from, to std.Address, amount uint64) error {
	if err := t.spendAllowance(from, spender, amount); err != nil {
		return err
	}
	return t.transfer(from, to, amount)
}

// Administration helpers implementation.
//

func (t *AdminToken) Mint(to std.Address, amount uint64) error {
	return t.mint(to, amount)
}

func (t *AdminToken) Burn(from std.Address, amount uint64) error {
	return t.burn(from, amount)
}

// private helpers
//

func (t *AdminToken) mint(address std.Address, amount uint64) error {
	if err := checkIsValidAddress(address); err != nil {
		return err
	}

	// TODO: check for overflow

	t.totalSupply += amount
	currentBalance, err := t.balanceOf(address)
	if err != nil {
		return err
	}
	newBalance := currentBalance + amount

	t.balances.Set(string(address), newBalance)

	event := TransferEvent{zeroAddress, address, amount}
	emit(&event)

	return nil
}

func (t *AdminToken) burn(address std.Address, amount uint64) error {
	if err := checkIsValidAddress(address); err != nil {
		return err
	}
	// TODO: check for overflow

	currentBalance, err := t.balanceOf(address)
	if err != nil {
		return err
	}
	if currentBalance < amount {
		return ErrInsufficientBalance
	}

	t.totalSupply -= amount
	newBalance := currentBalance - amount

	t.balances.Set(string(address), newBalance)

	event := TransferEvent{address, zeroAddress, amount}
	emit(&event)

	return nil
}

func (t *AdminToken) balanceOf(address std.Address) (uint64, error) {
	if err := checkIsValidAddress(address); err != nil {
		return 0, err
	}

	balance, found := t.balances.Get(address.String())
	if !found {
		return 0, nil
	}
	return balance.(uint64), nil
}

func (t *AdminToken) spendAllowance(owner, spender std.Address, amount uint64) error {
	if err := checkIsValidAddress(owner); err != nil {
		return err
	}
	if err := checkIsValidAddress(spender); err != nil {
		return err
	}

	currentAllowance, err := t.allowance(owner, spender)
	if err != nil {
		return err
	}
	if currentAllowance < amount {
		return ErrInsufficientAllowance
	}

	return nil
}

func (t *AdminToken) transfer(from, to std.Address, amount uint64) error {
	if err := checkIsValidAddress(from); err != nil {
		return err
	}
	if err := checkIsValidAddress(to); err != nil {
		return err
	}

	if from == to {
		return ErrCannotTransferToSelf
	}

	toBalance, err := t.balanceOf(to)
	if err != nil {
		return err
	}
	fromBalance, err := t.balanceOf(from)
	if err != nil {
		return err
	}

	// debug.
	// println("from", from, "to", to, "amount", amount, "fromBalance", fromBalance, "toBalance", toBalance)

	if fromBalance < amount {
		return ErrInsufficientBalance
	}

	newToBalance := toBalance + amount
	newFromBalance := fromBalance - amount

	t.balances.Set(string(to), newToBalance)
	t.balances.Set(string(from), newFromBalance)

	event := TransferEvent{from, to, amount}
	emit(&event)

	return nil
}

func (t *AdminToken) allowance(owner, spender std.Address) (uint64, error) {
	if err := checkIsValidAddress(owner); err != nil {
		return 0, err
	}
	if err := checkIsValidAddress(spender); err != nil {
		return 0, err
	}

	key := owner.String() + ":" + spender.String()

	allowance, found := t.allowances.Get(key)
	if !found {
		return 0, nil
	}

	return allowance.(uint64), nil
}

func (t *AdminToken) approve(owner, spender std.Address, amount uint64) error {
	if err := checkIsValidAddress(owner); err != nil {
		return err
	}
	if err := checkIsValidAddress(spender); err != nil {
		return err
	}

	key := owner.String() + ":" + spender.String()
	t.allowances.Set(key, amount)

	event := ApprovalEvent{owner, spender, amount}
	emit(&event)

	return nil
}

func (t *AdminToken) RenderHome() string {
	str := ""
	str += ufmt.Sprintf("# %s ($%s)\n\n", t.name, t.symbol)
	str += ufmt.Sprintf("* **Decimals**: %d\n", t.decimals)
	str += ufmt.Sprintf("* **Total supply**: %d\n", t.totalSupply)
	str += ufmt.Sprintf("* **Known accounts**: %d\n", t.balances.Size())
	return str
}

// GRC20 returns an instance that can be exposed to the end user.
func (t *AdminToken) GRC20() IGRC20 {
	return &userToken{admin: t}
}
