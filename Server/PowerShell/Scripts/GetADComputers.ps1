# Get AD Computer Names
$computers = (Get-ADComputer -Filter *).Name | Sort-Object -Unique

# Return object
class ADComputers {
    $Computers
    $Count

    ADComputers($Computers) {
        $this.Computers = $Computers
        $this.Count = $this.Computers.Count
    }
}

# Instantiate return object
$ADComputers = [ADComputers]::new($computers)

# Convert to json so Go can Consume it
$ADComputers | ConvertTo-Json -Compress 