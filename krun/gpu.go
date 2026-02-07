package krun

/*
#include <libkrun.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// SetGPUOptions enables and configures a virtio-gpu device.
// virglFlags is a bitmask of Virgl* flags.
func (c *Context) SetGPUOptions(virglFlags uint32) error {
	return checkRet(
		C.krun_set_gpu_options(C.uint32_t(c.id), C.uint32_t(virglFlags)),
		"krun_set_gpu_options",
	)
}

// SetGPUOptions2 enables and configures a virtio-gpu device with a custom
// host SHM window size (acting as vRAM in the guest).
func (c *Context) SetGPUOptions2(virglFlags uint32, shmSize uint64) error {
	return checkRet(
		C.krun_set_gpu_options2(C.uint32_t(c.id), C.uint32_t(virglFlags), C.uint64_t(shmSize)),
		"krun_set_gpu_options2",
	)
}

// AddDisplay configures a display output for the VM.
// A display backend must also be set via [Context.SetDisplayBackend].
// Returns the display ID (0 to [MaxDisplays]-1) on success.
func (c *Context) AddDisplay(width, height uint32) (uint32, error) {
	ret := C.krun_add_display(C.uint32_t(c.id), C.uint32_t(width), C.uint32_t(height))
	if ret < 0 {
		return 0, retError(ret, "krun_add_display")
	}
	return uint32(ret), nil
}

// DisplaySetEDID configures a custom EDID blob for a display.
// This replaces the generated EDID. libkrun does not verify the EDID
// matches the width/height from [Context.AddDisplay].
func (c *Context) DisplaySetEDID(displayID uint32, edidBlob []byte) error {
	return checkRet(
		C.krun_display_set_edid(
			C.uint32_t(c.id), C.uint32_t(displayID),
			(*C.uint8_t)(unsafe.Pointer(&edidBlob[0])), C.size_t(len(edidBlob)),
		),
		"krun_display_set_edid",
	)
}

// DisplaySetDPI configures the DPI of a display reported to the guest.
func (c *Context) DisplaySetDPI(displayID, dpi uint32) error {
	return checkRet(
		C.krun_display_set_dpi(C.uint32_t(c.id), C.uint32_t(displayID), C.uint32_t(dpi)),
		"krun_display_set_dpi",
	)
}

// DisplaySetPhysicalSize sets the physical display dimensions reported to the guest.
func (c *Context) DisplaySetPhysicalSize(displayID uint32, widthMM, heightMM uint16) error {
	return checkRet(
		C.krun_display_set_physical_size(
			C.uint32_t(c.id), C.uint32_t(displayID),
			C.uint16_t(widthMM), C.uint16_t(heightMM),
		),
		"krun_display_set_physical_size",
	)
}

// DisplaySetRefreshRate configures the refresh rate for a display (in Hz).
func (c *Context) DisplaySetRefreshRate(displayID, refreshRate uint32) error {
	return checkRet(
		C.krun_display_set_refresh_rate(
			C.uint32_t(c.id), C.uint32_t(displayID), C.uint32_t(refreshRate),
		),
		"krun_display_set_refresh_rate",
	)
}

// SetDisplayBackend configures the display backend.
// displayBackend should point to a krun_display_backend struct (from libkrun_display.h).
func (c *Context) SetDisplayBackend(displayBackend unsafe.Pointer, backendSize uintptr) error {
	return checkRet(
		C.krun_set_display_backend(C.uint32_t(c.id), displayBackend, C.size_t(backendSize)),
		"krun_set_display_backend",
	)
}

// AddInputDeviceFD creates a passthrough input device from a host /dev/input/* file descriptor.
func (c *Context) AddInputDeviceFD(inputFD int) error {
	return checkRet(
		C.krun_add_input_device_fd(C.uint32_t(c.id), C.int(inputFD)),
		"krun_add_input_device_fd",
	)
}

// AddInputDevice adds an input device with separate config and events backends.
// configBackend should point to a krun_input_config struct.
// eventsBackend should point to a krun_input_event_provider struct.
func (c *Context) AddInputDevice(configBackend unsafe.Pointer, configSize uintptr, eventsBackend unsafe.Pointer, eventsSize uintptr) error {
	ret := C.krun_add_input_device(
		C.uint32_t(c.id),
		configBackend, C.size_t(configSize),
		eventsBackend, C.size_t(eventsSize),
	)
	return checkRet(C.int32_t(ret), "krun_add_input_device")
}

// SetSndDevice enables or disables the virtio-snd device.
func (c *Context) SetSndDevice(enable bool) error {
	return checkRet(
		C.krun_set_snd_device(C.uint32_t(c.id), C.bool(enable)),
		"krun_set_snd_device",
	)
}
