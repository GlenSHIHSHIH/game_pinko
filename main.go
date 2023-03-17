package main

import (
	"math/rand"
	"mini/game/service"

	// _ "mini/game/model"
	"time"
)

var (
	ball map[int]int
)

type BallOdds map[int]float32

func main() {

	// for key, value := range service.GetPlinkoList() {
	// 	fmt.Printf("%v:%v \n", key, value)

	// }

	data := service.GetPlinkoList()
	plinko := service.NewPlinko(1, 8, 10)
	plinko.Calculate(data)

}

func CheckCount(row int, odds *BallOdds) {
	count := len(*odds)

	if count == row {
		return
	}

}

func RandomByNumber(number int) int {
	// 设置随机种子，确保每次运行都产生不同的随机数
	rand.Seed(time.Now().UnixNano())

	// 生成一个介于0和100之间的随机整数
	randomInt := rand.Intn(number + 1)

	return randomInt
}
