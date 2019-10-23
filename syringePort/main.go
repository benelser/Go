package main

import (
	"fmt"
	"syscall"
	"unsafe"
	"golang.org/x/sys/windows"
	"time"
)

const (

	REMOTE_ASSEMBLY_STUB_LENGTH_RELEASE = 32
	PROCESS_CREATE_THREAD = 0x0002
	PROCESS_QUERY_INFORMATION = 0x0400
	PROCESS_VM_OPERATION = 0x0008
	PROCESS_VM_WRITE = 0x0020
	PROCESS_VM_READ = 0x0010

	MEM_RESERVE = 0x00002000
	MEM_COMMIT = 0x00001000
	PAGE_READWRITE = 0x04
	PAGE_EXECUTE_READWRITE = 0x40

	
	CREATE_THREAD_ACCESS = (PROCESS_CREATE_THREAD | PROCESS_QUERY_INFORMATION | PROCESS_VM_OPERATION | PROCESS_VM_WRITE | PROCESS_VM_READ)
	PROCESS_ALL_ACCESS	=	(windows.SYNCHRONIZE | windows.STANDARD_RIGHTS_ALL | 0xfff | CREATE_THREAD_ACCESS)
)

func main()  {
	
	// Import DLL
	var kernel32 = syscall.NewLazyDLL("Kernel32.dll")
	var MSCorWks = syscall.NewLazyDLL("MSCorWks.dll")
	MSCorWks.Load()

	// Get procedures
	var OpenProcess = kernel32.NewProc("OpenProcess")
	var VirtualAllocEx = kernel32.NewProc("VirtualAllocEx")
	var WriteProcessMemory = kernel32.NewProc("WriteProcessMemory")
	var ReadProcessMemory = kernel32.NewProc("ReadProcessMemory")
	var CreateRemoteThread = kernel32.NewProc("CreateRemoteThread")
    var LoadLibraryA = kernel32.NewProc("LoadLibraryA")


	// LoadLibraryA local test
	// loadedDLLHandle, _, _ := syscall.Syscall(LoadLibraryA.Addr(), 
	// 								1, // number of args 
	// 								uintptr(unsafe.Pointer(syscall.StringBytePtr("C:\\Users\\bjelser-admin\\Downloads\\hello-world-x64.dll"))),
	// 								0,
	// 								0,
	// 							)
	// fmt.Printf("\nDLL load result: %v", loadedDLLHandle)

	// time.Sleep(time.Second * 240)
	
	hwnd, _, _ := syscall.Syscall(OpenProcess.Addr(), 
									3, // number of args 
									uintptr(PROCESS_ALL_ACCESS),
									0,
									uintptr(8700),
								)	

	// VirtualAllocEx
	const MyDll = "C:\\Users\\bjelser-admin\\Downloads\\hello-world-x64.dll"
	const (
		dataSize uint32 = 260
	) 
	//const dataSizePtr = uintptr(dataSize)
	
	remoteBufferBaseAddress, _, _ := syscall.Syscall6(VirtualAllocEx.Addr(), 
								5, // number of args 
								uintptr(hwnd),
								0,
								uintptr(dataSize),
								uintptr((MEM_RESERVE | MEM_COMMIT)), 
								uintptr(PAGE_EXECUTE_READWRITE),
								0,
							)
	
	fmt.Printf("Handle: %v\r\n", hwnd)
	fmt.Printf("Remote Buffer base address: %v", remoteBufferBaseAddress)

	
	// Write memory
	WriteMemResult, _, _ := syscall.Syscall6(WriteProcessMemory.Addr(), 
								5, // number of args 
								uintptr(hwnd),
								uintptr(remoteBufferBaseAddress),
								uintptr(unsafe.Pointer(syscall.StringBytePtr(MyDll))),
								uintptr(dataSize), 
								0,
								0,
	)
	
	fmt.Printf("\nWrite memory result: %v\r\n", WriteMemResult)

	// read memory
	var (
		data [dataSize]byte
		length uint32
	)
	
	ReadMemResult, _, _ := syscall.Syscall6(ReadProcessMemory.Addr(), 
								5, // number of args 
								uintptr(hwnd),
								uintptr(remoteBufferBaseAddress),
    							uintptr(unsafe.Pointer(&data)),
								uintptr(dataSize),
								uintptr(unsafe.Pointer(&length)),
								0,
	)


	fmt.Printf("\nRead memory result: %v\r\n%s\r\n", ReadMemResult, data)

	// Create remote thread	

	remoteThreadHandle, _, _ := syscall.Syscall9(CreateRemoteThread.Addr(), 
								7, // number of args 
								uintptr(hwnd),
								0,
								0,
								uintptr(LoadLibraryA.Addr()),	
								uintptr(remoteBufferBaseAddress),
								0,
								0, 
								0,	// satisfy syscall9
								0,	// satisfy syscall9
	)

	fmt.Printf("\nRemote Thread Handle: %v\r\n", remoteThreadHandle)
	fmt.Printf("\n%v", windows.GetLastError())
	windows.CloseHandle(windows.Handle(hwnd))
	time.Sleep(time.Second * 240)



	// this appears to work as well but to minimize dependecies wrap own function def
	// windows.CloseHandle(windows.Handle(hwnd))
	// // Close Handle
	// success, _, _ := syscall.Syscall(CloseHandle.Addr(), 
	// 							1, // number of args 
	// 							hwnd,
	// 							0,
	// 							0,
	// 						)	

	// fmt.Printf("Close Handle Result%v", success)
}
