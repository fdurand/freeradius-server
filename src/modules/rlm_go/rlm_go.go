package main

// #include <freeradius-devel/radiusd.h>
// #include <freeradius-devel/modules.h>
// #include <freeradius-devel/rad_assert.h>
import "C"

import (
	//"fmt"
	//"unsafe"

	"sync"

	"github.com/fdurand/freeradius-go"
	modules "github.com/fdurand/freeradius-go-modules"
)

var (
	instanceMap = make(map[string]freeradius.Module)
	mapLock     = &sync.Mutex{}
)

func insertInstance(path string, mod freeradius.Module) {
	mapLock.Lock()
	defer mapLock.Unlock()
	instanceMap[path] = mod
}

func getInstance(path string) freeradius.Module {
	mapLock.Lock()
	defer mapLock.Unlock()
	return instanceMap[path]
}

func getInstanceC(cs *C.char) freeradius.Module {
	path := C.GoString(cs)
	return getInstance(path)
}

//export go_instantiate
func go_instantiate(cconf *C.CONF_SECTION, pl *C.char) C.int {
	plugin := C.GoString(pl)
	instance, _ := modules.Create(plugin, "bob")
	radlogInstance.Radlog(freeradius.LogTypeInfo, "Initiating go plugin")
	if err := instance.Init(radlogInstance); err != nil {
		radlogInstance.Radlog(freeradius.LogTypeError, "Unable to initialize go module %s: %#v", plugin, err)
	}
	insertInstance(plugin, instance)
	return 0
}

//export go_authorize
func go_authorize(instancePath *C.char, request *C.REQUEST) C.int {

	mod := getInstanceC(instancePath)
	if mod == nil {
		// Should never happen
		return toRlmCodeT(freeradius.RlmCodeFail)
	}

	req := NewRequest(request)

	return toRlmCodeT(mod.Authorize(req))
}

func toRlmCodeT(in freeradius.RlmCode) C.int {
	return C.int(in)
}

func main() {}
