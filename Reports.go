/*
   This is the expenses report generation module.

   It will do a simple query on the expense table
   and produce a pdf file.

   The report is ordered by catagory then date so that the report
   is grouped by the catagory and by date within the catagory.

   The user can select a start date and an end date - both, either or none.

   Using the gofpdf by jung. Found out that the pdf must be re-instanced at
   each call because, unlike PDF and such, the app does not start over at
   every invocation.

   Note: Found an aparrent bug in gofpdf. If a new line is specified in a call
         to CellFormat() the app panics.
         The one thAT WORKS:
            pdf.CellFormat(2.0, 0.3, "Date", "", 0, "L", false, 0, "")

         The one that crashes:
            pdf.CellFormat(2.0, 0.3, "Date", "", 0, "L", false, 1, "")
                                                                ^
         The 1 is supposed to result moving to the next line after output.
         The one with the 0 is supposed to go nowhere and must have a Ln()
         after it to go to next line.
         
    The reports page does not use AJAX like the others. I just submits a form
    and extracts the form values.

*/

package main

import (
	"bytes"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
	"time"
)

//Instance the pdf generator
//Portrait, inches, letter paper size, no font directory
//Do it here so we can compile ok...
//It is reinstanced at the beginning of the generation.
var pdf = gofpdf.New("P", "in", "letter", "")

var reportpath string = "reports/expenses.pdf"

//Serve the reports page
func ReportsIndexHandler(w http.ResponseWriter, r *http.Request) {
	//Template data struct
	type tpldata struct {
		Title string
	}

	var data tpldata
	data.Title = "Expenses Report"

	//Parse and execute the template file
	tmpl, err := template.ParseFiles("./templates/Reports.html",
		"./templates/header.tpl",
		"./templates/nav.tpl",
		"./templates/footer.tpl")
	WebCheck(err, " Failed to parse reports template. ", w)

	err = tmpl.Execute(w, data)
	WebCheck(err, " Failed to execute reports template", w)
}

func ReportsHandler(w http.ResponseWriter, r *http.Request) {
	GenerateExpenseReport(w, r)
}

//Generate the PDF called from main.go
func GenerateExpenseReport(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//Extract the date range values
	var startdate string = r.PostFormValue("stdate")
	var enddate string = r.PostFormValue("nddate")

	//Load the expense data
	err, data := loadData(startdate, enddate)
	Check(err, "Report data load failed.")

	//Generate the pdf report
	GeneratePdf(data)

	//Send the pdf in the browser
	ShowPDF(w, r, reportpath)

}

/**************************** Fetch the Expense Data **************************/

//LoadData creates the query using the dates and loads the data.
//Returnsthe error and rows slice
func loadData(startdate, enddate string) (error, []Expense) {

	var qry string = "select * from expense " //order  by excatagory, exdate"

	//Date range setup

	//Quote the dates if given
	if startdate != "" {
		startdate = "'" + startdate + "'"
	}
	if enddate != "" {
		enddate = "'" + enddate + "'"
	}

	//Both given?
	if startdate != "" && enddate != "" {
		qry += "where exdate between " + startdate + " and " + enddate
	} else if startdate != "" {
		qry += " where exdate >= " + startdate
	} else if enddate != "" {
		qry += " where exdate <= " + enddate
	}
	qry += " order  by excatagory, exdate"

	rows, err := DB.Query(qry)
	Check(err, "Report data load failed.")

	var data []Expense

	for rows.Next() {
		err = rows.Scan(&ExRes.Id, &ExRes.Exdate, &ExRes.Examt, &ExRes.Excatagory, &ExRes.Exdescription)
		Check(err, " Load Data scan failed.")
		data = append(data, ExRes)
	}

	return nil, data
}

/************************* Begin PDF Generation ********************************/

//Add the catagory summary
func catagorySummary(catsum float32, catnum int) {
	var s string = fmt.Sprintf("Catagory Totals: %d transactions totaling $%6.2f", catnum, catsum)
	pdf.CellFormat(0, 0.3, s, "", 0, "L", false, 0, "")
	pdf.Ln(0.2)
	pdf.CellFormat(0, 0.3, "___________________________", "", 0, "L", false, 0, "")
	pdf.Ln(0.2)
}

func reportSummary(gtotal float32, gttlrecs int) {
	var s string = fmt.Sprintf("Report Grand Totals: %d transactions totaling $%6.2f", gttlrecs, gtotal)
	pdf.CellFormat(0, 0.3, "___________________________", "", 0, "L", false, 0, "")
	pdf.Ln(0.2)
	pdf.CellFormat(0, 0.3, s, "", 0, "L", false, 0, "")
	pdf.Ln(0.2)
}

