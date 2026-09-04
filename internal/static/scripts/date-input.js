(function () {
  function onlyDigits(text) {
    var digits = "";
    for (var i = 0; i < text.length; i++) {
      var ch = text.charAt(i);
      if (ch >= "0" && ch <= "9") {
        digits += ch;
      }
    }
    if (digits.length > 8) {
      digits = digits.slice(0, 8);
    }
    return digits;
  }

  function toInt(text) {
    var n = 0;
    for (var i = 0; i < text.length; i++) {
      n = n * 10 + (text.charCodeAt(i) - 48);
    }
    return n;
  }

  function pad2(n) {
    return n < 10 ? "0" + n : String(n);
  }

  // Convierte dígitos crudos a DDMMYYYY con clamp suave:
  // día máx 31, mes máx 12; 4-9 al inicio de día/mes se antepone 0.
  function normalizeDateDigits(digits) {
    if (digits.length === 0) {
      return "";
    }

    var out = "";
    var i = 0;

    // Día
    var d0 = digits.charAt(i);
    if (d0 >= "4") {
      out += "0" + d0;
      i += 1;
    } else if (digits.length === 1) {
      return d0;
    } else {
      var day = toInt(digits.slice(i, i + 2));
      if (day > 31) {
        day = 31;
      }
      out += pad2(day);
      i += 2;
    }

    if (i >= digits.length) {
      return out;
    }

    // Mes
    var m0 = digits.charAt(i);
    if (m0 >= "2") {
      out += "0" + m0;
      i += 1;
    } else if (i === digits.length - 1) {
      out += m0;
      return out;
    } else {
      var month = toInt(digits.slice(i, i + 2));
      if (month > 12) {
        month = 12;
      }
      out += pad2(month);
      i += 2;
    }

    // Año (hasta 4 dígitos)
    out += digits.slice(i, i + 4);
    return out;
  }

  function formatDateFromDigits(digits) {
    var normalized = normalizeDateDigits(digits);
    if (normalized.length <= 2) {
      return normalized;
    }
    if (normalized.length <= 4) {
      return normalized.slice(0, 2) + "/" + normalized.slice(2);
    }
    return normalized.slice(0, 2) + "/" + normalized.slice(2, 4) + "/" + normalized.slice(4);
  }

  function applyDateMask(input) {
    if (input.readOnly || input.disabled) {
      return;
    }
    var digits = onlyDigits(input.value);
    var formatted = formatDateFromDigits(digits);
    if (input.value !== formatted) {
      input.value = formatted;
    }
  }

  function bindDateInput(input) {
    if (input.dataset.dateInputBound === "true") {
      return;
    }
    input.dataset.dateInputBound = "true";

    input.addEventListener("input", function () {
      applyDateMask(input);
    });

    input.addEventListener("paste", function (event) {
      event.preventDefault();
      var clipboard = event.clipboardData || window.clipboardData;
      var pasted = clipboard ? clipboard.getData("text") : "";
      input.value = formatDateFromDigits(onlyDigits(pasted));
    });
  }

  function initDateInputs(root) {
    (root || document).querySelectorAll("input[data-date-input]:not([readonly])").forEach(bindDateInput);
  }

  window.initDateInputs = initDateInputs;

  if (window.__dateInputBound) {
    return;
  }
  window.__dateInputBound = true;

  initDateInputs(document);
  document.body.addEventListener("htmx:afterSettle", function () {
    initDateInputs(document);
  });
})();
