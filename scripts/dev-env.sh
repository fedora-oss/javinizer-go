#!/usr/bin/env bash
# scripts/dev-env.sh — Set up and manage the local development environment
#
# Usage:
#   ./scripts/dev-env.sh [command]
#
# Commands:
#   start   [db]   Start backend (and optional DB service). db: sqlite|postgres|mysql
#   stop           Stop all dev services
#   status         Show running services and DB connection info
#   logs    [svc]  Tail logs (svc: javinizer|postgres|mysql, default: all)
#   reset   [db]   Destroy volumes and restart fresh
#   help           Show this help message
#
# Examples:
#   ./scripts/dev-env.sh start           # SQLite (no external DB)
#   ./scripts/dev-env.sh start postgres  # Start with PostgreSQL
#   ./scripts/dev-env.sh start mysql     # Start with MySQL
#   ./scripts/dev-env.sh stop
#   ./scripts/dev-env.sh reset postgres
#
# Requirements:
#   - Go 1.25+  (for backend, CGO + SQLite)
#   - Node.js 20+  (for frontend hot-reload, optional)
#   - Docker + Docker Compose  (for postgres/mysql profiles)

set -euo pipefail

# ── Colours ────────────────────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

log_info()    { echo -e "${GREEN}[INFO]${NC}  $*"; }
log_warn()    { echo -e "${YELLOW}[WARN]${NC}  $*"; }
log_error()   { echo -e "${RED}[ERROR]${NC} $*"; }
log_section() { echo -e "\n${BOLD}${BLUE}━━━ $* ━━━${NC}"; }

# ── Prerequisite checks ────────────────────────────────────────────────────────
check_prerequisites() {
    local missing=0

    log_section "Checking Prerequisites"

    if command -v go >/dev/null 2>&1; then
        log_info "Go:     $(go version | awk '{print $3}')"
    else
        log_error "Go not found. Install from https://go.dev/dl/"
        missing=1
    fi

    if command -v node >/dev/null 2>&1; then
        log_info "Node:   $(node --version)"
    else
        log_warn "Node.js not found — frontend hot-reload unavailable"
    fi

    if command -v docker >/dev/null 2>&1 && docker version >/dev/null 2>&1; then
        log_info "Docker: $(docker --version | awk '{print $3}' | tr -d ',')"
    else
        log_warn "Docker not running — postgres/mysql profiles unavailable"
    fi

    if command -v air >/dev/null 2>&1; then
        log_info "Air:    $(air -v 2>/dev/null | head -1 || echo 'installed')"
    else
        log_warn "air not found (live reload disabled). Install: go install github.com/air-verse/air@latest"
    fi

    [[ "${missing}" -eq 0 ]] || { log_error "Fix missing prerequisites and retry."; exit 1; }
}

# ── SQLite local dev ───────────────────────────────────────────────────────────
start_sqlite() {
    log_section "Starting Local Dev — SQLite"
    cd "${PROJECT_ROOT}"

    # Ensure data directory exists
    mkdir -p data

    # Copy default config if absent
    if [[ ! -f data/config.yaml ]]; then
        log_info "Copying default config → data/config.yaml"
        cp configs/config.yaml.example data/config.yaml 2>/dev/null || \
        cp internal/config/config.yaml.example data/config.yaml 2>/dev/null || true
    fi

    export JAVINIZER_CONFIG="${PROJECT_ROOT}/data/config.yaml"
    export JAVINIZER_DB="${PROJECT_ROOT}/data/javinizer.db"
    export JAVINIZER_DB_TYPE="sqlite"

    log_info "Config: ${JAVINIZER_CONFIG}"
    log_info "DB:     ${JAVINIZER_DB}"
    log_info ""
    log_info "Starting API server (Ctrl+C to stop)..."
    log_info "Web UI:  http://localhost:8080"
    log_info "API:     http://localhost:8080/api/v1/"
    log_info "Docs:    http://localhost:8080/docs"
    echo ""

    if command -v air >/dev/null 2>&1; then
        exec air -c .air.toml
    else
        exec go run ./cmd/javinizer api
    fi
}

# ── PostgreSQL local dev ───────────────────────────────────────────────────────
start_postgres() {
    log_section "Starting Local Dev — PostgreSQL"
    cd "${PROJECT_ROOT}"

    # Source .env if it exists so variables like POSTGRES_PASSWORD are available
    [[ -f .env ]] && set -a && source .env && set +a

    local pg_user="${POSTGRES_USER:-javinizer}"
    local pg_pass="${POSTGRES_PASSWORD:-changeme}"
    local pg_db="${POSTGRES_DB:-javinizer}"
    local pg_port="${POSTGRES_HOST_PORT:-5432}"

    log_info "Starting PostgreSQL container..."
    docker compose --profile postgres up -d postgres

    log_info "Waiting for PostgreSQL to be ready..."
    local retries=30
    until docker compose exec -T postgres \
            pg_isready -U "${pg_user}" -d "${pg_db}" >/dev/null 2>&1; do
        retries=$((retries - 1))
        [[ "${retries}" -le 0 ]] && { log_error "PostgreSQL did not become ready in time."; exit 1; }
        sleep 1
    done
    log_info "PostgreSQL is ready ✓"

    mkdir -p data

    if [[ ! -f data/config.yaml ]]; then
        log_info "Copying default config → data/config.yaml"
        cp configs/config.yaml.example data/config.yaml 2>/dev/null || true
    fi

    export JAVINIZER_CONFIG="${PROJECT_ROOT}/data/config.yaml"
    export JAVINIZER_DB_TYPE="postgres"
    export JAVINIZER_DB_DSN="host=localhost user=${pg_user} password=${pg_pass} dbname=${pg_db} port=${pg_port} sslmode=disable TimeZone=UTC"

    log_info "DSN:    host=localhost user=${pg_user} dbname=${pg_db} port=${pg_port} sslmode=disable"
    log_info ""
    log_info "Starting API server (Ctrl+C to stop — Postgres container keeps running)..."
    log_info "Web UI:  http://localhost:8080"
    echo ""

    if command -v air >/dev/null 2>&1; then
        exec air -c .air.toml
    else
        exec go run ./cmd/javinizer api
    fi
}

