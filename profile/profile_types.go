package profile

import "database/sql"

// Profile stores general information about a user
type Profile struct {
	ID             int64
	UserID         int64
	Username       string
	FullName       string
	Bio            sql.NullString
	Location       sql.NullString
	PhotoUUID      sql.NullString
	NumFollowing   int64
	NumFollowers   int64
	NumCollections int64
	NumPieces      int64
	NumLikes       int64
	NumViews       int64
	NumComments    int64
	IsNew          bool
}

// Brief is used around the site for smaller things
type Brief struct {
	ID        int64
	UserID    int64
	Username  string
	FullName  string
	PhotoUUID sql.NullString
}
