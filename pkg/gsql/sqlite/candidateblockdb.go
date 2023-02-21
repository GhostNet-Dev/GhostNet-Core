package sqlite

import (
	"database/sql"
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

func (gSql *GSqlite3) GetMinPoolId() uint32 {
	var min uint32
	query, err := gSql.db.Prepare(`select min(TxIndex) from c_transactions`)
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer query.Close()

	if err = query.QueryRow().Scan(&min); err != nil {
		log.Printf("%s", err)
		return 0
	}
	return min
}

func (gSql *GSqlite3) GetMaxPoolId() uint32 {
	var max int
	query, err := gSql.db.Prepare(`select max(TxIndex) from c_transactions`)
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer query.Close()

	if err = query.QueryRow().Scan(&max); err == sql.ErrNoRows {
		return 0
	} else if err != nil {
		log.Printf("%s", err)
		return 0
	}

	return uint32(max)
}

// InsertTx ..
func (gSql *GSqlite3) InsertCandidateTx(tx *types.GhostTransaction, poolId uint32) {
	blockId, txIndexInBlock := poolId, poolId
	gSql.InsertQuery(`INSERT INTO "c_transactions" ("TxId", "InputCounter",
	"OutputCounter","Nonce","LockTime","TxIndex") VALUES (?,?,?,?, ?,?,?,?);
		`, tx.TxId, blockId, tx.Body.InputCounter, tx.Body.OutputCounter, tx.Body.Nonce, tx.Body.LockTime, txIndexInBlock)
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
func (gSql *GSqlite3) InsertCandidateDataTx(dataTx *types.GhostDataTransaction, poolId uint32) {
	blockId, txIndexInBlock := poolId, poolId
	gSql.InsertQuery(`INSERT INTO "c_data_transactions" ("TxId","BlockId","LogicalAddress","Data",
		"DataSize","TxIndex") VALUES (?,?,?,?, ?,?);`,
		dataTx.TxId, blockId, dataTx.LogicalAddress, dataTx.Data, dataTx.DataSize, txIndexInBlock)
}

func (gSql *GSqlite3) UpdatePoolId(oldPoolId uint32, newPoolId uint32) {
	gSql.InsertQuery(`update c_transactions set TxIndex = ? where TxIndex == ?;`,
		oldPoolId, newPoolId)
}

func (gSql *GSqlite3) SelectCandidateTxCount() uint32 {
	var count uint32
	query, err := gSql.db.Prepare("select count(*) from c_transactions")
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer query.Close()

	if err = query.QueryRow().Scan(&count); err != nil {
		log.Printf("%s", err)
		return 0
	}
	return count
}

func (gSql *GSqlite3) SelectTxsPool(poolId uint32) []types.GhostTransaction {
	rows, err := gSql.db.Query(`select TxId, InputCounter, OutputCounter, Nonce, LockTime 
	from c_transactions tx where TxIndex = ?`, poolId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	return gSql.GetCandidateTxRows(rows)
}

func (gSql *GSqlite3) SelectDataTxsPool(poolId uint32) []types.GhostDataTransaction {
	rows, err := gSql.db.Query(`select TxId, LogicalAddress, DataSize, Data from c_data_transactions 
		where TxIndex = ?`, poolId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	return gSql.GetCandidateDataTxRows(rows)
}

func (gSql *GSqlite3) SelectCandidateInputs(TxId []byte, count uint32) []types.TxInput {
	inputs := make([]types.TxInput, count)

	rows, err := gSql.db.Query(`select prev_TxId, prev_OutIndex, Sequence, ScriptSize, Script
	 	from c_inputs where TxId = ?`, TxId)
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

func (gSql *GSqlite3) SelectCandidateOutputs(TxId []byte, count uint32) []types.TxOutput {
	outputs := make([]types.TxOutput, count)

	rows, err := gSql.db.Query(`select ToAddr, BrokerAddr, Script, ScriptSize, Type, Value from c_outputs 
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

func (gSql *GSqlite3) GetCandidateTxRows(rows *sql.Rows) []types.GhostTransaction {
	txs := []types.GhostTransaction{}
	for rows.Next() {
		tx := types.GhostTransaction{}
		if err := rows.Scan(&tx.TxId, &tx.Body.InputCounter, &tx.Body.OutputCounter,
			&tx.Body.Nonce, &tx.Body.LockTime); err != nil {
			log.Fatal(err)
		}
		tx.Body.Vin = gSql.SelectCandidateInputs(tx.TxId, tx.Body.InputCounter)
		tx.Body.Vout = gSql.SelectCandidateOutputs(tx.TxId, tx.Body.OutputCounter)
		txs = append(txs, tx)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return txs
}

func (gSql *GSqlite3) GetCandidateDataTxRows(rows *sql.Rows) []types.GhostDataTransaction {
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
