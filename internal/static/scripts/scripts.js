document.onkeydown = function(e){
    if (e.key === 'Enter' ){
        console.log("se apreto enter")
        e.preventDefault()
    }
}

function enableInputs(){
    inputs = Array.from(document.getElementsByTagName("input"))
    inputs.forEach(input => {
        input.disabled = false
    }); 
}

function enableSelects(){
    selects = Array.from(document.getElementsByTagName("select"))
    selects.forEach(select => {
        select.disabled = false
    }); 
     if (document.getElementById('name-company-input')){
        document.getElementById('name-company-input').disabled = true
     }
}

function blankSearchInput(){
    document.getElementById("search-input").value = ""
}

function showConfirmButton(model){
    document.getElementById(`btn-${model}-confirm`).style.display = 'inline'
}



function showAddMemberBtn(){
    document.getElementById('add-member-btn').style.display = 'inline'
}

function showAddCompanyBtn(){
    document.getElementById('add-company-btn').style.display = 'inline'
}

function hideAddMemberBtn(){
    document.getElementById('add-member-btn').style.display = 'none'
}

function hideAddCompanyBtn(){
    document.getElementById('add-company-btn').style.display = 'none'
}

function showMemberSearchInput(){
    document.getElementById('member-search-input-div').style.display = 'inline'
}

function hideMemberSearchInput(){
    input = document.getElementById('member-search-input-div').style.display = 'none'
}

function showCompaniesearchInput(){
    document.getElementById('company-search-input-div').style.display = 'inline'
}

function hideCompaniesearchInput(){
    document.getElementById('company-search-input-div').style.display = 'none'
}

function showCompanyMemberSearchNav(){
    document.getElementById('nav-company-member-search').style.display = 'inline'
}

function hideCompanyMemberSearchNav(){
    document.getElementById('nav-company-member-search').style.display = 'none'
}

function hideCompanyFile(){
    document.getElementById('company-btn-file').style.display = 'none'
}


function showParentSearchInput(){
    document.getElementById('parent-search-input').style.display = 'inline'
}

function hideParentSearchInput(){
    document.getElementById('parent-search-input').style.display = 'none'
}


function showParentFile(){
    document.getElementById('table-div').style.display = 'inline'
}

function hideParentFile(){
    document.getElementById('table-div').style.display = 'none'
}

function showAddParentButton(){
    document.getElementById('add-parent-button').style.display = 'inline'
}

function hideAddParentButton(){
    document.getElementById('add-parent-button').style.display = 'none'
}

function showParentTable(){
    document.getElementById('table-div').style.display = 'inline'
}

function loadParentTable(){
    document.getElementById('parent-button').click()
}

function hideParentTable(){
    document.getElementById('parent-table').style.display = 'none'
}

function showCompanyButton(){
    document.getElementById('company-button').style.display = 'inline'
}

function hideCompanyButton(){
    document.getElementById('company-button').style.display = 'none'
}

function hideAddPaymentButton(){
    document.getElementById('add-payment-btn').style.display = 'none'
}

function showAddPaymentBtn(){
    document.getElementById('add-payment-btn').style.display = 'inline'
}

function showXCompany(){
    document.getElementById('x-company').style.display = 'inline'
}


function searchCompanyAgain(){
    document.getElementById('x-company').style.display = 'none'
    document.getElementById('name-company-input').style.display = 'none'
    document.getElementById('company-search-box').style.display = 'inline'
}

function selectCompany(companyID, companyName){
    document.getElementById("company-id-input").value = companyID
    document.getElementById("name-company-input").value = companyName
    document.getElementById('company-search-box').style.display = 'none'
    document.getElementById('optionsTable').innerHTML = ''
    document.getElementById('name-company-input').style.display = 'inline'
    document.getElementById('x-company').style.display = 'inline'
}

function disableCompanyInput(){
    document.getElementById('name-company-input').disabled = true
}

function reloadWeb(){
    location.reload();
}

function updatePermissions() {
  const Admin = document.getElementById('admin')
  const canWrite = document.getElementById('can-write')
  const canDelete = document.getElementById('can-delete')

  if (Admin.checked) {
    canWrite.checked = true
    canDelete.checked = true
  }

  if (!canWrite.checked || !canDelete.checked) {
    Admin.checked = false
  }
}

function onAdminChange(checkbox) {
    const id = checkbox.dataset.id;
    if (checkbox.checked) {
        const canWrite = document.getElementById(`can-write-${id}`);
        const canDelete = document.getElementById(`can-delete-${id}`);
        canWrite.checked = true;
        canDelete.checked = true;
        htmx.trigger(canWrite, 'change');
        htmx.trigger(canDelete, 'change');
    }
}

function onPermissionChange(checkbox) {
    const id = checkbox.dataset.id;
    if (!checkbox.checked) {
        const Admin = document.getElementById(`admin-${id}`);
        Admin.checked = false;
        htmx.trigger(Admin, 'change');
    }
}

function closeModal() {
    document.getElementById('modal-container').innerHTML = '';
}

