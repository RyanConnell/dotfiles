#!/usr/bin/env python3
import json, os, subprocess, sys, time

try:
    cluster = os.environ.get("AI_PROXY_CLUSTER", "infra")
    cache_dir = os.path.join(
		os.environ.get("XDG_CACHE_HOME", os.path.expanduser("~/.cache")),
		"zellij",
		"zjstatus")
    cache_file = os.path.join(cache_dir, f"ai-proxy-budget-{cluster}")

    if os.path.exists(cache_file) and time.time() - os.path.getmtime(cache_file) < 60:
        print(open(cache_file).read().strip())
        sys.exit(0)

    key = open(os.path.expanduser("~/.ai-proxy-api-key")).read().strip()
    r = subprocess.run(
        ["curl", "-fsSL", "--max-time", "3", "-K", "-",
         f"https://ai-proxy.{cluster}.corp.arista.io/key/info"],
        input=f'header = "Authorization: Bearer {key}"\n',
        capture_output=True, text=True, timeout=5,
    )
    info = json.loads(r.stdout).get("info", {})
    spend, budget = info["spend"], info["max_budget"]
    output = f"${spend:.2f}/${int(budget)}"
    os.makedirs(cache_dir, exist_ok=True)
    open(cache_file, "w").write(output + "\n")
    print(output)
except Exception:
    print("$?/?")
