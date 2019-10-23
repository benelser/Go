package main

import "fmt"

func main()  {

	fruits := []Fruit{Fruit{Name:"Apple", Amount: 33}, Fruit{Name:"Orange", Amount: 3}}
	me := NewPerson(31, fruits, "Benjamin", "Elser")
	fmt.Println(*me)
	// Get the type usig the %T Verb
	fmt.Printf("%T\r\n", *me)

	// Calling method on fruits inside of Person struct
	for _, fruit := range me.FavoriteFruits {
		fruit.Eat()
	}
}

type Person struct {
	FName string
	LName string
	Age int
	FavoriteFruits []Fruit
}

type Fruit struct {
	Name string
	Amount int
}

// Wire up receiver 
func (f Fruit) Eat()  {
	fmt.Printf("%ss are so yummy!\r\n", f.Name)
}

func NewPerson(age int, fruits []Fruit, firstName, lastName string) *Person {
	return &Person{
		FName : firstName,
		LName : lastName,
		Age : age,
		FavoriteFruits : fruits,
	}
}