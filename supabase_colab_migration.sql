-- ============================================================
-- Muziboo Colab Feature — Supabase Migration
-- Run this in: Supabase Dashboard → SQL Editor → New Query
-- ============================================================

-- 1. UPDATE TRACKS TABLE
-- Add support for branching and identifying Colab projects
alter table public.tracks 
add column if not exists parent_track_id uuid references public.tracks(id) on delete set null,
add column if not exists is_colab boolean default false not null;

create index if not exists tracks_parent_track_id_idx on public.tracks (parent_track_id);

-- 2. TRACK INVITES TABLE
-- Tracks which users are invited to download stems for a colab
create table if not exists public.track_invites (
  id uuid default gen_random_uuid() primary key,
  track_id uuid references public.tracks(id) on delete cascade not null,
  invited_user_id uuid references public.profiles(id) on delete cascade not null,
  created_at timestamptz default now() not null,
  unique(track_id, invited_user_id)
);

create index if not exists track_invites_track_id_idx on public.track_invites (track_id);
create index if not exists track_invites_user_id_idx on public.track_invites (invited_user_id);

alter table public.track_invites enable row level security;

create policy "Track owners and invitees can view invites"
  on public.track_invites for select
  using (
    invited_user_id = auth.uid() or 
    exists (select 1 from public.tracks where id = track_id and user_id = auth.uid())
  );

create policy "Track owners can insert invites"
  on public.track_invites for insert
  with check (exists (select 1 from public.tracks where id = track_id and user_id = auth.uid()));

create policy "Track owners can delete invites"
  on public.track_invites for delete
  using (exists (select 1 from public.tracks where id = track_id and user_id = auth.uid()));

-- 3. STEMS TABLE
-- Stores individual audio stems associated with a colab track
create table if not exists public.stems (
  id uuid default gen_random_uuid() primary key,
  track_id uuid references public.tracks(id) on delete cascade not null,
  name text not null,
  audio_url text not null,
  created_at timestamptz default now() not null
);

create index if not exists stems_track_id_idx on public.stems (track_id);

alter table public.stems enable row level security;

-- Only track owners or invited users can view stems
create policy "Owners and invitees can view stems"
  on public.stems for select
  using (
    exists (select 1 from public.tracks where id = track_id and user_id = auth.uid()) or
    exists (select 1 from public.track_invites where track_id = stems.track_id and invited_user_id = auth.uid())
  );

create policy "Track owners can insert stems"
  on public.stems for insert
  with check (exists (select 1 from public.tracks where id = track_id and user_id = auth.uid()));

create policy "Track owners can delete stems"
  on public.stems for delete
  using (exists (select 1 from public.tracks where id = track_id and user_id = auth.uid()));

-- 4. STORAGE BUCKET FOR STEMS
-- Note: We are making this bucket public so the Go backend can serve direct links.
-- However, the UI will only render these links for authorized users (owners/invitees).
-- If strict cryptographic access is needed later, we can change this to a private bucket.
insert into storage.buckets (id, name, public) 
values ('stems', 'stems', true)
on conflict (id) do nothing;

create policy "Anyone can read stems"
  on storage.objects for select
  using (bucket_id = 'stems');

create policy "Authenticated users can upload stems"
  on storage.objects for insert
  with check (bucket_id = 'stems' and auth.role() = 'authenticated');

create policy "Users can delete their own stems"
  on storage.objects for delete
  using (bucket_id = 'stems' and auth.uid()::text = (storage.foldername(name))[1]);
