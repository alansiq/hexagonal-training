package domain

type Hero struct {
	Name     string
	Lastname string
	Age      int
	Level    int
	Type     string
	Weapon   *Weapon
	WeaponId int
}
