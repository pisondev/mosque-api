BEGIN;

-- Mengisi Feature Catalog bawaan sistem
INSERT INTO feature_catalog (feature_type, name, category_label) VALUES 
('facility', 'Tempat Wudhu Pria', 'Sanitasi'),
('facility', 'Tempat Wudhu Wanita', 'Sanitasi'),
('facility', 'Toilet Difabel', 'Sanitasi'),
('facility', 'Area Parkir Luas', 'Infrastruktur'),
('facility', 'Ruang AC', 'Kenyamanan'),
('facility', 'Perpustakaan', 'Pendidikan'),
('service', 'Kajian Rutin', 'Pendidikan'),
('service', 'Layanan Jenazah', 'Sosial'),
('service', 'Ambulans Gratis', 'Sosial'),
('service', 'Taman Pendidikan Al-Qur''an (TPA)', 'Pendidikan')
ON CONFLICT (name) DO NOTHING;

COMMIT;