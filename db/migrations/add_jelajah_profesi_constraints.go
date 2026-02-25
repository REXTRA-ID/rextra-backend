package migrations

import (
	"gorm.io/gorm"
)

func AddJelajahProfesiConstraints(db *gorm.DB) error {
	// Constraints that cannot be generated automatically by GORM

	// Composite unique sub-category
	if err := db.Exec(`
		ALTER TABLE profession_sub_categories
		DROP CONSTRAINT IF EXISTS uq_sub_category_code;
		ALTER TABLE profession_sub_categories
		ADD CONSTRAINT uq_sub_category_code UNIQUE (main_category_id, code);
		
		ALTER TABLE profession_sub_categories
		DROP CONSTRAINT IF EXISTS uq_sub_category_name;
		ALTER TABLE profession_sub_categories
		ADD CONSTRAINT uq_sub_category_name UNIQUE (main_category_id, name);
	`).Error; err != nil {
		return err
	}

	// Unique alias per profesi
	if err := db.Exec(`
		ALTER TABLE profession_aliases
		DROP CONSTRAINT IF EXISTS uq_alias_per_profession;
		ALTER TABLE profession_aliases
		ADD CONSTRAINT uq_alias_per_profession UNIQUE (profession_id, alias_name);
	`).Error; err != nil {
		return err
	}

	// CHECK valid skill
	if err := db.Exec(`
		ALTER TABLE profession_skill_rels
		DROP CONSTRAINT IF EXISTS chk_skill_type;
		ALTER TABLE profession_skill_rels
		ADD CONSTRAINT chk_skill_type CHECK (skill_type IN ('hard', 'soft'));
		
		ALTER TABLE profession_skill_rels
		DROP CONSTRAINT IF EXISTS chk_priority;
		ALTER TABLE profession_skill_rels
		ADD CONSTRAINT chk_priority CHECK (priority IN ('wajib', 'dianjurkan'));
	`).Error; err != nil {
		return err
	}

	// CHECK valid tool
	if err := db.Exec(`
		ALTER TABLE profession_tool_rels
		DROP CONSTRAINT IF EXISTS chk_usage_type;
		ALTER TABLE profession_tool_rels
		ADD CONSTRAINT chk_usage_type CHECK (usage_type IN ('wajib', 'umum'));
	`).Error; err != nil {
		return err
	}

	return nil
}
