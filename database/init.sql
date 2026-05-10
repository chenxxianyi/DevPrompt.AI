-- DevPrompt AI Database Initialization
-- MySQL 8+

CREATE DATABASE IF NOT EXISTS devprompt_ai
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE devprompt_ai;

-- ============================================================
-- 1. users
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    username VARCHAR(64) NOT NULL,
    email VARCHAR(128) NOT NULL,
    password_hash VARCHAR(256) NOT NULL,
    avatar VARCHAR(512) NOT NULL DEFAULT '',
    role ENUM('user', 'admin') NOT NULL DEFAULT 'user',
    membership_level ENUM('free', 'pro', 'team', 'enterprise') NOT NULL DEFAULT 'free',
    membership_expired_at DATETIME NULL DEFAULT NULL,
    daily_generate_count INT NOT NULL DEFAULT 0,
    last_generate_date DATE NULL DEFAULT NULL,
    status ENUM('active', 'disabled') NOT NULL DEFAULT 'active',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_email (email),
    UNIQUE INDEX idx_username (username),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 2. prompt_categories
-- ============================================================
CREATE TABLE IF NOT EXISTS prompt_categories (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL,
    slug VARCHAR(64) NOT NULL,
    description VARCHAR(256) NOT NULL DEFAULT '',
    sort INT NOT NULL DEFAULT 0,
    status ENUM('active', 'disabled') NOT NULL DEFAULT 'active',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_slug (slug),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 3. prompt_templates
-- ============================================================
CREATE TABLE IF NOT EXISTS prompt_templates (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    category_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(256) NOT NULL,
    slug VARCHAR(256) NOT NULL,
    description VARCHAR(1024) NOT NULL DEFAULT '',
    content TEXT NOT NULL,
    tags JSON NOT NULL DEFAULT ('[]'),
    use_count INT NOT NULL DEFAULT 0,
    like_count INT NOT NULL DEFAULT 0,
    favorite_count INT NOT NULL DEFAULT 0,
    status ENUM('active', 'disabled') NOT NULL DEFAULT 'active',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_slug (slug),
    INDEX idx_category_id (category_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 4. prompt_favorites
-- ============================================================
CREATE TABLE IF NOT EXISTS prompt_favorites (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    prompt_template_id BIGINT UNSIGNED NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_user_prompt (user_id, prompt_template_id),
    INDEX idx_user_id (user_id),
    INDEX idx_prompt_template_id (prompt_template_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 5. prompt_likes
-- ============================================================
CREATE TABLE IF NOT EXISTS prompt_likes (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    prompt_template_id BIGINT UNSIGNED NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_user_prompt (user_id, prompt_template_id),
    INDEX idx_user_id (user_id),
    INDEX idx_prompt_template_id (prompt_template_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 6. generated_prompts
-- ============================================================
CREATE TABLE IF NOT EXISTS generated_prompts (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    type ENUM('project', 'cursor-rules', 'claude-code', 'optimize') NOT NULL,
    title VARCHAR(512) NOT NULL DEFAULT '',
    input JSON NOT NULL DEFAULT ('{}'),
    output TEXT NOT NULL,
    model VARCHAR(128) NOT NULL DEFAULT '',
    provider VARCHAR(64) NOT NULL DEFAULT '',
    tokens INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_user_id (user_id),
    INDEX idx_type (type),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 7. ai_models
-- ============================================================
CREATE TABLE IF NOT EXISTS ai_models (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    provider VARCHAR(64) NOT NULL,
    model_name VARCHAR(128) NOT NULL,
    display_name VARCHAR(256) NOT NULL,
    api_base_url VARCHAR(512) NOT NULL DEFAULT '',
    is_default TINYINT NOT NULL DEFAULT 0,
    status ENUM('active', 'disabled') NOT NULL DEFAULT 'active',
    priority INT NOT NULL DEFAULT 0,
    timeout_seconds INT NOT NULL DEFAULT 60,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_provider_model (provider, model_name),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 8. ai_call_logs
-- ============================================================
CREATE TABLE IF NOT EXISTS ai_call_logs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    provider VARCHAR(64) NOT NULL DEFAULT '',
    model VARCHAR(128) NOT NULL DEFAULT '',
    request_type VARCHAR(64) NOT NULL DEFAULT '',
    prompt_tokens INT NOT NULL DEFAULT 0,
    completion_tokens INT NOT NULL DEFAULT 0,
    total_tokens INT NOT NULL DEFAULT 0,
    status ENUM('success', 'failed') NOT NULL DEFAULT 'success',
    error_message VARCHAR(1024) NOT NULL DEFAULT '',
    latency_ms INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 9. membership_plans
-- ============================================================
CREATE TABLE IF NOT EXISTS membership_plans (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL,
    code ENUM('free', 'pro', 'team', 'enterprise') NOT NULL,
    price DECIMAL(10,2) NOT NULL DEFAULT 0,
    duration_days INT NOT NULL DEFAULT 0,
    daily_limit INT NOT NULL DEFAULT 5,
    features JSON NOT NULL DEFAULT ('[]'),
    status ENUM('active', 'disabled') NOT NULL DEFAULT 'active',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_code (code),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- 10. orders
-- ============================================================
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    plan_id BIGINT UNSIGNED NOT NULL,
    order_no VARCHAR(64) NOT NULL,
    amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    status ENUM('pending', 'paid', 'cancelled', 'refunded') NOT NULL DEFAULT 'pending',
    paid_at DATETIME NULL DEFAULT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_order_no (order_no),
    INDEX idx_user_id (user_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- Seed data: membership plans
-- ============================================================
INSERT INTO membership_plans (name, code, price, duration_days, daily_limit, features, status) VALUES
('Free', 'free', 0, 0, 5, '["基础 Prompt 模板", "每日 5 次生成", "社区支持"]', 'active'),
('Pro', 'pro', 29.00, 30, 100, '["所有 Prompt 模板", "每日 100 次生成", "支持 5 种 AI 模型", "优先支持"]', 'active'),
('Team', 'team', 99.00, 30, 500, '["所有 Pro 功能", "每日 500 次生成", "团队协作", "API 访问"]', 'active'),
('Enterprise', 'enterprise', 299.00, 365, 999999, '["所有 Team 功能", "无限生成", "专属支持", "私有化部署", "SLA 保障"]', 'active');

-- ============================================================
-- Seed data: prompt categories
-- ============================================================
INSERT INTO prompt_categories (name, slug, description, sort, status) VALUES
('项目开发', 'project-dev', '项目开发相关的 Prompt 模板', 1, 'active'),
('Cursor 规则', 'cursor-rules', 'Cursor 编辑器 .cursorrules 配置模板', 2, 'active'),
('Claude Code', 'claude-code', 'Claude Code CLI 任务 Prompt 模板', 3, 'active'),
('Prompt 优化', 'prompt-optimize', 'Prompt 优化与改进模板', 4, 'active'),
('代码生成', 'code-generation', '代码生成相关 Prompt 模板', 5, 'active'),
('代码审查', 'code-review', '代码审查相关 Prompt 模板', 6, 'active'),
('文档写作', 'documentation', '技术文档写作 Prompt 模板', 7, 'active'),
('测试', 'testing', '测试相关 Prompt 模板', 8, 'active');

-- ============================================================
-- Seed data: ai_models
-- ============================================================
INSERT INTO ai_models (provider, model_name, display_name, api_base_url, is_default, status, priority, timeout_seconds) VALUES
('openai', 'gpt-4o', 'GPT-4o', '', 1, 'active', 1, 60),
('openai', 'gpt-4o-mini', 'GPT-4o Mini', '', 0, 'active', 2, 30),
('claude', 'claude-sonnet-4-6', 'Claude Sonnet 4.6', '', 0, 'active', 3, 60),
('claude', 'claude-haiku-4-5', 'Claude Haiku 4.5', '', 0, 'active', 4, 30),
('deepseek', 'deepseek-chat', 'DeepSeek Chat', '', 0, 'active', 5, 60),
('gemini', 'gemini-2.0-flash', 'Gemini 2.0 Flash', '', 0, 'active', 6, 60),
('qwen', 'qwen-plus', 'Qwen Plus', '', 0, 'active', 7, 60);

-- ============================================================
-- Seed data: admin user (password: admin123)
-- ============================================================
INSERT INTO users (username, email, password_hash, role, membership_level, status) VALUES
('admin', 'admin@devprompt.ai', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin', 'enterprise', 'active');
