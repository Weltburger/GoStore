CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "public"."admins"(
    "uuid"	uuid PRIMARY KEY DEFAULT uuid_generate_v1mc(),
    "name"	VARCHAR(60) NOT NULL UNIQUE,
    "email"	VARCHAR(120) NOT NULL UNIQUE,
    "password" VARCHAR(200) NOT NULL,
    "create_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp NOT NULL DEFAULT now()
);

CREATE TABLE "public"."products"(
    "uuid"	uuid PRIMARY KEY DEFAULT uuid_generate_v1mc(),
    "user_uuid" uuid NOT NULL,
    "name"	VARCHAR(60) NOT NULL,
    "description"	VARCHAR(1000) NOT NULL,
    "price" INTEGER NOT NULL,
    "quantity" INTEGER NOT NULL,
    "create_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp NOT NULL DEFAULT now(),
    "deleted_at" timestamp
);

CREATE TABLE "public"."orders"(
    "uuid"	uuid PRIMARY KEY DEFAULT uuid_generate_v1mc(),
    "product_uuid"	uuid NOT NULL,
    "price" INTEGER NOT NULL,
    "quantity" INTEGER NOT NULL,
    "email"	VARCHAR(100) NOT NULL,
    "status" INTEGER NOT NULL,
    "create_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp NOT NULL DEFAULT now(),
    "deleted_at" timestamp
);

Alter TABLE "public"."orders" add constraint fk_order_uuid foreign key ("product_uuid")
    references "public"."products"(uuid) on delete restrict on update cascade;