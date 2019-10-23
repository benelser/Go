package main

import(
	"fmt"
	"encoding/json"
	"./PowerShell/Utils"
	"./PowerShell/Scripts"
)

// Build "files.go" using the following command: staticfiles -o PowerShell/PowerShellScripts.go PowerShell/Scripts

func main()  {
	
	// Read static Resource
	scriptReader, _ := PowerShellScripts.Open("ConvertToJson.ps1")

	// Create encoded command
	encoddedCommand := PowerShell.CreateEncodedPowerShellCommand(scriptReader)

	// Stage channels to read data back from threads
	scriptOutput1 := make(chan []byte)
	scriptOutput2 := make(chan []byte)
	scriptOutput3 := make(chan []byte)
	scriptOutput4 := make(chan []byte)
	scriptOutput5 := make(chan []byte)
	scriptOutput6 := make(chan []byte)
	scriptOutput7 := make(chan []byte)
	scriptOutput8 := make(chan []byte)
	scriptOutput9 := make(chan []byte)
	scriptOutput10 := make(chan []byte)
	scriptOutput11 := make(chan []byte)
	scriptOutput12 := make(chan []byte)
	scriptOutput13 := make(chan []byte)
	scriptOutput14 := make(chan []byte)
	scriptOutput15 := make(chan []byte)
	scriptOutput16 := make(chan []byte)
	scriptOutput17 := make(chan []byte)
	scriptOutput18 := make(chan []byte)
	scriptOutput19 := make(chan []byte)
	scriptOutput20 := make(chan []byte)
	scriptOutput21 := make(chan []byte)
	scriptOutput22 := make(chan []byte)
	scriptOutput23 := make(chan []byte)
	scriptOutput24 := make(chan []byte)
	scriptOutput25 := make(chan []byte)
	scriptOutput26 := make(chan []byte)
	scriptOutput27 := make(chan []byte)
	scriptOutput28 := make(chan []byte)
	scriptOutput29 := make(chan []byte)
	scriptOutput30 := make(chan []byte)
	scriptOutput31 := make(chan []byte)
	scriptOutput32 := make(chan []byte)

	// Fire off Go Routines 
	go func() { scriptOutput1 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput2 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput3 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput4 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput5 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput6 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput7 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput8 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput9 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput10 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput11 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput12 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput13 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput14 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput15 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput16 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput17 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput18 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput19 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput20 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput21 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput22 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput23 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput24 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput25 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput26 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput27 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput28 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput29 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput30 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput31 <- PowerShell.InvokePowerShell(encoddedCommand)}()
	go func() { scriptOutput32 <- PowerShell.InvokePowerShell(encoddedCommand)}()

	// Demo main thread still executing
	fmt.Println("Sitting here on the main thread")
	fmt.Println("Sitting here on the main thread")
	fmt.Println("Sitting here on the main thread")
	fmt.Println("doing more stuff")

	// Print Results
	writeOutPowerShell(<- scriptOutput1, 1)
	writeOutPowerShell(<- scriptOutput32, 32)

}

func writeOutPowerShell(output []byte, rn int) {
	fmt.Printf("OUTPUT FROM #%v GO ROUTINE..\r\n", rn)
	var processes []PowerShell.Process
	json.Unmarshal(output, &processes)
	fmt.Printf("%v\r\n",processes)
	fmt.Printf("%v\r\n",len(processes))
	for _, process := range processes {
		fmt.Printf("Process name: %v\r\n", process.ProcessName)
	}
}