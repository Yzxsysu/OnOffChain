package application

import (
	"math"
	"math/rand"
	"time"
)

type ZipfDistribution struct {
	Size       float64
	Skew       float64
	AccountNum float64
}

func NewZipfDistribution(size float64, skew float64, AccountNum float64) *ZipfDistribution {
	return &ZipfDistribution{
		size,
		skew,
		AccountNum,
	}
}

func H(n float64, s float64) float64 {
	if n == 1 {
		return 1.0 / math.Pow(n, s)
	} else {
		return (1.0 / math.Pow(n, s)) + H(n-1, s)
	}
}

func (z *ZipfDistribution) F(k float64) float64 {
	return (1 / math.Pow(k, z.Skew) / H(z.Size, z.Skew))
}

func (z *ZipfDistribution) Cdf(k float64) float64 {
	return H(k, z.Skew) / H(z.Size, z.Skew)
}

func RandWithout(start, end, without int) int {
	if end-start+1 <= 1 {
		return start
	}

	result := rand.Intn(end-start+1-1) + start
	if result >= without {
		result++
	}
	return result
}

func (z *ZipfDistribution) Uint64() ([]int, []int) {
	rand.Seed(time.Now().UnixNano())
	size := int(z.Size)
	from := make([]int, size)
	to := make([]int, size)
	var TxFromAccount int = 1

	count := 0
	for i := 1; i <= size; i++ {
		EachAccountFrequency := math.Ceil(z.F(float64(TxFromAccount)) * z.Size)

		if TxFromAccount == int(z.AccountNum) {
			TxFromAccount = 1
		}
		for j := 0; j < int(EachAccountFrequency); j++ {
			from[count] = TxFromAccount
			count++
			if count == size {
				break
			}
		}
		if count == size {
			break
		}
		TxFromAccount++
	}
	//fmt.Println("Original slice:", from)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(from), func(i, j int) {
		from[i], from[j] = from[j], from[i]
	})
	for i := 0; i < len(from); i++ {
		excludeNum := RandWithout(1, int(z.Size), from[i])
		to[i] = excludeNum
		for from[i] == to[i] {
			to[i] = RandWithout(1, int(z.Size), from[i])
		}
		if from[i] > to[i] {
			from[i], to[i] = to[i], from[i]
		}
	}
	return from, to
}

func GenerateTx(TxNum float64, TxAccountNum float64, Skew float64) []SmallBankTransaction {
	txList := make([]SmallBankTransaction, int(TxNum))
	z := NewZipfDistribution(TxNum, Skew, TxAccountNum)
	from, to := z.Uint64()
	var s SmallBankTransaction
	for i := 0; i < int(TxNum); i++ {
		s.F = IntToBytes(from[i])
		s.O = IntToBytes(to[i])
		rand.Seed(time.Now().UnixNano())
		TxType := rand.Intn(6) + 1
		s.T = uint8(TxType)
		Balance := rand.Intn(1000) + 1
		s.B = Balance
		s.I = uint16(i + 1)
		txList[i] = s
	}
	return txList
}

/*func main() {
	//m := make(map[int]int)
	tx := NewTxResult()
	fmt.Println(unsafe.Sizeof(tx))
	fmt.Println(tx)
	q := "1"

	a, _ := json.Marshal([]byte(q))
	fmt.Println(len(a))

	acclock := make(map[string]*Lock)
	acclock["1"] = NewLock()
	acclock["1"].lock.RLock()
	defer acclock["1"].lock.RUnlock()

	//smtx := GenerateTx(1000, 1000, 1)
	//fmt.Println(smtx)
	//BCstate := NewBlockchainState("cleveldb,", true, "nil")
	//for i := 1; i <= 1000; i++ {
	//	BCstate.CreateAccount(string(i), 1000, 1000)
	//}
	//BCstate.ExecuteSmallBankTransaction(smtx, 10)
}*/
