-- 1. Tambahkan kolom paket dan template di tabel tenants
ALTER TABLE tenants 
ADD COLUMN subscription_plan VARCHAR(50) NOT NULL DEFAULT 'free',
ADD COLUMN active_template_id VARCHAR(50) NOT NULL DEFAULT 'template_default';

-- 2. Buat tabel baru untuk melacak penggunaan (Usage Tracking)
CREATE TABLE tenant_usages (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    storage_used_mb NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_tenant_usage UNIQUE (tenant_id)
);

-- Buat index untuk mempercepat pencarian
CREATE INDEX idx_tenant_usages_tenant_id ON tenant_usages(tenant_id);

-- 3. (Opsional tapi Penting) Seed data awal untuk tenant yang sudah ada
-- Agar tenant lama (termasuk masjid testing kita) langsung punya record usage 0 MB
INSERT INTO tenant_usages (tenant_id, storage_used_mb)
SELECT id, 0.00 FROM tenants
ON CONFLICT (tenant_id) DO NOTHING;