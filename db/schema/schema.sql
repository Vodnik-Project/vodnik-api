CREATE TABLE "users" (
  "id" serial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "pass_hash" varchar NOT NULL,
  "join_date" timestamptz NOT NULL DEFAULT (now()),
  "bio" varchar,
  "profile_photo" varchar
);

CREATE TABLE "usersetting" (
  "id" serial PRIMARY KEY,
  "darkmode" bool
);

CREATE TABLE "projects" (
  "id" serial PRIMARY KEY,
  "title" varchar NOT NULL,
  "info" varchar,
  "owner_id" serial NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "tasks" (
  "id" serial PRIMARY KEY,
  "project_id" serial NOT NULL,
  "title" varchar NOT NULL,
  "info" varchar,
  "tag" varchar,
  "created_by" serial NOT NULL,
  "created_at" timestamptz DEFAULT (now()),
  "beggining" timestamptz DEFAULT (now()),
  "deadline" timestamptz,
  "color" varchar
);

CREATE TABLE "usersinproject" (
  "project_id" serial NOT NULL,
  "user_id" serial NOT NULL,
  "added_at" timestamptz DEFAULT (now())
);

CREATE TABLE "usersintask" (
  "task_id" serial NOT NULL,
  "user_id" serial NOT NULL,
  "added_at" timestamptz DEFAULT (now())
);

CREATE INDEX ON "users" ("id");

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "usersetting" ("id");

CREATE INDEX ON "projects" ("id");

CREATE INDEX ON "projects" ("owner_id");

CREATE INDEX ON "tasks" ("id");

CREATE INDEX ON "tasks" ("project_id");

ALTER TABLE "usersetting" ADD FOREIGN KEY ("id") REFERENCES "users" ("id");

ALTER TABLE "projects" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "usersinproject" ADD FOREIGN KEY ("project_id") REFERENCES "projects" ("id");

ALTER TABLE "usersinproject" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "usersintask" ADD FOREIGN KEY ("task_id") REFERENCES "tasks" ("id");

ALTER TABLE "usersintask" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "usersinproject" ADD CONSTRAINT "uq_usersinproject" UNIQUE("project_id", "user_id");

ALTER TABLE "usersintask" ADD CONSTRAINT "uq_usersintask" UNIQUE("task_id", "user_id");
