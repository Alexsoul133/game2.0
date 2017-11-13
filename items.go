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

const (
	slotUnequip TSlot = (1 << iota) >> 1
	slotRightHand
	slotLeftHand
	slotHelmet
	slotBody
	slotLegs
	slotFeets
	slotGloves
	slotNecklace
	slotBack
	slotAll = slotBack - 1
)

type TItems struct {
	TObject
	quality               int
	dmg, stamina, str, hp int
	slot                  TSlot //2hand? set?
}

type IItems interface {
	IObject
	// Dmg() int
	RespondToPick() bool
	RespondToEquip() bool
	GetQuality() string
	GetSlot() TSlot
}

type Invetory []IItems

type TSlot int //1 RightHand 2 lefthand 3 helmet 4 body 5 legs 6 feet 7 hands 8 necklace 9 back

var _ IObject = (*TItems)(nil)
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

func (o *TItems) GetSlot() TSlot {
	return slotUnequip
}

func (o *TItems) RespondToPick() bool {
	deleteitem(o)
	return true
}

func (o *TItems) RespondToEquip() bool {
	if o.__.(IItems).GetSlot() != slotUnequip {
		return true
	}
	log.Info(o.GetType(), " is unquip")
	return false
}

func newItem(x, y int) *TItems {
	o := &TItems{}
	// o.quality = randomInt(-1, 3)
	o.x = x
	o.y = y
	o.__ = o
	return o
}

func (o *Invetory) TotalDmg() int {
	dmg := 0
	for _, item := range *o {
		if item != nil {
			dmg += item.GetDmg()
		}
	}
	return dmg
}

func (o *Invetory) TotalMovementType() int {
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

func (o *TSword) GetSlot() TSlot {
	return slotRightHand | slotLeftHand
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

func (o *TBoots) GetSlot() TSlot {
	return slotFeets
}
