-- =====================================================
-- 安全的数据库初始化脚本
-- 支持多次执行，自动清理旧数据
-- =====================================================

-- 创建数据库
CREATE DATABASE IF NOT EXISTS aiops_cmdb DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE aiops_cmdb;

-- 添加授权
GRANT ALL PRIVILEGES ON aiops_cmdb.* TO 'ops_user'@'%';
FLUSH PRIVILEGES;

-- =====================================================
-- 创建表结构
-- =====================================================
CREATE TABLE IF NOT EXISTS `host_groups` (
                                             `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                                             `name` varchar(100) NOT NULL COMMENT '分组名称',
    `description` text COMMENT '描述',
    `parent_id` bigint unsigned DEFAULT '0' COMMENT '父分组ID',
    `sort` int DEFAULT '0' COMMENT '排序',
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_host_groups_name` (`name`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_deleted_at` (`deleted_at`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='主机分组表';

CREATE TABLE IF NOT EXISTS `hosts` (
                                       `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                                       `hostname` varchar(100) NOT NULL COMMENT '主机名',
    `ip` varchar(50) NOT NULL COMMENT '内网IP',
    `public_ip` varchar(50) DEFAULT NULL COMMENT '公网IP',
    `os` varchar(50) DEFAULT NULL COMMENT '操作系统',
    `os_version` varchar(50) DEFAULT NULL COMMENT '系统版本',
    `cpu` int DEFAULT NULL COMMENT 'CPU核数',
    `memory` int DEFAULT NULL COMMENT '内存大小(MB)',
    `disk` int DEFAULT NULL COMMENT '磁盘大小(GB)',
    `status` varchar(20) DEFAULT 'offline' COMMENT '状态: online, offline, maintenance',
    `environment` varchar(20) DEFAULT NULL COMMENT '环境: production, staging, development',
    `region` varchar(50) DEFAULT NULL COMMENT '区域',
    `idc` varchar(50) DEFAULT NULL COMMENT 'IDC机房',
    `cabinet` varchar(50) DEFAULT NULL COMMENT '机柜',
    `group_id` bigint unsigned DEFAULT NULL COMMENT '分组ID',
    `tags` json DEFAULT NULL COMMENT '标签',
    `labels` json DEFAULT NULL COMMENT '标签键值对',
    `remark` text COMMENT '备注',
    `agent_status` varchar(20) DEFAULT 'offline' COMMENT 'Agent状态',
    `last_seen_at` bigint DEFAULT NULL COMMENT '最后上报时间',
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_hosts_hostname` (`hostname`),
    KEY `idx_ip` (`ip`),
    KEY `idx_status` (`status`),
    KEY `idx_environment` (`environment`),
    KEY `idx_group_id` (`group_id`),
    KEY `idx_deleted_at` (`deleted_at`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='主机表';

CREATE TABLE IF NOT EXISTS `host_changes` (
                                              `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                                              `host_id` bigint unsigned NOT NULL COMMENT '主机ID',
                                              `change_type` varchar(50) NOT NULL COMMENT '变更类型',
    `old_value` text COMMENT '旧值',
    `new_value` text COMMENT '新值',
    `operator` varchar(100) DEFAULT NULL COMMENT '操作人',
    `remark` text COMMENT '备注',
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_host_id` (`host_id`),
    KEY `idx_deleted_at` (`deleted_at`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='主机变更记录表';

-- =====================================================
-- 清空旧数据（可选）
-- =====================================================
-- TRUNCATE TABLE hosts;
-- TRUNCATE TABLE host_groups;

-- =====================================================
-- 插入测试数据
-- =====================================================
INSERT IGNORE INTO `host_groups` (`name`, `description`, `parent_id`, `sort`, `created_at`, `updated_at`) VALUES
('Web服务器', 'Web应用服务器', 0, 1, NOW(), NOW()),
('数据库服务器', '数据库服务器', 0, 2, NOW(), NOW()),
('缓存服务器', 'Redis/Memcached', 0, 3, NOW(), NOW()),
('消息队列', 'RabbitMQ/Kafka', 0, 4, NOW(), NOW()),
('监控服务器', 'Prometheus/Grafana', 0, 5, NOW(), NOW());

INSERT IGNORE INTO `hosts` (`hostname`, `ip`, `os`, `os_version`, `cpu`, `memory`, `disk`, `status`, `environment`, `group_id`, `tags`, `created_at`, `updated_at`) VALUES
('web-server-01', '192.168.1.10', 'CentOS', '7.9', 8, 16384, 500, 'online', 'production', 1, '["web", "nginx"]', NOW(), NOW()),
('web-server-02', '192.168.1.11', 'CentOS', '7.9', 8, 16384, 500, 'online', 'production', 1, '["web", "nginx"]', NOW(), NOW()),
('web-server-03', '192.168.1.12', 'Ubuntu', '20.04', 8, 16384, 500, 'offline', 'staging', 1, '["web", "apache"]', NOW(), NOW()),
('db-master-01', '192.168.1.20', 'CentOS', '7.9', 16, 32768, 1000, 'online', 'production', 2, '["database", "mysql", "master"]', NOW(), NOW()),
('db-slave-01', '192.168.1.21', 'CentOS', '7.9', 16, 32768, 1000, 'online', 'production', 2, '["database", "mysql", "slave"]', NOW(), NOW()),
('db-slave-02', '192.168.1.22', 'CentOS', '7.9', 16, 32768, 1000, 'online', 'production', 2, '["database", "mysql", "slave"]', NOW(), NOW()),
('redis-01', '192.168.1.30', 'CentOS', '7.9', 4, 8192, 200, 'online', 'production', 3, '["cache", "redis"]', NOW(), NOW()),
('redis-02', '192.168.1.31', 'CentOS', '7.9', 4, 8192, 200, 'online', 'production', 3, '["cache", "redis"]', NOW(), NOW()),
('mq-01', '192.168.1.40', 'Ubuntu', '20.04', 8, 16384, 300, 'online', 'production', 4, '["mq", "rabbitmq"]', NOW(), NOW()),
('monitor-01', '192.168.1.50', 'Ubuntu', '20.04', 4, 8192, 200, 'online', 'production', 5, '["monitor", "prometheus"]', NOW(), NOW());

-- =====================================================
-- 显示统计信息
-- =====================================================
SELECT '✅ 数据库初始化完成！' AS message;
SELECT CONCAT('📊 主机分组: ', COUNT(*), ' 个') AS info FROM host_groups;
SELECT CONCAT('🖥️  主机总数: ', COUNT(*), ' 台') AS info FROM hosts;
SELECT CONCAT('✅ 在线主机: ', COUNT(*), ' 台') AS info FROM hosts WHERE status = 'online';
SELECT CONCAT('❌ 离线主机: ', COUNT(*), ' 台') AS info FROM hosts WHERE status = 'offline';
SELECT CONCAT('🔧 维护中: ', COUNT(*), ' 台') AS info FROM hosts WHERE status = 'maintenance';