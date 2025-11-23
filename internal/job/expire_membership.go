package job

import "rextra-backend/internal/api/service"

type ExpireMembershipJob struct {
	MembershipService service.MembershipService
}

func (j *ExpireMembershipJob) Run() {
	j.MembershipService.CheckExpiredMemberships()
}
