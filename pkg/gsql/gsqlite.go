package gsql

import (
	"database/sql"
	"io/ioutil"
	"log"

	types "github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	_ "github.com/mattn/go-sqlite3"
)

type GSqlite3 struct {
	db *sql.DB
}

// OpenSQL sql Open
func (gSql *GSqlite3) OpenSQL(path string) {
	db, err := sql.Open("sqlite3", path+"block.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatal(err)
		defer db.Close()
	}
	gSql.db = db
}

// CreateTable ..
func (gSql *GSqlite3) CreateTable(schemaFile string) {
	file, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Execute all
	query := string(file)
	_, err = gSql.db.Exec(query)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (gSql *GSqlite3) InsertBlock(pair types.PairedBlock) {
	header := pair.Block.Header
	dataHeader := pair.DataBlock.Header
	gSql.InsertQuery(`INSERT INTO "paired_block" 
		("Id", "Version","PreviousBlockHeaderHash","MerkleRoot","DataBlockHeaderHash","Timestamp","Bits",
		"Nonce", "AliceCount", "TransactionCount", "SignatureSize", "SigHash", 
		"Data_PreviousBlockHeaderHash", "Data_MerkleRoot", "Data_Nonce", "Data_TransactionCount") 
		VALUES (?,?,?,?, ?,?,?,? ,?,?,?,? ,?,?,?,?);
		`, header.Version, header.PreviousBlockHeaderHash, header.MerkleRoot, header.DataBlockHeaderHash,
		header.TimeStamp, header.Bits, header.Nonce, header.AliceCount, header.TransactionCount,
		header.SignatureSize, header.BlockSignature,
		dataHeader.PreviousBlockHeaderHash, dataHeader.MerkleRoot, dataHeader.Nonce, dataHeader.TransactionCount)

	for i, tx := range pair.Block.Alice {
		gSql.InsertTx(header.Id, tx, types.AliceTx, uint32(i))
	}
	for i, tx := range pair.Block.Transaction {
		gSql.InsertTx(header.Id, tx, types.NormalTx, uint32(i))
	}
	for i, tx := range pair.DataBlock.Transaction {
		gSql.InsertDataTx(header.Id, tx, uint32(i))
	}
}

func (gSql *GSqlite3) SelectBlock(blockId uint32) *types.PairedBlock {
	rows, err := gSql.db.Query(`select "Id", "Version","PreviousBlockHeaderHash","MerkleRoot","DataBlockHeaderHash","Timestamp","Bits",
		"Nonce", "AliceCount", "TransactionCount", "SignatureSize", "SigHash", 
		"Data_PreviousBlockHeaderHash", "Data_MerkleRoot", "Data_Nonce", "Data_TransactionCount"
		from paired_block 
		where BlockId = ?`, blockId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	header := types.GhostNetBlockHeader{}
	dataHeader := types.GhostNetDataBlockHeader{}
	for rows.Next() {
		if err = rows.Scan(&header.Id, &header.Version, &header.PreviousBlockHeaderHash, &header.MerkleRoot, &header.DataBlockHeaderHash, &header.TimeStamp, &header.Bits,
			&header.Nonce, &header.AliceCount, &header.TransactionCount, &header.SignatureSize, &header.BlockSignature,
			&dataHeader.PreviousBlockHeaderHash, &dataHeader.MerkleRoot, &dataHeader.Nonce, &dataHeader.TransactionCount); err != nil {
			log.Fatal(err)
		}
	}
	dataHeader.Id = header.Id
	dataHeader.Version = header.Version
	pair := types.PairedBlock{
		Block: types.GhostNetBlock{
			Header:      header,
			Alice:       gSql.SelectTxs(header.Id, types.AliceTx),
			Transaction: gSql.SelectTxs(header.Id, types.NormalTx),
		},
		DataBlock: types.GhostNetDataBlock{
			Header:      dataHeader,
			Transaction: gSql.SelectDataTx(header.Id),
		},
	}
	return &pair
}

// InsertTx ..
func (gSql *GSqlite3) InsertTx(blockId uint32, tx types.GhostTransaction, txType uint32, txIndexInBlock uint32) {
	gSql.InsertQuery(`INSERT INTO "transactions" ("TxId","BlockId","InputCounter","OutputCounter","Nonce","LockTime","TxIndex") 
		VALUES (?,?,?,?, ?,?,?);
		`, tx.TxId, blockId, tx.Body.InputCounter, tx.Body.OutputCounter, tx.Body.Nonce, tx.Body.LockTime, txIndexInBlock)
	for i, input := range tx.Body.Vin {
		gSql.InsertQuery(`INSERT INTO "inputs" 
			("TxId","BlockId","prev_TxId","prev_OutIndex","Sequence","ScriptSize", "Script", "Index") 
			VALUES (?,?,?,?, ?,?,?,?);
			`, tx.TxId, blockId, input.PrevOut.TxId, input.PrevOut.TxOutIndex, input.Sequence, input.ScriptSize, input.ScriptSig, i)
	}
	for i, output := range tx.Body.Vout {
		gSql.InsertQuery(`INSERT INTO "outputs" 
			("TxId","BlockId","ToAddr","BrokerAddr", "Type", "Value", "ScriptSize","Script","OutputIndex") 
			VALUES (?,?,?,? ,?,?,?,?, ?);
			`, tx.TxId, blockId, output.Addr, output.BrokerAddr, output.Type, output.Value, output.ScriptSize, output.ScriptPubKey, i)
	}
}

// InsertDataTx ..
func (gSql *GSqlite3) InsertDataTx(blockId uint32, dataTx types.GhostDataTransaction, txIndexInBlock uint32) {
	gSql.InsertQuery(`INSERT INTO "data_transactions" ("TxId","BlockId","LogicalAddress","Data","DataSize","TxIndex") 
		VALUES (?,?,?,?, ?,?);`,
		dataTx.TxId, blockId, dataTx.LogicalAddress, dataTx.Data, dataTx.DataSize, txIndexInBlock)
}

func (gSql *GSqlite3) SelectDataTx(blockId uint32) []types.GhostDataTransaction {
	rows, err := gSql.db.Query(`select TxId, LogicalAddress, DataSize, Data from data_transactions 
		where BlockId = ? order by TxIndex`, blockId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	dataTxs := []types.GhostDataTransaction{}
	for rows.Next() {
		dataTx := types.GhostDataTransaction{}
		if err = rows.Scan(&dataTx.TxId, &dataTx.LogicalAddress, &dataTx.DataSize, &dataTx.Data); err != nil {
			log.Fatal(err)
		}
		dataTxs = append(dataTxs, dataTx)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return dataTxs
}

func (gSql *GSqlite3) SelectData(TxId []byte) *types.GhostDataTransaction {
	dataTx := types.GhostDataTransaction{TxId: TxId}

	rows, err := gSql.db.Query(`select LogicalAddress, DataSize, Data from data_transactions 
		where TxId = ?`, TxId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&dataTx.LogicalAddress, &dataTx.DataSize, &dataTx.Data); err != nil {
			log.Fatal(err)
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return &dataTx
}

func (gSql *GSqlite3) SelectTxs(blockId uint32, txType uint32) []types.GhostTransaction {
	rows, err := gSql.db.Query(`select TxId, InputCounter, OutputCounter, Nonce, LockTime from transactions tx 
		where BlockId = ? and Type = ? order by TxIndex`, blockId, txType)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	txs := []types.GhostTransaction{}
	for rows.Next() {
		tx := types.GhostTransaction{}
		if err = rows.Scan(&tx.TxId, &tx.Body.InputCounter, &tx.Body.OutputCounter,
			&tx.Body.Nonce, &tx.Body.LockTime); err != nil {
			log.Fatal(err)
		}
		tx.Body.Vin = gSql.SelectInputs(tx.TxId, tx.Body.InputCounter)
		tx.Body.Vout = gSql.SelectOutputs(tx.TxId, tx.Body.OutputCounter)
		txs = append(txs, tx)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return txs
}

// SelectTx 역시나 여긴 인터페이스 역할을하고 아래 쿼리들은 하위 클래스에 할당해야할 것 같음
func (gSql *GSqlite3) SelectTx(TxId []byte, txType uint32) *types.GhostTransaction {
	tx := types.GhostTransaction{TxId: TxId}

	rows, err := gSql.db.Query(`select InputCounter, OutputCounter, Nonce, LockTime from transactions tx 
		where TxId = ? and Type = ? order by TxIndex`, TxId, txType)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&tx.Body.InputCounter, &tx.Body.OutputCounter,
			&tx.Body.Nonce, &tx.Body.LockTime); err != nil {
			log.Fatal(err)
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	tx.Body.Vin = gSql.SelectInputs(TxId, tx.Body.InputCounter)
	tx.Body.Vout = gSql.SelectOutputs(TxId, tx.Body.OutputCounter)

	return &tx
}

func (gSql *GSqlite3) SelectInputs(TxId []byte, count uint32) []types.TxInput {
	inputs := make([]types.TxInput, count)

	rows, err := gSql.db.Query(`select prev_TxId, prev_OutIndex, Sequence, ScriptSize, Script
	 	from inputs where TxId = ?`, TxId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for i, input := range inputs {
		rows.Next()
		var prev_TxId []byte
		var prev_OutIndex uint32
		if err = rows.Scan(&prev_TxId, &prev_OutIndex, &input.Sequence, &input.ScriptSize,
			&input.ScriptSig); err != nil {
			log.Fatal(err)
		}
		input.PrevOut = types.TxOutPoint{prev_TxId, prev_OutIndex}
		inputs[i] = input
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return inputs
}

func (gSql *GSqlite3) SelectOutputs(TxId []byte, count uint32) []types.TxOutput {
	outputs := make([]types.TxOutput, count)

	rows, err := gSql.db.Query(`select ToAddr, BrokerAddr, Script, ScriptSize, Type, Value from outputs where TxId = ?`, TxId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for i, output := range outputs {
		rows.Next()
		if err = rows.Scan(&output.Addr, &output.BrokerAddr, &output.ScriptPubKey, &output.ScriptSize,
			&output.Type, &output.Value); err != nil {
			log.Fatal(err)
		}
		outputs[i] = output
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return outputs
}

func (gSql *GSqlite3) InsertQuery(query string, args ...interface{}) {
	tx, err := gSql.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
