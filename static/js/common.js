/*
    Functions common to all of the web pages
    
    Using JQuery where possible
*/

    /*
        Send the Ajax request to the page server
        and handle the success or error response.
        The page specific js file calls 
        the $.ajaxSetup() to set the url
    */
    function sendRequest(formdata, op){
        $.ajax({
            type: "POST",
            contentType: 'application/json; charset=utf-8',
            //Sending the form's data
            data: formdata,
            //Expecting a json return
            dataType: 'json',
            //Success calls the response handler 
            //passing the response and operation,
            success: function(data, status, xhttp){
                handleResponse(data, op)
            },
            //Error return just alerts the response text
            error: function(r){
                doError(op, r)
            }
        });
    } 

    /*
        Handle request sets up operation code 
        and the record id field in the form hidden
        fields, serializes the form data
        and calls the Ajax handler.
    */
     function handleRequest(opcode){
         event.preventDefault();


        //Save the requested op code and record id
        //in the hidden inputs of the form.
        $('#crudop').val(opcode)        
        $('#recid').val($('#id').val())
        
        //Extract the form's values into a string
        var data = $("#form1").serializeArray()
        
        //To JSON string
        data = JSON.stringify(data)
        
        //and send it
        sendRequest(data,opcode )
    }

    /*
        Ajax success response handler:
        Stringify the response and call the
        appropriate handler per the opcode
        passing the response data.
   */
    function handleResponse(r, op){
            
            //Stringify the response
            r = JSON.stringify(r)
            
            switch (op){
                case 'readfirst':
                case 'readlast':
                case 'read':
                case "readnext":
                case "readprev":
                    doReadResponse(r);
                    break;
                 case 'insert':
                    doInsertResponse(r);
                    break;
                case 'update':
                    doUpdateResponse(r);
                    break;
                case 'lookup':
                    doLookupResponse(r);
                    break;
                case 'delete':
                    doDeleteResponse(r);
                    break;
                default:
                    doError(op)
            }
    }
    
    //Read response populates the form
    //for read, readnext and readprev
    function doReadResponse(r){
        //Parse the json to {k:v} pairs
        k = JSON.parse(r)
        
        //Iterate over the string/obj and fill in the form fields
        Object.keys(k).forEach(function(key) {
            $("#" + key).val(k[key])
        })
        
        //Iterate over all of the checkboxes and set them 
        //according to the values returned.
        $('input[type=checkbox]').each(function () {
            var id = $(this).attr('id')
            if ($('#' + id).val() == '1'){
                $('#' + id).prop('checked', true)
            }else{
                $('#' + id).prop('checked', false)
            }
        })
    }
 
   function doInsertResponse(r){
        alert("Record Created.")
    }
    
    function doUpdateResponse(r){
        alert("Record Updated")
    }
    
    function doDeleteResponse(r){
        alert("Record Deleted" + r)
    }
    
    function doLookupResponse(r){
        //Uncheck the radio button
        $("input[type='radio'].lb").each(function () {
            $(this).prop('checked', false);
        });
        
        var k = JSON.parse(r)
        
        //Insert the options into the select tag
        document.getElementById("slookup").innerHTML = k
        //See if anything came back
        var Len = document.getElementById("slookup").options.length
        if(Len <2){
            alert("No records matched.")
            return
        }
        toggleElement("lookuplist")
    }
    
    function lookupProc(sel){
        //Hide the select div
        toggleElement("lookuplist")
        //Get the value of the select
        var v = sel.value
        //User choose the cancel option?
        if(v == 'cancel') return
        $("#id").val(v)
        handleRequest("read")
        
    }    
  
    
    
    //Just alert the op and error message
    //when an error response is received
    function doError(op, r){
        cap = capitalizeFirstLetter(op)
        v = JSON.stringify(r)
        alert(cap + ": " + r.responseText)

    }

   //Capitalize op codes for display
    function capitalizeFirstLetter(string) {
        return string.charAt(0).toUpperCase() + string.slice(1);
}

//Simple element toggle show/hide at each call
function toggleElement(id){ 
    $("#"+id).toggle()
}


//Trickery to do a read on a table from the tablelist screen
function ReadFromSelect(tblname){
    //Hide the tablelist record list select
    toggleElement("recordlist")
 
    
    //Setup the ajax url for the table
    //This is the trick...
    $.ajaxSetup({
      url: "/" + tblname
    });  

    $('#crudop').val("read")        
    $('#recid').val($('#recordlist').val())
    
    //Extract the form's values into a string
    var data = $("#form1").serializeArray()
    
    //To JSON string
    data = JSON.stringify(data)
    
    //and send it
    sendRequest(data,"read" )
}    

//Call the current page's help page
function callHelp(){
    window.location.href = helpurl
}