package main

import "math"

func (o *TBasic) GetLvl() int {
	return o.lvl
}
func (o *TBasic) GetExp() int {
	return o.exp
}

func (o *TBasic) GetExpLvl() int {
	i := o.__.(IObject).GetLvl() * 10 * int(math.Pow(1.1, 2))
	return i
}
func (o *TBasic) IncreaseExp(i int) {
	o.exp += i
	if o.__.(IObject).GetExp() >= o.__.(IObject).GetExpLvl() {
		o.lvl++
		log.Info(o.__.(IObject).GetType(), " reach level up")
	}
}
func (o *TBasic) GetMaxHp() int {
	return o.__.(IObject).GetLvl() * o.__.(IObject).GetStamina()
}

func (o *TBasic) GetHp() int {
	return o.__.(IObject).GetMaxHp() - o.gainedDmg
}

func (o *TBasic) GetDmg() int {
	return o.__.(IObject).GetStr() * 3 / 5
}

func (o *TBasic) GetStamina() int {
	return o.__.(IObject).GetLvl() * int(math.Pow(1.1, 2))
}

func (o *TBasic) GetStr() int {
	return o.__.(IObject).GetLvl() * int(math.Pow(1.1, 2))
}
