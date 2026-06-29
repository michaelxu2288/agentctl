package orchestration

import (
	"strings"

	"github.com/michaelxu2288/swarmboard/internal/domain"
	"github.com/michaelxu2288/swarmboard/internal/providers"
)

type Router struct {
	PlannerSession  string
	CoderSession    string
	ReviewerSession string
}

func (r Router) RouteTask(task domain.Task) (session string, role providers.Role) {
	title := strings.ToLower(task.Title)
	desc := strings.ToLower(task.Description)

	if strings.Contains(title, "review") || strings.Contains(desc, "review") {
		return r.ReviewerSession, providers.RoleReviewer
	}
	if strings.Contains(title, "plan") || strings.Contains(desc, "decompose") {
		return r.PlannerSession, providers.RolePlanner
	}
	return r.CoderSession, providers.RoleCoder
}
