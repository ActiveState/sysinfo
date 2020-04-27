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
	"golang.org/x/sys/windows/registry"
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
		0: "Windows 10 / Windows Server",
	},
}

func winVersionName(major, minor, micro int) string {
	name := "Unknown"
	if subversions, ok := versions[major]; ok {
		if value, ok := subversions[minor]; ok {
			name = value
		}
	}
	return name
}

func newOSVersionInfo(major, minor, micro int) *OSVersionInfo {
	return &OSVersionInfo{
		fmt.Sprintf("%d.%d.%d", major, minor, micro),
		major,
		minor,
		micro,
		winVersionName(major, minor, micro),
	}
}

// OSVersion returns the system's OS version.
func OSVersion() (*OSVersionInfo, error) {
	major, minor, micro, err := getWindowsVersionFromRegistry()
	if err == nil {
		return newOSVersionInfo(major, minor, micro), nil
	}
	regErr := err

	major, minor, micro, err = getWindowsVersionFromDLL()
	if err != nil {
		fmt.Sprintf("From DLL error: %v. From Registry error: %v", err, regErr)
	}

	return newOSVersionInfo(major, minor, micro), nil
}

func getWindowsVersionFromRegistry() (major, minor, micro int, err error) {
	keyName := `SOFTWARE\Microsoft\Windows NT\CurrentVersion`
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, keyName, registry.QUERY_VALUE)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("Cannot open registry key %q: %w", keyName, err)
	}
	defer key.Close()

	keyEntryErrMsgFmt := "Cannot get entry %q at %q: %w"

	majorEntryName := "CurrentMajorVersionNumber"
	major64, _, err := key.GetIntegerValue(majorEntryName)
	if err != nil {
		return 0, 0, 0, fmt.Errorf(keyEntryErrMsgFmt, majorEntryName, keyName, err)
	}
	major = int(major64)

	minorEntryName := "CurrentMinorVersionNumber"
	minor64, _, err := key.GetIntegerValue(minorEntryName)
	if err != nil {
		return major, 0, 0, fmt.Errorf(keyEntryErrMsgFmt, minorEntryName, keyName, err)
	}
	minor = int(minor64)

	microEntryName := "CurrentBuild"
	microText, _, err := key.GetStringValue(microEntryName)
	if err != nil {
		return major, minor, 0, fmt.Errorf(keyEntryErrMsgFmt, microEntryName, keyName, err)
	}
	micro, err = strconv.Atoi(microText)
	if err != nil {
		atoiErr := fmt.Errorf("Cannot convert %q text to integer: %w", microEntryName, err)
		return major, minor, 0, atoiErr
	}

	return major, minor, micro, nil
}

func getWindowsVersionFromDLL() (major, minor, micro int, err error) {
	dll := windows.NewLazySystemDLL("kernel32.dll")
	version, _, err := dll.NewProc("GetVersion").Call()
	if err != nil {
		return 0, 0, 0, fmt.Errorf("'GetVersion' via kernel32.dll failed: %w")
	}

	major = int(byte(version))
	minor = int(uint8(version >> 8))
	micro = int(uint16(version >> 16))

	return major, minor, micro, nil
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
		return nil, errors.New("Unable to find system root; %SYSTEMROOT% undefined")
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
		"gcc.exe": Mingw,
	}
	// Search for MSVC locations and add their C++ compilers to the map.
	if key, err := registry.OpenKey(registry.LOCAL_MACHINE, "Software\\Wow6432Node\\Microsoft\\VisualStudio\\SxS\\VS7", registry.QUERY_VALUE); err == nil {
		// This registry technique works for all MSVC prior to 15.0 (VS2017).
		if valueNames, err := key.ReadValueNames(0); err == nil {
			for _, name := range valueNames {
				if _, err := strconv.ParseFloat(name, 32); err != nil {
					continue
				}
				path, _, err := key.GetStringValue(name)
				cl := filepath.Join(path, "VC", "bin", "cl.exe")
				if _, err = os.Stat(cl); err == nil {
					compilerMap[cl] = Msvc
				}
			}
		}
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
