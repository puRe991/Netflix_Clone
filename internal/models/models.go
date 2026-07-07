// Package models defines the domain types shared across the store and http layers.
package models

import "time"

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

type SubscriptionStatus string

const (
	SubStatusFree       SubscriptionStatus = "FREE"
	SubStatusTrialing   SubscriptionStatus = "TRIALING"
	SubStatusActive     SubscriptionStatus = "ACTIVE"
	SubStatusPastDue    SubscriptionStatus = "PAST_DUE"
	SubStatusCanceled   SubscriptionStatus = "CANCELED"
	SubStatusIncomplete SubscriptionStatus = "INCOMPLETE"
)

func (s SubscriptionStatus) CanStream() bool {
	return s == SubStatusActive || s == SubStatusTrialing
}

type MediaType string

const (
	MediaMovie  MediaType = "MOVIE"
	MediaSeries MediaType = "SERIES"
)

type Plan string

const (
	PlanFree    Plan = "FREE"
	PlanBasic   Plan = "BASIC"
	PlanPremium Plan = "PREMIUM"
)

type User struct {
	ID                 string
	Email              string
	PasswordHash       string
	Name               *string
	Role               Role
	StripeCustomerID   *string
	SubscriptionStatus SubscriptionStatus
	LastProfileID      *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (u User) CanStreamFullContent() bool {
	return u.Role == RoleAdmin || u.SubscriptionStatus.CanStream()
}

type Profile struct {
	ID            string
	UserID        string
	Name          string
	AvatarURL     *string
	IsKidsProfile bool
	CreatedAt     time.Time
}

type Media struct {
	ID           string
	Title        string
	Slug         string
	Description  string
	Type         MediaType
	ThumbnailURL string
	BannerURL    string
	TrailerURL   *string
	VideoURL     *string
	ReleaseYear  int
	Duration     *int
	AgeRating    string
	Language     string
	Tags         []string
	IsPublished  bool
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Genres []Genre
	Series *Series
}

type Series struct {
	ID          string
	MediaID     string
	Title       string
	Description string

	Seasons []Season
	Media   *Media
}

type Season struct {
	ID           string
	SeriesID     string
	SeasonNumber int
	Title        string

	Episodes []Episode
}

type Episode struct {
	ID            string
	SeasonID      string
	Title         string
	Description   string
	EpisodeNumber int
	Duration      int
	VideoURL      string
	ThumbnailURL  string
	IsPublished   bool
}

type Genre struct {
	ID   string
	Name string
	Slug string
}

type WatchProgress struct {
	ID              string
	UserID          string
	ProfileID       string
	MediaID         string
	EpisodeID       *string
	ProgressSeconds int
	Completed       bool
	UpdatedAt       time.Time
}

type WatchlistItem struct {
	ID        string
	UserID    string
	ProfileID string
	MediaID   string
	CreatedAt time.Time

	Media *Media
}

type Subscription struct {
	ID                   string
	UserID               string
	StripeSubscriptionID string
	Plan                 Plan
	Status               SubscriptionStatus
	CurrentPeriodEnd     *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time

	UserEmail *string
	UserName  *string
}

type PaymentLog struct {
	ID            string
	UserID        *string
	StripeEventID string
	Type          string
	RawData       []byte
	CreatedAt     time.Time
}

type Stats struct {
	Users               int64
	ActiveSubscriptions int64
	Media               int64
	WatchtimeSeconds    int64
}
