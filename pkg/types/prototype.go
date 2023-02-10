package types

import "github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"

func GhostDataTxToProtoType(d *GhostDataTransaction) *ptypes.GhostDataTransaction {
	return &ptypes.GhostDataTransaction{
		TxId:           d.TxId,
		LogicalAddress: d.LogicalAddress,
		DataSize:       d.DataSize,
		Data:           d.Data,
	}
}

func ProtoTypeToGhostDataTx(d *ptypes.GhostDataTransaction) *GhostDataTransaction {
	return &GhostDataTransaction{
		TxId:           d.TxId,
		LogicalAddress: d.LogicalAddress,
		DataSize:       d.DataSize,
		Data:           d.Data,
	}
}

func GhostTxOutputToProtoType(o *TxOutput) *ptypes.TxOutput {
	return &ptypes.TxOutput{
		Addr:         o.Addr,
		BrokerAddr:   o.BrokerAddr,
		Type:         uint64(o.Type),
		Value:        o.Value,
		ScriptSize:   o.ScriptSize,
		ScriptPubKey: o.ScriptPubKey,
	}
}

func ProtoTypeToGhostTxOutput(o *ptypes.TxOutput) *TxOutput {
	return &TxOutput{
		Addr:         o.Addr,
		BrokerAddr:   o.BrokerAddr,
		Type:         TxOutputType(o.Type),
		Value:        o.Value,
		ScriptSize:   o.ScriptSize,
		ScriptPubKey: o.ScriptPubKey,
	}
}

func GhostTxInputToProtoType(in *TxInput) *ptypes.TxInput {
	return &ptypes.TxInput{
		PrevOut: &ptypes.TxOutPoint{
			TxId:       in.PrevOut.TxId,
			TxOutIndex: in.PrevOut.TxOutIndex,
		},
		Sequence:   in.Sequence,
		ScriptSize: in.ScriptSize,
		ScriptSig:  in.ScriptSig,
	}
}

func ProtoTypeToGhostTxInput(in *ptypes.TxInput) *TxInput {
	return &TxInput{
		PrevOut: TxOutPoint{
			TxId:       in.PrevOut.TxId,
			TxOutIndex: in.PrevOut.TxOutIndex,
		},
		Sequence:   in.Sequence,
		ScriptSize: in.ScriptSize,
		ScriptSig:  in.ScriptSig,
	}
}

func GhostTxToProtoType(t *GhostTransaction) *ptypes.GhostTransaction {
	var txInput []*ptypes.TxInput
	var txOutput []*ptypes.TxOutput
	for _, in := range t.Body.Vin {
		txInput = append(txInput, GhostTxInputToProtoType(&in))
	}
	for _, out := range t.Body.Vout {
		txOutput = append(txOutput, GhostTxOutputToProtoType(&out))
	}
	return &ptypes.GhostTransaction{
		TxId: t.TxId,
		Body: &ptypes.TxBody{
			InputCounter:  t.Body.InputCounter,
			Vin:           txInput,
			OutputCounter: t.Body.OutputCounter,
			Vout:          txOutput,
			Nonce:         t.Body.Nonce,
			LockTime:      t.Body.LockTime,
		},
	}
}

func ProtoTypeToGhostTx(t *ptypes.GhostTransaction) *GhostTransaction {
	var txInput []TxInput
	var txOutput []TxOutput
	for _, in := range t.Body.Vin {
		txInput = append(txInput, *ProtoTypeToGhostTxInput(in))
	}
	for _, out := range t.Body.Vout {
		txOutput = append(txOutput, *ProtoTypeToGhostTxOutput(out))
	}
	return &GhostTransaction{
		TxId: t.TxId,
		Body: TxBody{
			InputCounter:  t.Body.InputCounter,
			Vin:           txInput,
			OutputCounter: t.Body.OutputCounter,
			Vout:          txOutput,
			Nonce:         t.Body.Nonce,
			LockTime:      t.Body.LockTime,
		},
	}
}

