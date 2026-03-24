package engagement

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	ListStaticPaymentMethods(ctx context.Context, tenantID string, q ListQuery) ([]StaticPaymentMethodResponse, int64, error)
	CreateStaticPaymentMethod(ctx context.Context, tenantID string, req StaticPaymentMethodPayload) (*StaticPaymentMethodResponse, error)
	GetStaticPaymentMethod(ctx context.Context, tenantID string, id int64) (*StaticPaymentMethodResponse, error)
	UpdateStaticPaymentMethod(ctx context.Context, tenantID string, id int64, req StaticPaymentMethodPayload) error
	DeleteStaticPaymentMethod(ctx context.Context, tenantID string, id int64) error
	ListPublicStaticPaymentMethods(ctx context.Context, hostname string, q ListQuery) ([]StaticPaymentMethodResponse, int64, error)

	ListSocialLinks(ctx context.Context, tenantID string, q ListQuery) ([]SocialLinkResponse, int64, error)
	CreateSocialLink(ctx context.Context, tenantID string, req SocialLinkPayload) (*SocialLinkResponse, error)
	GetSocialLink(ctx context.Context, tenantID string, id int64) (*SocialLinkResponse, error)
	UpdateSocialLink(ctx context.Context, tenantID string, id int64, req SocialLinkPayload) error
	DeleteSocialLink(ctx context.Context, tenantID string, id int64) error
	ListPublicSocialLinks(ctx context.Context, hostname string, q ListQuery) ([]SocialLinkResponse, int64, error)

	ListExternalLinks(ctx context.Context, tenantID string, q ListQuery) ([]ExternalLinkResponse, int64, error)
	CreateExternalLink(ctx context.Context, tenantID string, req ExternalLinkPayload) (*ExternalLinkResponse, error)
	GetExternalLink(ctx context.Context, tenantID string, id int64) (*ExternalLinkResponse, error)
	UpdateExternalLink(ctx context.Context, tenantID string, id int64, req ExternalLinkPayload) error
	DeleteExternalLink(ctx context.Context, tenantID string, id int64) error
	ListPublicExternalLinks(ctx context.Context, hostname string, q ListQuery) ([]ExternalLinkResponse, int64, error)

	ListFeatureCatalog(ctx context.Context) ([]FeatureCatalogResponse, error)
	ListWebsiteFeatures(ctx context.Context, tenantID string) ([]WebsiteFeatureResponse, error)
	UpsertWebsiteFeature(ctx context.Context, tenantID string, featureID int64, req WebsiteFeatureUpdateRequest) error
	BulkUpsertWebsiteFeatures(ctx context.Context, tenantID string, items []WebsiteFeatureBulkItem) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) ListStaticPaymentMethods(ctx context.Context, tenantID string, q ListQuery) ([]StaticPaymentMethodResponse, int64, error) {
	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM static_payment_methods WHERE tenant_id=$1`, tenantID).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (q.Page - 1) * q.Limit
	rows, err := r.db.Query(ctx, `SELECT id,channel_type::text,label,bank_name,bank_branch,account_number,account_holder_name,qris_image_url,merchant_id,description,sort_order,is_public
		FROM static_payment_methods WHERE tenant_id=$1 ORDER BY sort_order ASC,id DESC LIMIT $2 OFFSET $3`, tenantID, q.Limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []StaticPaymentMethodResponse
	for rows.Next() {
		var it StaticPaymentMethodResponse
		if err := rows.Scan(&it.ID, &it.ChannelType, &it.Label, &it.BankName, &it.BankBranch, &it.AccountNumber, &it.AccountHolderName, &it.QrisImageURL, &it.MerchantID, &it.Description, &it.SortOrder, &it.IsPublic); err != nil {
			return nil, 0, err
		}
		items = append(items, it)
	}
	return items, total, nil
}

func (r *repository) CreateStaticPaymentMethod(ctx context.Context, tenantID string, req StaticPaymentMethodPayload) (*StaticPaymentMethodResponse, error) {
	var out StaticPaymentMethodResponse
	err := r.db.QueryRow(ctx, `INSERT INTO static_payment_methods
		(tenant_id,channel_type,label,bank_name,bank_branch,account_number,account_holder_name,qris_image_url,merchant_id,description,sort_order,is_public)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		RETURNING id,channel_type::text,label,bank_name,bank_branch,account_number,account_holder_name,qris_image_url,merchant_id,description,sort_order,is_public`,
		tenantID, req.ChannelType, req.Label, req.BankName, req.BankBranch, req.AccountNumber, req.AccountHolderName, req.QrisImageURL, req.MerchantID, req.Description, req.SortOrder, req.IsPublic).
		Scan(&out.ID, &out.ChannelType, &out.Label, &out.BankName, &out.BankBranch, &out.AccountNumber, &out.AccountHolderName, &out.QrisImageURL, &out.MerchantID, &out.Description, &out.SortOrder, &out.IsPublic)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *repository) GetStaticPaymentMethod(ctx context.Context, tenantID string, id int64) (*StaticPaymentMethodResponse, error) {
	var out StaticPaymentMethodResponse
	err := r.db.QueryRow(ctx, `SELECT id,channel_type::text,label,bank_name,bank_branch,account_number,account_holder_name,qris_image_url,merchant_id,description,sort_order,is_public
		FROM static_payment_methods WHERE tenant_id=$1 AND id=$2`, tenantID, id).
		Scan(&out.ID, &out.ChannelType, &out.Label, &out.BankName, &out.BankBranch, &out.AccountNumber, &out.AccountHolderName, &out.QrisImageURL, &out.MerchantID, &out.Description, &out.SortOrder, &out.IsPublic)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &out, nil
}

func (r *repository) UpdateStaticPaymentMethod(ctx context.Context, tenantID string, id int64, req StaticPaymentMethodPayload) error {
	tag, err := r.db.Exec(ctx, `UPDATE static_payment_methods SET channel_type=$1,label=$2,bank_name=$3,bank_branch=$4,account_number=$5,account_holder_name=$6,qris_image_url=$7,merchant_id=$8,description=$9,sort_order=$10,is_public=$11,updated_at=now()
		WHERE tenant_id=$12 AND id=$13`,
		req.ChannelType, req.Label, req.BankName, req.BankBranch, req.AccountNumber, req.AccountHolderName, req.QrisImageURL, req.MerchantID, req.Description, req.SortOrder, req.IsPublic, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) DeleteStaticPaymentMethod(ctx context.Context, tenantID string, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM static_payment_methods WHERE tenant_id=$1 AND id=$2`, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) ListPublicStaticPaymentMethods(ctx context.Context, hostname string, q ListQuery) ([]StaticPaymentMethodResponse, int64, error) {
	tenantID, err := r.resolveTenantByHostname(ctx, hostname)
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM static_payment_methods WHERE tenant_id=$1 AND is_public=true`, tenantID).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (q.Page - 1) * q.Limit
	rows, err := r.db.Query(ctx, `SELECT id,channel_type::text,label,bank_name,bank_branch,account_number,account_holder_name,qris_image_url,merchant_id,description,sort_order,is_public
		FROM static_payment_methods WHERE tenant_id=$1 AND is_public=true ORDER BY sort_order ASC,id DESC LIMIT $2 OFFSET $3`, tenantID, q.Limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []StaticPaymentMethodResponse
	for rows.Next() {
		var it StaticPaymentMethodResponse
		if err := rows.Scan(&it.ID, &it.ChannelType, &it.Label, &it.BankName, &it.BankBranch, &it.AccountNumber, &it.AccountHolderName, &it.QrisImageURL, &it.MerchantID, &it.Description, &it.SortOrder, &it.IsPublic); err != nil {
			return nil, 0, err
		}
		items = append(items, it)
	}
	return items, total, nil
}

func (r *repository) ListSocialLinks(ctx context.Context, tenantID string, q ListQuery) ([]SocialLinkResponse, int64, error) {
	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM social_links WHERE tenant_id=$1`, tenantID).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (q.Page - 1) * q.Limit
	rows, err := r.db.Query(ctx, `SELECT id,platform,account_name,url,description,show_in_footer,show_in_contact_page,sort_order
		FROM social_links WHERE tenant_id=$1 ORDER BY sort_order ASC,id DESC LIMIT $2 OFFSET $3`, tenantID, q.Limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []SocialLinkResponse
	for rows.Next() {
		var it SocialLinkResponse
		if err := rows.Scan(&it.ID, &it.Platform, &it.AccountName, &it.URL, &it.Description, &it.ShowInFooter, &it.ShowInContactPage, &it.SortOrder); err != nil {
			return nil, 0, err
		}
		items = append(items, it)
	}
	return items, total, nil
}

