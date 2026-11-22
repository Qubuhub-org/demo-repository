#!/usr/bin/env bash
# Runtime transcript for the shared-session fix: N runs on one provider session through the real RunExecutionService and
# ProviderCommitBatch (TestClock, virtual time), with main's RunExecutionService.ts and with this PR's. Usage, with a
# checkout of the PR head whose dependencies are installed (pnpm install / vp i):
#   bash reproduce.sh <path-to-checkout>
# Appends block.ts to RunExecutionService.test.ts for the run only and restores both files on exit.
set -euo pipefail
WT=$(cd "$1" && pwd); HERE=$(cd "$(dirname "$0")" && pwd)
BASE=6fd14f46d667ff4db712478e8c877a1de80a80fe
SRC=apps/server/src/orchestration-v2
RAW=$(pwd)/transcript-raw.txt
cd "$WT"
git diff --quiet HEAD -- "$SRC/RunExecutionService.ts" "$SRC/RunExecutionService.test.ts" ||
  { echo "local changes in the two files; refusing" >&2; exit 1; }
restore() { git -C "$WT" checkout -q HEAD -- "$SRC/RunExecutionService.ts" "$SRC/RunExecutionService.test.ts"; }
trap restore EXIT
if command -v vp >/dev/null 2>&1; then VP=(vp); else VP=(npx --no-install vp); fi
: > "$RAW"
cat "$HERE/block.ts" >> "$SRC/RunExecutionService.test.ts"
for side in base fix; do
  if [ "$side" = base ]; then git show "$BASE:$SRC/RunExecutionService.ts" > "$SRC/RunExecutionService.ts"
  else git checkout -q HEAD -- "$SRC/RunExecutionService.ts"; fi
  echo "== $side"
  (cd apps/server && SIDE=$side MEASURE_OUT="$RAW" "${VP[@]}" test run src/orchestration-v2/RunExecutionService.test.ts -t TRANSCRIPT)
done
python3 -I - "$RAW" <<'PY'
import sys, re
rows = {}
for line in open(sys.argv[1]):
    kv = dict(re.findall(r"(\w+)=(\S+)", line)); side = line.split()[0]
    rows[(side, int(kv["runs"]), int(kv["minute"]))] = kv
print()
print("| runs | base: rounds behind at 1 / 5 / 10 min | fix: rounds behind at 1 / 5 / 10 min | base: commits, snapshot jobs (10 min) | fix: commits, snapshot jobs (10 min) |")
print("|---|---|---|---|---|")
for runs in sorted({k[1] for k in rows}):
    def behind(side):
        return " / ".join(
            (lambda r: r["behind_min"] if r["behind_min"] == r["behind_max"] else f'{r["behind_min"]}-{r["behind_max"]}')(rows[(side, runs, m)])
            for m in (1, 5, 10))
    def work(side):
        r = rows[(side, runs, 10)]; return f'{int(r["commits"]):,}, {int(r["jobs"]):,}'
    print(f"| {runs} | {behind('base')} | {behind('fix')} | {work('base')} | {work('fix')} |")
PY
