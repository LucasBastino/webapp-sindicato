import * as v from "./validators.js";
import { validateForm } from "./common.js";

export function validateParent() {
    const optional = (fn) => (input) => {
        if (!input || String(input).trim() === "") return { valid: true, error: "" };
        return fn(input);
    };

    const fields = {
        name:         { id: "name",         validate: v.validateNameOrLastName },
        lastName:     { id: "last-name",    validate: v.validateNameOrLastName },
        relationship: { id: "relationship", validate: optional(v.validateRelationship) },
        birthday:     { id: "birthday",     validate: optional(v.validateBirthday) },
        gender:       { id: "gender",       validate: optional(v.validateGender) },
        cuil:         { id: "cuil",         validate: optional(v.validateCuilCuit) },
        observations: { id: "observations", validate: v.validateObservations },
    };

    return validateForm(fields);
}
