package config

//go:generate go run github.com/NilFoundation/fastssz/sszgen --path params.go -include ../types/address.go,../types/uint256.go,../types/transaction.go --objs ListValidators,ParamValidators,ValidatorInfo,ParamGasPrice,ParamFees,WorkaroundToImportTypes
