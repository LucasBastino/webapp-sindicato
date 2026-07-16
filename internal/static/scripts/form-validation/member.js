import * as v from "./validators.js";
import {validateForm} from "common.js"

function validateMember(action){
  const fields = {
        name: { id: "name", validate: validateNameOrLastName },
        lastName: { id: "last-name", validate: validateNameOrLastName },
        dni: { id: "dni", validate: validateDni },
        birthday: { id: "birthday", validate: validateBirthday },
        gender: { id: "gender", validate: validateGender },
        maritalStatus: { id: "marital-status", validate: validateMaritalStatus },
        phone: { id: "phone", validate: validatePhone },
        email: { id: "email", validate: validateEmail },
        address: { id: "address", validate: validateAddress },
        postalCode: { id: "postal-code", validate: validatePostalCode },
        district: { id: "district", validate: validateDistrict },
        memberNumber: { id: "member-number", validate: validateMemberNumber },
        // chequeo companyID porque el afiliado puede cambiar de empresa
        companyID: { id: "company-id", validate: validateCompanyID },
        cuil: { id: "cuil", validate: validateCuilCuit },
        category: { id: "category", validate: validateCategory },
        entryDate: { id: "entry-date", validate: validateEntryDate },
        observations: { id: "observations", validate: validateObservations },
    };

    validateForm(fields, "member", action)
}

