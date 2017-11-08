package main

import (
	"math/rand"

	"github.com/macroblock/garbage/conio"
)

const (
	qPoor        = -1
	qRegular int = iota
	qUnusual
	qRare
	qLegendary
)

type TItems struct {
	TObject
	quality               int
	dmg, stamina, str, hp int
}

type IItems interface {
	IObject
	// Dmg() int
	RespondToPick() bool
	GetQuality() string
}
type ItemList []IItems

var _ IItems = (*TItems)(nil)

type TSword struct {
	TItems
}

type TBoots struct {
	TItems
}

func (o *TItems) GetMovementType() int {
	return 0
}

func (o *TItems) GetDmg() int {
	return o.dmg
}

func (o *TItems) GetStamina() int {
	return 0
}

func (o *TItems) GetStr() int {
	return 0
}
func (o *TItems) GetMaxHp() int {
	return 5
}

func (o *TItems) RecieveDmg(i int) {
}

func (o *TItems) Draw() {
	sprite := "ᴥ"
	fg := conio.ColorWhite
	bg := conio.ColorBlack
	drawObject(o.x, o.y, sprite, fg, bg)
}

func (o *TItems) GetType() string {
	return "Garbage"
}

func (o *TItems) GetQuality() string {
	switch o.quality {
	case qPoor:
		return "poor"
	case qRegular:
		return "regular"
	case qUnusual:
		return "unusual"
	case qRare:
		return "rare"
	case qLegendary:
		return "legendary"
	}
	return ""
}

func (o *TItems) RespondToPick() bool {
	deleteitem(o)
	return true
}

func newItem(x, y int) *TItems {
	o := &TItems{}
	// o.quality = randomInt(-1, 3)
	o.x = x
	o.y = y
	o.__ = o
	return o
}

func (o *ItemList) TotalDmg() int {
	dmg := 0
	for _, item := range *o {
		dmg += item.GetDmg()
	}
	return dmg
}

func (o *ItemList) TotalMovementType() int {
	i := 0
	for _, item := range *o {
		i += item.GetMovementType()
	}
	return i
}

/////////////////////////////////////////////////////////
func newSword(x, y int) *TSword {
	o := &TSword{TItems: *newItem(x, y)}
	o.quality = randomInt(-1, 3)
	switch o.quality {
	case qPoor:
		o.dmg = 1 + rand.Intn(1)
	case qRegular:
		o.dmg = 2 + rand.Intn(2)
	case qUnusual:
		o.dmg = 4 + rand.Intn(3)
	case qRare:
		o.dmg = 7 + rand.Intn(4)
	case qLegendary:
		o.dmg = 12 + rand.Intn(5)
	}
	o.__ = o
	return o
}
func (o *TSword) GetType() string {
	return "sword"
}

/////////////////////////////////////////////////////////
func newBoots(x, y int) *TBoots {
	o := &TBoots{TItems: *newItem(x, y)}
	o.__ = o
	return o
}

func (o *TBoots) GetMovementType() int {
	return surfaceWater
}

func (o *TBoots) GetType() string {
	return "boots"
}
