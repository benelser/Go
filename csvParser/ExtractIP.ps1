$regexString = "((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)"
$file = "C:\Temp\BensGo\csvParser\IPTest.csv"
$Regex = [regex]::new($regexString)

[System.Collections.Hashtable]$myIpAddresses = @{}

$streamReader = [System.IO.StreamReader]::new($file)
$valueCounter = 1
$sw = [System.Diagnostics.Stopwatch]::StartNew()
while ($streamReader.EndOfStream -ne $true) {
    
  $lineContent = $streamReader.ReadLine()
  $Regex = [regex]::new($regexString)
  $matchesArray = $Regex.Matches($lineContent)
  foreach ($match in $matchesArray) {
      if ($myIpAddresses.Contains($match.Value)) {
          continue
      }
      $myIpAddresses.Add($match.Value, $valueCounter)
      $valueCounter ++
  }
    
}

$streamReader.Close()
$streamReader.Dispose()
$sw.Elapsed
$myIpAddresses