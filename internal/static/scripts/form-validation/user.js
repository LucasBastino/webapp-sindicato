import * as v from "./validators.js";
import {validateForm} from "common.js"

function validateUser(action){

    const fields = {
        username: { id: "username", validate: validateUsername },
        password: { id: "password", validate: validatePassword },
        admin: { id: "admin", validate: validateAdmin },
        member: { id: "member", validate: validatePermissions },
        company: { id: "company", validate: validatePermissions },
        payment: { id: "payment", validate: validatePermissions },
    };

    validateForm(fields, "user", action);
}



