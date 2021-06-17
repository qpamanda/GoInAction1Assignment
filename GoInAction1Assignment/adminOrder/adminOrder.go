package adminOrder

import (
	aPizza "GoAdvancedAssignment/adminPizza"
	"errors"
	"fmt"
)

type OrderItem struct {
	PizzaNo  int
	OrderQty int
}

type Order struct {
	OrderNo    int
	OrderSlice []OrderItem
	TotalCost  float64
	UserName   string
}

type Node struct {
	Item Order
	Next *Node
}

type Queue struct {
	Front *Node
	Back  *Node
	Size  int
}

func (p *Queue) Enqueue(orderNo int, orderSlice []OrderItem, totalCost float64, userName string) error {

	newOrder := Order{
		OrderNo:    orderNo,
		OrderSlice: orderSlice,
		TotalCost:  totalCost,
		UserName:   userName,
	}

	newNode := &Node{
		Item: newOrder,
		Next: nil,
	}
	if p.Front == nil {
		p.Front = newNode

	} else {
		p.Back.Next = newNode

	}
	p.Back = newNode
	p.Size++

	return nil
}

func (p *Queue) Dequeue() (Order, error) {

	var item Order

	if p.Front == nil {
		return item, errors.New(">> No orders to delete")
	}

	item = p.Front.Item

	if p.Size == 1 {
		p.Front = nil
		p.Back = nil
	} else {
		p.Front = p.Front.Next
	}
	p.Size--

	return item, nil
}

func (p *Queue) PrintAllOrders(pizzaList *aPizza.Linkedlist) error {
	currentNode := p.Front
	if currentNode == nil {
		fmt.Println(">> No orders in the queue")
		return nil
	}

	p.PrintOrder(currentNode.Item, pizzaList)

	for currentNode.Next != nil {
		currentNode = currentNode.Next
		p.PrintOrder(currentNode.Item, pizzaList)
	}
	return nil
}

func (p *Queue) PrintOrder(item Order, pizzaList *aPizza.Linkedlist) error {

	fmt.Println("Order No: ", item.OrderNo)
	orderSlice := item.OrderSlice

	if len(orderSlice) > 0 {
		for _, v := range orderSlice {
			pizzaOrder, _ := pizzaList.SearchPizza(v.PizzaNo)
			pizzaTotal := float64(v.OrderQty) * pizzaOrder.PizzaPrice

			fmt.Printf("%d x %s @ $%.2f\n", v.OrderQty, pizzaOrder.PizzaName, pizzaTotal)
		}
		fmt.Printf("Total Amount = $%.2f\n", item.TotalCost)
		fmt.Println()
		fmt.Println()
	}

	return nil
}

func (p *Queue) SearchOrder(orderNo int) (Order, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">>", err)
		}
	}()

	currentNode := p.Front

	if currentNode != nil {
		if currentNode.Item.OrderNo == orderNo {
			return currentNode.Item, nil
		} else {
			for currentNode.Next != nil {
				currentNode = currentNode.Next
				if currentNode.Item.OrderNo == orderNo {
					return currentNode.Item, nil
				}
			}
		}
	} else {
		panic("No orders found")
	}

	return currentNode.Item, errors.New(">> No orders found")
}

func (p *Queue) SearchPizzaInOrder(pizzaNo int) (bool, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">>", err)
		}
	}()

	currentNode := p.Front

	if currentNode != nil {
		orderSlice := currentNode.Item.OrderSlice

		// check if pizza found in slice
		if SearchPizzaInSlice(pizzaNo, orderSlice) {
			return true, nil
		} else {
			for currentNode.Next != nil {
				currentNode = currentNode.Next
				orderSlice = currentNode.Item.OrderSlice

				// check if pizza found in slice
				if SearchPizzaInSlice(pizzaNo, orderSlice) {
					return true, nil
				}
			}
		}

	} else {
		panic("No orders found")
	}

	return false, errors.New(">> No orders found")
}

func SearchPizzaInSlice(pizzaNo int, orderSlice []OrderItem) bool {
	if len(orderSlice) > 0 {
		for _, v := range orderSlice {
			if v.PizzaNo == pizzaNo {
				return true // if pizza not found in any orders, return true
			}
		}
	}
	return false
}

func (p *Queue) UpdateOrder(orderNo int, orderSlice []OrderItem, totalCost float64) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">>", err)
		}
	}()

	currentNode := p.Front

	if currentNode != nil {
		if currentNode.Item.OrderNo == orderNo {
			currentNode.Item.OrderSlice = orderSlice
			currentNode.Item.TotalCost = totalCost
			return nil
		} else {
			for currentNode.Next != nil {
				currentNode = currentNode.Next
				if currentNode.Item.OrderNo == orderNo {
					currentNode.Item.OrderSlice = orderSlice
					currentNode.Item.TotalCost = totalCost
					return nil
				}
			}
		}
	} else {
		panic("No orders found")
	}

	return errors.New(">> No orders found")
}

func (p *Queue) IsEmpty() bool {
	return p.Size == 0
}

func (p *Queue) GetAllOrders(userName string) ([]Order, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> panic:", err)
		}
	}()

	orderList := make([]Order, 0)

	currentNode := p.Front

	if currentNode != nil {
		if userName != "admin" {
			if currentNode.Item.UserName == userName {
				orderList = append(orderList, currentNode.Item)
			}
			for currentNode.Next != nil {
				currentNode = currentNode.Next
				if currentNode.Item.UserName == userName {
					orderList = append(orderList, currentNode.Item)
				}
			}

		} else {
			orderList = append(orderList, currentNode.Item)
			for currentNode.Next != nil {
				currentNode = currentNode.Next
				orderList = append(orderList, currentNode.Item)
			}
		}
	} else {
		panic("No orders found")
	}

	return orderList, nil
}
