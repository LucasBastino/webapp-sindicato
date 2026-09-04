// Allow HTMX to swap server-rendered error modals on 4xx/5xx.
// Do not assign xhr.responseText via innerHTML (XSS surface).
(function () {
  document.body.addEventListener("htmx:beforeSwap", function (event) {
    var xhr = event.detail.xhr;
    if (!xhr || xhr.status < 400) {
      return;
    }

    // Session / auth redirects must not be swapped into the modal container.
    if (xhr.status === 401 || xhr.getResponseHeader("HX-Redirect")) {
      event.detail.shouldSwap = false;
      return;
    }

    var container = document.getElementById("app-modal-container");
    if (!container) {
      return;
    }

    event.detail.shouldSwap = true;
    event.detail.target = container;
  });
})();
