/*
  The file handles the main function of the app - processing
 the expense inputs.
*/
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	_"strconv"
	_ "strings"
	"text/template"
	_ "time"
)

//Database: expenses   Table: expense
type Expense struct {
	Id            int     `json:"id"`
	Exdate        string  `json:"exdate"`
	Examt         float32 `json:"examt"`
	Excatagory    string  `json:"excatagory"`
	Exdescription string  `json:"exdescription"`
}

//Field list for scan statements
var ExScanline string = "&ExRes.Id, &ExRes.Exdate, &ExRes.Examt, &ExRes.Excatagory, &ExRes.Exdescription"

//Table var to scan into
var ExRes Expense

//The table's primary key
const ExPk string = "id"

func ExpenseAjaxHandler(w http.ResponseWriter, r *http.Request) {
	//Get the body of the response - contains the table form values, record id and opcode
	//as json. data is read as byte type.
	data, err := ioutil.ReadAll(r.Body)
	WebCheck(err, "Error getting the data file from response body.", w)

	//Write it to a file
	var filename = "jsonfiles/expenses" + ".json"

	//Save the input for reference and further processing
	err = ioutil.WriteFile(filename, data, 0775)
	WebCheck(err, "Error saving the result data file", w)

	//Extract the input data - in utility.go
	err, s := DecodeAjaxData(filename)
	WebCheck(err, "Failed to decode the ajax data", w)

	//Save the table data
	Fvals = s

	//What to do, what to do...
	switch OpCode {
	case "readfirst":
		ReadFirstExpense(w)
	case "readprev":
		ReadPrevExpense(w)
	case "read":
		ReadExpense(w)
	case "readnext":
		ReadNextExpense(w)
	case "readlast":
		ReadLastExpense(w)
	case "readlike":
		ReadLikeExpense(w)
	case "lookup":
		LookupExpense(w)
	case "insert":
		InsertNewRecord(w, "expense")
	case "update":
		UpdateRecord(w, "expense")
	case "delete":
		DeleteRecord(w, "expense")
	default:
		log.Println("Invalid opcode sent: ", OpCode)
		http.Error(w, "Invalid op code"+OpCode, http.StatusInternalServerError)
	}
}

//Serve the spending table editor html page
func ExpenseIndexHandler(w http.ResponseWriter, r *http.Request) {
	//Template data struct
	type tpldata struct {
		Catagoryselect string
		Title          string
	}

	var data tpldata
	data.Catagoryselect = CatagoriesSelectList(0)
	data.Title = "Expense Editor"

	//Parse and execute the template file
	tmpl, err := template.ParseFiles("./templates/Expenses.html",
		"./templates/header.tpl",
		"./templates/nav.tpl",
		"./templates/buttons.tpl",
		"./templates/footer.tpl")
	WebCheck(err, " Failed to parse spending template. ", w)

	err = tmpl.Execute(w, data)
	WebCheck(err, " Failed to execute spending template", w)
}

/************************Begin DB Read Functiond********************************/
func ReadExpense(w http.ResponseWriter) {
	row := ReadExact("expense", RecId)
	CompleteExpenseRead(w, row)
}

func ReadNextExpense(w http.ResponseWriter) {
	row := ReadNext("expense", RecId)
	CompleteExpenseRead(w, row)
}

func ReadPrevExpense(w http.ResponseWriter) {
	row := ReadPrev("expense", RecId)
	CompleteExpenseRead(w, row)
}

func ReadFirstExpense(w http.ResponseWriter) {
	row := ReadFirst("expense")
	CompleteExpenseRead(w, row)
}

func ReadLastExpense(w http.ResponseWriter) {
	row := ReadLast("expense")
	CompleteExpenseRead(w, row)
}

//ReadLikeExpense does a select where x like %y%
//and returns a list of qualifiers for the user to choose from.
func ReadLikeExpense(w http.ResponseWriter) {
	//Field and value to query
	var fld string = ""
	var val string = ""

	//Look for a field with something in it
	for _, f := range Fvals {
		if f.Value != "" {
			fld = f.Name
			val = f.Value
			break
		}
	}

	//Return no field found
	if fld == "" {
		io.WriteString(w, "No non-blank field was found\n")
		return
	}

	//The usual sql stuff

	var qry string
	qry = "select * from spending where " + fld + ` like "%` + val + `%"`

	rows, err := DB.Query(qry)
	WebCheck(err, " Read Like query failed.", w)
	defer rows.Close()

	var list []Expense

	for rows.Next() {
		err = rows.Scan(&ExRes.Id, &ExRes.Exdate, &ExRes.Examt, &ExRes.Excatagory, &ExRes.Exdescription)
		WebCheck(err, "Read like row scan failed.", w)
		list = append(list, ExRes)
	}
	if len(list) == 0 {
		w.Write([]byte("No records Found"))
		return
	}

	var s string
	//Build a list of html options to return
	for _, ExRes := range list {
		s += fmt.Sprintf("<option value=%s>%s</option>\n", ExRes.Exdescription, ExRes.Exdescription)
	}
	w.Write([]byte(s))
}

//Code common to all read ops.
func CompleteExpenseRead(w http.ResponseWriter, row *sql.Row) {
	err := row.Scan(&ExRes.Id, &ExRes.Exdate, &ExRes.Examt, &ExRes.Excatagory, &ExRes.Exdescription)

	if err == sql.ErrNoRows {
		http.Error(w, " Record not found.", http.StatusInternalServerError)
		return
	} else if err != nil {
		WebCheck(err, " Failed to read the requested record.", w)
	}

	//Get the ExRes.Exdate value and extract the yyyy-mm-dd from it
	//and replace the Esdate value. Make it show up in the browser
	d := ExRes.Exdate
	d = d[0:10] //YYYY-MM-DD
	ExRes.Exdate = d

	//Marshal the db values to json
	var s []byte
	s, err = json.Marshal(ExRes)
	WebCheck(err, "Marshal Error", w)

	//Send the json to the browser
	w.Write([]byte(s))
}

/****************************** End of DB Read Functions ***********************/

//Do a lookup on the expense table.
//The utility service handles all of the lookup processing and
//responding to the request.
//The result is an HTML select list of records on screen. 
//When the user selects one, the read function is called.
func LookupExpense(w http.ResponseWriter) {
    err := LookupService(w, "expense", "excatagory, exdate", true) 
	WebCheck(err, " Lookup service failed.",w)
}    
