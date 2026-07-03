#!/usr/bin/env bash
# ArcDLC installer — skills for Claude Code / Codex / OpenCode + the arctool CLI.
#
# One-line install:
#   curl -fsSL https://raw.githubusercontent.com/FrogoAI/arcdlc/main/install.sh | bash
#
# From a local clone (uses the checked-out tree, no download):
#   ./install.sh
#
# Options (flags or environment):
#   --agents LIST   comma-separated: claude,codex,opencode,all,none (default: auto-detect)
#   --bindir DIR    where to put arctool                (default: ~/.local/bin)   [ARCDLC_BINDIR]
#   --ref REF       git tag/branch for the skills       (default: main)           [ARCDLC_REF]
#   --skills-only   install the skills, skip arctool
#   --tool-only     install arctool, skip the skills
#   --uninstall     remove everything this script installs
#
# Supported arctool platforms: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64.

set -euo pipefail

REPO="FrogoAI/arcdlc"
PLUGIN="arcdlc"
TOOL="arctool"
SUBSKILLS="aic archive examinate execute plan policy source-map"

BINDIR="${ARCDLC_BINDIR:-$HOME/.local/bin}"
REF="${ARCDLC_REF:-main}"
AGENTS="${ARCDLC_AGENTS:-auto}"
DO_SKILLS=1
DO_TOOL=1
UNINSTALL=0

info() { printf '\033[1;32m==>\033[0m %s\n' "$*"; }
warn() { printf '\033[1;33mwarning:\033[0m %s\n' "$*" >&2; }
die()  { printf '\033[1;31merror:\033[0m %s\n' "$*" >&2; exit 1; }

while [ $# -gt 0 ]; do
  case "$1" in
    --agents) AGENTS="$2"; shift 2 ;;
    --bindir) BINDIR="$2"; shift 2 ;;
    --ref)    REF="$2"; shift 2 ;;
    --skills-only) DO_TOOL=0; shift ;;
    --tool-only)   DO_SKILLS=0; shift ;;
    --uninstall)   UNINSTALL=1; shift ;;
    -h|--help) sed -n '2,19p' "$0" 2>/dev/null || true; exit 0 ;;
    *) die "unknown option: $1 (see --help)" ;;
  esac
done

# --- platform detection (arctool binaries exist for these four targets) ---
detect_platform() {
  local os arch
  case "${ARCDLC_OS:-$(uname -s)}" in
    Linux|linux)   os=linux ;;
    Darwin|darwin) os=darwin ;;
    *) return 1 ;;
  esac
  case "${ARCDLC_ARCH:-$(uname -m)}" in
    x86_64|amd64)  arch=amd64 ;;
    aarch64|arm64) arch=arm64 ;;
    *) return 1 ;;
  esac
  echo "$os-$arch"
}

# --- agent detection ---
claude_dir="$HOME/.claude"
codex_dir="$HOME/.codex"
opencode_dir="${XDG_CONFIG_HOME:-$HOME/.config}/opencode"

resolve_agents() {
  case "$AGENTS" in
    auto)
      local found=""
      [ -d "$claude_dir" ]   && found="claude"
      [ -d "$codex_dir" ]    && found="$found codex"
      [ -d "$opencode_dir" ] && found="$found opencode"
      echo "$found" ;;
    all)  echo "claude codex opencode" ;;
    none) echo "" ;;
    *)    echo "$AGENTS" | tr ',' ' ' ;;
  esac
}

# --- uninstall ---
if [ "$UNINSTALL" = 1 ]; then
  info "removing ArcDLC skills and $TOOL"
  rm -rf "$claude_dir/skills/$PLUGIN"
  for s in $SUBSKILLS; do
    rm -rf "$codex_dir/skills/$PLUGIN-$s" "$opencode_dir/skills/$PLUGIN-$s"
  done
  rm -f "$BINDIR/$TOOL"
  info "done"
  exit 0
fi

# --- locate the source tree: local clone, or download the ref tarball ---
tmpdir=""
cleanup() { [ -n "$tmpdir" ] && rm -rf "$tmpdir"; }
trap cleanup EXIT

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]:-.}")" 2>/dev/null && pwd || echo /)"
if [ -f "$script_dir/.claude-plugin/plugin.json" ]; then
  src="$script_dir"
  info "using local checkout: $src"
