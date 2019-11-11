package wrappy

import (
	"github.com/alexodle/wrappy/testdata/generated/animals"
	"github.com/alexodle/wrappy/testdata/generated/animals/util"
	orig_animals "github.com/alexodle/wrappy/testdata/input/animals"
	orig_util "github.com/alexodle/wrappy/testdata/input/animals/util"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_Generate(t *testing.T) {
	require.NoError(t, os.RemoveAll("testdata/generated"))
	GenerateWrappers("testdata/input", "testdata/generated", GenerateWrappersOptions{})
}

func Test_GeneratedCode(t *testing.T) {
	impl := orig_animals.NewAnimalStore()
	st := animals.NewAnimalStore(impl)

	require.Len(t, st.GetAllAnimals(), 0)
	_, ok := st.GetAnimalByName("Rover")
	require.False(t, ok)

	rover := &orig_animals.Animal{AnimalDescription: orig_animals.AnimalDescription{
		Name:   "Rover",
		Age:    4,
		Type:   "Dog",
		Weight: 60,
	}}
	st.AddAnimal(animals.NewAnimal(rover))

	rover_wrapper, ok := st.GetAnimalByName("Rover")
	require.True(t, ok)
	require.Equal(t, rover.Name, rover_wrapper.GetAnimalDescription().GetName())
	require.Equal(t, rover.Age, rover_wrapper.GetAnimalDescription().GetAge())

	rover_wrapper.GetAnimalDescription().SetAge(5)
	require.Equal(t, 5, rover_wrapper.GetAnimalDescription().GetAge())
	require.Equal(t, 5, rover.Age)

	allAnimals := st.GetAllAnimals()
	require.Len(t, allAnimals, 1)
	require.Equal(t, allAnimals[0].GetImpl(), rover)

	charlotte := &orig_animals.Animal{AnimalDescription: orig_animals.AnimalDescription{
		Name:   "Charlotte",
		Age:    7,
		Type:   "Pig",
		Weight: 80,
	}}
	allAnimals = append(allAnimals, animals.NewAnimal(charlotte))

	util.NewUtils(orig_util.NewUtils(st.GetImpl())).AddAllAnimals(allAnimals)
	allAnimals = st.GetAllAnimals()
	require.Len(t, allAnimals, 2)
}
