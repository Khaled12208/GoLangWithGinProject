#!/usr/bin/env bash

# Colors for output
if [[ -t 1 ]]; then
    GREEN='\033[0;32m'
    BLUE='\033[0;34m'
    RED='\033[0;31m'
    YELLOW='\033[1;33m'
    NC='\033[0m' # No Color
else
    GREEN=''
    BLUE=''
    RED=''
    YELLOW=''
    NC=''
fi

# Function to display KHALED ASCII art
display_khaled_ascii() {
    echo -e "${BLUE}"
    echo "██╗  ██╗██╗  ██╗ █████╗ ██╗     ███████╗██████╗"
    echo "██║ ██╔╝██║  ██║██╔══██╗██║     ██╔════╝██╔══██╗"
    echo "█████╔╝ ███████║███████║██║     █████╗  ██║  ██║"
    echo "██╔═██╗ ██╔══██║██╔══██║██║     ██╔══╝  ██║  ██║"
    echo "██║  ██╗██║  ██║██║  ██║███████╗███████╗██████╔╝"
    echo "╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚══════╝╚═════╝"
    echo -e "${NC}"
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin*)
            echo "macos"
            ;;
        Linux*)
            echo "linux"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            echo "windows"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

OS=$(detect_os)

# Function to display usage
show_usage() {
    echo -e "${BLUE}Usage: ./run.sh [command]${NC}"
    echo "Commands:"
    echo "  start     - Start the application with Swagger UI"
    echo "  local     - Run the project locally"
    echo "  docker    - Run the project in Docker containers"
    echo "  stop      - Stop all project containers"
    echo "  clean     - Clean up all Docker resources"
    echo "  swagger   - Generate Swagger documentation"
    echo "  test      - Run tests (use -v for verbose, -race for race detection)"
    echo "  coverage  - Run tests with coverage report"
    echo "  integration - Run integration tests"
    echo "  test-all  - Run all tests (unit, integration, race detection)"
    echo "  help      - Show this help message"
}

# Function to check if a port is in use
is_port_in_use() {
    local port=$1
    case $OS in
        "windows")
            netstat -ano | grep "LISTENING" | grep ":$port" > /dev/null 2>&1
            ;;
        "macos"|"linux")
            command -v lsof >/dev/null 2>&1 && lsof -i:$port > /dev/null 2>&1
            ;;
        *)
            echo -e "${RED}Unsupported operating system${NC}"
            exit 1
            ;;
    esac
    return $?
}

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        echo -e "${RED}Docker is not running. Please start Docker and try again.${NC}"
        exit 1
    fi
}

