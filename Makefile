# MongoDB Unit 测试工具

.PHONY: help up down test test-multi test-all coverage clean

# 默认目标
help:
	@echo "MongoDB Unit 测试工具"
	@echo ""
	@echo "可用命令:"
	@echo "  make up          - 启动 MongoDB Docker 服务"
	@echo "  make down        - 停止 MongoDB Docker 服务"
	@echo "  make test        - 运行单实例测试 (默认)"
	@echo "  make test-multi  - 运行多实例测试"
	@echo "  make test-all    - 运行所有测试 (单实例 + 多实例)"
	@echo "  make coverage    - 生成测试覆盖率报告"
	@echo "  make clean       - 清理测试环境"

# 启动 MongoDB 服务
up:
	docker compose up -d
	@echo "等待 MongoDB 服务启动..."
	@docker compose ps

# 停止 MongoDB 服务
down:
	docker compose down

# 单实例测试
test: up
	go test -v ./...

# 多实例测试
test-multi: up
	TAO_TEST_MULTI_INSTANCE=true go test -v ./...

# 运行所有测试
test-all: up
	@echo "=== 运行单实例测试 ==="
	go test -v ./...
	@echo ""
	@echo "=== 运行多实例测试 ==="
	TAO_TEST_MULTI_INSTANCE=true go test -v ./...

# 生成覆盖率报告
coverage: up
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 清理
clean: down
	rm -f coverage.out coverage.html
