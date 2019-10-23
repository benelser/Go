package main

// Import entire package..... This method may cause conflict with a larger project
// Standard way is to omit . 
// You can also specify an alias instead of . 
import("fmt")

func greet(f, l string){
	fmt.Printf("Hello there %s %s\r\n", f, l)
}

// You can omit a name of the return var here and just define a comma seperated list of types 
// For example: (int, string, err)
func greetString(f, l string) (greeting string) {
	return fmt.Sprintf("Hello there %s %s", f, l)
}

// The type of a function is sometimes called its signature
// If two functions share the same paramaters and return types they are said to have the same signature 
// Go has no concept of default paramaters 

// Example of multiple return values
func greetingAndAge(age int, f, l string) (fullname string, yourAge int){
	return fmt.Sprintf("%s %s", f, l), age
}

// Functions are first class values in Go

func SayHello(name string) string  {
	return fmt.Sprintf("Hello there %s", name)
}

// Anonymous Funcs
// This function returns a anonymous func that return a int
func counter() func() int {
	x := 0
	return func() int {
		x++
		return x
	}
}

// A function is a wrapper around a sequence of statements to be reused throughout a project
func main(){
	fmt.Printf("This is so much better %s\r\n", "Ben")
	greet("Benjamin", "Elser")
	greeting := greetString("Benjamin", "Elser")
	fmt.Printf(greeting)

	// Th result of calling a multi-valued functions is a tuple of values.
	// Values must be explicitly assigned or ignored using the "_" identifier
	greeting, age := greetingAndAge(31, "Benjamin", "Elser")
	fmt.Printf("\n%s. You are %d years old!\r\n", greeting, age)

	// Example of setting func to var and calling 
	myGreeter := SayHello
	fmt.Printf("%s", myGreeter("Benjamin"))

	// calling the function that implements anonymous func
	c := counter()
	fmt.Printf("\nCounter Value : %d\r\n", c())
	fmt.Printf("Counter Value : %d\r\n", c())

	// Calling Variadic function
	fmt.Printf(MyStringJoin("-", "Benjamin", "James", "Elser"))
}

func MyStringJoin(delimiter string, str ...string) string {

	myString := ""
	for index := 0; index < len(str); index++ {
		if index == (len(str) - 1) {
			myString += str[index]
			myString += "\r\n"
			break
			
		}
		myString += str[index]
		myString += delimiter
	}

	return myString
}