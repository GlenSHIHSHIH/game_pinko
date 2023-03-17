package service

import (
	"math"
	"math/rand"
	"mini/game/model"
	"sort"
	"time"

	"github.com/shopspring/decimal"
)

type Plinko struct {
	Type    int
	Row     int
	Money   float64
	Account float64
	Odd     float64
}

// (t)type:球種 , (r)row:階層行數 , (m)money:投注金額
func NewPlinko(t, r int, m float64) *Plinko {
	return &Plinko{
		Type:  t,
		Row:   r,
		Money: m,
	}
}

func GetPlinkoList() []map[string]any {
	var data []map[string]any
	db := model.GetMySqlDB()
	db.Model(&model.PlinkoBalls{}).Select("*").
		Joins("INNER JOIN plinko_odds ON plinko_odds.plinko_balls_id = plinko_balls.id").
		Scan(&data)
	return data
}

func (p *Plinko) Calculate(data []map[string]any) *Plinko {

	arryCount := p.Row + 1
	odds := p.GetPlinkoBytypeAndRow(arryCount, data)

	//排序
	sort.Float64s(odds)

	//求中間數值
	middleNumber := int((arryCount - 1) / 2)

	plinkoSort := GetPlinkoDistrByArray(arryCount, middleNumber, odds)

	intervalRate := GetPlinkoProbability(p.Row)

	probIndex := GetPlinkoProbIndex(arryCount, plinkoSort, intervalRate)

	p.Account, _ = decimal.NewFromFloat(p.Money).Mul(decimal.NewFromFloat(plinkoSort[probIndex])).Float64()
	// p.Account = p.Money * plinkoSort[probIndex]
	p.Odd = plinkoSort[probIndex]

	return p
}

func (p *Plinko) GetPlinkoBytypeAndRow(arryCount int, data []map[string]any) []float64 {
	odds := make([]float64, 0, arryCount)
	index := 0
	//取值
	for _, value := range data {
		if (int(value["type"].(int64)) == p.Type) &&
			(int(value["row"].(int64)) == p.Row) {
			valueOdds := float64(value["odd"].(float64))
			// fmt.Println(valueOdds)
			odds = append(odds, valueOdds)
			// odds[index] = float64(valueOdds)
			index++
		}
	}
	return odds
}

// 計算機率落於的區間值 回傳index
func GetPlinkoProbIndex(count int, plinkoSort, intervalRate []float64) int {

	times := time.Now().UnixNano()
	rand.Seed(times)
	// rand.Seed(time.Now().UnixNano())
	random := decimal.NewFromFloat(rand.Float64())

	totalProb := decimal.NewFromFloat(0)
	// fmt.Println("random:", random)
	for i := 0; i < count; i++ {
		totalProb = totalProb.Add(decimal.NewFromFloat(intervalRate[i]))
		// tP, _ := totalProb.Float64()
		// fmt.Printf("%.10f \n", intervalRate[i])
		// fmt.Printf("sum[%v] %.10f \n", i, tP)
		if random.LessThanOrEqual(totalProb) {
			return i
		}
	}
	return int((count - 1) / 2)
}

// Plinko 的排序從中間向外 由小到大
func GetPlinkoDistrByArray(count, middleNumber int, odds []float64) []float64 {
	data := make([]float64, len(odds))
	decrease := 0
	increment := 1
	for i := 0; i < count; i++ {
		if i%2 == 0 {
			data[middleNumber-(decrease)] = odds[i]
			decrease++
		} else {
			data[middleNumber+increment] = odds[i]
			increment++
		}
	}
	return data
}

// 計算區間機率
func GetPlinkoProbability(n int) []float64 {

	// 定義實驗次數n和成功概率p
	// n := 8
	p := 0.5

	arryCount := n + 1
	prob := make([]float64, 0, arryCount)

	// 計算成功次數為0到n的概率
	for k := 0; k <= n; k++ {
		value := binomialProb(n, p, k)
		prob = append(prob, value)
		// fmt.Printf("P(X = %d) = %.6f\n", k, value)
	}
	return prob
}

// 計算二項機率分布的概率質量函數
func binomialProb(n int, p float64, k int) float64 {
	// 計算組合數
	c := float64(comb(n, k))

	// 計算概率質量函數
	return c * math.Pow(p, float64(k)) * math.Pow(1-p, float64(n-k))
}

// 計算組合數
func comb(n, k int) int {
	if k > n {
		return 0
	}
	if k == 0 || k == n {
		return 1
	}
	return comb(n-1, k-1) + comb(n-1, k)
}