func (r *repository) CreateSocialLink(ctx context.Context, tenantID string, req SocialLinkPayload) (*SocialLinkResponse, error) {
	var out SocialLinkResponse
	err := r.db.QueryRow(ctx, `INSERT INTO social_links
		(tenant_id,platform,account_name,url,description,show_in_footer,show_in_contact_page,sort_order)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id,platform,account_name,url,description,show_in_footer,show_in_contact_page,sort_order`,
		tenantID, req.Platform, req.AccountName, req.URL, req.Description, req.ShowInFooter, req.ShowInContactPage, req.SortOrder).
		Scan(&out.ID, &out.Platform, &out.AccountName, &out.URL, &out.Description, &out.ShowInFooter, &out.ShowInContactPage, &out.SortOrder)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *repository) GetSocialLink(ctx context.Context, tenantID string, id int64) (*SocialLinkResponse, error) {
	var out SocialLinkResponse
	err := r.db.QueryRow(ctx, `SELECT id,platform,account_name,url,description,show_in_footer,show_in_contact_page,sort_order
		FROM social_links WHERE tenant_id=$1 AND id=$2`, tenantID, id).
		Scan(&out.ID, &out.Platform, &out.AccountName, &out.URL, &out.Description, &out.ShowInFooter, &out.ShowInContactPage, &out.SortOrder)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &out, nil
}

