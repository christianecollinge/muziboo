-- ============================================================
-- Muziboo Platform — Supabase Migration
-- Run this in: Supabase Dashboard → SQL Editor → New Query
-- ============================================================

-- 1. PROFILES TABLE
-- Linked to Supabase auth.users via id
create table public.profiles (
  id uuid references auth.users on delete cascade primary key,
  username text unique not null,
  display_name text,
  bio text,
  avatar_url text,
  created_at timestamptz default now() not null,
  updated_at timestamptz default now() not null
);

-- Index for fast username lookups
create index profiles_username_idx on public.profiles (username);

-- RLS: anyone can read, only owner can write
alter table public.profiles enable row level security;

create policy "Public profiles are viewable by everyone"
  on public.profiles for select
  using (true);

create policy "Users can insert their own profile"
  on public.profiles for insert
  with check (auth.uid() = id);

create policy "Users can update their own profile"
  on public.profiles for update
  using (auth.uid() = id);

-- 2. TRACKS TABLE
create table public.tracks (
  id uuid default gen_random_uuid() primary key,
  user_id uuid references public.profiles(id) on delete cascade not null,
  title text not null,
  description text,
  genre text,
  audio_url text not null,
  artwork_url text,
  duration integer, -- duration in seconds
  created_at timestamptz default now() not null,
  updated_at timestamptz default now() not null
);

-- Indexes
create index tracks_user_id_idx on public.tracks (user_id);
create index tracks_created_at_idx on public.tracks (created_at desc);

-- RLS: anyone can read, only owner can write/delete
alter table public.tracks enable row level security;

create policy "Public tracks are viewable by everyone"
  on public.tracks for select
  using (true);

create policy "Users can insert their own tracks"
  on public.tracks for insert
  with check (auth.uid() = user_id);

create policy "Users can update their own tracks"
  on public.tracks for update
  using (auth.uid() = user_id);

create policy "Users can delete their own tracks"
  on public.tracks for delete
  using (auth.uid() = user_id);

-- 3. AUTO-CREATE PROFILE ON SIGNUP
-- When a new user signs up, auto-create a profile row
create or replace function public.handle_new_user()
returns trigger as $$
begin
  insert into public.profiles (id, username, display_name)
  values (
    new.id,
    coalesce(new.raw_user_meta_data->>'username', split_part(new.email, '@', 1)),
    coalesce(new.raw_user_meta_data->>'display_name', split_part(new.email, '@', 1))
  );
  return new;
end;
$$ language plpgsql security definer;

create trigger on_auth_user_created
  after insert on auth.users
  for each row execute function public.handle_new_user();

-- 4. UPDATED_AT TRIGGER
create or replace function public.update_updated_at()
returns trigger as $$
begin
  new.updated_at = now();
  return new;
end;
$$ language plpgsql;

create trigger profiles_updated_at
  before update on public.profiles
  for each row execute function public.update_updated_at();

create trigger tracks_updated_at
  before update on public.tracks
  for each row execute function public.update_updated_at();

-- 5. STORAGE BUCKETS
-- Create public buckets for audio, artwork, and avatars
insert into storage.buckets (id, name, public) values ('audio', 'audio', true);
insert into storage.buckets (id, name, public) values ('artwork', 'artwork', true);
insert into storage.buckets (id, name, public) values ('avatars', 'avatars', true);

-- Storage policies: authenticated users can upload, anyone can read
create policy "Anyone can read audio files"
  on storage.objects for select
  using (bucket_id = 'audio');

create policy "Authenticated users can upload audio"
  on storage.objects for insert
  with check (bucket_id = 'audio' and auth.role() = 'authenticated');

create policy "Users can delete their own audio"
  on storage.objects for delete
  using (bucket_id = 'audio' and auth.uid()::text = (storage.foldername(name))[1]);

create policy "Anyone can read artwork"
  on storage.objects for select
  using (bucket_id = 'artwork');

create policy "Authenticated users can upload artwork"
  on storage.objects for insert
  with check (bucket_id = 'artwork' and auth.role() = 'authenticated');

create policy "Users can delete their own artwork"
  on storage.objects for delete
  using (bucket_id = 'artwork' and auth.uid()::text = (storage.foldername(name))[1]);

create policy "Anyone can read avatars"
  on storage.objects for select
  using (bucket_id = 'avatars');

create policy "Authenticated users can upload avatars"
  on storage.objects for insert
  with check (bucket_id = 'avatars' and auth.role() = 'authenticated');

create policy "Users can delete their own avatars"
  on storage.objects for delete
  using (bucket_id = 'avatars' and auth.uid()::text = (storage.foldername(name))[1]);
