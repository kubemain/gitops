-- ==================== 智能运维平台权限系统数据库 ====================
-- 版本: v1.0

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建数据库
CREATE DATABASE IF NOT EXISTS aiops_user_service DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 添加授权
GRANT ALL PRIVILEGES ON aiops_user_service.* TO 'ops_user'@'%';
FLUSH PRIVILEGES;

-- 切换数据库
USE aiops_user_service;

-- ==================== 1. 用户表 ====================
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
                         `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
                         `username` VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名（登录用）',
                         `nickname` VARCHAR(50) DEFAULT '' COMMENT '昵称（显示用）',
                         `email` VARCHAR(100) DEFAULT '' COMMENT '邮箱',
                         `phone` VARCHAR(20) DEFAULT '' COMMENT '手机号',
                         `avatar` VARCHAR(255) DEFAULT '' COMMENT '头像URL',
                         `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希（bcrypt）',
                         `status` TINYINT DEFAULT 1 COMMENT '状态：1=启用 0=禁用',
                         `last_login_at` TIMESTAMP NULL COMMENT '最后登录时间',
                         `last_login_ip` VARCHAR(50) DEFAULT '' COMMENT '最后登录IP',
                         `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                         `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                         INDEX idx_username (`username`),
                         INDEX idx_email (`email`),
                         INDEX idx_status (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ==================== 2. 角色表 ====================
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles` (
                         `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '角色ID',
                         `name` VARCHAR(50) NOT NULL COMMENT '角色名称',
                         `code` VARCHAR(50) NOT NULL UNIQUE COMMENT '角色编码（如：admin）',
                         `description` VARCHAR(255) DEFAULT '' COMMENT '角色描述',
                         `sort_order` INT DEFAULT 0 COMMENT '排序（数字越小越靠前）',
                         `status` TINYINT DEFAULT 1 COMMENT '状态：1=启用 0=禁用',
                         `is_system` TINYINT DEFAULT 0 COMMENT '是否系统角色（1=不可删除）',
                         `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                         `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                         INDEX idx_code (`code`),
                         INDEX idx_status (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- ==================== 3. 权限表 ====================
DROP TABLE IF EXISTS `permissions`;
CREATE TABLE `permissions` (
                               `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '权限ID',
                               `parent_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '父权限ID（0为顶级）',
                               `name` VARCHAR(50) NOT NULL COMMENT '权限名称',
                               `code` VARCHAR(100) NOT NULL UNIQUE COMMENT '权限编码（如：user:create）',
                               `type` VARCHAR(20) DEFAULT 'menu' COMMENT '类型：menu=菜单 button=按钮 api=接口',
                               `resource` VARCHAR(50) DEFAULT '' COMMENT '资源（如：user）',
                               `action` VARCHAR(50) DEFAULT '' COMMENT '操作（如：create/view/edit/delete）',
                               `path` VARCHAR(200) DEFAULT '' COMMENT 'API路径（如：/api/users）',
                               `method` VARCHAR(10) DEFAULT '' COMMENT 'HTTP方法（GET/POST/PUT/DELETE）',
                               `sort_order` INT DEFAULT 0 COMMENT '排序',
                               `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                               `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                               INDEX idx_code (`code`),
                               INDEX idx_parent (`parent_id`),
                               INDEX idx_type (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- ==================== 4. 菜单表 ====================
DROP TABLE IF EXISTS `menus`;
CREATE TABLE `menus` (
                         `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '菜单ID',
                         `parent_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '父菜单ID',
                         `name` VARCHAR(50) NOT NULL COMMENT '菜单名称（路由name）',
                         `title` VARCHAR(50) NOT NULL COMMENT '菜单标题（显示用）',
                         `path` VARCHAR(200) DEFAULT '' COMMENT '路由路径（如：/system/user）',
                         `component` VARCHAR(200) DEFAULT '' COMMENT '组件路径（如：system/user/index）',
                         `icon` VARCHAR(50) DEFAULT '' COMMENT '图标',
                         `permission_code` VARCHAR(100) DEFAULT '' COMMENT '关联权限编码',
                         `sort_order` INT DEFAULT 0 COMMENT '排序',
                         `visible` TINYINT DEFAULT 1 COMMENT '是否可见',
                         `status` TINYINT DEFAULT 1 COMMENT '状态',
                         `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                         INDEX idx_parent (`parent_id`),
                         INDEX idx_perm_code (`permission_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';

-- ==================== 5. 部门表 ====================
DROP TABLE IF EXISTS `departments`;
CREATE TABLE `departments` (
                               `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '部门ID',
                               `parent_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '父部门ID',
                               `name` VARCHAR(50) NOT NULL COMMENT '部门名称',
                               `code` VARCHAR(50) DEFAULT '' COMMENT '部门编码',
                               `leader_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '负责人ID',
                               `sort_order` INT DEFAULT 0 COMMENT '排序',
                               `status` TINYINT DEFAULT 1 COMMENT '状态',
                               `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                               INDEX idx_parent (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门表';

-- ==================== 6. 用户角色关联表 ====================
DROP TABLE IF EXISTS `user_roles`;
CREATE TABLE `user_roles` (
                              `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                              `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
                              `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
                              `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              UNIQUE KEY uk_user_role (`user_id`, `role_id`),
                              INDEX idx_user (`user_id`),
                              INDEX idx_role (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- ==================== 7. 角色权限关联表 ====================
DROP TABLE IF EXISTS `role_permissions`;
CREATE TABLE `role_permissions` (
                                    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                                    `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
                                    `permission_id` BIGINT UNSIGNED NOT NULL COMMENT '权限ID',
                                    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    UNIQUE KEY uk_role_perm (`role_id`, `permission_id`),
                                    INDEX idx_role (`role_id`),
                                    INDEX idx_perm (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- ==================== 8. 操作日志表 ====================
DROP TABLE IF EXISTS `audit_logs`;
CREATE TABLE `audit_logs` (
                              `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                              `user_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '操作用户ID',
                              `username` VARCHAR(50) DEFAULT '' COMMENT '用户名',
                              `action` VARCHAR(50) NOT NULL COMMENT '操作动作（如：login/create_user）',
                              `resource` VARCHAR(50) DEFAULT '' COMMENT '资源类型（如：user）',
                              `resource_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '资源ID',
                              `method` VARCHAR(10) DEFAULT '' COMMENT 'HTTP方法',
                              `path` VARCHAR(255) DEFAULT '' COMMENT '请求路径',
                              `ip_address` VARCHAR(50) DEFAULT '' COMMENT '客户端IP',
                              `user_agent` VARCHAR(500) DEFAULT '' COMMENT 'User Agent',
                              `request_data` TEXT COMMENT '请求数据（JSON）',
                              `response_status` INT DEFAULT 0 COMMENT 'HTTP响应状态码',
                              `error_msg` TEXT COMMENT '错误信息',
                              `duration` INT DEFAULT 0 COMMENT '耗时（毫秒）',
                              `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                              INDEX idx_user (`user_id`),
                              INDEX idx_action (`action`),
                              INDEX idx_created (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- ==================== 初始化数据 ====================
-- 插入默认角色
INSERT INTO `roles` (`id`, `name`, `code`, `description`, `sort_order`, `is_system`, `status`) VALUES
                                                                                                   (1, '超级管理员', 'super_admin', '拥有系统所有权限', 1, 1, 1),
                                                                                                   (2, '系统管理员', 'admin', '管理用户、角色、权限', 2, 1, 1),
                                                                                                   (3, '运维工程师', 'ops', '运维操作权限', 3, 1, 1),
                                                                                                   (4, '开发者', 'developer', '开发和查看权限', 4, 1, 1),
                                                                                                   (5, '访客', 'guest', '只读权限', 5, 1, 1)
    ON DUPLICATE KEY UPDATE `name`=VALUES(`name`);

-- 插入默认管理员用户
-- 用户名: admin
-- 密码: admin123
-- 密码哈希使用 bcrypt: $2a$10$bd1P/aPzdmZOwr6YB/Qo9.RLUFwOvuPsGwf3l7zblGRRBVfeY3X8.
INSERT INTO `users` (`id`, `username`, `nickname`, `email`, `password_hash`, `status`) VALUES
    (1, 'admin', '超级管理员', 'admin@ops-platform.com', '$2a$10$bd1P/aPzdmZOwr6YB/Qo9.RLUFwOvuPsGwf3l7zblGRRBVfeY3X8.', 1)
    ON DUPLICATE KEY UPDATE `nickname`=VALUES(`nickname`);

-- 给管理员分配超级管理员角色
INSERT IGNORE INTO `user_roles` (`user_id`, `role_id`) VALUES (1, 1);

-- 插入基础权限（树形结构）
INSERT INTO `permissions` (`id`, `parent_id`, `name`, `code`, `type`, `resource`, `action`, `path`, `method`, `sort_order`) VALUES
-- 系统管理（一级菜单）
(1, 0, '系统管理', 'system', 'menu', '', '', '', '', 1),

-- 用户管理（二级菜单）
(2, 1, '用户管理', 'system:user', 'menu', 'user', '', '/system/user', '', 1),
(3, 2, '查看用户', 'system:user:view', 'api', 'user', 'view', '/api/v1/users', 'GET', 1),
(4, 2, '新增用户', 'system:user:create', 'api', 'user', 'create', '/api/v1/users', 'POST', 2),
(5, 2, '编辑用户', 'system:user:edit', 'api', 'user', 'edit', '/api/v1/users/:id', 'PUT', 3),
(6, 2, '删除用户', 'system:user:delete', 'api', 'user', 'delete', '/api/v1/users/:id', 'DELETE', 4),
(7, 2, '重置密码', 'system:user:reset_pwd', 'api', 'user', 'reset_pwd', '/api/v1/users/:id/password', 'PUT', 5),

-- 角色管理（二级菜单）
(8, 1, '角色管理', 'system:role', 'menu', 'role', '', '/system/role', '', 2),
(9, 8, '查看角色', 'system:role:view', 'api', 'role', 'view', '/api/v1/roles', 'GET', 1),
(10, 8, '新增角色', 'system:role:create', 'api', 'role', 'create', '/api/v1/roles', 'POST', 2),
(11, 8, '编辑角色', 'system:role:edit', 'api', 'role', 'edit', '/api/v1/roles/:id', 'PUT', 3),
(12, 8, '删除角色', 'system:role:delete', 'api', 'role', 'delete', '/api/v1/roles/:id', 'DELETE', 4),
(13, 8, '分配权限', 'system:role:assign', 'api', 'role', 'assign', '/api/v1/roles/:id/permissions', 'POST', 5),

-- 部门管理（二级菜单）
(14, 1, '部门管理', 'system:dept', 'menu', 'department', '', '/system/dept', '', 3),
(15, 14, '查看部门', 'system:dept:view', 'api', 'department', 'view', '/api/v1/departments', 'GET', 1),
(16, 14, '新增部门', 'system:dept:create', 'api', 'department', 'create', '/api/v1/departments', 'POST', 2),
(17, 14, '编辑部门', 'system:dept:edit', 'api', 'department', 'edit', '/api/v1/departments/:id', 'PUT', 3),
(18, 14, '删除部门', 'system:dept:delete', 'api', 'department', 'delete', '/api/v1/departments/:id', 'DELETE', 4)
    ON DUPLICATE KEY UPDATE `name`=VALUES(`name`);

-- 给超级管理员角色分配所有权限
INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`)
SELECT 1, id FROM `permissions`;

-- 给系统管理员角色分配系统管理权限
INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`)
SELECT 2, id FROM `permissions` WHERE `code` LIKE 'system%';

-- 插入菜单数据
INSERT INTO `menus` (`id`, `parent_id`, `name`, `title`, `path`, `component`, `icon`, `permission_code`, `sort_order`, `visible`, `status`) VALUES
                                                                                                                                                (1, 0, 'system', '系统管理', '/system', 'Layout', '⚙️', 'system', 99, 1, 1),
                                                                                                                                                (2, 1, 'user', '用户管理', '/system/user', 'system/user/index', '👥', 'system:user:view', 1, 1, 1),
                                                                                                                                                (3, 1, 'role', '角色管理', '/system/role', 'system/role/index', '🔐', 'system:role:view', 2, 1, 1),
                                                                                                                                                (4, 1, 'dept', '部门管理', '/system/dept', 'system/dept/index', '🏢', 'system:dept:view', 3, 1, 1)
    ON DUPLICATE KEY UPDATE `title`=VALUES(`title`);

-- 插入默认部门
INSERT INTO `departments` (`id`, `parent_id`, `name`, `code`, `sort_order`, `status`) VALUES
                                                                                          (1, 0, '总公司', 'ROOT', 1, 1),
                                                                                          (2, 1, '技术部', 'TECH', 1, 1),
                                                                                          (3, 2, '后端组', 'BACKEND', 1, 1),
                                                                                          (4, 2, '前端组', 'FRONTEND', 2, 1),
                                                                                          (5, 1, '运维部', 'OPS', 2, 1)
    ON DUPLICATE KEY UPDATE `name`=VALUES(`name`);

SET FOREIGN_KEY_CHECKS = 1;

-- ==================== 完成提示 ====================
SELECT '========================================' AS '';
SELECT '✅ 数据库初始化完成！' AS 'Status';
SELECT '========================================' AS '';
SELECT '默认管理员账号:' AS '';
SELECT '  用户名: admin' AS '';
SELECT '  密码: admin123' AS '';
SELECT '========================================' AS '';
SELECT CONCAT('当前数据库: ', DATABASE()) AS '';
SELECT CONCAT('表数量: ', COUNT(*)) AS '' FROM information_schema.tables WHERE table_schema = DATABASE();
SELECT '========================================' AS '';
