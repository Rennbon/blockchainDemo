package errors

import "errors"

var (
	ERR_OPEN_SOURCE          = errors.New("open source error.")
	ERR_DATA_EXISTS          = errors.New("data already exists.")
	ERR_DATA_INCONSISTENCIES = errors.New("data inconsistency.")
	ERR_NOPE                 = errors.New("nope")
	ERR_NOT_ENOUGH_COIN      = errors.New("not enougt coin")
)
