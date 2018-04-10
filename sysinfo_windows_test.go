package sysinfo

import (
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/sys/windows/registry"
)

// By default, an MSVC C++ compiler is probably not on %PATH%. In order for
// tests to pass, search for a MSVC installation and add its C++ compiler to
// %PATH%.
// Note: Go's exec.Command() looks in both %PATH% and %PATHEXT%. Set the latter
// in order to avoid modifying %PATH% directly.
func init() {
	if key, err := registry.OpenKey(registry.LOCAL_MACHINE, "Software\\Wow6432Node\\Microsoft\\VisualStudio\\SxS\\VS7", registry.QUERY_VALUE); err == nil {
		// This registry technique works for all MSVC up to 15.0 (VS2017).
		if valueNames, err := key.ReadValueNames(0); err == nil {
			for _, name := range valueNames {
				if _, err := strconv.ParseFloat(name, 32); err == nil {
					path, _, err := key.GetStringValue(name)
					if _, err = os.Stat(filepath.Join(path, "VC", "bin", "cl.exe")); err == nil {
						os.Setenv("PATHEXT", filepath.Join(path, "VC", "bin"))
					}
				}
			}
		}
	}
}
