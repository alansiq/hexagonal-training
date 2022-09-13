package models

type CreateArmDTO struct {
	Name string `json:"name"`
}

type ArmDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
