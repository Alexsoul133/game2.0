package main

import (
	"fmt"
	"math/rand"

	"github.com/macroblock/garbage/conio"
)

type TObject struct {
	TBasic
	items ItemList
	fov   int
}

type THero struct {
	THuman
	moved bool
}

var _ IObject = (*TObject)(nil)

type IObject interface {
	IBasic
	Look(d TDirection) (*TCell, IObject)
	LookItem(d TDirection) (*TCell, IItems)
	GetMovementType() int
	// PickUp(d TDirection) bool
	IsDead() bool
	// LookAround() []IObject
	FindTarget() IObject
}

type IFighter interface {
	Attack() bool
	SpinAttack() bool
}

type IMove interface {
	Move() bool
}

var _ IPicker = (*THero)(nil)

type IPicker interface {
	PickUp() bool
}

var _ IPickable = (*TObject)(nil)

type IPickable interface {
	RespondToPick() bool
}

var _ IFighter = (*THero)(nil)
var _ IObject = (*THero)(nil)
var _ IBasic = (*TBasic)(nil)

type TCat struct {
	TObject
	dir TDirection
}

type TDog struct {
	TCat
}

type THuman struct {
	TCat
}

type TSpinner struct {
	TObject
	phase int
}

type TPlant struct {
	TSpinner
	dir TDirection
}

type TPiranha struct {
	TCat
}

func random(i int) bool {
	if rand.Intn(100) < i {
		return true
	}
	return false
}

///////////////////////////////////////////

func newObject(x, y int) *TObject {
	o := &TObject{}
	o.x = x
	o.y = y
	o.items = nil
	o.__ = o
	// placeObject(o)
	return o
}

func (o *TObject) Draw() {
	drawObject(o.x, o.y, "", conio.ColorWhite, conio.ColorBlack)
}

func (o *TObject) GetType() string {
	return "Nothing"
}

func (o *TObject) Look(d TDirection) (*TCell, IObject) {
	dx, dy := d.GetOffset()
	offset := (o.x + dx) + (o.y+dy)*mapW
	if (o.x+dx) < 0 || (o.x+dx) > mapW-1 || (o.y+dy) < 0 || (o.y+dy) > mapH-1 {
		return nil, nil
	}
	// log.Debug(o.__.(IObject).GetType(), " ", dx, " ", dy, ":", o.x, " ", o.y)
	object := findNpc(o.x+dx, o.y+dy)
	cell := gameMap.data[offset]
	return cell, object
}

func (o *TObject) LookItem(d TDirection) (*TCell, IItems) {
	dx, dy := d.GetOffset()
	offset := (o.x + dx) + (o.y+dy)*mapW
	if (o.x+dx) < 0 || (o.x+dx) > mapW-1 || (o.y+dy) < 0 || (o.y+dy) > mapH-1 {
		return nil, nil
	}
	object := findItem(o.x+dx, o.y+dy)
	cell := gameMap.data[offset]
	return cell, object
}

// func (o *TObject) LookAround() IObject {
// 	// target := []IObject{}
// 	for i := 0; i <= o.fov; i++ {
// 		o.Look()
// 	}
// }

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Distance(o1, o2 IObject) int {
	x1, y1 := o1.GetPos()
	x2, y2 := o2.GetPos()
	return Abs(x2 - x1 + y2 - y1)
}

func (o *TObject) FindTarget() IObject {
	minTarget := IObject(gameMap.hero)
	minDistance := Distance(o, minTarget)
	for _, npc := range gameMap.npc {
		d := Distance(o, npc)
		if d == 0 || npc.IsDead() {
			continue
		}
		if d < minDistance {
			minTarget = npc
			minDistance = d
		}
	}
	return minTarget
}

func (o *TObject) Do() {
	if o.__.(IObject).IsDead() {
		return
	}

	fighter, ok := o.__.(IFighter)
	if ok && fighter.Attack() {
		return
	}
	pick, ok := o.__.(IPicker)
	if ok && pick.PickUp() {
		return
	}
	move, ok := o.__.(IMove)
	if ok && move.Move() {
		return
	}
}

