package gvm

import (
	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
)

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
	OP_0         = 0 // OP_FALSE
	OP_PUSH      = 0x4c
	OP_PUSHDATA1 = 0x4c
	OP_PUSHDATA2 = 0x4d
	OP_PUSHDATA4 = 0x4e
	OP_1NEGATE   = 0x4f
	OP_PUSHSIG   = 0x50
	OP_PUSHTOKEN = 0x51
	OP_1         = 0x52 // OP_TRUE
	OP_POP       = 0x53

	OP_TRANSFER_TOKEN = 0x60
	OP_CREATE_TOKEN   = 0x61
	OP_PAY            = 0x62
	OP_MAPPING        = 0x63

	//stack
	OP_TOALTSTACK = 0x6b
	OP_DUP        = 0x76
	OP_RETURN     = 0x77

	//Bitwise
	OP_INVERT = 0x83
	OP_AND
	OP_OR
	OP_XOR
	OP_EQUAL       = 0x87
	OP_EQUALVERIFY = 0x88

	//cripto
	OP_RIPEMD160           = 0xA7
	OP_SHA1                = 0xA8
	OP_SHA256              = 0xA9
	OP_HASH160             = 0xAA
	OP_HASH256             = 0xAB
	OP_CHECKSIG            = 0xAC
	OP_CHECKSIGVERIFY      = 0xAD
	OP_CHECKMULTISIG       = 0xAE
	OP_CHECKMULTISIGVERIFY = 0xAF
)

type SigHash struct {
	DER           byte
	SignatureSize byte
	RType         byte
	RBuf          ghostBytes.HashBytes
	SType         byte
	SSize         byte
	SBuf          ghostBytes.HashBytes
	SignatureType uint32
	PubKeySize    byte
	PubKey        ghostBytes.HashBytes
}
