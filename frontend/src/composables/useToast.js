import iziToast from 'izitoast'

// ─── Global defaults ──────────────────────────────────────────────────────────
// Applied once at module load. Individual calls can still override these.
iziToast.settings({
  transitionIn:  'fadeInDown',
  transitionOut: 'fadeOutUp',
  progressBar:   true,
  closeOnClick:  true,
  displayMode:   1,       // 1 = replace same-type toasts (no stacking)
  position:      'topRight',
  timeout:       4000,
  zindex:        99999,
  maxWidth:      380,
  layout:        2,       // 2 = icon left, title/message right
})

// ─── Helpers ──────────────────────────────────────────────────────────────────
function isDark () {
  return document.documentElement.classList.contains('dark')
}

const palettes = {
  success: () => isDark()
    ? { bg: 'rgba(5,46,22,0.95)',    title: '#4ade80', msg: '#86efac', bar: '#22c55e', icon: '#4ade80' }
    : { bg: 'rgba(240,253,244,0.97)', title: '#15803d', msg: '#166534', bar: '#22c55e', icon: '#16a34a' },
  error: () => isDark()
    ? { bg: 'rgba(62,7,7,0.95)',      title: '#f87171', msg: '#fca5a5', bar: '#ef4444', icon: '#f87171' }
    : { bg: 'rgba(255,241,242,0.97)', title: '#991b1b', msg: '#b91c1c', bar: '#ef4444', icon: '#dc2626' },
  warning: () => isDark()
    ? { bg: 'rgba(62,36,3,0.95)',     title: '#fbbf24', msg: '#fcd34d', bar: '#f59e0b', icon: '#fbbf24' }
    : { bg: 'rgba(255,251,235,0.97)', title: '#92400e', msg: '#b45309', bar: '#f59e0b', icon: '#d97706' },
  info: () => isDark()
    ? { bg: 'rgba(9,16,45,0.95)',     title: '#818cf8', msg: '#a5b4fc', bar: '#6366f1', icon: '#818cf8' }
    : { bg: 'rgba(239,246,255,0.97)', title: '#3730a3', msg: '#4338ca', bar: '#6366f1', icon: '#4f46e5' },
}

function makeOpts (type, extra) {
  const p = palettes[type]()
  return {
    backgroundColor:  p.bg,
    titleColor:       p.title,
    messageColor:     p.msg,
    progressBarColor: p.bar,
    iconColor:        p.icon,
    class:            `izi-custom izi-custom--${type}`,
    ...extra,
  }
}

// ─── Composable ───────────────────────────────────────────────────────────────
export function useToast () {
  return {
    success: (message, title = '¡Listo!')  => iziToast.success(makeOpts('success', { title, message, timeout: 4000 })),
    error:   (message, title = 'Error')    => iziToast.error  (makeOpts('error',   { title, message, timeout: 6000 })),
    warning: (message, title = 'Aviso')    => iziToast.warning(makeOpts('warning', { title, message, timeout: 5000 })),
    info:    (message, title = '')         => iziToast.info   (makeOpts('info',    { title, message, timeout: 4000 })),
  }
}
