CREATE TABLE feed_post (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  type TEXT NOT NULL,
  title TEXT NOT NULL,
  body TEXT NOT NULL,
  image_url TEXT,
  tags TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  author TEXT NOT NULL,
  cta TEXT,
  likes INT NOT NULL DEFAULT 0,
  comments INT NOT NULL DEFAULT 0,
  read_time TEXT
);
