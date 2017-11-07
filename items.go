package main

import "github.com/macroblock/garbage/conio"

type TItems struct {
	TObject
}

type IItems interface {
	IObject
	Dmg() int
	RespondToPick() bool
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
	return 0
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

func (o *TItems) GetHp() int {
	return o.__.(IObject).GetMaxHp() - o.gainedDmg
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

func (o *TObject) Dmg() int {
	return 0
}
func (o *TItems) Dmg() int {
	dmg := 0
	for _, item := range o.items {
		dmg += item.GetDmg()
	}
	return dmg
}

func (o *TItems) RespondToPick() bool {
	deleteitem(o)
	return true
}

func newItem(x, y int) *TItems {
	o := &TItems{}
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
	o.__ = o
	return o
}
func (o *TSword) GetDmg() int {
	return 4
}
func (o *TSword) GetType() string {
	return "Sword"
}

/////////////////////////////////////////////////////////
func newBoots(x, y int) *TBoots {
	o := &TBoots{TItems: *newItem(x, y)}
	o.__ = o
	return o
}

func (o *TBoots) GetMovementType() int {
	return surfaceWater | ^surfaceGrass
}

func (o *TBoots) GetType() string {
	return "Boots"
}
