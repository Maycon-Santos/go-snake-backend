package container

import (
	"fmt"
	"reflect"
)

type Container interface {
	Inject(dependencies ...interface{}) error
	Retrieve(dependenciesAbstraction ...interface{}) error
}

type container map[reflect.Type]reflect.Value

func New() Container {
	return make(container)
}

func (c container) Inject(dependencies ...interface{}) error {
	for i, dependency := range dependencies {
		dependencyType := reflect.TypeOf(dependency)
		dependencyValue := reflect.ValueOf(dependency)

		if dependencyType == nil {
			return fmt.Errorf("container: dependency %d is <nil>", i)
		}

		if dependencyType.Kind() == reflect.Ptr {
			if dependencyValue.IsNil() {
				return fmt.Errorf("container: dependency %v is a <nil> value", dependencyType)
			}

			c[dependencyType.Elem()] = dependencyValue.Elem()
			continue
		}

		return fmt.Errorf("container: dependency %v is not a pointer", dependencyType)
	}

	return nil
}

func (c container) Retrieve(dependenciesAbstraction ...interface{}) error {
	for i, dependencyAbstraction := range dependenciesAbstraction {
		abstractionType := reflect.TypeOf(dependencyAbstraction)
		abstractionValue := reflect.ValueOf(dependencyAbstraction)

		if abstractionType == nil {
			return fmt.Errorf("container: dependency abstraction %d is <nil>", i)
		}

		if abstractionType.Kind() == reflect.Ptr {
			if abstractionValue.IsNil() {
				return fmt.Errorf("container: dependency abstraction %v is a <nil> value", abstractionType)
			}

			abstractionElem := abstractionType.Elem()

			if dependency, ok := c[abstractionElem]; ok {
				abstractionValue.Elem().Set(dependency)
				continue
			}

			if retrieveHard(c, abstractionValue) {
				continue
			}

			return fmt.Errorf("container: dependency %s has not been implemented", abstractionElem)
		}

		return fmt.Errorf("container: dependency abstraction %v is not a pointer", abstractionType)
	}

	return nil
}

func retrieveHard(c container, abstractionValue reflect.Value) bool {
	for _, dependency := range c {
		found := func() bool {
			defer func() {
				recover()
			}()
			abstractionValue.Elem().Set(dependency)
			return true
		}()

		if found {
			return true
		}
	}

	return false
}
