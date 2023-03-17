package test

import (
	"fmt"
	"mini/game/service"
	"sync"
	"testing"

	"github.com/shopspring/decimal"
)

// get db data
func TestPlinko(t *testing.T) {
	data := service.GetPlinkoList()

	for key, value := range data {
		fmt.Printf("%v:%v \n", key, value)
	}

	plinko := service.NewPlinko(1, 9, 1)
	arryCount := plinko.Row + 1
	middleNumber := int((arryCount - 1) / 2)
	odds := plinko.GetPlinkoBytypeAndRow(arryCount, data)
	dataSort := service.GetPlinkoDistrByArray(arryCount, middleNumber, odds)

	for key, value := range dataSort {
		fmt.Printf("odd[%d] = %.4f\n", key, value)
	}

	t.Log("success")
}

// get 機率(P) data
func TestPlinkoCalculate(t *testing.T) {
	intervalRate := service.GetPlinkoProbability(9)
	for key, value := range intervalRate {
		fmt.Printf("P(X = %d) = %.20f\n", key, value)
	}

	t.Log("success")
}

func TestPlinkoRandom(t *testing.T) {
	data := service.GetPlinkoList()

	count := 10

	c := NewCustomer()
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			plinko := service.NewPlinko(1, 9, 1)
			p := plinko.Calculate(data)
			c.increment(p)

			fmt.Printf("%+v \n", p)
			defer wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("%+v \n", c)

}

type customer struct {
	SpendMoney float64
	// GenMoney   float64
	TotalMoney float64
}

func NewCustomer() *customer {
	return &customer{
		// SpendMoney: 0,
		// TotalMoney: 0,
	}
}

var lock sync.Mutex

func (c *customer) increment(p *service.Plinko) {
	lock.Lock()
	defer lock.Unlock()
	c.SpendMoney, _ = decimal.NewFromFloat(p.Money).Add(decimal.NewFromFloat(c.SpendMoney)).Float64()
	c.TotalMoney, _ = decimal.NewFromFloat(p.Money * p.Odd).Add(decimal.NewFromFloat(c.TotalMoney)).Float64()
}
