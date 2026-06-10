-- ============================================================
-- Muziboo Private Audio — Supabase Migration
-- Run this in: Supabase Dashboard → SQL Editor → New Query
--
-- Makes the audio and stems buckets private. The Go backend now
-- serves time-limited signed URLs instead of permanent public
-- links, so draft tracks and invite-only stems can no longer be
-- streamed by anyone who has the raw storage URL.
--
-- IMPORTANT: deploy the signed-URL backend BEFORE running this,
-- otherwise playback breaks until the deploy finishes.
-- Artwork and avatars intentionally stay public (low-sensitivity
-- images; keeps CDN caching).
-- ============================================================

-- 1. MAKE BUCKETS PRIVATE
update storage.buckets set public = false where id in ('audio', 'stems');

-- 2. DROP THE WORLD-READ POLICIES
-- With the bucket private these only governed authenticated REST
-- downloads, but any logged-in user could still fetch any file.
drop policy if exists "Anyone can read audio files" on storage.objects;
drop policy if exists "Anyone can read stems" on storage.objects;

-- 3. OWNERS KEEP DIRECT READ ACCESS WITH THEIR OWN JWT
-- (The app itself uses service-role signed URLs and doesn't need
-- these, but they keep the rules honest for direct API use.)
drop policy if exists "Owners can read their own audio" on storage.objects;
create policy "Owners can read their own audio"
  on storage.objects for select
  using (bucket_id = 'audio' and auth.uid()::text = (storage.foldername(name))[1]);

drop policy if exists "Owners can read their own stems" on storage.objects;
create policy "Owners can read their own stems"
  on storage.objects for select
  using (bucket_id = 'stems' and auth.uid()::text = (storage.foldername(name))[1]);
