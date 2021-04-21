package container

import (
	"fmt"
	"reflect"
)

type Container map[reflect.Type]reflect.Value

func New() Container {
	return make(Container)
}

func (c Container) Inject(dependencies ...interface{}) error {
	for i, dependency := range dependencies {
		dependencyType := reflect.TypeOf(dependency)
		if dependencyType == nil {
			return fmt.Errorf("container: dependency %d is %v", i, dependencyType)
		}

		dependencyValue := reflect.ValueOf(dependency)

		if dependencyType.Kind() == reflect.Ptr {
			c[dependencyType.Elem()] = dependencyValue.Elem()
			continue
		}

		c[dependencyType] = dependencyValue
	}

	return nil
}

func (c Container) Retrieve(dependenciesAbstraction ...interface{}) error {
	for i, dependencyAbstraction := range dependenciesAbstraction {
		abstractionType := reflect.TypeOf(dependencyAbstraction)
		if abstractionType == nil {
			return fmt.Errorf("container: dependency abstraction %d is %v", i, abstractionType)
		}

		if abstractionType.Kind() == reflect.Ptr {
			abstractionElem := abstractionType.Elem()

			if dependency, ok := c[abstractionElem]; ok {
				reflect.ValueOf(dependencyAbstraction).Elem().Set(dependency)
				continue
			}

			return fmt.Errorf("container: dependency %s is not implemented", abstractionElem)
		}

		return fmt.Errorf("container: dependency abstraction %v is not a pointer", abstractionType)
	}

	return nil
}
