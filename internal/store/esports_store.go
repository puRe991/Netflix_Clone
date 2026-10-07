package store

import (
	"context"
	"time"

	"github.com/pure991/streamflix/internal/models"
)

// --- Games ---

func (s *Store) ListGames(ctx context.Context) ([]models.Game, error) {
	rows, err := s.Pool.Query(ctx, `select id, name, slug from games order by name asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []models.Game
	for rows.Next() {
		var g models.Game
		if err := rows.Scan(&g.ID, &g.Name, &g.Slug); err != nil {
			return nil, err
		}
		games = append(games, g)
	}
	return games, rows.Err()
}

func (s *Store) EnsureGame(ctx context.Context, name, slug string) (models.Game, error) {
	var g models.Game
	err := s.Pool.QueryRow(ctx, `
		insert into games (name, slug) values ($1, $2)
		on conflict (slug) do update set name = excluded.name
		returning id, name, slug
	`, name, slug).Scan(&g.ID, &g.Name, &g.Slug)
	return g, err
}

// --- Teams ---

const teamColumns = `id, name, slug, short_name, country, logo_url, created_at`

func scanTeam(row interface{ Scan(dest ...any) error }) (models.Team, error) {
	var t models.Team
	err := row.Scan(&t.ID, &t.Name, &t.Slug, &t.ShortName, &t.Country, &t.LogoURL, &t.CreatedAt)
	return t, err
}

func (s *Store) ListTeams(ctx context.Context) ([]models.Team, error) {
	rows, err := s.Pool.Query(ctx, `select `+teamColumns+` from teams order by name asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		t, err := scanTeam(rows)
		if err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, rows.Err()
}

func (s *Store) GetTeamBySlug(ctx context.Context, slug string) (models.Team, error) {
	t, err := scanTeam(s.Pool.QueryRow(ctx, `select `+teamColumns+` from teams where slug = $1`, slug))
	return t, mapNotFound(err)
}

func (s *Store) GetTeamByID(ctx context.Context, id string) (models.Team, error) {
	t, err := scanTeam(s.Pool.QueryRow(ctx, `select `+teamColumns+` from teams where id = $1`, id))
	return t, mapNotFound(err)
}

type TeamInput struct {
	Name      string
	Slug      string
	ShortName *string
	Country   *string
	LogoURL   *string
}

func (s *Store) CreateTeam(ctx context.Context, in TeamInput) (models.Team, error) {
	row := s.Pool.QueryRow(ctx, `
		insert into teams (name, slug, short_name, country, logo_url) values ($1, $2, $3, $4, $5)
		returning `+teamColumns, in.Name, in.Slug, in.ShortName, in.Country, in.LogoURL)
	return scanTeam(row)
}

// EnsureTeam returns the team with in.Slug, creating it if missing (used by
// the seed so it can be re-run safely).
func (s *Store) EnsureTeam(ctx context.Context, in TeamInput) (models.Team, error) {
	t, err := s.GetTeamBySlug(ctx, in.Slug)
	if err == ErrNotFound {
		return s.CreateTeam(ctx, in)
	}
	return t, err
}

func (s *Store) DeleteTeam(ctx context.Context, id string) error {
	tag, err := s.Pool.Exec(ctx, `delete from teams where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// TeamHasMatches reports whether a team is referenced by any match (teams
// with matches can't be deleted, the FK is ON DELETE RESTRICT).
func (s *Store) TeamHasMatches(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := s.Pool.QueryRow(ctx, `select exists(select 1 from matches where team_a_id = $1 or team_b_id = $1)`, id).Scan(&exists)
	return exists, err
}

// --- Matches ---

const matchSelect = `
	select mt.episode_id, mt.game_id, mt.team_a_id, mt.team_b_id, mt.stage, mt.best_of,
		mt.score_a, mt.score_b, mt.played_on, mt.extra_video_urls,
		g.id, g.name, g.slug,
		ta.id, ta.name, ta.slug, ta.short_name, ta.country, ta.logo_url, ta.created_at,
		tb.id, tb.name, tb.slug, tb.short_name, tb.country, tb.logo_url, tb.created_at
	from matches mt
	left join games g on g.id = mt.game_id
	join teams ta on ta.id = mt.team_a_id
	join teams tb on tb.id = mt.team_b_id`

func scanMatch(row interface{ Scan(dest ...any) error }, extra ...any) (models.Match, error) {
	var m models.Match
	var gameID, gameName, gameSlug *string
	dest := []any{&m.EpisodeID, &m.GameID, &m.TeamAID, &m.TeamBID, &m.Stage, &m.BestOf,
		&m.ScoreA, &m.ScoreB, &m.PlayedOn, &m.ExtraVideoURLs,
		&gameID, &gameName, &gameSlug,
		&m.TeamA.ID, &m.TeamA.Name, &m.TeamA.Slug, &m.TeamA.ShortName, &m.TeamA.Country, &m.TeamA.LogoURL, &m.TeamA.CreatedAt,
		&m.TeamB.ID, &m.TeamB.Name, &m.TeamB.Slug, &m.TeamB.ShortName, &m.TeamB.Country, &m.TeamB.LogoURL, &m.TeamB.CreatedAt,
	}
	if err := row.Scan(append(dest, extra...)...); err != nil {
		return m, err
	}
	if gameID != nil {
		m.Game = &models.Game{ID: *gameID, Name: *gameName, Slug: *gameSlug}
	}
	return m, nil
}

func (s *Store) GetMatchByEpisodeID(ctx context.Context, episodeID string) (models.Match, error) {
	m, err := scanMatch(s.Pool.QueryRow(ctx, matchSelect+` where mt.episode_id = $1`, episodeID))
	return m, mapNotFound(err)
}

// attachMatches loads match metadata for the given episodes in one query
// and sets Episode.Match where present.
func (s *Store) attachMatches(ctx context.Context, episodes []models.Episode) error {
	if len(episodes) == 0 {
		return nil
	}
	ids := make([]string, len(episodes))
	for i, e := range episodes {
		ids[i] = e.ID
	}

	rows, err := s.Pool.Query(ctx, matchSelect+` where mt.episode_id = any($1::uuid[])`, ids)
	if err != nil {
		return err
	}
	defer rows.Close()

	byEpisode := make(map[string]models.Match)
	for rows.Next() {
		m, err := scanMatch(rows)
		if err != nil {
			return err
		}
		byEpisode[m.EpisodeID] = m
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for i := range episodes {
		if m, ok := byEpisode[episodes[i].ID]; ok {
			episodes[i].Match = &m
		}
	}
	return nil
}

type MatchInput struct {
	GameID         *string
	TeamAID        string
	TeamBID        string
	Stage          string
	BestOf         int
	ScoreA         *int
	ScoreB         *int
	PlayedOn       *time.Time
	ExtraVideoURLs []string
}

func (s *Store) UpsertMatch(ctx context.Context, episodeID string, in MatchInput) error {
	if in.ExtraVideoURLs == nil {
		in.ExtraVideoURLs = []string{}
	}
	_, err := s.Pool.Exec(ctx, `
		insert into matches (episode_id, game_id, team_a_id, team_b_id, stage, best_of, score_a, score_b, played_on, extra_video_urls)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		on conflict (episode_id) do update set
			game_id = excluded.game_id, team_a_id = excluded.team_a_id, team_b_id = excluded.team_b_id,
			stage = excluded.stage, best_of = excluded.best_of, score_a = excluded.score_a,
			score_b = excluded.score_b, played_on = excluded.played_on, extra_video_urls = excluded.extra_video_urls
	`, episodeID, in.GameID, in.TeamAID, in.TeamBID, in.Stage, in.BestOf, in.ScoreA, in.ScoreB, in.PlayedOn, in.ExtraVideoURLs)
	return err
}

func (s *Store) DeleteMatch(ctx context.Context, episodeID string) error {
	_, err := s.Pool.Exec(ctx, `delete from matches where episode_id = $1`, episodeID)
	return err
}

func scanTeamMatches(rows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
}) ([]models.TeamMatch, error) {
	var list []models.TeamMatch
	for rows.Next() {
		var tm models.TeamMatch
		m, err := scanMatch(rows, &tm.EpisodeTitle, &tm.MediaTitle, &tm.MediaSlug, &tm.ThumbnailURL)
		if err != nil {
			return nil, err
		}
		tm.Match = m
		list = append(list, tm)
	}
	return list, rows.Err()
}

// The extra columns scanTeamMatches expects after the match columns.
const teamMatchExtraColumns = `, e.title, md.title, md.slug, md.thumbnail_url`

func (s *Store) ListMatchesForTeam(ctx context.Context, teamID string) ([]models.TeamMatch, error) {
	return s.listTeamMatchesWithExtras(ctx, `mt.team_a_id = $1 or mt.team_b_id = $1`, teamID)
}

func (s *Store) ListMatchesForFollowedTeams(ctx context.Context, profileID string) ([]models.TeamMatch, error) {
	return s.listTeamMatchesWithExtras(ctx, `exists (
		select 1 from profile_followed_teams f
		where f.profile_id = $1 and f.team_id in (mt.team_a_id, mt.team_b_id))`, profileID)
}

// listTeamMatchesWithExtras returns published matches (published episode
// of a published media title) matching the given where-clause on mt, newest
// first, together with where they live in the catalog.
func (s *Store) listTeamMatchesWithExtras(ctx context.Context, where string, args ...any) ([]models.TeamMatch, error) {
	query := `select mt.episode_id, mt.game_id, mt.team_a_id, mt.team_b_id, mt.stage, mt.best_of,
		mt.score_a, mt.score_b, mt.played_on, mt.extra_video_urls,
		g.id, g.name, g.slug,
		ta.id, ta.name, ta.slug, ta.short_name, ta.country, ta.logo_url, ta.created_at,
		tb.id, tb.name, tb.slug, tb.short_name, tb.country, tb.logo_url, tb.created_at` + teamMatchExtraColumns + `
	from matches mt
	left join games g on g.id = mt.game_id
	join teams ta on ta.id = mt.team_a_id
	join teams tb on tb.id = mt.team_b_id
	join episodes e on e.id = mt.episode_id
	join seasons sn on sn.id = e.season_id
	join series se on se.id = sn.series_id
	join media md on md.id = se.media_id
	where e.is_published = true and md.is_published = true and (` + where + `)
	order by coalesce(mt.played_on, make_date(md.release_year, 1, 1)) desc, e.episode_number asc
	limit 500`
	rows, err := s.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTeamMatches(rows)
}

// ListEsportsMedia returns published media titles that contain at least one
// match, for the "eSports" browse row.
func (s *Store) ListEsportsMedia(ctx context.Context) ([]models.Media, error) {
	rows, err := s.Pool.Query(ctx, `
		select `+mediaColumnsQualified+` from media m
		where m.is_published = true and exists (
			select 1 from series se
			join seasons sn on sn.series_id = se.id
			join episodes e on e.season_id = sn.id
			join matches mt on mt.episode_id = e.id
			where se.media_id = m.id)
		order by m.release_year desc, m.created_at desc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Media
	for rows.Next() {
		m, err := scanMedia(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

// --- Followed teams ---

func (s *Store) FollowTeam(ctx context.Context, profileID, teamID string) error {
	_, err := s.Pool.Exec(ctx, `
		insert into profile_followed_teams (profile_id, team_id) values ($1, $2)
		on conflict do nothing
	`, profileID, teamID)
	return err
}

func (s *Store) UnfollowTeam(ctx context.Context, profileID, teamID string) error {
	_, err := s.Pool.Exec(ctx, `delete from profile_followed_teams where profile_id = $1 and team_id = $2`, profileID, teamID)
	return err
}

// FollowedTeamIDs returns the set of team IDs the profile follows.
func (s *Store) FollowedTeamIDs(ctx context.Context, profileID string) (map[string]bool, error) {
	rows, err := s.Pool.Query(ctx, `select team_id from profile_followed_teams where profile_id = $1`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make(map[string]bool)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids[id] = true
	}
	return ids, rows.Err()
}

// EpisodeArtwork is what the generated episode graphic shows.
type EpisodeArtwork struct {
	Episode     models.Episode
	SeasonTitle string
	MediaTitle  string
	MediaSlug   string
	ReleaseYear int
}

// GetEpisodeArtwork loads an episode (with its match, if any) together with
// the season and title it belongs to, in one round trip for the episode
// row plus the match lookup.
func (s *Store) GetEpisodeArtwork(ctx context.Context, episodeID string) (EpisodeArtwork, error) {
	var a EpisodeArtwork
	err := s.Pool.QueryRow(ctx, `
		select sn.title, m.title, m.slug, m.release_year
		from episodes e
		join seasons sn on sn.id = e.season_id
		join series se on se.id = sn.series_id
		join media m on m.id = se.media_id
		where e.id = $1
	`, episodeID).Scan(&a.SeasonTitle, &a.MediaTitle, &a.MediaSlug, &a.ReleaseYear)
	if err != nil {
		return a, mapNotFound(err)
	}
	a.Episode, err = s.GetEpisodeByID(ctx, episodeID)
	return a, err
}
