-- ============================================================
-- Muziboo Engagement & Publishing — Supabase Migration
-- Run this in: Supabase Dashboard → SQL Editor → New Query
--
-- Adds: votes, comments, messages tables (votes/comments were
-- referenced by the API but never created), a draft/publish
-- flag on tracks, and a collision-safe signup trigger.
-- ============================================================

-- 1. VOTES TABLE
create table if not exists public.votes (
  id uuid default gen_random_uuid() primary key,
  track_id uuid references public.tracks(id) on delete cascade not null,
  user_id uuid references public.profiles(id) on delete cascade not null,
  created_at timestamptz default now() not null,
  unique(track_id, user_id)
);

create index if not exists votes_track_id_idx on public.votes (track_id);
create index if not exists votes_user_id_idx on public.votes (user_id);

alter table public.votes enable row level security;

drop policy if exists "Votes are viewable by everyone" on public.votes;
create policy "Votes are viewable by everyone"
  on public.votes for select
  using (true);

drop policy if exists "Users can insert their own votes" on public.votes;
create policy "Users can insert their own votes"
  on public.votes for insert
  with check (auth.uid() = user_id);

drop policy if exists "Users can delete their own votes" on public.votes;
create policy "Users can delete their own votes"
  on public.votes for delete
  using (auth.uid() = user_id);

-- 2. COMMENTS TABLE
create table if not exists public.comments (
  id uuid default gen_random_uuid() primary key,
  track_id uuid references public.tracks(id) on delete cascade not null,
  user_id uuid references public.profiles(id) on delete cascade not null,
  content text not null,
  created_at timestamptz default now() not null,
  updated_at timestamptz default now() not null
);

create index if not exists comments_track_id_idx on public.comments (track_id);

alter table public.comments enable row level security;

drop policy if exists "Comments are viewable by everyone" on public.comments;
create policy "Comments are viewable by everyone"
  on public.comments for select
  using (true);

drop policy if exists "Users can insert their own comments" on public.comments;
create policy "Users can insert their own comments"
  on public.comments for insert
  with check (auth.uid() = user_id);

drop policy if exists "Users can delete their own comments" on public.comments;
create policy "Users can delete their own comments"
  on public.comments for delete
  using (auth.uid() = user_id);

drop trigger if exists comments_updated_at on public.comments;
create trigger comments_updated_at
  before update on public.comments
  for each row execute function public.update_updated_at();

-- 3. MESSAGES TABLE
-- Already exists in production but was never recorded in a
-- migration file. Kept here so the schema is reproducible.
create table if not exists public.messages (
  id uuid default gen_random_uuid() primary key,
  sender_id uuid references public.profiles(id) on delete cascade not null,
  recipient_id uuid references public.profiles(id) on delete cascade not null,
  content text not null,
  created_at timestamptz default now() not null,
  read_at timestamptz
);

create index if not exists messages_sender_id_idx on public.messages (sender_id);
create index if not exists messages_recipient_id_idx on public.messages (recipient_id);

alter table public.messages enable row level security;

drop policy if exists "Participants can view their messages" on public.messages;
create policy "Participants can view their messages"
  on public.messages for select
  using (auth.uid() = sender_id or auth.uid() = recipient_id);

drop policy if exists "Users can send messages as themselves" on public.messages;
create policy "Users can send messages as themselves"
  on public.messages for insert
  with check (auth.uid() = sender_id);

-- 4. DRAFT / PUBLISH FLAG ON TRACKS
-- New uploads are private drafts until the owner sets them live.
-- Existing tracks were uploaded when everything was public, so
-- they are backfilled to live to avoid surprising owners.
alter table public.tracks
add column if not exists is_public boolean default false not null;

update public.tracks set is_public = true;

-- Replace the world-readable select policy: drafts are visible
-- only to their owner.
drop policy if exists "Public tracks are viewable by everyone" on public.tracks;
create policy "Live tracks are viewable by everyone"
  on public.tracks for select
  using (is_public = true or auth.uid() = user_id);

create index if not exists tracks_is_public_idx on public.tracks (is_public);

-- 5. COLLISION-SAFE SIGNUP TRIGGER
-- The previous version failed the whole signup when two users
-- shared an email prefix (unique violation on username).
create or replace function public.handle_new_user()
returns trigger as $$
declare
  base_username text;
  final_username text;
  attempts int := 0;
begin
  base_username := coalesce(
    new.raw_user_meta_data->>'username',
    split_part(new.email, '@', 1)
  );
  if base_username is null or base_username = '' then
    base_username := 'user';
  end if;

  final_username := base_username;
  while exists (select 1 from public.profiles where username = final_username) loop
    attempts := attempts + 1;
    final_username := base_username || '_' || floor(random() * 100000)::int;
    if attempts > 10 then
      final_username := base_username || '_' || replace(new.id::text, '-', '');
      exit;
    end if;
  end loop;

  insert into public.profiles (id, username, display_name)
  values (
    new.id,
    final_username,
    coalesce(new.raw_user_meta_data->>'display_name', base_username)
  );
  return new;
end;
$$ language plpgsql security definer;
