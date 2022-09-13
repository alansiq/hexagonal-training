package models

type HeroDto struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	Lastname string     `json:"last_name"`
	Age      int        `json:"age"`
	Level    int        `json:"level"`
	Type     string     `json:"type"`
	WeaponID int        `json:"weapon_id"`
	Weapon   *WeaponDTO `json:"weapon,omitempty"`
}

type CreateHeroDto struct {
	Name     string `json:"name"`
	Lastname string `json:"last_name"`
	Age      int    `json:"age"`
	Level    int    `json:"level"`
	Type     string `json:"type"`
	WeaponID int    `json:"weapon_id"`
}
