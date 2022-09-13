package models

type CreateWeaponDTO struct {
	Name string `json:"name"`
}

type WeaponDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
