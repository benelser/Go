package PowerShell

import(
	"encoding/base64"
	"bytes"
	"log"
	"io"
	"os/exec"
	"io/ioutil"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"../Scripts"
	"encoding/json"
	"fmt"
)

func CreateEncodedPowerShellCommand(embeddedResource io.ReadCloser) (encodedCommand string) {

	// Create New Buffer for Posh script
	scriptBuf := new(bytes.Buffer)
	scriptBuf.ReadFrom(embeddedResource)
	
	// Process to Encode bytes to utf-16 to satisfy -encoded command switch with PowerShell
	win16be := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)

	// Create transformer
	utf16bom := unicode.BOMOverride(win16be.NewEncoder())

	// Create new reader using transformer 
	unicodeReader := transform.NewReader(bytes.NewReader(scriptBuf.Bytes()), utf16bom)

	// Read the entire UnicodeReader
	script, _ := ioutil.ReadAll(unicodeReader)

	// Convert To Base64
	return base64.StdEncoding.EncodeToString(script)
	
}

func InvokePowerShellEncodedCommand(encodedCommand string)  (output []byte){
	// Spawn Child Process
	cmd := exec.Command("powershell", "-encodedCommand", encodedCommand, "-NonInteractive", "-NoProfile")
	stdin, err := cmd.StdinPipe()

	// Ensures we close up when we are done
	defer stdin.Close()

	if err != nil{
		log.Fatal(err)
	}

	// Get output
	out, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	return out
}

// func WriteOutPowerShell(output []byte, rn int) {
// 	fmt.Printf("OUTPUT FROM #%v GO ROUTINE..\r\n", rn)
// 	// var processes []PowerShell.Process
// 	// Define type here to marshall
// 	// json.Unmarshal(output, &processes)
// 	// fmt.Printf("%v\r\n",processes)
// 	// fmt.Printf("%v\r\n",len(processes))
// 	// for _, process := range processes {
// 	// 	fmt.Printf("Process name: %v\r\n", process.ProcessName)
// 	// }
// }

func GetADComputers() (computers ADComputers) {
	
	// Return var
	adComputers := ADComputers{}

	// Read static Resource
	scriptReader, _ := PowerShellScripts.Open("GetADComputers.ps1")

	// Create encoded command
	encoddedCommand := CreateEncodedPowerShellCommand(scriptReader)

	// Invoke and marshal
	adComputersBytes := InvokePowerShellEncodedCommand(encoddedCommand)
	err := json.Unmarshal(adComputersBytes, &adComputers)
	if err != nil{
		errMessage := fmt.Sprintf("Failed to decode JSON Version with: %s\n", err)
		// myErr := errors.New(errMessage)
		// return myErr
		fmt.Printf("%v", errMessage)
	}

	return adComputers
}

type ADComputers struct {
	Computers []string
	Count int 
}