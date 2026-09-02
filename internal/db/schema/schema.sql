CREATE EXTENSION pgcrypto;

CREATE TABLE "fizzbuzz_stats" (
  "id"         uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "created_at" timestamptz NOT NULL DEFAULT (NOW()),
  "updated_at" timestamptz NOT NULL DEFAULT (NOW()),
  "int1"       integer NOT NULL,
  "int2"       integer NOT NULL,
  "limit"      integer NOT NULL,
  "str1"       text NOT NULL,
  "str2"       text NOT NULL,
  "hits"       bigint NOT NULL DEFAULT 1,
  UNIQUE ("int1", "int2", "limit", "str1", "str2")
);
