-- Tabel 1: Konfigurasi Payment Gateway per Masjid
CREATE TABLE pg_configs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    use_central_pg BOOLEAN NOT NULL DEFAULT true,
    provider VARCHAR(50) NOT NULL CHECK (provider IN ('midtrans', 'xendit', 'tripay', 'etakmir_pusat')),
    client_key TEXT, -- Bisa null jika pakai pusat
    server_key TEXT, -- Disimpan dalam bentuk terenkripsi
    is_production BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_tenant_pg_config UNIQUE (tenant_id)
);

-- Tabel 2: Kampanye Donasi (Crowdfunding)
CREATE TABLE donation_campaigns (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT,
    image_url TEXT,
    target_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    collected_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_campaign_slug_per_tenant UNIQUE (tenant_id, slug)
);

-- Tabel 3: Transaksi Donasi
CREATE TABLE donation_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- ID ini yang dikirim ke PG sebagai Order ID
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    campaign_id BIGINT NOT NULL REFERENCES donation_campaigns(id) ON DELETE CASCADE,
    donor_name VARCHAR(255),
    is_anonymous BOOLEAN NOT NULL DEFAULT false,
    amount NUMERIC(15, 2) NOT NULL,
    payment_method VARCHAR(100), -- cth: qris, bank_transfer_bca
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'paid', 'expired', 'failed', 'refunded')),
    payment_url TEXT, -- URL Checkout dari Midtrans/Xendit
    snap_token VARCHAR(255), -- Token dari Midtrans
    pg_reference_id VARCHAR(255), -- Transaction ID dari pihak PG (berguna untuk rekonsiliasi)
    paid_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes untuk mempercepat pencarian
CREATE INDEX idx_donation_campaigns_tenant ON donation_campaigns(tenant_id);
CREATE INDEX idx_donation_transactions_tenant ON donation_transactions(tenant_id);
CREATE INDEX idx_donation_transactions_campaign ON donation_transactions(campaign_id);
CREATE INDEX idx_donation_transactions_status ON donation_transactions(status);