# Function to check required tools
check_requirements() {
    local missing_tools=()

    # Check Go
    if ! command -v go >/dev/null 2>&1; then
        missing_tools+=("Go")
    fi

    # Check Docker
    if ! command -v docker >/dev/null 2>&1; then
        missing_tools+=("Docker")
    fi

    # Check Docker Compose
    if ! command -v docker-compose >/dev/null 2>&1; then
        missing_tools+=("Docker Compose")
    fi

    # Check swag
    if ! command -v swag >/dev/null 2>&1; then
        echo -e "${BLUE}Installing swag...${NC}"
        go install github.com/swaggo/swag/cmd/swag@latest
    fi

    if [ ${#missing_tools[@]} -ne 0 ]; then
        echo -e "${RED}Missing required tools: ${missing_tools[*]}${NC}"
        echo "Please install them and try again."
        exit 1
    fi
}

# Function to wait for MySQL
wait_for_mysql() {
    local retries=30
    local wait_time=2
    echo -e "${BLUE}Waiting for MySQL to be ready...${NC}"
    
    for ((i=1; i<=$retries; i++)); do
        if docker exec mysql-local mysqladmin ping -h"localhost" -P"3306" -u"root" -p"mysecretpassword" --silent >/dev/null 2>&1; then
            echo -e "${GREEN}MySQL is ready!${NC}"
            return 0
        fi
        echo -n "."
        sleep $wait_time
    done
    
    echo -e "\n${RED}MySQL did not become ready in time${NC}"
    return 1
}

# Function to run unit tests
run_tests() {
    echo -e "${BLUE}Running unit tests...${NC}"
    local test_flags="$1"
    if [ -z "$test_flags" ]; then
        test_flags="-v"
    fi
    
    go test $test_flags ./...
    local test_exit_code=$?
    
    if [ $test_exit_code -eq 0 ]; then
        echo -e "${GREEN}All tests passed!${NC}"
    else
        echo -e "${RED}Tests failed!${NC}"
        return $test_exit_code
    fi
}

# Function to run tests with coverage
run_coverage() {
    echo -e "${BLUE}Running tests with coverage...${NC}"
    
    # Create coverage output directory if it doesn't exist
    mkdir -p coverage
    
    # Run tests with coverage
    go test -coverprofile=coverage/coverage.out ./...
    local test_exit_code=$?
    
    if [ $test_exit_code -eq 0 ]; then
        # Generate coverage report
        go tool cover -html=coverage/coverage.out -o coverage/coverage.html
        echo -e "${GREEN}Coverage report generated: coverage/coverage.html${NC}"
        
        # Show coverage percentage
        local coverage_percent=$(go tool cover -func=coverage/coverage.out | grep total | awk '{print $3}')
        echo -e "${BLUE}Total coverage: ${YELLOW}$coverage_percent${NC}"
    else
        echo -e "${RED}Tests failed!${NC}"
        return $test_exit_code
    fi
}

# Function to run integration tests
run_integration_tests() {
    echo -e "${BLUE}Running integration tests...${NC}"
    
    # Start test database
    echo -e "${BLUE}Starting test database...${NC}"
    docker-compose -f docker-compose.test.yml up -d mysql-test
    
    # Wait for test database to be ready
    sleep 5
    
    # Run integration tests
    go test -tags=integration -v ./...
    local test_exit_code=$?
    
    # Clean up test database
    docker-compose -f docker-compose.test.yml down
    
    if [ $test_exit_code -eq 0 ]; then
        echo -e "${GREEN}All integration tests passed!${NC}"
    else
        echo -e "${RED}Integration tests failed!${NC}"
        return $test_exit_code
    fi
}

# Function to run all tests
run_all_tests() {
    echo -e "${BLUE}Running all tests...${NC}"
    
    # Run unit tests with race detection
    echo -e "\n${YELLOW}Running unit tests with race detection...${NC}"
    run_tests "-v -race"
    local unit_exit_code=$?
    
    # Run integration tests
    echo -e "\n${YELLOW}Running integration tests...${NC}"
    run_integration_tests
    local integration_exit_code=$?
    
    # Run coverage
    echo -e "\n${YELLOW}Generating coverage report...${NC}"
    run_coverage
    local coverage_exit_code=$?
    
    # Check if any tests failed
    if [ $unit_exit_code -eq 0 ] && [ $integration_exit_code -eq 0 ] && [ $coverage_exit_code -eq 0 ]; then
        echo -e "\n${GREEN}All tests passed successfully!${NC}"
        return 0
    else
        echo -e "\n${RED}Some tests failed!${NC}"
        return 1
    fi
}

# Function to run locally
run_local() {
    echo -e "${BLUE}Starting local development...${NC}"
    display_khaled_ascii
    check_requirements
    
    # Check if ports are available
    if is_port_in_use 8888; then
        echo -e "${RED}Port 8888 is already in use. Please free it up first.${NC}"
        exit 1
    fi
    
    if is_port_in_use 3306; then
        echo -e "${RED}Port 3306 is already in use. Please free it up first.${NC}"
        exit 1
    fi

    # Start MySQL in Docker
    echo -e "${BLUE}Starting MySQL container...${NC}"
    check_docker
    docker run --name mysql-local \
        -e MYSQL_ROOT_PASSWORD=mysecretpassword \
        -e MYSQL_DATABASE=golangwithgin \
        -p 3306:3306 \
        -d mysql:8.0 \
        --default-authentication-plugin=mysql_native_password

    # Wait for MySQL
    wait_for_mysql

    # Initialize database
    echo -e "${BLUE}Initializing database...${NC}"
    docker cp scripts/init.sql mysql-local:/init.sql
    docker exec mysql-local mysql -uroot -pmysecretpassword golangwithgin < scripts/init.sql

    # Generate Swagger docs
    echo -e "${BLUE}Generating Swagger documentation...${NC}"
    swag init -g cmd/api/main.go

    # Set environment variables based on OS
    case $OS in
        "windows")
            export GIN_MODE=release
            ;;
        "macos"|"linux")
            export GIN_MODE=release
            ;;
    esac

    # Run the application
    echo -e "${GREEN}Starting the application...${NC}"
    echo -e "${BLUE}Swagger UI will be available at: http://localhost:8888/swagger/index.html${NC}"
    go run cmd/api/main.go
}

# Function to run in Docker
run_docker() {
    echo -e "${BLUE}Starting Docker deployment...${NC}"
    display_khaled_ascii
    check_requirements
    
    # Check if ports are available
    if is_port_in_use 8888; then
        echo -e "${RED}Port 8888 is already in use. Please free it up first.${NC}"
        exit 1
    fi
    
    if is_port_in_use 3306; then
        echo -e "${RED}Port 3306 is already in use. Please free it up first.${NC}"
        exit 1
    fi

    # Generate Swagger docs
    echo -e "${BLUE}Generating Swagger documentation...${NC}"
    swag init -g cmd/api/main.go

    # Start containers
    echo -e "${BLUE}Starting containers...${NC}"
    echo -e "${BLUE}Swagger UI will be available at: http://localhost:8888/swagger/index.html${NC}"
    docker-compose up --build
}

