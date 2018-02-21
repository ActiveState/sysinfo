package sysinfo

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
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
	}
	log.Printf("OSVersion: unsupported OS '%s'", os)
	return ""
}

// Architecture returns the system's architecture (e.g. "amd64", "386", etc.).
func Architecture() string {
	if architectureOverride != "" {
		return architectureOverride
	}
	os := OS()
	if os == "linux" {
		archBytes, err := exec.Command("uname", "-m").Output()
		if err != nil {
			log.Printf("Architecture: unable to run 'uname -m': %s", err)
			return ""
		}
		arch := string(archBytes)
		if strings.HasSuffix(arch, "64") {
			return "amd64"
		} else if strings.HasPrefix(arch, "i") {
			return "386"
		} else if strings.HasPrefix(arch, "arm") {
			return "arm"
		}
		log.Printf("Architecture: unknown architecture '%s'", arch)
		return ""
	} else if os == "windows" {
		// Development on Windows is either amd64 or 386.
		// Use a trick suggested on the golang-nuts mailing list to get integer
		// size.
		const intSize = 32 + int(^uintptr(0)>>63<<5)
		if intSize == 64 {
			return "amd64"
		}
		return "386"
	}
	log.Printf("Architecture: unsupported OS '%s'", os)
	return ""
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
