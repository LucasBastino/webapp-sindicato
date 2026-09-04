import { validatePaidAt } from './form-validation/validators.js';

function closeOverlay(overlay, onKeyDown) {
  document.removeEventListener('keydown', onKeyDown, true);
  overlay.remove();
}

export function openDatePickerModal({ title = 'Fecha de pago', initialValue = '', onConfirm }) {
  const overlay = document.createElement('div');
  overlay.className = 'date-picker-overlay show';
  overlay.setAttribute('role', 'presentation');

  const initial = (initialValue || '').trim();

  overlay.innerHTML = `
    <div class="date-picker-modal" role="dialog" aria-modal="true" aria-labelledby="date-picker-title">
      <div class="date-picker-header">
        <h4 id="date-picker-title">${title}</h4>
      </div>
      <div class="date-picker-body">
        <div class="input-wrapper has-icon">
          <span class="input-icon"><i data-lucide="calendar"></i></span>
          <input type="text" id="date-picker-input" value="${initial}" placeholder="DD/MM/AAAA" data-date-input inputmode="numeric" autocomplete="off" maxlength="10">
        </div>
        <p class="date-picker-error field-error" aria-live="polite"></p>
      </div>
      <div class="date-picker-actions">
        <button type="button" class="btn btn-secondary date-picker-cancel">Cancelar</button>
        <button type="button" class="btn btn-primary date-picker-confirm">Confirmar</button>
      </div>
    </div>
  `;

  document.body.appendChild(overlay);

  const modal = overlay.querySelector('.date-picker-modal');
  const input = overlay.querySelector('#date-picker-input');
  const errorEl = overlay.querySelector('.date-picker-error');
  const cancelBtn = overlay.querySelector('.date-picker-cancel');
  const confirmBtn = overlay.querySelector('.date-picker-confirm');

  if (typeof lucide !== 'undefined') {
    lucide.createIcons({ nodes: [overlay] });
  }

  if (typeof window.initDateInputs === 'function') {
    window.initDateInputs(overlay);
  }

  function dismiss() {
    closeOverlay(overlay, onKeyDown);
  }

  function confirm() {
    const value = input.value.trim();
    const result = validatePaidAt(value);
    if (!result.valid) {
      errorEl.textContent = result.error || 'Fecha no válida.';
      input.focus();
      return;
    }
    errorEl.textContent = '';
    if (typeof onConfirm === 'function') {
      onConfirm(value);
    }
    dismiss();
  }

  function onKeyDown(e) {
    if (e.key === 'Escape') {
      e.stopPropagation();
      e.preventDefault();
      dismiss();
    }
  }

  overlay.addEventListener('click', function (e) {
    if (e.target === overlay) {
      dismiss();
    }
  });

  modal.addEventListener('click', function (e) {
    e.stopPropagation();
  });

  cancelBtn.addEventListener('click', dismiss);
  confirmBtn.addEventListener('click', confirm);

  input.addEventListener('keydown', function (e) {
    if (e.key === 'Enter') {
      e.preventDefault();
      confirm();
    }
  });

  document.addEventListener('keydown', onKeyDown, true);
  input.focus();
  if (input.value) {
    input.select();
  }
}
