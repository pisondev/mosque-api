# URL koneksi ke database. Perhatikan port 5435 sesuai docker-compose milikmu.
DB_URL=postgres://root:secretpassword@localhost:5435/mosque_saas?sslmode=disable

.PHONY: migrate-up migrate-down migrate-force new-migration run

# Menjalankan migrasi ke atas (apply schema)
migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

# Menurunkan 1 migrasi terakhir (rollback)
migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down 1

# Memaksa versi migrasi jika terjadi dirty state (misal: make migrate-force v=1)
migrate-force:
	migrate -path db/migrations -database "$(DB_URL)" force $(v)

# Membuat file migrasi baru (misal: make new-migration name=add_billing_table)
new-migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

# Menjalankan server Go lokal
run:
	go run main.go