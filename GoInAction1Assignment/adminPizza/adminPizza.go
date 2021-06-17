package adminPizza

import (
	"errors"
	"fmt"
)

type Pizza struct {
	PizzaNo    int
	PizzaName  string
	PizzaPrice float64
}

type Node struct {
	Item Pizza
	Next *Node
}

// Create a linkedlist for the pizza menu
type Linkedlist struct {
	Head *Node
	Size int
}

func (p *Linkedlist) CreateStartMenu(standardPizza []string, standardPrice float64) error {

	for pizzaNo, pizzaName := range standardPizza {
		pizzaNo = pizzaNo + 1
		p.AddPizza(pizzaNo, pizzaName, standardPrice)
	}

	return nil
}

func (p *Linkedlist) AddPizza(pizzaNo int, pizzaName string, pizzaPrice float64) error {

	newPizza := Pizza{
		PizzaNo:    pizzaNo,
		PizzaName:  pizzaName,
		PizzaPrice: pizzaPrice,
	}

	newNode := &Node{
		Item: newPizza,
		Next: nil,
	}

	if p.Head == nil {
		p.Head = newNode
	} else {
		currentNode := p.Head
		for currentNode.Next != nil {
			currentNode = currentNode.Next
		}
		currentNode.Next = newNode
	}
	p.Size++

	return nil
}

func (p *Linkedlist) EditPizza(pizzaNo int, pizzaName string, pizzaPrice float64) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">>", err)
		}
	}()

	currentNode := p.Head

	if currentNode != nil {
		if pizzaNo == currentNode.Item.PizzaNo {
			currentNode.Item.PizzaName = pizzaName
			currentNode.Item.PizzaPrice = pizzaPrice
		} else {
			for currentNode.Next != nil {
				currentNode = currentNode.Next

				if pizzaNo == currentNode.Item.PizzaNo {
					currentNode.Item.PizzaName = pizzaName
					currentNode.Item.PizzaPrice = pizzaPrice
				}
			}
		}
	} else {
		panic("No pizza found")
	}
	return errors.New(">> Invalid pizza no")
}

func (p *Linkedlist) DeletePizza(pizzaNo int) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">>", err)
		}
	}()

	currentNode := p.Head

	if currentNode != nil {

		for i := 0; i < p.Size; i++ {
			if currentNode.Item.PizzaNo == pizzaNo {
				if i > 0 {
					prevNode := p.GetAt(i - 1)
					prevNode.Next = p.GetAt(i).Next
				} else {
					p.Head = currentNode.Next
				}
				p.Size--
				return nil
			}
			currentNode = currentNode.Next
		}
	} else {
		panic("No pizza found")
	}
	return errors.New(">> Invalid pizza no")
}

func (p *Linkedlist) GetAt(pos int) *Node {
	currentNode := p.Head
	if pos < 0 {
		return currentNode
	}
	if pos > (p.Size - 1) {
		return nil
	}
	for i := 0; i < pos; i++ {
		currentNode = currentNode.Next
	}
	return currentNode
}

func (p *Linkedlist) PrintPizzaMenu() error {
	currentNode := p.Head

	if currentNode == nil {
		return errors.New(">> Sorry. No pizza on the menu today")
	}

	fmt.Printf("(%d) - %+v - $%.2f\n", currentNode.Item.PizzaNo, currentNode.Item.PizzaName, currentNode.Item.PizzaPrice)
	for currentNode.Next != nil {
		currentNode = currentNode.Next
		fmt.Printf("(%d) - %+v - $%.2f\n", currentNode.Item.PizzaNo, currentNode.Item.PizzaName, currentNode.Item.PizzaPrice)
	}

	return nil
}

func (p *Linkedlist) SearchPizza(pizzaNo int) (Pizza, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">>", err)
		}
	}()

	currentNode := p.Head

	if currentNode != nil {
		if pizzaNo == currentNode.Item.PizzaNo {
			return currentNode.Item, nil
		} else {
			for currentNode.Next != nil {
				currentNode = currentNode.Next

				if pizzaNo == currentNode.Item.PizzaNo {
					return currentNode.Item, nil
				}
			}
		}
	} else {
		panic("No pizza found")
	}
	return currentNode.Item, errors.New(">> Invalid pizza no")
}

func (p *Linkedlist) GetAllPizza() ([]Pizza, error) {
	pizzaSlice := make([]Pizza, 0)

	currentNode := p.Head

	if currentNode == nil {
		return pizzaSlice, errors.New(">> Sorry. No pizza on the menu today")
	}

	pizzaSlice = append(pizzaSlice, currentNode.Item)
	for currentNode.Next != nil {
		currentNode = currentNode.Next
		pizzaSlice = append(pizzaSlice, currentNode.Item)
	}

	return pizzaSlice, nil
}
