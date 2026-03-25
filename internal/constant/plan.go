package constant

// 1. Definisi Key Fitur (Agar tidak terjadi typo saat mengetik string manual)
const (
	FeatureProfile       = "profile"
	FeatureSchedules     = "schedules"
	FeatureFacilities    = "facilities"
	FeatureStaticPayment = "static_payment"
	FeatureManagement    = "management"
	FeatureEvents        = "events"
	FeatureArticles      = "articles"
	FeatureGallery       = "gallery"
	FeaturePGDigital     = "pg_digital"
	FeatureCustomDomain  = "custom_domain"
)

// 2. Definisi Nama Plan
const (
	PlanFree        = "free"
	PlanPremiumPlus = "premium_plus"
	PlanProPlus     = "pro_plus"
	PlanMaxPlus     = "max_plus"
)

// 3. Struktur Detail Paket
type PlanDetail struct {
	Name                  string
	Price                 float64
	MaxTemplates          int
	StorageLimitMB        float64
	PlatformFeePercentage float64
	AttributionEnabled    bool
	FeaturesUnlocked      []string
}

// 4. THE SINGLE SOURCE OF TRUTH (Kamus Utama SaaS Kita)
var SubscriptionPlans = map[string]PlanDetail{
	PlanFree: {
		Name:                  "FREE",
		Price:                 0,
		MaxTemplates:          1,
		StorageLimitMB:        0,
		PlatformFeePercentage: 0, // Tidak ada PG
		AttributionEnabled:    true,
		FeaturesUnlocked: []string{
			FeatureProfile,
		},
	},
	PlanPremiumPlus: {
		Name:                  "PREMIUM+",
		Price:                 24900,
		MaxTemplates:          3,
		StorageLimitMB:        0,
		PlatformFeePercentage: 0, // Tidak ada PG
		AttributionEnabled:    false,
		FeaturesUnlocked: []string{
			FeatureProfile,
			FeatureSchedules,
			FeatureStaticPayment,
			FeatureFacilities,
		},
	},
	PlanProPlus: {
		Name:                  "PRO++",
		Price:                 79900,
		MaxTemplates:          5,
		StorageLimitMB:        500, // 500 MB
		PlatformFeePercentage: 0.5, // Potongan 0.5% dari tiap transaksi donasi digital
		AttributionEnabled:    false,
		FeaturesUnlocked: []string{
			FeatureProfile,
			FeatureSchedules,
			FeatureStaticPayment,
			FeatureFacilities,
			FeatureManagement,
			FeatureEvents,
			FeatureArticles,
			FeatureGallery,
			FeaturePGDigital,
		},
	},
	PlanMaxPlus: {
		Name:                  "MAX+++",
		Price:                 149000,
		MaxTemplates:          5,    // + Customable (Nanti di-handle logic UI frontend)
		StorageLimitMB:        1000, // 1 GB
		PlatformFeePercentage: 0.0,  // Fee 0% untuk paket tertinggi
		AttributionEnabled:    false,
		FeaturesUnlocked: []string{
			FeatureProfile,
			FeatureSchedules,
			FeatureStaticPayment,
			FeatureFacilities,
			FeatureManagement,
			FeatureEvents,
			FeatureArticles,
			FeatureGallery,
			FeaturePGDigital,
			FeatureCustomDomain,
		},
	},
}

// 5. Helper Function: Cek apakah paket memiliki fitur tertentu
func HasFeature(planKey string, featureKey string) bool {
	plan, exists := SubscriptionPlans[planKey]
	if !exists {
		return false // Kalau nama paket ngawur, tolak otomatis
	}

	for _, f := range plan.FeaturesUnlocked {
		if f == featureKey {
			return true
		}
	}
	return false
}
