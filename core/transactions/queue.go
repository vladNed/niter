package transactions


type TxPool struct {
	transactions map[uint8]Tx
	queue        []uint8
	index        uint8
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[uint8]Tx),
		queue:        make([]uint8, 0),
		index:        0,
	}
}

func (tp *TxPool) Add(tx Tx) uint8 {
	tp.index++
	tp.transactions[tp.index] = tx
	tp.queue = append(tp.queue, tp.index)

	return tp.index
}

func (tp *TxPool) Next() Tx {
	if len(tp.queue) == 0 {
		return nil
	}
	tx := tp.transactions[tp.queue[0]]
	return tx
}

func (tp *TxPool) Mark() {
	if len(tp.queue) == 1 {
		tp.queue = make([]uint8, 0)
		return
	}
	tp.queue = tp.queue[1:]
}


func RunTxPoolHandler(tp *TxPool, tpChan chan Tx) {
	for {
		newTx := <-tpChan
		tp.Add(newTx)
	}
}