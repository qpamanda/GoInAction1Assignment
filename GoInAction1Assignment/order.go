package main

import (
	aOrder "GoInAction1Assignment/adminOrder"
	aPizza "GoInAction1Assignment/adminPizza"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
)

func generateOrderNo(orderSlice []aOrder.OrderItem) {

	defer wg.Done()

	// Increment orderNo global variable by 1 if there are OrderItem in the slice
	if len(orderSlice) > 0 {
		mutex.Lock()
		runtime.Gosched()
		newOrderNo = newOrderNo + 1
		mutex.Unlock()
	}
}

func addOrder(orderSlice []aOrder.OrderItem, pizzaNo int, orderQty int) []aOrder.OrderItem {
	orderItem := aOrder.OrderItem{
		PizzaNo:  pizzaNo,
		OrderQty: orderQty,
	}

	orderSlice = append(orderSlice, orderItem)

	return orderSlice
}

func addorder(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	clientMsg := ""
	bValidOrder := true

	viewOrderItemSlice := make([]viewOrderItem, 0)
	orderItemSlice := make([]aOrder.OrderItem, 0)

	pizzaSlice, err := pizzaList.GetAllPizza()

	if err != nil {
		clientMsg = "Error getting pizza menu. "
		fmt.Println(">> Error getting pizza menu. ")
	} else {
		for idx1, val1 := range pizzaSlice {
			pizzaOrder := viewOrderItem{idx1 + 1, val1.PizzaNo, val1.PizzaName, fmt.Sprintf("%.2f", val1.PizzaPrice), 0, "", ""}
			viewOrderItemSlice = append(viewOrderItemSlice, pizzaOrder)
		}
	}

	if req.Method == http.MethodPost {
		for _, val1 := range viewOrderItemSlice {
			errMsg := ""

			selectedPizza := req.FormValue(strconv.Itoa(val1.PizzaNo))

			if selectedPizza != "" {
				selectedQty := req.FormValue("orderqty" + strconv.Itoa(val1.ItemNo))

				pizzaNo, _ := strconv.Atoi(selectedPizza)
				orderQty, errQty := strconv.Atoi(selectedQty)

				if orderQty <= 0 || orderQty > maxOrderQty || errQty != nil {
					errMsg = "Enter a valid quantity"
					bValidOrder = false
				}

				viewOrderItemSlice[val1.ItemNo-1].OrderQty = orderQty
				viewOrderItemSlice[val1.ItemNo-1].Checked = "checked"
				viewOrderItemSlice[val1.ItemNo-1].ErrorMsg = errMsg

				orderItemSlice = addOrder(orderItemSlice, pizzaNo, orderQty)
			}
		}

		if len(orderItemSlice) > 0 && bValidOrder {
			wg.Add(1)
			go generateOrderNo(orderItemSlice)
			wg.Wait()

			orderNo := newOrderNo

			totalCost := getTotalCost(orderItemSlice)

			orderQueue.Enqueue(orderNo, orderItemSlice, totalCost, myUser.Username)

			printDividerLine()
			fmt.Println("* RECEIPT *")
			fmt.Println()

			printOrder(orderNo, orderItemSlice, totalCost)

			clientMsg = "Order " + strconv.Itoa(orderNo) + " added successfully. Total payment is $" + fmt.Sprintf("%.2f", totalCost)

		} else {
			clientMsg = clientMsg + "No orders made. "
		}
	}

	data := struct {
		User            user
		OrderSlice      []viewOrderItem
		CntCurrentItems int
		MaxOrder        int
		ClientMsg       string
	}{
		myUser,
		viewOrderItemSlice,
		len(viewOrderItemSlice),
		maxOrderQty,
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "addorder.gohtml", data)
}

