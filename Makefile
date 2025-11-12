.PHONY: generate generate-types generate-server clean-gen

generate: clean-gen generate-types generate-server

generate-types:
	@mkdir -p internal/api/generated
	oapi-codegen \
		-package generated \
		-generate types \
		-o internal/api/generated/types.gen.go \
		spec/openapi.yml

generate-server:
	oapi-codegen \
		-package generated \
		-generate server,strict-server \
		-o internal/api/generated/server.gen.go \
		spec/openapi.yml

clean-gen:
	@rm -rf internal/api/generated/*.gen.go