else
  command -v curl >/dev/null || die "curl is required"
  command -v tar  >/dev/null || die "tar is required"
  tmpdir="$(mktemp -d)"
  info "downloading $REPO@$REF"
  curl -fsSL "https://github.com/$REPO/archive/refs/heads/$REF.tar.gz" -o "$tmpdir/src.tgz" \
    || curl -fsSL "https://github.com/$REPO/archive/refs/tags/$REF.tar.gz" -o "$tmpdir/src.tgz" \
    || die "cannot download $REPO@$REF"
  tar -xzf "$tmpdir/src.tgz" -C "$tmpdir"
  src="$(find "$tmpdir" -maxdepth 1 -type d -name "$PLUGIN-*" | head -1)"
  if [ -z "$src" ] || [ ! -f "$src/.claude-plugin/plugin.json" ]; then
    die "unexpected archive layout"
  fi
fi

# --- skills ---
if [ "$DO_SKILLS" = 1 ]; then
  agents="$(resolve_agents)"
  if [ -z "$agents" ]; then
    warn "no agent directories found (~/.claude, ~/.codex, ~/.config/opencode) — skipping skills."
    warn "re-run with --agents claude|codex|opencode|all to force."
  fi
  for agent in $agents; do
    case "$agent" in
      claude)
        # Prefer the official plugin CLI (Claude Code >= 2.1.157) when installing
        # from GitHub; otherwise drop the skills-directory plugin into
        # ~/.claude/skills, which Claude Code auto-loads as /arcdlc:<name>.
        if [ -z "${ARCDLC_NO_PLUGIN_CLI:-}" ] && [ ! -f "$script_dir/.claude-plugin/plugin.json" ] \
           && command -v claude >/dev/null && claude plugin --help >/dev/null 2>&1 \
           && claude plugin marketplace add "$REPO" >/dev/null 2>&1 \
           && claude plugin install "$PLUGIN@$PLUGIN" >/dev/null 2>&1; then
          info "Claude Code: installed via 'claude plugin' (marketplace $REPO, commands: /$PLUGIN:<name>)"
        else
          dest="$claude_dir/skills/$PLUGIN"
          rm -rf "$dest" && mkdir -p "$dest"
          cp -R "$src/.claude-plugin" "$src/skills" "$dest/"
          info "Claude Code: $dest (commands: /$PLUGIN:<name>)"
        fi ;;
      codex|opencode)
        # No plugin namespace: flatten each sub-skill to <plugin>-<name>.
        [ "$agent" = codex ] && root="$codex_dir/skills" || root="$opencode_dir/skills"
        mkdir -p "$root"
        for s in $SUBSKILLS; do
          rm -rf "${root:?}/$PLUGIN-$s"
          cp -R "$src/skills/$s" "$root/$PLUGIN-$s"
        done
        info "$agent: $root/$PLUGIN-<name> (invoke by skill name)" ;;
      *) warn "unknown agent '$agent' — skipped" ;;
    esac
  done
fi

# --- arctool ---
if [ "$DO_TOOL" = 1 ]; then
  platform="$(detect_platform)" || die "unsupported platform $(uname -s)/$(uname -m); supported: linux/darwin × amd64/arm64"
  mkdir -p "$BINDIR"
  installed=""

  # Prefer a released static binary (checksum-verified); fall back to building from source.
  if [ -z "$tmpdir" ]; then tmpdir="$(mktemp -d)"; fi
  base="https://github.com/$REPO/releases/latest/download"
  if command -v curl >/dev/null \
     && curl -fsSL "$base/$TOOL-$platform" -o "$tmpdir/$TOOL" 2>/dev/null \
     && curl -fsSL "$base/SHA256SUMS" -o "$tmpdir/SHA256SUMS" 2>/dev/null; then
    if command -v sha256sum >/dev/null; then sumcmd="sha256sum"; else sumcmd="shasum -a 256"; fi
    want="$(grep " $TOOL-$platform\$" "$tmpdir/SHA256SUMS" | awk '{print $1}')"
    got="$($sumcmd "$tmpdir/$TOOL" | awk '{print $1}')"
    if [ -z "$want" ] || [ "$want" != "$got" ]; then
      die "checksum mismatch for $TOOL-$platform"
    fi
    install -m 0755 "$tmpdir/$TOOL" "$BINDIR/$TOOL"
    installed="release binary ($platform, sha256 verified)"
  elif command -v go >/dev/null; then
    info "no release binary available — building from source with $(go version | awk '{print $3}')"
    (cd "$src" && CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o "$BINDIR/$TOOL" ./cmd/$TOOL)
    installed="built from source"
  else
    die "cannot install $TOOL: no release binary reachable and Go is not installed"
  fi

  info "$TOOL → $BINDIR/$TOOL ($installed)"
  "$BINDIR/$TOOL" version >/dev/null || die "$BINDIR/$TOOL does not run on this system"
  case ":$PATH:" in
    *":$BINDIR:"*) ;;
    *) warn "$BINDIR is not on PATH — add:  export PATH=\"$BINDIR:\$PATH\"" ;;
  esac
fi

info "ArcDLC install complete. Restart your agent session to pick up the skills."
