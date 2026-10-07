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
	HideSpoilers  bool
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
	Rights *RightsInfo
}

type RightsSource string

const (
	RightsPublicDomain    RightsSource = "PUBLIC_DOMAIN"
	RightsCreativeCommons RightsSource = "CREATIVE_COMMONS"
	RightsRevenueShare    RightsSource = "REVENUE_SHARE"
	RightsOwned           RightsSource = "OWNED"
	RightsLicensed        RightsSource = "LICENSED"
	// RightsOfficialEmbed marks titles played through the rights holder's
	// own embeddable upload (e.g. a tournament organizer's YouTube channel).
	// The platform never hosts or paywalls these videos.
	RightsOfficialEmbed RightsSource = "OFFICIAL_EMBED"
)

// RightsInfo tracks who holds the streaming rights to a Media title and for
// how long. It's a 1:1 extension of Media, so it's optional (nil means no
// rights tracked yet) rather than embedded directly.
type RightsInfo struct {
	ID              string
	MediaID         string
	Source          RightsSource
	RightsHolder    *string
	LicenseDocURL   *string
	RevenueSharePct *float64
	TerritoryLimit  *string
	ValidFrom       time.Time
	ValidUntil      *time.Time // nil = unbefristet
	Notes           *string
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Media *Media
}

// Expired reports whether this rights grant has passed its ValidUntil date.
// Unbefristete (ValidUntil == nil) grants never expire.
func (r RightsInfo) Expired() bool {
	return r.ValidUntil != nil && r.ValidUntil.Before(time.Now())
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

	Match *Match
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
	PartIndex       int
	Completed       bool
	UpdatedAt       time.Time
}

type Game struct {
	ID   string
	Name string
	Slug string
}

type Team struct {
	ID        string
	Name      string
	Slug      string
	ShortName *string
	Country   *string
	LogoURL   *string
	CreatedAt time.Time
}

// Match is the optional eSports extension of an Episode: who played, in
// which stage and format, and the result. The result must only be rendered
// for viewers who opted into spoilers (see Profile.HideSpoilers).
type Match struct {
	EpisodeID      string
	GameID         *string
	TeamAID        string
	TeamBID        string
	Stage          string
	BestOf         int
	ScoreA         *int
	ScoreB         *int
	PlayedOn       *time.Time
	ExtraVideoURLs []string

	Game  *Game
	TeamA Team
	TeamB Team
}

// HasResult reports whether both scores are known.
func (m Match) HasResult() bool {
	return m.ScoreA != nil && m.ScoreB != nil
}

// WinnerID returns the winning team's ID, or "" when the result is unknown
// or a draw.
func (m Match) WinnerID() string {
	if !m.HasResult() || *m.ScoreA == *m.ScoreB {
		return ""
	}
	if *m.ScoreA > *m.ScoreB {
		return m.TeamAID
	}
	return m.TeamBID
}

// VideoParts returns all VOD parts in playback order, starting with the
// episode's own video URL.
func (e Episode) VideoParts() []string {
	parts := []string{e.VideoURL}
	if e.Match != nil {
		for _, u := range e.Match.ExtraVideoURLs {
			if u != "" {
				parts = append(parts, u)
			}
		}
	}
	return parts
}

// TeamMatch is a match row joined with where it lives in the catalog, for
// team pages and "followed teams" rows.
type TeamMatch struct {
	Match        Match
	EpisodeTitle string
	MediaTitle   string
	MediaSlug    string
	ThumbnailURL string
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
