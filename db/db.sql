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
create index if not exists users_nickname_nickname_email on users (nickname, email);

--create index if not exists user_forum_forum on user_forum (forum);
create index if not exists user_forum_nickname on user_forum (nickname);
create index if not exists user_forum_all on user_forum (forum, nickname);

--create index if not exists threads_slug on threads (forum);
create index if not exists threads_created on threads (created);
create index if not exists threads_forum_created on threads (forum, created);

--create index if not exists posts_id_thread on posts (thread, id);
create index if not exists posts_id_thread on posts (thread, id, parent NULLS FIRST);
create index if not exists posts_id_path_path1 on posts (path, (path[1]), id);
--create index if not exists posts_path_path1 on posts (path, (path[1]));
create index if not exists posts_id_thread_parent_path1 on posts ((path[1]), thread, id, parent NULLS FIRST);
--create index if not exists posts_thread on posts (thread);
create index if not exists posts_thread_past on posts (thread, path);

create unique index if not exists votes_key on votes (thread, nickname);
