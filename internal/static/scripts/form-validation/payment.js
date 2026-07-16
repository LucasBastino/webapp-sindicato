import * as v from "./validators.js";
import {validateForm} from "common.js"

function validatePayment(action){

    const fields = {
        // no chequeo companyID porque el payment no puede cambiar de company
        // month: { id: "month", validate: validateMonth },
        // year: { id: "year", validate: validateYear },
        isPaid: { id: "is-paid", validate: validateIsPaid },
        status: { id: "status", validate: validateStatus },
        amount: { id: "amount", validate: validateAmount },
        paidAt: { id: "paid-at", validate: validatePaidAt },
        observations: { id: "observations", validate: validateObservations }
    };

    validateForm(fields, "payment", action);
}



