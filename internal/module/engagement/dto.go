package engagement

type DonationChannelPayload struct {
	ChannelType       string  `json:"channel_type"`
	Label             string  `json:"label"`
	BankName          *string `json:"bank_name"`
	BankBranch        *string `json:"bank_branch"`
	AccountNumber     *string `json:"account_number"`
	AccountHolderName *string `json:"account_holder_name"`
	QrisImageURL      *string `json:"qris_image_url"`
	MerchantID        *string `json:"merchant_id"`
	Description       *string `json:"description"`
	SortOrder         int     `json:"sort_order"`
	IsPublic          bool    `json:"is_public"`
}

type DonationChannelResponse struct {
	ID                int64   `json:"id"`
	ChannelType       string  `json:"channel_type"`
	Label             string  `json:"label"`
	BankName          *string `json:"bank_name"`
	BankBranch        *string `json:"bank_branch"`
	AccountNumber     *string `json:"account_number"`
	AccountHolderName *string `json:"account_holder_name"`
	QrisImageURL      *string `json:"qris_image_url"`
	MerchantID        *string `json:"merchant_id"`
	Description       *string `json:"description"`
	SortOrder         int     `json:"sort_order"`
	IsPublic          bool    `json:"is_public"`
}

type SocialLinkPayload struct {
	Platform          string  `json:"platform"`
	AccountName       *string `json:"account_name"`
	URL               string  `json:"url"`
	Description       *string `json:"description"`
	ShowInFooter      bool    `json:"show_in_footer"`
	ShowInContactPage bool    `json:"show_in_contact_page"`
	SortOrder         int     `json:"sort_order"`
}

type SocialLinkResponse struct {
	ID                int64   `json:"id"`
	Platform          string  `json:"platform"`
	AccountName       *string `json:"account_name"`
	URL               string  `json:"url"`
	Description       *string `json:"description"`
	ShowInFooter      bool    `json:"show_in_footer"`
	ShowInContactPage bool    `json:"show_in_contact_page"`
	SortOrder         int     `json:"sort_order"`
}

type ExternalLinkPayload struct {
	LinkType   string  `json:"link_type"`
	Label      string  `json:"label"`
	URL        string  `json:"url"`
	Note       *string `json:"note"`
	Visibility string  `json:"visibility"`
	SortOrder  int     `json:"sort_order"`
}

type ExternalLinkResponse struct {
	ID         int64   `json:"id"`
	LinkType   string  `json:"link_type"`
	Label      string  `json:"label"`
	URL        string  `json:"url"`
	Note       *string `json:"note"`
	Visibility string  `json:"visibility"`
	SortOrder  int     `json:"sort_order"`
}

type FeatureCatalogResponse struct {
	ID            int64   `json:"id"`
	FeatureType   string  `json:"feature_type"`
	Name          string  `json:"name"`
	CategoryLabel *string `json:"category_label"`
}

type WebsiteFeatureResponse struct {
	ID        int64   `json:"id"`
	FeatureID int64   `json:"feature_id"`
	Name      string  `json:"name"`
	Enabled   bool    `json:"enabled"`
	IsActive  bool    `json:"is_active"`
	Detail    *string `json:"detail"`
	Note      *string `json:"note"`
}

type WebsiteFeatureUpdateRequest struct {
	Enabled  bool    `json:"enabled"`
	IsActive bool    `json:"is_active"`
	Detail   *string `json:"detail"`
	Note     *string `json:"note"`
}

type WebsiteFeatureBulkRequest struct {
	Items []WebsiteFeatureBulkItem `json:"items"`
}

type WebsiteFeatureBulkItem struct {
	FeatureID int64   `json:"feature_id"`
	Enabled   bool    `json:"enabled"`
	IsActive  bool    `json:"is_active"`
	Detail    *string `json:"detail"`
	Note      *string `json:"note"`
}

type ListQuery struct {
	Page  int
	Limit int
}
