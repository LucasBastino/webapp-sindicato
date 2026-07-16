export function validateForm(fields, model, action) {
    let isValid = true;

    clearErrors();

    for (const key in fields){
        const { id, validate } = fields[key];

        const value = (getInputValue(id) || "").trim();
        const result = validate(value);

        if (!result.valid) {
            isValid = false;
            showError(id, result.error);
        }
    }

    if (isValid){
        postForm(model, action);
    }
}

function showError(id, error){
    const errorDiv = document.querySelector(`.${id}-error`);
    if (!errorDiv) return;
    errorDiv.style.display = 'inline'
    errorDiv.innerHTML = error
}

function clearErrors() {
    // obtiene todos los elementos que terminarn con "-error"
    const errors = document.querySelectorAll('[class$="-error"]');
    errors.forEach(e => {
        e.style.display = 'none';
        e.innerHTML = '';
    });
}

function postForm(model, action){
    document.getElementById(`submit-${action}-${model}-btn`).click()
}

