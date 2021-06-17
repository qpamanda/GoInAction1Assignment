package main

import (
	aPizza "GoInAction1Assignment/adminPizza"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
)

func generatePizzaNo() {

	defer wg.Done()

	mutex.Lock()
	runtime.Gosched()
	// Increment PizzaNo global variable by 1
	newPizzaNo = newPizzaNo + 1
	mutex.Unlock()
}

func addpizza(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	clientMsg := ""
	pizzaName := ""
	inputPizzaPrice := ""

	if req.Method == http.MethodPost {
		pizzaName = req.FormValue("pizzaname")

		if pizzaName != "" {
			inputPizzaPrice = req.FormValue("pizzaprice")

			pizzaPrice, err := validatePizzaPrice(inputPizzaPrice)

			if err != nil || pizzaPrice == 0 {
				clientMsg = "Please enter a valid Pizza Price."
			} else {
				wg.Add(1)
				go generatePizzaNo()
				wg.Wait()

				pizzaNo := newPizzaNo
				pizzaList.AddPizza(pizzaNo, pizzaName, pizzaPrice)

				clientMsg = fmt.Sprintf("%s @ $%.2f added successfully.\n", pizzaName, pizzaPrice)
			}
		} else {
			clientMsg = "Please enter Pizza Name."
		}
	}

	data := struct {
		User       user
		PizzaName  string
		PizzaPrice string
		ClientMsg  string
	}{
		myUser,
		pizzaName,
		inputPizzaPrice,
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "addpizza.gohtml", data)
}

func editpizza(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	clientMsg := ""
	pizzaNo := 0
	fPizzaName := ""
	fPizzaPrice := ""
	pizzaPrice := 0.0

	viewPizzaSlice := make([]viewPizza, 0)
	pizzaSlice, err := pizzaList.GetAllPizza()

	if err != nil {
		clientMsg = "There are no pizza on the menu today. "
	}

	if req.Method == http.MethodPost {
		fPizzaNo := req.FormValue("pizzano")
		pizzaNo, _ := strconv.Atoi(fPizzaNo)

		bPizzaInOrder, _ := checkPizzaInOrder(pizzaNo)
		if bPizzaInOrder {
			clientMsg = "Orders have been made on the selected pizza. You are not allowed to edit this pizza."
		} else {
			selectedPizza, err := pizzaList.SearchPizza(pizzaNo)

			if err != nil {
				clientMsg = "Cannot edit this pizza."
			} else {
				fPizzaName = req.FormValue("pizzaname")
				fPizzaPrice = req.FormValue("pizzaprice")

				if fPizzaName == "" {
					fPizzaName = selectedPizza.PizzaName
				}

				if fPizzaPrice == "" {
					fPizzaPrice = fmt.Sprintf("%.2f", selectedPizza.PizzaPrice)
				}

				if fPizzaName == selectedPizza.PizzaName && fPizzaPrice == fmt.Sprintf("%.2f", selectedPizza.PizzaPrice) {
					clientMsg = "No changes made on the selected pizza."
				} else {
					pizzaPrice, err := validatePizzaPrice(fPizzaPrice)

					if err != nil || pizzaPrice == 0 {
						clientMsg = "Please enter a valid Pizza Price."
					} else {
						pizzaList.EditPizza(pizzaNo, fPizzaName, pizzaPrice)

						clientMsg = fmt.Sprintf("%s @ $%s updated successfully.\n", fPizzaName, fPizzaPrice)
					}
				}

				for _, v := range pizzaSlice {
					if pizzaNo == v.PizzaNo {
						viewPizza := viewPizza{pizzaNo, fPizzaName, pizzaPrice, fPizzaPrice, "Selected"}
						viewPizzaSlice = append(viewPizzaSlice, viewPizza)
					} else {
						viewPizza := viewPizza{v.PizzaNo, v.PizzaName, v.PizzaPrice, fmt.Sprintf("%.2f", v.PizzaPrice), ""}
						viewPizzaSlice = append(viewPizzaSlice, viewPizza)
					}
				}
			}
		}
	}

	if len(viewPizzaSlice) == 0 {
		for _, v := range pizzaSlice {
			viewPizza := viewPizza{v.PizzaNo, v.PizzaName, v.PizzaPrice, fmt.Sprintf("%.2f", v.PizzaPrice), ""}
			viewPizzaSlice = append(viewPizzaSlice, viewPizza)
		}
	}

	data := struct {
		User           user
		ViewPizzaSlice []viewPizza
		CntPizza       int
		PizzaNo        int
		PizzaName      string
		PizzaPrice     string
		ClientMsg      string
	}{
		myUser,
		viewPizzaSlice,
		len(viewPizzaSlice),
		pizzaNo,
		fPizzaName,
		fPizzaPrice,
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "editpizza.gohtml", data)
}

func deletepizza(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	clientMsg := ""
	pizzaNo := 0

	if req.Method == http.MethodPost {
		fPizzaNo := req.FormValue("pizzano")
		pizzaNo, _ := strconv.Atoi(fPizzaNo)

		bPizzaInOrder, _ := checkPizzaInOrder(pizzaNo)
		if bPizzaInOrder {
			clientMsg = "Orders have been made on the selected pizza. You are not allowed to delete this pizza."
		} else {
			selectedPizza, err := pizzaList.SearchPizza(pizzaNo)

			if err != nil {
				clientMsg = "Cannot edit this pizza."
			} else {
				pizzaList.DeletePizza(pizzaNo)
				clientMsg = fmt.Sprintf("%s @ $%.2f deleted successfully.\n", selectedPizza.PizzaName, selectedPizza.PizzaPrice)
			}
		}
	}

	pizzaSlice, err := pizzaList.GetAllPizza()

	if err != nil {
		clientMsg = "There are no pizza on the menu today. "
	}

	data := struct {
		User           user
		ViewPizzaSlice []aPizza.Pizza
		CntPizza       int
		PizzaNo        int
		ClientMsg      string
	}{
		myUser,
		pizzaSlice,
		len(pizzaSlice),
		pizzaNo,
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "deletepizza.gohtml", data)
}

func viewpizza(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	clientMsg := ""

	pizzaSlice, err := pizzaList.GetAllPizza()
	viewPizzaSlice := make([]viewPizza, 0)

	if err != nil {
		clientMsg = "There are no pizza on the menu today. "
	} else {
		for _, v := range pizzaSlice {
			viewPizza := viewPizza{v.PizzaNo, v.PizzaName, v.PizzaPrice, fmt.Sprintf("%.2f", v.PizzaPrice), ""}
			viewPizzaSlice = append(viewPizzaSlice, viewPizza)
		}
	}

	data := struct {
		User           user
		ViewPizzaSlice []viewPizza
		CntPizza       int
		ClientMsg      string
	}{
		myUser,
		viewPizzaSlice,
		len(viewPizzaSlice),
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "viewpizza.gohtml", data)
}

func validatePizzaPrice(pizzaPrice string) (float64, error) {
	retValue, err := strconv.ParseFloat(pizzaPrice, 64)

	if err != nil {
		return retValue, errors.New(">> Invalid pizza price")
	}

	return retValue, nil
}

func checkPizzaInOrder(pizzaNo int) (bool, error) {
	// If there are no orders made means there are no pizza in any order, thus return false
	if orderQueue.IsEmpty() {
		return false, nil
	} else {
		return orderQueue.SearchPizzaInOrder(pizzaNo)
	}
}