//Label for a new catagory
func catagoryHeading(s string) {
	var c string = fmt.Sprintf("Catagagory: %s", s)
	pdf.CellFormat(0, 0.3, c, "", 0, "L", false, 0, "")
	pdf.Ln(0.2)
	//Add the column headers
	pdf.CellFormat(2.0, 0.3, "Date", "", 0, "L", false, 0, "")
	pdf.CellFormat(2.0, 0.3, "Amount", "", 0, "L", false, 0, "")
	pdf.CellFormat(3.0, 0.3, "Description", "", 0, "L", false, 0, "")
	pdf.Ln(0.2)
}

//Generate the page head and footer functions for the pdf
//The package leaves the function blank so it can be overridden and customized.

func setupPDFHeader() {
	pdf.SetHeaderFunc(func() {
		pdf.SetY(.2)
		pdf.SetFont("Arial", "B", 15)
		//pdf.Cell(2.2, 0, "")
		pdf.CellFormat(0, .4, "Expense Report", "", 0, "C", false, 0, "")
		pdf.Ln(.2)
		dt := time.Now()
		dts := dt.Format("01-02-2006 15:04")
		pdf.SetFont("Arial", "", 12)
		pdf.CellFormat(0, .4, dts, "", 0, "C", false, 0, "")
		pdf.Ln(.5)
	})
}

func setupPDFFooter() {
	pdf.SetFooterFunc(func() {
		pdf.SetY(-.5)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, .4, fmt.Sprintf("Page %d /{nb}", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})
}

//Create and output the PDF.
//The data parameter is a slice of the expense records
//created by the loader.
func GeneratePdf(data []Expense) {
	//Must declare a new instance so the pdf so it
	//does not retain context between calls.
	pdf = gofpdf.New("P", "in", "letter", "")

	//Setup our page head and footer
	setupPDFHeader()
	setupPDFFooter()

	pdf.AliasNbPages("")         //Set page of pages
	pdf.AddPage()                //Put in the first page
	pdf.SetFont("Arial", "", 10) //Set the document font

	//Catagory totals and grand totals
	var thiscat string = "" //Current catagory for group control
	var catsum float32 = 0  //Catagory total amount
	var catnum int = 0      //Number of recs in this catagory
	var gtotal float32 = 0  //Grand total amount
	var gttlrecs int = 0    //Grand total number of records

	//Range over the DB records adding to
	//the PDF body.
	for _, d := range data {

		//Manage Catagory changes.
		if thiscat != d.Excatagory {
			//If the current catagory is blank (first record) then just
			//Save the new catagory and output the catagory heading
			if thiscat == "" {
				thiscat = d.Excatagory
				catagoryHeading(d.Excatagory)

				//For subsequent headers save the new catagory, output the catagory summary
				//clear the cat sums and output the new cat heading
			} else {
				//Update the current catagory name
				thiscat = d.Excatagory

				//Output the catagory summary
				catagorySummary(catsum, catnum)

				//Init. the catagory totals
				catsum = 0
				catnum = 0

				//Output the new catagory heading
				catagoryHeading(d.Excatagory)
			}
		}

		//Maintain the area totals
		gttlrecs++
		gtotal += d.Examt
		catnum++
		catsum += d.Examt

		//Amount to string xxxxxx.yy
		samt := fmt.Sprintf("$%6.2f", d.Examt)

		//Copy the YYYY-MM-DD from the go date/time
		sdate := d.Exdate[0:10]

		//Add the transaction to the pdf
		pdf.CellFormat(2.0, 0.3, sdate, "", 0, "L", false, 0, "")
		pdf.CellFormat(2.0, 0.3, samt, "", 0, "L", false, 0, "")
		pdf.CellFormat(3.0, 0.3, d.Exdescription, "", 0, "L", false, 0, "")
		pdf.Ln(0.2)
	}

	//Output the grand totals
	reportSummary(gtotal, gttlrecs)

	//Store the pdf file and cleanup
	err := pdf.OutputFileAndClose("reports/expenses.pdf")
	Check(err, " PDF write and close failed")

}

//Render the pdf to the browser.
//Called from GenerateExpeseReport
func ShowPDF(w http.ResponseWriter, r *http.Request, filename string) {
	//Load the PDF file
	streamPDFbytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//To a byte buffer
	b := bytes.NewBuffer(streamPDFbytes)

	//Let 'em know what's coming
	w.Header().Set("Content-type", "application/pdf")

	//Write the file bytes to the brower
	if _, err := b.WriteTo(w); err != nil {
		fmt.Fprintf(w, "%s", err)
	}
}
