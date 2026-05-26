// components.jsx — shared UI atoms

const { useState, useEffect, useRef, useCallback, useMemo, createContext, useContext } = React;

/* ─── Popover ─────────────────────────────────────────────────────── */
function Popover({ anchorRect, onClose, align = 'right', children, offsetY = 6, minWidth }) {
  const ref = useRef(null);
  useEffect(() => {
    function onDoc(e) { if (ref.current && !ref.current.contains(e.target)) onClose(); }
    function onKey(e) { if (e.key === 'Escape') onClose(); }
    document.addEventListener('mousedown', onDoc);
    document.addEventListener('keydown', onKey);
    return () => { document.removeEventListener('mousedown', onDoc); document.removeEventListener('keydown', onKey); };
  }, [onClose]);
  if (!anchorRect) return null;
  const style = {
    top: anchorRect.bottom + offsetY,
    minWidth: minWidth || 'auto',
  };
  if (align === 'right') style.right = window.innerWidth - anchorRect.right;
  else if (align === 'left') style.left = anchorRect.left;
  else style.left = anchorRect.left + (anchorRect.width / 2);
  return (
    <div ref={ref} className="pop" style={style}>{children}</div>
  );
}

/* ─── Avatar ──────────────────────────────────────────────────────── */
function Avatar({ user, size = 24 }) {
  if (!user) return null;
  const initials = user.name.split(' ').map(s => s[0]).slice(0, 2).join('');
  return (
    <span className="av" style={{ background: user.color, width: size, height: size, fontSize: Math.round(size * 0.42) }}>{initials}</span>
  );
}

/* ─── Status pill ─────────────────────────────────────────────────── */
function StatusPill({ status }) {
  const cls = status === 'active' ? 'pill-on' : (status === 'inactive' ? 'pill-off' : 'pill-warn');
  return (
    <span className={'pill ' + cls}>
      <span className="dot"></span>{status === 'active' ? 'Active' : status === 'inactive' ? 'Inactive' : status}
    </span>
  );
}

/* ─── Toast ───────────────────────────────────────────────────────── */
const ToastCtx = createContext({ push: () => {} });
function ToastProvider({ children }) {
  const [items, setItems] = useState([]);
  const push = useCallback((msg, kind = 'ok') => {
    const id = Math.random().toString(36).slice(2);
    setItems(arr => [...arr, { id, msg, kind }]);
    setTimeout(() => setItems(arr => arr.filter(x => x.id !== id)), 3000);
  }, []);
  return (
    <ToastCtx.Provider value={{ push }}>
      {children}
      <div className="toast-wrap">
        {items.map(t => (
          <div key={t.id} className={'toast ' + (t.kind === 'err' ? 'err' : '')}>
            <span className="t-ico"><Icon name={t.kind === 'err' ? 'x' : 'check'} size={12} stroke={3} /></span>
            <span>{t.msg}</span>
          </div>
        ))}
      </div>
    </ToastCtx.Provider>
  );
}
function useToast() { return useContext(ToastCtx); }

/* ─── Color helpers ───────────────────────────────────────────────── */
function hexToRgb(hex) {
  const m = hex.replace('#', '').match(/.{2}/g);
  if (!m) return { r: 0, g: 0, b: 0 };
  return { r: parseInt(m[0], 16), g: parseInt(m[1], 16), b: parseInt(m[2], 16) };
}
function isLight(hex) {
  const { r, g, b } = hexToRgb(hex);
  // perceived luminance
  const l = (0.299 * r + 0.587 * g + 0.114 * b) / 255;
  return l > 0.7;
}
function swatchStyle(hex) {
  return { background: hex, boxShadow: isLight(hex) ? 'inset 0 0 0 1px rgba(11,18,32,.08)' : 'inset 0 0 0 1px rgba(255,255,255,.18)' };
}

/* ─── Checkbox ────────────────────────────────────────────────────── */
function Checkbox({ checked, indeterminate, onChange }) {
  const ref = useRef(null);
  useEffect(() => { if (ref.current) ref.current.indeterminate = !!indeterminate; }, [indeterminate]);
  return <input ref={ref} type="checkbox" className="cb" checked={!!checked} onChange={e => onChange && onChange(e.target.checked)} />;
}

/* ─── Empty state ─────────────────────────────────────────────────── */
function Empty({ title, sub, action }) {
  return (
    <div className="empty">
      <div style={{ width: 56, height: 56, borderRadius: 14, background: '#eef1f7', display: 'inline-flex', alignItems: 'center', justifyContent: 'center', color: '#9aa3b8', marginBottom: 12 }}>
        <Icon name="search" size={22} />
      </div>
      <div style={{ color: '#0b1220', fontWeight: 600, marginBottom: 4 }}>{title}</div>
      {sub && <div style={{ fontSize: 13 }}>{sub}</div>}
      {action && <div style={{ marginTop: 14 }}>{action}</div>}
    </div>
  );
}

Object.assign(window, { Popover, Avatar, StatusPill, ToastProvider, useToast, hexToRgb, isLight, swatchStyle, Checkbox, Empty });
