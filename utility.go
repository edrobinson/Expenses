/*
   Spending utilities

   Contains a number of useful functions for
   database related processes.

   Note the function Check. It does the checks on err and
   reports the file, line, dev's message and the err text.


   Func signature                                     Func function
   Check(e error, msg string)                         Error check with file, line and error report
   SnameToCamel(str string) (camelCase string         Field name to camel case
   ToLength(s string, sz int) string                  Pad a string to a length
   dbAccess()                                         Open  aglobal db connection
   Mmddyy2yymmdd(mmddyy string)                       Date conversion to mysql format
*/

package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	//	_ "github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

//Form field struct to capture form fields and values from the client.
//See DecodeAjaxData() below.
type Formval struct {
	Name  string
	Value string
}

//Instance for json unmarshalling
var fval Formval

//Inputs from the ajax request
var Fvals []Formval

//Requested record id and crud operation
//from the client form.
var RecId string
var OpCode string

//Global db reference
var DB *sql.DB

var ConString string

//Init accesses the database
func init() {
	//Open a DB connection
	dbAccess()
}

/*********************** Error checking functions ***************************/

// Error check provides expanded error
//information to the log and exits the program.
func Check(e error, msg string) {
	//If no error just return
    if e == nil {
		return
	}
    
	//Get the run time info
    _, fn, line, _ := runtime.Caller(1)
	_, file := path.Split(fn)
    
	//Log it to the command tool
	log.Printf("Error in File: %s\n", file)
	log.Println("Dev Msg", msg)
	log.Printf("At Line #: %d\n", line-1)
	log.Printf("Error: %v\n", e)
	os.Exit(1)
}

// Same as Check(), logginh the info then also sends the info to the browser causing an alert.
func WebCheck(e error, msg string, w http.ResponseWriter) {
    if e == nil {
		return
	}
    
	//Get the run time info
    _, fn, line, _ := runtime.Caller(1)
	_, file := path.Split(fn)
	//Log it to the command tool
	log.Printf("Error in File: %s\n", file)
	log.Println("Dev Msg", msg)
	log.Printf("At Line #: %d\n", line-1)
	log.Printf("Error: %v\n", e)
    //Log it to the browser
	var s string
	s += fmt.Sprintf("Error in File: %s\n", file)
	s += fmt.Sprintf("Dev Msg: %s\n", msg)
	s += fmt.Sprintf("At Line #: %d\n", line-1)
	s += fmt.Sprintf("Error Msg: %s\n", e.Error())
	w.Write([]byte(s))
    //and exit
    os.Exit(1)
    
}

/***********************************************************************************/

//Convert Snake case - underscores to Camel Case - Capitalized words run together
//See: https://www.socketloop.com/tutorials/golang-underscore-or-snake-case-to-camel-case-example
func SnameToCamel(str string) (camelCase string) {
	return strcase.ToCamel(str)
}

//Pad a string with spaces on the right
//to the length specified by sz.
func ToLength(s string, sz int) string {
	for len(s) < sz {
		s += " "
	}
	return s
}

//Make connection to the database
//Save it in the global DB so everything can use it.
//Called from init()
func dbAccess() {

	//Make the connection string depending on
	//the presence of a password
	var con string
	con = Userid + ":" + Password + "@tcp(127.0.0.1:3306)/" + Dbname // + "?parseTime=true"

	//Replace the fake password with empty string
	con = strings.Replace(con, "zzz", "", 1)

	//Pass it to the template
	ConString = con

	// Open the database and assign the db pointer to the global DB
	//db, err := sql.Open("mysql", con)
	db, err := sql.Open("sqlite3", "database/expenses.sqlite3")
	Check(err, "Failed to open database ")

	//Save to global db reference
	DB = db
}

//Convery mm/dd/yyyy date to yyyy-mm-dd for mysql tables
func Mmddyy2yymmdd(mmddyy string) string {
	var p []string
	p = strings.Split(mmddyy, "/")
	return p[2] + "-" + p[0] + "-" + p[1]
}

//DecodeAjaxData reads tha ajax json data from the saved file,
//unmarshalls it, saves the op code and record id. Then
//it converts the remaining json into set of form field names and values
//for processing.
func DecodeAjaxData(filename string) (error, []Formval) {
	var fvals []Formval //Slice of the form values struct receives unmarshal
	var tvals []Formval //Slice of the formvals to return
	var pfval Formval   //A Formval struct object

	//Load the saved json inputs
	file, err := ioutil.ReadFile(filename)
	Check(err, "Error loading result json file")

	//Extract the measurement records
	err = json.Unmarshal(file, &fvals)
	if err != nil {
		return errors.New("No ajax data received"), nil
	}

	//Scan the json and build the struct of field names and values
	//Each field is {"name":"value"}
	//The 0th item is the crud op code name
	//The 1st item is the record id as appropriate
	//The rest are the input values from the form.
	OpCode = ""
	RecId = ""
	for i := range fvals {
		if i == 0 {
			OpCode = fvals[i].Value
		} else if i == 1 {
			RecId = fvals[i].Value
		} else {
			//Break out the name and input value
			pfval.Name = fvals[i].Name
			pfval.Value = fvals[i].Value
			//Add to the slice of field defs
			tvals = append(tvals, pfval)
		}
	}
	return nil, tvals
}

