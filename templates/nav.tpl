    <nav class="navbar navbar-expand-lg navbar-light bg-light">
       <ul class="navbar-nav">
        <li>
        <a class="navbar-brand" href="#" style="margin-right: 50px;">{{.Title}}</a>
        </li>
        <li class="nav-item dropdown" style="margin-right: 50px;">
        <button type="button" class="btn btn-primary btn-sm dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
        Page Menu
        </button>
            <div class="dropdown-menu">
              <a class="dropdown-item" href="/expenseindex">Expenses</a>
              <a class="dropdown-item" href="/catagoriesindex">Catagories</a>
              <a class="dropdown-item" href="/groupsindex">Groups</a>
              <a class="dropdown-item" href="/reportsindex">Expense Report</a>
            </div>
        </li>
        <li class="nav-item dropdown">
        <button type="button" class="btn btn-primary btn-sm dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false" style="margin-left: 100px">
        Help Menu
        </button>        
            <div class="dropdown-menu dropdown-menu-right">
              <a class="dropdown-item" href="/static/help/ExpensesIntroHelp.html">Introduction Help</a>
              <a class="dropdown-item" href="/static/help/ExpensesEditorHelp.html">Expenses Help</a>
              <a class="dropdown-item" href="/static/help/CataGoryEditorHelp.html">Catagories Help</a>
              <a class="dropdown-item" href="/static/help/GroupEditorHelp.html">Groups Help</a>
              <a class="dropdown-item" href="/static/help/Reports.html">Report Help</a>
            </div>
        </li>    
       </ul>
    </nav>

