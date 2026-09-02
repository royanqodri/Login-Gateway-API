#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -euo pipefail

# =========================
# COLORS
# =========================
GREEN='\033[0;32m'
CYAN='\033[1;36m'
RESET='\033[0m'

# === CONFIGURABLE VARIABLES ===
REGISTRY_HOST=""     # Your private registry
GITHUB_REPOSITORY=""  # Your GitHub repo path
GITHUB_REF="$1"                                           # Passed tag ref, e.g., refs/tags/v1.0.0
GITHUB_SHA=$(git rev-parse HEAD)
GITHUB_ACTOR=$(git config user.name || echo "unknown")

# === DERIVED VARIABLES ===
REPO_NAME=$(basename "$GITHUB_REPOSITORY" | tr '[:upper:]' '[:lower:]')
VERSION=$(echo "$GITHUB_REF" | sed 's#refs/tags/##')
VERSION_TAG="${REGISTRY_HOST}/${REPO_NAME}:${VERSION}"
LATEST_TAG="${REGISTRY_HOST}/${REPO_NAME}:latest"

# === FUNCTIONS ===
check_tag_format() {
  if [[ ! "$GITHUB_REF" =~ ^refs/tags/v.* ]]; then
    echo "Not a version tag (v*). Skipping build."
    exit 0
  fi
}

setup_buildx() {
  docker buildx use default || docker buildx create --use
}

build_and_push() {
  docker buildx build \
    --platform linux/amd64 \
    --push \
    --tag "$VERSION_TAG" \
    --tag "$LATEST_TAG" \
    --label "org.opencontainers.image.source=https://github.com/${GITHUB_REPOSITORY}" \
    --label "org.opencontainers.image.revision=${GITHUB_SHA}" \
    --label "org.opencontainers.image.version=${VERSION}" \
    --label "org.opencontainers.image.created=$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
    --label "org.opencontainers.image.authors=${GITHUB_ACTOR}" \
    .
}

log() {
  echo -e "${CYAN}[INFO]${RESET} $1"
}

success() {
  echo -e "${GREEN}[DONE]${RESET} $1"
}

send_discord_notification() {
  local status="$1"
  local description="$2"
  local color
  local emoji

  if [ -z "${DISCORD_WEBHOOK_URL:-}" ]; then
    log "DISCORD_WEBHOOK_URL is not set. Skipping Discord notification."
    return
  fi

  if [ "$status" = "success" ]; then
    color=3066993
    emoji="✅"
  else
    color=15158332
    emoji="❌"
  fi

  curl -s -X POST -H "Content-Type: application/json" \
    -d "{
      \"content\": \"$emoji **$status**\",
      \"embeds\": [
        {
          \"title\": \"Deploy $status 🚀\",
          \"description\": \"$description\",
          \"color\": $color,
          \"footer\": {
            \"text\": \"Deploy Script\"
          },
          \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
        }
      ]
    }" \
    "$DISCORD_WEBHOOK_URL" >/dev/null
}

# cleanup_docker() {
  # Uncomment if needed in local runner
  # docker ps -aq | xargs -r docker rm -f || true
  # docker rmi "$VERSION_TAG" "$LATEST_TAG" || true
  # docker system prune -af --volumes || true
  # docker buildx prune -af || true
  # docker volume prune -af || true
# }

# =========================
# TRAP FAILURE
# =========================
trap 'send_discord_notification "failure" "Error Docker image built and pushed: ${REGISTRY_HOST}/${REPO_NAME}:${VERSION} and ${REGISTRY_HOST}/${REPO_NAME}:latest."' ERR


# === MAIN ===
check_tag_format
setup_buildx
build_and_push
# cleanup_docker

# =========================
# SUCCESS
# =========================
success "✅ Docker image built and pushed: ${REGISTRY_HOST}/${REPO_NAME}:${VERSION} and ${REGISTRY_HOST}/${REPO_NAME}:latest"
send_discord_notification "success" "✅ Docker image built and pushed: ${REGISTRY_HOST}/${REPO_NAME}:${VERSION} and ${REGISTRY_HOST}/${REPO_NAME}:latest"
