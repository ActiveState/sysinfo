package sysinfo

// OsInfo represents an OS returned by OS().
type OsInfo int

const (
	// Linux represents the Linux operating system.
	Linux OsInfo = iota
	// Windows represents the Windows operating system.
	Windows
	// Mac represents the Macintosh operating system.
	Mac
	// UnknownOs represents an unknown operating system.
	UnknownOs
)

// OSVersionInfo represents an OS version returned by OSVersion().
type OSVersionInfo struct {
	Version string // raw version string
	Major   int    // major version number
	Minor   int    // minor version number
	Micro   int    // micro version number
	Name    string // free-form name string (varies by OS)
}

// ArchInfo represents an architecture returned by Architecture().
type ArchInfo int

const (
	// I386 represents the Intel x86 (32-bit) architecture.
	I386 ArchInfo = iota
	// Amd64 represents the x86_64 (64-bit) architecture.
	Amd64
	// Arm represents the ARM architecture.
	Arm
	// UnknownArch represents an unknown architecture.
	UnknownArch
)

// LibcNameInfo represents a C library name.
type LibcNameInfo int

const (
	// Glibc represents the GNU C library.
	Glibc LibcNameInfo = iota
	// Msvcrt represents the Microsoft Visual C++ runtime library.
	Msvcrt
	// UnknownLibc represents an unknown C library.
	UnknownLibc
)

// LibcInfo represents a LibC returned by Libc().
type LibcInfo struct {
	Name  LibcNameInfo // C library name
	Major int          // major version number
	Minor int          // minor version number
}

// CompilerNameInfo reprents a compiler toolchain name.
type CompilerNameInfo int

const (
	// Gcc represents the GNU C Compiler toolchain.
	Gcc CompilerNameInfo = iota
	// Clang represents the LLVM/Clang toolchain.
	Clang
	// Msvc represents the Microsoft Visual C++ toolchain.
	Msvc
	// Mingw represents the Minimalist GNU for Windows toolchain.
	Mingw
	// Cygwin represents the Cygwin toolchain.
	Cygwin
)

// CompilerInfo represents a compiler toolchain returned by Compiler().
type CompilerInfo struct {
	Name  CompilerNameInfo // C compiler name
	Major int              // major version number
	Minor int              // minor version number
}
