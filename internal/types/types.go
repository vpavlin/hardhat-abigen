package types

type ABI struct {
	ContractName string      `json:"contractName"`
	ABI          interface{} `json:"abi"`
	ByteCode     string      `json:"bytecode"`
}
