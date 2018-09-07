package tompson

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"fmt"
)

type actioner interface {
	action(s *Storage) []byte
}

type Router struct {
	router *gin.Engine
}

func (r *Router) ListenAndServe(port string) {
	r.router.Run(port)
}

func errorToByte(err error) []byte {
	e, _ := json.Marshal(&struct {
		Error string `json:"error"`
	}{
		fmt.Sprint(err),
	})
	return e
}

func newHandler(obj actioner) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		data, err := ctx.GetRawData()
		if err != nil {
			ctx.JSON(503, err)
		}
		json.Unmarshal(data, obj)
		storage, ok := ctx.Get("Storage")
		if !ok {
			panic("Отсутствует хранилище игровых комнат!")
		}
		b := obj.action(storage.(*Storage))
		ctx.String(200, string(b))
	}
}

//Создает новый роутер с подготовленной таблицей маршрутизации
func NewRouter(storage *Storage) *Router {
	router := gin.Default()
	router.Use(func(ctx *gin.Context) {
		ctx.Set("Storage", storage)
	})
	router.POST("/init", newHandler(&JsonRoom{}))
	router.POST("/win", newHandler(&Win{}))
	router.POST("/loose", newHandler(&Loose{}))
	router.POST("/game", newHandler(&JsonGame{}))

	router.POST("/save", func(context *gin.Context) {
	})

	return &Router{
		router: router,
	}
}

type JsonAction struct {
	Method string `json:"method"`
}

type JsonRoom struct {
	ID       string `json:"id"`
	OutCount int    `json:"out_count"`
	Inputs   []int  `json:"inputs"`
}

func (jr *JsonRoom) action(s *Storage) []byte {
	iv := make([]InputVec, len(jr.Inputs))

	for c, i := range jr.Inputs {
		iv[c] = InputVec{i}
	}
	_, err := s.Create(jr.ID, jr.OutCount, iv)
	if err != nil {
		return errorToByte(err)
	}
	return []byte("ok")
}

// Для Win и Loose прие игре в обычного однорукого бандита с несколькими слотами
// Значения intput и output должны быть равны 0
type Win struct {
	ID      string `json:"id"`
	Input   int    `json:"input"`
	Machine int    `json:"machine"`
	Output  int    `json:"output"`
}

type Loose struct {
	ID      string `json:"id"`
	Input   int    `json:"input"`
	Machine int    `json:"machine"`
	Output  int    `json:"output"`
}

func (w *Win) action(s *Storage) []byte {
	room := s.Get(w.ID)
	if room == nil {
		return errorToByte(fmt.Errorf("игровая комната не загружена"))
	}
	room.Win(w.Input, w.Machine, w.Output)
	return []byte("ok")
}

func (l *Loose) action(s *Storage) []byte {
	room := s.Get(l.ID)
	if room == nil {
		return errorToByte(fmt.Errorf("игровая комната не загружена"))
	}
	room.Lose(l.Input, l.Machine, l.Output)
	return []byte("ok")
}

type JsonGame struct {
	ID string `json:"id"`
}

func (jg *JsonGame) action(s *Storage) []byte {
	room := s.Get(jg.ID)
	if room == nil {
		return errorToByte(fmt.Errorf("игровая комната не загружена"))
	}
	res := room.Game()
	byte, _ := json.Marshal(&struct {
		Slots []int
	}{res})
	return byte
}
