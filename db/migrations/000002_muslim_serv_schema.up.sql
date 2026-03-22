DO $$ BEGIN
  CREATE TYPE place_kind AS ENUM ('masjid','musala','surau','langgar');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE parking_type AS ENUM ('none','motor_only','car_and_motor');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE domain_type AS ENUM ('subdomain','custom_domain');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE domain_status AS ENUM ('pending','verifying','active','disabled');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE prayer_calc_location_mode AS ENUM ('city','coordinates');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE prayer_name AS ENUM ('subuh','dzuhur','ashar','maghrib','isya');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE duty_category AS ENUM ('jumat','fardhu','tarawih','id');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE event_time_mode AS ENUM ('exact_time','after_prayer');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE event_status AS ENUM ('upcoming','ongoing','finished','cancelled');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE audience_target AS ENUM ('umum','pria','wanita','remaja','anak','pengurus');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE post_category AS ENUM ('news_activity','announcement','reflection','static_page');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE post_status AS ENUM ('draft','published','archived');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE gallery_media_type AS ENUM ('image','video');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE album_media_kind AS ENUM ('photo','video','mix');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE payment_channel_type AS ENUM ('bank_account','qris');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE feature_type AS ENUM ('facility','service');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE special_day_kind AS ENUM ('idul_fitri','idul_adha','other');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE event_category AS ENUM ('kajian_rutin','tabligh_akbar','rapat_pengurus','kegiatan_sosial','phbi','pelatihan','lainnya');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

DO $$ BEGIN
  CREATE TYPE tag_scope AS ENUM ('post','gallery','event');
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

CREATE TABLE IF NOT EXISTS website_domains (
  id BIGSERIAL PRIMARY KEY,
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  domain_type domain_type NOT NULL,
  hostname TEXT NOT NULL,
  status domain_status NOT NULL DEFAULT 'pending',
  verified_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT uq_website_domains_hostname UNIQUE (hostname)
);

CREATE INDEX IF NOT EXISTS idx_website_domains_tenant_id ON website_domains(tenant_id);
CREATE INDEX IF NOT EXISTS idx_website_domains_status ON website_domains(status);

