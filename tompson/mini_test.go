package tompson

import (
	"fmt"
	"net/http"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"testing"
	"time"
	"math/rand"
)

type mockDB struct{}

func (mockDB) Put([]byte, []byte) error {
	return nil
}

func (mockDB) Get([]byte) ([]byte, error) {
	return []byte{}, nil
}

func (mockDB) Delete([]byte) error {
	return nil
}
func init() {
	go startServer()
	time.Sleep(1 * time.Second)
}

var storage *Storage

func startServer() {
	db := &mockDB{}
	storage = NewStorage(db)
	router := NewRouter(storage)
	router.ListenAndServe(":8081")
}

type Slots struct {
	Slots []int `json:"Slots"`
}

func createRoom() {
	b, _ := json.Marshal(&JsonRoom{
		"firstGame",
		1,
		[]int{2},
	})
	reader := bytes.NewReader(b)
	resp, err := http.Post("http://127.0.0.1:8081/init", "", reader)
	if err != nil {
		fmt.Print(err)
		return
	}

	data, _ := ioutil.ReadAll(resp.Body)

	fmt.Print(string(data))
}

func TestListen(t *testing.T) {
	createRoom()
	//time.Sleep(time.Second * 1)
	toGame()
}

func toGame() {
	gameQuery := func() int {
		b, _ := json.Marshal(&JsonGame{
			"firstGame",
		})
		reader := bytes.NewReader(b)
		resp, err := http.Post("http://127.0.0.1:8081/game", "", reader)
		if err != nil {
			fmt.Print(err)
		}
		data, _ := ioutil.ReadAll(resp.Body)
		slots := &Slots{}
		err = json.Unmarshal(data, slots)
		return slots.Slots[0]
	}
	winQuery := func(i int) {
		b, _ := json.Marshal(&Win{
			ID:      "firstGame",
			Machine: i,
		})
		reader := bytes.NewReader(b)
		http.Post("http://127.0.0.1:8081/win", "", reader)
	}
	looseQuery := func(i int) {
		b, _ := json.Marshal(&Win{
			ID:      "firstGame",
			Machine: i,
			Input:   0,
			Output:  0,
		})
		reader := bytes.NewReader(b)
		http.Post("http://127.0.0.1:8081/loose", "", reader)
	}
	firstWinRate := 0.1
	secondWinRate := 0.15
	rand.Seed(42)
	firtGamesCount, secondGamesCount := 0, 0

	for i := 0; i < 6000; i++ {
		p := rand.Float64()
		machine := gameQuery()
		if machine == 0 {
			firtGamesCount++
			if p < firstWinRate {
				winQuery(0)
			} else {
				looseQuery(0)
			}
		} else {
			secondGamesCount++
			if p < secondWinRate {
				winQuery(1)
			} else {
				looseQuery(1)
			}
		}
	}
	fmt.Printf("Первая машина сыграла %d игр\n Вторая сыграла %d игр", firtGamesCount, secondGamesCount)
	room := storage.Get("firstGame")
	room.Results()
	}
