package main

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func signup(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	clientMsg := ""
	username := ""
	password := ""
	cmfpassword := ""
	firstname := ""
	lastname := ""

	var myUser user

	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		username = req.FormValue("username")
		password = req.FormValue("password")
		cmfpassword = req.FormValue("cmfpassword")
		firstname = req.FormValue("firstname")
		lastname = req.FormValue("lastname")

		if username != "" {
			// check if username exist/ taken
			if _, ok := mapUsers[username]; ok {
				clientMsg = "Username already taken. "
			} else {
				if password != cmfpassword {
					clientMsg = "Confirm Password must be the same as Password. "
				} else {
					// create session
					id, _ := uuid.NewV4()
					myCookie := &http.Cookie{
						Name:  "myCookie",
						Value: id.String(),
					}
					http.SetCookie(res, myCookie)
					mapSessions[myCookie.Value] = username

					bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
					if err != nil {
						clientMsg = "Internal server error. "
					} else {
						myUser = user{username, bPassword, firstname, lastname}
						mapUsers[username] = myUser

						// redirect to main index
						http.Redirect(res, req, "/", http.StatusSeeOther)
						return
					}
				}
			}
		}
	}

	data := struct {
		User      user
		UserName  string
		FirstName string
		LastName  string
		ClientMsg string
	}{
		myUser,
		username,
		firstname,
		lastname,
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "signup.gohtml", data)
}

func logout(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")
	// delete the session
	delete(mapSessions, myCookie.Value)
	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func getUser(res http.ResponseWriter, req *http.Request) user {
	// get current session cookie
	myCookie, err := req.Cookie("myCookie")
	// no cookie found
	if err != nil {
		// add new cookie
		id, _ := uuid.NewV4()
		myCookie = &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}

	}
	http.SetCookie(res, myCookie)

	// if the user exists already, get user
	var myUser user
	if username, ok := mapSessions[myCookie.Value]; ok {
		myUser = mapUsers[username]
	}

	return myUser
}

func alreadyLoggedIn(req *http.Request) bool {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false
	}
	username := mapSessions[myCookie.Value]
	_, ok := mapUsers[username]
	return ok
}

func edituser(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	clientMsg := ""
	username := myUser.Username
	password := ""
	cmfpassword := ""
	firstname := myUser.First
	lastname := myUser.Last

	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		username = req.FormValue("username")
		password = req.FormValue("password")
		cmfpassword = req.FormValue("cmfpassword")
		firstname = req.FormValue("firstname")
		lastname = req.FormValue("lastname")

		if password != cmfpassword {
			clientMsg = "Confirm Password must be the same as Password. "
		} else {
			bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				clientMsg = "Internal server error. "
			} else {
				myUser = user{username, bPassword, firstname, lastname}
				mapUsers[username] = myUser

				// redirect to main index
				http.Redirect(res, req, "/", http.StatusSeeOther)
				return
			}
		}
	}

	data := struct {
		User      user
		UserName  string
		FirstName string
		LastName  string
		ClientMsg string
	}{
		myUser,
		username,
		firstname,
		lastname,
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "edituser.gohtml", data)
}

func deleteuser(res http.ResponseWriter, req *http.Request) {
	myUser := getUser(res, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	clientMsg := ""

	if req.Method == http.MethodPost {
		username := req.FormValue("username")

		if username == "admin" {
			clientMsg = "You cannot delete this account. "
		} else {
			delete(mapUsers, username)

			clientMsg = "User [" + username + "] deleted successfully. "
		}
	}

	data := struct {
		User      user
		MapUsers  map[string]user
		CntUsers  int
		ClientMsg string
	}{
		myUser,
		mapUsers,
		len(mapUsers),
		clientMsg,
	}

	tpl.ExecuteTemplate(res, "deleteuser.gohtml", data)
}
