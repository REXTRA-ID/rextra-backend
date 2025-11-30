package dto_request

const DescriptionPathfinder = "Sebagai Pathfinder, kamu sedang dalam proses mencari dan merencanakan karier digital impian. Kamu perlu mengeksplorasi dan mengenali potensi diri terlebih dahulu. Sobat Rexi, mulai perjalananmu dengan fitur KENALI DIRI, untuk membantu kamu menemukan profesi yang sesuai dengan minat dan bakat"
const DescriptionBuilder = "Sebagai Builder, kamu akan mulai membangun fondasi karier dengan membuat portofolio profesional. Di tahap ini, fokusmu adalah menambah pengalaman, meningkatkan keterampilan, serta memastikan setiap aktivitas dan proyek selaras dengan tujuan karier yang sudah kamu tetapkan"
const DescriptionAchiever = "Sebagai Achiever, kamu semakin dekat dengan karier impian! Di tahap ini, kamu sudah siap menghadapi proses seleksi kerja atau magang. Karir Lab disiapkan oleh REXTRA untuk memaksimalkan persiapan kamu dengan berbagai konten edukasi karier, CV Generator, dan AI Interviewer"

type (
	CreatePersonaRequest struct {
		UserID               string `json:"user_id" binding:"required"`
		HasCareerGoal        *bool  `json:"has_career_goal" binding:"required"`
		BuildingPortofolio   *bool  `json:"building_portofolio" binding:"required"`
		InRecruitmentProcess *bool  `json:"in_recruitment_process" binding:"required"`
	}
	MissionPersonaCompleteRequest struct {
		UserID     string `json:"user_id" binding:"required"`
		MissionKey string `json:"mission_key" binding:"required"`
	}
)
