package tompson

import (
	"math/rand"
	"math"
)
/*
	Реализация взятия случайного числа из бета распределения
*/
type BetaSampler struct {
	R *rand.Rand
}

func (bs *BetaSampler) New(seed int64) {
	bs.R = rand.New(rand.NewSource(seed))
}

func (bs *BetaSampler) Sample(a, b float64) float64{
	alpha := a + b
	beta := .0
	u1, u2, w, v := .0, .0, .0, .0
	if math.Min(a, b) <= 1.0 {
		beta = math.Max(1 / a, 1 / b)
	} else {
		beta = math.Sqrt((alpha - 2) / (2 * a * b - alpha))
	}
	gamma := a + 1 / beta
	for true {
		u1 = bs.R.Float64()
		u2 = bs.R.Float64()
		v = beta * math.Log(u1 / (1 - u1))
		w = a * math.Exp(v)
		tmp := math.Log(alpha / (b + w))
		if alpha * tmp + (gamma * v) - 1.3862944 >= math.Log(u1 * u1 * u2){
			break
		}
	}
	return w / (b + w)
}
