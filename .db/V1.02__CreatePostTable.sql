CREATE SEQUENCE post_seq START 100001;
CREATE TABLE post (
  id            TEXT                     NOT NULL DEFAULT 'POST-'::text || nextval('post_seq'::regclass),
  username      TEXT                     NOT NULL,
  image         BYTEA                    NOT NULL,
  description   TEXT                     NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
  PRIMARY KEY   (id)
);
ALTER SEQUENCE post_seq OWNED BY post.id;

CREATE INDEX idx_post_username ON post (username);
CREATE INDEX idx_post_created_at ON post (created_at);
