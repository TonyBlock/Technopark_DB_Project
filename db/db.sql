CREATE EXTENSION IF NOT EXISTS citext;

-- TABLES
CREATE UNLOGGED TABLE IF NOT EXISTS users
(
    nickname citext COLLATE "ucs_basic" NOT NULL UNIQUE PRIMARY KEY,
    fullname text                       NOT NULL,
    about    text,
    email    citext                     NOT NULL UNIQUE
);

CREATE UNLOGGED TABLE IF NOT EXISTS forums
(
    title          text   NOT NULL,
    user_          citext NOT NULL REFERENCES users (nickname),
    slug           citext NOT NULL PRIMARY KEY,
    posts          bigint DEFAULT 0,
    threads        int    DEFAULT 0
);

CREATE UNLOGGED TABLE IF NOT EXISTS threads
(
    id      bigserial             NOT NULL PRIMARY KEY,
    title   text                  NOT NULL,
    author  citext                NOT NULL REFERENCES users (nickname),
    forum   citext                NOT NULL REFERENCES forums (slug),
    message text                  NOT NULL,
    votes   int                   DEFAULT 0,
    slug    citext,
    created timestamp with time zone DEFAULT now()
);

CREATE UNLOGGED TABLE IF NOT EXISTS posts
(
    id        bigserial             NOT NULL UNIQUE PRIMARY KEY ,
    parent    int                   REFERENCES posts (id),
    author    citext                NOT NULL REFERENCES users (nickname),
    message   text                  NOT NULL,
    is_edited bool                  DEFAULT FALSE,
    forum     citext                NOT NULL REFERENCES forums (slug),
    thread    int                   NOT NULL REFERENCES threads (id),
    created   timestamp with time zone DEFAULT now(),
    path      bigint[]              DEFAULT ARRAY []::INTEGER[]
);

CREATE UNLOGGED TABLE IF NOT EXISTS votes
(
    nickname  citext NOT NULL REFERENCES users (nickname),
    thread    int    NOT NULL REFERENCES threads (id),
    voice     int    NOT NULL,
    constraint user_thread_key unique (nickname, thread)
);

CREATE UNLOGGED TABLE IF NOT EXISTS user_forum
(
    nickname citext COLLATE "ucs_basic" NOT NULL REFERENCES users (nickname),
    forum    citext NOT NULL REFERENCES forums (slug),
    constraint user_forum_key unique (nickname, forum)
);

-- TRIGGERS AND PROCEDURES
CREATE OR REPLACE FUNCTION insert_votes_proc()
    RETURNS TRIGGER AS
$$
BEGIN
UPDATE threads
SET votes = threads.votes + NEW.voice
WHERE id = NEW.thread;
RETURN NEW;
END;
$$ language plpgsql;

CREATE TRIGGER insert_votes
    AFTER INSERT
    ON votes
    FOR EACH ROW
    EXECUTE PROCEDURE insert_votes_proc();


CREATE OR REPLACE FUNCTION update_votes_proc()
    RETURNS TRIGGER AS
$$
BEGIN
UPDATE threads
SET votes = threads.votes + NEW.voice - OLD.voice
WHERE id = NEW.thread;
RETURN NEW;
END;
$$ language plpgsql;

CREATE TRIGGER update_votes
    AFTER UPDATE
    ON votes
    FOR EACH ROW
    EXECUTE PROCEDURE update_votes_proc();


CREATE OR REPLACE FUNCTION insert_post_before_proc()
    RETURNS TRIGGER AS
$$
DECLARE
parent_post_id posts.id%type := 0;
BEGIN
    NEW.path = (SELECT path FROM posts WHERE id = new.parent) || NEW.id;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_post_before
    BEFORE INSERT
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE insert_post_before_proc();


CREATE OR REPLACE FUNCTION insert_post_after_proc()
    RETURNS TRIGGER AS
$$
BEGIN
UPDATE forums
SET posts = forums.posts + 1
WHERE slug = NEW.forum;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_post_after
    AFTER INSERT
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE insert_post_after_proc();


CREATE OR REPLACE FUNCTION insert_threads_proc()
    RETURNS TRIGGER AS
$$
BEGIN
UPDATE forums
SET threads = forums.threads + 1
WHERE slug = NEW.forum;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_threads
    AFTER INSERT
    ON threads
    FOR EACH ROW
    EXECUTE PROCEDURE insert_threads_proc();


CREATE OR REPLACE FUNCTION add_user()
    RETURNS TRIGGER AS
