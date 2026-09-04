(function () {
  var COOKIE_NAME = "csrf_";

  function readCookie(name) {
    var parts = ("; " + document.cookie).split("; " + name + "=");
    if (parts.length < 2) return "";
    return decodeURIComponent(parts.pop().split(";").shift() || "");
  }

  function csrfToken() {
    return readCookie(COOKIE_NAME);
  }

  document.addEventListener("htmx:configRequest", function (e) {
    var token = csrfToken();
    if (token) {
      e.detail.headers["X-Csrf-Token"] = token;
    }
  });

  document.addEventListener(
    "submit",
    function (e) {
      var form = e.target;
      if (!form || form.tagName !== "FORM") return;

      var method = (form.getAttribute("method") || "get").toLowerCase();
      if (method === "get") return;

      var token = csrfToken();
      if (!token) return;

      var input = form.querySelector('input[name="csrf"]');
      if (!input) {
        input = document.createElement("input");
        input.type = "hidden";
        input.name = "csrf";
        form.appendChild(input);
      }
      input.value = token;
    },
    true
  );
})();
