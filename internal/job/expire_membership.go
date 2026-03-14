package job

import "rextra-backend/internal/modules/membership/service"

type ExpireMembershipJob struct {
	MembershipService service.MembershipService
}

func (j *ExpireMembershipJob) Run() {
	j.MembershipService.CheckExpiredMemberships()
}
