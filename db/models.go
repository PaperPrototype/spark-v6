package db

import "time"

func migrate() {
	gormDB.AutoMigrate(
		// auth
		&User{},
		&Session{},

		// course
		&Course{},
		&Release{},
		&Version{},
		&Section{},
		&Content{},
		&Media{},
		&MediaChunk{},

		// purchases
		&Purchase{},
		&BuyRelease{},

		// posts
		&Post{},
		&PostToRelease{},
	)
}

type BuyRelease struct {
	ID           string
	ReleaseID    uint64
	UserID       uint64
	AmountPaying uint16
	ExpiresAt    time.Time
}

type Purchase struct {
	ID     uint64 `gorm:"primaryKey"`
	UserID uint64 `gorm:"not null"`

	CreatedAt     time.Time `gorm:"not null"`
	AmountPaid    uint16    `gorm:"default 0"`
	PercentageDue float32   `gorm:"default 0"`

	// a specific course release
	ReleaseID uint64

	// not a required parameter but used to keep track of version user is currently taking
	// also set to newest version when user first buys a course
	VersionID uint64
}

type User struct {
	ID       uint64 `gorm:"primaryKey"`
	Username string `sql:"UNIQUE"` // unique identifer used in the url
	Name     string // real name
	Hash     string // password hash
	Email    string `sql:"UNIQUE"`
	Bio      string

	Purchases []Purchase `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Session struct {
	TokenUUID string `sql:"UNIQUE" gorm:"primaryKey"` // this is the session id
	DeleteAt  time.Time
	UserID    uint64 `gorm:"not null"`
}

/* COURSE */
type Course struct {
	ID    uint64 `gorm:"primaryKey"`
	Title string `gorm:"not null"`         // a short title of the course
	Name  string `gorm:"UNIQUE; NOT NULL"` // the courses unique url name (eg. spark.com/minecraftcourse)
	Desc  string

	UserID uint64 `gorm:"not null"`

	Releases []Release `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Versions []Version `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Release struct {
	ID       uint64 `gorm:"primaryKey"`
	Price    uint16 `gorm:"default:0"`
	Num      uint16 `gorm:"default:0"`
	Desc     string
	CourseID uint64 `gorm:"not null"`

	Versions []Version `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Version struct {
	ID        uint64 `gorm:"primaryKey"`
	Num       uint16
	Patch     uint16 `gorm:"default:0"`
	CourseID  uint64 `gorm:"not null"`
	ReleaseID uint64 `gorm:"not null"`
	Error     string `gorm:"default:null"`

	Sections []Section `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Section struct {
	ID   uint64 `gorm:"primaryKey"`
	Name string
	// version this section is connected to
	VersionID uint64 `gorm:"not null"`
	ParentID  uint64 // parent section ID

	// children contents
	// special ORM parameter that can be preloaded with data
	Contents []Content `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Content struct {
	ID        uint64 `gorm:"primaryKey"`
	Language  string
	SectionID uint64 `gorm:"not null"`
	Markdown  string
}

// media files for the course
type Media struct {
	ID        uint64 `gorm:"primaryKey"`
	VersionID uint64 `gorm:"not null"`
	Name      string
	Length    uint32
	Type      string // the "type" of file (.zip .png .gif)

	MediaChunks []MediaChunk `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type MediaChunk struct {
	MediaID uint64 `gorm:"not null"`
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
	UserID    uint64    `gorm:"not null"`
	Markdown  string

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UserMedia struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64
	Name      string
	Length    uint32
	Type      string // the "type" of file (.zip .png .gif)
	CreatedAt time.Time

	UserMediaChunks []UserMediaChunk `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UserMediaChunk struct {
	UserMediaID uint64 `gorm:"not null"`
	Data        []byte
	Position    uint16
}

type Channel struct {
	ID       uint64 `gorm:"primaryKey"`
	CourseID uint64 `gorm:"not null"`
	Name     string

	Course   Course    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Messages []Message `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Message struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64 `gorm:"not null"`
	ChannelID uint64 `gorm:"not null"`
	CreatedAt time.Time
	Markdown  string
}

/* RELATIONS */
// relate posts to a course release
type PostToRelease struct {
	PostID    uint64 `gorm:"not null"`
	ReleaseID uint64 `gorm:"not null"`
	SectionID uint64
}

// allow for "thread-like" conversations to continue from messages
type Thread struct {
	// the new channel ID
	ChannelID uint64 `gorm:"not null"`

	// the parent message
	ParentMessageID uint64 `gorm:"not null"`
	// the parent channel ID
	ParentChannelID uint64 `gorm:"not null"`
}