func editorder(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	inputOrderNo := newOrderNo
	clientMsg := ""
	bValidOrder := true
	bValidOrderNo := false

	viewOrderItemSlice := make([]viewOrderItem, 0)
	newOrderItemSlice := make([]aOrder.OrderItem, 0)
	pizzaSlice, err := pizzaList.GetAllPizza()

	if err != nil {
		clientMsg = "Error getting pizza menu. "
		fmt.Println(">> Error getting pizza menu")
	} else {
		for i, v := range pizzaSlice {
			pizzaOrder := viewOrderItem{i + 1, v.PizzaNo, v.PizzaName, fmt.Sprintf("%.2f", v.PizzaPrice), 0, "", ""}
			viewOrderItemSlice = append(viewOrderItemSlice, pizzaOrder)
		}
	}

	if req.Method == http.MethodPost {
		inputOrderNo, _ = strconv.Atoi(req.FormValue("orderno"))

		myOrder, err := orderQueue.SearchOrder(inputOrderNo)
		myOrderSlice := myOrder.OrderSlice

		if err != nil || len(myOrderSlice) == 0 {
			clientMsg = "Order not found. "
			bValidOrderNo = false
		} else {
			if myOrder.UserName != myUser.Username && myUser.Username != "admin" {
				clientMsg = "Order not found. "
				bValidOrderNo = false
			} else {
				bValidOrderNo = true

				if bFirst {
					for _, val1 := range viewOrderItemSlice {
						for _, val2 := range myOrderSlice {
							if val1.PizzaNo == val2.PizzaNo {
								viewOrderItemSlice[val1.ItemNo-1].OrderQty = val2.OrderQty
								viewOrderItemSlice[val1.ItemNo-1].Checked = "checked"
								viewOrderItemSlice[val1.ItemNo-1].ErrorMsg = ""
							}
						}
					}
					bFirst = false
				}

				for _, val1 := range viewOrderItemSlice {
					errMsg := ""

					selectedPizza := req.FormValue(strconv.Itoa(val1.PizzaNo))

					if selectedPizza != "" {
						selectedQty := req.FormValue("orderqty" + strconv.Itoa(val1.ItemNo))

						pizzaNo, _ := strconv.Atoi(selectedPizza)
						orderQty, errQty := strconv.Atoi(selectedQty)

						if orderQty <= 0 || orderQty > maxOrderQty || errQty != nil {
							errMsg = "Enter a valid quantity"
							bValidOrder = false
						}

						viewOrderItemSlice[val1.ItemNo-1].OrderQty = orderQty
						viewOrderItemSlice[val1.ItemNo-1].Checked = "checked"
						viewOrderItemSlice[val1.ItemNo-1].ErrorMsg = errMsg

						newOrderItemSlice = addOrder(newOrderItemSlice, pizzaNo, orderQty)
					}
				}

				if len(newOrderItemSlice) > 0 && bValidOrder {
					totalCost := getTotalCost(newOrderItemSlice)

					orderQueue.UpdateOrder(inputOrderNo, newOrderItemSlice, totalCost)

					printDividerLine()
					fmt.Println("* RECEIPT *")
					fmt.Println()

					printOrder(inputOrderNo, newOrderItemSlice, totalCost)

					clientMsg = "Order [" + strconv.Itoa(inputOrderNo) + "] updated successfully. Total payment is $" + fmt.Sprintf("%.2f", totalCost)
				} else {
					clientMsg = clientMsg + "No orders updated. "
				}
			}
		}
	}

	data := struct {
		User            user
		OrderNo         int
		OrderSlice      []viewOrderItem
		CntCurrentItems int
		MaxOrder        int
		ClientMsg       string
		ValidOrderNo    bool
	}{
		myUser,
		inputOrderNo,
		viewOrderItemSlice,
		len(viewOrderItemSlice),
		maxOrderQty,
		clientMsg,
		bValidOrderNo,
	}

	tpl.ExecuteTemplate(res, "editorder.gohtml", data)
}

func vieworders(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	clientMsg := ""

	viewOrderSlice, err := getCurrentOrders(myUser.Username)

	if err != nil {
		clientMsg = "No orders found. "
		fmt.Println(">> Error getting current orders")
	}

	myCompletedOrderSlice := getCompletedOrders(myUser.Username)

	data := struct {
		User                user
		ViewOrderSlice      []viewOrder
		CntCurrentItems     int
		CompletedOrderSlice []viewOrder
		CntCompletedItems   int
		ClientMsg           string
	}{
		myUser,
		viewOrderSlice,
		len(viewOrderSlice),
		myCompletedOrderSlice,
		len(myCompletedOrderSlice),
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "vieworders.gohtml", data)
}

