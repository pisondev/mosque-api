package community

type EventPayload struct {
	Title                string  `json:"title"`
	Category             string  `json:"category"`
	SpeakerName          *string `json:"speaker_name"`
	PersonInCharge       *string `json:"person_in_charge"`
	Description          *string `json:"description"`
	NotePublic           *string `json:"note_public"`
	NoteInternal         *string `json:"note_internal"`
	StartDate            string  `json:"start_date"`
	EndDate              *string `json:"end_date"`
	TimeMode             string  `json:"time_mode"`
	StartTime            *string `json:"start_time"`
	EndTime              *string `json:"end_time"`
	AfterPrayer          *string `json:"after_prayer"`
	AfterPrayerOffsetMin *int16  `json:"after_prayer_offset_min"`
	RepeatPattern        *string `json:"repeat_pattern"`
	RepeatWeekdays       []int16 `json:"repeat_weekdays"`
	Audience             *string `json:"audience"`
	Capacity             *int    `json:"capacity"`
	FeeType              *string `json:"fee_type"`
	FeeAmount            *string `json:"fee_amount"`
	ContactPhone         *string `json:"contact_phone"`
	LocationInside       *string `json:"location_inside"`
	LocationOutside      *string `json:"location_outside"`
	Status               string  `json:"status"`
	PosterImageURL       *string `json:"poster_image_url"`
}

type EventResponse struct {
	ID                   int64   `json:"id"`
	Title                string  `json:"title"`
	Category             string  `json:"category"`
	SpeakerName          *string `json:"speaker_name"`
	PersonInCharge       *string `json:"person_in_charge"`
	Description          *string `json:"description"`
	NotePublic           *string `json:"note_public"`
	NoteInternal         *string `json:"note_internal"`
	StartDate            string  `json:"start_date"`
	EndDate              *string `json:"end_date"`
	TimeMode             string  `json:"time_mode"`
	StartTime            *string `json:"start_time"`
	EndTime              *string `json:"end_time"`
	AfterPrayer          *string `json:"after_prayer"`
	AfterPrayerOffsetMin *int16  `json:"after_prayer_offset_min"`
	RepeatPattern        *string `json:"repeat_pattern"`
	RepeatWeekdays       []int16 `json:"repeat_weekdays"`
	Audience             *string `json:"audience"`
	Capacity             *int    `json:"capacity"`
	FeeType              *string `json:"fee_type"`
	FeeAmount            *string `json:"fee_amount"`
	ContactPhone         *string `json:"contact_phone"`
	LocationInside       *string `json:"location_inside"`
	LocationOutside      *string `json:"location_outside"`
	Status               string  `json:"status"`
	PosterImageURL       *string `json:"poster_image_url"`
}

type EventListQuery struct {
	Status   string
	Category string
	Search   string
	Page     int
	Limit    int
}

type UpdateEventStatusRequest struct {
	Status string `json:"status"`
}

type GalleryAlbumPayload struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	MediaKind   string  `json:"media_kind"`
}

type GalleryAlbumResponse struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	MediaKind   string  `json:"media_kind"`
}

type GalleryItemPayload struct {
	AlbumID      *int64  `json:"album_id"`
	MediaType    string  `json:"media_type"`
	MediaURL     string  `json:"media_url"`
	Caption      *string `json:"caption"`
	TakenAt      *string `json:"taken_at"`
	LocationNote *string `json:"location_note"`
	IsHighlight  bool    `json:"is_highlight"`
	SortOrder    int     `json:"sort_order"`
}

type GalleryItemResponse struct {
	ID           int64   `json:"id"`
	AlbumID      *int64  `json:"album_id"`
	MediaType    string  `json:"media_type"`
	MediaURL     string  `json:"media_url"`
	Caption      *string `json:"caption"`
	TakenAt      *string `json:"taken_at"`
	LocationNote *string `json:"location_note"`
	IsHighlight  bool    `json:"is_highlight"`
	SortOrder    int     `json:"sort_order"`
}

type ManagementMemberPayload struct {
	FullName        string  `json:"full_name"`
	RoleTitle       string  `json:"role_title"`
	PhoneWhatsapp   *string `json:"phone_whatsapp"`
	ProfileImageURL *string `json:"profile_image_url"`
	ShowPublic      bool    `json:"show_public"`
	SortOrder       int     `json:"sort_order"`
}

type ManagementMemberResponse struct {
	ID              int64   `json:"id"`
	FullName        string  `json:"full_name"`
	RoleTitle       string  `json:"role_title"`
	PhoneWhatsapp   *string `json:"phone_whatsapp"`
	ProfileImageURL *string `json:"profile_image_url"`
	ShowPublic      bool    `json:"show_public"`
	SortOrder       int     `json:"sort_order"`
}

type BaseListQuery struct {
	Page   int
	Limit  int
	Search string
}
