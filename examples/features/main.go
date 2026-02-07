// features queries and prints the capabilities of the installed libkrun library.
//
// Usage:
//
//	go run .
package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/mishushakov/libkrun-go/krun"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Query max vCPUs supported by the hypervisor.
	maxVCPUs, err := krun.GetMaxVCPUs()
	if err != nil {
		return fmt.Errorf("get max vcpus: %w", err)
	}
	fmt.Printf("Max vCPUs: %d\n", maxVCPUs)

	// Check nested virtualization support (macOS only).
	nested, err := krun.CheckNestedVirt()
	if err != nil {
		// Not an error on Linux — just unsupported.
		if errors.Is(err, syscall.ENOSYS) {
			fmt.Println("Nested virtualization: not supported on this platform")
		} else {
			return fmt.Errorf("check nested virt: %w", err)
		}
	} else {
		fmt.Printf("Nested virtualization: %v\n", nested)
	}

	// Query compile-time features.
	features := []struct {
		name    string
		feature krun.Feature
	}{
		{"Networking", krun.FeatureNet},
		{"Block devices", krun.FeatureBLK},
		{"GPU", krun.FeatureGPU},
		{"Sound", krun.FeatureSND},
		{"Input", krun.FeatureInput},
		{"EFI", krun.FeatureEFI},
		{"TEE", krun.FeatureTEE},
		{"AMD SEV", krun.FeatureAMDSEV},
		{"Intel TDX", krun.FeatureIntelTDX},
		{"AWS Nitro", krun.FeatureAWSNitro},
		{"Virgl Resource Map2", krun.FeatureVirglResourceMap2},
	}

	fmt.Println("\nCompile-time features:")
	for _, f := range features {
		supported, err := krun.HasFeature(f.feature)
		if err != nil {
			// EINVAL means the feature constant is unknown to this libkrun version.
			if errors.Is(err, syscall.EINVAL) {
				fmt.Printf("  %-25s unknown (older libkrun?)\n", f.name)
			} else {
				fmt.Printf("  %-25s error: %v\n", f.name, err)
			}
			continue
		}
		status := "no"
		if supported {
			status = "yes"
		}
		fmt.Printf("  %-25s %s\n", f.name, status)
	}

	return nil
}
