package util

import "github.com/alexodle/wrappy/testdata/input/animals"

type Utils struct {
	Store *animals.AnimalStore
}

func NewUtils(st *animals.AnimalStore) *Utils {
	return &Utils{Store: st}
}

func (u *Utils) AddAllAnimals(animals []*animals.Animal) {
	for _, an := range animals {
		u.Store.AddAnimal(an)
	}
}
