-- Hapus tabel usage
DROP TABLE IF EXISTS tenant_usages;

-- Hapus kolom yang ditambahkan di tabel tenants
ALTER TABLE tenants
DROP COLUMN IF EXISTS active_template_id,
DROP COLUMN IF EXISTS subscription_plan;