package accounts

import (
	"errors"
	"fmt"
)

// Account struct
type Account struct {
	owner string
	balance int
}

// NewAccount creates Account
func NewAccount(owner string) *Account{
	// 이름을 인자로 받아 구조체를 생성 
	account:= Account{owner:owner, balance: 0}
	// 주소값으로 전달
	return &account
}

// Deposit x amount on your account
func (a *Account) Deposit(amount int) {
	// 구조체 methods 등록하는 방법!!!
	// (a *Account) reciver 라고 함
	// 컨벤션은 struct에 첫 글자를 따서 소문자로 지어야 함

	a.balance += amount
}

// Balance of your account
func (a *Account) Balance() int {
	return a.balance
}

// errNoMoney not enough money
var errNoMoney = errors.New("can't withdraw")

// Withdraw amount from your account
func (a *Account) WithDraw(amount int) error {
	if a.balance < amount {
		return errNoMoney
	}
	a.balance -= amount
	// nil == null , None 과 같음
	return nil
}

// ChangeOwner of the account
func (a *Account) ChangeOwner(newOwner string){
	a.owner = newOwner
} 


// Owner of the account
func (a *Account) Owner() string {
	return a.owner
}

// String like __str__
func(a *Account) String() string {
	// return "whatever you want"

	return fmt.Sprint(a.Owner(), "'s account.\nHas: ",a.Balance())
}