import * as v from "./helpers.js"
import { getInputValue } from "./common.js"


// PERSONAL INFO
export function validateNameOrLastName(input){
    // no pongo trim porque borro los espacios despues

    // borro los espacios y despues verifico que no este vacio
    input = v.removeChars(input, " '")

    if (input === "") return { valid: false, error: "Campo requerido." }

    let result = v.isAlphabetic(input)
    if (!result.valid) return result

    return v.isNotLongerThan(input, 50)
}

export function validateDni(input){
    
    if (input === "") return { valid: true, error: "" }

    // verifico que sean solo numeros    
    let result = v.isNumeric(input)
    if (!result.valid) return result

    return v.isNotLongerThan(input, 8)
}

export function validateBirthday(input){
    if (input === "") return { valid: false, error: "Campo requerido." }

    const min = new Date(1900, 0, 1)
    const max = new Date()

    return v.isValidDate(input, min, max)
}

export function validateGender(input){

    if (input === "") return { valid: false, error: "Campo requerido." }

    const optionsSet = new Set([
        "Masculino",
        "Femenino",
        "No binario",
        "Masculino transgénero",
        "Femenino transgénero",
        "Otra identidad",
        "Sin especificar"
    ]);

    return v.isValidOption(input, optionsSet)
}

export function validateRelationship(input){

    if (input === "") return { valid: false, error: "Campo requerido." }

    let result = v.isAlphabetic(input)
    if (!result.valid) return result

    return v.isNotLongerThan(input, 20)
}

// export function validateIsPaid(input){
//     if (input === "") return { valid: false, error: "Campo requerido." }

//     if (input !== "true" && input !== "false"){
//         return { valid: false, error: "El valor ingresado no es válido." }
//     }
//     return { valid: true, error: ""}
// }


export function validateMaritalStatus(input){

    if (input === "") return { valid: false, error: "Campo requerido." }

    // creo el set, un array de strings, solo keys
    const optionsSet = new Set([
        "Soltero",
        "Casado",
        "Separado",
        "Viudo",
        "Separado",
        "Divorciado"])

    return v.isValidOption(input, optionsSet)
}


// ADDRESS INFO
export function validateAddress(input){
    if (input === "") return { valid: true, error: "" }

    input = v.removeChars(input, " #-'&~,.º")

    let result = v.isAlphanumeric(input)
    if (!result.valid) return result

    return v.isNotLongerThan(input, 100)
}

export function validateDistrict(input){
    // borro los espacios y caracteres especiales permitidos
    // y recien despues me fijo si esta vacio
    if (input === "") return { valid: true, error: "" }
    
    input = v.removeChars(input, " #-'&~,.º")

    let result = v.isAlphanumeric(input)
    if (!result.valid) return result

    return v.isNotLongerThan(input, 50)
}


export function validatePostalCode(input){
    
    if (input === "") return { valid: true, error: "" }
    
    if (input.length === 4){
        let result = v.isNumeric(input)
        if (!result.valid) return result
        return { valid: true, error: "" }
    }
    if (input.length === 8){
        let result = v.isAlphabetic(input[0])
        if (!result.valid) return result

        for (let i = 1; i<= 4; i++){
            let result = v.isNumeric(input[i])
            if (!result.valid) return result
        }

        for (let i = 5; i<= 7; i++){
            let result = v.isAlphabetic(input[i])
            if (!result.valid) return result
        }

        return { valid: true, error: "" }
    }

    return { valid: false, error: "Debe contener 4 u 8 caracteres."}
}



// CONTACT INFO
export function validatePhone(input){
    if (input === "") return { valid: true, error: "" }

    input = v.removeChars(input, " ")

    let result = v.isNumeric(input)
    if (!result.valid) return result

    return v.isNotLongerThan(input, 20)
}

export function validateEmail(input){
    
    if (input === "") return { valid: true, error: "" }

    const parts = input.split("@")
    if (parts.length != 2) return { valid: false, error: "Formato no válido." }
    let [local, domain] = parts
    if (
        local.length === 0 ||
        !domain.includes(".") ||
        domain.length < 3 ||
        domain[0] === "." ||
        domain[domain.length - 1] === "."
    ){
        return { valid: false, error: "Formato no válido." }
    }

    const domainParts = domain.split(".")
    if (domainParts.includes("")){
        return { valid: false, error: "Formato no válido." } // cubre "".gmail.com", "gmail."", "gmail..com"
    }
    if (
        domainParts.length < 2 ||
        // si la ultima parte (el TLD) no tiene minimo 2 caracteres no es válido
        domainParts[domainParts.length -1].length < 2
    ){
        return { valid: false, error: "Formato no válido." }  
        }

    domain = v.removeChars(domain, ".")
    let result = v.isAlphanumeric(domain)
    if (!result.valid) return result

    input = v.removeChars(input, "!#$%&'*+-/=?^_`{|}~.@")
    result = v.isAlphanumeric(input)
    if (!result.valid) return result

    return v.isNotLongerThan(input, 50)
}

// SOCIAL SECURITY INFO
export function validateMemberNumber(input){
    
    if (input === "") return { valid: true, error: "" }

    let result = v.isNumeric(input)
    if (!result.valid) return result

    return v.isNotLongerThan(input, 20)
}


export function validateCompanyID(input){
    if (input === "") return { valid: false, error: "Campo requerido." }

    let result = v.isNumeric(input)
    if (!result.valid) return result

    if (input == "0") return { valid: false, error: "El valor ingresado no es válido." }

    return v.isNotLongerThan(input, 10)
}


