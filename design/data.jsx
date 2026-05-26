// data.jsx — mock domain data

const USERS = [
  { id: 'u1', name: 'Aarav Sharma',   email: 'aarav@fleet.local',   role: 'Admin',         color: '#ff7a59' },
  { id: 'u2', name: 'Priya Tan',      email: 'priya@fleet.local',   role: 'Fleet Manager', color: '#3aaed8' },
  { id: 'u3', name: 'Marcus Lee',     email: 'marcus@fleet.local',  role: 'Operator',      color: '#7a5af0' },
  { id: 'u4', name: 'Hui Min Goh',    email: 'huimin@fleet.local',  role: 'Operator',      color: '#16a34a' },
  { id: 'u5', name: 'Daniel Koh',     email: 'daniel@fleet.local',  role: 'Viewer',        color: '#d97706' },
  { id: 'u6', name: 'Siti Rahman',    email: 'siti@fleet.local',    role: 'Fleet Manager', color: '#ec4899' },
];

const ROLES = [
  { id: 'r1', name: 'Admin',         users: 2, permissions: 42, system: true,  desc: 'Full access to all modules and settings.' },
  { id: 'r2', name: 'Fleet Manager', users: 6, permissions: 28, system: false, desc: 'Manage vehicles, drivers, bookings, telematics.' },
  { id: 'r3', name: 'Operator',      users: 14, permissions: 16, system: false, desc: 'Day-to-day vehicle dispatch and maintenance log.' },
  { id: 'r4', name: 'Viewer',        users: 9, permissions: 7,  system: false, desc: 'Read-only access to fleet dashboards.' },
];

const VEHICLE_COLORS_SEED = [
  { id: 'c01', name: 'Pearl White',    hex: '#F4F4F2', vehicles: 142, status: 'active',   author: 'u2', updatedAt: '2026-05-24T09:18:00Z' },
  { id: 'c02', name: 'Midnight Black', hex: '#0E0E10', vehicles: 96,  status: 'active',   author: 'u1', updatedAt: '2026-05-22T14:02:00Z' },
  { id: 'c03', name: 'Storm Grey',     hex: '#5C636B', vehicles: 73,  status: 'active',   author: 'u3', updatedAt: '2026-05-21T07:44:00Z' },
  { id: 'c04', name: 'Cobalt Blue',    hex: '#1E5EF7', vehicles: 54,  status: 'active',   author: 'u2', updatedAt: '2026-05-20T16:30:00Z' },
  { id: 'c05', name: 'Crimson Red',    hex: '#C8232C', vehicles: 41,  status: 'active',   author: 'u4', updatedAt: '2026-05-18T11:11:00Z' },
  { id: 'c06', name: 'Forest Green',   hex: '#1F6F3D', vehicles: 22,  status: 'active',   author: 'u6', updatedAt: '2026-05-16T08:50:00Z' },
  { id: 'c07', name: 'Sand Beige',     hex: '#D9C8A5', vehicles: 18,  status: 'active',   author: 'u3', updatedAt: '2026-05-15T13:28:00Z' },
  { id: 'c08', name: 'Sunset Orange',  hex: '#E8782D', vehicles: 12,  status: 'active',   author: 'u5', updatedAt: '2026-05-14T18:02:00Z' },
  { id: 'c09', name: 'Sky Cyan',       hex: '#36B1D6', vehicles: 9,   status: 'inactive', author: 'u1', updatedAt: '2026-05-12T10:00:00Z' },
  { id: 'c10', name: 'Plum Purple',    hex: '#6B3FA0', vehicles: 6,   status: 'active',   author: 'u2', updatedAt: '2026-05-10T15:45:00Z' },
  { id: 'c11', name: 'Champagne Gold', hex: '#C9A86B', vehicles: 4,   status: 'inactive', author: 'u4', updatedAt: '2026-05-08T09:30:00Z' },
  { id: 'c12', name: 'Slate Teal',     hex: '#1F7A6A', vehicles: 3,   status: 'active',   author: 'u6', updatedAt: '2026-05-06T11:15:00Z' },
  { id: 'c13', name: 'Brick Maroon',   hex: '#7A1F2B', vehicles: 2,   status: 'inactive', author: 'u3', updatedAt: '2026-05-04T17:00:00Z' },
  { id: 'c14', name: 'Mint Frost',     hex: '#B8E4D2', vehicles: 1,   status: 'active',   author: 'u5', updatedAt: '2026-05-02T08:12:00Z' },
  { id: 'c15', name: 'Charcoal',       hex: '#2A2D34', vehicles: 28,  status: 'active',   author: 'u1', updatedAt: '2026-04-29T12:36:00Z' },
  { id: 'c16', name: 'Ice Silver',     hex: '#C9CFD6', vehicles: 33,  status: 'active',   author: 'u2', updatedAt: '2026-04-27T14:48:00Z' },
  { id: 'c17', name: 'Lava',           hex: '#A12E1B', vehicles: 0,   status: 'inactive', author: 'u3', updatedAt: '2026-04-25T10:21:00Z' },
  { id: 'c18', name: 'Lime',           hex: '#9ACD32', vehicles: 7,   status: 'active',   author: 'u6', updatedAt: '2026-04-22T09:00:00Z' },
  { id: 'c19', name: 'Onyx',           hex: '#1B1B1F', vehicles: 11,  status: 'active',   author: 'u5', updatedAt: '2026-04-20T19:30:00Z' },
  { id: 'c20', name: 'Rose Quartz',    hex: '#E8B4BC', vehicles: 2,   status: 'inactive', author: 'u4', updatedAt: '2026-04-18T16:00:00Z' },
  { id: 'c21', name: 'Navy',           hex: '#0A2540', vehicles: 19,  status: 'active',   author: 'u2', updatedAt: '2026-04-15T08:10:00Z' },
  { id: 'c22', name: 'Coffee',         hex: '#5B3924', vehicles: 5,   status: 'active',   author: 'u1', updatedAt: '2026-04-12T11:50:00Z' },
];

