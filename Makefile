# DATABASE_URL="host=localhost port=5432 user=postgres password=postgres dbname=ktaxes sslmode=disable"
DATABASE_URL=postgres://postgres:postgres@localhost:5432/ktaxes?sslmode=disable
PORT=8080
ADMIN_USERNAME=adminTax
ADMIN_PASSWORD=admin!

## How to run

DATABASE_URL=postgres://postgres:postgres@localhost:5432/ktaxes?sslmode=disable PORT=8080 ADMIN_USERNAME=adminTax ADMIN_PASSWORD=admin! go run main.go