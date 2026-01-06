#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

usage() {
    echo "Usage: $0 <command>"
    echo ""
    echo "Commands:"
    echo "  start       Start all services"
    echo "  stop        Stop all services"
    echo "  restart     Restart all services"
    echo "  status      Show services status"
    echo "  logs        Show services logs"
    echo "  clean       Stop and remove all data"
    echo ""
    exit 1
}

check_docker() {
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}Error: Docker is not installed${NC}"
        exit 1
    fi
    if ! docker info &> /dev/null; then
        echo -e "${RED}Error: Docker daemon is not running${NC}"
        exit 1
    fi
}

start_services() {
    echo -e "${YELLOW}Starting development environment...${NC}"
    cd "$PROJECT_DIR"
    docker compose up -d
    echo ""
    echo -e "${GREEN}Services started!${NC}"
    echo ""
    echo "Service URLs:"
    echo "  MySQL:  localhost:3306 (root/root)"
    echo "  Redis:  localhost:6379"
    echo "  Nacos:  http://localhost:8848/nacos (nacos/nacos)"
    echo "  Apollo: http://localhost:8070 (apollo/admin)"
    echo ""
    echo -e "${YELLOW}Note: Apollo first startup may take 1-2 minutes to initialize database${NC}"
}

stop_services() {
    echo -e "${YELLOW}Stopping development environment...${NC}"
    cd "$PROJECT_DIR"
    docker compose down
    echo -e "${GREEN}Services stopped!${NC}"
}

restart_services() {
    stop_services
    start_services
}

show_status() {
    cd "$PROJECT_DIR"
    docker compose ps
}

show_logs() {
    cd "$PROJECT_DIR"
    docker compose logs -f
}

clean_all() {
    echo -e "${RED}Warning: This will remove all data!${NC}"
    read -p "Are you sure? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        cd "$PROJECT_DIR"
        docker compose down -v
        echo -e "${GREEN}All services and data removed!${NC}"
    else
        echo "Cancelled."
    fi
}

# Main
check_docker

case "${1:-}" in
    start)
        start_services
        ;;
    stop)
        stop_services
        ;;
    restart)
        restart_services
        ;;
    status)
        show_status
        ;;
    logs)
        show_logs
        ;;
    clean)
        clean_all
        ;;
    *)
        usage
        ;;
esac
