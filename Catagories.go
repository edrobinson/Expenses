/*
   This go handles the catagories table in the expenses app.
   It provides a list of expense catagories for the expense
   table's desctiption field.

*/
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

//Database: expenses   Table: catagories
type Catagories struct {
	Id           int    `json:"id"`
	Cdescription string `json:"cdescription"`
}

//Field list for scan statements
var Scanline string = "&CatRes.Id, &CatRes.Cdescription"

//Table var to scan into
var CatRes Catagories

var CatPk string = "id"

func CatagoriesAjaxHandler(w http.ResponseWriter, r *http.Request) {
	//Get the body of the response - contains the table form values, record id and opcode
	//as json. data is read as byte type.
	data, err := ioutil.ReadAll(r.Body)

	WebCheck(err, "Error getting the data file from response body.", w)

	//Write it to a file
	var filename = "jsonfiles/catagories" + ".json"

	//Save the input for reference and further processing
	err = ioutil.WriteFile(filename, data, 0775)
	WebCheck(err, "Error saving the result data file", w)

	//Extract the input data
	err, s := DecodeAjaxData(filename)
	WebCheck(err, "Failed to decode the ajax data", w)

	//Save the table data
	Fvals = s
	//log.Println("Extracted:", string(s))

	switch OpCode {
	case "readfirst":
		ReadFirstCatagories(w)
	case "readprev":
		ReadPrevCatagories(w)
	case "read":
		ReadCatagories(w)
	case "readnext":
		ReadNextCatagories(w)
	case "readlast":
		ReadLastCatagories(w)
	case "lookup":
		LookupCatagories(w)
	case "insert":
		InsertCatagories(w)
	case "update":
		UpdateCatagories(w)
	case "delete":
		DeleteCatagories(w)
	default:
		log.Println("Invalid opcode sent: ", OpCode)
		http.Error(w, "Invalid op code"+OpCode, http.StatusInternalServerError)
	}
}

//Page server func
func CatagoriesIndexHandler(w http.ResponseWriter, r *http.Request) {
	//Template data struct
	type tpldata struct {
		Groupsselect string
		Title        string
	}

	var data tpldata
	data.Title = "Catagories Editor"
	data.Groupsselect = GroupsSelectList(0)
	//The account id anf product id are FKs of the table
	//so we get a set of html select options for both tables

	//Parse and execute the template file
	tmpl, err := template.ParseFiles("./templates/Catagorys.html",
		"./templates/header.tpl",
		"./templates/nav.tpl",
		"./templates/buttons.tpl",
		"./templates/footer.tpl")

	WebCheck(err, " Failed to parse catagories template. ", w)

	err = tmpl.Execute(w, data)
	WebCheck(err, " Failed to execute classes page template", w)
}

//Utility read function.
//Called by each crud function for validation
func CatagoriesLoadRecord() error {
	var whereclause string = CatPk + " = " + RecId
	var qry = "SELECT * FROM classes WHERE " + whereclause
	row := DB.QueryRow(qry)
	err := row.Scan(&CatRes.Id, &CatRes.Cdescription)
	return err
}

//CRUD Skeletons
//ReadProducts does a select where id = xx
func ReadCatagories(w http.ResponseWriter) {
	row := ReadExact("catagories", RecId)
	CompleteCatagoriesRead(w, row)
}

func ReadNextCatagories(w http.ResponseWriter) {
	row := ReadNext("catagories", RecId)
	CompleteCatagoriesRead(w, row)
}

func ReadPrevCatagories(w http.ResponseWriter) {
	row := ReadPrev("catagories", RecId)
	CompleteCatagoriesRead(w, row)
}

func ReadFirstCatagories(w http.ResponseWriter) {
	row := ReadFirst("catagories")
	CompleteCatagoriesRead(w, row)
}

func ReadLastCatagories(w http.ResponseWriter) {
	row := ReadLast("catagories")
	CompleteCatagoriesRead(w, row)
}

