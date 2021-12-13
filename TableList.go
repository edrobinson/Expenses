/*
    This is a utility page that calls up a list of records from a selected 
    area - expenses, catagories or groups.
    
    The list responds to a click and calls for a read by primary key.
*/
package main

import (
	"database/sql"
	_"encoding/json"
	"fmt"
    _"io/ioutil"
	_"log"
	"net/http"
	_"strconv"
	_"strings"
    "text/template"
)

//Template data struct
type tpldata struct {
    Title string
    RecList string
}

var data tpldata

//Page server func
func TableListIndexHandler(w http.ResponseWriter, r *http.Request) {
    data.Title = "Table Lister"
    data.RecList = ""
	//Parse and execute the template file
	tmpl, err := template.ParseFiles("./templates/TableList.html",
                                     "./templates/header.tpl",
                                     "./templates/nav.tpl",
                                     "./templates/footer.tpl")

	WebCheck(err, " Failed to parse table list template. ",w)

	err = tmpl.Execute(w, data)
	WebCheck(err, " Failed to execute table list page template",w)
}

//Create and Display the list
func TableListHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    //Get the selected table name
    tbl := r.PostForm.Get("listtbl")
    qry := "select * from " + tbl
    rows,_ := DB.Query(qry)
        switch tbl{
            case "expense":
                scanExpenses(rows, w)
            case "catagories":
                scanCatagories(rows, w)
            case "groups":
                scanGroups(rows, w)
        }
 }                        

//Scan the rows from the table query and make a select of them
func scanExpenses(rows *sql.Rows, w http.ResponseWriter){
	var list []Expense
	for rows.Next() {
		err := rows.Scan(&ExRes.Id, &ExRes.Exdate, &ExRes.Examt, &ExRes.Excatagory, &ExRes.Exdescription)
		WebCheck(err, "Expense row scan failed.", w)
		list = append(list, ExRes)
	}
    
    //Construct an HTML select.
    //The select onchange hides the select element and sets up a read on the table
    var s string = fmt.Sprintf("<select class='form-select' id=%q onchange=ReadFromSelect(%q)>\n","recordlist","expense")
    s += fmt.Sprintf("<option> </option>\n")
    var rd string
    for _, res := range list {
        dt := res.Exdate[0:10]
        rd = fmt.Sprintf("%s  %5.2f  %s  %s </option>\n",dt, res.Examt,res.Excatagory, res.Exdescription)
        s += fmt.Sprintf("<option value=%d>%s",res.Id, rd)
    }
    s += "</select>"
    
    //Process the page template
    data.Title = "Record Selection"
    data.RecList = s
	//Parse and execute the template file
	tmpl, err := template.ParseFiles("./templates/RecordsList.html",
                                     "./templates/header.tpl",
                                     "./templates/nav.tpl",
                                     "./templates/footer.tpl")

	WebCheck(err, " Failed to parse table list template. ",w)

	err = tmpl.Execute(w, data)
	WebCheck(err, " Failed to execute table list page template",w)
}

func scanCatagories(rows *sql.Rows, w http.ResponseWriter) string{
	var clses []Catagories
	//Table var to scan into
	var res Catagories

	for rows.Next() {
		err := rows.Scan(&res.Id, &res.Cdescription)
		WebCheck(err, "Row scan failed",w)
		clses = append(clses, res)
	}
    var s string
    return s
}

func scanGroups(rows *sql.Rows, w http.ResponseWriter) string{
	var clses []Groups
	//Table var to scan into
	var res Groups
	for rows.Next() {
		err := rows.Scan(&res.Id, &res.Gdescription)
		Check(err, "Row scan failed")
		clses = append(clses, res)
	}
    var s string
    return s
}


