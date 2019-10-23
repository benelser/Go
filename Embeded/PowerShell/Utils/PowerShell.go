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

func InvokePowerShell(encodedCommand string)  (output []byte){
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

type Process struct {
	Handles int `json:"Handles"`
	Id int `json:"Id"`
	ProcessName string `json:"ProcessName"`
}