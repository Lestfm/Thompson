package tompson

import (
	"fmt"
	"time"
	"math/rand"
)

var bs = new(BetaSampler)

func init() {
	bs.New(time.Now().UnixNano())
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

type MachineName struct {
	X, Y int
}

type Input struct {
	Machines          []Machine
	TopMachines []MachineName
}

//Примитивная функция выбора лучшего результата для многоруких бандитов
//Оптимизировать функцию можно по различным критериям, тк значимость оптимизаций по руке A может быть больше,
//чем поруке B, но для демонстрации в этом нет смысла и по этому победитель выбирается случайным образом из показавших
//лучший результат на одной из рук. Для однорукого бандита соответствено это просто
//победитель.

func selectTop(m []MachineName, f []float64, highProbAv float64, mnAv MachineName) MachineName {
	win := rand.Intn(len(m) + 1)
	if win == len(m) {
		return mnAv
	}
	return m[win]
}

func (inp *Input) topMachine(outCount int, selectTop func(m []MachineName, f []float64, highProbAv float64, mnAv MachineName) MachineName) MachineName {
	inp.sample()
	//fmt.Printf("%+v\n",ii)
	highProb := make([]float64, outCount)
	mns := make([]MachineName, outCount)
	highProbAv := .0
	mnAv := MachineName{}
	j := 0
	i := 0
AGAIN:
	for i < len(inp.Machines) {
		mn := MachineName{0, 0}
		o := &inp.Machines[i][j]

		if o.Probs > highProb[j] {
			highProb[j] = o.Probs
			mn.X = i
			mn.Y = j
		}

		if j == 0 {
			var sumProb float64
			for _, e := range inp.Machines[i] {
				sumProb += e.Probs
			}
			if highProbAv < sumProb {
				highProbAv = sumProb
				mnAv.X = i
				mnAv.Y = j
			}
		}
		mns[j] = mn
		j++
		i++
		if j >= outCount {
			j = 0
			goto AGAIN
		}
	}
	//fmt.Printf(">>>%+v\n",mnAv)
	return selectTop(mns, highProb, highProbAv, mnAv)
}

func (inp *Input) sample() {
	for _, m := range inp.Machines {
		m.sample()
	}
}

func newInp(countMachine, outCount int) Input {
	inp := Input{}
	inp.Machines = make([]Machine, countMachine)
	for i := 0; i < countMachine; i++ {
		inp.Machines[i] = make([]Output, outCount)
		for j := 0; j < outCount; j++ {
			inp.Machines[i][j] = Output{0, 0, 1.0, 1.0, .0}
		}
	}
	return inp
}

type Machine []Output

func (m *Machine) sample() {
	for i := 0; i < len(*m); i++ {
		(*m)[i].sample()
	}
}
//Параметры alpha и beta можно было бы использовать при добавлении нового бандита к уже проходящим испытания
//Но в данный момент этого функционала нет
type Output struct {
	Succ, Fail int
	Alpha, Beta,
	Probs float64
}

func (o *Output) win() {
	o.Succ++
}

func (o *Output) lose() {
	o.Fail++
}

func (o *Output) sample() {
	sm := bs.Sample(float64(o.Succ)+o.Alpha, float64(o.Fail)+o.Beta)
	o.Probs = sm
}

type InputVec struct {
	Count int
}


//Возможно создать машину имеющую multi-input & multi-output.
//Из-за функция определения лучшей машины она будет вести себя несколько сюрреалистично
//Но логически верно, именно так как ожидается
func InitRoom(id string, outCount int, in []InputVec) (*Room, error) {
	r := new(Room)
	r.Id = id
	r.InCount = len(in)
	r.OutCount = outCount
	r.In = make([]Input, len(in))

	if outCount == 0 || len(in) == 0 {
		return nil, fmt.Errorf("error initialze input or output == 0")
	}

	for k, v := range in {
		r.In[k] = newInp(v.Count, outCount)
	}

	return r, nil
}

type Room struct {
	Id       string
	OutCount int
	InCount  int
	In       []Input
}

func (r *Room) Win(input, numberMach, NumberOutput int) {
	r.In[input].Machines[numberMach][NumberOutput].win()
	if r := recover(); r != nil {
		fmt.Println("Recovered:",r)
	}
}

func (r *Room) Lose(input, numberMach, NumberOutput int) {
	r.In[input].Machines[numberMach][NumberOutput].lose()
	if r := recover(); r != nil {
		fmt.Println("Recovered:", r)
	}
}

//Позиция элемента в возвращенном массиве соответствует номеру input
//значение номеру машины
func (r *Room) Game() []int {
	res := make([]int, len(r.In))
	for i, input := range r.In {
		res[i] = input.topMachine(r.OutCount, selectTop).X
	}
	return res
}

func (r *Room) Results() {
	for i := range r.In {
		fmt.Println("")
		for j := range r.In[i].Machines {
			for k := range r.In[i].Machines[j] {
				fmt.Printf("input: %d, machine %d, output %d: succes -> %d, fail -> %d\n", i, j, k, r.In[i].Machines[j][k].Succ, r.In[i].Machines[j][k].Fail)
			}
		}
	}
}

