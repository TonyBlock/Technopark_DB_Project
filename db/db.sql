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
    parent    int                   DEFAULT 0 REFERENCES posts,
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
    voice     int    NOT NULL
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
SET votes = votes + NEW.voice
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
SET votes = votes + NEW.voice - OLD.voice
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
BEGIN
    new.path = (SELECT path FROM posts WHERE id = new.parent) || new.id;
RETURN new;
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
SET posts = forum.posts + 1
WHERE slug = NEW.forum;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_post_after
    AFTER INSERT
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE insert_post_after_proc();


CREATE OR REPLACE FUNCTION update_post_proc()
    RETURNS TRIGGER AS
$$
BEGIN
IF OLD.is_edited = false and NEW.message <> OLD.message then
    UPDATE posts
    SET is_edited = true
    WHERE id = NEW.id;
RETURN NEW;
END;
$$ language plpgsql;

CREATE TRIGGER update_posts
    AFTER UPDATE
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE update_post_proc();


CREATE OR REPLACE FUNCTION insert_threads_proc()
    RETURNS TRIGGER AS
$$
BEGIN
UPDATE forums
SET threads = forum.threads + 1
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
    ON thread
    FOR EACH ROW
    EXECUTE PROCEDURE add_user();

CREATE TRIGGER insert_new_post
    AFTER INSERT
    ON post
    FOR EACH ROW
    EXECUTE PROCEDURE add_user();

-- INDEXES
create index if not exists users_nickname on users (nickname);
create index if not exists users_all on users (nickname, fullname, about, email);
create index if not exists forums_slug_hash on forums using hash (slug);
create index if not exists forums_user_hash on forums using hash (user_);
create index if not exists user_forum_forum on user_forum (forum);
create index if not exists user_forum_nickname on user_forum (nickname);
create index if not exists user_forum_all on user_forum (forum, nickname);
create index if not exists threads_slug_id on threads using hash (id);
create index if not exists threads_slug_id on threads using hash (id, slug);
create index if not exists threads_slug on threads using hash (slug);
create index if not exists threads_user on threads using hash (author);
create index if not exists threads_created on threads (created);
create index if not exists threads_forum on threads using hash (forum);
create index if not exists threads_forum_created on threads (forum, created);
create index if not exists posts_thread_created on posts (thread, created, id);
create index if not exists posts_sorting_desc on posts ((path[1]) desc, path, id);
create index if not exists posts_sorting_asc on posts ((path[1]) asc, path, id);
create index if not exists posts_thread on posts using hash (thread);
create index if not exists posts_parent on posts (thread, id, (path[1]), parent);
create index if not exists posts_thread_created_id ON posts (id, thread, created);
create index if not exists posts_path_first_path ON posts ((path[1]), path);
create index if not exists posts_thread_path_id on posts (thread, path, id);
create unique index if not exists votes_all on votes (nickname, thread, voice);
create unique index if not exists votes on votes (nickname, thread);

VACUUM;
VACUUM ANALYSE;
