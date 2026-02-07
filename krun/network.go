package krun

/*
#include <libkrun.h>
#include <stdlib.h>
*/
import "C"

// SetPortMap configures host-to-guest TCP port mappings.
// Each entry has the format "host_port:guest_port".
//
// Pass nil to expose all listening guest ports to the host.
// Pass an empty slice to expose no ports.
func (c *Context) SetPortMap(portMap []string) error {
	cArr := stringsToCArray(portMap)
	if portMap != nil {
		defer freeCStringArray(cArr, len(portMap))
	}
	return checkRet(C.krun_set_port_map(C.uint32_t(c.id), cArr), "krun_set_port_map")
}
