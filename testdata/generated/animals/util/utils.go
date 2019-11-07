package util

import animals "github.com/alexodle/wrappy/testdata/generated/animals"
import orig_animals "github.com/alexodle/wrappy/testdata/input/animals"
import orig_util "github.com/alexodle/wrappy/testdata/input/animals/util"

type Utils interface {
	GetImpl() *orig_util.Utils
	AddAllAnimals(animals []animals.Animal)
	GetStore() animals.AnimalStore
	SetStore(v animals.AnimalStore)
}

func NewUtils(impl *orig_util.Utils) Utils {
	return &utilsWrapper{impl: impl}
}

type utilsWrapper struct {
	impl *orig_util.Utils
}

func (o *utilsWrapper) GetImpl() *orig_util.Utils {
	return o.impl
}

func (wrapperRcvr *utilsWrapper) GetStore() animals.AnimalStore {
	retval := wrapperRcvr.impl.Store
	retval_1 := animals.NewAnimalStore(retval)
	return retval_1
}

func (wrapperRcvr *utilsWrapper) SetStore(v animals.AnimalStore) {
	v_1 := v.GetImpl()
	wrapperRcvr.impl.Store = v_1
}

func (wrapperRcvr *utilsWrapper) AddAllAnimals(animals []animals.Animal) {
	var animals_1 []*orig_animals.Animal
	for _, it := range animals {
		it_1 := it.GetImpl()
		animals_1 = append(animals_1, it_1)
	}
	wrapperRcvr.impl.AddAllAnimals(animals_1)
}