# Function to stop containers
stop_containers() {
    echo -e "${BLUE}Stopping containers...${NC}"
    check_docker
    docker-compose down
    docker stop mysql-local 2>/dev/null || true
    docker rm mysql-local 2>/dev/null || true
    echo -e "${GREEN}Containers stopped${NC}"
}

# Function to clean Docker resources
clean_docker() {
    echo -e "${BLUE}Cleaning Docker resources...${NC}"
    check_docker
    stop_containers
    docker volume prune -f
    docker network prune -f
    echo -e "${GREEN}Docker resources cleaned${NC}"
}

# Function to generate Swagger docs
generate_swagger() {
    echo -e "${BLUE}Generating Swagger documentation...${NC}"
    check_requirements
    swag init -g cmd/api/main.go
    echo -e "${GREEN}Swagger documentation generated${NC}"
}

# Function to kill process using a specific port
kill_port_process() {
    local port=$1
    case $OS in
        "windows")
            local pid=$(netstat -ano | grep ":$port" | grep "LISTENING" | awk '{print $5}')
            if [ ! -z "$pid" ]; then
                taskkill //PID $pid //F >/dev/null 2>&1
            fi
            ;;
        "macos")
            local pid=$(lsof -i:$port -t)
            if [ ! -z "$pid" ]; then
                kill -9 $pid >/dev/null 2>&1
            fi
            ;;
        "linux")
            local pid=$(lsof -i:$port -t)
            if [ ! -z "$pid" ]; then
                kill -9 $pid >/dev/null 2>&1
            fi
            ;;
    esac
}

# Function to start the application with Swagger
start_application() {
    echo -e "${BLUE}Starting the application with Swagger UI...${NC}"
    display_khaled_ascii
    
    # Kill any process using port 8888
    if is_port_in_use 8888; then
        echo -e "${YELLOW}Port 8888 is in use. Stopping the process...${NC}"
        kill_port_process 8888
        sleep 2
    fi
    
    # Generate Swagger documentation
    echo -e "${BLUE}Generating Swagger documentation...${NC}"
    swag init -g cmd/api/main.go
    
    # Stop any existing MySQL container and network
    echo -e "${BLUE}Cleaning up existing MySQL resources...${NC}"
    docker stop mysql-local >/dev/null 2>&1 || true
    docker rm mysql-local >/dev/null 2>&1 || true
    docker network rm mysql-network >/dev/null 2>&1 || true
    
    # Create a Docker network for MySQL
    echo -e "${BLUE}Creating MySQL network...${NC}"
    docker network create mysql-network >/dev/null 2>&1 || true
    
    # Start MySQL container
    echo -e "${BLUE}Starting MySQL container...${NC}"
    docker run --name mysql-local \
        --network mysql-network \
        -e MYSQL_ROOT_PASSWORD=mysecretpassword \
        -e MYSQL_DATABASE=golangwithgin \
        -p 3306:3306 \
        -d mysql:8.0 \
        --default-authentication-plugin=mysql_native_password \
        --character-set-server=utf8mb4 \
        --collation-server=utf8mb4_unicode_ci
        
    # Wait for MySQL to be ready
    wait_for_mysql
    
    # Initialize the database
    echo -e "${BLUE}Initializing database...${NC}"
    docker exec mysql-local mysql -uroot -pmysecretpassword -e "ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'mysecretpassword';" >/dev/null 2>&1
    docker exec mysql-local mysql -uroot -pmysecretpassword -e "FLUSH PRIVILEGES;" >/dev/null 2>&1
    
    # Start the application
    echo -e "${GREEN}Starting the application...${NC}"
    echo -e "${BLUE}Swagger UI will be available at: ${GREEN}http://localhost:8888/swagger/index.html${NC}"
    
    # Export database configuration
    export DB_HOST=127.0.0.1
    export DB_PORT=3306
    export DB_USER=root
    export DB_PASSWORD=mysecretpassword
    export DB_NAME=golangwithgin
    export SERVER_PORT=8888
    
    # Add a small delay to ensure MySQL is fully ready
    sleep 5
    
    go run cmd/api/main.go
}

# Main script logic
case "$1" in
    "start")
        check_requirements
        start_application
        ;;
    "local")
        check_requirements
        run_local
        ;;
    "docker")
        run_docker
        ;;
    "stop")
        stop_containers
        ;;
    "clean")
        clean_docker
        ;;
    "swagger")
        generate_swagger
        ;;
    "test")
        shift
        run_tests "$*"
        ;;
    "coverage")
        run_coverage
        ;;
    "integration")
        run_integration_tests
        ;;
    "test-all")
        run_all_tests
        ;;
    "help"|"")
        show_usage
        ;;
    *)
        echo -e "${RED}Unknown command: $1${NC}"
        show_usage
        exit 1
        ;;
esac 