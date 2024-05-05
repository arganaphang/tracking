CREATE TABLE IF NOT EXISTS "orders"(
    "id" VARCHAR PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "description" TEXT NOT NULL,
    "note" TEXT,
    "status" VARCHAR NOT NULL,
    "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by" VARCHAR NOT NULL,
    "updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_by" VARCHAR NOT NULL
);