package follow

import (
	"time"

	"github.com/promptist/web/profile"
)

// Relationship shows a follower relationship bewteen a follower and a leader
type Relationship struct {
	ID            int64
	LeaderBrief   profile.Brief
	FollowerBrief profile.Brief
	CreatedAt     time.Time
}

// FollowOption represents the options a user has for following someone
// else
// 0 - Unable to follow or unfollow
// 1 - Can follow
// 2 - Can unfollow

type FollowOption int32
