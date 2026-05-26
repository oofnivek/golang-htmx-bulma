// pages.jsx — Vehicle Colors (main) + generic list pages + stubs

const { useState: uS, useEffect: uE, useMemo: uM, useRef: uR, useCallback: uCb } = React;

/* ─── Pagination Footer (shared) ──────────────────────────────────── */
function PagFooter({ total, page, perPage, setPage, setPerPage, tz, setTz }) {
  const pages = Math.max(1, Math.ceil(total / perPage));
  const start = total === 0 ? 0 : (page - 1) * perPage + 1;
  const end = Math.min(total, page * perPage);

  const [tzOpen, setTzOpen] = uS(false);
  const [rppOpen, setRppOpen] = uS(false);
  const tzRef = uR(null), rppRef = uR(null);
  const [tzRect, setTzRect] = uS(null), [rppRect, setRppRect] = uS(null);

  function buildPages() {
    const list = [];
    const max = pages;
    const add = n => list.push(n);
    if (max <= 7) { for (let i = 1; i <= max; i++) add(i); return list; }
    add(1);
    if (page > 3) add('…');
    for (let i = Math.max(2, page - 1); i <= Math.min(max - 1, page + 1); i++) add(i);
    if (page < max - 2) add('…');
    add(max);
    return list;
  }

  return (
    <div className="pag">
      <div className="left" style={{ display: 'flex', alignItems: 'center', gap: 14 }}>
        <span>Showing <b style={{ color: 'var(--ink)' }}>{start}–{end}</b> of <b style={{ color: 'var(--ink)' }}>{total}</b></span>

        <span style={{ display: 'inline-flex', alignItems: 'center', gap: 6, color: 'var(--ink-3)' }}>
          <span>Rows</span>
          <button
            ref={rppRef}
            className="btn btn-sm"
            onClick={() => { setRppRect(rppRef.current.getBoundingClientRect()); setRppOpen(true); }}
          >
            {perPage}<Icon name="chevdown" size={12} />
          </button>
        </span>

        <span style={{ display: 'inline-flex', alignItems: 'center', gap: 6, color: 'var(--ink-3)' }}>
          <Icon name="globe" size={13} />
          <span>Timezone</span>
          <button
            ref={tzRef}
            className="btn btn-sm"
            onClick={() => { setTzRect(tzRef.current.getBoundingClientRect()); setTzOpen(true); }}
          >
            {tz.label}<Icon name="chevdown" size={12} />
          </button>
        </span>
      </div>

      <div className="pg-ctrls">
        <button className="pag-btn" disabled={page === 1} onClick={() => setPage(1)} title="First"><Icon name="chevsleft" size={14} /></button>
        <button className="pag-btn" disabled={page === 1} onClick={() => setPage(p => Math.max(1, p - 1))} title="Previous"><Icon name="chevleft" size={14} /></button>
        {buildPages().map((p, i) => (
          p === '…' ? <span key={i} style={{ padding: '0 6px', color: 'var(--ink-3)' }}>…</span>
          : <button key={i} className={'pag-btn' + (p === page ? ' on' : '')} onClick={() => setPage(p)}>{p}</button>
        ))}
        <button className="pag-btn" disabled={page === pages} onClick={() => setPage(p => Math.min(pages, p + 1))} title="Next"><Icon name="chevright" size={14} /></button>
        <button className="pag-btn" disabled={page === pages} onClick={() => setPage(pages)} title="Last"><Icon name="chevsright" size={14} /></button>
      </div>

      {tzOpen && (
        <Popover anchorRect={tzRect} onClose={() => setTzOpen(false)} align="right" minWidth={200}>
          <div className="pop-hd">Display times in</div>
          {TIMEZONES.map(z => (
            <div key={z.id} className={'opt' + (z.id === tz.id ? ' on' : '')} onClick={() => { setTz(z); setTzOpen(false); }}>
              <span>{z.label}</span>
              {z.id === tz.id && <Icon name="check" size={14} className="chk" stroke={3} />}
            </div>
          ))}
        </Popover>
      )}
      {rppOpen && (
        <Popover anchorRect={rppRect} onClose={() => setRppOpen(false)} align="left" minWidth={120}>
          <div className="pop-hd">Rows per page</div>
          {[10, 20, 50, 100].map(n => (
            <div key={n} className={'opt' + (n === perPage ? ' on' : '')} onClick={() => { setPerPage(n); setPage(1); setRppOpen(false); }}>
              <span>{n}</span>
              {n === perPage && <Icon name="check" size={14} className="chk" stroke={3} />}
            </div>
          ))}
        </Popover>
      )}
    </div>
  );
}

