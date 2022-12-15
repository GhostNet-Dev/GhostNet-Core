package gsql

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	_ "github.com/mattn/go-sqlite3"
)

type GSqlite3 struct {
	db       *sql.DB
	filepath string
}

// OpenSQL sql Open
func (gSql *GSqlite3) OpenSQL(path string) error {
	db, err := sql.Open("sqlite3", path+"block.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatal(err)
		defer db.Close()
	}
	gSql.db = db
	gSql.filepath = path + "block.db"
	return err
}

func (gSql *GSqlite3) CloseSQL() {
	gSql.db.Close()
}

// CreateTable ..
func (gSql *GSqlite3) CreateTable(schemaFile string) error {
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

	return err
}

func (gSql *GSqlite3) DropTable() {
	gSql.db.Close()
	e := os.Remove(gSql.filepath)
	if e != nil {
		log.Fatal(e)
	}
}

func (gSql *GSqlite3) InsertBlock(pair types.PairedBlock) {
	header := pair.Block.Header
	dataHeader := pair.DataBlock.Header
	gSql.InsertQuery(`INSERT INTO "paired_block" 
		("Id", "Version","PreviousBlockHeaderHash","MerkleRoot","DataBlockHeaderHash"
		,"Timestamp","Bits", "Nonce", "AliceCount", "TransactionCount", "SignatureSize"
		, "SigHash", "Data_PreviousBlockHeaderHash", "Data_MerkleRoot", "Data_Nonce", 
		"Data_TransactionCount") VALUES (?,?,?,?, ?,?,?,? ,?,?,?,? ,?,?,?,?);
		`, header.Id, header.Version, header.PreviousBlockHeaderHash, header.MerkleRoot,
		header.DataBlockHeaderHash, header.TimeStamp, header.Bits, header.Nonce,
		header.AliceCount, header.TransactionCount, header.SignatureSize,
		header.BlockSignature.SerializeToByte(), dataHeader.PreviousBlockHeaderHash,
		dataHeader.MerkleRoot, dataHeader.Nonce, dataHeader.TransactionCount)

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
		where Id = ?`, blockId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	header := types.GhostNetBlockHeader{}
	dataHeader := types.GhostNetDataBlockHeader{}
	for rows.Next() {
		blockSigHash := []byte{}
		if err = rows.Scan(&header.Id, &header.Version, &header.PreviousBlockHeaderHash,
			&header.MerkleRoot, &header.DataBlockHeaderHash, &header.TimeStamp, &header.Bits,
			&header.Nonce, &header.AliceCount, &header.TransactionCount, &header.SignatureSize,
			&blockSigHash, &dataHeader.PreviousBlockHeaderHash, &dataHeader.MerkleRoot,
			&dataHeader.Nonce, &dataHeader.TransactionCount); err != nil {
			log.Fatal(err)
		}
		header.BlockSignature.DeserializeSigHashFromByte(blockSigHash)
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
			Transaction: gSql.SelectDataTxs(header.Id),
		},
	}
	return &pair
}

// InsertTx ..
func (gSql *GSqlite3) InsertTx(blockId uint32, tx types.GhostTransaction, txType uint32, txIndexInBlock uint32) {
	gSql.InsertQuery(`INSERT INTO "transactions" ("TxId", "Type", "BlockId","InputCounter",
	"OutputCounter","Nonce","LockTime","TxIndex") VALUES (?,?,?,?, ?,?,?,?);
		`, tx.TxId, txType, blockId, tx.Body.InputCounter, tx.Body.OutputCounter, tx.Body.Nonce, tx.Body.LockTime, txIndexInBlock)
	for i, input := range tx.Body.Vin {
		gSql.InsertQuery(`INSERT INTO "inputs" 
			("TxId","BlockId","prev_TxId","prev_OutIndex","Sequence","ScriptSize", "Script", 
			"Index") VALUES (?,?,?,?, ?,?,?,?);
			`, tx.TxId, blockId, input.PrevOut.TxId, input.PrevOut.TxOutIndex, input.Sequence, input.ScriptSize, input.ScriptSig, i)
	}
	for i, output := range tx.Body.Vout {
		gSql.InsertQuery(`INSERT INTO "outputs" 
			("TxId","BlockId","ToAddr","BrokerAddr", "Type", "Value", "ScriptSize","Script",
			"OutputIndex") VALUES (?,?,?,? ,?,?,?,?, ?);
			`, tx.TxId, blockId, output.Addr, output.BrokerAddr, output.Type, output.Value, output.ScriptSize, output.ScriptPubKey, i)
	}
}

// InsertDataTx ..
func (gSql *GSqlite3) InsertDataTx(blockId uint32, dataTx types.GhostDataTransaction, txIndexInBlock uint32) {
	gSql.InsertQuery(`INSERT INTO "data_transactions" ("TxId","BlockId","LogicalAddress","Data",
		"DataSize","TxIndex") VALUES (?,?,?,?, ?,?);`,
		dataTx.TxId, blockId, dataTx.LogicalAddress, dataTx.Data, dataTx.DataSize, txIndexInBlock)
}

func (gSql *GSqlite3) SelectDataTxs(blockId uint32) []types.GhostDataTransaction {
	rows, err := gSql.db.Query(`select TxId, LogicalAddress, DataSize, Data from data_transactions 
		where BlockId = ? order by TxIndex`, blockId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	return gSql.GetDataTxRows(rows)
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
	rows, err := gSql.db.Query(`select TxId, InputCounter, OutputCounter, Nonce, LockTime 
	from transactions tx where BlockId = ? and Type = ? order by TxIndex`, blockId, txType)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	return gSql.GetTxRows(rows)
}

// SelectTx 역시나 여긴 인터페이스 역할을하고 아래 쿼리들은 하위 클래스에 할당해야할 것 같음
func (gSql *GSqlite3) SelectTx(TxId []byte) *types.GhostTransaction {
	if ok := gSql.CheckExistTxId(TxId); ok == false {
		return nil
	}

	tx := types.GhostTransaction{TxId: TxId}

	rows, err := gSql.db.Query(`select InputCounter, OutputCounter, Nonce, LockTime from transactions tx 
		where TxId = ? order by TxIndex`, TxId)
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
		input.PrevOut = types.TxOutPoint{TxId: prev_TxId, TxOutIndex: prev_OutIndex}
		inputs[i] = input
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return inputs
}

func (gSql *GSqlite3) SelectOutputs(TxId []byte, count uint32) []types.TxOutput {
	outputs := make([]types.TxOutput, count)

	rows, err := gSql.db.Query(`select ToAddr, BrokerAddr, Script, ScriptSize, Type, Value from outputs 
		where TxId = ? limit 0, ?`, TxId, count)
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

func (gSql *GSqlite3) SelectUnusedOutputs(txType uint32, toAddr []byte) []types.PrevOutputParam {
	outputs := []types.PrevOutputParam{}

	rows, err := gSql.db.Query(`select outputs.TxId, outputs.ToAddr, outputs.BrokerAddr, outputs.Script, outputs.ScriptSize, 
		outputs.Type, outputs.Value, outputs.OutputIndex from outputs 
		left outer join inputs  on inputs.prev_TxId = outputs.TxId and inputs.prev_OutIndex = outputs.OutputIndex 
		where outputs.ToAddr = ? and  outputs.Type = ?  and  inputs.Id is NULL
		order by outputs.BlockId ASC`, toAddr, txType)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for i, output := range outputs {
		rows.Next()
		if err = rows.Scan(&output.VOutPoint.TxId, &output.Vout.Addr, &output.Vout.BrokerAddr, &output.Vout.ScriptPubKey,
			&output.Vout.ScriptSize, &output.Vout.Type, &output.Vout.Value, &output.VOutPoint.TxOutIndex); err != nil {
			log.Fatal(err)
		}
		outputs[i] = output
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return outputs
}

func (gSql *GSqlite3) CheckExistTxId(txId []byte) bool {
	var count uint32
	query, err := gSql.db.Prepare("select count(*) from transactions where TxId = ?")
	if err != nil {
		log.Printf("%s", err)
	}
	defer query.Close()

	err = query.QueryRow(txId).Scan(&count)
	return count > 0
}

func (gSql *GSqlite3) CheckExistRefOutout(refTxId []byte, outIndex uint32, notTxId []byte) bool {
	var count uint32
	query, err := gSql.db.Prepare(`select count(*) from inputs where 
		prev_TxId == ? and prev_OutputIndex == ? and TxId != ?`)
	if err != nil {
		log.Printf("%s", err)
	}
	defer query.Close()

	err = query.QueryRow(refTxId, outIndex, notTxId).Scan(&count)
	return count > 0
}

func (gSql *GSqlite3) GetBlockHeight() uint32 {
	var id uint32
	query, err := gSql.db.Prepare(`select Id from paired_block order by Id desc limit 0, 1`)
	if err != nil {
		log.Printf("%s", err)
	}
	defer query.Close()

	err = query.QueryRow().Scan(&id)
	return id
}

func (gSql *GSqlite3) GetMinPoolId() uint32 {
	var min uint32
	query, err := gSql.db.Prepare(`select min(TxIndex) from c_transactions`)
	if err != nil {
		log.Printf("%s", err)
	}
	defer query.Close()

	err = query.QueryRow().Scan(&min)
	return min
}

func (gSql *GSqlite3) GetMaxPoolId() uint32 {
	var min uint32
	query, err := gSql.db.Prepare(`select max(TxIndex) from c_transactions`)
	if err != nil {
		log.Printf("%s", err)
	}
	defer query.Close()

	err = query.QueryRow().Scan(&min)
	return min
}

func (gSql *GSqlite3) UpdatePoolId(oldPoolId uint32, newPoolId uint32) {
	gSql.InsertQuery(`update c_transactions set TxIndex = ? where TxIndex == ?;`,
		oldPoolId, newPoolId)
}

func (gSql *GSqlite3) SelectTxsPool(poolId uint32) []types.GhostTransaction {
	rows, err := gSql.db.Query(`select TxId, InputCounter, OutputCounter, Nonce, LockTime 
	from c_transactions tx where TxIndex = ?`, poolId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	return gSql.GetTxRows(rows)
}

func (gSql *GSqlite3) SelectDataTxsPool(poolId uint32) []types.GhostDataTransaction {
	rows, err := gSql.db.Query(`select TxId, LogicalAddress, DataSize, Data from data_transactions 
		where TxIndex = ?`, poolId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	return gSql.GetDataTxRows(rows)
}

func (gSql *GSqlite3) GetTxRows(rows *sql.Rows) []types.GhostTransaction {
	txs := []types.GhostTransaction{}
	for rows.Next() {
		tx := types.GhostTransaction{}
		if err := rows.Scan(&tx.TxId, &tx.Body.InputCounter, &tx.Body.OutputCounter,
			&tx.Body.Nonce, &tx.Body.LockTime); err != nil {
			log.Fatal(err)
		}
		tx.Body.Vin = gSql.SelectInputs(tx.TxId, tx.Body.InputCounter)
		tx.Body.Vout = gSql.SelectOutputs(tx.TxId, tx.Body.OutputCounter)
		txs = append(txs, tx)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return txs
}

func (gSql *GSqlite3) GetDataTxRows(rows *sql.Rows) []types.GhostDataTransaction {
	dataTxs := []types.GhostDataTransaction{}
	for rows.Next() {
		dataTx := types.GhostDataTransaction{}
		if err := rows.Scan(&dataTx.TxId, &dataTx.LogicalAddress, &dataTx.DataSize, &dataTx.Data); err != nil {
			log.Fatal(err)
		}
		dataTxs = append(dataTxs, dataTx)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return dataTxs
}
