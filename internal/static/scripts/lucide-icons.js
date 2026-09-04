(function () {
  'use strict';

  function refreshIcons(event) {
    if (!window.lucide) return;

    var root = event && event.detail && event.detail.target;
    if (root && root.nodeType === 1) {
      lucide.createIcons({ root: root });
      return;
    }

    lucide.createIcons();
  }

  if (window.__lucideIconsBound) return;
  window.__lucideIconsBound = true;

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', refreshIcons);
  } else {
    refreshIcons();
  }

  document.body.addEventListener('htmx:afterSwap', refreshIcons);
  document.body.addEventListener('htmx:afterSettle', refreshIcons);
})();
