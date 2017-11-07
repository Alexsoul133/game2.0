package main

import (
	"fmt"

	"github.com/macroblock/garbage/conio"
	"github.com/macroblock/garbage/utils"
	"github.com/macroblock/zl/core/zlogger"
)

var (
	key           int
	ch            rune
	width, height int
	canClose      bool
	msgLog        = make([]string, 0)
)

const (
	mapX     = 2
	mapY     = 2
	mapW     = 15
	mapH     = 15
	logWidth = 100
)

var CustomWriter TCustomWriter

type TCustomWriter struct {
}

func (o *TCustomWriter) Write(p []byte) (n int, err error) {
	msgLog = append(msgLog, string(p[:]))
	return len(p), nil
}
func draw() {
	scr := conio.Screen()
	winFg := conio.ColorWhite
	winBg := conio.ColorBlack
	logFg := conio.ColorWhite
	logBg := conio.ColorDarkGray
	scr.Clear('░', conio.ColorWhite, conio.ColorBlack)
	scr.SelectBorder("Single")
	scr.SetColor(logFg, logBg)
	drawWindow(width-logWidth-1, 0, logWidth, height, "[ Log ]", func(x, y, w, h int) {
		for i := 0; i < utils.Min(len(msgLog), h); i++ {
			scr.DrawAlignedString(x, y+h-1-i, w, msgLog[len(msgLog)-1-i])
		}
	})

	scr.SelectBorder("Double")
	scr.SetColor(winFg, winBg)
	drawWindow(mapX-1, mapY-1, mapW*2+2, mapH+2, "[ Map ]", func(x, y, w, h int) {
		for i := 0; i < mapW*mapH; i++ {
			gameMap.data[i%mapH+i/mapH*mapW].ground.Draw()
			// drawCell(scr, i%mapH, i/mapH)
			//scr.DrawAlignedString(x+offsX, y+i, w-offsX, name)
		}
		for _, npc := range gameMap.npc {
			if npc.IsDead() {
				npc.Draw()
			}
		}
		for _, npc := range gameMap.npc {
			if !npc.IsDead() {
				npc.Draw()
			}
		}
		for _, items := range gameMap.item {
			items.Draw()
		}
		gameMap.hero.__.(IObject).Draw()
	})
	drawWindow(mapX-1, mapY-1+mapH+2+1, width-logWidth-1-1, len(gameMap.npc)*2+2+2, "[ NPC ]", func(x, y, w, h int) {
		drawStatus(x, y, w, gameMap.hero.__.(IObject))
		for i, npc := range gameMap.npc {
			drawStatus(x, y+(i+1)*2, w, npc)
		}
	})
	scr.Flush()
}

func drawStatus(x, y, w int, o IObject) {
	dead := "alive"
	if o.IsDead() {
		dead = "dead"
	}
	status := fmt.Sprintf(" %v LVL:%v (%v/%v) HP:%v/%v %v", o.GetType(), o.GetLvl(), o.GetExp(), o.GetExpLvl(), o.GetHp(), o.GetMaxHp(), dead)
	conio.Screen().DrawAlignedString(x, y, w, status)
	status = fmt.Sprintf("   Target - %v", o.FindTarget().GetType())
	conio.Screen().DrawAlignedString(x, y+1, w, status)
}

func drawWindow(x, y, w, h int, title string, draw func(x, y, w, h int)) {
	scr := conio.Screen()
	scr.DrawBorder(x, y, w, h)
	scr.FillRect(x+1, y+1, w-2, h-2, ' ')
	scr.DrawAlignedString(x+1, y, w-2, title)
	draw(x+1, y+1, w-2, h-2)
}

func initialize() {
	log.Add(zlogger.Build().Format("(~m) ~l~s~x~e").Writer(&CustomWriter).Done())

	conio.Screen().Flush()
	width = conio.Screen().Width()
	height = conio.Screen().Height() - 26

	gameInit()
	conio.NewKeyboardAction("Exit", "`", "", func(ev conio.TKeyboardEvent) bool {
		canClose = true
		return true
	})
	conio.NewKeyboardAction("Exitё", "ё", "", func(ev conio.TKeyboardEvent) bool {
		canClose = true
		return true
	})

	conio.ActionMap.Apply()
}
