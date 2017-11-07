package main

import (
	"math"

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

func (o TDirection) GetOffset() (int, int) {
	o %= 4
	if o < 0 {
		o = 4 + o
	}
	delta := dirArray4[o]
	return delta.x, delta.y
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
	return true
}

func pathfinder(o *TObject, x, y int) TDirection {
	defer throw.Catch()
	throw.Panic(o == nil, "Object is nil")
	dx := o.x - x
	dy := o.y - y
	throw.Return(dx == 0 && dy == 0)

	if math.Abs(float64(dx)) <= math.Abs(float64(dy)) { //move horizontal
		if dx > 0 {
			return dirRight
		}
		return dirLeft
	}
	if dy > 0 { //move vertical
		return dirDown
	}
	return dirUp
}
