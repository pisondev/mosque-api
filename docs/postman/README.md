# Postman Workspace

Kumpulan Postman untuk `mosque-api` disusun per module agar import, review, dan maintenance tidak menumpuk dalam satu collection besar.

## Struktur

- `environments/etakmir-local.postman_environment.json`: environment lokal bersama.
- `collections/00-system.postman_collection.json`: health check dasar server.
- `collections/01-auth.postman_collection.json`: module auth sebagai entry point untuk mendapat `bearerToken`.

## Urutan Pakai

1. Import environment `environments/etakmir-local.postman_environment.json`.
2. Jalankan `collections/00-system.postman_collection.json` untuk verifikasi server hidup.
3. Jalankan `collections/01-auth.postman_collection.json`.
4. Request `Register Admin` otomatis mengisi `registerEmail`, `email`, dan `bearerToken` jika berhasil.
5. Request `Login Email Password` mengisi ulang `bearerToken` dari response login.

## Catatan Auth

- `Forgot Password` hanya akan benar-benar mengirim email jika `RESEND_API_KEY` aktif di backend.
- `Reset Password` membutuhkan `resetToken` valid dari email reset atau sumber dev lain.
- Secara default `autoGenerateRegisterEmail=true`, jadi request register aman dijalankan berulang tanpa bentrok email lama.

## Rencana Module Berikutnya

- `02-management`
- `03-worship`
- `04-community`
- `05-engagement`
- `06-finance`
- `07-public`
- `08-webhook`
