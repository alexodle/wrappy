package animals

import dog "github.com/alexodle/wrappy/testdata/generated/animals/dog"
import orig_animals "github.com/alexodle/wrappy/testdata/input/animals"
import orig_context "context"
import orig_dog "github.com/alexodle/wrappy/testdata/input/animals/dog"

type AnimalDescription interface {
	GetImpl() *orig_animals.AnimalDescription
	GetBreed() string
	GetName() string
	GetWeight() int
	SetBreed(v string)
	SetName(v string)
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

func (wrapperRcvr *animalDescriptionWrapper) GetBreed() string {
	retval := wrapperRcvr.impl.Breed
	return retval
}

func (wrapperRcvr *animalDescriptionWrapper) SetBreed(v string) {
	wrapperRcvr.impl.Breed = v
}

func (wrapperRcvr *animalDescriptionWrapper) GetName() string {
	retval := wrapperRcvr.impl.Name
	return retval
}

func (wrapperRcvr *animalDescriptionWrapper) SetName(v string) {
	wrapperRcvr.impl.Name = v
}

func (wrapperRcvr *animalDescriptionWrapper) GetWeight() int {
	retval := wrapperRcvr.impl.Weight
	return retval
}

func (wrapperRcvr *animalDescriptionWrapper) SetWeight(v int) {
	wrapperRcvr.impl.Weight = v
}

type Animals interface {
	GetImpl() *orig_animals.Animals
	AddAnimals(animals *[]interface{}) error
	AddDogs(dogs []dog.Dog) map[string]Animals
	GetAllDogs(ctx orig_context.Context) []dog.Dog
	GetAnimalDescription() AnimalDescription
	GetDogByName(name string) dog.Dog
	GetDogs() *[]dog.Dog
	GetDogsByNameField() map[string]dog.Dog
	GetDogsByNames(names []string) map[string]dog.Dog
	GetLocation() orig_animals.Location
	GetLocations() []orig_animals.Location
	SetAnimalDescription(v AnimalDescription)
	SetDogs(v *[]dog.Dog)
	SetDogsByNameField(v map[string]dog.Dog)
	SetLocation(v orig_animals.Location)
	SetLocations(v []orig_animals.Location)
}

func NewAnimals(impl *orig_animals.Animals) Animals {
	return &animalsWrapper{impl: impl}
}

type animalsWrapper struct {
	impl *orig_animals.Animals
}

func (o *animalsWrapper) GetImpl() *orig_animals.Animals {
	return o.impl
}

func (wrapperRcvr *animalsWrapper) GetLocations() []orig_animals.Location {
	retval := wrapperRcvr.impl.Locations
	return retval
}

func (wrapperRcvr *animalsWrapper) SetLocations(v []orig_animals.Location) {
	wrapperRcvr.impl.Locations = v
}

func (wrapperRcvr *animalsWrapper) GetLocation() orig_animals.Location {
	retval := wrapperRcvr.impl.Location
	return retval
}

func (wrapperRcvr *animalsWrapper) SetLocation(v orig_animals.Location) {
	wrapperRcvr.impl.Location = v
}

func (wrapperRcvr *animalsWrapper) GetDogs() *[]dog.Dog {
	retval := wrapperRcvr.impl.Dogs
	var retval_1 *[]dog.Dog
	if retval != nil {
		for _, it := range *retval {
			it_1 := dog.NewDog(it)
			*retval_1 = append(*retval_1, it_1)
		}
	}
	return retval_1
}

func (wrapperRcvr *animalsWrapper) SetDogs(v *[]dog.Dog) {
	var v_1 *[]*orig_dog.Dog
	if v != nil {
		for _, it := range *v {
			it_1 := it.GetImpl()
			*v_1 = append(*v_1, it_1)
		}
	}
	wrapperRcvr.impl.Dogs = v_1
}

func (wrapperRcvr *animalsWrapper) GetDogsByNameField() map[string]dog.Dog {
	retval := wrapperRcvr.impl.DogsByNameField
	var retval_1 map[string]dog.Dog
	retval_1 = map[string]dog.Dog{}
	for k, it := range retval {
		it_1 := dog.NewDog(it)
		retval_1[k] = it_1
	}
	return retval_1
}

func (wrapperRcvr *animalsWrapper) SetDogsByNameField(v map[string]dog.Dog) {
	var v_1 map[string]*orig_dog.Dog
	v_1 = map[string]*orig_dog.Dog{}
	for k, it := range v {
		it_1 := it.GetImpl()
		v_1[k] = it_1
	}
	wrapperRcvr.impl.DogsByNameField = v_1
}

func (wrapperRcvr *animalsWrapper) GetAnimalDescription() AnimalDescription {
	retval_1 := NewAnimalDescription(&wrapperRcvr.impl.AnimalDescription)
	return retval_1
}

func (wrapperRcvr *animalsWrapper) SetAnimalDescription(v AnimalDescription) {
	v_1 := *v.GetImpl()
	wrapperRcvr.impl.AnimalDescription = v_1
}

func (wrapperRcvr *animalsWrapper) GetAllDogs(ctx orig_context.Context) []dog.Dog {
	retval0 := wrapperRcvr.impl.GetAllDogs(ctx)
	var retval0_1 []dog.Dog
	for _, it := range retval0 {
		it_1 := dog.NewDog(it)
		retval0_1 = append(retval0_1, it_1)
	}
	return retval0_1
}

func (wrapperRcvr *animalsWrapper) GetDogsByNames(names []string) map[string]dog.Dog {
	retval0 := wrapperRcvr.impl.GetDogsByNames(names)
	var retval0_1 map[string]dog.Dog
	retval0_1 = map[string]dog.Dog{}
	for k, it := range retval0 {
		it_1 := dog.NewDog(it)
		retval0_1[k] = it_1
	}
	return retval0_1
}

func (wrapperRcvr *animalsWrapper) GetDogByName(name string) dog.Dog {
	retval0 := wrapperRcvr.impl.GetDogByName(name)
	retval0_1 := dog.NewDog(retval0)
	return retval0_1
}

func (wrapperRcvr *animalsWrapper) AddAnimals(animals *[]interface{}) error {
	retval0 := wrapperRcvr.impl.AddAnimals(animals)
	return retval0
}

func (wrapperRcvr *animalsWrapper) AddDogs(dogs []dog.Dog) map[string]Animals {
	var dogs_1 []orig_dog.Dog
	for _, it := range dogs {
		it_1 := *it.GetImpl()
		dogs_1 = append(dogs_1, it_1)
	}
	retval0 := wrapperRcvr.impl.AddDogs(dogs_1)
	var retval0_1 map[string]Animals
	retval0_1 = map[string]Animals{}
	for k, it := range retval0 {
		it_1 := NewAnimals(&it)
		retval0_1[k] = it_1
	}
	return retval0_1
}
