import * as v from "./validators.js";
import { validateForm } from "./common.js";

export function validateParent() {
    const fields = {
        name:         { id: "name",         validate: v.validateNameOrLastName },
        lastName:     { id: "last-name",    validate: v.validateNameOrLastName },
        relationship: { id: "relationship", validate: v.validateRelationship },
        birthday:     { id: "birthday",     validate: v.validateBirthday },
        gender:       { id: "gender",       validate: v.validateGender },
        cuil:         { id: "cuil",         validate: v.validateCuilCuit },
        observations: { id: "observations", validate: v.validateObservations },
    };

    return validateForm(fields);
}