//Code common to all read ops after the utility db queries.
func CompleteCatagoriesRead(w http.ResponseWriter, row *sql.Row) {
	err := row.Scan(&CatRes.Id, &CatRes.Cdescription)

	if err == sql.ErrNoRows {
		http.Error(w, " Record not found.", http.StatusInternalServerError)
		return
	} else if err != nil {
		WebCheck(err, " Failed to read the requested record.", w)
	}

	//Marshal the db values to json
	var s []byte
	s, err = json.Marshal(CatRes)
	WebCheck(err, "Marshal Error", w)

	//Send the json to the browser
	w.Write([]byte(s))
}

func LookupCatagories(w http.ResponseWriter) {
    err := LookupService(w, "catagories", "cdescription", true) 
	WebCheck(err, " Lookup service failed.",w)
}

func InsertCatagories(w http.ResponseWriter) {
	//Prefix the Group name to the description
	addGroupToDescription()
	//Insert list of field names
	var fields string = ""

	//Insert list of field values
	var values string = ""

	for i := range Fvals {
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
	var qry = "insert into " + "catagories" + " " + fields + " values " + values

	//GO for it
	result, err := DB.Exec(qry)
	WebCheck(err, " Insert statement failed.", w)
	if result == nil {
		return
	}

	a, err := result.RowsAffected()
	WebCheck(err, " RowsAffected call failed.", w)

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

/*
   //AddGroupToDescription extracts the
   //group and description from the form values,
   //and makes group description and inserts back into
   //the form description field and emptys the group field
   // so it is not considered later.
   //The purpose of the group is to be able to group
   //the catagories for later processing
*/
func addGroupToDescription() {
	var gp, desc string //Group and description values
	var gpi, desci int  //Group and description index
	//
	for i := range Fvals {
		//Save the description and index
		if Fvals[i].Name == "cdescription" {
			desc = Fvals[i].Value
			desci = i
		}
		//Save the group and index
		if Fvals[i].Name == "group" {
			gp = Fvals[i].Value
			gpi = i
		}
	}
	//Form the combined description
	var newdesc string = gp + " " + desc

	Fvals[desci].Value = newdesc
	Fvals[gpi].Value = ""
}

func UpdateCatagories(w http.ResponseWriter) {
	//Get the record that the update call specifies
	//for comparison
	err := CatagoriesLoadRecord()

	//WebCheck for not found
	if err == sql.ErrNoRows {
		http.Error(w, " Record not found to update.", http.StatusInternalServerError)
		return
	} else {
		WebCheck(err, " Unable to find update candidate record", w)
	}

	//Iterate over the fields and create the colname = expression clause
	var fields string
	for i := range Fvals {
		if Fvals[i].Name == CatPk {
			continue
		} //Skip the primary keyfield
		fields += Fvals[i].Name + " = '" + Fvals[i].Value + "',"
	}

	//Trim off the trailing comma
	fields = strings.TrimRight(fields, ",")

	//Create the update statement
	var qry string = "update catagories set " + fields + " where " + CatPk + " = " + RecId

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

func DeleteCatagories(w http.ResponseWriter) {

	//Make certain the record exists
	err := CatagoriesLoadRecord()

	//WebCheck for not found
	if err == sql.ErrNoRows {
		http.Error(w, " Record not found to delete.", http.StatusInternalServerError)
		return
	} else {
		WebCheck(err, " Unable to find update candidate record", w)
	}

	stmt, err := DB.Prepare("delete from classes where " + CatPk + " = " + RecId)
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

//Make a list of select options
//The sel is the  value of the option to be selected
//Pass 0 to select none
func CatagoriesSelectList(sel int) string {
	var clses []Catagories
	//Table var to scan into
	var res Catagories

	//Query the accounts table
	var qry string = "select * from catagories order by cdescription"
	rows, err := DB.Query(qry)
	Check(err, " Product select list query failed. ")
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&res.Id, &res.Cdescription)
		Check(err, "Row scan failed")
		clses = append(clses, res)
	}

	//Construct the select statement
	var s string

	for _, res := range clses {
		if res.Id == sel {
			s += fmt.Sprintf("<option value=%q selected>%s</option>\n", res.Cdescription, res.Cdescription)
		} else {
			s += fmt.Sprintf("<option value=%q>%s</option>\n", res.Cdescription, res.Cdescription)
		}
	}

	return s
}
