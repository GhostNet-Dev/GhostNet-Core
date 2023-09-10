package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "embed"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	_ "github.com/mattn/go-sqlite3"
)

type GSqlite3 struct {
	db       *sql.DB
	filepath string
}

var (
	GSqlite      = new(GSqlite3)
	MergeGSqlite = new(GSqlite3)
	//go:embed db.sqlite3.sql
	schemaFileByte []byte
)

func NewGSqlite3() *GSqlite3 {
	return GSqlite
}

func NewMergeGSqlite() *GSqlite3 {
	return MergeGSqlite
}

// OpenSQL sql Open
func (gSql *GSqlite3) OpenSQL(path string, filename string) error {
	filepath := path + filename
	db, err := sql.Open("sqlite3", filepath+"?cache=shared&mode=rwc")
	if err != nil {
		log.Fatal(err)
		defer db.Close()
	}
	gSql.db = db
	gSql.filepath = filepath
	return err
}

func (gSql *GSqlite3) CloseSQL() {
	gSql.db.Close()
}

// CreateTable ..
func (gSql *GSqlite3) CreateTable() error {
	query := string(schemaFileByte)
	if _, err := gSql.db.Exec(query); err != nil {
		log.Fatal(err.Error())
	}
	return nil
}

func (gSql *GSqlite3) DropTable() {
	gSql.db.Close()
	e := os.Remove(gSql.filepath)
	if e != nil {
		log.Fatal(e)
	}
}

func (gSql *GSqlite3) deletePairedBlockAfterTargetId(blockId uint32) {
	_, err := gSql.db.Exec(fmt.Sprint("delete from paired_block where Id >=", blockId))
	if err != nil {
		log.Fatal(err)
	}
}

