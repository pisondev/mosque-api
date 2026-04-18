ALTER TABLE tenants
ADD COLUMN onboarding_completed BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN onboarding_payment_status VARCHAR(20) NOT NULL DEFAULT 'pending';

CREATE TABLE subscription_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    order_id VARCHAR(120) NOT NULL UNIQUE,
    plan_code VARCHAR(50) NOT NULL,
    amount NUMERIC(14,2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'IDR',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    payment_method VARCHAR(100),
    midtrans_transaction_id VARCHAR(120),
    snap_token TEXT,
    payment_url TEXT,
    raw_notification JSONB,
    paid_at TIMESTAMPTZ,
    expired_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT subscription_status_check CHECK (status IN ('pending', 'paid', 'failed', 'expired'))
);

CREATE INDEX idx_subscription_transactions_tenant_id ON subscription_transactions(tenant_id);
CREATE INDEX idx_subscription_transactions_order_id ON subscription_transactions(order_id);
CREATE INDEX idx_subscription_transactions_status ON subscription_transactions(status);

UPDATE tenants
SET status = 'pending',
    onboarding_completed = false,
    onboarding_payment_status = 'pending',
    updated_at = NOW();
