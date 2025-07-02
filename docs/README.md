# API Documentation

The OpenAPI (Swagger) spec for Reservio is in `swagger.yaml`.

## How to View

### Option 1: Swagger Editor (Web)
- Go to https://editor.swagger.io/
- Click "File" > "Import File" and select `swagger.yaml` from this repo.

### Option 2: Swagger UI Docker
- Run Swagger UI locally with Docker:
  ```sh
  docker run -p 8081:8080 -v $(pwd)/swagger.yaml:/usr/share/nginx/html/swagger.yaml swaggerapi/swagger-ui
  ```
- Open http://localhost:8081 and enter `/swagger.yaml` as the spec URL.

---

For updates, edit `docs/swagger.yaml` and re-import as needed.

## Automated Swagger Generation

- You can generate OpenAPI docs from Go code comments using [swaggo/swag](https://github.com/swaggo/swag):
  ```sh
  go install github.com/swaggo/swag/cmd/swag@latest
  swag init -g ../../cmd/main.go -o ./generated
  ```
- This will create `docs/generated/swagger.yaml` from your Go code.
- You can view or merge this with the main `swagger.yaml` as needed. 