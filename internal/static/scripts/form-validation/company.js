import * as v from "./validators.js";
import {validateForm} from "common.js"


function validateCompany(action){

    const fields = {
        name: { id: "name", validate: v.validateCompanyName },
        companyNumber: { id: "company-number", validate: v.validateCompanyNumber },
        address: { id: "address", validate: v.validateAddress },
        cuit: { id: "cuit", validate: v.validateCuilCuit },
        district: { id: "district", validate: v.validateDistrict },
        postalCode: { id: "postal-code", validate: v.validatePostalCode },
        phone: { id: "phone", validate: v.validatePhone },
        contact: { id: "contact", validate: v.validateContact },
        observations: { id: "observations", validate: v.validateObservations }
    };
    validateForm(fields, "company", action)
}