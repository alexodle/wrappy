package animals

import orig_animals "github.com/alexodle/wrappy/testdata/input/animals"

type Animal interface {
	GetImpl() *orig_animals.Animal
	GetAnimalDescription() AnimalDescription
	SetAnimalDescription(v AnimalDescription)
}

func NewAnimal(impl *orig_animals.Animal) Animal {
	return &animalWrapper{impl: impl}
}

type animalWrapper struct {
	impl *orig_animals.Animal
}

func (o *animalWrapper) GetImpl() *orig_animals.Animal {
	return o.impl
}

func (wrapperRcvr *animalWrapper) GetAnimalDescription() AnimalDescription {
	retval := &wrapperRcvr.impl.AnimalDescription
	retval_1 := NewAnimalDescription(retval)
	return retval_1
}

func (wrapperRcvr *animalWrapper) SetAnimalDescription(v AnimalDescription) {
	v_1 := *v.GetImpl()
	wrapperRcvr.impl.AnimalDescription = v_1
}

type AnimalDescription interface {
	GetImpl() *orig_animals.AnimalDescription
	GetAge() int
	GetName() string
	GetType() string
	GetWeight() int
	SetAge(v int)
	SetName(v string)
	SetType(v string)
	SetWeight(v int)
}

func NewAnimalDescription(impl *orig_animals.AnimalDescription) AnimalDescription {
	return &animalDescriptionWrapper{impl: impl}
}

type animalDescriptionWrapper struct {
	impl *orig_animals.AnimalDescription
}

func (o *animalDescriptionWrapper) GetImpl() *orig_animals.AnimalDescription {
	return o.impl
}

func (wrapperRcvr *animalDescriptionWrapper) GetName() string {
	retval := wrapperRcvr.impl.Name
	return retval
}

func (wrapperRcvr *animalDescriptionWrapper) SetName(v string) {
	wrapperRcvr.impl.Name = v
}

func (wrapperRcvr *animalDescriptionWrapper) GetAge() int {
	retval := wrapperRcvr.impl.Age
	return retval
}

func (wrapperRcvr *animalDescriptionWrapper) SetAge(v int) {
	wrapperRcvr.impl.Age = v
}

func (wrapperRcvr *animalDescriptionWrapper) GetType() string {
	retval := wrapperRcvr.impl.Type
	return retval
}

func (wrapperRcvr *animalDescriptionWrapper) SetType(v string) {
	wrapperRcvr.impl.Type = v
}

func (wrapperRcvr *animalDescriptionWrapper) GetWeight() int {
	retval := wrapperRcvr.impl.Weight
	return retval
}

func (wrapperRcvr *animalDescriptionWrapper) SetWeight(v int) {
	wrapperRcvr.impl.Weight = v
}

type AnimalStore interface {
	GetImpl() *orig_animals.AnimalStore
	AddAnimal(animal Animal)
	GetAllAnimals() []Animal
	GetAnimalByName(name string) (Animal, bool)
}

func NewAnimalStore(impl *orig_animals.AnimalStore) AnimalStore {
	return &animalStoreWrapper{impl: impl}
}

type animalStoreWrapper struct {
	impl *orig_animals.AnimalStore
}

func (o *animalStoreWrapper) GetImpl() *orig_animals.AnimalStore {
	return o.impl
}

func (wrapperRcvr *animalStoreWrapper) GetAnimalByName(name string) (Animal, bool) {
	retval0, retval1 := wrapperRcvr.impl.GetAnimalByName(name)
	retval0_1 := NewAnimal(retval0)
	return retval0_1, retval1
}

func (wrapperRcvr *animalStoreWrapper) GetAllAnimals() []Animal {
	retval0 := wrapperRcvr.impl.GetAllAnimals()
	var retval0_1 []Animal
	for _, it := range retval0 {
		it_1 := NewAnimal(it)
		retval0_1 = append(retval0_1, it_1)
	}
	return retval0_1
}

func (wrapperRcvr *animalStoreWrapper) AddAnimal(animal Animal) {
	animal_1 := animal.GetImpl()
	wrapperRcvr.impl.AddAnimal(animal_1)
}
