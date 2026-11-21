import { ArrowUpRight, Bell, CreditCard, Search, Users, WalletCards } from 'lucide-react';
import { MetricCard } from '../components/metric-card';
import { Sidebar } from '../components/sidebar';

const activityRows = [
  ['Mobile app redesign', 'UI Sprint', '72%', 'bg-cyan-300'],
  ['Identity wallet audit', 'Security', '91%', 'bg-violet-300'],
  ['Revenue automation', 'Finance', '64%', 'bg-amber-300'],
] as const;

export function DashboardPage() {
  return (
    <main className="min-h-screen overflow-hidden bg-[#07111f] p-4 text-slate-100 md:p-8">
      <div className="pointer-events-none fixed inset-0 bg-[radial-gradient(circle_at_top_left,_rgba(34,211,238,0.18),_transparent_34%),radial-gradient(circle_at_75%_15%,_rgba(168,85,247,0.16),_transparent_30%)]" />
      <div className="relative mx-auto flex max-w-7xl gap-6">
        <Sidebar />
        <section className="min-w-0 flex-1">
          <header className="flex flex-col gap-5 rounded-[2rem] border border-white/10 bg-white/[0.06] p-5 backdrop-blur-xl md:flex-row md:items-center md:justify-between">
            <div>
              <p className="text-sm font-semibold uppercase tracking-[0.35em] text-cyan-200/80">Dashboard</p>
              <h1 className="mt-2 text-4xl font-bold tracking-tight text-white md:text-5xl">Business overview</h1>
            </div>
            <div className="flex items-center gap-3">
              <label className="flex min-w-0 items-center gap-2 rounded-2xl bg-slate-950/50 px-4 py-3 ring-1 ring-white/10">
                <Search aria-hidden="true" className="h-5 w-5 text-slate-500" />
                <input className="w-full bg-transparent text-sm outline-none placeholder:text-slate-500" placeholder="Search" type="search" />
              </label>
              <button aria-label="Notifications" className="rounded-2xl bg-white/10 p-3 ring-1 ring-white/10">
                <Bell aria-hidden="true" className="h-5 w-5" />
              </button>
            </div>
          </header>

          <div className="mt-6 grid gap-6 xl:grid-cols-[1fr_21rem]">
            <div className="space-y-6">
              <section className="grid gap-4 md:grid-cols-3">
                <MetricCard icon={WalletCards} label="Total balance" value="$128.4K" delta="+12.5% from last month" tone="cyan" />
                <MetricCard icon={Users} label="Active users" value="24.8K" delta="+8.2% user growth" tone="violet" />
                <MetricCard icon={CreditCard} label="Transactions" value="8,942" delta="+18.1% completed" tone="amber" />
              </section>

              <section className="rounded-[2rem] border border-white/10 bg-white/[0.07] p-6 shadow-2xl shadow-black/20 backdrop-blur">
                <div className="flex items-center justify-between">
                  <div>
                    <h2 className="text-2xl font-semibold text-white">Performance</h2>
                    <p className="mt-1 text-sm text-slate-400">Revenue, customers, and system health</p>
                  </div>
                  <ArrowUpRight aria-hidden="true" className="h-6 w-6 text-cyan-200" />
                </div>
                <div className="mt-8 flex h-72 items-end gap-3 rounded-[1.5rem] bg-slate-950/35 p-5">
                  {[42, 68, 54, 81, 63, 96, 74, 88, 69, 92, 78, 100].map((height, index) => (
                    <div className="flex flex-1 flex-col items-center gap-3" key={height + index}>
                      <div
                        className="w-full rounded-t-2xl bg-gradient-to-t from-cyan-500 to-cyan-200 shadow-lg shadow-cyan-500/20"
                        style={{ height: `${height}%` }}
                      />
                      <span className="text-xs font-medium text-slate-500">{index + 1}</span>
                    </div>
                  ))}
                </div>
              </section>
            </div>

            <aside className="space-y-6">
              <section className="rounded-[2rem] border border-white/10 bg-cyan-300 p-6 text-slate-950 shadow-2xl shadow-cyan-900/30">
                <p className="text-sm font-bold uppercase tracking-[0.25em] text-slate-700">Aura Score</p>
                <div className="mt-8 grid place-items-center">
                  <div className="grid h-44 w-44 place-items-center rounded-full bg-slate-950 text-center text-white ring-[18px] ring-cyan-100/80">
                    <span className="text-5xl font-black">86%</span>
                  </div>
                </div>
                <p className="mt-8 text-sm font-semibold text-slate-700">Operational confidence is strong across product, security, and revenue channels.</p>
              </section>

              <section className="rounded-[2rem] border border-white/10 bg-white/[0.07] p-6 backdrop-blur">
                <h2 className="text-xl font-semibold text-white">Current work</h2>
                <div className="mt-5 space-y-4">
                  {activityRows.map(([name, group, progress, color]) => (
                    <div className="rounded-2xl bg-slate-950/35 p-4" key={name}>
                      <div className="flex items-center justify-between gap-3">
                        <div>
                          <p className="font-semibold text-white">{name}</p>
                          <p className="text-sm text-slate-500">{group}</p>
                        </div>
                        <span className="text-sm font-bold text-cyan-200">{progress}</span>
                      </div>
                      <div className="mt-3 h-2 rounded-full bg-white/10">
                        <div className={`h-2 rounded-full ${color}`} style={{ width: progress }} />
                      </div>
                    </div>
                  ))}
                </div>
              </section>
            </aside>
          </div>
        </section>
      </div>
    </main>
  );
}
