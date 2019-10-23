package main

import (
	"fmt"
	"syscall"
	"unsafe"
	. "golang.org/x/sys/windows"
	"time"
)

// Api Docs
// https://docs.microsoft.com/en-us/windows/desktop/api/winuser/nf-winuser-messageboxw

func main() {

	// Import DLL
	var user32 = syscall.NewLazyDLL("user32.dll")
	var kernel32 = syscall.NewLazyDLL("Kernel32.dll")

	// // call on wrapper 
	CreateDirectoryA(kernel32, "C:\\Temp\\BensGo\\Load_DLL\\MyTest")
	MessageBoxW(user32, "This Can Get Dirty", "Go and Win32",)
	CreateProcessW(kernel32, `C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`)
	time.Sleep(time.Second * 120)

}
	

func MessageBoxW(DLL *syscall.LazyDLL, message, title string) {

	// Define constants for MessageBox and any others to satisfy win32 api
	const(
		MB_YESNOCANCEL = 0x00000003
		MB_ICONEXCLAMATION = 0x00000030
	)
	// Get pointer to function/procedure
	var MessageBoxW = DLL.NewProc("MessageBoxW")
	returnValue, _, _ := MessageBoxW.Call(0,
		// StringToUTF16Ptr == LPCWSTR 
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(message))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		// unitptr == UINT which are defined as constants
		uintptr(MB_YESNOCANCEL | MB_ICONEXCLAMATION))
	fmt.Printf("Return: %d\n", returnValue)
}

func CreateDirectoryA(DLL *syscall.LazyDLL, pathName string)  {
	// Get pointer to function/procedure
	var CreateDirectory = DLL.NewProc("CreateDirectoryA")
	// StringBytePtr == LPCSTR
	CreateDirectory.Call(uintptr(unsafe.Pointer(syscall.StringBytePtr(pathName))),0)
}

func CreateProcessW(DLL *syscall.LazyDLL, pathToExe string){

	// Imported types from windows package
	// All are null empty structs satisfying the api
	procSecurity := &SecurityAttributes{}
	threadSecurity := &SecurityAttributes{}
	startupInfo := &StartupInfo{}
	outProcInfo := &ProcessInformation{}
	proc := DLL.NewProc("CreateProcessW")
	// Addr returns the address of the procedure represented by p.
	// The return value can be passed to Syscall to run the procedure.
	// nargs is 10 representing number of args
	var CREATE_NO_WINDOW = 0x08000000
	r1, _, _ := syscall.Syscall12(proc.Addr(), 
									10, // number of args 
									uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(pathToExe))), 
									0, 
									uintptr(unsafe.Pointer(procSecurity)), 
									uintptr(unsafe.Pointer(threadSecurity)), 
									0,
									uintptr(CREATE_NO_WINDOW), 
									0, 
									0, 
									uintptr(unsafe.Pointer(startupInfo)), 
									uintptr(unsafe.Pointer(outProcInfo)),
									0, // indicates null for 11th arg of Syscall12
									0, // indicates null for 12th arg of Syscall12
								)
	fmt.Printf("%v\r\n", r1)
	fmt.Printf("PowerShell process ID: %v\r\n", outProcInfo.ProcessId)	
}