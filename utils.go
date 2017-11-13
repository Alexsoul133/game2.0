package main

import (
	"math/rand"
	"time"

	"github.com/macroblock/garbage/conio"
)

type TDirection int //0Up 1Right 2Down 3Left

var dirArray4 = [...]struct{ x, y int }{{x: 0, y: -1}, {x: 1, y: 0}, {x: 0, y: 1}, {x: -1, y: 0}}
var dirArray8 = [...]struct{ x, y int }{{x: 0, y: -1}, {x: 1, y: -1}, {x: 1, y: 0}, {x: 1, y: 1}, {x: 0, y: 1}, {x: -1, y: 1}, {x: -1, y: 0}, {x: -1, y: -1}}

const (
	dirUp TDirection = iota
	dirRight
	dirDown
	dirLeft
)

func random(i int) bool {
	if rand.Intn(100) < i {
		return true
	}
	return false
}

func randomInt(min, max int) int {
	rand.Seed(int64(time.Now().Nanosecond() + len(gameMap.item) + time.Now().Nanosecond()))
	// rand.Seed(time.Now().UTC().UnixNano())
	res := min + rand.Intn(max-min)
	log.Info("RandomInt ", res)
	return res
}

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

func (o TDirection) GetOffset() (int, int) {
	o %= 4
	if o < 0 {
		o = 4 + o
	}
	delta := dirArray4[o]
	return delta.x, delta.y
}

func NewDir(dx, dy int) TDirection {
	if Abs(dx) > Abs(dy) {
		if dx > 0 {
			return dirRight
		}
		return dirLeft
	}
	if dy > 0 {
		return dirDown
	}
	return dirUp
}

func drawObject(x, y int, str string, color, bgcolor conio.TColor) {
	scr := conio.Screen()
	c := scr.FgColor()
	bg := scr.BgColor()
	scr.SetFgColor(color)
	scr.SetBgColor(bgcolor)
	scr.DrawAlignedString(mapX+x*2+1, mapY+y, mapW*2, str)
	scr.SetFgColor(c)
	scr.SetBgColor(bg)
}

func drawGround(x, y int, str string, color, bgcolor conio.TColor) {
	scr := conio.Screen()
	c := scr.FgColor()
	bg := scr.BgColor()
	scr.SetFgColor(color)
	scr.SetBgColor(bgcolor)
	conio.Screen().DrawAlignedString(mapX+x*2, mapY+y, mapW*2, str)
	scr.SetFgColor(c)
	scr.SetBgColor(bg)
}

func placeGround(o IGround) {
	x, y := o.GetPos()
	// log.Debug("+Ground ", x, y)
	cell := TCell{}
	cell.ground = o
	gameMap.data[x+y*mapW] = &cell
}

func addNpc(o IObject) {
	gameMap.npc = append(gameMap.npc, o)
}

func deleteNpc(o IObject) {
	x, y := o.GetPos()
	for i, npc := range gameMap.npc {
		if npcx, npcy := npc.GetPos(); npcx == x && npcy == y {
			gameMap.npc[i] = nil
		}
	}
	newNpc := []IObject{}
	for _, npc := range gameMap.npc {
		if npc != nil {
			newNpc = append(newNpc, npc)
		}
	}
	gameMap.npc = newNpc
}

func findNpc(x, y int) IObject {
	if herox, heroy := gameMap.hero.__.(IObject).GetPos(); herox == x && heroy == y && !gameMap.hero.__.(IObject).IsDead() {
		return gameMap.hero
	}
	for _, npc := range gameMap.npc {
		if npcx, npcy := npc.GetPos(); npcx == x && npcy == y && !npc.IsDead() {
			return npc
		}
	}
	return nil
}

func addItem(o IItems) {
	gameMap.item = append(gameMap.item, o)
}

func deleteitem(o IItems) {
	x, y := o.GetPos()
	for i, item := range gameMap.item {
		if itemx, itemy := item.GetPos(); itemx == x && itemy == y {
			gameMap.item[i] = nil
		}
	}
	newItem := []IItems{}
	for _, item := range gameMap.item {
		if item != nil {
			newItem = append(newItem, item)
		}
	}
	gameMap.item = newItem
}

