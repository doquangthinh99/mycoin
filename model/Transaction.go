package model
// type TXOutput struct {
// 	Value      int
// 	PubKeyHash []byte
// }

// type TXOutputs struct {
// 	Outputs []TXOutput
// }

// func (out *TXOutput) Lock(address []byte) {
// 	pubKeyHash := Base58Decode(address)
// 	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
// 	out.PubKeyHash = pubKeyHash
// }

// func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
// 	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
// }

// func NewTXOutput(value int, address string) *TXOutput {
// 	txo := &TXOutput{value, nil}
// 	txo.Lock([]byte(address))

// 	return txo
// }

// type TXInput struct {
// 	Txid      []byte
// 	Vout      int
// 	Signature []byte
// 	PubKey    []byte
// }

// type Transaction struct {
// 	ID   []byte
// 	Vin  []TXInput
// 	Vout []TXOutput
// }

// func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
// 	lockingHash := HashPubKey(in.PubKey)

// 	return bytes.Compare(lockingHash, pubKeyHash) == 0
// }

// func NewCoinBaseTX(to, data string) *Transaction {
// 	if data == "" {
// 		randData := make([]byte, 20)
// 		_, err := rand.Read(randData)
// 		if err != nil {
// 			log.Panic(err)
// 		}

// 		data = fmt.Sprintf("%x", randData)
// 	}

// 	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
// 	txout := NewTXOutput(10, to)
// 	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
// 	tx.ID = tx.Hash()

// 	fmt.Println(tx)

// 	return &tx
// }

// func (tx *Transaction) SetID() {
// 	var encoded bytes.Buffer
// 	var hash [32]byte

// 	enc := gob.NewEncoder(&encoded)
// 	err := enc.Encode(tx)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	hash = sha256.Sum256(encoded.Bytes())
// 	tx.ID = hash[:]
// }

// func (tx Transaction) IsCoinBase() bool {
// 	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
// }

// func (tx Transaction) Serialize() []byte {
// 	var encoded bytes.Buffer

// 	enc := gob.NewEncoder(&encoded)
// 	err := enc.Encode(tx)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	return encoded.Bytes()
// }

// func (tx *Transaction) Hash() []byte {
// 	var hash [32]byte

// 	txCopy := *tx
// 	txCopy.ID = []byte{}

// 	hash = sha256.Sum256(txCopy.Serialize())

// 	return hash[:]
// }

// func NewUTXOTransaction(from, to string, amount int, UTXOSet *UTXOSet) *Transaction {
// 	var inputs []TXInput
// 	var outputs []TXOutput

// 	wallets, err := NewWallets()
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	wallet := wallets.GetWallet(from)
// 	pubKeyHash := HashPubKey(wallet.PublicKey)
// 	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)

// 	if acc < amount {
// 		log.Panic("ERROR: Not enough funds")
// 	}

// 	// Build a list of inputs
// 	for txid, outs := range validOutputs {
// 		txID, err := hex.DecodeString(txid)
// 		if err != nil {
// 			log.Panic(err)
// 		}

// 		for _, out := range outs {
// 			input := TXInput{txID, out, nil, wallet.PublicKey}
// 			inputs = append(inputs, input)
// 		}
// 	}

// 	// Build a list of outputs
// 	outputs = append(outputs, *NewTXOutput(amount, to))
// 	if acc > amount {
// 		outputs = append(outputs, *NewTXOutput(acc-amount, from)) // a change
// 	}

// 	tx := Transaction{nil, inputs, outputs}
// 	tx.ID = tx.Hash()
// 	UTXOSet.BlockChain.SignTransaction(&tx, wallet.PrivateKey)

// 	return &tx
// }

// func (tx *Transaction) TrimmedCopy() Transaction {
// 	var inputs []TXInput
// 	var outputs []TXOutput

// 	for _, vin := range tx.Vin {
// 		inputs = append(inputs, TXInput{vin.Txid, vin.Vout, nil, nil})
// 	}

// 	for _, vout := range tx.Vout {
// 		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
// 	}

// 	txCopy := Transaction{tx.ID, inputs, outputs}

// 	return txCopy
// }

// func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
// 	if tx.IsCoinBase() {
// 		return
// 	}

// 	for _, vin := range tx.Vin {
// 		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
// 			log.Panic("ERROR: Previous transaction is not correct")
// 		}
// 	}

// 	txCopy := tx.TrimmedCopy()

// 	for inID, vin := range txCopy.Vin {
// 		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
// 		txCopy.Vin[inID].Signature = nil
// 		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
// 		txCopy.ID = txCopy.Hash()
// 		txCopy.Vin[inID].PubKey = nil

// 		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
// 		if err != nil {
// 			log.Panic(err)
// 		}
// 		signature := append(r.Bytes(), s.Bytes()...)

// 		tx.Vin[inID].Signature = signature
// 	}
// }

// func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
// 	txCopy := tx.TrimmedCopy()
// 	curve := elliptic.P256()

// 	for inID, vin := range tx.Vin {
// 		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
// 		txCopy.Vin[inID].Signature = nil
// 		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
// 		txCopy.ID = txCopy.Hash()
// 		txCopy.Vin[inID].PubKey = nil

// 		r := big.Int{}
// 		s := big.Int{}
// 		sigLen := len(vin.Signature)
// 		r.SetBytes(vin.Signature[:(sigLen / 2)])
// 		s.SetBytes(vin.Signature[(sigLen / 2):])

// 		x := big.Int{}
// 		y := big.Int{}
// 		keyLen := len(vin.PubKey)
// 		x.SetBytes(vin.PubKey[:(keyLen / 2)])
// 		y.SetBytes(vin.PubKey[(keyLen / 2):])

// 		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
// 		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
// 			return false
// 		}
// 	}

// 	return true
// }

// func (outs TXOutputs) Serialize() []byte {
// 	var buff bytes.Buffer

// 	enc := gob.NewEncoder(&buff)
// 	err := enc.Encode(outs)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	return buff.Bytes()
// }

// func DeserializeOutputs(data []byte) TXOutputs {
// 	var outputs TXOutputs

// 	dec := gob.NewDecoder(bytes.NewReader(data))
// 	err := dec.Decode(&outputs)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	return outputs
// }

type Transaction struct{
	Sign 		[]byte
	Miner		[]byte
	Sender		[]byte
	Receiver	[]byte
	Amount		int
	Status		string
}

func NewTransaction(sender []byte,receiver []byte,miner []byte,amount int) Transaction{
	return Transaction{nil,nil,sender,receiver,amount,"Unconfirm"}
}

func CreateBaseTransaction() Transaction{
	return Transaction{nil,nil,nil,nil,0,"BaseTX"}
}

func (tx Transaction)SignTransaction(sign []byte) {
	tx.Sign = sign;
}