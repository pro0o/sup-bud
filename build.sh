#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}==== Olaf Interpreter - Build Script ====${NC}"

# Check Go installation
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go first.${NC}"
    exit 1
fi

# Ensure proper directory structure
echo -e "${GREEN}Creating directory structure...${NC}"
mkdir -p web/js
mkdir -p web/css

# Build the WASM module
echo -e "${GREEN}Building WebAssembly module...${NC}"
GOOS=js GOARCH=wasm go build -o web/olaf.wasm main.go

# Copy the JavaScript glue code
echo -e "${GREEN}Copying WebAssembly support files...${NC}"
GO_ROOT=$(go env GOROOT)
cp "${GO_ROOT}/misc/wasm/wasm_exec.js" web/js/

# Copy our app.js
echo -e "${GREEN}Setting up application JavaScript...${NC}"
# This assumes app.js is in the current directory and will be moved to web/js
# If app.js is already in the right place, you can remove or adjust this line
if [ -f "app.js" ]; then
    cp app.js web/js/
fi

# Copy index.html if needed
if [ -f "index.html" ]; then
    echo -e "${GREEN}Copying HTML files...${NC}"
    cp index.html web/
fi

# Build server
echo -e "${GREEN}Building development server...${NC}"
go build -o olaf-server server.go

echo -e "${YELLOW}==== Build Complete ====${NC}"
echo -e "To start the server, run: ${GREEN}./olaf-server${NC}"
echo -e "Then open ${GREEN}http://localhost:8080${NC} in your browser"