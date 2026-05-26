// modals.jsx — Add/Edit/View/Delete dialogs for Vehicle Colors

const PRESET_COLORS = [
  '#0E0E10', '#1B1B1F', '#2A2D34', '#5C636B', '#C9CFD6', '#F4F4F2',
  '#1E5EF7', '#0A2540', '#36B1D6', '#1F7A6A', '#1F6F3D', '#9ACD32',
  '#C8232C', '#7A1F2B', '#E8782D', '#C9A86B', '#D9C8A5', '#5B3924',
  '#6B3FA0', '#E8B4BC', '#7A5AE0', '#FF7A59', '#FFB86B', '#B8E4D2',
];

function ColorPickerField({ value, onChange }) {
  return (
    <div className="field">
      <label>Color value</label>
      <div style={{ display: 'flex', gap: 10, alignItems: 'center' }}>
        <label
          style={{
            position: 'relative', width: 56, height: 40, borderRadius: 8,
            border: '1px solid var(--border)', cursor: 'default', overflow: 'hidden',
            ...swatchStyle(value),
          }}
          title="Pick color"
        >
          <input
            type="color"
            value={value}
            onChange={e => onChange(e.target.value.toUpperCase())}
            style={{ position: 'absolute', inset: 0, opacity: 0, cursor: 'default' }}
          />
        </label>
        <input
          className="input"
          style={{ width: 140, fontFamily: 'var(--mono)', textTransform: 'uppercase' }}
          value={value}
          onChange={e => {
            let v = e.target.value.trim();
            if (v && !v.startsWith('#')) v = '#' + v;
            onChange(v.toUpperCase());
          }}
          maxLength={7}
        />
        <div style={{ display: 'flex', gap: 6, flexWrap: 'wrap', flex: 1 }}>
          {PRESET_COLORS.slice(0, 12).map(c => (
            <button
              key={c}
              type="button"
              onClick={() => onChange(c)}
              title={c}
              style={{
                width: 22, height: 22, borderRadius: 6, border: 0, padding: 0, cursor: 'default',
                outline: value.toUpperCase() === c.toUpperCase() ? '2px solid var(--primary)' : 'none',
                outlineOffset: 2,
                ...swatchStyle(c),
              }}
            />
          ))}
        </div>
      </div>
    </div>
  );
}

function ColorPreview({ name, hex }) {
  const txt = isLight(hex) ? '#0b1220' : '#ffffff';
  return (
    <div style={{
      borderRadius: 10, padding: 14, background: hex, color: txt,
      display: 'flex', alignItems: 'center', justifyContent: 'space-between',
      border: '1px solid rgba(11,18,32,.06)',
    }}>
      <div>
        <div style={{ fontSize: 12, opacity: .7, fontWeight: 500 }}>Preview</div>
        <div style={{ fontSize: 18, fontWeight: 600, letterSpacing: '-.01em', marginTop: 2 }}>{name || 'Color name'}</div>
      </div>
      <div style={{ fontFamily: 'var(--mono)', fontSize: 13, opacity: .8 }}>{hex}</div>
    </div>
  );
}

