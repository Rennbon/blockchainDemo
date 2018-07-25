package errors

import "errors"

var (
	ERR_OPEN_SOURCE          = errors.New("Open source error.")
	ERR_DATA_EXISTS          = errors.New("Data already exists.")
	ERR_DATA_INCONSISTENCIES = errors.New("Data inconsistency.")
	ERR_NOPE                 = errors.New("Nope.")
	ERR_NOT_ENOUGH_COIN      = errors.New("Not enough coin.")
	ERR_UNCONFIRMED          = errors.New("Tx unconfirmed.")
	ERR_Param_Fail           = errors.New("Parameter validation failed.")
)
