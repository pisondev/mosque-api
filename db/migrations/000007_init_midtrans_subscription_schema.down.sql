DROP INDEX IF EXISTS idx_subscription_transactions_status;
DROP INDEX IF EXISTS idx_subscription_transactions_order_id;
DROP INDEX IF EXISTS idx_subscription_transactions_tenant_id;

DROP TABLE IF EXISTS subscription_transactions;

ALTER TABLE tenants
DROP COLUMN IF EXISTS onboarding_payment_status,
DROP COLUMN IF EXISTS onboarding_completed;
