package router

// @Summary Health check API
// @Tags System
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /health [get]
func swaggerHealthAPI() {}

// @Summary API docs links
// @Tags System
// @Success 200 {object} map[string]interface{}
// @Router /docs [get]
func swaggerAPIDocsLinks() {}

// @Summary Register tenant admin
// @Tags Auth
// @Param payload body object true "payload"
// @Success 201 {object} map[string]interface{}
// @Router /auth/register [post]
func swaggerAuthRegister() {}

// @Summary Login with email and password
// @Tags Auth
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /auth/login [post]
func swaggerAuthLogin() {}

// @Summary Request forgot password email
// @Tags Auth
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /auth/forgot-password [post]
func swaggerAuthForgotPassword() {}

// @Summary Reset password
// @Tags Auth
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /auth/reset-password [post]
func swaggerAuthResetPassword() {}

// @Summary Login with Google token
// @Tags Auth
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /auth/google [post]
func swaggerAuthGoogleLogin() {}

// @Summary Get payment gateway config
// @Tags Finance
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tenant/pg-config [get]
func swaggerTenantPGConfigGet() {}

// @Summary Upsert payment gateway config
// @Tags Finance
// @Security BearerAuth
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/pg-config [put]
func swaggerTenantPGConfigUpsert() {}

// @Summary List tenant campaigns
// @Tags Finance
// @Security BearerAuth
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/campaigns [get]
func swaggerTenantCampaignsList() {}

// @Summary Create tenant campaign
// @Tags Finance
// @Security BearerAuth
// @Param payload body object true "payload"
// @Success 201 {object} map[string]interface{}
// @Router /tenant/campaigns [post]
func swaggerTenantCampaignsCreate() {}

// @Summary Get tenant campaign detail
// @Tags Finance
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/campaigns/{id} [get]
func swaggerTenantCampaignsGet() {}

// @Summary Update tenant campaign
// @Tags Finance
// @Security BearerAuth
// @Param id path int true "id"
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/campaigns/{id} [put]
func swaggerTenantCampaignsUpdate() {}

// @Summary List tenant campaign transactions
// @Tags Finance
// @Security BearerAuth
// @Param id path int true "id"
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/campaigns/{id}/transactions [get]
func swaggerTenantCampaignTransactionsList() {}

// @Summary Get subscription quote
// @Tags Billing
// @Security BearerAuth
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/subscription/quote [post]
func swaggerTenantSubscriptionQuote() {}

// @Summary Create subscription checkout from quote
// @Tags Billing
// @Security BearerAuth
// @Param payload body object true "payload"
// @Success 201 {object} map[string]interface{}
// @Router /tenant/subscription/checkout-v2 [post]
func swaggerTenantSubscriptionCheckoutV2Create() {}

// @Summary List subscription transactions
// @Tags Billing
// @Security BearerAuth
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/subscription/transactions [get]
func swaggerTenantSubscriptionTransactionsList() {}

// @Summary Public campaigns
// @Tags Public
// @Param hostname path string true "hostname"
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {object} map[string]interface{}
// @Router /public/{hostname}/campaigns [get]
func swaggerPublicCampaignsList() {}

// @Summary Public campaign detail
// @Tags Public
// @Param hostname path string true "hostname"
// @Param slug path string true "slug"
// @Success 200 {object} map[string]interface{}
// @Router /public/{hostname}/campaigns/{slug} [get]
func swaggerPublicCampaignsGetBySlug() {}

// @Summary Public campaign donors
// @Tags Public
// @Param hostname path string true "hostname"
// @Param id path int true "id"
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {object} map[string]interface{}
// @Router /public/{hostname}/campaigns/{id}/donors [get]
func swaggerPublicCampaignDonorsList() {}

// @Summary Create public donation checkout
// @Tags Public
// @Param hostname path string true "hostname"
// @Param id path int true "id"
// @Param payload body object true "payload"
// @Success 201 {object} map[string]interface{}
// @Router /public/{hostname}/campaigns/{id}/donate [post]
func swaggerPublicCampaignDonate() {}

// @Summary Midtrans webhook callback
// @Tags Webhook
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /webhook/midtrans [post]
func swaggerMidtransWebhook() {}
