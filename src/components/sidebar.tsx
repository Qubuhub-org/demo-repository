import { BarChart3, CircleDollarSign, LayoutDashboard, Settings, ShieldCheck, Users } from 'lucide-react';

const navigationItems = [
  { label: 'Overview', icon: LayoutDashboard, active: true },
  { label: 'Analytics', icon: BarChart3, active: false },
  { label: 'Clients', icon: Users, active: false },
  { label: 'Revenue', icon: CircleDollarSign, active: false },
  { label: 'Security', icon: ShieldCheck, active: false },
  { label: 'Settings', icon: Settings, active: false },
] as const;

export function Sidebar() {
  return (
    <aside className="hidden w-72 shrink-0 rounded-[2rem] border border-white/10 bg-white/[0.06] p-5 shadow-2xl shadow-black/30 backdrop-blur-xl lg:block">
      <div className="flex items-center gap-3 px-2 py-3">
        <div className="grid h-11 w-11 place-items-center rounded-2xl bg-cyan-400 font-black text-slate-950">A</div>
        <div>
          <p className="text-lg font-bold text-white">Aura</p>
          <p className="text-xs uppercase tracking-[0.3em] text-cyan-200/70">Ecosystem</p>
        </div>
      </div>
      <nav className="mt-10 space-y-2" aria-label="Dashboard navigation">
        {navigationItems.map(({ label, icon: Icon, active }) => (
          <a
            className={`flex items-center gap-3 rounded-2xl px-4 py-3 text-sm font-semibold transition ${
              active ? 'bg-cyan-300 text-slate-950 shadow-lg shadow-cyan-500/20' : 'text-slate-400 hover:bg-white/10 hover:text-white'
            }`}
            href="/dashboard"
            key={label}
          >
            <Icon aria-hidden="true" className="h-5 w-5" />
            {label}
          </a>
        ))}
      </nav>
    </aside>
  );
}
