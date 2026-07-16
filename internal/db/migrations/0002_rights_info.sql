create table rights_info (
  id uuid primary key default gen_random_uuid(),
  media_id uuid not null unique references media(id) on delete cascade,
  source text not null check (source in ('PUBLIC_DOMAIN', 'CREATIVE_COMMONS', 'REVENUE_SHARE', 'OWNED', 'LICENSED')),
  rights_holder text,
  license_doc_url text,
  revenue_share_pct double precision,
  territory_limit text,
  valid_from timestamptz not null,
  valid_until timestamptz,
  notes text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create index idx_rights_info_valid_until on rights_info(valid_until);
