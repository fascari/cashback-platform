#!/usr/bin/env bash

set -euo pipefail

root_dir="$(git rev-parse --show-toplevel)"
agents_dir="$root_dir/.agents"
skills_link="$agents_dir/skills"
desired_target="../.github/skills"

mkdir -p "$agents_dir"

current_target=""
if [ -L "$skills_link" ]; then
  current_target="$(readlink "$skills_link")"
fi

if [ "$current_target" != "$desired_target" ]; then
  ln -sfn "$desired_target" "$skills_link"
fi

echo "Codex skills linked at .agents/skills"
