package main

import (
	"fmt"

	"github.com/macroblock/garbage/conio"
)

const (
	surfaceNil int = (1 << iota) >> 1
	surfaceGround
	surfaceWall
	surfaceGrass
	surfaceWater
	surfaceLast
	surfaceAll = surfaceLast - 1
)

///////////////////////////////////////////////////////
type TBasic struct {
	x, y      int
	lvl, exp  int
	gainedDmg int
	__        interface{}
}

type TGround struct {
	TBasic
}

type IBasic interface {
	GetType() string
	Status() string
	Draw()
	GetPos() (int, int)
	Do()
	GetLvl() int
	GetMaxHp() int
	GetHp() int
	GetDmg() int
	GetStamina() int
	GetStr() int
	// RecieveDmg(i int)
	GetExp() int
	GetExpLvl() int
	IncreaseExp(lvl int)

	// placeGround(o *TGround)
}

type IDamagable interface {
	RecieveDmg(i int)
}

type IGround interface {
	IBasic
	GetSurfaceType() int
}

type TWall struct {
	TGround
}

type TGrass struct {
	TGround
}

type TWater struct {
	TGround
}

func (o *TBasic) GetType() string {
	return "Ground"
}

func (o *TBasic) Status() string {
	// return "status"
	return fmt.Sprintf("%v", o.__.(IBasic).GetType())
}

func (o *TBasic) Draw() {
	drawGround(o.x, o.y, ".\"", conio.ColorDarkGray, conio.ColorBlack)
}

func (o *TBasic) GetPos() (int, int) {
	return o.x, o.y
}

func (o *TBasic) Do() {

}

/////////////////////////////////////////////////////
func newGround(x, y int) *TGround {
	o := &TGround{}
	o.x = x
	o.y = y
	o.__ = o
	return o
}

/////////////////////////////////////////////////////

func (o *TGround) GetSurfaceType() int {
	return surfaceGround
}

/////////////////////////////////////////////////////
func newWater(x, y int) *TWater {
	o := &TWater{}
	o.x = x
	o.y = y
	o.__ = o
	return o
}

func (o *TWater) GetType() string {
	return "Water"
}

func (o *TWater) Draw() {
	drawGround(o.x, o.y, "~~", conio.ColorCyan, conio.ColorBlue)
}
func (o *TWater) GetSurfaceType() int {
	return surfaceWater
}

/////////////////////////////////////////////////////

func newWall(x, y int) *TWall {
	o := &TWall{}
	o.x = x
	o.y = y
	o.__ = o
	return o
}

func (o *TWall) GetType() string {
	return "Wall"
}

func (o *TWall) Draw() {
	drawGround(o.x, o.y, "##", conio.ColorBlack, conio.ColorRed)
}
func (o *TWall) GetSurfaceType() int {
	return surfaceWall
}

/////////////////////////////////////////////////////
func newGrass(x, y int) *TGrass {
	o := &TGrass{TGround: *newGround(x, y)}
	o.__ = o
	return o
}

func (o *TGrass) GetType() string {
	return "Grass"
}

func (o *TGrass) Draw() {
	drawGround(o.x, o.y, ",,", conio.ColorGreen, conio.ColorBlack)
}

func (o *TGrass) GetSurfaceType() int {
	return surfaceGrass
}
