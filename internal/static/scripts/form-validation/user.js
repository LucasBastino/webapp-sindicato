import * as v from "./validators.js";
import { validateForm, getInputValue, showError } from "./common.js";

export function validateUser() {
    const fields = {
        username: { id: "username", validate: v.validateUsername },
        password: { id: "password", validate: v.validatePassword },
        confirmPassword: { id: "confirm-password", validate: v.validateConfirmPassword },
        member: { id: "member", validate: v.validatePermissions },
        company: { id: "company", validate: v.validatePermissions },
    };

    return validateForm(fields);
}

function validateRequired(input) {
    if (input === "") return { valid: false, error: "Campo requerido." };
    return { valid: true, error: "" };
}

export function validateChangePassword() {
    const fields = {};

    if (document.getElementById("current-password")) {
        fields.currentPassword = { id: "current-password", validate: validateRequired };
    }

    fields.password = { id: "password", validate: v.validatePassword };
    fields.confirmPassword = { id: "confirm-password", validate: v.validateConfirmPassword };

    let ok = validateForm(fields);

    const password = (getInputValue("password") || "").trim();
    const confirm = (getInputValue("confirm-password") || "").trim();
    const current = (getInputValue("current-password") || "").trim();

    if (ok && password !== confirm) {
        showError("confirm-password", "Las contraseñas no coinciden.");
        showError("password", "Las contraseñas no coinciden.");
        ok = false;
    }

    return ok;
}
