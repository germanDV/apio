create table if not exists notes (
  id uuid not null primary key,
  title text not null,
  content text not null,
  created_by uuid not null,
  created_at timestamp not null default (now() at time zone 'utc'),
  updated_at timestamp not null default (now() at time zone 'utc')
);

create index if not exists notes_created_by_idx on notes (created_by);

create table if not exists tags (
  id uuid not null primary key,
  name text unique not null
);

create table notes_tags (
  note_id uuid references notes(id) on delete cascade,
  tag_id uuid references tags(id) on delete cascade,
  primary key (note_id, tag_id)
);

---- create above / drop below ----

drop table if exists notes;
drop index if exists notes_created_by_idx;
drop table if exists tags;
drop table if exists notes_tags;