func (o *TObject) GetMovementType() int {
	return surfaceAll &^ surfaceWater &^ surfaceWall
}

func (o *TObject) IsDead() bool {
	return o.__.(IObject).GetHp() <= 0
}

func (o *TObject) RespondToPick() bool {
	return false
}

/////////////////////////////////////////////////////

func newHero(x, y int) *THero {
	o := &THero{THuman: *newHuman(x, y)}
	o.items = nil
	o.__ = o

	return o
}

func (o *THero) GetDmg() int {
	return o.__.(IObject).GetStr()*3/5 + o.items.TotalDmg()
}

func (o *THero) Move() bool {
	if !move(&o.TObject, o.dir) {
		log.Debug(o.__.(IObject).GetType(), " cant move")
		return false
	}
	return true
}

func (o *THero) Do() {
	if o.__.(IObject).IsDead() {
		return
	}
	if o.moved {
		return
	}
	o.moved = true
	o.TObject.Do()
}

func (o *THero) SetDir(d TDirection) {
	o.moved = false
	o.dir = d
}

func (o *THero) Draw() {
	sprite := "H"
	fg, bg := conio.Screen().CellColor((o.x*2)+mapX, o.y+mapY)
	fg = conio.ColorRed
	// bg := conio.ColorDefault
	if o.__.(IObject).IsDead() {
		fg = conio.ColorDarkGray
	}
	drawObject(o.x, o.y, sprite, fg, bg)
}

func (o *THero) GetType() string {
	return "Hero"
}

func (o *THero) GetInvetory() []string {
	var items []string
	for i := range o.items {
		items = append(items, o.items[i].GetType())
	}
	return items
}

/////////////////////////////////////////////////////
func newCat(x, y int) *TCat {
	o := &TCat{TObject: *newObject(x, y)}
	o.lvl = 1
	o.exp = 1
	o.fov = 4
	o.__ = o
	log.Debug(fmt.Sprintf("%+v", o.__.(IObject)), " created ", x, " ", y)
	return o
}

func (o *TCat) GetType() string {
	return "Cat"
}

func (o *TCat) Draw() {
	sprite := "C"
	fg := conio.ColorWhite
	bg := conio.ColorBlack
	if o.__.(IObject).IsDead() {
		fg = conio.ColorDarkGray
	}
	drawObject(o.x, o.y, sprite, fg, bg)
}
func (o *TCat) GetStamina() int {
	return 3 + o.__.(IObject).GetLvl()
}

func (o *TCat) GetStr() int {
	return 2 + o.__.(IObject).GetLvl()
}

func (o *TCat) RecieveDmg(i int) {
	o.gainedDmg += i
}

func (o *TCat) Move() bool {
	for i := 0; i < 4; i++ {
		if move(&o.TObject, o.dir) {
			return true
		}
		o.dir++
	}
	log.Debug(o.__.(IObject).GetType(), " cant move")
	return false
}

func (o *TCat) Attack() bool {
	if !attack(&o.TObject, o.dir) {
		// log.Debug(o.__.(IObject).GetType(), " cant attack")
		return false
	}
	return true
}

func (o *TCat) SpinAttack() bool {
	return false
}
func (o *TCat) PickUp() bool {
	if !pickup(&o.TObject, o.dir) {
		return false
	}
	return true
}

func (o *TCat) Do() {
	// o.dir = TDirection(rand.Intn(4))
	x, y := o.FindTarget().GetPos()
	dx := x - o.x
	dy := y - o.y
	o.dir = NewDir(dx, dy)
	o.TObject.Do()
}

/////////////////////////////////////////////////////
func newDog(x, y int) *TDog {
	o := &TDog{TCat: *newCat(x, y)}
	o.__ = o
	log.Debug(o.__.(IObject).GetType(), " created ", x, " ", y)
	return o
}

/////////////////////////////////////////////////////

func newHuman(x, y int) *THuman {
	o := &THuman{TCat: *newCat(x, y)}
	o.__ = o
	log.Debug(o.__.(IObject).GetType(), " created ", x, " ", y)
	return o
}

