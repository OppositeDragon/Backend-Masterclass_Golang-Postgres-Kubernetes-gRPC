package api

import (
	"fmt"
	"reflect"
	db "simple_bank/db/sqlc"
	"simple_bank/util"

	"go.uber.org/mock/gomock"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, isOk := x.(db.CreateUserParams)
	if !isOk {
		return false
	}
	err := util.CheckPasswordHash(e.password, arg.HashedPassword)
	if err != nil {
		return false

	}
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}
