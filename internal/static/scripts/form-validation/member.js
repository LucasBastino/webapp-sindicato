import * as v from "./validators.js";
import { validateForm } from "./common.js";

export function validateMember() {
    const fields = {
        name:          { id: "name",           validate: v.validateNameOrLastName },
        lastName:      { id: "last-name",      validate: v.validateNameOrLastName },
        dni:           { id: "dni",            validate: v.validateDni },
        birthday:      { id: "birthday",       validate: v.validateBirthday },
        gender:        { id: "gender",         validate: v.validateGender },
        maritalStatus: { id: "marital-status", validate: v.validateMaritalStatus },
        phone:         { id: "phone",          validate: v.validatePhone },
        email:         { id: "email",          validate: v.validateEmail },
        address:       { id: "address",        validate: v.validateAddress },
        postalCode:    { id: "postal-code",    validate: v.validatePostalCode },
        district:      { id: "district",       validate: v.validateDistrict },
        memberNumber:  { id: "member-number",  validate: v.validateMemberNumber },
        companyID:     { id: "company-id",     validate: v.validateCompanyID },
        cuil:          { id: "cuil",           validate: v.validateCuilCuit },
        category:      { id: "category",       validate: v.validateCategory },
        entryDate:     { id: "entry-date",     validate: v.validateEntryDate },
        observations:  { id: "observations",   validate: v.validateObservations },
    };

    return validateForm(fields);
}
