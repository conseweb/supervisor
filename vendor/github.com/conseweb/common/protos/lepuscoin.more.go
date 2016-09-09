/*
Copyright Mojing Inc. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package protos

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/golang/protobuf/proto"
)

// NewTx returns a new tx message that conforms to the Message
// interface.  The return instance has a default version of TxVersion and there
// are no transaction inputs or outputs.  Also, the lock time is set to zero
// to indicate the transaction is valid immediately as opposed to some time in
// future.
func NewTx() *TX {
	return &TX{
		Version: 1,
		Txin:    make([]*TX_TXIN, 0, 5),
		Txout:   make([]*TX_TXOUT, 0, 5),
	}
}

// NewTxIn returns a new transaction input with the provided
// previous outpoint point and signature script with a default sequence of
// MaxTxInSequenceNum.
func NewTxIn(prevHash []byte, prevIdx uint32, signatureScript []byte) *TX_TXIN {
	return &TX_TXIN{
		SourceHash: prevHash,
		Ix:         prevIdx,
		Script:     signatureScript,
		Sequence:   0xffffffff,
	}
}

// NewTxOut returns a new transaction output with the provided
// transaction value and public key script.
func NewTxOut(value uint64, pkScript []byte) *TX_TXOUT {
	return &TX_TXOUT{
		Value:  value,
		Script: pkScript,
	}
}

func (tx *TX) Base64Bytes() ([]byte, error) {
	txBytes, err := tx.Bytes()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, base64.StdEncoding.EncodedLen(len(txBytes)))
	base64.StdEncoding.Encode(buf, txBytes)

	return buf, nil
}

func (tx *TX) Bytes() ([]byte, error) {
	return proto.Marshal(tx)
}

// AddTxIn adds a transaction input to the message.
func (tx *TX) AddTxIn(ti *TX_TXIN) {
	tx.Txin = append(tx.Txin, ti)
}

// AddTxOut adds a transaction output to the message.
func (msg *TX) AddTxOut(to *TX_TXOUT) {
	msg.Txout = append(msg.Txout, to)
}

// TxHash generates the Hash for the transaction.
func (msg *TX) TxHash() []byte {
	txBytes, err := msg.Bytes()
	if err != nil {
		return nil
	}

	fHash := sha256.Sum256(txBytes)
	lHash := sha256.Sum256(fHash[:])
	txHash := make([]byte, 0)
	return append(txHash, lHash[:]...)
}

// Copy creates a deep copy of a transaction so that the original does not get
// modified when the copy is manipulated.
func (msg *TX) Copy() *TX {
	// Create new tx and start by copying primitive values and making space
	// for the transaction inputs and outputs.
	newTx := TX{
		Version:  msg.Version,
		Txin:     make([]*TX_TXIN, 0, len(msg.Txin)),
		Txout:    make([]*TX_TXOUT, 0, len(msg.Txout)),
		LockTime: msg.LockTime,
		Blocks:   msg.Blocks,
		Fee:      msg.Fee,
	}

	// Deep copy the old TxIn data.
	for _, oldTxIn := range msg.Txin {
		// Deep copy the old signature script.
		var newScript []byte
		oldScript := oldTxIn.Script
		oldScriptLen := len(oldScript)
		if oldScriptLen > 0 {
			newScript = make([]byte, oldScriptLen, oldScriptLen)
			copy(newScript, oldScript[:oldScriptLen])
		}

		// Create new txIn with the deep copied data and append it to
		// new Tx.
		newTxIn := &TX_TXIN{
			SourceHash: oldTxIn.SourceHash,
			Ix:         oldTxIn.Ix,
			Script:     newScript,
			Sequence:   oldTxIn.Sequence,
		}
		newTx.AddTxIn(newTxIn)
	}

	// Deep copy the old TxOut data.
	for _, oldTxOut := range msg.Txout {
		// Deep copy the old PkScript
		var newScript []byte
		oldScript := oldTxOut.Script
		oldScriptLen := len(oldScript)
		if oldScriptLen > 0 {
			newScript = make([]byte, oldScriptLen, oldScriptLen)
			copy(newScript, oldScript[:oldScriptLen])
		}

		// Create new txOut with the deep copied data and append it to
		// new Tx.
		newTxOut := &TX_TXOUT{
			Value:    oldTxOut.Value,
			Script:   newScript,
			Color:    oldTxOut.Color,
			Quantity: oldTxOut.Quantity,
			Addr:     oldTxOut.Addr,
		}
		newTx.AddTxOut(newTxOut)
	}

	return &newTx
}

func (e *ExecResult) Bytes() ([]byte, error) {
	return proto.Marshal(e)
}

func (r *QueryAddrResult) Bytes() ([]byte, error) {
	return proto.Marshal(r)
}

// ParseTXBytes unmarshal txData into TX object
func ParseTXBytes(txData []byte) (*TX, error) {
	tx := new(TX)
	err := proto.Unmarshal(txData, tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