# ── MySQL local dev ────────────────────────────────────────────────────────────
start_mysql() {
    log_section "Starting Local Dev — MySQL"
    cd "${PROJECT_ROOT}"

    [[ -f .env ]] && set -a && source .env && set +a

    local my_user="${MYSQL_USER:-javinizer}"
    local my_pass="${MYSQL_PASSWORD:-changeme}"
    local my_db="${MYSQL_DB:-javinizer}"
    local my_port="${MYSQL_HOST_PORT:-3306}"

    log_info "Starting MySQL container..."
    docker compose --profile mysql up -d mysql

    log_info "Waiting for MySQL to be ready..."
    local retries=40
    until docker compose exec -T mysql \
            mysqladmin ping -h localhost -u "${my_user}" "--password=${my_pass}" --silent >/dev/null 2>&1; do
        retries=$((retries - 1))
        [[ "${retries}" -le 0 ]] && { log_error "MySQL did not become ready in time."; exit 1; }
        sleep 2
    done
    log_info "MySQL is ready ✓"

    mkdir -p data

    if [[ ! -f data/config.yaml ]]; then
        log_info "Copying default config → data/config.yaml"
        cp configs/config.yaml.example data/config.yaml 2>/dev/null || true
    fi

    export JAVINIZER_CONFIG="${PROJECT_ROOT}/data/config.yaml"
    export JAVINIZER_DB_TYPE="mysql"
    export JAVINIZER_DB_DSN="${my_user}:${my_pass}@tcp(localhost:${my_port})/${my_db}?parseTime=True&loc=UTC&charset=utf8mb4"

    log_info "DSN:    ${my_user}:***@tcp(localhost:${my_port})/${my_db}?parseTime=True&loc=UTC"
    log_info ""
    log_info "Starting API server (Ctrl+C to stop — MySQL container keeps running)..."
    log_info "Web UI:  http://localhost:8080"
    echo ""

    if command -v air >/dev/null 2>&1; then
        exec air -c .air.toml
    else
        exec go run ./cmd/javinizer api
    fi
}

# ── Stop ───────────────────────────────────────────────────────────────────────
cmd_stop() {
    log_section "Stopping Dev Services"
    cd "${PROJECT_ROOT}"
    docker compose --profile postgres --profile mysql --profile flaresolverr down
    log_info "All Compose services stopped."
}

# ── Status ─────────────────────────────────────────────────────────────────────
cmd_status() {
    log_section "Dev Environment Status"
    cd "${PROJECT_ROOT}"
    echo ""
    docker compose --profile postgres --profile mysql ps 2>/dev/null || true
    echo ""

    if curl -sf http://localhost:8080/health >/dev/null 2>&1; then
        log_info "API server: ${GREEN}running${NC} at http://localhost:8080"
    else
        log_warn "API server: not reachable (start with: ./scripts/dev-env.sh start)"
    fi
}

# ── Logs ───────────────────────────────────────────────────────────────────────
cmd_logs() {
    local svc="${1:-}"
    cd "${PROJECT_ROOT}"
    if [[ -n "${svc}" ]]; then
        docker compose --profile postgres --profile mysql logs -f "${svc}"
    else
        docker compose --profile postgres --profile mysql logs -f
    fi
}

# ── Reset ──────────────────────────────────────────────────────────────────────
cmd_reset() {
    local db="${1:-sqlite}"
    log_section "Resetting Dev Environment (db=${db})"
    cd "${PROJECT_ROOT}"

    log_warn "This will destroy all DB volumes and data/javinizer.db!"
    read -r -p "Are you sure? [y/N] " confirm
    [[ "${confirm}" =~ ^[Yy]$ ]] || { log_info "Aborted."; exit 0; }

    docker compose --profile postgres --profile mysql down -v 2>/dev/null || true
    rm -f data/javinizer.db
    log_info "Volumes and SQLite file removed."
    cmd_start "${db}"
}

# ── Start dispatcher ───────────────────────────────────────────────────────────
cmd_start() {
    local db="${1:-sqlite}"
    check_prerequisites

    case "${db}" in
        sqlite|"")  start_sqlite   ;;
        postgres|pg) start_postgres ;;
        mysql|mariadb) start_mysql  ;;
        *)
            log_error "Unknown database type: '${db}'. Valid: sqlite, postgres, mysql"
            exit 1
            ;;
    esac
}

# ── Help ───────────────────────────────────────────────────────────────────────
cmd_help() {
    sed -n '2,20p' "${BASH_SOURCE[0]}" | grep '^#' | sed 's/^# \?//'
}

# ── Main ───────────────────────────────────────────────────────────────────────
main() {
    local cmd="${1:-help}"
    shift || true

    case "${cmd}" in
        start)   cmd_start  "${1:-sqlite}" ;;
        stop)    cmd_stop   ;;
        status)  cmd_status ;;
        logs)    cmd_logs   "${1:-}" ;;
        reset)   cmd_reset  "${1:-sqlite}" ;;
        help|-h|--help) cmd_help ;;
        *)
            log_error "Unknown command: '${cmd}'"
            cmd_help
            exit 1
            ;;
    esac
}

main "$@"
