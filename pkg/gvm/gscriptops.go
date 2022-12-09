package gvm

const (
	SIGHASH_ALL          = 0x1
	SIGHASH_NONE         = 0x2
	SIGHASH_SINGLE       = 0x3
	SIGHASH_ANYONECANPAY = 0x80
)

const ( //enum Register : byte
	R1 = 1
	R2 = 2
	R3 = 3
	R4 = 4
	R5 = 5
	R6 = 6
	R7 = 7
	R8 = 8
)

const ( // enum OperationCode : byte
	OP_0         uint16 = 0 // OP_FALSE
	OP_PUSH      uint16 = 0x4c
	OP_PUSHDATA1 uint16 = 0x4c
	OP_PUSHDATA2 uint16 = 0x4d
	OP_PUSHDATA4 uint16 = 0x4e
	OP_1NEGATE   uint16 = 0x4f
	OP_PUSHSIG   uint16 = 0x50
	OP_PUSHTOKEN uint16 = 0x51
	OP_1         uint16 = 0x52 // OP_TRUE
	OP_POP       uint16 = 0x53

	OP_TRANSFER_TOKEN uint16 = 0x60
	OP_CREATE_TOKEN   uint16 = 0x61
	OP_PAY            uint16 = 0x62
	OP_MAPPING        uint16 = 0x63

	//stack
	OP_TOALTSTACK uint16 = 0x6b
	OP_DUP        uint16 = 0x76
	OP_RETURN     uint16 = 0x77

	//Bitwise
	OP_INVERT      uint16 = 0x83
	OP_AND         uint16 = 0x84
	OP_OR          uint16 = 0x85
	OP_XOR         uint16 = 0x86
	OP_EQUAL       uint16 = 0x87
	OP_EQUALVERIFY uint16 = 0x88

	//cripto
	OP_RIPEMD160           uint16 = 0xA7
	OP_SHA1                uint16 = 0xA8
	OP_SHA256              uint16 = 0xA9
	OP_HASH160             uint16 = 0xAA
	OP_HASH256             uint16 = 0xAB
	OP_CHECKSIG            uint16 = 0xAC
	OP_CHECKSIGVERIFY      uint16 = 0xAD
	OP_CHECKMULTISIG       uint16 = 0xAE
	OP_CHECKMULTISIGVERIFY uint16 = 0xAF
)
