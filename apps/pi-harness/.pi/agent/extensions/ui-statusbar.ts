import type { ExtensionAPI } from "@earendil-works/pi-coding-agent";
import * as path from "path";

/**
 * State object to hold the session telemetry.
 */
interface TelemetryState {
  modelName: string;
  contextUsage: string;
  maxContextWindow: number;
  readonly tokensPerSecond: string;
}

const telemetry: TelemetryState = {
  modelName: "Unknown Model",
  contextUsage: "0/0 (0%)",
  maxContextWindow: 0,
  get tokensPerSecond() {
    return (globalThis as any).__pi_tokens_per_second || "";
  }
};

/**
 * Removes ANSI escape sequences from a string and calculates its visible length.
 */
function stripAnsi(str: string): string {
  return str.replace(/\x1b\[[0-9;]*m/g, "");
}

/**
 * Formats a token count to k-notation if >= 1000 (e.g. 8200 -> 8.2k, 64000 -> 64k).
 */
function formatTokens(tokens: number): string {
  if (tokens >= 1000) {
    const kVal = tokens / 1000;
    return kVal % 1 === 0 ? `${kVal}k` : `${kVal.toFixed(1)}k`;
  }
  return tokens.toString();
}

/**
 * Reads model and context usage from the context object and updates telemetry.
 */
function updateTelemetry(ctx: any) {
  const model = ctx.model;
  if (model) {
    // NOTE: If we want the full modelName we can use '${model.name}' instead of '${model.id}'
    telemetry.modelName = `${model.id}`;
    telemetry.maxContextWindow = model.contextWindow || 0;
  }

  const usage = ctx.getContextUsage();
  if (usage) {
    const used = usage.tokens ?? 0;
    const max = usage.contextWindow ?? telemetry.maxContextWindow ?? 0;
    const percentage = usage.percent !== null
      ? usage.percent.toFixed(0)
      : (max > 0 ? ((used / max) * 100).toFixed(0) : "0");

    const usedStr = formatTokens(used);
    const maxStr = formatTokens(max);
    telemetry.contextUsage = `${usedStr}/${maxStr} (${percentage}%)`;
  } else {
    const maxStr = formatTokens(telemetry.maxContextWindow);
    telemetry.contextUsage = `0 / ${maxStr} (0%)`;
  }
}

export default function (pi: ExtensionAPI) {
  pi.on("session_start", async (event, ctx) => {
    updateTelemetry(ctx);

    ctx.ui.setFooter((_tui, theme) => {
      return {
        render: (width: number): string[] => {
          const cwd = process.cwd();
          const displayCwd = cwd.length > 30 ? "..." + path.basename(cwd) : cwd;

          // --- SEGMENT 1: SESSION & PERFORMANCE (Left) ---
          let sessionPart = theme.fg("accent", `🤖 ${telemetry.modelName}`) +
                            theme.fg("dim", ` | 🧠 ${telemetry.contextUsage}`);

          if (telemetry.tokensPerSecond) {
            sessionPart += theme.fg("dim", ` | ⚡ ${telemetry.tokensPerSecond}`);
          }

          // --- SEGMENT 2: ENVIRONMENT (Right) ---
          const envPart = theme.fg("dim", `📂 ${displayCwd}`);

          const fullLeftAndMiddle = sessionPart;
          // Calculate visible width of the left part
          const visibleWidthLeft = stripAnsi(fullLeftAndMiddle).length;
          // Calculate visible width of the right part
          const visibleWidthRight = stripAnsi(envPart).length;

          const totalAvailableWidth = width - 2; // Subtracting small margin for safety
          const paddingNeeded = Math.max(0, totalAvailableWidth - (visibleWidthLeft + visibleWidthRight) - 1);
          const padding = " ".repeat(paddingNeeded);

          return [`${fullLeftAndMiddle}${padding}${envPart}`];
        }
      };
    });
  });

  pi.on("model_select", async (event, ctx) => {
    updateTelemetry(ctx);
  });

  pi.on("message_start", async (event, ctx) => {
    if (event.message.role === "assistant") {
      (globalThis as any).__pi_message_start_time = Date.now();
    }
    updateTelemetry(ctx);
  });

  pi.on("message_end", async (event, ctx) => {
    const startTime = (globalThis as any).__pi_message_start_time || 0;
    if (event.message.role === "assistant" && startTime > 0) {
      const durationS = (Date.now() - startTime) / 1000;
      const usage = (event.message as any).usage;
      const outputTokens = usage?.output ?? 0;
      if (outputTokens > 0 && durationS > 0) {
        const tps = (outputTokens / durationS).toFixed(1);
        (globalThis as any).__pi_tokens_per_second = `${tps} t/s`;
      }
      (globalThis as any).__pi_message_start_time = 0;
    }
    updateTelemetry(ctx);
  });

  pi.on("session_compact", async (event, ctx) => {
    updateTelemetry(ctx);
  });
}