func GhostBlockToProtoType(p *PairedBlock) *ptypes.PairedBlocks {
	h := p.Block.Header
	dh := p.DataBlock.Header

	var atxs []*ptypes.GhostTransaction
	var txs []*ptypes.GhostTransaction
	var dataTxs []*ptypes.GhostDataTransaction

	for _, tx := range p.Block.Alice {
		atxs = append(atxs, GhostTxToProtoType(&tx))
	}

	for _, tx := range p.Block.Transaction {
		txs = append(txs, GhostTxToProtoType(&tx))
	}

	for _, tx := range p.DataBlock.Transaction {
		dataTxs = append(dataTxs, GhostDataTxToProtoType(&tx))
	}

	return &ptypes.PairedBlocks{
		Block: &ptypes.GhostNetBlock{
			Header: &ptypes.GhostNetBlockHeader{
				Id:                      h.Id,
				Version:                 h.Version,
				PreviousBlockHeaderHash: h.PreviousBlockHeaderHash,
				MerkleRoot:              h.MerkleRoot,
				DataBlockHeaderHash:     h.DataBlockHeaderHash,
				TimeStamp:               h.TimeStamp,
				Bits:                    h.Bits,
				Nonce:                   h.Nonce,
				AliceCount:              h.AliceCount,
				TransactionCount:        h.TransactionCount,
				SignatureSize:           h.SignatureSize,
				BlockSignature:          h.BlockSignature.SerializeToByte(),
			},
			Alice:       atxs,
			Transaction: txs,
		},
		DataBlock: &ptypes.GhostNetDataBlock{
			Header: &ptypes.GhostNetDataBlockHeader{
				Id:                      dh.Id,
				Version:                 dh.Version,
				PreviousBlockHeaderHash: dh.PreviousBlockHeaderHash,
				MerkleRoot:              dh.MerkleRoot,
				Nonce:                   dh.Nonce,
				TransactionCount:        dh.TransactionCount,
			},
			Transaction: dataTxs,
		},
	}
}

func ProtoTypeToGhostBlock(p *ptypes.PairedBlocks) *PairedBlock {
	h := p.Block.Header
	dh := p.DataBlock.Header
	sig := &SigHash{}
	sig.DeserializeSigHashFromByte(h.BlockSignature)

	var atxs []GhostTransaction
	var txs []GhostTransaction
	var dataTxs []GhostDataTransaction

	for _, tx := range p.Block.Alice {
		atxs = append(atxs, *ProtoTypeToGhostTx(tx))
	}

	for _, tx := range p.Block.Transaction {
		txs = append(txs, *ProtoTypeToGhostTx(tx))
	}

	for _, tx := range p.DataBlock.Transaction {
		dataTxs = append(dataTxs, *ProtoTypeToGhostDataTx(tx))
	}

	return &PairedBlock{
		Block: GhostNetBlock{
			Header: GhostNetBlockHeader{
				Id:                      h.Id,
				Version:                 h.Version,
				PreviousBlockHeaderHash: h.PreviousBlockHeaderHash,
				MerkleRoot:              h.MerkleRoot,
				DataBlockHeaderHash:     h.DataBlockHeaderHash,
				TimeStamp:               h.TimeStamp,
				Bits:                    h.Bits,
				Nonce:                   h.Nonce,
				AliceCount:              h.AliceCount,
				TransactionCount:        h.TransactionCount,
				SignatureSize:           h.SignatureSize,
				BlockSignature:          *sig,
			},
			Alice:       atxs,
			Transaction: txs,
		},
		DataBlock: GhostNetDataBlock{
			Header: GhostNetDataBlockHeader{
				Id:                      dh.Id,
				Version:                 dh.Version,
				PreviousBlockHeaderHash: dh.PreviousBlockHeaderHash,
				MerkleRoot:              dh.MerkleRoot,
				Nonce:                   dh.Nonce,
				TransactionCount:        dh.TransactionCount,
			},
			Transaction: dataTxs,
		},
	}
}
