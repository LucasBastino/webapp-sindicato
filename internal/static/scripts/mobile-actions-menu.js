(function () {
  function closeAllActionsMenus(exceptMenu) {
    document.querySelectorAll('.actions-menu.show').forEach(function (menu) {
      if (menu !== exceptMenu) menu.classList.remove('show');
    });
    document.querySelectorAll('.actions-trigger[aria-expanded="true"]').forEach(function (trigger) {
      if (!exceptMenu || !exceptMenu.closest('.actions-dropdown')?.contains(trigger)) {
        trigger.setAttribute('aria-expanded', 'false');
      }
    });
  }

  function initActionsDropdowns(root) {
    (root || document).querySelectorAll('.actions-dropdown').forEach(function (dropdown) {
      if (dropdown.dataset.actionsInit === 'true') return;
      dropdown.dataset.actionsInit = 'true';

      var toggle = dropdown.querySelector('.actions-trigger');
      var menu = dropdown.querySelector('.actions-menu');
      if (!toggle || !menu) return;

      // Ensure HTMX processes hx-* attrs inside the menu
      if (window.htmx) htmx.process(menu);

      toggle.addEventListener('click', function (e) {
        e.stopPropagation();
        var willOpen = !menu.classList.contains('show');
        closeAllActionsMenus(willOpen ? menu : null);
        menu.classList.toggle('show', willOpen);
        toggle.setAttribute('aria-expanded', willOpen ? 'true' : 'false');
        if (window.lucide) lucide.createIcons();
      });

      menu.querySelectorAll('button, a').forEach(function (el) {
        el.addEventListener('click', function () {
          menu.classList.remove('show');
          toggle.setAttribute('aria-expanded', 'false');
        });
      });
    });
  }

  document.addEventListener('click', function (e) {
    if (!e.target.closest('.actions-dropdown')) {
      closeAllActionsMenus(null);
    }
  });

  initActionsDropdowns(document);

  document.body.addEventListener('htmx:afterSwap', function () {
    initActionsDropdowns(document);
    if (window.lucide) lucide.createIcons();
  });

  window.initActionsDropdowns = initActionsDropdowns;
})();
