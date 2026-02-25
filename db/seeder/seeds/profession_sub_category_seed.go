package seeds

import (
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

func SeedProfessionSubCategory(db *gorm.DB) error {
	var engineering, data, product, design entity.ProfessionMainCategory

	if err := db.Where("code = ?", "ENGINEERING").First(&engineering).Error; err != nil {
		return err
	}
	if err := db.Where("code = ?", "DATA").First(&data).Error; err != nil {
		return err
	}
	if err := db.Where("code = ?", "PRODUCT").First(&product).Error; err != nil {
		return err
	}
	if err := db.Where("code = ?", "DESIGN").First(&design).Error; err != nil {
		return err
	}

	subCategories := []entity.ProfessionSubCategory{
		// ── ENGINEERING ──────────────────────────────────────────────────────
		{MainCategoryID: engineering.ID, Code: "BACKEND", Name: "Backend",
			Description: "Profesi engineering yang membangun dan merawat sisi server aplikasi, termasuk logika bisnis, basis data, dan API."},
		{MainCategoryID: engineering.ID, Code: "FRONTEND", Name: "Frontend",
			Description: "Profesi engineering yang membangun antarmuka aplikasi di sisi pengguna dan menerjemahkan rancangan tampilan menjadi kode interaktif dan responsif."},
		{MainCategoryID: engineering.ID, Code: "MOBILE", Name: "Mobile",
			Description: "Profesi engineering yang mengembangkan aplikasi pada perangkat bergerak (Android/iOS) serta menjaga kompatibilitas dan performa lintas perangkat."},
		{MainCategoryID: engineering.ID, Code: "DEVOPS", Name: "DevOps",
			Description: "Profesi engineering yang mengelola proses rilis dan penyebaran aplikasi, termasuk otomasi CI/CD dan pemantauan stabilitas sistem."},
		{MainCategoryID: engineering.ID, Code: "QUALITY_ASSURANCE", Name: "Quality Assurance",
			Description: "Profesi engineering yang memverifikasi kualitas perangkat lunak sebelum rilis melalui pengujian fungsional dan otomatis."},
		{MainCategoryID: engineering.ID, Code: "SECURITY", Name: "Security",
			Description: "Profesi engineering yang melindungi sistem dan data dari ancaman siber melalui penilaian kerentanan dan respons insiden."},
		{MainCategoryID: engineering.ID, Code: "NETWORK", Name: "Network",
			Description: "Profesi engineering yang merancang dan memelihara jaringan komunikasi agar konektivitas stabil dan aman."},
		{MainCategoryID: engineering.ID, Code: "GAME", Name: "Game",
			Description: "Profesi engineering yang membangun sistem permainan digital, termasuk logika interaksi dan integrasi aset grafis/audio."},
		{MainCategoryID: engineering.ID, Code: "EMBEDDED", Name: "Embedded",
			Description: "Profesi engineering yang mengembangkan perangkat lunak pada sistem tertanam (firmware) dan perangkat IoT."},
		{MainCategoryID: engineering.ID, Code: "ROBOTICS", Name: "Robotics",
			Description: "Profesi engineering yang merancang sistem robotik, termasuk algoritma kontrol dan integrasi sensor untuk otomasi industri."},
		{MainCategoryID: engineering.ID, Code: "CLOUD", Name: "Cloud",
			Description: "Profesi engineering yang merancang dan mengelola layanan berbasis cloud, mencakup arsitektur, skalabilitas, dan efisiensi biaya."},
		{MainCategoryID: engineering.ID, Code: "AUTOMATION", Name: "Automation",
			Description: "Profesi engineering yang membangun otomasi proses teknis untuk mengurangi human error dan mempercepat waktu rilis."},
		{MainCategoryID: engineering.ID, Code: "HARDWARE", Name: "Hardware",
			Description: "Profesi engineering yang merancang dan mengembangkan perangkat keras untuk produk teknologi, termasuk pengujian prototipe dan integrasi software."},

		// ── DATA ─────────────────────────────────────────────────────────────
		{MainCategoryID: data.ID, Code: "INFRASTRUCTURE", Name: "Infrastructure",
			Description: "Kategori profesi yang berfokus pada pembangunan dan pengelolaan infrastruktur data, termasuk data pipeline dan penyimpanan data warehouse/data lake."},
		{MainCategoryID: data.ID, Code: "GOVERNANCE", Name: "Governance",
			Description: "Kategori profesi yang berfokus pada pengelolaan kualitas data, kepatuhan regulasi, dan keamanan data."},
		{MainCategoryID: data.ID, Code: "ANALYTICS", Name: "Analytics",
			Description: "Kategori profesi yang berfokus pada analisis data historis untuk menghasilkan wawasan pengambilan keputusan bisnis."},
		{MainCategoryID: data.ID, Code: "SCIENCE", Name: "Science",
			Description: "Kategori profesi yang berfokus pada prediksi berbasis data menggunakan teknik statistik dan machine learning."},
		{MainCategoryID: data.ID, Code: "AI", Name: "AI",
			Description: "Kategori profesi yang berfokus pada pengembangan sistem machine learning dan kemampuan kognitif buatan."},

		// ── PRODUCT ──────────────────────────────────────────────────────────
		{MainCategoryID: product.ID, Code: "MANAGEMENT", Name: "Management",
			Description: "Kategori profesi yang bertanggung jawab atas visi, strategi, prioritas fitur, dan eksekusi peluncuran produk."},
		{MainCategoryID: product.ID, Code: "GROWTH", Name: "Growth",
			Description: "Kategori profesi yang fokus pada akuisisi pengguna, retensi, dan optimasi produk untuk meningkatkan konversi dan pendapatan."},
		{MainCategoryID: product.ID, Code: "OPERATIONS", Name: "Operations",
			Description: "Kategori profesi yang berfokus pada efisiensi proses dan ritme kerja dalam tim produk."},
		{MainCategoryID: product.ID, Code: "ANALYTICS", Name: "Analytics",
			Description: "Kategori profesi yang bertanggung jawab untuk analisis perilaku pengguna dan wawasan data produk."},

		// ── DESIGN ───────────────────────────────────────────────────────────
		{MainCategoryID: design.ID, Code: "UI_UX", Name: "UI/UX",
			Description: "Subkategori profesi yang bertanggung jawab atas perancangan antarmuka pengguna (UI) dan pengalaman pengguna (UX)."},
		{MainCategoryID: design.ID, Code: "VISUAL", Name: "Visual",
			Description: "Subkategori yang berfokus pada komunikasi visual, branding, dan desain grafis untuk memperkuat identitas merek digital."},
		{MainCategoryID: design.ID, Code: "RESEARCH", Name: "Research",
			Description: "Subkategori yang bertanggung jawab untuk meneliti perilaku pengguna melalui wawancara, survei, dan usability testing."},
		{MainCategoryID: design.ID, Code: "CONTENT", Name: "Content",
			Description: "Subkategori yang berfokus pada penulisan microcopy dan teks dalam produk digital untuk membantu navigasi pengguna."},
	}

	for _, sc := range subCategories {
		if err := db.Where("main_category_id = ? AND code = ?", sc.MainCategoryID, sc.Code).FirstOrCreate(&sc).Error; err != nil {
			return err
		}
	}
	return nil
}
