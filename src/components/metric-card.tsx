import type { LucideIcon } from 'lucide-react';

interface MetricCardProps {
  readonly icon: LucideIcon;
  readonly label: string;
  readonly value: string;
  readonly delta: string;
  readonly tone: 'cyan' | 'violet' | 'amber';
}

const toneStyles: Record<MetricCardProps['tone'], string> = {
  cyan: 'from-cyan-400/25 to-sky-500/10 text-cyan-200 ring-cyan-300/25',
  violet: 'from-violet-400/25 to-fuchsia-500/10 text-violet-200 ring-violet-300/25',
  amber: 'from-amber-400/25 to-orange-500/10 text-amber-200 ring-amber-300/25',
};

export function MetricCard({ icon: Icon, label, value, delta, tone }: MetricCardProps) {
  return (
    <article className="rounded-[1.75rem] border border-white/10 bg-white/[0.07] p-5 shadow-2xl shadow-black/20 backdrop-blur">
      <div className="flex items-center justify-between gap-4">
        <div>
          <p className="text-sm font-medium text-slate-400">{label}</p>
          <p className="mt-3 text-3xl font-semibold tracking-tight text-white">{value}</p>
        </div>
        <div className={`rounded-2xl bg-gradient-to-br p-3 ring-1 ${toneStyles[tone]}`}>
          <Icon aria-hidden="true" className="h-6 w-6" />
        </div>
      </div>
      <p className="mt-4 text-sm font-medium text-emerald-300">{delta}</p>
    </article>
  );
}
