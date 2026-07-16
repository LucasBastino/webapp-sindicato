 function showErrorModal(errorMsg) {
      const modalContainer = document.getElementById("error-modal-container");
      const modal = document.getElementById("error-modal");
      modal.innerHTML = errorMsg;
      modalContainer.classList.remove("hidden");
    }

    function closeErrorModal() {
      document.getElementById("error-modal-container").classList.add("hidden");
    }

    // listen to htmx errors
    document.body.addEventListener("htmx:responseError", (event) => {
      showErrorModal(event.detail.xhr.responseText);
    });