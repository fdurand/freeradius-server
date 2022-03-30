package main

// #include <freeradius-devel/radiusd.h>
// #include <freeradius-devel/modules.h>
// #include <freeradius-devel/rad_assert.h>
import "C"

import (
	"fmt"
	"unsafe"
	//"plugin"
	"sync"

	"github.com/fdurand/freeradius-go"
	"github.com/davecgh/go-spew/spew"
)

const (
	createModuleSymbol = "CreateModule"
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
func go_instantiate(cconf *C.CONF_SECTION, instance unsafe.Pointer) C.int {
	fmt.Println("go_instantiate called!")
	instance = unsafe.Pointer(&C.struct_rlm_go_cache_eap_t{})
	return 0
}

//export go_instantiate
//func go_instantiate(cconf *C.CONF_SECTION, pl *C.char) C.int {
//	redisServer := C.GoString(pl)
//	radlogInstance.Radlog(freeradius.LogTypeInfo, "using redis server %s", redisServer)
	//gomodule, err := plugin.Open(pluginPath)
	//if err != nil || gomodule == nil {
	//	radlogInstance.Radlog(freeradius.LogTypeError, "Failed to load plugin %s: %#v", pluginPath, err)
	//	return -1
	//}

	//radlogInstance.Radlog(freeradius.LogTypeInfo, "Looking up plugin symbol CreateModule")
	//createModule, err := gomodule.Lookup(createModuleSymbol)
	//if err != nil {
	//	radlogInstance.Radlog(freeradius.LogTypeError, "Unable to lookup symbol %s: %#v", createModuleSymbol, err)
	//	return -1
	//}

//	radlogInstance.Radlog(freeradius.LogTypeInfo, "Calling CreateModule")
//	instance := createModule.(freeradius.ModuleFunc)()
//	if instance == nil {
//		radlogInstance.Radlog(freeradius.LogTypeError, "Created go module instance is nil")
//		return -1
//	}

//	radlogInstance.Radlog(freeradius.LogTypeInfo, "Initiating go plugin")
//	if err := instance.Init(radlogInstance); err != nil {
//		radlogInstance.Radlog(freeradius.LogTypeError, "Unable to initialize go module %s: %#v", redisServer, err)
//	}

//	insertInstance(redisServer, instance)
//	return 0
//}

//export go_authorize
func go_authorize(instancePath *C.char, request *C.REQUEST) C.int {

	mod := getInstanceC(instancePath)
	if mod == nil {
		// Should never happen
		return toRlmCodeT(freeradius.RlmCodeFail)
	}

	req := NewRequest(request)
	spew.Dump(request)
	return toRlmCodeT(mod.Authorize(req))
}

func toRlmCodeT(in freeradius.RlmCode) C.int {
	return C.int(in)
}

func main() {}