/* ─── Add / Edit Color modal ──────────────────────────────────────── */
function ColorFormModal({ mode, initial, onClose, onSave }) {
  const [name, setName] = useState(initial?.name || '');
  const [hex, setHex] = useState(initial?.hex || '#1E5EF7');
  const [status, setStatus] = useState(initial?.status || 'active');
  const [errors, setErrors] = useState({});

  function submit() {
    const errs = {};
    if (!name.trim()) errs.name = 'Name is required';
    if (!/^#[0-9A-F]{6}$/i.test(hex)) errs.hex = 'Use a 6-digit HEX like #1E5EF7';
    setErrors(errs);
    if (Object.keys(errs).length) return;
    onSave({ name: name.trim(), hex: hex.toUpperCase(), status });
  }

  return (
    <div className="modal-back" onMouseDown={e => { if (e.target === e.currentTarget) onClose(); }}>
      <div className="modal" role="dialog" aria-modal="true">
        <div className="modal-hd">
          <div>
            <h3>{mode === 'edit' ? 'Edit vehicle color' : 'Add new vehicle color'}</h3>
            <p>Colors are used for vehicle records and rider visibility in the app.</p>
          </div>
          <button className="tb-icon-btn" onClick={onClose} aria-label="Close"><Icon name="x" size={18} /></button>
        </div>
        <div className="modal-bd">
          <ColorPreview name={name} hex={hex} />
          <div className="field">
            <label>Color name</label>
            <input
              className="input"
              value={name}
              autoFocus
              placeholder="e.g. Pearl White"
              onChange={e => setName(e.target.value)}
            />
            {errors.name && <div className="hint" style={{ color: 'var(--danger)' }}>{errors.name}</div>}
          </div>
          <ColorPickerField value={hex} onChange={setHex} />
          {errors.hex && <div className="hint" style={{ color: 'var(--danger)' }}>{errors.hex}</div>}
          <div className="field">
            <label>Status</label>
            <div style={{ display: 'inline-flex', background: '#f1f3f8', borderRadius: 8, padding: 3 }}>
              {['active', 'inactive'].map(s => (
                <button
                  key={s}
                  type="button"
                  onClick={() => setStatus(s)}
                  className="btn-sm"
                  style={{
                    background: status === s ? '#fff' : 'transparent',
                    boxShadow: status === s ? '0 1px 2px rgba(11,18,32,.08)' : 'none',
                    border: 0, color: status === s ? '#0b1220' : '#5b6478',
                    fontWeight: status === s ? 600 : 500, padding: '6px 14px',
                    borderRadius: 6, cursor: 'default', textTransform: 'capitalize',
                  }}
                >{s}</button>
              ))}
            </div>
          </div>
        </div>
        <div className="modal-ft">
          <button className="btn btn-ghost" onClick={onClose}>Cancel</button>
          <button className="btn btn-primary" onClick={submit}>
            <Icon name="check" size={14} stroke={3} />
            {mode === 'edit' ? 'Save changes' : 'Add color'}
          </button>
        </div>
      </div>
    </div>
  );
}

/* ─── View Color modal ────────────────────────────────────────────── */
function ColorViewModal({ row, onClose, onEdit, tz }) {
  const author = userById(row.author);
  const created = formatWhen('2026-03-12T10:00:00Z', tz.offset);
  const updated = formatWhen(row.updatedAt, tz.offset);
  return (
    <div className="modal-back" onMouseDown={e => { if (e.target === e.currentTarget) onClose(); }}>
      <div className="modal" role="dialog" aria-modal="true">
        <div className="modal-hd">
          <div>
            <h3>{row.name}</h3>
            <p>Vehicle color details</p>
          </div>
          <button className="tb-icon-btn" onClick={onClose} aria-label="Close"><Icon name="x" size={18} /></button>
        </div>
        <div className="modal-bd">
          <ColorPreview name={row.name} hex={row.hex} />
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 12 }}>
            <Detail k="HEX" v={<span style={{ fontFamily: 'var(--mono)' }}>{row.hex}</span>} />
            <Detail k="Status" v={<StatusPill status={row.status} />} />
            <Detail k="Vehicles assigned" v={<b style={{ fontVariantNumeric: 'tabular-nums' }}>{row.vehicles}</b>} />
            <Detail k="Color ID" v={<span style={{ fontFamily: 'var(--mono)', color: '#5b6478' }}>{row.id.toUpperCase()}</span>} />
            <Detail k="Created" v={<><div>{created.abs}</div><div className="meta-sub">10 weeks ago</div></>} />
            <Detail k="Last updated" v={<><div>{updated.abs}</div><div className="meta-sub">{updated.rel}</div></>} />
            <Detail wide k="Updated by" v={
              <div className="who-cell">
                <Avatar user={author} size={28} />
                <div>
                  <div style={{ fontWeight: 500 }}>{author.name}</div>
                  <div className="meta-sub">{author.role} · {author.email}</div>
                </div>
              </div>
            } />
          </div>
        </div>
        <div className="modal-ft">
          <button className="btn btn-ghost" onClick={onClose}>Close</button>
          <button className="btn" onClick={onEdit}><Icon name="pencil" size={14} />Edit</button>
        </div>
      </div>
    </div>
  );
}

function Detail({ k, v, wide }) {
  return (
    <div style={{ gridColumn: wide ? '1 / -1' : 'auto', background: '#fbfcfe', border: '1px solid var(--border-2)', borderRadius: 10, padding: '10px 12px' }}>
      <div style={{ fontSize: 11.5, fontWeight: 600, textTransform: 'uppercase', letterSpacing: '.05em', color: 'var(--ink-3)' }}>{k}</div>
      <div style={{ marginTop: 4 }}>{v}</div>
    </div>
  );
}

/* ─── Delete confirm modal ────────────────────────────────────────── */
function DeleteModal({ rows, onClose, onConfirm }) {
  const isMulti = rows.length > 1;
  return (
    <div className="modal-back" onMouseDown={e => { if (e.target === e.currentTarget) onClose(); }}>
      <div className="modal" style={{ width: 460 }} role="dialog" aria-modal="true">
        <div className="modal-hd">
          <div style={{ display: 'flex', gap: 12 }}>
            <div style={{ width: 36, height: 36, borderRadius: 10, background: 'var(--danger-50)', color: 'var(--danger)', display: 'flex', alignItems: 'center', justifyContent: 'center', flex: '0 0 36px' }}>
              <Icon name="trash" size={18} />
            </div>
            <div>
              <h3>{isMulti ? `Delete ${rows.length} colors?` : `Delete "${rows[0].name}"?`}</h3>
              <p>This action cannot be undone. {rows.reduce((s, r) => s + r.vehicles, 0)} vehicle{rows.reduce((s, r) => s + r.vehicles, 0) === 1 ? '' : 's'} currently use{rows.reduce((s, r) => s + r.vehicles, 0) === 1 ? 's' : ''} {isMulti ? 'these colors' : 'this color'}.</p>
            </div>
          </div>
        </div>
        {isMulti && (
          <div className="modal-bd" style={{ paddingTop: 0 }}>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: 6 }}>
              {rows.slice(0, 8).map(r => (
                <span key={r.id} className="chip">
                  <span className="swatch" style={{ width: 12, height: 12, borderRadius: 3, ...swatchStyle(r.hex) }}></span>
                  {r.name}
                </span>
              ))}
              {rows.length > 8 && <span className="chip">+{rows.length - 8} more</span>}
            </div>
          </div>
        )}
        <div className="modal-ft">
          <button className="btn btn-ghost" onClick={onClose}>Cancel</button>
          <button className="btn btn-danger" onClick={onConfirm}>
            <Icon name="trash" size={14} />Delete {isMulti ? `${rows.length} colors` : 'color'}
          </button>
        </div>
      </div>
    </div>
  );
}

Object.assign(window, { ColorFormModal, ColorViewModal, DeleteModal, ColorPickerField, PRESET_COLORS });
