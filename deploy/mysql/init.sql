-- 初始化数据库（可选：MySQL 镜像会用环境变量 MYSQL_DATABASE 自动建库，此文件兜底/额外权限）
CREATE DATABASE IF NOT EXISTS fortune CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 若使用独立账号，可在下方授权（默认由 docker-compose 的 MYSQL_USER 自动创建）
-- GRANT ALL PRIVILEGES ON fortune.* TO 'fortune'@'%';
-- FLUSH PRIVILEGES;