package feedback

const (
	TestCategoryCareerProfile = "CAREER_PROFILE"

	SortBySubmittedAt = "submitted_at"
	SortByName        = "name"

	SortDirAsc  = "asc"
	SortDirDesc = "desc"

	RespondentTypeStudent = "STUDENT"
	RespondentTypeExpert  = "EXPERT"
)

type Top5Status string

const (
	Top5StatusP1         Top5Status = "P1"
	Top5StatusP2         Top5Status = "P2"
	Top5StatusP3_5       Top5Status = "P3_5"
	Top5StatusNotPresent Top5Status = "NOT_PRESENT"
)
