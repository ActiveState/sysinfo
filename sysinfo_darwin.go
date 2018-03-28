package sysinfo

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

// OS returns the system's OS
func OS() OsInfo {
	return Mac
}

// OSVersion returns the system's OS version.
func OSVersion() (OSVersionInfo, error) {
	// Fetch OS version.
	version, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return OSVersionInfo{}, fmt.Errorf("Unable to determine OS version: %s", err)
	}
	version = bytes.TrimSpace(version)
	// Parse OS version parts.
	regex := regexp.MustCompile("^(\\d+)\\D(\\d+)\\D(\\d+)")
	parts := regex.FindStringSubmatch(string(version))
	if len(parts) != 4 {
		return OSVersionInfo{}, fmt.Errorf("Unable to parse version string '%s'", version)
	}
	for i := 1; i < len(parts); i++ {
		if _, err := strconv.Atoi(parts[i]); err != nil {
			return OSVersionInfo{}, fmt.Errorf("Unable to parse part '%s' of version string '%s'", parts[i], version)
		}
	}
	major, _ := strconv.Atoi(parts[1])
	minor, _ := strconv.Atoi(parts[2])
	micro, _ := strconv.Atoi(parts[3])
	// Fetch OS name.
	name, err := exec.Command("sw_vers", "-productName").Output()
	return OSVersionInfo{string(version), major, minor, micro, string(name)}, nil
}

// Libc returns the system's C library.
func Libc() (LibcInfo, error) {
	version, err := exec.Command("clang", "--version")
	if err != nil {
		return LibcInfo{}, fmt.Errorf("Unable to fetch libc version: %s", err)
	}
	regex := regexp.MustCompile("(\\d+)\\D(\\d+)")
	parts := regex.FindStringSubmatch(string(version))
	if len(parts) != 3 {
		return LibcInfo{}, fmt.Errorf("Unable to parse libc string '%s'", libc)
	}
	for i := 1; i < len(parts); i++ {
		if _, err := strconv.Atoi(parts[i]); err != nil {
			return LibcInfo{}, fmt.Errorf("Unable to parse part '%s' of libc string '%s'", parts[i], libc)
		}
	}
	major, _ := strconv.Atoi(parts[1])
	minor, _ := strconv.Atoi(parts[2])
	return LibcInfo{BsdLibc, major, minor}, nil
}

// Compilers returns the system's available compilers.
func Compilers() ([]CompilerInfo, error) {
	compilers := []CompilerInfo{}
	major, minor, err := getCompilerVersion([]string{"clang", "--version"})
	if err != nil {
		return compilers, err
	} else if major > 0 {
		compilers = append(compilers, CompilerInfo{nameInfo, major, minor})
	}

	return compilers, nil
}
