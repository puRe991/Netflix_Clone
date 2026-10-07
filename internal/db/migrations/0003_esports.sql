-- eSports: games, teams, match metadata on episodes, spoiler mode per profile.
--
-- Mapping onto the existing catalog: a tournament is a SERIES media title,
-- a stage (e.g. "Playoffs") is a season and every match is an episode. The
-- match row below is an optional 1:1 extension of an episode, so regular
-- series stay untouched.

create table games (
  id uuid primary key default gen_random_uuid(),
  name text not null unique,
  slug text not null unique
);

create table teams (
  id uuid primary key default gen_random_uuid(),
  name text not null,
  slug text not null unique,
  short_name text,
  country text,
  logo_url text,
  created_at timestamptz not null default now()
);

create table matches (
  episode_id uuid primary key references episodes(id) on delete cascade,
  game_id uuid references games(id) on delete set null,
  team_a_id uuid not null references teams(id) on delete restrict,
  team_b_id uuid not null references teams(id) on delete restrict,
  stage text not null,
  best_of int not null check (best_of in (1, 3, 5)),
  -- Scores are optional (unknown/not yet played) and only ever shown when
  -- the viewer has spoilers enabled.
  score_a int check (score_a >= 0),
  score_b int check (score_b >= 0),
  played_on date,
  -- Additional VOD parts (e.g. map 2, map 3) played after episodes.video_url.
  extra_video_urls text[] not null default '{}',
  check (team_a_id <> team_b_id)
);

create index idx_matches_team_a on matches(team_a_id);
create index idx_matches_team_b on matches(team_b_id);

-- New profiles default to spoiler-free viewing: for eSports, knowing the
-- result before watching is the single biggest way to ruin a VOD.
alter table profiles add column hide_spoilers boolean not null default true;

create table profile_followed_teams (
  profile_id uuid not null references profiles(id) on delete cascade,
  team_id uuid not null references teams(id) on delete cascade,
  created_at timestamptz not null default now(),
  primary key (profile_id, team_id)
);

-- Which VOD part (map) the saved progress_seconds refers to.
alter table watch_progress add column part_index int not null default 0 check (part_index >= 0);

-- Officially published, embeddable videos (e.g. a tournament organizer's own
-- YouTube upload) played through the platform's embed player. These must
-- never sit behind the paywall (see docs/esports-vods.md).
alter table rights_info drop constraint rights_info_source_check;
alter table rights_info add constraint rights_info_source_check
  check (source in ('PUBLIC_DOMAIN', 'CREATIVE_COMMONS', 'REVENUE_SHARE', 'OWNED', 'LICENSED', 'OFFICIAL_EMBED'));
