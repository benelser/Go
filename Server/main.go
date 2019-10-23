package main

import(
	. "./Server"

)

// Build "files.go" using the following command: staticfiles -o PowerShell/PowerShellScripts.go PowerShell/Scripts
// Move to Scripts dir and rename package to PowerShellScripts
func main()  {
	StartServer()
}