/************************* Begin DB Procedures ********************************/

//Generic read first
func ReadFirst(tbl string) *sql.Row {
	qry := fmt.Sprintf("select * from %s limit 1", tbl)
	row := DB.QueryRow(qry)
	return row
}

//Generic read last
func ReadLast(tbl string) *sql.Row {
	qry := fmt.Sprintf("select * from %s order by id desc limit 1", tbl)
	row := DB.QueryRow(qry)
	return row
}

//Generic read by id
func ReadExact(tbl string, recid string) *sql.Row {
	qry := fmt.Sprintf("select * from %s where id = %s", tbl, recid)
	row := DB.QueryRow(qry)
	return row
}

//Generid read next record
func ReadNext(tbl string, recid string) *sql.Row {
	qry := fmt.Sprintf("select * from %s where id = (select min(id) from %s  where id > %s)", tbl, tbl, recid)
	row := DB.QueryRow(qry)
	return row
}

//Generic read previous record query
func ReadPrev(tbl string, recid string) *sql.Row {
	qry := fmt.Sprintf("select * from %s where id = (select max(id) from %s where id <  %s)", tbl, tbl, recid)
	row := DB.QueryRow(qry)
	return row
}

//Generic Query
func ReadRows(qry string) (*sql.Rows, error) {
	rows, err := DB.Query(qry)
	return rows, err
}

//Global sql insert
func InsertNewRecord(w http.ResponseWriter, tblname string) {
	//Insert list of field names
	var fields string = ""

	//Insert list of field values
	var values string = ""

	for i := range Fvals { //Fvals were created by the Ajax decoder
		if Fvals[i].Name == "id" {
			continue
		} //Skip the primary keyfield
		if Fvals[i].Value == "" {
			continue
		} //Skip empty fields
		fields += Fvals[i].Name + ","         //Field names list
		values += "'" + Fvals[i].Value + "'," //Value list
	}

	//Remove the trailing comma in the fields list
	fields = strings.TrimRight(fields, ",")
	//Add the surrounding prens
	fields = "(" + fields + ")"

	if fields == "()" {
		http.Error(w, " No fields were filled out. Please retry.", http.StatusInternalServerError)
		return
	}

	//Remove the trailing comma in the values list
	values = strings.TrimRight(values, ",")
	//Add prens
	values = "(" + values + ")"

	//Form the query string
	var qry = "insert into " + tblname + " " + fields + " values " + values
	// log.Println("Insert Query: ",qry)
	//GO for it
	result, err := DB.Exec(qry)
	WebCheck(err, " Insert statement failed.", w)

	a, _ := result.RowsAffected()

	if a < 1 {
		var msg string = "Insert Failed"
		w.Write([]byte(msg))
		return
	}

	//Send back the new record id
	lid, err := result.LastInsertId()
	msg := "Record inserted.. Id: " + strconv.FormatInt(lid, 10)
	w.Write([]byte(msg))
}

//Global sql record delete
func DeleteRecord(w http.ResponseWriter, tblname string) {

	stmt, err := DB.Prepare("delete from spending where id = " + RecId)
	WebCheck(err, " Prepare statement failed in delete handler", w)

	//Execute the delete
	result, err := stmt.Exec()
	WebCheck(err, " Exec statement failed in delete handler", w)

	// affected rows
	a, err := result.RowsAffected()
	WebCheck(err, "RowsAffected failed in delete handler", w)
	if a > 0 {
		w.Write([]byte("Record " + RecId + " Was Deleted."))
	} else {
		w.Write([]byte("Delete Failed - No records were affected."))
	}
}

func UpdateRecord(w http.ResponseWriter, tblname string) {
	//Iterate over the fields and create the colname = expression clause
	var fields string
	for i := range Fvals {
		if Fvals[i].Name == "id" {
			continue
		} //Skip the primary keyfield
		fields += Fvals[i].Name + " = '" + Fvals[i].Value + "',"
	}

	//Trim off the trailing comma
	fields = strings.TrimRight(fields, ",")

	//Create the update statement
	var qry string = "update " + tblname + " set " + fields + " where id = " + RecId

	//Execute the query
	result, err := DB.Exec(qry)
	WebCheck(err, " Update statement failed.", w)

	//Did anything happen?
	a, _ := result.RowsAffected()

	if a < 1 {
		var msg string = "Update Failed to affect any rows. "
		w.Write([]byte(msg))
		return
	} else {
		msg := "Record updated."
		w.Write([]byte(msg))
	}

}

