import * as v from "./validators.js";
import {validateForm} from "common.js"

function validateParent(action){

    const fields = {
        // no chequeo memberID porque el parent no puede cambiar de member
        name: { id: "name", validate: validateNameOrLastName },
        lastName: { id: "last-name", validate: validateNameOrLastName },
        relationship: { id: "relationship", validate: validateRelationship },
        birthday: { id: "birthday", validate: validateBirthday },
        gender: { id: "gender", validate: validateGender },
        cuil: { id: "cuil", validate: validateCuilCuit }
    };

    validateForm(fields, "parent", action)
}