$$
BEGIN
INSERT INTO user_forum (nickname, forum)
VALUES (NEW.author, NEW.forum)
    ON CONFLICT do nothing;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_new_thread
    AFTER INSERT
    ON threads
    FOR EACH ROW
    EXECUTE PROCEDURE add_user();

CREATE TRIGGER insert_new_post
    AFTER INSERT
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE add_user();

-- INDEXES
--create index if not exists users_nickname_lower on users (LOWER(nickname));
create index if not exists users_nickname_nickname_email on users (nickname, email);

--create index if not exists forums_slug_hash on forums (LOWER(slug));

create index if not exists user_forum_forum on user_forum (forum);
create index if not exists user_forum_nickname on user_forum (nickname);----
create index if not exists user_forum_all on user_forum (forum, nickname);

--create index if not exists threads_slug on threads (LOWER(slug));
create index if not exists threads_slug on threads (forum);
create index if not exists threads_created on threads (created);----
create index if not exists threads_forum_created on threads (forum, created);

-- create index if not exists posts_thread_created on posts (id, thread, created);----
-- create index if not exists posts_thread_parent on posts (id, thread, thread, path[1]); --?
-- create index if not exists posts_thread_parent on posts (id, parent, thread); --?
-- create index if not exists posts_sorting_desc on posts (id, path, (path[1])); --?
-- create index if not exists posts_sorting_asc on posts ((path[1]) asc, path, id);----
-- create index if not exists posts_thread on posts (thread);----
-- create index if not exists posts_parent on posts (id, path, thread, (path[1]), parent); --?
-- create index if not exists posts_thread_created_id ON posts (id, thread, created); --
-- create index if not exists posts_path_id ON posts (id, path);----
-- create index if not exists posts_path_first_path ON posts ((path[1]), path);----
-- create index if not exists posts_thread_path_id on posts (thread, path, id) --
-- create index if not exists posts_thread_path on posts (thread, path); --

--CREATE INDEX IF NOT EXISTS posts_thread_id_path1_id_idx ON posts (thread, (path[1]), id); --?
CREATE INDEX IF NOT EXISTS posts_thread_id_path1_id_idx ON posts ((path[1]), thread, (path[2:])); --
CREATE INDEX IF NOT EXISTS posts_thread_id_path2_id_idx ON posts (path, (path[1]), thread); --

CREATE INDEX IF NOT EXISTS posts_thread_id_path_idx ON posts (thread, path); --

CREATE INDEX IF NOT EXISTS posts_thread_id_id_idx ON posts (thread, id); --

--CREATE INDEX IF NOT EXISTS posts_thread_id_parent_path_idx ON posts (thread, parent, path); --?
CREATE INDEX IF NOT EXISTS posts_thread_id_parent_path_idx ON posts (path[1], thread, (parent NULLS FIRST)); --

--CREATE INDEX IF NOT EXISTS posts_parent_id_idx ON posts (parent, id); --?

--CREATE INDEX IF NOT EXISTS posts_id_created_thread_id_idx ON posts (id, created, thread); --?

--CREATE INDEX IF NOT EXISTS posts_id_path_idx ON posts (id, path); --?

--create index if not exists posts_thread on posts (thread); --example

-- create index if not exists posts_id_thread on posts (id, thread);
-- create index if not exists posts_id_thread on posts (id, thread, parent NULLS FIRST);
-- create index if not exists posts_id_path_path1 on posts (id, path, (path[1]));
-- create index if not exists posts_path_path1 on posts (path, (path[1]));
-- create index if not exists posts_id_thread_parent_path1 on posts (id, thread, (path[1]), parent NULLS FIRST);
-- create index if not exists posts_thread on posts (thread) INCLUDE (path);


-- CREATE INDEX IF NOT EXISTS post_id_path ON posts(id, (path[1]));
-- CREATE INDEX IF NOT EXISTS post_thread_path_id ON posts(thread, path, id);
-- CREATE INDEX IF NOT EXISTS post_thread_id_path1_parent ON posts(thread, id, (path[1]), parent);
-- CREATE INDEX IF NOT EXISTS post_path1 ON posts((path[1]));
-- CREATE INDEX IF NOT EXISTS post_thr_id ON posts(thread);
-- CREATE INDEX IF NOT EXISTS post_thread_id ON posts(thread, id);

create unique index if not exists votes_key on votes (thread, nickname);

VACUUM;
VACUUM ANALYSE;
