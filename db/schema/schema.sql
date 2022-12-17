CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
  "user_id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "username" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "pass_hash" varchar NOT NULL,
  "join_date" timestamptz NOT NULL DEFAULT (now()),
  "bio" varchar,
  "profile_photo" varchar
);

CREATE TABLE "usersetting" (
  "user_id" uuid PRIMARY KEY,
  "darkmode" bool
);

CREATE TABLE "projects" (
  "project_id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "title" varchar NOT NULL,
  "info" varchar,
  "owner_id" uuid NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "tasks" (
  "task_id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "project_id" uuid NOT NULL,
  "title" varchar NOT NULL,
  "info" varchar,
  "tag" varchar,
  "created_by" uuid NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "beggining" timestamptz DEFAULT (now()),
  "deadline" timestamptz,
  "color" varchar
);

CREATE TABLE "usersinproject" (
  "project_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "added_at" timestamptz DEFAULT (now())
);

CREATE TABLE "usersintask" (
  "task_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "added_at" timestamptz DEFAULT (now())
);

CREATE TABLE "refresh_token" (
  "token" varchar PRIMARY KEY,
  "username" varchar NOT NULL,
  "fingerprint" varchar NOT NULL,
  "device" varchar NOT NULL
);

CREATE INDEX ON "users" ("user_id");

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "usersetting" ("user_id");

CREATE INDEX ON "projects" ("project_id");

CREATE INDEX ON "projects" ("owner_id");

CREATE INDEX ON "tasks" ("task_id");

CREATE INDEX ON "tasks" ("project_id");

CREATE INDEX ON "refresh_token" ("token");

ALTER TABLE "usersetting" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "projects" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("user_id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("project_id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("user_id");

ALTER TABLE "usersinproject" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("project_id");

ALTER TABLE "usersinproject" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "usersintask" ADD FOREIGN KEY ("task_id") REFERENCES "tasks" ("task_id");

ALTER TABLE "usersintask" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "usersinproject" ADD CONSTRAINT "uq_usersinproject" UNIQUE("project_id", "user_id");

ALTER TABLE "usersintask" ADD CONSTRAINT "uq_usersintask" UNIQUE("task_id", "user_id");

ALTER TABLE "refresh_token" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "refresh_token" ADD CONSTRAINT "uq_refresh_token" UNIQUE("username", "fingerprint");
