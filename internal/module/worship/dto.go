package worship

type PrayerTimeSettingsRequest struct {
	Timezone      string   `json:"timezone"`
	LocationMode  string   `json:"location_mode"`
	CityName      *string  `json:"city_name"`
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
	CalcMethod    *string  `json:"calc_method"`
	AsrMadhhab    *string  `json:"asr_madhhab"`
	AdjSubuhMin   int16    `json:"adj_subuh_min"`
	AdjDzuhurMin  int16    `json:"adj_dzuhur_min"`
	AdjAsharMin   int16    `json:"adj_ashar_min"`
	AdjMaghribMin int16    `json:"adj_maghrib_min"`
	AdjIsyaMin    int16    `json:"adj_isya_min"`
}

type PrayerTimeSettingsResponse struct {
	ID            int64    `json:"id"`
	Timezone      string   `json:"timezone"`
	LocationMode  string   `json:"location_mode"`
	CityName      *string  `json:"city_name"`
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
	CalcMethod    *string  `json:"calc_method"`
	AsrMadhhab    *string  `json:"asr_madhhab"`
	AdjSubuhMin   int16    `json:"adj_subuh_min"`
	AdjDzuhurMin  int16    `json:"adj_dzuhur_min"`
	AdjAsharMin   int16    `json:"adj_ashar_min"`
	AdjMaghribMin int16    `json:"adj_maghrib_min"`
	AdjIsyaMin    int16    `json:"adj_isya_min"`
}

type PrayerTimesDailyQuery struct {
	From  string
	To    string
	Page  int
	Limit int
}

type PrayerTimesDailyPayload struct {
	DayDate     string  `json:"day_date"`
	SubuhTime   string  `json:"subuh_time"`
	DzuhurTime  string  `json:"dzuhur_time"`
	AsharTime   string  `json:"ashar_time"`
	MaghribTime string  `json:"maghrib_time"`
	IsyaTime    string  `json:"isya_time"`
	SunriseTime *string `json:"sunrise_time"`
	DhuhaTime   *string `json:"dhuha_time"`
	SourceLabel *string `json:"source_label"`
	FetchedAt   *string `json:"fetched_at"`
}

type PrayerTimesDailyResponse struct {
	ID          int64   `json:"id"`
	DayDate     string  `json:"day_date"`
	SubuhTime   string  `json:"subuh_time"`
	DzuhurTime  string  `json:"dzuhur_time"`
	AsharTime   string  `json:"ashar_time"`
	MaghribTime string  `json:"maghrib_time"`
	IsyaTime    string  `json:"isya_time"`
	SunriseTime *string `json:"sunrise_time"`
	DhuhaTime   *string `json:"dhuha_time"`
	SourceLabel *string `json:"source_label"`
	FetchedAt   *string `json:"fetched_at"`
}

type PrayerDutiesQuery struct {
	From     string
	To       string
	Category string
	Prayer   string
	Page     int
	Limit    int
}

type PrayerDutyPayload struct {
	Category         string  `json:"category"`
	DutyDate         string  `json:"duty_date"`
	Prayer           *string `json:"prayer"`
	KhatibName       *string `json:"khatib_name"`
	ImamName         *string `json:"imam_name"`
	MuadzinName      *string `json:"muadzin_name"`
	FirstAdhanTime   *string `json:"first_adhan_time"`
	KhutbahStartTime *string `json:"khutbah_start_time"`
	KhutbahTopic     *string `json:"khutbah_topic"`
	Note             *string `json:"note"`
}

type PrayerDutyResponse struct {
	ID               int64   `json:"id"`
	Category         string  `json:"category"`
	DutyDate         string  `json:"duty_date"`
	Prayer           *string `json:"prayer"`
	KhatibName       *string `json:"khatib_name"`
	ImamName         *string `json:"imam_name"`
	MuadzinName      *string `json:"muadzin_name"`
	FirstAdhanTime   *string `json:"first_adhan_time"`
	KhutbahStartTime *string `json:"khutbah_start_time"`
	KhutbahTopic     *string `json:"khutbah_topic"`
	Note             *string `json:"note"`
}

type SpecialDaysQuery struct {
	Year  string
	Kind  string
	From  string
	To    string
	Page  int
	Limit int
}

type SpecialDayPayload struct {
	Kind         string  `json:"kind"`
	Title        string  `json:"title"`
	DayDate      string  `json:"day_date"`
	LocationNote *string `json:"location_note"`
	StartTime    *string `json:"start_time"`
	Note         *string `json:"note"`
	ImamName     *string `json:"imam_name"`
	KhatibName   *string `json:"khatib_name"`
	MuadzinName  *string `json:"muadzin_name"`
}

type SpecialDayResponse struct {
	ID           int64   `json:"id"`
	Kind         string  `json:"kind"`
	Title        string  `json:"title"`
	DayDate      string  `json:"day_date"`
	LocationNote *string `json:"location_note"`
	StartTime    *string `json:"start_time"`
	Note         *string `json:"note"`
	ImamName     *string `json:"imam_name"`
	KhatibName   *string `json:"khatib_name"`
	MuadzinName  *string `json:"muadzin_name"`
}
