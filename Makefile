mocks:
# repository
	mockgen -source=internal/domain/repositories/products.go -destination=internal/domain/repositories/mocks/mock_products.go -package=mocks

test:
	go test ./...