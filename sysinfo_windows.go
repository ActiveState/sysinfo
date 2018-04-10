package sysinfo

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"golang.org/x/sys/windows"
)

// OS returns the system's OS
func OS() OsInfo {
	return Windows
}

// From https://msdn.microsoft.com/en-us/library/windows/desktop/ms724832%28v=vs.85%29.aspx
// Note: cannot differentiate between some versions, hence the '/'. Also, unless
// the program using this package is "manifested" (see above link), Windows will
// not report higher than 6.2 (Windows 8 / Windows Server 2012).
var versions = map[int]map[int]string{
	5: map[int]string{
		0: "Windows 2000",
		1: "Windows XP",
		2: "Windows XP / Windows Server 2003",
	},
	6: map[int]string{
		0: "Windows Vista / Windows Server 2008",
		1: "Windows 7 / Windows Server 2008 R2",
		2: "Windows 8 / Windows Server 2012",
		3: "Windows 8.1 / Windows Server 2012 R2",
	},
	10: map[int]string{
		0: "Windows 10 / Windows Server 2016",
	},
}

// OSVersion returns the system's OS version.
func OSVersion() (*OSVersionInfo, error) {
	dll := windows.NewLazySystemDLL("kernel32.dll")
	version, _, _ := dll.NewProc("GetVersion").Call()
	major := int(byte(version))
	minor := int(uint8(version >> 8))
	micro := int(uint16(version >> 16))
	name := "Unknown"
	if subversion, ok := versions[major]; ok {
		if value, ok := subversion[minor]; ok {
			name = value
		}
	}
	return &OSVersionInfo{
		fmt.Sprintf("%d.%d.%d", major, minor, micro),
		major,
		minor,
		micro,
		name,
	}, nil
}

// Libc returns the system's C library.
func Libc() (*LibcInfo, error) {
	// Use Windows powershell in order to query the version information from
	// msvcrt.dll. This works on Windows 7 and higher.
	// Note: cannot easily use version.dll's GetFileVersionInfo function since its
	// return value is a pointer and VerQueryValue is needed in order to fetch a C
	// struct with version info. Manipulating C structs with Go is an exercise in
	// patience.
	windir := os.Getenv("SYSTEMROOT")
	if windir == "" {
		return nil, errors.New("Unable to find system root")
	}
	msvcrt := filepath.Join(windir, "System32", "msvcrt.dll")
	if _, err := os.Stat(msvcrt); err != nil {
		return nil, nil // no libc found
	}
	versionInfo, err := exec.Command("powershell", "-command", "(Get-Item "+msvcrt+").VersionInfo").Output()
	if err != nil {
		return nil, fmt.Errorf("Unable to determine libc version: %s", err)
	}
	regex := regexp.MustCompile("(\\d+)\\D(\\d+)")
	parts := regex.FindStringSubmatch(string(versionInfo))
	if len(parts) != 3 {
		return nil, fmt.Errorf("Unable to parse versionInfo string '%s'", versionInfo)
	}
	major, _ := strconv.Atoi(parts[1])
	minor, _ := strconv.Atoi(parts[2])
	return &LibcInfo{Msvcrt, major, minor}, nil
}

// Compilers returns the system's available compilers.
func Compilers() ([]*CompilerInfo, error) {
	compilers := []*CompilerInfo{}

	// Map of compiler commands to CompilerNameInfos.
	var compilerMap = map[string]CompilerNameInfo{
		"cl": Msvc,
	}
	for command, nameInfo := range compilerMap {
		major, minor, err := getCompilerVersion([]string{command})
		if err != nil {
			return compilers, err
		} else if major > 0 {
			compilers = append(compilers, &CompilerInfo{nameInfo, major, minor})
		}
	}

	return compilers, nil
}
