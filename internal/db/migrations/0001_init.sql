create extension if not exists pgcrypto;

create table users (
  id uuid primary key default gen_random_uuid(),
  email text not null unique,
  password_hash text not null,
  name text,
  role text not null default 'USER' check (role in ('USER', 'ADMIN')),
  stripe_customer_id text unique,
  subscription_status text not null default 'FREE'
    check (subscription_status in ('FREE', 'TRIALING', 'ACTIVE', 'PAST_DUE', 'CANCELED', 'INCOMPLETE')),
  last_profile_id uuid,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table profiles (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references users(id) on delete cascade,
  name text not null,
  avatar_url text,
  is_kids_profile boolean not null default false,
  created_at timestamptz not null default now()
);

create table media (
  id uuid primary key default gen_random_uuid(),
  title text not null,
  slug text not null unique,
  description text not null,
  type text not null check (type in ('MOVIE', 'SERIES')),
  thumbnail_url text not null,
  banner_url text not null,
  trailer_url text,
  video_url text,
  release_year int not null,
  duration int,
  age_rating text not null,
  language text not null default 'Deutsch',
  tags text[] not null default '{}',
  is_published boolean not null default false,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table series (
  id uuid primary key default gen_random_uuid(),
  media_id uuid not null unique references media(id) on delete cascade,
  title text not null,
  description text not null
);

create table seasons (
  id uuid primary key default gen_random_uuid(),
  series_id uuid not null references series(id) on delete cascade,
  season_number int not null,
  title text not null,
  unique (series_id, season_number)
);

create table episodes (
  id uuid primary key default gen_random_uuid(),
  season_id uuid not null references seasons(id) on delete cascade,
  title text not null,
  description text not null,
  episode_number int not null,
  duration int not null,
  video_url text not null,
  thumbnail_url text not null,
  is_published boolean not null default false,
  unique (season_id, episode_number)
);

create table genres (
  id uuid primary key default gen_random_uuid(),
  name text not null unique,
  slug text not null unique
);

create table media_genres (
  media_id uuid not null references media(id) on delete cascade,
  genre_id uuid not null references genres(id) on delete cascade,
  primary key (media_id, genre_id)
);

create table watch_progress (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references users(id) on delete cascade,
  profile_id uuid not null references profiles(id) on delete cascade,
  media_id uuid not null references media(id) on delete cascade,
  episode_id uuid references episodes(id) on delete cascade,
  progress_seconds int not null default 0,
  completed boolean not null default false,
  updated_at timestamptz not null default now()
);

-- Partial unique indexes instead of a single UNIQUE(profile_id, media_id, episode_id):
-- Postgres treats every NULL as distinct in a plain unique constraint, so a naive
-- UNIQUE(profile_id, media_id, episode_id) would silently allow duplicate rows for
-- movies (episode_id IS NULL). These two indexes make both cases race-free upsert targets.
create unique index watch_progress_movie_uniq on watch_progress(profile_id, media_id) where episode_id is null;
create unique index watch_progress_episode_uniq on watch_progress(profile_id, media_id, episode_id) where episode_id is not null;

create table watchlist (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references users(id) on delete cascade,
  profile_id uuid not null references profiles(id) on delete cascade,
  media_id uuid not null references media(id) on delete cascade,
  created_at timestamptz not null default now(),
  unique (profile_id, media_id)
);

create table subscriptions (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references users(id) on delete cascade,
  stripe_subscription_id text not null unique,
  plan text not null check (plan in ('FREE', 'BASIC', 'PREMIUM')),
  status text not null check (status in ('FREE', 'TRIALING', 'ACTIVE', 'PAST_DUE', 'CANCELED', 'INCOMPLETE')),
  current_period_end timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table payment_logs (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references users(id) on delete set null,
  stripe_event_id text not null unique,
  type text not null,
  raw_data jsonb not null,
  created_at timestamptz not null default now()
);

create index idx_media_published on media(is_published);
create index idx_watchlist_user on watchlist(user_id);
create index idx_watch_progress_user_media on watch_progress(user_id, media_id);
create index idx_profiles_user on profiles(user_id);
