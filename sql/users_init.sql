-- users_init.sql — 用户表初始化脚本
-- 幂等：可重复执行（CREATE TABLE IF NOT EXISTS，COMMENT 本身可重复）
-- 手动执行：psql "postgres://user:password@localhost:5432/beid?sslmode=disable" -f users_init.sql
-- 注意：gen_random_uuid() 需要 PostgreSQL 13+（更高版本已内置，无需扩展）

CREATE TABLE IF NOT EXISTS users (
    id            uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    username      text        NOT NULL UNIQUE,
    email         text        NOT NULL UNIQUE,
    password_hash text        NOT NULL,
    status        smallint    NOT NULL DEFAULT 0
                  CHECK (status IN (0, 1, 2, 3)),
    created_at    timestamptz NOT NULL DEFAULT now()
);

COMMENT ON TABLE  users        IS '用户表';
COMMENT ON COLUMN users.status IS '用户状态：0=未激活，1=正常，2=被封禁，3=已删除';