func completeorder(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	clientMsg := ""

	if req.Method == http.MethodPost {
		pizzaOrder, err := orderQueue.Dequeue()

		if err != nil {
			clientMsg = "No order to dequeue."
		} else {

			updateCompletedOrders(pizzaOrder, myUser)
			clientMsg = "Order [" + strconv.Itoa(pizzaOrder.OrderNo) + "] completed and added to pizza sales."
		}
	}

	viewOrderSlice, err := getCurrentOrders(myUser.Username)

	if err != nil {
		clientMsg = "No orders found. "
		//fmt.Println(">> Error getting current orders")
	}

	data := struct {
		User            user
		ViewOrderSlice  []viewOrder
		CntCurrentItems int
		ClientMsg       string
	}{
		myUser,
		viewOrderSlice,
		len(viewOrderSlice),
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "completeorder.gohtml", data)
}

func pizzasales(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	clientMsg := ""

	viewOrderSlice, err := getCurrentOrders(myUser.Username)

	if err != nil {
		clientMsg = "No orders found. "
		fmt.Println(">> Error getting current orders")
	}

	myCompletedOrderSlice := getCompletedOrders(myUser.Username)

	viewCurrentPizzaSales := getPizzaSales(viewOrderSlice)
	viewCompletedPizzaSales := getPizzaSales(myCompletedOrderSlice)

	currentPizzaSalesTotal := getTotalSales(viewCurrentPizzaSales)
	completedPizzaSalesTotal := getTotalSales(viewCompletedPizzaSales)

	data := struct {
		User                     user
		CurrentPizzaSales        []viewPizzaSales
		CntCurrentItems          int
		CurrentPizzaSalesTotal   string
		CompletedPizzaSales      []viewPizzaSales
		CntCompletedItems        int
		CompletedPizzaSalesTotal string
		ClientMsg                string
	}{
		myUser,
		viewCurrentPizzaSales,
		len(viewCurrentPizzaSales),
		fmt.Sprintf("%.2f", currentPizzaSalesTotal),
		viewCompletedPizzaSales,
		len(viewCompletedPizzaSales),
		fmt.Sprintf("%.2f", completedPizzaSalesTotal),
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "pizzasales.gohtml", data)
}

func printOrder(orderNo int, orderSlice []aOrder.OrderItem, totalCost float64) {

	fmt.Println("Order No: ", orderNo)
	fmt.Println()

	pizzaTotal := 0.0

	for _, v := range orderSlice {
		pizzaOrder, _ := pizzaList.SearchPizza(v.PizzaNo)
		pizzaTotal = float64(v.OrderQty) * pizzaOrder.PizzaPrice

		fmt.Printf("%d x %s\t$%.2f\n", v.OrderQty, pizzaOrder.PizzaName, pizzaTotal)
	}

	fmt.Println("\t\t\t--------")
	fmt.Printf("TOTAL PAYMENT\t\t$%.2f\n", totalCost)
	fmt.Println("\t\t\t--------")
}

func getTotalCost(orderSlice []aOrder.OrderItem) float64 {

	orderTotal := 0.0
	pizzaTotal := 0.0

	var pizzaOrder aPizza.Pizza

	for _, v := range orderSlice {
		pizzaOrder, _ = pizzaList.SearchPizza(v.PizzaNo)
		pizzaTotal = float64(v.OrderQty) * pizzaOrder.PizzaPrice
		orderTotal = orderTotal + pizzaTotal
	}

	return orderTotal
}

func updateCompletedOrders(completedOrder aOrder.Order, myUser user) {

	orderSlice := completedOrder.OrderSlice
	viewOrderItemSlice := make([]viewOrderItem, 0)
	pizzaSlice, _ := pizzaList.GetAllPizza()

	for idx1, val1 := range orderSlice {
		for _, val2 := range pizzaSlice {
			if val1.PizzaNo == val2.PizzaNo {
				pizzaOrder := viewOrderItem{idx1 + 1, val1.PizzaNo, val2.PizzaName, fmt.Sprintf("%.2f", val2.PizzaPrice), val1.OrderQty, "", ""}
				viewOrderItemSlice = append(viewOrderItemSlice, pizzaOrder)
			}
		}
	}

	viewOrder := viewOrder{len(completedOrderSlice) + 1, completedOrder.OrderNo, viewOrderItemSlice, fmt.Sprintf("%.2f", completedOrder.TotalCost), completedOrder.UserName}
	completedOrderSlice = append(completedOrderSlice, viewOrder)
}