/* ─── Vehicle Colors page ─────────────────────────────────────────── */
function VehicleColorsPage({ tz, setTz, density, setDensity }) {
  const toast = useToast();
  const [rows, setRows] = uS(VEHICLE_COLORS_SEED);
  const [query, setQuery] = uS('');
  const [statusFilter, setStatusFilter] = uS('all'); // all | active | inactive
  const [sort, setSort] = uS({ key: 'updatedAt', dir: 'desc' });
  const [page, setPage] = uS(1);
  const [perPage, setPerPage] = uS(10);
  const [selected, setSelected] = uS(new Set());
  const [visibleCols, setVisibleCols] = uS({
    color: true, vehicles: true, status: true, updatedBy: true, updatedAt: true,
  });

  const [modal, setModal] = uS(null); // {kind:'add'|'edit'|'view'|'delete', row?, rows?}
  const [colsOpen, setColsOpen] = uS(false);
  const colsRef = uR(null);
  const [colsRect, setColsRect] = uS(null);

  const filtered = uM(() => {
    const q = query.trim().toLowerCase();
    let arr = rows.filter(r =>
      (statusFilter === 'all' || r.status === statusFilter) &&
      (!q || r.name.toLowerCase().includes(q) || r.hex.toLowerCase().includes(q))
    );
    const { key, dir } = sort;
    const m = dir === 'asc' ? 1 : -1;
    arr = [...arr].sort((a, b) => {
      const va = key === 'updatedBy' ? userById(a.author).name : a[key];
      const vb = key === 'updatedBy' ? userById(b.author).name : b[key];
      if (typeof va === 'number') return (va - vb) * m;
      return String(va).localeCompare(String(vb)) * m;
    });
    return arr;
  }, [rows, query, statusFilter, sort]);

  const total = filtered.length;
  const paged = filtered.slice((page - 1) * perPage, page * perPage);

  // Reset page on filter
  uE(() => { setPage(1); }, [query, statusFilter]);
  // Sync selection size limited to filtered set
  uE(() => { setSelected(new Set([...selected].filter(id => filtered.find(r => r.id === id)))); /* eslint-disable-next-line */ }, [filtered]);

  function toggleSort(key) {
    setSort(s => s.key === key ? { key, dir: s.dir === 'asc' ? 'desc' : 'asc' } : { key, dir: 'asc' });
  }

  function toggleAll() {
    const allOn = paged.every(r => selected.has(r.id)) && paged.length > 0;
    const ns = new Set(selected);
    if (allOn) paged.forEach(r => ns.delete(r.id));
    else paged.forEach(r => ns.add(r.id));
    setSelected(ns);
  }
  function toggleOne(id) {
    const ns = new Set(selected);
    if (ns.has(id)) ns.delete(id); else ns.add(id);
    setSelected(ns);
  }

  function handleSave(data) {
    if (modal.kind === 'add') {
      const id = 'c' + String(rows.length + 1).padStart(2, '0');
      setRows([{ id, ...data, vehicles: 0, author: 'u1', updatedAt: new Date().toISOString() }, ...rows]);
      toast.push(`Color "${data.name}" added`);
    } else {
      setRows(rows.map(r => r.id === modal.row.id ? { ...r, ...data, updatedAt: new Date().toISOString(), author: 'u1' } : r));
      toast.push(`Color "${data.name}" updated`);
    }
    setModal(null);
  }
  function handleDelete() {
    const ids = new Set(modal.rows.map(r => r.id));
    setRows(rows.filter(r => !ids.has(r.id)));
    const ns = new Set(selected); ids.forEach(id => ns.delete(id)); setSelected(ns);
    toast.push(`${modal.rows.length} color${modal.rows.length === 1 ? '' : 's'} deleted`, 'ok');
    setModal(null);
  }

  const allOnPage = paged.length > 0 && paged.every(r => selected.has(r.id));
  const someOnPage = paged.some(r => selected.has(r.id));

  const selectedRows = rows.filter(r => selected.has(r.id));

  return (
    <>
      <div className="ph">
        <div>
          <div className="crumbs">
            <span>Vehicle Management</span><Icon name="chevright" size={12} /><b>Vehicle Colors</b>
          </div>
          <h1>Vehicle Colors</h1>
          <p>Manage the palette of colors available for the fleet · {rows.length} total · {rows.filter(r => r.status === 'active').length} active</p>
        </div>
        <div className="ph-actions">
          <button className="btn"><Icon name="download" size={14} />Export</button>
          <button className="btn"><Icon name="upload" size={14} />Import</button>
          <button className="btn btn-primary" onClick={() => setModal({ kind: 'add' })}>
            <Icon name="plus" size={14} stroke={3} />Add new color
          </button>
        </div>
      </div>

      <div className="card">
        <div className="toolbar">
          <div className="tb-search-2">
            <Icon name="search" size={14} />
            <input placeholder="Search by name or HEX…" value={query} onChange={e => setQuery(e.target.value)} />
            {query && <button className="tb-icon-btn" style={{ width: 22, height: 22 }} onClick={() => setQuery('')}><Icon name="x" size={12} /></button>}
          </div>

          <div className="chips">
            {['all', 'active', 'inactive'].map(s => (
              <button
                key={s}
                className={'chip ' + (statusFilter === s ? 'chip-on' : '')}
                onClick={() => setStatusFilter(s)}
                style={{ border: 0, cursor: 'default' }}
              >
                <span style={{ textTransform: 'capitalize' }}>{s === 'all' ? 'All' : s}</span>
                <span style={{ fontVariantNumeric: 'tabular-nums', opacity: .7 }}>
                  {s === 'all' ? rows.length : rows.filter(r => r.status === s).length}
                </span>
              </button>
            ))}
          </div>

          <div style={{ flex: 1 }}></div>

          <div style={{ display: 'inline-flex', background: '#f1f3f8', borderRadius: 8, padding: 3 }}>
            {[
              { id: 'compact', icon: 'layers' },
              { id: 'regular', icon: 'columns' },
              { id: 'comfy',   icon: 'dashboard' },
            ].map(o => (
              <button
                key={o.id}
                onClick={() => setDensity(o.id)}
                title={`Density: ${o.id}`}
                className="tb-icon-btn"
                style={{
                  width: 30, height: 28, borderRadius: 6,
                  background: density === o.id ? '#fff' : 'transparent',
                  boxShadow: density === o.id ? '0 1px 2px rgba(11,18,32,.08)' : 'none',
                  color: density === o.id ? 'var(--ink)' : 'var(--ink-3)',
                }}
              >
                <Icon name={o.icon} size={14} />
              </button>
            ))}
          </div>

          <button
            ref={colsRef}
            className="btn btn-sm"
            onClick={() => { setColsRect(colsRef.current.getBoundingClientRect()); setColsOpen(true); }}
          >
            <Icon name="columns" size={14} />Columns<Icon name="chevdown" size={12} />
          </button>
        </div>

        {selectedRows.length > 0 && (
          <div className="bulk">
            <Icon name="check" size={14} stroke={3} />
            <b>{selectedRows.length}</b> selected
            <span className="grow"></span>
            <button className="btn btn-sm" onClick={() => setSelected(new Set())}>Clear</button>
            <button className="btn btn-sm"><Icon name="archive" size={13} />Archive</button>
            <button className="btn btn-sm" style={{ color: 'var(--danger)', borderColor: 'var(--danger-50)' }}
                    onClick={() => setModal({ kind: 'delete', rows: selectedRows })}>
              <Icon name="trash" size={13} />Delete
            </button>
          </div>
        )}

        <div style={{ overflowX: 'auto' }}>
          <table className="tbl">
            <thead>
              <tr>
                <th className="check"><Checkbox checked={allOnPage} indeterminate={!allOnPage && someOnPage} onChange={toggleAll} /></th>
                {visibleCols.color && (
                  <th><span className="sort" onClick={() => toggleSort('name')}>Color {sort.key === 'name' && <Icon name="chevdown" size={11} style={{ transform: sort.dir === 'asc' ? 'rotate(180deg)' : 'none' }} />}</span></th>
                )}
                {visibleCols.vehicles && (
                  <th className="col-num"><span className="sort" onClick={() => toggleSort('vehicles')}>Vehicles {sort.key === 'vehicles' && <Icon name="chevdown" size={11} style={{ transform: sort.dir === 'asc' ? 'rotate(180deg)' : 'none' }} />}</span></th>
                )}
                {visibleCols.status && (<th>Status</th>)}
                {visibleCols.updatedBy && (
                  <th><span className="sort" onClick={() => toggleSort('updatedBy')}>Updated by {sort.key === 'updatedBy' && <Icon name="chevdown" size={11} style={{ transform: sort.dir === 'asc' ? 'rotate(180deg)' : 'none' }} />}</span></th>
                )}
                {visibleCols.updatedAt && (
                  <th><span className="sort" onClick={() => toggleSort('updatedAt')}>Updated at ({tz.label.match(/\(([^)]+)\)/)?.[1] || tz.id}) {sort.key === 'updatedAt' && <Icon name="chevdown" size={11} style={{ transform: sort.dir === 'asc' ? 'rotate(180deg)' : 'none' }} />}</span></th>
                )}
                <th className="col-act">Actions</th>
              </tr>
            </thead>
            <tbody>
              {paged.length === 0 && (
                <tr><td colSpan={7}><Empty title="No colors match your filters" sub="Try clearing the search or status filter." action={
                  <button className="btn" onClick={() => { setQuery(''); setStatusFilter('all'); }}>Clear filters</button>
                } /></td></tr>
              )}
              {paged.map(r => {
                const u = userById(r.author);
                const when = formatWhen(r.updatedAt, tz.offset);
                return (
                  <tr key={r.id} className={selected.has(r.id) ? 'selected' : ''}>
                    <td className="check"><Checkbox checked={selected.has(r.id)} onChange={() => toggleOne(r.id)} /></td>
                    {visibleCols.color && (
                      <td>
                        <div className="colorname">
                          <span className="swatch" style={swatchStyle(r.hex)}></span>
                          <span className="nm">{r.name}</span>
                          <span className="hex">{r.hex}</span>
                        </div>
                      </td>
                    )}
                    {visibleCols.vehicles && (
                      <td className="col-num">{r.vehicles.toLocaleString()}</td>
                    )}
                    {visibleCols.status && (<td><StatusPill status={r.status} /></td>)}
                    {visibleCols.updatedBy && (
                      <td>
                        <div className="who-cell">
                          <Avatar user={u} size={26} />
                          <div>
                            <div style={{ fontWeight: 500 }}>{u.name}</div>
                            <div className="meta-sub">{u.role}</div>
                          </div>
                        </div>
                      </td>
                    )}
                    {visibleCols.updatedAt && (
                      <td>
                        <div>{when.abs}</div>
                        <div className="meta-sub">{when.rel}</div>
                      </td>
                    )}
                    <td className="col-act">
                      <div className="row-act">
                        <button className="ra-btn" title="View" onClick={() => setModal({ kind: 'view', row: r })}><Icon name="eye" size={15} /></button>
                        <button className="ra-btn" title="Edit" onClick={() => setModal({ kind: 'edit', row: r })}><Icon name="pencil" size={14} /></button>
                        <button className="ra-btn danger" title="Delete" onClick={() => setModal({ kind: 'delete', rows: [r] })}><Icon name="trash" size={14} /></button>
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>

        <PagFooter total={total} page={page} perPage={perPage} setPage={setPage} setPerPage={setPerPage} tz={tz} setTz={setTz} />
      </div>

      {colsOpen && (
        <Popover anchorRect={colsRect} onClose={() => setColsOpen(false)} align="right" minWidth={200}>
          <div className="pop-hd">Toggle columns</div>
          {[
            { k: 'color',     l: 'Color' },
            { k: 'vehicles',  l: 'Vehicles' },
            { k: 'status',    l: 'Status' },
            { k: 'updatedBy', l: 'Updated by' },
            { k: 'updatedAt', l: 'Updated at' },
          ].map(c => (
            <label key={c.k} className="opt" style={{ cursor: 'default' }}>
              <Checkbox checked={visibleCols[c.k]} onChange={v => setVisibleCols({ ...visibleCols, [c.k]: v })} />
              <span>{c.l}</span>
            </label>
          ))}
        </Popover>
      )}

      {modal?.kind === 'add' && <ColorFormModal mode="add" onClose={() => setModal(null)} onSave={handleSave} />}
      {modal?.kind === 'edit' && <ColorFormModal mode="edit" initial={modal.row} onClose={() => setModal(null)} onSave={handleSave} />}
      {modal?.kind === 'view' && <ColorViewModal row={modal.row} tz={tz} onClose={() => setModal(null)} onEdit={() => setModal({ kind: 'edit', row: modal.row })} />}
      {modal?.kind === 'delete' && <DeleteModal rows={modal.rows} onClose={() => setModal(null)} onConfirm={handleDelete} />}
    </>
  );
}

/* ─── Generic list page (for makes, types, fuels, statuses) ───────── */
function GenericListPage({ title, crumb, subtitle, addLabel, rows: rowsIn, columns, tz, setTz, density, setDensity }) {
  const [query, setQuery] = uS('');
  const [page, setPage] = uS(1);
  const [perPage, setPerPage] = uS(10);
  const [selected, setSelected] = uS(new Set());

  const filtered = uM(() => {
    const q = query.trim().toLowerCase();
    return rowsIn.filter(r => !q || r.name.toLowerCase().includes(q));
  }, [rowsIn, query]);

  const paged = filtered.slice((page - 1) * perPage, page * perPage);
  const allOnPage = paged.length > 0 && paged.every(r => selected.has(r.id));
  const someOnPage = paged.some(r => selected.has(r.id));

  return (
    <>
      <div className="ph">
        <div>
          <div className="crumbs"><span>{crumb}</span><Icon name="chevright" size={12} /><b>{title}</b></div>
          <h1>{title}</h1>
          <p>{subtitle}</p>
        </div>
        <div className="ph-actions">
          <button className="btn"><Icon name="download" size={14} />Export</button>
          <button className="btn btn-primary"><Icon name="plus" size={14} stroke={3} />{addLabel}</button>
        </div>
      </div>

      <div className="card">
        <div className="toolbar">
          <div className="tb-search-2">
            <Icon name="search" size={14} />
            <input placeholder="Search…" value={query} onChange={e => setQuery(e.target.value)} />
          </div>
          <div style={{ flex: 1 }}></div>
          <button className="btn btn-sm"><Icon name="filter" size={13} />Filters</button>
        </div>

        <div style={{ overflowX: 'auto' }}>
          <table className="tbl">
            <thead>
              <tr>
                <th className="check"><Checkbox checked={allOnPage} indeterminate={!allOnPage && someOnPage} onChange={() => {
                  const ns = new Set(selected);
                  if (allOnPage) paged.forEach(r => ns.delete(r.id));
                  else paged.forEach(r => ns.add(r.id));
                  setSelected(ns);
                }} /></th>
                {columns.map(c => <th key={c.key} className={c.numeric ? 'col-num' : ''}>{c.label}</th>)}
                <th>Updated by</th>
                <th>Updated at ({tz.label.match(/\(([^)]+)\)/)?.[1] || tz.id})</th>
                <th className="col-act">Actions</th>
              </tr>
            </thead>
            <tbody>
              {paged.length === 0 && <tr><td colSpan={columns.length + 4}><Empty title="No records" sub="Try a different search." /></td></tr>}
              {paged.map(r => {
                const u = userById(r.author);
                const when = formatWhen(r.updatedAt, tz.offset);
                return (
                  <tr key={r.id} className={selected.has(r.id) ? 'selected' : ''}>
                    <td className="check"><Checkbox checked={selected.has(r.id)} onChange={() => {
                      const ns = new Set(selected);
                      if (ns.has(r.id)) ns.delete(r.id); else ns.add(r.id);
                      setSelected(ns);
                    }} /></td>
                    {columns.map(c => <td key={c.key} className={c.numeric ? 'col-num' : ''}>{c.render(r)}</td>)}
                    <td>
                      <div className="who-cell">
                        <Avatar user={u} size={26} />
                        <div>
                          <div style={{ fontWeight: 500 }}>{u.name}</div>
                          <div className="meta-sub">{u.role}</div>
                        </div>
                      </div>
                    </td>
                    <td>
                      <div>{when.abs}</div>
                      <div className="meta-sub">{when.rel}</div>
                    </td>
                    <td className="col-act">
                      <div className="row-act">
                        <button className="ra-btn" title="View"><Icon name="eye" size={15} /></button>
                        <button className="ra-btn" title="Edit"><Icon name="pencil" size={14} /></button>
                        <button className="ra-btn danger" title="Delete"><Icon name="trash" size={14} /></button>
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>

        <PagFooter total={filtered.length} page={page} perPage={perPage} setPage={setPage} setPerPage={setPerPage} tz={tz} setTz={setTz} />
      </div>
    </>
  );
}

/* ─── Pages built on GenericListPage ─────────────────────────────── */
function VehicleMakesPage(props) {
  return <GenericListPage
    title="Vehicle Makes" crumb="Vehicle Management"
    subtitle={`Manufacturers in your fleet · ${MAKES.length} total`}
    addLabel="Add new make"
    rows={MAKES}
    columns={[
      { key: 'name', label: 'Make', render: r => <b style={{ fontWeight: 600 }}>{r.name}</b> },
      { key: 'models', label: 'Models', numeric: true, render: r => r.models },
      { key: 'vehicles', label: 'Vehicles', numeric: true, render: r => r.vehicles.toLocaleString() },
    ]}
    {...props}
  />;
}

function VehicleTypesPage(props) {
  return <GenericListPage
    title="Vehicle Types" crumb="Vehicle Management"
    subtitle="Body classifications used in vehicle records"
    addLabel="Add new type"
    rows={VEHICLE_TYPES}
    columns={[
      { key: 'name', label: 'Type', render: r => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
          <span style={{ width: 28, height: 28, borderRadius: 7, background: '#eef1f7', color: 'var(--ink-2)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}><Icon name={r.icon} size={15} /></span>
          <b style={{ fontWeight: 600 }}>{r.name}</b>
        </div>
      ) },
      { key: 'vehicles', label: 'Vehicles', numeric: true, render: r => r.vehicles.toLocaleString() },
    ]}
    {...props}
  />;
}

function FuelTypesPage(props) {
  return <GenericListPage
    title="Fuel Types" crumb="Vehicle Management"
    subtitle="Energy sources tracked on vehicles and refuelling logs"
    addLabel="Add new fuel type"
    rows={FUEL_TYPES}
    columns={[
      { key: 'name', label: 'Fuel', render: r => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
          <span style={{ width: 28, height: 28, borderRadius: 7, background: '#eef1f7', color: 'var(--ink-2)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}><Icon name={r.icon} size={15} /></span>
          <b style={{ fontWeight: 600 }}>{r.name}</b>
        </div>
      ) },
      { key: 'unit', label: 'Unit', render: r => <span style={{ fontFamily: 'var(--mono)', color: 'var(--ink-2)' }}>{r.unit}</span> },
      { key: 'vehicles', label: 'Vehicles', numeric: true, render: r => r.vehicles.toLocaleString() },
    ]}
    {...props}
  />;
}

function VehicleStatusesPage(props) {
  return <GenericListPage
    title="Vehicle Statuses" crumb="Vehicle Management"
    subtitle="Operational states used across the fleet"
    addLabel="Add new status"
    rows={VEHICLE_STATUSES}
    columns={[
      { key: 'name', label: 'Status', render: r => (
        <span className="pill" style={{ background: r.color + '22', color: r.color }}>
          <span className="dot" style={{ background: r.color }}></span>{r.name}
        </span>
      ) },
      { key: 'vehicles', label: 'Vehicles', numeric: true, render: r => r.vehicles.toLocaleString() },
    ]}
    {...props}
  />;
}

function UsersPage(props) {
  const rows = USERS.map(u => ({ ...u, vehicles: 0, author: u.id, updatedAt: '2026-05-20T10:00:00Z' }));
  return <GenericListPage
    title="Users" crumb="Access Management"
    subtitle={`Team members with access to the console · ${USERS.length} total`}
    addLabel="Invite user"
    rows={rows}
    columns={[
      { key: 'name', label: 'User', render: r => (
        <div className="who-cell">
          <Avatar user={r} size={32} />
          <div>
            <div style={{ fontWeight: 600 }}>{r.name}</div>
            <div className="meta-sub">{r.email}</div>
          </div>
        </div>
      ) },
      { key: 'role', label: 'Role', render: r => <span className="pill" style={{ background: '#eef1f7', color: 'var(--ink-2)' }}>{r.role}</span> },
      { key: 'status', label: 'Status', render: () => <StatusPill status="active" /> },
    ]}
    {...props}
  />;
}

function RolesPage(props) {
  const rows = ROLES.map(r => ({ ...r, author: 'u1', updatedAt: '2026-05-15T10:00:00Z' }));
  return <GenericListPage
    title="Roles & Permissions" crumb="Access Management"
    subtitle="Permission groups assigned to users"
    addLabel="Add new role"
    rows={rows}
    columns={[
      { key: 'name', label: 'Role', render: r => (
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <b style={{ fontWeight: 600 }}>{r.name}</b>
            {r.system && <span className="pill" style={{ background: '#fef6e7', color: '#b45309' }}>System</span>}
          </div>
          <div className="meta-sub">{r.desc}</div>
        </div>
      ) },
      { key: 'users', label: 'Users', numeric: true, render: r => r.users },
      { key: 'permissions', label: 'Permissions', numeric: true, render: r => r.permissions },
    ]}
    {...props}
  />;
}

/* ─── Dashboard (simple overview placeholder) ─────────────────────── */
function Dashboard() {
  const stats = [
    { label: 'Total vehicles', v: '331', d: '+12 this week', icon: 'car' },
    { label: 'Available now',   v: '162', d: '49% of fleet',   icon: 'check' },
    { label: 'In maintenance',  v: '23',  d: '-3 vs last week', icon: 'cog' },
    { label: 'Active users',    v: '31',  d: '6 online',       icon: 'users' },
  ];
  return (
    <>
      <div className="ph">
        <div>
          <div className="crumbs"><b>Dashboard</b></div>
          <h1>Good afternoon, Aarav</h1>
          <p>Here's what's happening across the fleet today.</p>
        </div>
      </div>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 14, marginBottom: 18 }}>
        {stats.map(s => (
          <div key={s.label} className="card" style={{ padding: 16, display: 'flex', flexDirection: 'column', gap: 10 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <span style={{ fontSize: 12, color: 'var(--ink-3)', fontWeight: 500 }}>{s.label}</span>
              <span style={{ width: 30, height: 30, borderRadius: 8, background: '#eef3ff', color: 'var(--primary)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}><Icon name={s.icon} size={15} /></span>
            </div>
            <div style={{ fontSize: 28, fontWeight: 700, letterSpacing: '-.02em' }}>{s.v}</div>
            <div style={{ fontSize: 12, color: 'var(--ink-3)' }}>{s.d}</div>
          </div>
        ))}
      </div>
      <div className="card stub">
        <Icon name="dashboard" size={26} />
        <h2>Operations overview</h2>
        <div>This area shows live telemetry, alerts and recent events.</div>
      </div>
    </>
  );
}

/* ─── Generic stub for unbuilt pages ─────────────────────────────── */
function StubPage({ title, crumb }) {
  return (
    <>
      <div className="ph">
        <div>
          <div className="crumbs"><span>{crumb}</span><Icon name="chevright" size={12} /><b>{title}</b></div>
          <h1>{title}</h1>
          <p>This module is under construction.</p>
        </div>
      </div>
      <div className="card stub">
        <Icon name="layers" size={26} />
        <h2>Coming soon</h2>
        <div>Follow the same pattern as Vehicle Colors to wire up CRUD for {title}.</div>
      </div>
    </>
  );
}

Object.assign(window, {
  VehicleColorsPage, VehicleMakesPage, VehicleTypesPage, FuelTypesPage, VehicleStatusesPage,
  UsersPage, RolesPage, Dashboard, StubPage,
});
