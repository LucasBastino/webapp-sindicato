import * as v from "./validators.js";
import { validateForm, getInputValue } from "./common.js";

export function validatePaymentPlan() {
    const fields = {
        amount:                { id: "amount",                  validate: validatePaymentPlanAmount },
        numberOfInstallments:  { id: "number-of-installments",  validate: v.validateNumberOfInstallments },
        firstDueDate:          { id: "first-due-date",          validate: v.validateFirstDueDate },
        observations:          { id: "observations",            validate: v.validateObservations },
    };

    let ok = validateForm(fields);

    const checked = document.querySelectorAll(".payment-id-checkbox:checked");
    const paymentIdsError = document.querySelector(".payment-ids-error");
    if (paymentIdsError) paymentIdsError.innerHTML = "";
    if (!checked.length) {
        ok = false;
        if (paymentIdsError) {
            paymentIdsError.innerHTML = "Debés seleccionar al menos un pago vencido.";
        }
        const list = document.querySelector(".payment-plan-checklist");
        if (list) list.classList.add("input-invalid");
    }

    return ok;
}

function validatePaymentPlanAmount(input) {
    if (input === "") return { valid: false, error: "Campo requerido." };
    const result = v.validateAmount(input);
    if (!result.valid) return result;
    const n = Number(v.normalizeAmountInput(String(input)));
    if (!(n > 0)) return { valid: false, error: "El monto debe ser mayor a cero." };
    return { valid: true, error: "" };
}

function validatePaymentAmount(input) {
    if (input === "") return { valid: false, error: "Campo requerido." };
    const result = v.validateAmount(input);
    if (!result.valid) return result;
    const n = Number(v.normalizeAmountInput(String(input)));
    if (!(n > 0)) return { valid: false, error: "El monto debe ser mayor a cero." };
    return { valid: true, error: "" };
}

function validatePaymentPaidAt(input) {
    if (input === "") return { valid: false, error: "Ingresá la fecha de pago." };
    return v.validatePaidAt(input);
}

export function validatePayment() {
    const isPaid = (getInputValue("is-paid") || "").trim() === "true";

    const fields = {
        amount: { id: "amount", validate: validatePaymentAmount },
        observations: { id: "observations", validate: v.validateObservations },
    };

    if (isPaid) {
        fields.paidAt = { id: "paid-at", validate: validatePaymentPaidAt };
    }

    return validateForm(fields);
}

export function validateInstallment() {
    const isPaid = (getInputValue("is-paid") || "").trim() === "true";

    const fields = {
        observations: { id: "observations", validate: v.validateObservations },
    };

    if (isPaid) {
        fields.paidAt = { id: "paid-at", validate: validatePaymentPaidAt };
    }

    return validateForm(fields);
}
