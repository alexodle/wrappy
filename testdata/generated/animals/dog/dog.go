package dog

import food "github.com/alexodle/wrappy/testdata/generated/animals/food"
import orig_dog "github.com/alexodle/wrappy/testdata/input/animals/dog"

type Dog interface {
	GetImpl() *orig_dog.Dog
	Barks() bool
	Clone() (Dog, error)
	Eat(f food.Food) int
	GetName() string
	Meows() bool
	SetName(v string)
}

func NewDog(impl *orig_dog.Dog) Dog {
	return &dogWrapper{impl: impl}
}

type dogWrapper struct {
	impl *orig_dog.Dog
}

func (o *dogWrapper) GetImpl() *orig_dog.Dog {
	return o.impl
}

func (o *dogWrapper) GetName() string {
	retval := o.impl.Name
	return retval
}

func (o *dogWrapper) SetName(v string) {
	o.impl.Name = v
}

func (o *dogWrapper) Barks() bool {
	retval0 := o.impl.Barks()
	return retval0
}

func (o *dogWrapper) Meows() bool {
	retval0 := o.impl.Meows()
	return retval0
}

func (o *dogWrapper) Eat(f food.Food) int {
	f_1 := f.GetImpl()
	retval0 := o.impl.Eat(f_1)
	return retval0
}

func (o *dogWrapper) Clone() (Dog, error) {
	retval0, retval1 := o.impl.Clone()
	retval0_1 := NewDog(retval0)
	return retval0_1, retval1
}