package main

import(
	"fmt"
)

func main()  {
	me := Person{
		FName : "Me",
		LName : "You",
		Age : 31,
		Married: true,
	}
	fmt.Println(me)
	me.UpdateEmail("Me.You@gmail.com")
	fmt.Println(me)
	me.Introduce()

	her := NewPerson("Her", "Name")
	her.Introduce()
	her.UpdateEmail("Her.Name@gmail.com")
	fmt.Println(*her)
	her.UpdateEmail("Test@gmail.com")
	fmt.Println(*her)
}

type Person struct {
	FName, LName string
	Age int
	Married bool
	Email string
}

// This method must take pointer value for its receiver because it is going to modify the object
func (p *Person)UpdateEmail(email string)  {
	
	p.Email = email
}

// Wires up method to Person struct
// In go p is called the methods receiver "sending a message to an object"
// In go there is no self or this for the receiver 
func (p Person)Introduce(){
	
	fmt.Printf("Hi! My name is %s %s!\r\n", p.FName, p.LName)
}

// New Person factory like func
func NewPerson(f, l string)  *Person {
	
	return &Person{
		FName : f,
		LName : l,
	}
}
