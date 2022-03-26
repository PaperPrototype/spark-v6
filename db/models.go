package db

import "time"

func migrate() {
	gormDB.AutoMigrate(
		// auth
		&User{},
		&Session{},

		// purchases
		&Purchase{},

		// course
		&Course{},
		&Release{},
		&Version{},
		&Section{},
		&Content{},
		&Media{},
		&MediaChunk{},

		// posts
		&Post{},
		&PostToRelease{},
	)
}

type Purchase struct {
	ID     uint64 `gorm:"primaryKey"`
	UserID uint64

	// a specific course release
	ReleaseID  uint64
	CreatedAt  time.Time
	AmountPaid uint16
}

type User struct {
	ID       uint64 `gorm:"primaryKey"`
	Username string `sql:"UNIQUE"` // unique identifer used in the url
	Name     string // real name
	Hash     string // password hash
	Email    string `sql:"UNIQUE"`
	Bio      string
}

type Session struct {
	TokenUUID string `sql:"UNIQUE" gorm:"primaryKey"` // this is the session id
	DeleteAt  time.Time
	UserID    uint64 `sql:"REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`
}

/* COURSE */
type Course struct {
	ID    uint64 `gorm:"primaryKey"`
	Title string `sql:"NOT NULL"`        // a short title of the course
	Name  string `sql:"UNIQUE NOT NULL"` // the courses unique url name (eg. spark.com/minecraftcourse)
	Desc  string

	UserID uint64 `sql:"REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`
}

type Release struct {
	ID       uint64 `gorm:"primaryKey"`
	Price    uint16 `sql:"DEFAULT 0"`
	Num      uint16
	Desc     string
	CourseID uint64 `sql:"REFERENCES courses(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`
}

type Version struct {
	ID        uint64 `gorm:"primaryKey"`
	Num       uint16
	CourseID  uint64 `sql:"REFERENCES courses(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`
	ReleaseID uint64 `sql:"REFERENCES releases(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`
}

type Section struct {
	ID   uint64 `gorm:"primaryKey"`
	Name string
	// version this section is connected to
	VersionID uint64 `sql:"REFERENCES versions(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`
	ParentID  uint64 // parent section ID

	// children contents
	// special ORM parameter that can be preloaded with data
	Contents []Content
}

type Content struct {
	ID        uint64 `gorm:"primaryKey"`
	Language  string
	SectionID uint64 `sql:"REFERENCES sections(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`
	Markdown  string
}

// media files for the course
type Media struct {
	ID          uint64 `gorm:"primaryKey"`
	VersionID   uint64 `sql:"REFERENCES versions(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`
	Name        string
	Length      uint32
	Type        string // the "type" of file (.zip .png .gif)
	MediaChunks []MediaChunk
}

type MediaChunk struct {
	MediaID uint64 `sql:"REFERENCES media(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`
	// the raw bytes from the file
	Data []byte
	// Order to load chunks from db
	Position uint16
}

/* SOCIAL */
type Post struct {
	ID        uint64    `gorm:"primaryKey"`
	CreatedAt time.Time // special param name gorm automaically sets time
	UpdatedAt time.Time // special param name gorm automaically sets time
	UserID    uint64    `sql:"REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`
	Markdown  string

	User User
}

type Channel struct {
	ID       uint64 `gorm:"primaryKey"`
	CourseID uint64 `gorm:"REFERENCES courses(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`
	Name     string
}

type Message struct {
	UserID    uint64 `sql:"REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`
	ChannelID uint64 `sql:"REFERENCES channels(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`
	CreatedAt time.Time
	Markdown  string
}

/* RELATIONS */
// relate posts to a course release
type PostToRelease struct {
	PostID    uint64 `sql:"REFERENCES posts(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`
	ReleaseID uint64 `sql:"REFERENCES releases(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`
	SectionID uint64
}

// allow for "thread-like" conversations to continue from messages
type Thread struct {
	// the new channel ID
	ChannelID uint64 `sql:"REFERENCES channels(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`

	// the parent message
	ParentMessageID uint64 `sql:"REFERENCES messages(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`
	// the parent channel ID
	ParentChannelID uint64 `sql:"REFERENCES channels(id) ON UPDATE CASCADE ON DELETE CASCADE; NOT NULL"`
}
