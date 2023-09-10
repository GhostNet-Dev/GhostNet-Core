package sqlite

import (
	"database/sql"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

func InitTest(t *testing.T) *GSqlite3 {
	gSql := NewGSqlite3()
	err := gSql.OpenSQL("../../../", "block.db")
	if err != nil {
		t.Error(err)
	}
	return gSql
}

func TestSqlLatest(t *testing.T) {
	gSql := InitTest(t)
	defer gSql.CloseSQL()

	outputs := outputQuery(t, `select outputs.TxId, outputs.ToAddr, outputs.BrokerAddr, outputs.Script, outputs.ScriptSize, 
		outputs.ScriptEx, outputs.ScriptExSize, outputs.Type, outputs.Value, outputs.OutputIndex from outputs 
		where outputs.Type = ?
		order by outputs.BlockId DESC`, types.TxTypeScriptStore)
	if len(outputs) == 0 {
		t.Error("there is no output")
	}
	addr := outputs[0]

	newOutputs := gSql.SelectOutputLatests(types.TxTypeScriptStore, addr.Vout.Addr,
		[]byte("member"), 0, 2)
	if len(newOutputs) == 0 {
		t.Error("there is no output = ", string(addr.Vout.ScriptPubKey))
	}
}
func TestSqlSearch(t *testing.T) {
	gSql := InitTest(t)

	outputs := outputQuery(t, `select outputs.TxId, outputs.ToAddr, outputs.BrokerAddr, outputs.Script, outputs.ScriptSize, 
		outputs.ScriptEx, outputs.ScriptExSize, outputs.Type, outputs.Value, outputs.OutputIndex from outputs 
		where outputs.Type = ?
		order by outputs.BlockId DESC`, types.TxTypeScriptStore)
	if len(outputs) == 0 {
		t.Error("there is no output")
	}
	addr := outputs[0]
	length := len(addr.Vout.ScriptPubKey)
	searchText := addr.Vout.ScriptPubKey[length-2 : length]
	//searchText = append(searchText, '%')

	newOutputs := gSql.SearchStringOutputs(types.TxTypeScriptStore, addr.Vout.Addr,
		searchText)
	if len(newOutputs) == 0 {
		t.Error("there is no output", string(searchText), " = ", string(addr.Vout.ScriptPubKey))
	}
}

func outputQuery(t *testing.T, query string, txType types.TxOutputType) []types.PrevOutputParam {
	outputs := []types.PrevOutputParam{}
	rows, err := GSqlite.db.Query(query, txType)
	if err != nil {
		t.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		output := types.PrevOutputParam{}
		if err = rows.Scan(&output.VOutPoint.TxId, &output.Vout.Addr, &output.Vout.BrokerAddr, &output.Vout.ScriptPubKey,
			&output.Vout.ScriptSize, &output.Vout.ScriptEx, &output.Vout.ScriptExSize, &output.Vout.Type, &output.Vout.Value, &output.VOutPoint.TxOutIndex); err == sql.ErrNoRows {
			t.Error(err)
		} else if err != nil {
			t.Error(err)
		}
		output.TxType = txType
		outputs = append(outputs, output)
	}

	if err = rows.Err(); err != nil {
		t.Error(err)
	}
	return outputs
}
