package animals

type Animal struct {
	AnimalDescription
}

type AnimalDescription struct {
	Name string
	Age int
	Type string
	Weight int
}

type AnimalStore struct {
	animals map[string]*Animal
}
func NewAnimalStore() *AnimalStore {
	return &AnimalStore{animals: map[string]*Animal{}}
}

func (a *AnimalStore) GetAnimalByName(name string) (*Animal, bool) {
	animal, ok := a.animals[name]
	return animal, ok
}

func (a *AnimalStore) GetAllAnimals() []*Animal {
	var ans []*Animal
	for _, an := range a.animals {
		ans = append(ans, an)
	}
	return ans
}

func (a *AnimalStore) AddAnimal(animal *Animal) {
	a.animals[animal.Name] = animal
}
