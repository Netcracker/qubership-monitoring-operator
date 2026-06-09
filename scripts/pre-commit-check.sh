#!/bin/bash
# Pre-commit validation script for qubership-monitoring-operator
# TODO File is temporary cretaed. Content of file should be moved to Makefile. 

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

log_section() {
    echo -e "\n${BLUE}🔍 $1${NC}"
    echo "═══════════════════════════════════════════════════════"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Environment checks
check_environment() {
    log_section "Environment Validation"
    
    # Check Go version
    if command_exists go; then
        GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
        log_info "Go version: $GO_VERSION"
        if [[ "$GO_VERSION" < "1.23" ]]; then
            log_error "Go version must be >= 1.23"
            exit 1
        fi
        log_success "Go version is compatible"
    else
        log_error "Go is not installed"
        exit 1
    fi
    
    # Check Docker
    if command_exists docker; then
        DOCKER_VERSION=$(docker --version | cut -d' ' -f3 | sed 's/,//')
        log_info "Docker version: $DOCKER_VERSION"
        log_success "Docker is available"
    else
        log_warning "Docker is not available - skipping Docker builds"
    fi
    
    # Check Python3
    if command_exists python3; then
        PYTHON_VERSION=$(python3 --version | cut -d' ' -f2)
        log_info "Python version: $PYTHON_VERSION"
        log_success "Python3 is available"
    else
        log_error "Python3 is not installed"
        exit 1
    fi
}

# Go code validation
go_validation() {
    log_section "Go Code Validation"
    
    log_info "Generating CRDs and deepcopy files..."
    make generate
    log_success "Code generation completed"
    
    log_info "Running go fmt..."
    make fmt
    log_success "Go formatting completed"
    
    log_info "Running go vet..."
    make vet
    log_success "Static analysis completed"
    
    log_info "Running unit tests..."
    make test
    log_success "Unit tests passed"
    
    log_info "Building binary..."
    make build-binary
    log_success "Binary build completed"
}

# Documentation validation
docs_validation() {
    log_section "Documentation Validation"
    
    log_info "Copying CRDs to documentation..."
    make docs
    log_success "CRDs copied to docs"
    
    log_info "Installing MkDocs dependencies..."
    if [[ -f "requirements_mkdocs.txt" ]]; then
        pip install -r requirements_mkdocs.txt > /dev/null 2>&1 || {
            log_warning "Failed to install requirements, trying with --user"
            pip install --user -r requirements_mkdocs.txt > /dev/null 2>&1
        }
        log_success "MkDocs dependencies installed"
    elif [[ -f "requirements.txt" ]]; then
        pip install -r requirements.txt > /dev/null 2>&1 || {
            log_warning "Failed to install requirements, trying with --user"
            pip install --user -r requirements.txt > /dev/null 2>&1
        }
        log_success "MkDocs dependencies installed"
    else
        log_warning "requirements_mkdocs.txt or requirements.txt not found, skipping pip install"
    fi
    
    if [[ "${ALLOW_DOCS_WARNINGS:-}" == "true" ]]; then
        log_info "Building MkDocs site (warnings allowed)..."
        if mkdocs build --quiet; then
            log_success "MkDocs site built successfully"
        else
            log_warning "MkDocs build completed with warnings"
        fi
        
        log_info "Testing MkDocs in strict mode (non-fatal)..."
        if mkdocs build --strict --quiet; then
            log_success "MkDocs strict validation passed"
        else
            log_warning "MkDocs has warnings - check with: mkdocs build --strict --verbose"
        fi
    else
        log_info "Building MkDocs site in strict mode..."
        if mkdocs build --strict > /dev/null 2>&1; then
            log_success "MkDocs strict validation passed - no warnings found"
        else
            log_error "MkDocs validation failed - broken links or warnings detected"
            echo -e "\n${RED}Detailed error output:${NC}"
            mkdocs build --strict --verbose 2>&1 | grep -E "(WARNING|ERROR)" || true
            echo -e "\n${RED}Fix all warnings before committing!${NC}"
            echo -e "${YELLOW}Tip: Use --no-strict-docs to allow warnings (not recommended)${NC}"
            return 1
        fi
    fi
}

# Docker validation (optional)
docker_validation() {
    if ! command_exists docker; then
        log_warning "Docker not available, skipping Docker validation"
        return 0
    fi
    
    log_section "Docker Validation"
    
    log_info "Building main operator image..."
    if docker build -t qubership-monitoring-operator -f Dockerfile . > /dev/null 2>&1; then
        log_success "Main operator image built"
    else
        log_error "Failed to build main operator image"
        return 1
    fi
    
    # Optional: Build other images if time permits
    log_info "Building transfer image..."
    if timeout 60 docker build -t qubership-monitoring-transfer -f docker-transfer/Dockerfile . > /dev/null 2>&1; then
        log_success "Transfer image built"
    else
        log_warning "Transfer image build skipped (timeout or error)"
    fi
}

# Helm chart validation
helm_validation() {
    log_section "Helm Chart Validation"
    
    if command_exists helm; then
        log_info "Validating Helm chart..."
        if helm lint charts/qubership-monitoring-operator/ > /dev/null 2>&1; then
            log_success "Helm chart validation passed"
        else
            log_warning "Helm chart has issues"
        fi
    else
        log_warning "Helm not available, skipping chart validation"
    fi
}

# Git status check
git_status_check() {
    log_section "Git Status Check"
    
    if git diff --quiet && git diff --staged --quiet; then
        log_success "Working directory is clean"
    else
        log_info "Changes detected:"
        git status --short
        log_warning "You have uncommitted changes"
    fi
    
    # Check if we're on a feature branch
    CURRENT_BRANCH=$(git branch --show-current)
    if [[ "$CURRENT_BRANCH" == "main" || "$CURRENT_BRANCH" == "master" ]]; then
        log_warning "You're on the main branch. Consider creating a feature branch."
    else
        log_success "Working on feature branch: $CURRENT_BRANCH"
    fi
}

# Summary and next steps
print_summary() {
    log_section "Pre-commit Check Summary"
    
    echo "✅ Environment validation"
    echo "✅ Go code validation"
    echo "✅ Documentation validation"
    echo "✅ Docker validation (if available)"
    echo "✅ Helm validation (if available)"
    echo "✅ Git status check"
    
    echo -e "\n${GREEN}🎉 Pre-commit validation completed successfully!${NC}"
    echo -e "\n${BLUE}Next steps:${NC}"
    echo "1. Review any warnings above"
    echo "2. Commit your changes: git add . && git commit -m 'Your message'"
    echo "3. Push to your fork: git push origin your-branch"
    echo "4. Create a Pull Request"
    
    if [[ -d "site" ]]; then
        echo -e "\n${BLUE}📚 Documentation preview available:${NC}"
        echo "   Run: mkdocs serve"
        echo "   Visit: http://127.0.0.1:8000"
    fi
}

# Main execution
main() {
    echo -e "${BLUE}🚀 Qubership Monitoring Operator - Pre-commit Validation${NC}"
    echo "═══════════════════════════════════════════════════════════════"
    
    # Change to project root if script is run from scripts/ directory
    if [[ "$(basename "$(pwd)")" == "scripts" ]]; then
        cd ..
    fi
    
    # Validate we're in the right directory
    if [[ ! -f "go.mod" ]] || [[ ! -f "mkdocs.yml" ]]; then
        log_error "This script must be run from the project root directory"
        exit 1
    fi
    
    # Run all validation steps
    check_environment || exit 1
    go_validation || exit 1
    docs_validation || exit 1
    docker_validation || exit 1
    helm_validation || exit 1
    git_status_check
    print_summary
}

# Command line options
case "${1:-}" in
    --help|-h)
        echo "Usage: $0 [options]"
        echo "Options:"
        echo "  --help, -h       Show this help message"
        echo "  --no-docker      Skip Docker validation"
        echo "  --docs-only      Run only documentation validation"
        echo "  --go-only        Run only Go validation"
        echo "  --no-strict-docs Allow MkDocs warnings (not recommended)"
        exit 0
        ;;
    --no-docker)
        SKIP_DOCKER=true
        ;;
    --no-strict-docs)
        ALLOW_DOCS_WARNINGS=true
        ;;
    --docs-only)
        docs_validation
        exit 0
        ;;
    --go-only)
        go_validation
        exit 0
        ;;
esac

# Run main if no special options
if [[ "${SKIP_DOCKER:-}" == "true" ]]; then
    docker_validation() { log_info "Docker validation skipped by user request"; }
fi

main "$@" 