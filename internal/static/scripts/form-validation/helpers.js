export function removeChars(input, charsToRemove){
    const removeSet = new Set(charsToRemove);

    let finalString = "";

    for (const r of input){
        if (!removeSet.has(r)) {
            finalString += r;
        }
    }

    return finalString;
}

export function hasXLength(input, length){
    if (input.length !== length){
        return { valid: false, error: `Debe contener ${length} caracteres.` }
    }
    return { valid: true, error: "" }
}

export function isNotLongerThan(input, limit){
    if (input.length>limit){
        return { valid: false, error: `No puede contener más de ${limit} caracteres.` }
    }
    return { valid: true, error: ""}
}

export function hasAtLeast(input, limit){
    if (input.length<limit){
        return { valid: false, error: `Debe contener al menos ${limit} caracteres.` }
    }
    return { valid: true, error: ""}
}

export function isNumeric(input){
    let validChars = "0123456789"
    return isValidCharacter(input, validChars)
}

export function isAlphabetic(input){
    let validChars = "abcdefghijklmnñopqrstuvwxyzáéíóúüÃ"
    return isValidCharacter(input, validChars)
}

export function isAlphanumeric(input){
    let validChars = "abcdefghijklmnñopqrstuvwxyzáéíóúüÃ0123456789"
    return isValidCharacter(input, validChars)
}

// export function isValidCharacter(input, validChars){
//     input = input.toLowerCase()
//     // el for tipo range recorre por caracter completo de unicode
//     // el for incremental recorre por unidad de codigo utf-16
//     // y con algunos caracteres como los emojis o simbolos raros te los puede romper porque ocupan 2 unidades de codigo utf-16 
//     for (const r of input){
//         if (!validChars.includes(r)){
//             return { valid: false, error: "El valor ingresado contiene un caracter no válido." }
//         }
//     }
//     return { valid: true, error: ""}
// }

function isValidCharacter(input, validChars){
    input = input.toLowerCase();
    const validSet = new Set(validChars);

    // corta al primer falso
    const isValid = [...input].every(r => validSet.has(r));

    if (isValid){
        return { valid: true, error: "" }
    }
    return { valid: false, error: "Caracter no válido."};
}

export function isValidOption(input, optionsSet){
    // si existe la key dentro del set, todo ok
    if (optionsSet.has(input)){
        return { valid: true, error: ""}
    }
    return { valid: false, error: "Opción no válida."}
}


export function isValidDate(input, min, max){
    input = input.trim()


    // separo la fecha por el caracter "/" y verifico que haya 3 partes
    var arrayStr = input.split("/")
    if (arrayStr.length !== 3){
        return { valid: false, error: "Fecha no válida." }
    }

    // convierto el array de strings en array de numeros, esto hace que "01" o "1" de como resultado "1"
    var arrayInt = arrayStr.map(Number)

    // creo variables obteniendo su valor del array de numeros
    const [day, month, year] = arrayInt

    // chequeo que no sean falsy (si un numero es 0 js lo toma como falsy)
    if (!day || !month || !year) {
        return { valid: false, error: "Fecha no válida." };
    }

    // creo la fecha (el mes se toma por index, enero es 0)
    // tener en cuenta que js aca normaliza la fecha, si pusiste 30 de febrero, se suman 2 dias y lo cuenta como 2 de marzo
    // por lo tanto despues al comparar el dia mes y año del array con el de date, van a dar numeros distintos y la validacion falla
    const date = new Date(year, month -1, day)

    // getDate devuelve el numero del dia del mes, getDay devuelve el numero de dia de la semana 
    // compara que las fechas sean las mismas despues de normalizar
    if (day === date.getDate() &&
        month === date.getMonth() +1 &&
        year === date.getFullYear() &&
        date >= min &&
        date <= max)
        {
            return { valid: true, error: "" }
        }
        return { valid: false, error: "Fecha no válida." }

}

export function isBoolean(input){
    if (input === "true" || input === "false"){
        return { valid: true, error: "" }
    }
    return { valid: false, error: "El valor ingresado no es válido" }
}