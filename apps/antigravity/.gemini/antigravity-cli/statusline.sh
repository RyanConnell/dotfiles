#!/bin/bash

# Read JSON payload from stdin
payload=$(cat)

# Extract basic model information
model=$(echo "$payload" | jq -r '.model.display_name // "N/A"')

# Extract remaining quotas for Gemini, convert to percentages, and round
gem_5h=$(echo "$payload" | jq -r '.quota["gemini-5h"].remaining_fraction // 0 | (. * 100) | round')
#gem_wk=$(echo "$payload" | jq -r '.quota["gemini-weekly"].remaining_fraction // 0 | (. * 100) | round')

# Extract remaining quotas for Claude/GPT (3p), convert to percentages, and round
p3_5h=$(echo "$payload" | jq -r '.quota["3p-5h"].remaining_fraction // 0 | (. * 100) | round')
#p3_wk=$(echo "$payload" | jq -r '.quota["3p-weekly"].remaining_fraction // 0 | (. * 100) | round')

# Build the final status bar string showing all limits
echo "🤖 ${model} | Quota(5-hour): Gemini: ${gem_5h}% | Claude/GPT: ${p3_5h}%"
