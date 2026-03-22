package management

import "time"

type DomainListQuery struct {
	Status     string
	DomainType string
	Page       int
	Limit      int
}

type CreateDomainRequest struct {
	DomainType string `json:"domain_type"`
	Hostname   string `json:"hostname"`
}

type UpdateDomainRequest struct {
	Status string `json:"status"`
}

type DomainResponse struct {
	ID         int64      `json:"id"`
	DomainType string     `json:"domain_type"`
	Hostname   string     `json:"hostname"`
	Status     string     `json:"status"`
	VerifiedAt *time.Time `json:"verified_at"`
}

type ProfileRequest struct {
	OfficialName string `json:"official_name"`
	Kind         string `json:"kind"`
	ShortName    string `json:"short_name"`
	City         string `json:"city"`
	AddressFull  string `json:"address_full"`
	PhoneWA      string `json:"phone_whatsapp"`
	Email        string `json:"email"`
}

type ProfileResponse struct {
	OfficialName string `json:"official_name"`
	Kind         string `json:"kind"`
	ShortName    string `json:"short_name"`
	City         string `json:"city"`
	AddressFull  string `json:"address_full"`
	PhoneWA      string `json:"phone_whatsapp"`
	Email        string `json:"email"`
}

type CreateTagRequest struct {
	Scope string `json:"scope"`
	Name  string `json:"name"`
}

type UpdateTagRequest struct {
	Name string `json:"name"`
}

type TagResponse struct {
	ID    int64  `json:"id"`
	Scope string `json:"scope"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
}

type PostListQuery struct {
	Category  string
	Status    string
	Search    string
	Page      int
	Limit     int
	SortBy    string
	SortOrder string
}

type PostPayload struct {
	Title           string     `json:"title"`
	Category        string     `json:"category"`
	Excerpt         string     `json:"excerpt"`
	ContentMarkdown string     `json:"content_markdown"`
	ThumbnailURL    string     `json:"thumbnail_url"`
	AuthorName      string     `json:"author_name"`
	PublishedAt     *time.Time `json:"published_at"`
	ExpiredAt       *time.Time `json:"expired_at"`
	Status          string     `json:"status"`
	ShowOnHomepage  bool       `json:"show_on_homepage"`
	SortOrder       int        `json:"sort_order"`
	TagIDs          []int64    `json:"tag_ids"`
}

type UpdatePostStatusRequest struct {
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"published_at"`
}

type PostResponse struct {
	ID              int64         `json:"id"`
	Title           string        `json:"title"`
	Slug            string        `json:"slug"`
	Category        string        `json:"category"`
	Excerpt         string        `json:"excerpt"`
	ContentMarkdown string        `json:"content_markdown"`
	ThumbnailURL    string        `json:"thumbnail_url"`
	AuthorName      string        `json:"author_name"`
	PublishedAt     *time.Time    `json:"published_at"`
	ExpiredAt       *time.Time    `json:"expired_at"`
	Status          string        `json:"status"`
	ShowOnHomepage  bool          `json:"show_on_homepage"`
	SortOrder       int           `json:"sort_order"`
	Tags            []TagResponse `json:"tags"`
}

type StaticPagePayload struct {
	Title           string `json:"title"`
	ContentMarkdown string `json:"content_markdown"`
	Excerpt         string `json:"excerpt"`
}

type SetupTenantRequest struct {
	Name      string `json:"name"`
	Subdomain string `json:"subdomain"`
}
