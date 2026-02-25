package seeds

import (
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
)

func SeedProfessionMainCategory(db *gorm.DB) error {
	categories := []entity.ProfessionMainCategory{
		{Code: "ENGINEERING", Name: "Engineering",
			Description: "Kategori profesi teknis yang bertanggung jawab atas perancangan, pengembangan, pengujian, dan pemeliharaan sistem digital."},
		{Code: "DATA", Name: "Data",
			Description: "Kategori profesi analitis yang berfokus pada pengumpulan, pengolahan, dan analisis data untuk menghasilkan wawasan strategis."},
		{Code: "PRODUCT", Name: "Product",
			Description: "Kategori profesi strategis yang menjembatani aspek teknis, bisnis, dan desain dalam pengelolaan produk digital."},
		{Code: "DESIGN", Name: "Design",
			Description: "Kategori profesi kreatif yang berfokus pada perancangan antarmuka dan pengalaman pengguna."},
		{Code: "MARKETING", Name: "Marketing",
			Description: "Kategori profesi komunikasi yang berfokus pada promosi, branding, dan akuisisi pengguna."},
		{Code: "BUSINESS", Name: "Business",
			Description: "Kategori profesi komersial yang berfokus pada strategi pertumbuhan dan kemitraan."},
		{Code: "FINANCE", Name: "Finance",
			Description: "Kategori profesi manajerial yang mengelola kesehatan finansial dan strategi pendanaan perusahaan."},
		{Code: "PEOPLE", Name: "People",
			Description: "Kategori profesi organisasi yang mengelola siklus hidup talenta dan budaya kerja."},
		{Code: "OPERATIONS", Name: "Operations",
			Description: "Kategori profesi eksekusi yang memastikan efisiensi proses bisnis dan operasional."},
		{Code: "LEGAL", Name: "Legal",
			Description: "Kategori profesi hukum yang menangani kepatuhan regulasi dan mitigasi risiko hukum."},
	}

	for _, cat := range categories {
		if err := db.Where("code = ?", cat.Code).FirstOrCreate(&cat).Error; err != nil {
			return err
		}
	}
	return nil
}
