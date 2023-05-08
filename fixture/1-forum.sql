-- forum database schema

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE OR REPLACE FUNCTION UUID() RETURNS uuid AS $$
    BEGIN
            RETURN uuid_generate_v4();
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION updated_datetime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE IF NOT EXISTS forum_users (
    id SERIAL PRIMARY KEY,
    nickname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    oauth_data JSONB NOT NULL DEFAULT '{}',
	created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS forum_sessions (
	uuid uuid NOT NULL DEFAULT uuid(),
    user_id INTEGER NOT NULL,
    data JSONB NOT NULL DEFAULT '{}',
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (uuid)
);

CREATE TABLE IF NOT EXISTS forum_topics (
    id SERIAL PRIMARY KEY,
    parent_id INTEGER NOT NULL DEFAULT 0,
    zorder INTEGER NOT NULL DEFAULT 0,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_by INTEGER NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES forum_users(id)
);

CREATE TABLE IF NOT EXISTS forum_threads (
    id SERIAL PRIMARY KEY,
    topic_id INTEGER NOT NULL,
    zorder INTEGER NOT NULL DEFAULT 0,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_by INTEGER NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (topic_id) REFERENCES forum_topics(id),
    FOREIGN KEY (created_by) REFERENCES forum_users(id)
);

CREATE TABLE IF NOT EXISTS forum_posts (
    id SERIAL PRIMARY KEY,
    thread_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_by INTEGER NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (thread_id) REFERENCES forum_threads(id),
    FOREIGN KEY (created_by) REFERENCES forum_users(id)
);

CREATE TRIGGER forum_users_updated_at
    BEFORE UPDATE
    ON forum_users
    FOR EACH ROW
    EXECUTE PROCEDURE updated_datetime();

CREATE TRIGGER forum_sessions_updated_at
    BEFORE UPDATE
    ON forum_sessions
    FOR EACH ROW
    EXECUTE PROCEDURE updated_datetime();

CREATE TRIGGER forum_topics_updated_at
    BEFORE UPDATE
    ON forum_topics
    FOR EACH ROW
    EXECUTE PROCEDURE updated_datetime();

CREATE TRIGGER forum_threads_updated_at
    BEFORE UPDATE
    ON forum_threads
    FOR EACH ROW
    EXECUTE PROCEDURE updated_datetime();

CREATE TRIGGER forum_posts_updated_at
    BEFORE UPDATE
    ON forum_posts
    FOR EACH ROW
    EXECUTE PROCEDURE updated_datetime();


