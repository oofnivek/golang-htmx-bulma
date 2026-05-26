// app.jsx — Fleet Console shell (sidebar + topbar + routing + tweaks)

const NAV = [
  {
    label: 'Overview',
    items: [
      { id: 'dashboard', label: 'Dashboard', icon: 'dashboard' },
    ],
  },
  {
    label: 'Vehicle Management',
    items: [
      { id: 'vehicles',         label: 'Vehicles',         icon: 'car',    badge: '331' },
      { id: 'vehicle-makes',    label: 'Vehicle Makes',    icon: 'building' },
      { id: 'vehicle-models',   label: 'Vehicle Models',   icon: 'layers' },
      { id: 'vehicle-types',    label: 'Vehicle Types',    icon: 'truck' },
      { id: 'vehicle-colors',   label: 'Vehicle Colors',   icon: 'palette' },
      { id: 'vehicle-fuels',    label: 'Vehicle Fuels',    icon: 'fuel' },
      { id: 'vehicle-statuses', label: 'Vehicle Statuses', icon: 'pulse' },
      { id: 'fuel-types',       label: 'Fuel Types',       icon: 'battery' },
    ],
  },
  {
    label: 'Operations',
    items: [
      { id: 'bookings',     label: 'Bookings',     icon: 'calendar' },
      { id: 'maintenance',  label: 'Maintenance',  icon: 'cog' },
      { id: 'fuel-logs',    label: 'Fuel Logs',    icon: 'fuel' },
      { id: 'trips',        label: 'Trip History', icon: 'history' },
    ],
  },
  {
    label: 'Access Management',
    items: [
      { id: 'users', label: 'Users', icon: 'users', badge: '31' },
      { id: 'roles', label: 'Roles', icon: 'shield' },
    ],
  },
  {
    label: 'Settings',
    items: [
      { id: 'preferences', label: 'Preferences', icon: 'sliders' },
      { id: 'org',         label: 'Organisation', icon: 'building' },
    ],
  },
];

const ROUTE_RENDERERS = {
  'dashboard':         (p) => <Dashboard />,
  'vehicle-colors':    (p) => <VehicleColorsPage {...p} />,
  'vehicle-makes':     (p) => <VehicleMakesPage {...p} />,
  'vehicle-types':     (p) => <VehicleTypesPage {...p} />,
  'vehicle-fuels':     (p) => <FuelTypesPage {...p} />, // alias used by user's menu
  'fuel-types':        (p) => <FuelTypesPage {...p} />,
  'vehicle-statuses':  (p) => <VehicleStatusesPage {...p} />,
  'users':             (p) => <UsersPage {...p} />,
  'roles':             (p) => <RolesPage {...p} />,
};

function Sidebar({ current, onNav }) {
  return (
    <aside className="sb">
      <div className="sb-brand">
        <div className="sb-brand-mark">F</div>
        <div style={{ flex: 1, minWidth: 0 }}>
          <div className="sb-brand-name">Fleet Console</div>
          <div className="sb-brand-sub">Admin · Production</div>
        </div>
        <button className="tb-icon-btn" style={{ width: 28, height: 28 }} title="Pin"><Icon name="pin" size={14} /></button>
      </div>
      <div className="sb-search">
        <div className="sb-search-inp">
          <Icon name="search" size={14} />
          <input placeholder="Search menu…" />
          <span className="kbd">⌘K</span>
        </div>
      </div>
      <nav className="sb-nav">
        {NAV.map(group => (
          <div key={group.label} className="sb-section">
            <div className="sb-section-h"><span>{group.label}</span></div>
            {group.items.map(item => {
              const active = current === item.id || (current === 'vehicle-fuels' && item.id === 'fuel-types');
              return (
                <div
                  key={item.id}
                  className={'sb-item' + (active ? ' active' : '')}
                  onClick={() => onNav(item.id)}
                >
                  <Icon name={item.icon} size={16} className="sb-ico" />
                  <span>{item.label}</span>
                  {item.badge && <span className="sb-badge">{item.badge}</span>}
                </div>
              );
            })}
          </div>
        ))}
      </nav>
      <div className="sb-foot">
        <div className="avatar">AS</div>
        <div className="who">
          <div className="who-name">Aarav Sharma</div>
          <div className="who-role">Admin</div>
        </div>
        <button className="tb-icon-btn" style={{ width: 28, height: 28 }} title="Settings"><Icon name="settings" size={14} /></button>
      </div>
    </aside>
  );
}