func (r *repository) UpdateSocialLink(ctx context.Context, tenantID string, id int64, req SocialLinkPayload) error {
	tag, err := r.db.Exec(ctx, `UPDATE social_links SET platform=$1,account_name=$2,url=$3,description=$4,show_in_footer=$5,show_in_contact_page=$6,sort_order=$7,updated_at=now()
		WHERE tenant_id=$8 AND id=$9`,
		req.Platform, req.AccountName, req.URL, req.Description, req.ShowInFooter, req.ShowInContactPage, req.SortOrder, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) DeleteSocialLink(ctx context.Context, tenantID string, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM social_links WHERE tenant_id=$1 AND id=$2`, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) ListPublicSocialLinks(ctx context.Context, hostname string, q ListQuery) ([]SocialLinkResponse, int64, error) {
	tenantID, err := r.resolveTenantByHostname(ctx, hostname)
	if err != nil {
		return nil, 0, err
	}
	return r.ListSocialLinks(ctx, tenantID, q)
}

func (r *repository) ListExternalLinks(ctx context.Context, tenantID string, q ListQuery) ([]ExternalLinkResponse, int64, error) {
	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM external_links WHERE tenant_id=$1`, tenantID).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (q.Page - 1) * q.Limit
	rows, err := r.db.Query(ctx, `SELECT id,link_type,label,url,note,visibility,sort_order
		FROM external_links WHERE tenant_id=$1 ORDER BY sort_order ASC,id DESC LIMIT $2 OFFSET $3`, tenantID, q.Limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []ExternalLinkResponse
	for rows.Next() {
		var it ExternalLinkResponse
		if err := rows.Scan(&it.ID, &it.LinkType, &it.Label, &it.URL, &it.Note, &it.Visibility, &it.SortOrder); err != nil {
			return nil, 0, err
		}
		items = append(items, it)
	}
	return items, total, nil
}

func (r *repository) CreateExternalLink(ctx context.Context, tenantID string, req ExternalLinkPayload) (*ExternalLinkResponse, error) {
	var out ExternalLinkResponse
	err := r.db.QueryRow(ctx, `INSERT INTO external_links (tenant_id,link_type,label,url,note,visibility,sort_order)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id,link_type,label,url,note,visibility,sort_order`,
		tenantID, req.LinkType, req.Label, req.URL, req.Note, req.Visibility, req.SortOrder).
		Scan(&out.ID, &out.LinkType, &out.Label, &out.URL, &out.Note, &out.Visibility, &out.SortOrder)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *repository) GetExternalLink(ctx context.Context, tenantID string, id int64) (*ExternalLinkResponse, error) {
	var out ExternalLinkResponse
	err := r.db.QueryRow(ctx, `SELECT id,link_type,label,url,note,visibility,sort_order
		FROM external_links WHERE tenant_id=$1 AND id=$2`, tenantID, id).
		Scan(&out.ID, &out.LinkType, &out.Label, &out.URL, &out.Note, &out.Visibility, &out.SortOrder)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &out, nil
}

func (r *repository) UpdateExternalLink(ctx context.Context, tenantID string, id int64, req ExternalLinkPayload) error {
	tag, err := r.db.Exec(ctx, `UPDATE external_links SET link_type=$1,label=$2,url=$3,note=$4,visibility=$5,sort_order=$6,updated_at=now() WHERE tenant_id=$7 AND id=$8`,
		req.LinkType, req.Label, req.URL, req.Note, req.Visibility, req.SortOrder, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) DeleteExternalLink(ctx context.Context, tenantID string, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM external_links WHERE tenant_id=$1 AND id=$2`, tenantID, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) ListPublicExternalLinks(ctx context.Context, hostname string, q ListQuery) ([]ExternalLinkResponse, int64, error) {
	tenantID, err := r.resolveTenantByHostname(ctx, hostname)
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM external_links WHERE tenant_id=$1 AND visibility='public'`, tenantID).Scan(&total); err != nil {
		return nil, 0, err
	}
	offset := (q.Page - 1) * q.Limit
	rows, err := r.db.Query(ctx, `SELECT id,link_type,label,url,note,visibility,sort_order
		FROM external_links WHERE tenant_id=$1 AND visibility='public' ORDER BY sort_order ASC,id DESC LIMIT $2 OFFSET $3`, tenantID, q.Limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []ExternalLinkResponse
	for rows.Next() {
		var it ExternalLinkResponse
		if err := rows.Scan(&it.ID, &it.LinkType, &it.Label, &it.URL, &it.Note, &it.Visibility, &it.SortOrder); err != nil {
			return nil, 0, err
		}
		items = append(items, it)
	}
	return items, total, nil
}

func (r *repository) ListFeatureCatalog(ctx context.Context) ([]FeatureCatalogResponse, error) {
	rows, err := r.db.Query(ctx, `SELECT id,feature_type::text,name,category_label FROM feature_catalog ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeatureCatalogResponse
	for rows.Next() {
		var it FeatureCatalogResponse
		if err := rows.Scan(&it.ID, &it.FeatureType, &it.Name, &it.CategoryLabel); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, nil
}

func (r *repository) ListWebsiteFeatures(ctx context.Context, tenantID string) ([]WebsiteFeatureResponse, error) {
	rows, err := r.db.Query(ctx, `SELECT wf.id,wf.feature_id,fc.name,wf.enabled,wf.is_active,wf.detail,wf.note
		FROM website_features wf
		JOIN feature_catalog fc ON fc.id=wf.feature_id
		WHERE wf.tenant_id=$1 ORDER BY wf.id ASC`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []WebsiteFeatureResponse
	for rows.Next() {
		var it WebsiteFeatureResponse
		if err := rows.Scan(&it.ID, &it.FeatureID, &it.Name, &it.Enabled, &it.IsActive, &it.Detail, &it.Note); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, nil
}

func (r *repository) UpsertWebsiteFeature(ctx context.Context, tenantID string, featureID int64, req WebsiteFeatureUpdateRequest) error {
	_, err := r.db.Exec(ctx, `INSERT INTO website_features (tenant_id,feature_id,enabled,is_active,detail,note)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (tenant_id,feature_id) DO UPDATE SET enabled=EXCLUDED.enabled,is_active=EXCLUDED.is_active,detail=EXCLUDED.detail,note=EXCLUDED.note,updated_at=now()`,
		tenantID, featureID, req.Enabled, req.IsActive, req.Detail, req.Note)
	return err
}

func (r *repository) BulkUpsertWebsiteFeatures(ctx context.Context, tenantID string, items []WebsiteFeatureBulkItem) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for _, item := range items {
		if _, err := tx.Exec(ctx, `INSERT INTO website_features (tenant_id,feature_id,enabled,is_active,detail,note)
			VALUES ($1,$2,$3,$4,$5,$6)
			ON CONFLICT (tenant_id,feature_id) DO UPDATE SET enabled=EXCLUDED.enabled,is_active=EXCLUDED.is_active,detail=EXCLUDED.detail,note=EXCLUDED.note,updated_at=now()`,
			tenantID, item.FeatureID, item.Enabled, item.IsActive, item.Detail, item.Note); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *repository) resolveTenantByHostname(ctx context.Context, hostname string) (string, error) {
	var tenantID string
	err := r.db.QueryRow(ctx, `SELECT tenant_id::text FROM website_domains WHERE hostname=$1 AND status='active' LIMIT 1`, hostname).Scan(&tenantID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", ErrNotFound
		}
		return "", err
	}
	return tenantID, nil
}

func (r *repository) debugString(id int64) string {
	return fmt.Sprint(id)
}
