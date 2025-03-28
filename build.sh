set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${YELLOW}==== Olaf Interpreter - Build Script ====${NC}"

if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go first.${NC}"
    exit 1
fi

echo -e "${GREEN}Creating directory structure...${NC}"
mkdir -p web/js

echo -e "${GREEN}Building WebAssembly module...${NC}"
GOOS=js GOARCH=wasm go build -o web/sup-bud.wasm cmd/sup-bud/main.go

echo -e "${GREEN}Copying WebAssembly support files...${NC}"
GO_ROOT=$(go env GOROOT)
cp "${GO_ROOT}/misc/wasm/wasm_exec.js" web/js/

echo -e "${GREEN}Setting up application JavaScript...${NC}"
if [ -f "app.js" ]; then
    cp app.js web/js/
fi

if [ -f "index.html" ]; then
    echo -e "${GREEN}Copying HTML files...${NC}"
    cp index.html web/
fi

echo -e "${GREEN}Building development server...${NC}"
go build -o sup-bud-server server.go

echo -e "${YELLOW}==== Build Complete ====${NC}"
echo -e "To start the server, run: ${GREEN}./web/sup-bud-server${NC}"
echo -e "Then open ${GREEN}http://localhost:8080${NC} in your browser"

./sup-bud-server