package Utils

import(
	"bufio"
	"os"
    "log"
)

func SortUniqueStringSlice(slice []string) []string {
    keys := make(map[string]bool)
    list := []string{} 
    for _, entry := range slice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }    
    return list
}

// difference returns the elements in `a` that aren't in `b`.
func Difference(ref, dif []string) []string {
    diffMap := map[string]bool{}
    for _, x := range dif {
        diffMap[x] = true
    }
    var newDiff []string
    for _, x := range ref {
        if _, ok := diffMap[x]; !ok {
            newDiff = append(newDiff, x)
        }
    }
    return newDiff
}

func ReadCSV(path string, header bool) (items []string)  {

    var returnItems []string
    csvFile, _ := os.Open(path)
	defer csvFile.Close()
	scanner := bufio.NewScanner(csvFile)
	scanner.Split(bufio.ScanLines)
	for{
        // Read to next token...In this case line
        if header == true{
            scanner.Scan() 
            header = false
            continue
        }
		line := scanner.Scan()
		if line == false {
			// returns false on error or EOF check err
			err := scanner.Err()
			if err == nil {
				break
			} else {
				log.Fatal(err)
			}
		}
		
		// Get data from scan with Bytes() or Text()
		returnItems = append(returnItems, scanner.Text())
    }
    return returnItems
}

func SliceToIntMap(elements []string) map[string]int {
    elementMap := make(map[string]int)
    for _, s := range elements {
        elementMap[s]++
    }
    return elementMap
}