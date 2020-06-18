package main

import (
	"fmt"
	"log"

	"github.com/ActiveState/sysinfo"
)

func main() {
	fmt.Printf("OS Name: %s\n", sysinfo.OS().String())
	osVersion, err := sysinfo.OSVersion()
	if err != nil {
		log.Panicf("Could not retrieve OSVersion: %v", err)
	}
	fmt.Printf("OS Version: %s\n", osVersion.Version)
	fmt.Printf("Architecture: %s\n", sysinfo.Architecture().String())

	compilers, err := sysinfo.Compilers()
	if err != nil {
		log.Panicf("Could not retrieve Compilers: %v", err)
	}
	for _, compiler := range compilers {
		fmt.Printf("Compiler: %s\n", compiler.Name.String())
	}

	libc, err := sysinfo.Libc()
	if err != nil {
		log.Panicf("Could not retrieve Libc: %v", err)
	}
	fmt.Printf("Libc: %s\n", libc.Name.String())
}