const MAKES = [
  'Toyota', 'Honda', 'Hyundai', 'KIA', 'Mazda', 'Nissan', 'Mitsubishi', 'BMW',
  'Mercedes-Benz', 'Volkswagen', 'Tesla', 'BYD', 'Polestar', 'Renault', 'SsangYong', 'Volvo',
].map((n, i) => ({ id: 'mk' + i, name: n, models: 3 + (i % 8), vehicles: 12 + ((i * 17) % 80), updatedAt: '2026-05-' + String(10 + i % 18).padStart(2, '0') + 'T10:00:00Z', author: USERS[i % USERS.length].id }));

const VEHICLE_TYPES = [
  { id: 'vt1', name: 'Sedan',         icon: 'car',    vehicles: 184, updatedAt: '2026-05-22T10:00:00Z', author: 'u2' },
  { id: 'vt2', name: 'SUV',           icon: 'car',    vehicles: 112, updatedAt: '2026-05-21T10:00:00Z', author: 'u3' },
  { id: 'vt3', name: 'MPV',           icon: 'truck',  vehicles: 64,  updatedAt: '2026-05-20T10:00:00Z', author: 'u4' },
  { id: 'vt4', name: 'Hatchback',     icon: 'car',    vehicles: 41,  updatedAt: '2026-05-19T10:00:00Z', author: 'u1' },
  { id: 'vt5', name: 'Station Wagon', icon: 'car',    vehicles: 24,  updatedAt: '2026-05-18T10:00:00Z', author: 'u5' },
  { id: 'vt6', name: 'Van',           icon: 'truck',  vehicles: 17,  updatedAt: '2026-05-17T10:00:00Z', author: 'u6' },
  { id: 'vt7', name: 'Pickup',        icon: 'truck',  vehicles: 8,   updatedAt: '2026-05-16T10:00:00Z', author: 'u2' },
  { id: 'vt8', name: 'Coupe',         icon: 'car',    vehicles: 3,   updatedAt: '2026-05-15T10:00:00Z', author: 'u3' },
];

