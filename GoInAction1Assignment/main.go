package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"golang.org/x/crypto/bcrypt"

	aOrder "GoInAction1Assignment/adminOrder"
	aPizza "GoInAction1Assignment/adminPizza"
)

var (
	wg    sync.WaitGroup
	mutex sync.Mutex

	tpl         *template.Template
	mapUsers    = map[string]user{}
	mapSessions = map[string]string{}

	standardPizza = []string{"Hawaiian Pizza", "Cheese Pizza", "Pepperoni Pizza", "Chicken Pizza", "Vegan Pizza"}

	newPizzaNo = len(standardPizza)

	newOrderNo = 1000

	bFirst = true

	pizzaList = &aPizza.Linkedlist{
		Head: nil,
		Size: 0,
	}

	orderQueue = &aOrder.Queue{
		Front: nil,
		Back:  nil,
		Size:  0,
	}

	completedOrderSlice = make([]viewOrder, 0)
)

const (
	standardPrice = 10.90 // standard price of all pizza start at $10.90
	maxOrderQty   = 5     // max order is 5 per selected pizza
)

type user struct {
	Username string
	Password []byte
	First    string
	Last     string
}

type viewOrderItem struct {
	ItemNo     int
	PizzaNo    int
	PizzaName  string
	PizzaPrice string
	OrderQty   int
	Checked    string
	ErrorMsg   string
}

type viewOrder struct {
	IdxNo          int
	OrderNo        int
	ViewOrderItems []viewOrderItem
	TotalCost      string
	UserName       string
}

type viewPizzaSales struct {
	PizzaNo     int
	PizzaName   string
	OrderQty    int
	TotalSales  float64
	STotalSales string
}

type viewPizza struct {
	PizzaNo     int
	PizzaName   string
	PizzaPrice  float64
	SPizzaPrice string
	Selected    string
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	mapUsers["admin"] = user{"admin", bPassword, "admin", "admin"}
	pizzaList.CreateStartMenu(standardPizza, standardPrice)
}

/* Print a divider line to segregate the sections for easy viewing */
func printDividerLine() {
	fmt.Println("------------------------------------------------------------")
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/edituser", edituser)
	http.HandleFunc("/deleteuser", deleteuser)
	http.HandleFunc("/addorder", addorder)
	http.HandleFunc("/editorder", editorder)
	http.HandleFunc("/vieworders", vieworders)
	http.HandleFunc("/completeorder", completeorder)
	http.HandleFunc("/pizzasales", pizzasales)
	http.HandleFunc("/addpizza", addpizza)
	http.HandleFunc("/editpizza", editpizza)
	http.HandleFunc("/deletepizza", deletepizza)
	http.HandleFunc("/viewpizza", viewpizza)
	http.HandleFunc("/logout", logout)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	// set listen port
	err := http.ListenAndServe(":5221", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func index(res http.ResponseWriter, req *http.Request) {
	bFirst = true
	clientMsg := ""

	// process form submission
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")

		// check if user exist with username
		myUser, ok := mapUsers[username]

		if !ok {
			clientMsg = "Username and/or password do not match."
		} else {
			// Matching of password entered
			err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
			if err != nil {
				clientMsg = "Username and/or password do not match."
			} else {
				myCookie, err := req.Cookie("myCookie")
				if err != nil {
					clientMsg = "No cookie found."
				} else {
					http.SetCookie(res, myCookie)
					mapSessions[myCookie.Value] = username
				}
			}
		}
	}

	myUser := getUser(res, req)

	data := struct {
		User      user
		ClientMsg string
	}{
		myUser,
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "index.gohtml", data)
}
