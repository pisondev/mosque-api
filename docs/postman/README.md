# Postman Workspace

Kumpulan Postman untuk `mosque-api` disusun per module agar import, review, dan maintenance tidak menumpuk dalam satu collection besar.

## Struktur

- `environments/etakmir-local.postman_environment.json`: environment lokal bersama.
- `collections/00-system.postman_collection.json`: health check dasar server.
- `collections/01-auth.postman_collection.json`: module auth sebagai entry point untuk mendapat `bearerToken`.
- `collections/02-management.postman_collection.json`: module onboarding tenant, profile, tags, posts, static pages, dan media management.
- `collections/03-worship.postman_collection.json`: module pengaturan waktu salat, jadwal harian, petugas ibadah, hari besar, dan kalender ibadah.
- `collections/04-community.postman_collection.json`: module event, album galeri, item galeri, dan anggota pengurus tenant.
- `collections/05-engagement.postman_collection.json`: module channel donasi statis, social links, external links, feature catalog, dan website features.
- `collections/06-finance.postman_collection.json`: module subscription, konfigurasi payment gateway, campaign donasi, dan transaksi campaign tenant.
- `collections/07-public.postman_collection.json`: endpoint public website untuk event, galeri, pengurus, channel donasi, social link, external link, campaign, donor, dan donate.
- `collections/08-webhook.postman_collection.json`: simulasi callback Midtrans untuk skenario order tidak dikenal dan order valid dengan signature terhitung otomatis.

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

## Module Tersedia

- `00-system`
- `01-auth`
- `02-management`
- `03-worship`
- `04-community`
- `05-engagement`
- `06-finance`
- `07-public`
- `08-webhook`

## Rencana Module Berikutnya

- `swagger` opsional

## Catatan Finance dan Webhook

- Route tenant di `06-finance` yang terkait `pg_digital` tetap bisa ditolak middleware jika tenant belum berada pada plan yang membuka fitur tersebut.
- Route public `donate` memerlukan campaign aktif dan konfigurasi payment gateway yang valid.
- Collection `08-webhook` menyediakan skenario `unknown order` yang aman untuk smoke test, serta template `existing order` yang menghitung `signature_key` otomatis dari `order_id + status_code + gross_amount + server_key`.