//This util function provides an html options list
//for the passed table name ordered on the orderby param.
//It is to be used as a record lookup device
func TableList(tbl string, orderby string) string {
	type qres struct {
		id   string
		desc string
	}
	var results []qres
	var res qres

	qry := fmt.Sprintf("select id, %s from %s order by %s", orderby, tbl, orderby)
	rows, err := DB.Query(qry)
	Check(err, " Table list select query failed. ")
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&res.id, &res.desc)
		Check(err, "Row scan failed")
		results = append(results, res)
	}
	//Construct the options list
	var s string

	for _, res := range results {
		s += fmt.Sprintf("<option value=%d>%s</option>\n", res.id, res.desc)
	}

	return s
}

/******************************* Lookup Service ****************************************
    This is a consolidated record lookup function. It looks up records per the passed parameters
    and sends the response to the browser. The callin function is reduced to just 2 lines.
    
    Parameters:
    ResponseWriter to send the response to the browser.
    
    tablename is the name of the table being queried.
    
    orderby string is on or more comma seperated field names to order the query.
    
    dolike is a boolean specifying if an sql like query is to be run.
    
    Flow:
    1. The form data is scanned for a possible field and value to use.
    2. The query is constructed and run.
    3. The rows are passed to a table specific scan routine that builds
       the html select.
    4. The data is marshalled to json and written to the browser.
    
    Error checking is done at various points as needed.
    
*/
func LookupService(w http.ResponseWriter, tablename string, orderby string, dolike bool) error{
	//Was a luby - "lookup by" - radio even given?
	var f Formval
	var fld string = ""

	//Scan for a "luby" radio button
	for _, f = range Fvals {
		if f.Name == "luby" {
			fld = f.Value //The field to query on
			break
		}
	}

	//Scan for the field and get it's value
	var fldVal string
	for _, f = range Fvals {
		if f.Name == fld {
			fldVal = f.Value //The value to query on
			break
		}
	}

	//Query for matching rows
	var qry string
    if fld == "" {
        //Do a select like if no field is selected
        qry = fmt.Sprintf("select * from %s  order by %s", tablename, orderby)
    //Query like called for?
    }else if dolike{
		fldVal = fmt.Sprintf("%%%s%%", fldVal) //In %s for like query
		qry = fmt.Sprintf("select * from %s where %s like %q",tablename, fld, fldVal)
	//Query on the selected field
    } else {
		qry = fmt.Sprintf("select * from %s where %s = %q",tablename, fld, fldVal)
	}

	//Run the query
    rows, err := ReadRows(qry)
    WebCheck(err, " Lookup query failed.",w)
    
    //Generate the select for the passed table
    var s string
    switch tablename{
        case "expense":
            s = MakeExpenseSelect(rows,w)
        case "catagories":
            s = MakeCatagorySelect(rows,w)
        case "groups":
            s = MakeGroupsSelect(rows,w)
    }
    
    //Add the blank and cancel options to the table options
    var ops string = fmt.Sprintf("<option value=''></option>", 0)
    //Add the cancel option
    ops += "<option value='cancel'>Cancel</option>"
    s = ops + s
    
    //Marshall the output to json
    var j []byte
	j, err = json.Marshal(s)
	WebCheck(err, " json marshal failed", w)
    
    //Send it to the browser
	w.Write([]byte(j))
    
    return  nil
 }

/*****************************************************************
    The next 3 functions accept the rows from the lookup service
    and create an HTML <select> specific to the table.
    
    All 3 are the same except for the scan and option build.
******************************************************************/
 
func MakeExpenseSelect(rows *sql.Rows,w http.ResponseWriter) string{
    var s string = ""
    //Scan the rows and add the options
    for rows.Next() {
		err := rows.Scan(&ExRes.Id, &ExRes.Exdate, &ExRes.Examt, &ExRes.Excatagory, &ExRes.Exdescription)
		WebCheck(err, " Lookup rows scan failed.",w)
        var sdate = ExRes.Exdate[0:10]
        
		s += fmt.Sprintf("<option value=%s>%s  $%.2f %s:  %s</option>", strconv.Itoa(ExRes.Id),sdate, ExRes.Examt, ExRes.Excatagory, ExRes.Exdescription)
	}
    //Complete the select
    s += "</select>"
    
    return s
 }

func MakeCatagorySelect(rows *sql.Rows, w http.ResponseWriter) string {
    var s string = ""
    for rows.Next() {
        err := rows.Scan(&CatRes.Id, &CatRes.Cdescription)
        WebCheck(err, " Lookup rows scan failed.",w)
        s += fmt.Sprintf("<option value=%s>%s</option>", strconv.Itoa(CatRes.Id), CatRes.Cdescription)
    }
    s += "</select>"
    return s
}

func MakeGroupsSelect(rows *sql.Rows, w http.ResponseWriter) string {
    var s string = ""
    for rows.Next() {
        err := rows.Scan(&GRes.Id, &GRes.Gdescription)
        WebCheck(err, " Lookup rows scan failed.",w)
        s += fmt.Sprintf("<option value=%s>%s</option>", strconv.Itoa(GRes.Id), GRes.Gdescription)
    }
    s += "</select>"
    return s
}
