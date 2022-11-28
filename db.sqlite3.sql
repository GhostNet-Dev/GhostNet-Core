BEGIN TRANSACTION;
-- name: create-paired_block
CREATE TABLE IF NOT EXISTS "paired_block" (
	"Id"	integer,
	"Version"	integer,
	"PreviousBlockHeaderHash"	blob,
	"MerkleRoot"	blob,
	"DataBlockHeaderHash"	blob,
	"Timestamp"	integer,
	"Bits"	integer,
	"Nonce"	integer,
	"AliceCount"	integer,
	"TransactionCount"	integer,
	"SignatureSize"	integer,
	"SigHash"	blob,
	"Data_PreviousBlockHeaderHash"	blob,
	"Data_MerkleRoot"	blob,
	"Data_Nonce"	integer,
	"Data_TransactionCount"	integer
);
-- name: create-transaction
CREATE TABLE IF NOT EXISTS "transactions" (
	"Id"	integer NOT NULL,
	"TxId"	blob,
	"Type"	integer,
	"BlockId"	integer,
	"InputCounter"	integer,
	"OutputCounter"	integer,
	"Nonce"	integer,
	"LockTime"	integer,
	"TxIndex"	integer,
	PRIMARY KEY("Id" AUTOINCREMENT)
);
-- name: create-inputs
CREATE TABLE IF NOT EXISTS "inputs" (
	"Id"	integer NOT NULL,
	"TxId"	blob,
	"BlockId"	integer,
	"prev_TxId"	blob,
	"prev_OutIndex"	integer,
	"Sequence"	integer,
	"ScriptSize"	integer,
	"Script"	blob,
	"Index"	integer,
	PRIMARY KEY("Id" AUTOINCREMENT)
);
-- name: create-outputs
CREATE TABLE IF NOT EXISTS "outputs" (
	"Id"	integer NOT NULL,
	"TxId"	blob,
	"BlockId"	integer,
	"ToAddr"	blob,
	"BrokerAddr"	blob,
	"Type"	integer,
	"Value"	integer,
	"ScriptSize"	integer,
	"Script"	blob,
	"OutputIndex"	integer,
	PRIMARY KEY("Id" AUTOINCREMENT)
);
-- name: create-data-transactions
CREATE TABLE IF NOT EXISTS "data_transactions" (
	"Id"	integer NOT NULL,
	"TxId"	blob,
	"BlockId"	integer,
	"LogicalAddress"	integer,
	"Data"	blob,
	"DataSize"	integer,
	"TxIndex"	integer,
	PRIMARY KEY("Id" AUTOINCREMENT)
);
-- name: create-id
CREATE INDEX IF NOT EXISTS "paired_block_Id" ON "paired_block" (
	"Id"
);
-- name: create-id
CREATE INDEX IF NOT EXISTS "transactions_TxId" ON "transactions" (
	"TxId"
);
-- name: create-id
CREATE INDEX IF NOT EXISTS "transactions_BlockId" ON "transactions" (
	"BlockId"
);
-- name: create-id
CREATE INDEX IF NOT EXISTS "inputs_TxId" ON "inputs" (
	"TxId"
);
-- name: create-id
CREATE INDEX IF NOT EXISTS "inputs_BlockId" ON "inputs" (
	"BlockId"
);
-- name: create-id
CREATE INDEX IF NOT EXISTS "inputs_prev_TxId" ON "inputs" (
	"prev_TxId"
);
-- name: create-id
CREATE INDEX IF NOT EXISTS "outputs_TxId" ON "outputs" (
	"TxId"
);
-- name: create-id
CREATE INDEX IF NOT EXISTS "outputs_BlockId" ON "outputs" (
	"BlockId"
);
-- name: create-id
CREATE INDEX IF NOT EXISTS "outputs_ToAddr" ON "outputs" (
	"ToAddr"
);
-- name: create-id
CREATE INDEX IF NOT EXISTS "data_transactions_TxId" ON "data_transactions" (
	"TxId"
);
-- name: create-id
CREATE INDEX IF NOT EXISTS "data_transactions_BlockId" ON "data_transactions" (
	"BlockId"
);
COMMIT;