func findItem(x, y int) IItems {
	for _, item := range gameMap.item {
		if itemx, itemy := item.GetPos(); itemx == x && itemy == y {
			return item
		}
	}
	// log.Debug("not found item")
	return nil
}

func fromInvetory(i int, object *TObject) {
	object.items[i] = nil
	newItem := []IItems{}
	for _, item := range object.items {
		if item != nil {
			newItem = append(newItem, item)
		}
	}
	object.items = newItem
}

func fromEquipment(i int, object *TObject) {
	object.equipment[i] = nil
	newItem := []IItems{}
	for _, item := range object.equipment {
		if item != nil {
			newItem = append(newItem, item)
		}
	}
	object.equipment = newItem
}

func equip(o *TObject, i int) bool {
	defer throw.Catch()
	throw.Panic(o == nil, "Object is nil")
	throw.Return(getCap(o.items)-1 < i)

	// for _, invetory := range o.items {
	// 	if invetory == o.items[i] {
	// 		log.Debug("Found item")
	// 		break
	// 	}
	// }
	throw.Return(!o.items[i].RespondToEquip())

	for i1, equipment := range o.__.(IPicker).GetEquipment() {
		if equipment.GetSlot() == o.items[i].GetSlot() {
			o.items = append(o.items, equipment)
			fromEquipment(i1, o)
			// slice := o.items[i+1:]
			// o.items = append(o.items[:i-1], o.items[i+1:])
		}

	}
	o.equipment = append(o.__.(IPicker).GetEquipment(), o.items[i])
	fromInvetory(i, o)
	log.Info(o.__.(IObject).GetType(), " equip ", o.items[i].GetType())
	// o.items[i] = nil
	return true
}

func getLen(o Invetory) int {
	for i, item := range o {
		if item == nil {
			return i + 1
		}
	}
	return len(o)
}

func getCap(o Invetory) int {
	return cap(o)
}

func move(o *TObject, d TDirection) bool {
	defer throw.Catch()
	throw.Panic(o == nil, "Object is nil")

	dst, npc := o.Look(d)
	throw.Return(dst == nil)
	throw.Return(npc != nil)

	throw.Return(o.__.(IObject).GetMovementType()&dst.ground.GetSurfaceType() == 0)
	o.x, o.y = dst.ground.GetPos()
	return true
}

func attack(o *TObject, d TDirection) bool {
	defer throw.Catch()
	throw.Panic(o == nil, "Object is nil")

	dst, npc := o.Look(d)
	throw.Return(dst == nil)
	throw.Return(npc == nil)

	dmg := o.__.(IObject).GetDmg()
	damagable := IDamagable(nil)
	damagable, ok := npc.(IDamagable)
	throw.Return(!ok)

	damagable.RecieveDmg(dmg)
	if npc.GetHp() <= 0 {
		log.Info(o.__.(IObject).GetType(), " kill ", npc.GetType(), ". Inflicted ", o.__.(IObject).GetDmg(), " dmg")
		o.__.(IObject).IncreaseExp(3)
		return true
	}
	log.Info(o.__.(IObject).GetType(), " attack ", npc.GetType(), ". Inflicted ", o.__.(IObject).GetDmg(), " dmg")
	return true
}

func pickup(o *TObject, d TDirection) bool {
	defer throw.Catch()
	throw.Panic(o == nil, "Object is nil")

	dst, item := o.LookItem(d)
	throw.Return(dst == nil)
	throw.Return(item == nil)
	pickable := IPickable(nil)
	pickable, ok := item.(IPickable)
	throw.Return(!ok)
	throw.Return(!pickable.RespondToPick())
	o.items = append(o.items, item)
	log.Info(o.__.(IObject).GetType(), " pick up ", item.GetType())
	// log.Info(fmt.Sprintf("%v", item))
	return true
}
