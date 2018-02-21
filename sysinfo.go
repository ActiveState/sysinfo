package sysinfo

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"syscall"
)

// Mainly for testing.
var osOverride, osVersionOverride, architectureOverride, libcOverride, compilerOverride string

// OS returns the system's OS name (e.g. "linux", "windows", etc.)
func OS() string {
	if osOverride != "" {
		return osOverride
	}
	return runtime.GOOS
}

// OSVersion returns the system's OS version.
func OSVersion() string {
	if osVersionOverride != "" {
		return osVersionOverride
	}
	os := OS()
	if os == "linux" {
		version, err := exec.Command("uname", "-r").Output()
		if err != nil {
			log.Printf("OSVersion: unable run 'uname -r': %s", err)
			return ""
		}
		return string(version)
	} else if os == "windows" {
		dll, err := syscall.LoadDLL("kernel32.dll")
		if err != nil {
			log.Println("OSVersion: cannot find 'kernel32.dll'")
			return ""
		}
		proc, err := dll.FindProc("GetVersion")
		if err != nil {
			log.Println("OSVersion: cannot find 'GetVersion' in 'kernel32.dll'")
			return ""
		}
		version, _, _ := proc.Call()
		return fmt.Sprintf("%d.%d.%d", byte(version), uint8(version>>8), uint16(version>>16))
	} else {
		log.Printf("OSVersion: unsupported OS '%s'", os)
	}
	return ""
}

// Architecture returns the system's architecture (e.g. "amd64", "386", etc.).
func Architecture() string {
	if architectureOverride != "" {
		return architectureOverride
	}
	return runtime.GOARCH
}

// Libc returns the system's libc (e.g. glibc, msvc, etc.) version.
func Libc() string {
	if libcOverride != "" {
		return libcOverride
	}
	return "" // TODO
}

// Compiler returns the system's compiler (e.g. gcc, msvc, etc.) version.
func Compiler() string {
	if compilerOverride != "" {
		return compilerOverride
	}
	return "" // TODO
}
