CREATE TABLE IF NOT EXISTS admin(
   id UUID NOT NULL PRIMARY KEY,
   admin_name TEXT, 
   admin_password TEXT,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   deleted_at TIMESTAMP
);