func (o *THero) GetMovementType() int {
	return surfaceAll&^surfaceWater&^surfaceWall | o.items.TotalMovementType()
}

func (o *THuman) GetType() string {
	return "Human"
}

func (o *THuman) Draw() {
	sprite := "☺"
	fg := conio.ColorWhite
	bg := conio.ColorBlack
	if o.__.(IObject).IsDead() {
		fg = conio.ColorDarkGray
	}
	drawObject(o.x, o.y, sprite, fg, bg)
}

/////////////////////////////////////////////////////
func newSpinner(x, y int) *TSpinner {
	o := &TSpinner{TObject: *newObject(x, y)}
	o.lvl = 1
	o.__ = o
	log.Debug(o.GetType(), " created ", x, " ", y)
	return o
}

func (o *TSpinner) GetType() string {
	return "Spinner"
}

const spinnerPhase = "-\\|/"

func (o *TSpinner) Draw() {
	sprite := string(spinnerPhase[o.phase])
	fg := conio.ColorWhite
	bg := conio.ColorBlack
	if o.__.(IObject).IsDead() {
		fg = conio.ColorDarkGray
	}
	drawObject(o.x, o.y, sprite, fg, bg)
}
func (o *TSpinner) GetStamina() int {
	return 3 + o.__.(IObject).GetLvl()
}

func (o *TSpinner) GetStr() int {
	return 2 + o.__.(IObject).GetLvl()
}
func (o *TSpinner) SpinAttack() bool {
	var d TDirection
	for i := 0; i < 4; i++ {
		attack(&o.TObject, d)
		d++
	}
	return true
}

func (o *TSpinner) RecieveDmg(i int) {
	o.gainedDmg += i
}
func (o *TSpinner) Do() {
	if o.__.(IObject).IsDead() {
		return
	}
	o.phase++
	o.phase %= len(spinnerPhase)
	o.SpinAttack()
}

/////////////////////////////////////////////////////
func newPlant(x, y int) *TPlant {
	o := &TPlant{TSpinner: *newSpinner(x, y)}
	o.__ = o
	log.Debug(o.GetType(), " created ", x, " ", y)
	return o
}

func (o *TPlant) GetType() string {
	return "plant"
}
func (o *TPlant) IsMaterial() bool {
	return true
}

const plantPhase = "+xO"

func (o *TPlant) Draw() {
	sprite := string(plantPhase[o.phase])
	fg := conio.ColorGreen
	bg := conio.ColorBlack
	if o.__.(IObject).IsDead() {
		fg = conio.ColorDarkGray
	}
	drawObject(o.x, o.y, sprite, fg, bg)
}

func (o *TPlant) Mul() {
	for i := 0; i < 4; i++ {
		dst, npc := o.Look(o.dir)
		dx, dy := o.dir.GetOffset()
		if dst == nil || npc != nil || dst.ground.GetSurfaceType() != (surfaceAll&^surfaceWall&^surfaceWater) {
			o.dir++
			continue
		}
		addNpc(newPlant(o.x+dx, o.y+dy))

		o.dir++
	}
}

func (o *TPlant) Do() {
	if o.__.(IObject).IsDead() {
		return
	}
	if random(50) {
		o.phase++
	}
	if o.phase == 3 {
		o.Mul()
		addNpc(newObject(o.x, o.y))
	}
	o.phase %= len(plantPhase)

}

/////////////////////////////////////////////////////
func newPiranha(x, y int) *TPiranha {
	o := &TPiranha{TCat: *newCat(x, y)}
	o.x = x
	o.y = y
	o.__ = o
	log.Debug(o.GetType(), " created ", x, " ", y)
	return o
}

func (o *TPiranha) GetMovementType() int {
	return surfaceWater
}
func (o *TPiranha) Draw() {
	sprite := "F"
	fg := conio.ColorBlack
	bg := conio.ColorBlue
	if o.__.(IObject).IsDead() {
		fg = conio.ColorDarkGray
	}
	drawObject(o.x, o.y, sprite, fg, bg)
}

func (o *TPiranha) GetType() string {
	return "Piranha"
}
