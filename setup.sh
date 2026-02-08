#!/bin/bash

# Assetra Security Setup Script
# This script helps you set up the project securely

set -e

echo "🔐 Assetra Security Setup"
echo "========================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

PROTO_SRC="messages/auth.proto"
PB_GO="pb/auth.pb.go"
PB_GRPC_GO="pb/auth_grpc.pb.go"
REGEN_PROTO=false

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Docker not found! Please install Docker: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if docker compose is available
if ! docker compose version &> /dev/null && ! command -v docker-compose &> /dev/null; then
    echo "docker compose not found! Please update Docker or install docker-compose: https://docs.docker.com/compose/install/"
    exit 1
fi

echo ""
echo "🔄 Checking if protobuf Go files need to be (re)generated..."

if [ ! -f "$PB_GO" ] || [ ! -f "$PB_GRPC_GO" ]; then
    echo -e "${YELLOW}Protobuf Go files missing. Will generate...${NC}"
    REGEN_PROTO=true
elif [ "$PROTO_SRC" -nt "$PB_GO" ] || [ "$PROTO_SRC" -nt "$PB_GRPC_GO" ]; then
    echo -e "${YELLOW}Protobuf source is newer than generated files. Will regenerate...${NC}"
    REGEN_PROTO=true
fi

if [ "$REGEN_PROTO" = true ]; then
    if protoc --go_out=pb --go-grpc_out=pb --proto_path=messages messages/auth.proto; then
        echo -e "${GREEN}✓ Protobuf Go files generated in pb/${NC}"
    else
        echo -e "${RED}✗ Failed to generate protobuf Go files. Check protoc and protoc-gen-go installation.${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}✓ Protobuf Go files are up to date${NC}"
fi

echo "📝 Creating .env file from template..."

# Backend .env
if [ -f .env ]; then
    echo -e "${YELLOW}⚠️  .env file already exists!${NC}"
    read -p "Do you want to overwrite it? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Keeping existing .env file"
    else
        echo "📝 Creating .env file from template..."
        cp .env.example .env
    fi
else
    echo "📝 Creating .env file from template..."
    cp .env.example .env
fi

# Frontend .env
if [ -f frontend/.env ]; then
    echo -e "${YELLOW}⚠️  frontend/.env file already exists!${NC}"
    read -p "Do you want to overwrite frontend/.env? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "📝 Creating frontend/.env from template..."
        cp frontend/.env.example frontend/.env
    else
        echo "Keeping existing frontend/.env file"
    fi
else
    if [ -f frontend/.env.example ]; then
        echo "📝 Creating frontend/.env from template..."
        cp frontend/.env.example frontend/.env
    fi
fi

# Generate secure SECURITY_KEY
echo ""
echo "🔑 Generating secure SECURITY_KEY..."
if command -v openssl &> /dev/null; then
    SECURITY_KEY=$(openssl rand -base64 64 | tr -d '\n')
    echo -e "${GREEN}✓ Generated 64-character random key${NC}"
else
    echo -e "${RED}✗ openssl not found. Please install openssl or generate a key manually${NC}"
    echo "  You can use: openssl rand -base64 64"
    exit 1
fi

# Update .env with generated key
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s|SECURITY_KEY=.*|SECURITY_KEY=$SECURITY_KEY|" .env
else
    # Linux
    sed -i "s|SECURITY_KEY=.*|SECURITY_KEY=$SECURITY_KEY|" .env
fi

echo ""
echo "📋 Please provide the following information:"
echo ""

# Database configuration
read -p "Database host (default: localhost): " DB_HOST
DB_HOST=${DB_HOST:-localhost}

read -p "Database port (default: 5432): " DB_PORT
DB_PORT=${DB_PORT:-5432}

read -p "Database name (default: assetra): " DB_NAME
DB_NAME=${DB_NAME:-assetra}

read -p "Database user (default: postgres): " DB_USER
DB_USER=${DB_USER:-postgres}

echo ""
read -sp "Database password: " DB_PASSWORD
echo ""

# Update .env with database config
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s|DB_HOST=.*|DB_HOST=$DB_HOST|" .env
    sed -i '' "s|DB_PORT=.*|DB_PORT=$DB_PORT|" .env
    sed -i '' "s|DB_NAME=.*|DB_NAME=$DB_NAME|" .env
    sed -i '' "s|DB_USER=.*|DB_USER=$DB_USER|" .env
    sed -i '' "s|DB_PASSWORD=.*|DB_PASSWORD=$DB_PASSWORD|" .env
else
    sed -i "s|DB_HOST=.*|DB_HOST=$DB_HOST|" .env
    sed -i "s|DB_PORT=.*|DB_PORT=$DB_PORT|" .env
    sed -i "s|DB_NAME=.*|DB_NAME=$DB_NAME|" .env
    sed -i "s|DB_USER=.*|DB_USER=$DB_USER|" .env
    sed -i "s|DB_PASSWORD=.*|DB_PASSWORD=$DB_PASSWORD|" .env
fi

echo ""
echo -e "${GREEN}✓ .env file created successfully!${NC}"
echo ""

# Install/Update Go dependencies
echo "📦 Installing Go dependencies..."
if go mod download && go mod tidy; then
    echo -e "${GREEN}✓ Go dependencies installed${NC}"
else
    echo -e "${RED}✗ Failed to install Go dependencies${NC}"
    exit 1
fi

# Check if npm is available for frontend
if command -v npm &> /dev/null; then
    echo ""
    read -p "Do you want to install frontend dependencies? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        cd frontend
        echo "📦 Installing frontend dependencies..."
        if npm install; then
            echo -e "${GREEN}✓ Frontend dependencies installed${NC}"
        else
            echo -e "${RED}✗ Failed to install frontend dependencies${NC}"
        fi
        cd ..
    fi
fi

echo ""
echo "🎉 Setup complete!"
echo ""
echo "Next steps:"
echo "  1. Review your .env file: nano .env"
echo "  2. Start: make build-start"
echo ""
echo -e "${YELLOW}⚠️  IMPORTANT: Never commit your .env file to git!${NC}"
