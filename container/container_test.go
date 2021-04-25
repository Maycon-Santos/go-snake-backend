package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeDependency1 = string
type fakeDependency2 = int
type fakeDependency3 = bool
type fakeDependency4 = uint8

func Test_container_Inject(t *testing.T) {
	instance := New()

	t.Run("should inject dependencies without errors", func(t *testing.T) {
		dep1 := fakeDependency1("dependency")
		dep2 := fakeDependency2(1)
		dep3 := fakeDependency3(true)

		err := instance.Inject(dep1, &dep2, dep3)
		assert.NoError(t, err)
	})

	t.Run("should return an error when receive nil", func(t *testing.T) {
		err := instance.Inject(nil, fakeDependency3(false))
		assert.EqualError(t, err, "container: dependency 0 is <nil>")
	})

	t.Run("should return an error when receive a dependency without initial value", func(t *testing.T) {
		var dep *fakeDependency3
		err := instance.Inject(fakeDependency2(1), dep)
		assert.EqualError(t, err, "container: dependency *bool is a <nil> value")
	})
}

func Test_container_Retrieve(t *testing.T) {
	dep1 := fakeDependency1("dependency")
	dep2 := fakeDependency2(1)
	dep3 := fakeDependency3(true)
	instance := New()

	instance.Inject(dep1, dep2, &dep3)

	t.Run("should retrieve dependencies without errors", func(t *testing.T) {
		var (
			abs1 fakeDependency1
			abs2 fakeDependency2
			abs3 fakeDependency3
		)

		err := instance.Retrieve(&abs1, &abs3, &abs2)
		assert.NoError(t, err)

		assert.Equal(t, fakeDependency1("dependency"), abs1)
		assert.Equal(t, fakeDependency2(1), abs2)
		assert.Equal(t, fakeDependency3(true), abs3)
	})

	t.Run("should return an error when trying retrieve abstraction without initial value", func(t *testing.T) {
		var (
			abs1 fakeDependency1
			abs2 fakeDependency2
			abs3 *fakeDependency3
		)

		err := instance.Retrieve(&abs1, abs3, &abs2)
		assert.EqualError(t, err, "container: dependency abstraction *bool is a <nil> value")
	})

	t.Run("should return an error when trying retrieve with nil value", func(t *testing.T) {
		var (
			abs1 fakeDependency1
			abs2 fakeDependency2
		)

		err := instance.Retrieve(&abs1, nil, &abs2)
		assert.EqualError(t, err, "container: dependency abstraction 1 is <nil>")
	})

	t.Run("should return an error when trying retrieve a unimplemented dependency", func(t *testing.T) {
		var (
			abs1 fakeDependency1
			abs2 fakeDependency2
			abs3 fakeDependency3
			abs4 fakeDependency4
		)

		err := instance.Retrieve(&abs1, &abs3, &abs2, &abs4)
		assert.EqualError(t, err, "container: dependency uint8 has not been implemented")
	})

	t.Run("should return an error when trying retrieve by passing a non-pointer", func(t *testing.T) {
		var (
			abs1 fakeDependency1
			abs2 fakeDependency2
		)

		err := instance.Retrieve(&abs1, abs2)
		assert.EqualError(t, err, "container: dependency abstraction int is not a pointer")
	})
}