func getCurrentOrders(userName string) ([]viewOrder, error) {

	viewOrderSlice := make([]viewOrder, 0)
	pizzaSlice, _ := pizzaList.GetAllPizza()
	orderQSlice, err := orderQueue.GetAllOrders(userName)

	if err != nil {
		return viewOrderSlice, errors.New(">> Error getting orders")
	} else {
		for idx1, val1 := range orderQSlice {
			orderSlice := val1.OrderSlice
			viewOrderItemSlice := make([]viewOrderItem, 0)

			for idx2, val2 := range orderSlice {
				for _, val3 := range pizzaSlice {
					if val2.PizzaNo == val3.PizzaNo {
						pizzaOrder := viewOrderItem{idx2 + 1, val2.PizzaNo, val3.PizzaName, fmt.Sprintf("%.2f", val3.PizzaPrice), val2.OrderQty, "", ""}
						viewOrderItemSlice = append(viewOrderItemSlice, pizzaOrder)
					}
				}
			}

			viewOrder := viewOrder{idx1 + 1, val1.OrderNo, viewOrderItemSlice, fmt.Sprintf("%.2f", val1.TotalCost), val1.UserName}
			viewOrderSlice = append(viewOrderSlice, viewOrder)
		}
	}

	return viewOrderSlice, nil
}

func getCompletedOrders(userName string) []viewOrder {

	myCompletedOrderSlice := make([]viewOrder, 0)

	if userName != "admin" {
		i := 0
		for _, val1 := range completedOrderSlice {
			if val1.UserName == userName {
				myCompletedOrderSlice = append(myCompletedOrderSlice, val1)
				myCompletedOrderSlice[i].IdxNo = i + 1
				i++
			}
		}
	} else {
		return completedOrderSlice
	}

	return myCompletedOrderSlice
}

func getPizzaSales(viewOrderSlice []viewOrder) []viewPizzaSales {

	viewPizzaSalesSlice := make([]viewPizzaSales, 0)

	for _, val1 := range viewOrderSlice {
		viewOrderItemSlice := val1.ViewOrderItems
		for _, val2 := range viewOrderItemSlice {
			viewPizzaSalesSlice = updatePizzaInSlice(val2, viewPizzaSalesSlice)
		}
	}

	return viewPizzaSalesSlice
}

func updatePizzaInSlice(vOrderItem viewOrderItem, viewPizzaSalesSlice []viewPizzaSales) []viewPizzaSales {
	bUpdate := false
	pizzaPrice, _ := strconv.ParseFloat(vOrderItem.PizzaPrice, 64)
	totalSales := float64(vOrderItem.OrderQty) * pizzaPrice

	if len(viewPizzaSalesSlice) > 0 {
		for i, v := range viewPizzaSalesSlice {
			if v.PizzaNo == vOrderItem.PizzaNo {
				viewPizzaSalesSlice[i].OrderQty = viewPizzaSalesSlice[i].OrderQty + vOrderItem.OrderQty
				viewPizzaSalesSlice[i].TotalSales = viewPizzaSalesSlice[i].TotalSales + totalSales
				viewPizzaSalesSlice[i].STotalSales = fmt.Sprintf("%.2f", viewPizzaSalesSlice[i].TotalSales)
				bUpdate = true
			}
		}
	}

	if !bUpdate {
		viewPizzaSales := viewPizzaSales{vOrderItem.PizzaNo, vOrderItem.PizzaName, vOrderItem.OrderQty, totalSales, fmt.Sprintf("%.2f", totalSales)}
		viewPizzaSalesSlice = append(viewPizzaSalesSlice, viewPizzaSales)
	}

	return viewPizzaSalesSlice
}

func getTotalSales(viewPizzaSalesSlice []viewPizzaSales) float64 {

	totalSales := 0.0
	for _, v := range viewPizzaSalesSlice {
		totalSales = totalSales + v.TotalSales
	}

	return totalSales
}
