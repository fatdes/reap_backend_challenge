CREATE SEQUENCE auth_seq START 100001;
CREATE TABLE auth (
  id            TEXT                     NOT NULL DEFAULT 'AUTH-'::text || nextval('auth_seq'::regclass),
  username      TEXT                     NOT NULL UNIQUE,
  password      TEXT                     NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
  PRIMARY KEY   (id)
);
ALTER SEQUENCE auth_seq OWNED BY auth.id;

CREATE INDEX idx_auth_username ON auth (username);
