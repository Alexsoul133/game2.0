package main

import (
	"fmt"
	"time"

	"github.com/macroblock/garbage/conio"
	"github.com/macroblock/garbage/utils"
	"github.com/macroblock/zl/core/zlog"
)

var (
	log   = zlog.Instance("game")
	throw = log.Catcher()
)
var gameover = false
var gameMap TMap
var protoMap = ",,,, ,,,,,,,,,," +
	", h      ,,,,,," +
	",~F~     ,,,s,," +
	" ,~~   , ,,b,,," +
	",S,       ,,M,," +
	" C,,     ,,,,,," +
	" ,,  ,   ,,,,,," +
	" ,S,     ,,,,,," +
	" , ,     ,,,,,," +
	",,,,   , ,h,,,," +
	",, ,   , ,,,,,," +
	",, ,   , ,,,,,," +
	",, ,   , ,,,,,," +
	",, ,   , ,,,,,," +
	",, ,   , ,,,,,,"

type TMap struct {
	data [mapW * mapH]*TCell
	hero *THero
	npc  []IObject
	item []IItems
}

type TCell struct {
	ground IGround
}

func create(r byte, x, y int) {
	object := IObject(nil)
	ground := IGround(nil)
	item := IItems(nil)
	switch rune(r) {
	case 'S':
		object = newSpinner(x, y)
		ground = newGround(x, y)
	case 'h':
		object = newHuman(x, y)
		ground = newGround(x, y)
	case 'C':
		object = newCat(x, y)
		ground = newGrass(x, y)
	case 'P':
		object = newPlant(x, y)
		ground = newGrass(x, y)
	case 'F':
		object = newPiranha(x, y)
		ground = newWater(x, y)
	case 'M':
		gameMap.hero = newHero(x, y)
		ground = newGround(x, y)
	case ' ':
		ground = newGround(x, y)

	case '#':
		ground = newWall(x, y)

	case ',':
		ground = newGrass(x, y)

	case '~':
		ground = newWater(x, y)

	case 's':
		ground = newGround(x, y)
		item = newSword(x, y)

	case 'b':
		ground = newGround(x, y)
		item = newBoots(x, y)
	} //end
	if object != nil {
		addNpc(object)
	}
	if item != nil {
		addItem(item)
	}
	if ground != nil {
		placeGround(ground)
		return
	}
	panic(fmt.Sprintf("Unknown object %q", r))
}

// func createGround(r byte, x, y int) {
// 	switch rune(r) {

// 	}
// 	panic(fmt.Sprintf("Unknown ground %q", r))
// }

func newMap() {
	for i := 0; i < mapH*mapW; i++ {
		create(protoMap[i], i%mapW, i/mapW)
	}
}

////////////////////////////////////////////////////
func ai() {

	// gameMap.hero.Do()
	for _, npc := range gameMap.npc {
		npc.Do()
	}
	for _, cell := range gameMap.data {
		cell.ground.Do()
	}

	gameover = true
	for _, npc := range gameMap.npc {
		gameover = gameover && npc.IsDead()
	}
	if gameMap.hero.__.(IObject).IsDead() {
		gameover = true
	}
	if gameover {
		s := "You win"
		if gameMap.hero.__.(IObject).IsDead() {
			s = "You lose"
		}
		log.Info(s)
	}
}

func gameInit() {
	conio.NewKeyboardAction("Status", "i", "", func(ev conio.TKeyboardEvent) bool {
		log.Info(gameMap.hero.GetInvetory())
		return true
	})
	conio.NewKeyboardAction("Up", "w", "", func(ev conio.TKeyboardEvent) bool {
		// LookAt(dirUp)
		gameMap.hero.SetDir(dirUp)
		gameMap.hero.Do()
		return true
	})
	conio.NewKeyboardAction("Down", "s", "", func(ev conio.TKeyboardEvent) bool {
		// LookAt(dirDown)
		gameMap.hero.SetDir(dirDown)
		gameMap.hero.Do()
		return true
	})
	conio.NewKeyboardAction("Left", "a", "", func(ev conio.TKeyboardEvent) bool {
		// LookAt(dirLeft)
		gameMap.hero.SetDir(dirLeft)
		gameMap.hero.Do()
		return true
	})
	conio.NewKeyboardAction("Right", "d", "", func(ev conio.TKeyboardEvent) bool {
		// LookAt(dirRight)
		gameMap.hero.SetDir(dirRight)
		gameMap.hero.Do()
		return true
	})
}

func main() {
	// log.Add(zlogger.Build().Format("(~m) ~l:~s~x~e\n").Done())
	log.Info("PING")
	err := conio.Init()
	utils.Assert(err == nil, "conio init failed")
	defer conio.Close()
	evs := conio.NewEventStream()
	utils.Assert(evs != nil, "eventStream init failed")
	defer evs.Close()
	scr := conio.NewScreen()
	utils.Assert(scr != nil, "screen init failed")
	defer scr.Close()
	nextTime := time.Now()
	// registerScreen(scr)
	initialize()
	newMap()

	for !canClose {
		draw()
		if evs.HasEvent() {
			ev := evs.ReadEvent()
			conio.HandleEvent(ev)
		}
		if time.Now().Before(nextTime) {
			continue
		}
		if !gameover {
			ai()
			nextTime = nextTime.Add(1000 * time.Millisecond)
		}
	}
}
