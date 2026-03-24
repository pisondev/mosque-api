-- Mengubah nama tabel
ALTER TABLE donation_channels RENAME TO static_payment_methods;

-- Mengubah nama sequence auto-increment (penting agar insert tidak error)
ALTER SEQUENCE donation_channels_id_seq RENAME TO static_payment_methods_id_seq;

-- Mengubah nama index
ALTER INDEX idx_donation_channels_tenant RENAME TO idx_static_payment_methods_tenant;

-- Mengubah nama constraint (opsional tapi disarankan agar rapi)
ALTER TABLE static_payment_methods RENAME CONSTRAINT chk_donation_bank_fields TO chk_static_bank_fields;
ALTER TABLE static_payment_methods RENAME CONSTRAINT chk_donation_qris_fields TO chk_static_qris_fields;