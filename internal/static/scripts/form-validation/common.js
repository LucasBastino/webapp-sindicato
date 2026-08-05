export function validateForm(fields) {
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

    return isValid;
}

export function getInputValue(id) {
    const el = document.getElementById(id);
    return el ? (el.value || "") : "";
}

export function showError(id, error) {
    const input = document.getElementById(id);
    const formGroup = input ? input.closest(".form-group") : null;
    const errorEl = (formGroup && formGroup.querySelector(".field-error"))
        || document.querySelector(`.${id}-error`);
    if (errorEl) errorEl.innerHTML = error;
    if (input) {
        if (formGroup) formGroup.classList.add("input-invalid");
        else input.classList.add("input-invalid");
    }
}

function clearErrors() {
    document.querySelectorAll('.field-error').forEach(e => { e.innerHTML = ''; });
    document.querySelectorAll('.input-invalid').forEach(el => el.classList.remove('input-invalid'));
}
