package reg

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[reflect.Type]map[string]any)
)

func Types() (ret []string) {
	driversMu.RLock()
	for typ := range drivers {
		ret = append(ret, typ.String())
	}
	driversMu.RUnlock()
	return
}

func Drivers[T any]() (ret []string) {
	typ := reflect.TypeOf((*T)(nil)).Elem()

	driversMu.RLock()
	registry, ok := drivers[typ]
	if !ok {
		registry = make(map[string]any)
		drivers[typ] = registry
	}
	for name := range registry {
		ret = append(ret, name)
	}
	driversMu.RUnlock()

	return
}

func Register[T any](name string, driver T) {
	driversMu.Lock()
	defer driversMu.Unlock()

	typ := reflect.TypeOf((*T)(nil)).Elem()

	registry, ok := drivers[typ]
	if !ok {
		registry = make(map[string]any)
		drivers[typ] = registry
	}

	if _, dup := registry[name]; dup {
		panic("reg: register called twice for driver " + name)
	}
	registry[name] = driver
}

func Open[T any](name string) (ret T, err error) {
	typ := reflect.TypeOf((*T)(nil)).Elem()

	driversMu.RLock()
	registry, ok := drivers[typ]
	if !ok {
		return ret, fmt.Errorf("reg: unknown type %q (forgotten import?)", typ)
	}
	driver, ok := registry[name]
	driversMu.RUnlock()
	if !ok {
		return ret, fmt.Errorf("reg: unknown driver %q (forgotten import?)", name)
	}

	return driver.(T), nil
}
