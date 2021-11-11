-- +goose Up

create table locations
(
    id         bigserial primary key,
    latitude   double precision not null,
    longitude  double precision not null,
    title      text             not null,
    removed    boolean          not null default false,
    created_at timestamptz      not null default now(),
    updated_at timestamptz      not null default now()
);

create index locations_removed_index on locations (removed);

create type location_event_type as enum ('Created', 'Updated', 'Removed');
create type location_event_status as enum ('Deferred', 'Processed');

create table locations_events
(
    id          bigserial primary key,
    location_id bigint                not null references locations (id),
    type        location_event_type   not null,
    status      location_event_status not null default 'Deferred',
    payload     jsonb                 not null,
    updated_at  timestamptz           not null default now()
);

create index locations_events_status_index on locations_events (status);

-- +goose Down

drop index locations_events_status_index;
drop index locations_removed_index;

drop type location_event_status;
drop type location_event_type;

drop table locations_events;
drop table locations;
