/*
   This is the main go for an expense tracking app.

   It runs as a web server.

   The mysql tables are in the db named "expenses."

   The transactions are in table "expense."

   The spending catagories are in the table "catagories."

   The main go defines the handlers and starts the server.
   The handlers actions are the table names.

*/

package main

import (
	"log"
	"net/http"
)

//Setup all of the net handlers
func init() {
	//Handle js and css files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	//Serve html files
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("html"))))

	//Home page display
	http.HandleFunc("/", defaultHandler)

	//Two handlers for each of the tables
	//1. The table's index handler to display the template page
	//2. The table's processes handler

	http.HandleFunc("/catagoriesindex", CatagoriesIndexHandler)
	http.HandleFunc("/catagories", CatagoriesAjaxHandler)

	http.HandleFunc("/expenseindex", ExpenseIndexHandler)
	http.HandleFunc("/expense", ExpenseAjaxHandler)

	http.HandleFunc("/groupsindex", GroupsIndexHandler)
	http.HandleFunc("/groups", GroupsAjaxHandler)

	http.HandleFunc("/reportsindex", ReportsIndexHandler)
	http.HandleFunc("/reports", ReportsHandler)
    
	http.HandleFunc("/tablelistindex", TableListIndexHandler)
	http.HandleFunc("/tablelist", TableListHandler)
    
}

func main() {
	log.Println("Listening on port 8080.")
	http.ListenAndServe(":8080", nil)
}

//The home page is the spending table html file
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	ExpenseIndexHandler(w, r)
}
