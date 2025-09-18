CREATE TABLE "accounts" (
  "id" bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,  -- 自增主键
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now()       -- 正确的时间默认值
);

CREATE TABLE "entries" (
  "id" bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,  -- 可正可负（收入/支出）
  "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE "transfers" (
  "id" bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,  -- 必须为正（转账金额）
  "created_at" timestamptz NOT NULL DEFAULT now()
);

-- 显式命名索引
CREATE INDEX "accounts_owner_idx" ON "accounts" ("owner");
CREATE INDEX "entries_account_id_idx" ON "entries" ("account_id");
CREATE INDEX "transfers_from_account_id_idx" ON "transfers" ("from_account_id");
CREATE INDEX "transfers_to_account_id_idx" ON "transfers" ("to_account_id");
CREATE INDEX "transfers_from_to_account_id_idx" ON "transfers" ("from_account_id", "to_account_id");

-- 修正拼写错误
COMMENT ON COLUMN "entries"."amount" IS 'can be positive or negative';
COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

-- 显式命名外键约束
ALTER TABLE "entries" ADD CONSTRAINT "entries_account_id_fkey"
  FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD CONSTRAINT "transfers_from_account_id_fkey"
  FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD CONSTRAINT "transfers_to_account_id_fkey"
  FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");