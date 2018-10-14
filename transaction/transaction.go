package transaction

type Transaction struct {
	Amount    uint64 `json:"amount"`
	Recipient string `json:"recipient"`
	Sender    string `json:"sender"`
	Data      []byte `json:"data"`
}

type TxPool struct {
	AllTx []Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		AllTx: make([]Transaction, 0),
	}
}

func (p *TxPool) Clear() bool {
	if len(p.AllTx) == 0 {
		return true
	}
	p.AllTx = make([]Transaction, 0)
	return true
}

func (p *TxPool) AddTx(tx *Transaction) int {
	p.AllTx = append(p.AllTx, *tx)
	return len(p.AllTx)
}
