(function () {
  var toggle = document.getElementById('sidebar-menu-toggle');
  var overlay = document.getElementById('sidebar-overlay');
  if (!toggle || !overlay) return;

  function setOpen(open) {
    document.body.classList.toggle('sidebar-open', open);
    toggle.setAttribute('aria-expanded', open ? 'true' : 'false');
    toggle.setAttribute('aria-label', open ? 'Cerrar menú' : 'Abrir menú');
    overlay.hidden = !open;
    var icon = toggle.querySelector('i');
    if (icon) {
      icon.setAttribute('data-lucide', open ? 'x' : 'menu');
      if (window.lucide) lucide.createIcons();
    }
  }

  function closeSidebar() {
    setOpen(false);
  }

  toggle.addEventListener('click', function () {
    setOpen(!document.body.classList.contains('sidebar-open'));
  });

  overlay.addEventListener('click', closeSidebar);

  document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape') closeSidebar();
  });

  document.querySelectorAll('.sidebar-drawer nav a').forEach(function (link) {
    link.addEventListener('click', closeSidebar);
  });

  document.querySelector('.sidebar-drawer .support')?.addEventListener('click', closeSidebar);

  document.body.addEventListener('htmx:afterSettle', closeSidebar);
})();
