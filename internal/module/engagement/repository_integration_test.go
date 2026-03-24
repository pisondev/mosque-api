package engagement

import (
	"context"
	"testing"

	"github.com/pisondev/mosque-api/internal/testutil"
)

func TestRepositoryIntegrationEngagement(t *testing.T) {
	db := testutil.OpenTestDB(t)
	defer db.Close()

	tenant := testutil.SeedTenant(t, db)
	defer testutil.CleanupTenant(t, db, tenant.ID)

	repo := NewRepository(db)
	ctx := context.Background()

	donationRes, err := repo.CreateStaticPaymentMethod(ctx, tenant.ID, StaticPaymentMethodPayload{
		ChannelType:       "bank_account",
		Label:             "BSI Integrasi",
		BankName:          ptr("BSI"),
		AccountNumber:     ptr("123456"),
		AccountHolderName: ptr("DKM Integrasi"),
		SortOrder:         1,
		IsPublic:          true,
	})
	if err != nil {
		t.Fatalf("CreateStaticPaymentMethod failed: %v", err)
	}
	if donationRes.ID == 0 {
		t.Fatal("CreateStaticPaymentMethod returned empty id")
	}

	donationList, totalDonation, err := repo.ListStaticPaymentMethods(ctx, tenant.ID, ListQuery{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("ListStaticPaymentMethods failed: %v", err)
	}
	if totalDonation < 1 || len(donationList) < 1 {
		t.Fatalf("ListStaticPaymentMethods expected data, total=%d len=%d", totalDonation, len(donationList))
	}

	gotDonation, err := repo.GetStaticPaymentMethod(ctx, tenant.ID, donationRes.ID)
	if err != nil {
		t.Fatalf("GetStaticPaymentMethod failed: %v", err)
	}
	if gotDonation.Label != "BSI Integrasi" {
		t.Fatalf("unexpected donation label: %s", gotDonation.Label)
	}

	if err := repo.UpdateStaticPaymentMethod(ctx, tenant.ID, donationRes.ID, StaticPaymentMethodPayload{
		ChannelType:       "bank_account",
		Label:             "BSI Integrasi Updated",
		BankName:          ptr("BSI"),
		AccountNumber:     ptr("123456"),
		AccountHolderName: ptr("DKM Integrasi"),
		SortOrder:         1,
		IsPublic:          true,
	}); err != nil {
		t.Fatalf("UpdateStaticPaymentMethod failed: %v", err)
	}

	socialRes, err := repo.CreateSocialLink(ctx, tenant.ID, SocialLinkPayload{
		Platform:          "instagram",
		URL:               "https://instagram.com/masjid-it",
		ShowInFooter:      true,
		ShowInContactPage: true,
		SortOrder:         1,
	})
	if err != nil {
		t.Fatalf("CreateSocialLink failed: %v", err)
	}
	if socialRes.ID == 0 {
		t.Fatal("CreateSocialLink returned empty id")
	}

	externalRes, err := repo.CreateExternalLink(ctx, tenant.ID, ExternalLinkPayload{
		LinkType:   "registration",
		Label:      "Pendaftaran",
		URL:        "https://example.com/register",
		Visibility: "public",
		SortOrder:  1,
	})
	if err != nil {
		t.Fatalf("CreateExternalLink failed: %v", err)
	}
	if externalRes.ID == 0 {
		t.Fatal("CreateExternalLink returned empty id")
	}

	publicDonation, _, err := repo.ListPublicStaticPaymentMethods(ctx, tenant.Hostname, ListQuery{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("ListPublicStaticPaymentMethods failed: %v", err)
	}
	if len(publicDonation) < 1 {
		t.Fatal("ListPublicStaticPaymentMethods expected data")
	}

	publicSocial, _, err := repo.ListPublicSocialLinks(ctx, tenant.Hostname, ListQuery{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("ListPublicSocialLinks failed: %v", err)
	}
	if len(publicSocial) < 1 {
		t.Fatal("ListPublicSocialLinks expected data")
	}

	publicExternal, _, err := repo.ListPublicExternalLinks(ctx, tenant.Hostname, ListQuery{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("ListPublicExternalLinks failed: %v", err)
	}
	if len(publicExternal) < 1 {
		t.Fatal("ListPublicExternalLinks expected data")
	}

	catalog, err := repo.ListFeatureCatalog(ctx)
	if err != nil {
		t.Fatalf("ListFeatureCatalog failed: %v", err)
	}
	if len(catalog) < 1 {
		t.Fatal("ListFeatureCatalog expected seeded data")
	}

	featureID := catalog[0].ID
	if err := repo.UpsertWebsiteFeature(ctx, tenant.ID, featureID, WebsiteFeatureUpdateRequest{
		Enabled:  true,
		IsActive: true,
		Detail:   ptr("on"),
		Note:     ptr("integration"),
	}); err != nil {
		t.Fatalf("UpsertWebsiteFeature failed: %v", err)
	}

	if err := repo.BulkUpsertWebsiteFeatures(ctx, tenant.ID, []WebsiteFeatureBulkItem{
		{
			FeatureID: featureID,
			Enabled:   true,
			IsActive:  true,
			Detail:    ptr("bulk"),
			Note:      ptr("bulk"),
		},
	}); err != nil {
		t.Fatalf("BulkUpsertWebsiteFeatures failed: %v", err)
	}

	websiteFeatures, err := repo.ListWebsiteFeatures(ctx, tenant.ID)
	if err != nil {
		t.Fatalf("ListWebsiteFeatures failed: %v", err)
	}
	if len(websiteFeatures) < 1 {
		t.Fatal("ListWebsiteFeatures expected data")
	}

	if err := repo.DeleteSocialLink(ctx, tenant.ID, socialRes.ID); err != nil {
		t.Fatalf("DeleteSocialLink failed: %v", err)
	}
	if err := repo.DeleteExternalLink(ctx, tenant.ID, externalRes.ID); err != nil {
		t.Fatalf("DeleteExternalLink failed: %v", err)
	}
	if err := repo.DeleteStaticPaymentMethod(ctx, tenant.ID, donationRes.ID); err != nil {
		t.Fatalf("DeleteStaticPaymentMethod failed: %v", err)
	}
}

func ptr(v string) *string {
	return &v
}
