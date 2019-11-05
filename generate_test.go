package destructor

import (
	"github.com/alexodle/wrappy/testdata/generated/animals"
	orig_animals "github.com/alexodle/wrappy/testdata/input/animals"
	orig_dog "github.com/alexodle/wrappy/testdata/input/animals/dog"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_Generate(t *testing.T) {
	require.NoError(t, os.RemoveAll("testdata/generated"))
	GenerateWrappers("testdata/input", "testdata/generated")
}

func Test_GeneratedCode(t *testing.T) {
	desc := animals.NewAnimalDescription(&orig_animals.AnimalDescription{
		Breed:  "Labrador",
		Name:   "Rooney",
		Weight: 60,
	})
	require.Equal(t, "Labrador", desc.GetBreed())
	desc.SetBreed("Shepherd")
	require.Equal(t, "Shepherd", desc.GetBreed())

	dog1 := &orig_dog.Dog{Name: "Dog1"}
	animal := animals.NewAnimals(&orig_animals.Animals{
		Locations:       []orig_animals.Location{{}},
		Dogs:            &[]*orig_dog.Dog{dog1},
		DogsByNameField: map[string]*orig_dog.Dog{dog1.Name: dog1},
		AnimalDescription: &orig_animals.AnimalDescription{
			Breed:  "Labrador",
			Name:   "Rooney",
			Weight: 60,
		},
	})

	d := animal.GetDogByName(dog1.Name)
	require.Equal(t, dog1, d.GetImpl())

	d = animal.GetDogsByNameField()[dog1.Name]
	require.Equal(t, dog1, d.GetImpl())
}