function TopBar({ tz, onOpenTweaks }) {
  return (
    <header className="tb">
      <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
        <button className="tb-icon-btn" title="Back"><Icon name="chevleft" size={16} /></button>
        <div className="tb-tz" title="Current time"><Icon name="clock" size={13} /> 12:00 · {tz.label}</div>
      </div>
      <div className="tb-search">
        <Icon name="search" size={14} />
        <input placeholder="Search vehicles, users, makes, models…" />
        <span style={{ fontSize: 11, color: 'var(--ink-3)', fontFamily: 'var(--mono)' }}>⌘K</span>
      </div>
      <button className="btn btn-sm"><Icon name="plus" size={13} stroke={3} />Create</button>
      <button className="tb-icon-btn" title="Tweaks" onClick={onOpenTweaks}><Icon name="sliders" size={16} /></button>
      <button className="tb-icon-btn" title="Notifications"><Icon name="bell" size={16} /><span className="dot"></span></button>
      <button className="tb-icon-btn" title="Help"><Icon name="info" size={16} /></button>
      <div className="tb-avatar" title="Aarav Sharma">AS</div>
    </header>
  );
}

const TWEAK_DEFAULTS = /*EDITMODE-BEGIN*/{
  "primary": "#1E5EF7",
  "density": "regular",
  "tz": "Asia/Singapore",
  "navMode": "light",
  "perPage": 10
}/*EDITMODE-END*/;

function App() {
  const [route, setRoute] = React.useState('vehicle-colors');
  const [t, setTweak] = useTweaks(TWEAK_DEFAULTS);
  const [tz, setTzState] = React.useState(() => TIMEZONES.find(z => z.id === t.tz) || TIMEZONES[0]);

  // Apply primary color
  React.useEffect(() => {
    document.documentElement.style.setProperty('--primary', t.primary);
    // derive a 600 (slightly darker) — simple darken
    document.documentElement.style.setProperty('--primary-600', shade(t.primary, -0.12));
    document.documentElement.style.setProperty('--primary-50', tint(t.primary, 0.92));
    document.documentElement.style.setProperty('--primary-100', tint(t.primary, 0.82));
  }, [t.primary]);

  React.useEffect(() => {
    const z = TIMEZONES.find(x => x.id === t.tz);
    if (z) setTzState(z);
  }, [t.tz]);

  function setTz(z) { setTzState(z); setTweak('tz', z.id); }

  const renderer = ROUTE_RENDERERS[route];
  const sharedProps = { tz, setTz, density: t.density, setDensity: (d) => setTweak('density', d) };

  return (
    <div className={'app density-' + t.density}>
      <Sidebar current={route} onNav={setRoute} />
      <TopBar tz={tz} onOpenTweaks={() => window.parent.postMessage({ type: '__edit_mode_available' }, '*')} />
      <main className="main">
        {renderer ? renderer(sharedProps) : <StubPage title={route.replace(/-/g, ' ').replace(/\b\w/g, c => c.toUpperCase())} crumb="Module" />}
      </main>

      <TweaksPanel title="Tweaks">
        <TweakSection label="Theme" />
        <TweakColor
          label="Primary accent"
          value={t.primary}
          options={['#1E5EF7', '#0066FF', '#0EA5E9', '#16A34A', '#7A5AE0', '#EC4899', '#F97316', '#0F172A']}
          onChange={v => setTweak('primary', v)}
        />
        <TweakRadio
          label="Density"
          value={t.density}
          options={['compact', 'regular', 'comfy']}
          onChange={v => setTweak('density', v)}
        />
        <TweakSection label="Defaults" />
        <TweakSelect
          label="Default timezone"
          value={t.tz}
          options={TIMEZONES.map(z => ({ value: z.id, label: z.label }))}
          onChange={v => setTweak('tz', v)}
        />
      </TweaksPanel>
    </div>
  );
}

// Tiny color helpers (no external lib)
function clamp(n, a, b) { return Math.max(a, Math.min(b, n)); }
function parseHex(h) {
  const s = h.replace('#', '');
  return { r: parseInt(s.slice(0, 2), 16), g: parseInt(s.slice(2, 4), 16), b: parseInt(s.slice(4, 6), 16) };
}
function toHex(r, g, b) {
  const h = n => clamp(Math.round(n), 0, 255).toString(16).padStart(2, '0');
  return '#' + h(r) + h(g) + h(b);
}
function shade(hex, amt) { // amt negative = darker
  const { r, g, b } = parseHex(hex);
  const f = 1 + amt;
  return toHex(r * f, g * f, b * f);
}
function tint(hex, amt) { // amt 0..1 toward white
  const { r, g, b } = parseHex(hex);
  return toHex(r + (255 - r) * amt, g + (255 - g) * amt, b + (255 - b) * amt);
}

ReactDOM.createRoot(document.getElementById('root')).render(
  <ToastProvider><App /></ToastProvider>
);