func (gSql *GSqlite3) DeleteBlock(blockId uint32) {
	_, err := gSql.db.Exec(fmt.Sprint("delete from paired_block where Id == ", blockId))
	if err != nil {
		log.Fatal(err)
	}
	tables := []string{"transactions", "data_transactions", "inputs", "outputs"}
	for _, table := range tables {
		_, err := gSql.db.Exec(fmt.Sprint("delete from ", table, " where BlockId == ", blockId))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (gSql *GSqlite3) DeleteTx(txId []byte) {
	tables := []string{"transactions", "inputs", "outputs"}
	for _, table := range tables {
		_, err := gSql.db.Exec(fmt.Sprint("delete from ", table, " where TxId == ", txId))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (gSql *GSqlite3) DeleteAfterTargetId(blockId uint32) {
	gSql.deletePairedBlockAfterTargetId(blockId)

	tables := []string{"transactions", "data_transactions", "inputs", "outputs"}
	for _, table := range tables {
		_, err := gSql.db.Exec(fmt.Sprint("delete from ", table, " where BlockId >=", blockId))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (gSql *GSqlite3) InsertBlock(pair *types.PairedBlock) {
	if gSql.CheckExistBlockId(pair.BlockId()) {
		return
	}
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
		gSql.InsertTx(header.Id, &tx, types.AliceTx, uint32(i))
	}
	for i, tx := range pair.Block.Transaction {
		gSql.InsertTx(header.Id, &tx, types.NormalTx, uint32(i))
	}
	for i, tx := range pair.DataBlock.Transaction {
		gSql.InsertDataTx(header.Id, &tx, uint32(i))
	}
}

func (gSql *GSqlite3) SelectBlockHeader(blockId uint32) (*types.GhostNetBlockHeader, *types.GhostNetDataBlockHeader) {
	rows, err := gSql.db.Query(`select "Id", "Version","PreviousBlockHeaderHash","MerkleRoot","DataBlockHeaderHash","Timestamp","Bits",
		"Nonce", "AliceCount", "TransactionCount", "SignatureSize", "SigHash", 
		"Data_PreviousBlockHeaderHash", "Data_MerkleRoot", "Data_Nonce", "Data_TransactionCount"
		from paired_block 
		where Id = ?`, blockId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var blockSig []byte
	header := types.GhostNetBlockHeader{}
	dataHeader := types.GhostNetDataBlockHeader{}
	count := 0
	for rows.Next() {
		if err = rows.Scan(&header.Id, &header.Version, &header.PreviousBlockHeaderHash,
			&header.MerkleRoot, &header.DataBlockHeaderHash, &header.TimeStamp, &header.Bits,
			&header.Nonce, &header.AliceCount, &header.TransactionCount, &header.SignatureSize,
			&blockSig, &dataHeader.PreviousBlockHeaderHash, &dataHeader.MerkleRoot,
			&dataHeader.Nonce, &dataHeader.TransactionCount); err != nil {
			log.Fatal(err)
		}
		count++
	}
	if count == 0 {
		return nil, nil
	}
	header.BlockSignature = types.SigHash{}
	header.BlockSignature.DeserializeSigHashFromByte(blockSig)
	dataHeader.Id = header.Id
	dataHeader.Version = header.Version
	return &header, &dataHeader
}

func (gSql *GSqlite3) SelectBlock(blockId uint32) *types.PairedBlock {
	header, dataHeader := gSql.SelectBlockHeader(blockId)
	if header == nil || dataHeader == nil {
		return nil
	}

	pair := types.PairedBlock{
		Block: types.GhostNetBlock{
			Header:      *header,
			Alice:       gSql.SelectTxs(header.Id, types.AliceTx),
			Transaction: gSql.SelectTxs(header.Id, types.NormalTx),
		},
		DataBlock: types.GhostNetDataBlock{
			Header:      *dataHeader,
			Transaction: gSql.SelectDataTxs(header.Id),
		},
	}
	if pair.Block.Alice == nil || pair.Block.Transaction == nil {
		gSql.DeleteBlock(pair.BlockId())
		return nil
	}
	return &pair
}

// InsertTx ..
func (gSql *GSqlite3) InsertTx(blockId uint32, tx *types.GhostTransaction, txType uint32, txIndexInBlock uint32) {
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
			("TxId","BlockId","ToAddr","BrokerAddr", "Type", "Value", "ScriptSize","Script", "ScriptExSize", "ScriptEx",
			"OutputIndex") VALUES (?,?,?,? ,?,?,?,?, ?,?,?);
			`, tx.TxId, blockId, output.Addr, output.BrokerAddr, output.Type, output.Value, output.ScriptSize, output.ScriptPubKey,
			output.ScriptExSize, output.ScriptEx, i)
	}
}

// InsertDataTx ..
func (gSql *GSqlite3) InsertDataTx(blockId uint32, dataTx *types.GhostDataTransaction, txIndexInBlock uint32) {
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
	txs := []types.GhostTransaction{}
	for rows.Next() {
		tx := types.GhostTransaction{}
		if err := rows.Scan(&tx.TxId, &tx.Body.InputCounter, &tx.Body.OutputCounter,
			&tx.Body.Nonce, &tx.Body.LockTime); err != nil {
			log.Fatal(err)
		}
		txs = append(txs, tx)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	rows.Close()

	return gSql.GetTxRows(txs)
}

// SelectTx 역시나 여긴 인터페이스 역할을하고 아래 쿼리들은 하위 클래스에 할당해야할 것 같음
func (gSql *GSqlite3) SelectTx(TxId []byte) (tx *types.GhostTransaction, blockId uint32) {
	if ok := gSql.CheckExistTxId(TxId); !ok {
		return nil, 0
	}

	tx = &types.GhostTransaction{TxId: TxId}
	rows, err := gSql.db.Query(`select BlockId, InputCounter, OutputCounter, Nonce, LockTime from transactions tx 
		where TxId = ? order by TxIndex`, TxId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&blockId, &tx.Body.InputCounter, &tx.Body.OutputCounter,
			&tx.Body.Nonce, &tx.Body.LockTime); err != nil {
			log.Fatal(err)
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	input := gSql.SelectInputs(TxId, tx.Body.InputCounter)
	if input == nil {
		return nil, 0
	}
	tx.Body.Vin = input
	output := gSql.selectOutputs(TxId, tx.Body.OutputCounter)
	if output == nil {
		return nil, 0
	}
	tx.Body.Vout = output

	return tx, blockId
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
		if !rows.Next() {
			return nil
		}
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

func (gSql *GSqlite3) selectOutputs(TxId []byte, count uint32) []types.TxOutput {
	outputs := make([]types.TxOutput, count)

	rows, err := gSql.db.Query(`select ToAddr, BrokerAddr, Script, ScriptSize, ScriptEx, ScriptExSize, Type, Value from outputs 
		where TxId = ? limit 0, ?`, TxId, count)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for i, output := range outputs {
		if !rows.Next() {
			return nil
		}
		if err = rows.Scan(&output.Addr, &output.BrokerAddr, &output.ScriptPubKey, &output.ScriptSize,
			&output.ScriptEx, &output.ScriptExSize, &output.Type, &output.Value); err != nil {
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

func (gSql *GSqlite3) SelectOutputs(txType types.TxOutputType, start, count int) []types.PrevOutputParam {
	outputs := []types.PrevOutputParam{}

	rows, err := gSql.db.Query(`select outputs.TxId, outputs.ToAddr, outputs.BrokerAddr, outputs.Script, outputs.ScriptSize, 
		outputs.ScriptEx, outputs.ScriptExSize, outputs.Type, outputs.Value, outputs.OutputIndex from outputs 
		where outputs.Type = ?
		order by outputs.BlockId DESC limit ?, ?`, txType, start, count)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		output := types.PrevOutputParam{}
		if err = rows.Scan(&output.VOutPoint.TxId, &output.Vout.Addr, &output.Vout.BrokerAddr, &output.Vout.ScriptPubKey,
			&output.Vout.ScriptSize, &output.Vout.ScriptEx, &output.Vout.ScriptExSize, &output.Vout.Type, &output.Vout.Value, &output.VOutPoint.TxOutIndex); err == sql.ErrNoRows {
			return nil
		} else if err != nil {
			log.Fatal(err)
		}
		output.TxType = txType
		outputs = append(outputs, output)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return outputs
}

func (gSql *GSqlite3) SearchOutput(txType types.TxOutputType, toAddr, uniqKey []byte) []types.PrevOutputParam {
	outputs := []types.PrevOutputParam{}

	rows, err := gSql.db.Query(`select outputs.TxId, outputs.ToAddr, outputs.BrokerAddr, outputs.Script, outputs.ScriptSize, 
		outputs.ScriptEx, outputs.ScriptExSize, outputs.Type, outputs.Value, outputs.OutputIndex from outputs 
		where outputs.ToAddr = ? and  outputs.Type = ? and outputs.Script = ?
		order by outputs.BlockId DESC`, toAddr, txType, uniqKey)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		output := types.PrevOutputParam{}
		if err = rows.Scan(&output.VOutPoint.TxId, &output.Vout.Addr, &output.Vout.BrokerAddr, &output.Vout.ScriptPubKey,
			&output.Vout.ScriptSize, &output.Vout.ScriptEx, &output.Vout.ScriptExSize, &output.Vout.Type, &output.Vout.Value, &output.VOutPoint.TxOutIndex); err == sql.ErrNoRows {
			return nil
		} else if err != nil {
			log.Print(err)
			return nil
		}
		output.TxType = txType
		outputs = append(outputs, output)
	}

	if err = rows.Err(); err != nil {
		log.Print(err)
		return nil
	}

	return outputs
}

func (gSql *GSqlite3) SearchOutputs(txType types.TxOutputType, toAddr []byte) []types.PrevOutputParam {
	outputs := []types.PrevOutputParam{}

	rows, err := gSql.db.Query(`select outputs.TxId, outputs.ToAddr, outputs.BrokerAddr, outputs.Script, outputs.ScriptSize, 
		outputs.ScriptEx, outputs.ScriptExSize, outputs.Type, outputs.Value, outputs.OutputIndex from outputs 
		where outputs.ToAddr = ? and  outputs.Type = ?
		order by outputs.BlockId DESC`, toAddr, txType)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		output := types.PrevOutputParam{}
		if err = rows.Scan(&output.VOutPoint.TxId, &output.Vout.Addr, &output.Vout.BrokerAddr, &output.Vout.ScriptPubKey,
			&output.Vout.ScriptSize, &output.Vout.ScriptEx, &output.Vout.ScriptExSize, &output.Vout.Type, &output.Vout.Value, &output.VOutPoint.TxOutIndex); err == sql.ErrNoRows {
			return nil
		} else if err != nil {
			log.Fatal(err)
		}
		output.TxType = txType
		outputs = append(outputs, output)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return outputs
}

func (gSql *GSqlite3) SearchStringOutputs(txType types.TxOutputType,
	toAddr, keyword []byte) []types.PrevOutputParam {
	outputs := []types.PrevOutputParam{}

	rows, err := gSql.db.Query(`select outputs.TxId, outputs.ToAddr, outputs.BrokerAddr, outputs.Script, outputs.ScriptSize, 
		outputs.ScriptEx, outputs.ScriptExSize, outputs.Type, outputs.Value, outputs.OutputIndex from outputs 
		where outputs.ToAddr = ? and  outputs.Type = ? 
		and outputs.Script like '%'||?
		order by outputs.BlockId DESC, outputs.Script DESC`, toAddr, txType, keyword)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		output := types.PrevOutputParam{}
		if err = rows.Scan(&output.VOutPoint.TxId, &output.Vout.Addr, &output.Vout.BrokerAddr, &output.Vout.ScriptPubKey,
			&output.Vout.ScriptSize, &output.Vout.ScriptEx, &output.Vout.ScriptExSize, &output.Vout.Type, &output.Vout.Value, &output.VOutPoint.TxOutIndex); err == sql.ErrNoRows {
			return nil
		} else if err != nil {
			log.Fatal(err)
		}
		output.TxType = txType
		outputs = append(outputs, output)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return outputs
}
func (gSql *GSqlite3) SelectOutputLatests(txType types.TxOutputType,
	toAddr []byte, start, count int) []types.PrevOutputParam {
	outputs := []types.PrevOutputParam{}

	rows, err := gSql.db.Query(`select outputs.TxId, outputs.ToAddr, outputs.BrokerAddr, outputs.Script, outputs.ScriptSize, 
		outputs.ScriptEx, outputs.ScriptExSize, outputs.Type, outputs.Value, outputs.OutputIndex from outputs 
		where outputs.Type = ? and outputs.ToAddr = ?
		order by outputs.Script DESC limit ?, ?`, txType, toAddr, start, count)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		output := types.PrevOutputParam{}
		if err = rows.Scan(&output.VOutPoint.TxId, &output.Vout.Addr, &output.Vout.BrokerAddr, &output.Vout.ScriptPubKey,
			&output.Vout.ScriptSize, &output.Vout.ScriptEx, &output.Vout.ScriptExSize, &output.Vout.Type, &output.Vout.Value, &output.VOutPoint.TxOutIndex); err == sql.ErrNoRows {
			return nil
		} else if err != nil {
			log.Fatal(err)
		}
		output.TxType = txType
		outputs = append(outputs, output)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return outputs
}

func (gSql *GSqlite3) SelectUnusedOutputs(txType types.TxOutputType, toAddr []byte) []types.PrevOutputParam {
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

	for rows.Next() {
		output := types.PrevOutputParam{}
		if err = rows.Scan(&output.VOutPoint.TxId, &output.Vout.Addr, &output.Vout.BrokerAddr, &output.Vout.ScriptPubKey,
			&output.Vout.ScriptSize, &output.Vout.Type, &output.Vout.Value, &output.VOutPoint.TxOutIndex); err == sql.ErrNoRows {
			return nil
		} else if err != nil {
			log.Fatal(err)
		}
		output.TxType = txType
		outputs = append(outputs, output)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return outputs
}

func (gSql *GSqlite3) GetNicknameToAddress(nickname []byte) (toAddr []byte) {
	query, err := gSql.db.Prepare("select outputs.ToAddr from transactions where Type = ? and ScriptEx = ?")
	if err != nil {
		log.Printf("%s", err)
	}
	defer query.Close()

	if err := query.QueryRow(types.TxTypeFSRoot, nickname).Scan(&toAddr); err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		log.Print(err)
	}
	return toAddr
}

func (gSql *GSqlite3) GetMaxLogicalAddress(toAddr []byte) (maxLogicalAddr uint64, err error) {
	rows, err := gSql.db.Query(`select LogicalAddress from data_transactions 
		left outer join outputs  on data_transactions.TxId = outputs.TxId 
		where outputs.ToAddr = ? and  outputs.Type = ?  and  inputs.Id is NULL
		order by outputs.BlockId ASC`, toAddr, types.TxTypeDataStore)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&maxLogicalAddr); err == sql.ErrNoRows {
			return 0, errors.New("there is no data tx")
		} else if err != nil {
			log.Fatal(err)
		}
	}
	return maxLogicalAddr, nil
}

func (gSql *GSqlite3) GetIssuedCoin(blockId uint32) uint64 {
	var coin uint64
	query, err := gSql.db.Prepare(`select sum(Value) from outputs 
	left outer join transactions on outputs.TxId = transactions.TxId
	where transactions.Type = ? and transactions.BlockId = ?`)
	if err != nil {
		log.Printf("%s", err)
	}
	defer query.Close()

	if err := query.QueryRow(types.AliceTx, blockId).Scan(&coin); err == sql.ErrNoRows {
		return 0
	} else if err != nil {
		log.Print(err)
	}
	return coin
}

func (gSql *GSqlite3) CheckExistBlockId(blockId uint32) bool {
	var count uint32
	query, err := gSql.db.Prepare("select count(*) from paired_block where Id = ?")
	if err != nil {
		log.Printf("%s", err)
	}
	defer query.Close()

	if err := query.QueryRow(blockId).Scan(&count); err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Print(err)
	}
	return count > 0
}

func (gSql *GSqlite3) CheckExistTxId(txId []byte) bool {
	var count uint32
	query, err := gSql.db.Prepare("select count(*) from transactions where TxId = ?")
	if err != nil {
		log.Printf("%s", err)
	}
	defer query.Close()

	if err := query.QueryRow(txId).Scan(&count); err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Print(err)
	}
	return count > 0
}

func (gSql *GSqlite3) CheckExistRefOutput(refTxId []byte, outIndex uint32, notTxId []byte) bool {
	var count uint32
	query, err := gSql.db.Prepare(`select count(*) from outputs where 
		TxId == ? and OutputIndex == ?`)
	if err != nil {
		log.Printf("%s", err)
	}
	defer query.Close()

	if err := query.QueryRow(refTxId, outIndex).Scan(&count); err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Print(err)
	}
	return count > 0
}

func (gSql *GSqlite3) CheckExistFsRoot(nickname []byte) bool {
	var count uint32
	query, err := gSql.db.Prepare(`select count(*) from outputs where 
		ScriptEx == ?`)
	if err != nil {
		log.Printf("%s", err)
	}
	defer query.Close()

	if err := query.QueryRow(nickname).Scan(&count); err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Print(err)
	}
	return count > 0
}

func (gSql *GSqlite3) GetBlockHeight() uint32 {
	var id uint32
	query, err := gSql.db.Prepare(`select Id from paired_block order by Id desc limit 0, 1`)
	if err != nil {
		log.Printf("%s", err)
	}
	defer query.Close()

	if err := query.QueryRow().Scan(&id); err == sql.ErrNoRows {
		return 0
	} else if err != nil {
		log.Print(err)
	}
	return id
}

func (gSql *GSqlite3) GetTxRows(txs []types.GhostTransaction) []types.GhostTransaction {
	for i, tx := range txs {
		input := gSql.SelectInputs(tx.TxId, tx.Body.InputCounter)
		if input == nil {
			return nil
		}
		txs[i].Body.Vin = input
		output := gSql.selectOutputs(tx.TxId, tx.Body.OutputCounter)
		if output == nil {
			return nil
		}
		txs[i].Body.Vout = output
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