const FUEL_TYPES = [
  { id: 'f1', name: 'Petrol',   unit: 'L',   icon: 'fuel',    vehicles: 213, updatedAt: '2026-05-23T10:00:00Z', author: 'u2' },
  { id: 'f2', name: 'Diesel',   unit: 'L',   icon: 'fuel',    vehicles: 87,  updatedAt: '2026-05-22T10:00:00Z', author: 'u4' },
  { id: 'f3', name: 'Hybrid',   unit: 'L',   icon: 'pulse',   vehicles: 64,  updatedAt: '2026-05-21T10:00:00Z', author: 'u3' },
  { id: 'f4', name: 'Electric', unit: 'kWh', icon: 'battery', vehicles: 78,  updatedAt: '2026-05-20T10:00:00Z', author: 'u1' },
  { id: 'f5', name: 'CNG',      unit: 'kg',  icon: 'fuel',    vehicles: 6,   updatedAt: '2026-05-19T10:00:00Z', author: 'u5' },
];

const VEHICLE_STATUSES = [
  { id: 's1', name: 'Available',     color: '#16a34a', vehicles: 162, updatedAt: '2026-05-24T10:00:00Z', author: 'u1' },
  { id: 's2', name: 'In Use',        color: '#1e5ef7', vehicles: 88,  updatedAt: '2026-05-24T10:00:00Z', author: 'u2' },
  { id: 's3', name: 'Maintenance',   color: '#d97706', vehicles: 23,  updatedAt: '2026-05-23T10:00:00Z', author: 'u3' },
  { id: 's4', name: 'Out of Service',color: '#dc2626', vehicles: 11,  updatedAt: '2026-05-22T10:00:00Z', author: 'u4' },
  { id: 's5', name: 'Reserved',      color: '#7a5af0', vehicles: 30,  updatedAt: '2026-05-21T10:00:00Z', author: 'u5' },
  { id: 's6', name: 'Cleaning',      color: '#3aaed8', vehicles: 14,  updatedAt: '2026-05-21T10:00:00Z', author: 'u6' },
  { id: 's7', name: 'Retired',       color: '#6b7280', vehicles: 3,   updatedAt: '2026-05-19T10:00:00Z', author: 'u1' },
];

const TIMEZONES = [
  { id: 'Asia/Singapore', label: 'Singapore (GMT+8)',  offset: 8 },
  { id: 'Asia/Tokyo',     label: 'Tokyo (GMT+9)',      offset: 9 },
  { id: 'Asia/Kolkata',   label: 'India (GMT+5:30)',   offset: 5.5 },
  { id: 'Europe/London',  label: 'London (GMT+1)',     offset: 1 },
  { id: 'America/New_York',label:'New York (GMT-4)',   offset: -4 },
  { id: 'UTC',            label: 'UTC',                offset: 0 },
];

function userById(id) { return USERS.find(u => u.id === id) || USERS[0]; }

// Pretty relative + absolute time
function formatWhen(iso, tzOffset) {
  const d = new Date(iso);
  const now = new Date('2026-05-26T12:00:00Z'); // pinned "now" for stable demo
  const diffMs = now - d;
  const mins = Math.round(diffMs / 60000);
  let rel;
  if (mins < 1)          rel = 'just now';
  else if (mins < 60)    rel = mins + ' min ago';
  else if (mins < 1440)  rel = Math.round(mins / 60) + ' h ago';
  else if (mins < 10080) rel = Math.round(mins / 1440) + ' d ago';
  else                   rel = Math.round(mins / 10080) + ' w ago';

  const shifted = new Date(d.getTime() + tzOffset * 3600 * 1000);
  const pad = n => String(n).padStart(2, '0');
  const abs = `${shifted.getUTCFullYear()}-${pad(shifted.getUTCMonth() + 1)}-${pad(shifted.getUTCDate())} ${pad(shifted.getUTCHours())}:${pad(shifted.getUTCMinutes())}`;
  return { rel, abs };
}

Object.assign(window, {
  USERS, ROLES, VEHICLE_COLORS_SEED, MAKES, VEHICLE_TYPES, FUEL_TYPES, VEHICLE_STATUSES,
  TIMEZONES, userById, formatWhen,
});
