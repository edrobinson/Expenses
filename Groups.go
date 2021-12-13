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

//Database: expenses   Table: groups
type Groups struct {
	Id           int    `json:"id"`
	Gdescription string `json:"gdescription"`
}

//Field list for scan statements
var GScanline string = "&GRes.Id, &GRes.Gdescription"

//Table var to scan into
var GRes Groups

var GPk string = "id"

func GroupsAjaxHandler(w http.ResponseWriter, r *http.Request) {
	//Get the body of the response - contains the table form values, record id and opcode
	//as json. data is read as byte type.
	data, err := ioutil.ReadAll(r.Body)

	WebCheck(err, "Error getting the data file from response body.", w)

	//Write it to a file
	var filename = "jsonfiles/groups" + ".json"

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
		ReadFirstGroups(w)
	case "readprev":
		ReadPrevGroups(w)
	case "read":
		ReadGroups(w)
	case "readnext":
		ReadNextGroups(w)
	case "readlast":
		ReadLastGroups(w)
	case "lookup":
		LookupGroups(w)
	case "insert":
		InsertGroups(w)
	case "update":
		UpdateGroups(w)
	case "delete":
		DeleteGroups(w)
	default:
		log.Println("Invalid opcode sent: ", OpCode)
		http.Error(w, "Invalid op code"+OpCode, http.StatusInternalServerError)
	}
}

//Page server func
func GroupsIndexHandler(w http.ResponseWriter, r *http.Request) {
	//Template data struct
	type tpldata struct {
		Title string
	}

	var data tpldata
	data.Title = "Groups Editor"
	//The account id anf product id are FKs of the table
	//so we get a set of html select options for both tables

	//Parse and execute the template file
	tmpl, err := template.ParseFiles("./templates/Groups.html",
		"./templates/header.tpl",
		"./templates/nav.tpl",
		"./templates/buttons.tpl",
		"./templates/footer.tpl")

	WebCheck(err, " Failed to parse products template. ", w)

	err = tmpl.Execute(w, data)
	WebCheck(err, " Failed to execute classes page template", w)
}

//Utility read function.
//Called by each crud function for validation
func GroupsLoadRecord() error {
	var whereclause string = CatPk + " = " + RecId
	var qry = "SELECT * FROM groups WHERE " + whereclause
	row := DB.QueryRow(qry)
	err := row.Scan(&CatRes.Id, &CatRes.Cdescription)
	return err
}

//CRUD Skeletons
//ReadProducts does a select where id = xx
func ReadGroups(w http.ResponseWriter) {
	row := ReadExact("groups", RecId)
	CompleteGroupsRead(w, row)
}

func ReadNextGroups(w http.ResponseWriter) {
	row := ReadNext("groups", RecId)
	CompleteGroupsRead(w, row)
}

func ReadPrevGroups(w http.ResponseWriter) {
	row := ReadPrev("groups", RecId)
	CompleteGroupsRead(w, row)
}

func ReadFirstGroups(w http.ResponseWriter) {
	row := ReadFirst("groups")
	CompleteGroupsRead(w, row)
}

func ReadLastGroups(w http.ResponseWriter) {
	row := ReadLast("groups")
	CompleteGroupsRead(w, row)
}

//Code common to all read ops after the utility db queries.
func CompleteGroupsRead(w http.ResponseWriter, row *sql.Row) {
	err := row.Scan(&GRes.Id, &GRes.Gdescription)

	if err == sql.ErrNoRows {
		http.Error(w, " Record not found.", http.StatusInternalServerError)
		return
	} else if err != nil {
		WebCheck(err, " Failed to read the requested record.", w)
	}

	//Marshal the db values to json
	var s []byte
	s, err = json.Marshal(GRes)
	WebCheck(err, "Marshal Error", w)

	//Send the json to the browser
	w.Write([]byte(s))
}

func LookupGroups(w http.ResponseWriter) {
    err := LookupService(w, "groups", "gdescription", true) 
	WebCheck(err, " Lookup service failed.",w)
}

/*	//Was a luby - "lookup by" - radio even given?
	var f Formval
	var fld string = ""

	//Scan for a "luby" radio button
	for _, f = range Fvals {
		if f.Name == "luby" {
			fld = f.Value //The field to query on
			break
		}
	}
	if fld == "" {
		http.Error(w, "No lookup field was selected.", http.StatusInternalServerError)
		return
	}

	//Scan for the field and get it's value
	var fldVal string
	for _, f = range Fvals {
		if f.Name == fld {
			fldVal = f.Value //The value to query on
			break
		}
	}
	if fldVal == "" {
		http.Error(w, "The lookup field was empty.", http.StatusInternalServerError)
		return
	}

	//Query for matching rows
	fldVal = "%" + fldVal + "%"
	var qry string = fmt.Sprintf("select * from groups where %s like %q", fld, fldVal)

	rows, err := ReadRows(qry)
	WebCheck(err, "Lookup - no matching records.", w)

	//Collect the rows data and build an option set.
	//Start with an empty option.
	var s string = fmt.Sprintf("<option value=%s></option>", 0)

	for rows.Next() {
		err := rows.Scan(&GRes.Id, &GRes.Gdescription)
		WebCheck(err, " Lookup rows scan failed.", w)
		s += fmt.Sprintf("<option value=%s>%s</option>", strconv.Itoa(GRes.Id), GRes.Gdescription)
	}

	var j []byte
	j, err = json.Marshal(s)
	WebCheck(err, " json marshal failed", w)
	w.Write([]byte(j))
}*/

func InsertGroups(w http.ResponseWriter) {
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
	var qry = "insert into " + "groups" + " " + fields + " values " + values
	log.Println("Groups insert: ", qry)
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

func UpdateGroups(w http.ResponseWriter) {
	//Get the record that the update call specifies
	//for comparison
	err := GroupsLoadRecord()

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
	var qry string = "update groups set " + fields + " where " + CatPk + " = " + RecId

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

func DeleteGroups(w http.ResponseWriter) {

	//Make certain the record exists
	err := GroupsLoadRecord()

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
func GroupsSelectList(sel int) string {
	var clses []Groups
	//Table var to scan into
	var res Groups

	//Query the accounts table
	var qry string = "select * from groups order by gdescription"
	rows, err := DB.Query(qry)
	Check(err, " Product select list query failed. ")
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&res.Id, &res.Gdescription)
		Check(err, "Row scan failed")
		clses = append(clses, res)
	}

	//Construct the select statement
	var s string

	for _, res := range clses {
		if res.Id == sel {
			s += fmt.Sprintf("<option value=%q selected>%s</option>\n", res.Gdescription, res.Gdescription)
		} else {
			s += fmt.Sprintf("<option value=%q>%s</option>\n", res.Gdescription, res.Gdescription)
		}
	}
	return s
}