CREATE TABLE IF NOT EXISTS masjid_profiles (
  id BIGSERIAL PRIMARY KEY,
  tenant_id UUID NOT NULL UNIQUE REFERENCES tenants(id) ON DELETE CASCADE,
  official_name TEXT NOT NULL,
  short_name TEXT NULL,
  header_image_url TEXT NULL,
  kind place_kind NOT NULL,
  category_label TEXT NULL,
  ownership_status TEXT NULL,
  established_year SMALLINT NULL CHECK (established_year BETWEEN 1000 AND 3000),
  short_description TEXT NULL,
  country TEXT NOT NULL DEFAULT 'Indonesia',
  province TEXT NULL,
  city TEXT NULL,
  district TEXT NULL,
  village TEXT NULL,
  postal_code TEXT NULL,
  address_full TEXT NULL,
  landmark TEXT NULL,
  google_maps_url TEXT NULL,
  latitude NUMERIC(10,7) NULL,
  longitude NUMERIC(10,7) NULL,
  phone_whatsapp TEXT NULL,
  office_phone TEXT NULL,
  email TEXT NULL,
  office_days SMALLINT[] NULL,
  office_open_time TIME NULL,
  office_close_time TIME NULL,
  building_area_m2 NUMERIC(12,2) NULL,
  land_area_m2 NUMERIC(12,2) NULL,
  capacity_male INTEGER NULL CHECK (capacity_male >= 0),
  capacity_female INTEGER NULL CHECK (capacity_female >= 0),
  floors_count INTEGER NULL CHECK (floors_count >= 0),
  parking parking_type NULL,
  disability_access BOOLEAN NULL,
  disability_note TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_masjid_profiles_city ON masjid_profiles(city);

CREATE TABLE IF NOT EXISTS prayer_time_settings (
  id BIGSERIAL PRIMARY KEY,
  tenant_id UUID NOT NULL UNIQUE REFERENCES tenants(id) ON DELETE CASCADE,
  timezone TEXT NOT NULL DEFAULT 'Asia/Jakarta',
  location_mode prayer_calc_location_mode NOT NULL DEFAULT 'city',
  city_name TEXT NULL,
  latitude NUMERIC(10,7) NULL,
  longitude NUMERIC(10,7) NULL,
  calc_method TEXT NULL,
  asr_madhhab TEXT NULL,
  adj_subuh_min SMALLINT NOT NULL DEFAULT 0,
  adj_dzuhur_min SMALLINT NOT NULL DEFAULT 0,
  adj_ashar_min SMALLINT NOT NULL DEFAULT 0,
  adj_maghrib_min SMALLINT NOT NULL DEFAULT 0,
  adj_isya_min SMALLINT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_prayer_setting_location_city
    CHECK (location_mode <> 'city' OR city_name IS NOT NULL),
  CONSTRAINT chk_prayer_setting_location_coord
    CHECK (location_mode <> 'coordinates' OR (latitude IS NOT NULL AND longitude IS NOT NULL))
);

CREATE TABLE IF NOT EXISTS prayer_times_daily (
  id BIGSERIAL PRIMARY KEY,
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  day_date DATE NOT NULL,
  subuh_time TIME NOT NULL,
  dzuhur_time TIME NOT NULL,
  ashar_time TIME NOT NULL,
  maghrib_time TIME NOT NULL,
  isya_time TIME NOT NULL,
  sunrise_time TIME NULL,
  dhuha_time TIME NULL,
  source_label TEXT NULL,
  fetched_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT uq_prayer_times_daily UNIQUE (tenant_id, day_date)
);

CREATE INDEX IF NOT EXISTS idx_prayer_times_daily_tenant_date ON prayer_times_daily(tenant_id, day_date);

CREATE TABLE IF NOT EXISTS prayer_duties (
  id BIGSERIAL PRIMARY KEY,
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  category duty_category NOT NULL,
  duty_date DATE NOT NULL,
  prayer prayer_name NULL,
  khatib_name TEXT NULL,
  imam_name TEXT NULL,
  muadzin_name TEXT NULL,
  first_adhan_time TIME NULL,
  khutbah_start_time TIME NULL,
  khutbah_topic TEXT NULL,
  note TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_duty_fardhu_prayer
    CHECK (category <> 'fardhu' OR prayer IS NOT NULL)
);

CREATE INDEX IF NOT EXISTS idx_prayer_duties_tenant_date ON prayer_duties(tenant_id, duty_date);

CREATE TABLE IF NOT EXISTS special_days (
  id BIGSERIAL PRIMARY KEY,
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  kind special_day_kind NOT NULL,
  title TEXT NOT NULL,
  day_date DATE NOT NULL,
  location_note TEXT NULL,
  start_time TIME NULL,
  note TEXT NULL,
  imam_name TEXT NULL,
  khatib_name TEXT NULL,
  muadzin_name TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT uq_special_days UNIQUE (tenant_id, kind, day_date)
);

CREATE INDEX IF NOT EXISTS idx_special_days_tenant_date ON special_days(tenant_id, day_date);

CREATE TABLE IF NOT EXISTS events (
  id BIGSERIAL PRIMARY KEY,
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  category event_category NOT NULL DEFAULT 'lainnya',
  speaker_name TEXT NULL,
  person_in_charge TEXT NULL,
  description TEXT NULL,
  note_public TEXT NULL,
  note_internal TEXT NULL,
  start_date DATE NOT NULL,
  end_date DATE NULL,
  time_mode event_time_mode NOT NULL DEFAULT 'exact_time',
  start_time TIME NULL,
  end_time TIME NULL,
  after_prayer prayer_name NULL,
  after_prayer_offset_min SMALLINT NULL,
  repeat_pattern TEXT NULL,
  repeat_weekdays SMALLINT[] NULL,
  audience audience_target NULL,
  capacity INTEGER NULL CHECK (capacity IS NULL OR capacity >= 0),
  fee_type TEXT NULL,
  fee_amount NUMERIC(12,2) NULL CHECK (fee_amount IS NULL OR fee_amount >= 0),
  contact_phone TEXT NULL,
  location_inside TEXT NULL,
  location_outside TEXT NULL,
  status event_status NOT NULL DEFAULT 'upcoming',
  poster_image_url TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_event_end_date
    CHECK (end_date IS NULL OR end_date >= start_date),
  CONSTRAINT chk_event_time_exact
    CHECK (time_mode <> 'exact_time' OR start_time IS NOT NULL),
  CONSTRAINT chk_event_time_after_prayer
    CHECK (time_mode <> 'after_prayer' OR after_prayer IS NOT NULL)
);

CREATE INDEX IF NOT EXISTS idx_events_tenant_start_date ON events(tenant_id, start_date);

CREATE TABLE IF NOT EXISTS management_members (
  id BIGSERIAL PRIMARY KEY,
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  full_name TEXT NOT NULL,
  role_title TEXT NOT NULL,
  phone_whatsapp TEXT NULL,
  profile_image_url TEXT NULL,
  show_public BOOLEAN NOT NULL DEFAULT TRUE,
  sort_order INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_management_members_tenant_sort ON management_members(tenant_id, sort_order);

CREATE TABLE IF NOT EXISTS feature_catalog (
  id BIGSERIAL PRIMARY KEY,
  feature_type feature_type NOT NULL,
  name TEXT NOT NULL UNIQUE,
  category_label TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS website_features (
  id BIGSERIAL PRIMARY KEY,
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  feature_id BIGINT NOT NULL REFERENCES feature_catalog(id) ON DELETE RESTRICT,
  enabled BOOLEAN NOT NULL DEFAULT FALSE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  detail TEXT NULL,
  note TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT uq_website_feature UNIQUE (tenant_id, feature_id)
);

CREATE INDEX IF NOT EXISTS idx_website_features_tenant ON website_features(tenant_id);

CREATE TABLE IF NOT EXISTS tags (
  id BIGSERIAL PRIMARY KEY,
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  scope tag_scope NOT NULL,
  name TEXT NOT NULL,
  slug TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT uq_tags UNIQUE (tenant_id, scope, slug)
);

CREATE INDEX IF NOT EXISTS idx_tags_tenant_scope ON tags(tenant_id, scope);

CREATE TABLE IF NOT EXISTS posts (
  id BIGSERIAL PRIMARY KEY,
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  slug TEXT NOT NULL,
  category post_category NOT NULL,
  excerpt TEXT NULL,
  content_markdown TEXT NOT NULL,
  thumbnail_url TEXT NULL,
  author_name TEXT NULL,
  published_at TIMESTAMPTZ NULL,
  expired_at TIMESTAMPTZ NULL,
  status post_status NOT NULL DEFAULT 'draft',
  show_on_homepage BOOLEAN NOT NULL DEFAULT FALSE,
  sort_order INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT uq_posts_slug UNIQUE (tenant_id, slug),
  CONSTRAINT chk_posts_expired_after_published
    CHECK (expired_at IS NULL OR published_at IS NULL OR expired_at >= published_at)
);

CREATE INDEX IF NOT EXISTS idx_posts_tenant_category ON posts(tenant_id, category);
CREATE INDEX IF NOT EXISTS idx_posts_tenant_status ON posts(tenant_id, status);

CREATE TABLE IF NOT EXISTS post_tags (
  post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  tag_id BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
  PRIMARY KEY (post_id, tag_id)
);

CREATE TABLE IF NOT EXISTS gallery_albums (
  id BIGSERIAL PRIMARY KEY,
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  description TEXT NULL,
  start_date DATE NULL,
  end_date DATE NULL,
  media_kind album_media_kind NOT NULL DEFAULT 'photo',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_album_end_date CHECK (end_date IS NULL OR start_date IS NULL OR end_date >= start_date)
);

CREATE INDEX IF NOT EXISTS idx_gallery_albums_tenant ON gallery_albums(tenant_id);

CREATE TABLE IF NOT EXISTS gallery_items (
  id BIGSERIAL PRIMARY KEY,
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  album_id BIGINT NULL REFERENCES gallery_albums(id) ON DELETE SET NULL,
  media_type gallery_media_type NOT NULL,
  media_url TEXT NOT NULL,
  caption TEXT NULL,
  taken_at TIMESTAMPTZ NULL,
  location_note TEXT NULL,
  is_highlight BOOLEAN NOT NULL DEFAULT FALSE,
  sort_order INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_gallery_items_tenant_album ON gallery_items(tenant_id, album_id);

CREATE TABLE IF NOT EXISTS gallery_item_tags (
  gallery_item_id BIGINT NOT NULL REFERENCES gallery_items(id) ON DELETE CASCADE,
  tag_id BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
  PRIMARY KEY (gallery_item_id, tag_id)
);

CREATE TABLE IF NOT EXISTS donation_channels (
  id BIGSERIAL PRIMARY KEY,
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  channel_type payment_channel_type NOT NULL,
  label TEXT NOT NULL,
  bank_name TEXT NULL,
  bank_branch TEXT NULL,
  account_number TEXT NULL,
  account_holder_name TEXT NULL,
  qris_image_url TEXT NULL,
  merchant_id TEXT NULL,
  description TEXT NULL,
  sort_order INTEGER NOT NULL DEFAULT 0,
  is_public BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_donation_bank_fields
    CHECK (channel_type <> 'bank_account' OR (bank_name IS NOT NULL AND account_number IS NOT NULL AND account_holder_name IS NOT NULL)),
  CONSTRAINT chk_donation_qris_fields
    CHECK (channel_type <> 'qris' OR qris_image_url IS NOT NULL)
);

CREATE INDEX IF NOT EXISTS idx_donation_channels_tenant ON donation_channels(tenant_id);

CREATE TABLE IF NOT EXISTS social_links (
  id BIGSERIAL PRIMARY KEY,
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  platform TEXT NOT NULL,
  account_name TEXT NULL,
  url TEXT NOT NULL,
  description TEXT NULL,
  show_in_footer BOOLEAN NOT NULL DEFAULT TRUE,
  show_in_contact_page BOOLEAN NOT NULL DEFAULT TRUE,
  sort_order INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_social_links_tenant ON social_links(tenant_id);

CREATE TABLE IF NOT EXISTS external_links (
  id BIGSERIAL PRIMARY KEY,
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  link_type TEXT NOT NULL,
  label TEXT NOT NULL,
  url TEXT NOT NULL,
  note TEXT NULL,
  visibility TEXT NOT NULL DEFAULT 'public',
  sort_order INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_external_links_tenant ON external_links(tenant_id);
