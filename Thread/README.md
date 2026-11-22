# Live check: a resumed Codex thread's restored usage report (T3 Code, CodexAdapterV2)

One live run per checkout (`main` and the PR) of the shape in the recorded `provider_thread_resume` fixture:
a Codex thread runs one turn, the server restarts, and the next turn on the same thread resumes the Codex
thread in a new app-server. On resume, Codex sends `thread/tokenUsage/updated` for the thread's last,
finished turn. The question is what that report costs the app-server client's serial read loop.

- `ps011-live-probe.diff`: a logging-only probe, applied identically to both checkouts and off unless
  `PS011_LIVE_PROBE=1` and `PS011_LIVE_PROBE_FILE` are set. It writes one JSON line per Codex frame as the
  read loop takes it (incoming frames are stamped just before routing; requests as they are queued for writing), with a
  running count of `awaitRegisteredTurn` lookups and `Effect.yieldNow` calls, and one line per
  `thread/tokenUsage/updated` handler with the lookups and yields it spent and its time on the loop.
- `analyze.py`: reads the probe file and prints the frame order around the restored report and the counts.
- `result.txt`: its output for this run. Turn ids are shown as T1 (first runtime's turn) and T2.

## Steps

From each checkout, with dependencies installed (`vp i`), Codex installed and signed in, and Node 26.13.1:

```sh
git apply ps011-live-probe.diff
mkdir <REPO> && git -C <REPO> init && (echo scratch > <REPO>/README.md) && git -C <REPO> add . && git -C <REPO> commit -m init

# Runtime 1: one thread, one short turn.
PS011_LIVE_PROBE=1 PS011_LIVE_PROBE_FILE=<BASE_DIR>/ps011-probe.ndjson \
  T3CODE_HOME=<BASE_DIR> T3CODE_DISABLE_AUTO_UPDATE=1 T3CODE_TELEMETRY_ENABLED=false T3CODE_TRACE_MAX_FILES=60 \
  node apps/server/src/bin.ts serve --base-dir <BASE_DIR> --port <PORT> --host 127.0.0.1 --no-browser <REPO> > server-1.log 2>&1 &
SERVER_PID=$!   # wait for "T3 Code server is ready." in server-1.log
node apps/server/src/bin.ts project add --base-dir <BASE_DIR> <REPO>
node apps/server/src/bin.ts thread start --base-dir <BASE_DIR> --provider codex --model <MODEL> \
  --option reasoningEffort=low --wait --json <REPO> "Reply with the single word ok."   # prints <THREAD_ID>
kill $SERVER_PID; wait $SERVER_PID

# Runtime 2: same <BASE_DIR>; the next turn resumes the Codex thread in a new app-server.
PS011_LIVE_PROBE=1 PS011_LIVE_PROBE_FILE=<BASE_DIR>/ps011-probe.ndjson \
  T3CODE_HOME=<BASE_DIR> T3CODE_DISABLE_AUTO_UPDATE=1 T3CODE_TELEMETRY_ENABLED=false T3CODE_TRACE_MAX_FILES=60 \
  node apps/server/src/bin.ts serve --base-dir <BASE_DIR> --port <PORT> --host 127.0.0.1 --no-browser <REPO> > server-2.log 2>&1 &
SERVER_PID=$!   # wait for "T3 Code server is ready." in server-2.log
node apps/server/src/bin.ts thread send --base-dir <BASE_DIR> --wait --json <THREAD_ID> "Reply with the single word ok."
kill $SERVER_PID; wait $SERVER_PID

python3 analyze.py main=<BASE_DIR_MAIN>/ps011-probe.ndjson PR=<BASE_DIR_PR>/ps011-probe.ndjson
git checkout apps/server/src/orchestration-v2/Adapters/CodexAdapterV2.ts   # remove the probe
```

Use a fresh `<BASE_DIR>` for each checkout and one server at a time. The frames are also in the server's own
provider logs, `<BASE_DIR>/userdata/logs/provider/` (millisecond timestamps, transient deltas left out).

## Result

macOS 27.0.1 arm64 (18 cores, load average about 7 during the runs), Node 26.13.1, Codex CLI 0.169.0 via an
OpenAI-compatible proxy, model `gpt-6-luna` at reasoning effort low. One run per checkout. Both turns completed
with the reply `ok` on both checkouts.

Codex 0.169.0 sent the restored report for T1 right after the `thread/resume` response, before T3 sent
`thread/name/set` and `turn/start` (the recorded 0.156.1 fixture has it after `turn/start`). Frame order in the
second runtime, ms from the read loop taking the restored report (`mcpServer/startupStatus/updated` frames
left out):

| main | PR |
|---|---|
| -1.5 in `thread/resume` response (3) | -1.4 in `thread/resume` response (3) |
| **0.0 in `thread/tokenUsage/updated` T1** | **0.0 in `thread/tokenUsage/updated` T1** |
| 16.8 out `thread/name/set` (4) | 0.4 in `thread/goal/cleared` |
| 29.4 in `thread/goal/cleared` | 16.9 out `thread/name/set` (4) |
| 30.1 in response 4 | 19.2 in response 4 |
| 31.3 out `config/read` (5) | 19.8 out `config/read` (5) |
| 31.5 in `thread/name/updated` | 20.0 in `thread/name/updated` |
| 36.3 in response 5 | 26.1 in response 5 |
| 38.3 out `turn/start` (6) | 27.5 out `turn/start` (6) |
| 44.9 in `thread/settings/updated` | 33.5 in `thread/settings/updated` |
| 46.1 in response 6 (T2) | 34.3 in response 6 (T2) |
| 51.0 in `thread/status/changed` | 62.9 in `thread/status/changed` |
| 51.4 in `turn/started` T2 | 63.3 in `turn/started` T2 |

| Restored report for T1 | main | PR |
|---|---|---|
| `awaitRegisteredTurn` lookups in its handler | 1,001 | 0 (one direct lookup, not counted) |
| `Effect.yieldNow` calls in its handler | 1,000 | 0 |
| handler time on the read loop (wall) | 28.7 ms | 0.1 ms |
| restored report to the read loop taking the next frame (`thread/goal/cleared`) (wall) | 29.4 ms | 0.4 ms |
| `thread/name/set` sent to its response being read (wall) | 13.4 ms | 2.3 ms |
| `config/read` / `turn/start` sent to response read (wall) | 5.0 / 7.8 ms | 6.3 / 6.8 ms |
| restored report to `turn/start` sent (wall) | 38.3 ms | 27.5 ms |
| restored report to `turn/started` read (wall) | 51.4 ms | 63.3 ms |

The other usage reports (T1's in the first runtime, T2's own) took 1 lookup and 0 yields on main and 0 / 0 on
the PR.

Reading:
- On main, the restored report's handler spent the full 1,000 yields and did not find T1; the read loop took no
  other frame for 29.4 ms. `thread/goal/cleared`, which the PR's read loop took 0.4 ms after the report, waited
  behind it, and so did the `thread/name/set` response that T3 waits for before it sends `turn/start`.
- `turn/start` went out 10.8 ms later on main than on the PR here. Time to `turn/started` does not separate the
  two in one run each: on the PR, `turn/started` came 29.0 ms after the `turn/start` response, against 5.3 ms
  on main, a gap the restored report's handler is not part of on either checkout.
- Wall times are from a lightly loaded machine; the yield count does not depend on load, the time per yield does.

Not covered: a fork, and other threads sharing the app-server during the spin (one thread here).
