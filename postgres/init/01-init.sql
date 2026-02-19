create type task_status as enum ('not processed', 'processing', 'completed', 'failed') ;

alter type task_status owner to postgres;

create table if not exists files_tasks
(
    id        serial
        primary key,
    filename  text                                                          not null
        unique,
    status    task_status              default 'not processed'::task_status not null,
    create_at timestamp with time zone default now()
);

alter table files_tasks
    owner to postgres;

create unique index if not exists idx_files_tasks_filename
    on files_tasks (filename);

create table if not exists reports
(
    timestamp timestamp with time zone default now() not null,
    filename  text                                   not null
        references files_tasks (filename)
            on delete cascade,
    msg       varchar(255)                           not null,
    primary key (timestamp, filename)
);

alter table reports
    owner to postgres;

create table if not exists  files_data
(
    id         bigserial
        primary key,
    filename   text
        references files_tasks (filename)
            on delete cascade,
    number     varchar(255),
    mqtt       varchar(255),
    inv_id     varchar(255),
    unit_guid  varchar(255),
    msg_id     varchar(255),
    msg_text   text,
    context    varchar(255),
    class      varchar(155),
    level      varchar(255),
    addr       varchar(255),
    area       varchar(255),
    block      varchar(255),
    type       varchar(255),
    bit        varchar(255),
    invert_bit varchar(255)
);

create index if not exists idx_files_data_filename
    on files_data (filename);

alter table files_data
    owner to postgres;