export function validateCategory(input){

    if (input === "") return { valid: false, error: "Campo requerido." }

    const optionsSet = new Set([
        "Nivel 1: Oficial Múltiple",
        "Nivel 2: Oficial Especializado",
        "Nivel 3: Oficial General",
        "Nivel 4: Medio Oficial",
        "Nivel 5: Ayudante",
        "Nivel 6: Operario Act. Industrial"])
    return v.isValidOption(input, optionsSet)
}

export function validateEntryDate(input){

    if (input === "") return { valid: true, error: "" }

    const min = new Date(1900, 0, 1)
    const max = new Date()

    return v.isValidDate(input, min, max)
}



// v.Company INFO
export function validateCompanyName(input){
    if (input === "") return { valid: false, error: "Campo requerido." }
    return v.isNotLongerThan(input, 150)
}


export function validateContact(input){

    if (input === "") return { valid: true, error: "" }
    
    return v.isNotLongerThan(input, 200)
}

export function validateCompanyNumber(input){

    if (input === "") return { valid: true, error: "" }
    
    let result = v.isNumeric(input)
    if (!result.valid) return result
    
    return v.isNotLongerThan(input, 20)
}


export function validateCuilCuit(input){
    
    if (input === "") return { valid: true, error: "" }

    input = v.removeChars(input, " ")
    // verifico si tiene guiones o no
    if (input.includes("-")){
        // si tiene lo separo en partes para chequear individualmente
        const parts = input.split("-")
        if (parts.length !== 3){
            return { valid: false, error: "Formato no válido."}
        }
        const [prefix, dni, digit] = parts
        if (prefix.length !== 2 || dni.length !== 8 || digit.length !== 1){
                return { valid: false, error: "Formato no válido."}
        } else{
            return { valid: true, error: "" }
        }
    }

    let result = v.isNumeric(input)
    if (!result.valid) return result

    return v.hasXLength(input, 11)
}



// PAYMENT INFO
// export function validateMonth(input){

//     if (input === "") return { valid: false, error: "Campo requerido." }

//     let result = v.isNumeric(input)
//     if (!result.valid) return result

//     if (input < 1 || input > 12) return { valid: false, error: "El valor ingresado no es válido" }
//     return { valid: true, error: ""}
// }


// export function validateYear(input){
    
//     // verifico que no este vacio
//     if (input === "") return { valid: false, error: "Campo requerido." }
    
//     let result = v.isNumeric(input)
//     if (!result.valid) return result

//     // si es mayor al año actual o menor a 2000 es un año invalido
//     if ((input > new Date().getFullYear()) || (input < 2000)) {
//         return { valid: false, error: "El valor ingresado no es válido"}
//     }
//     return { valid: true, error: ""}
// }


export function validateStatus(input){
   if (input === "") return { valid: false, error: "Campo requerido." }

    optionsSet = new Set([
        "Completado",
        "Pendiente",
        "Vencido",
        "En plan de pago"])

    return v.isValidOption(input, optionsSet)
}

export function validateAmount(input){
    if (input === "") return { valid: true, error: "" }

    const normalized = normalizeAmountInput(String(input))

    if (normalized.includes(".")) {
        if (normalized.startsWith(".") || normalized.endsWith(".")) {
            return { valid: false, error: "Formato no válido." }
        }

        const parts = normalized.split(".")
        if (parts.length !== 2) {
            return { valid: false, error: "Formato no válido." }
        }

        if (parts[1].length > 2) {
            return { valid: false, error: "Formato no válido." }
        }
    }

    const digits = v.removeChars(normalized, ".")

    const result = v.isNumeric(digits)
    if (!result.valid) return result

    return v.isNotLongerThan(digits, 20)
}

// With comma: Argentine style (dots = thousands). Without: dot is decimal.
function normalizeAmountInput(input) {
    const trimmed = input.trim()
    if (trimmed.includes(",")) {
        return trimmed.replaceAll(".", "").replace(",", ".")
    }
    return trimmed
}

export { normalizeAmountInput }

export function validatePaidAt(input){
    if (input === "") return { valid: true, error: "" }

    const min = new Date(1900, 0, 1)
    const max = new Date()
    
    return v.isValidDate(input, min, max)
}

export function validateNumberOfInstallments(input){

    if (input === "") return { valid: false, error: "Campo requerido." }
    
    let result = v.isNumeric(input)
    if (!result.valid) return result
    
    return v.isNotLongerThan(input, 3)
}

export function validateFirstDueDate(input){
    if (input === "") return { valid: false, error: "Campo requerido." }

    const min = new Date(1900, 0, 1)
    const max = new Date()
    max.setFullYear(max.getFullYear() +10)

    return v.isValidDate(input, min, max)
}


export function validateObservations(input){
    return v.isNotLongerThan(input, 2000)
}

export function validateUsername(input){
    if (input === "") return { valid: false, error: "Campo requerido." }

	let result = v.isAlphanumeric(input)
    if (!result.valid) return result

    result = v.hasAtLeast(input, 3)
    if (!result.valid) return result

	return v.isNotLongerThan(input, 20)
}

export function validatePassword(input){
    if (input === "") return { valid: false, error: "Campo requerido." }

    let result = v.isNotLongerThan(input, 64)
    if (!result.valid) return result

    result = v.hasAtLeast(input, 8)
    if (!result.valid) return result

    input = v.removeChars(input, " @#-'&,.!?*+")

    return v.isAlphanumeric(input)
}

export function validateConfirmPassword(input){
    if (input === "") return { valid: false, error: "Campo requerido." }

    const password = getInputValue("password").trim()

    if (input !== password) {
        return { valid: false, error: "Las contraseñas no coinciden." }
    }

    return { valid: true, error: "" }
}


export function validatePermissions(input){
    return v.isValidOption(input, permissionOptionsSet)
}

const permissionOptionsSet = new Set([
        "editor",
        "viewer",
        "",
    ]);