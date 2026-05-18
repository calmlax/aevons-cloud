/*
 Navicat MySQL Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80044
 Source Host           : 127.0.0.1:3306
 Source Schema         : aevons_db

 Target Server Type    : MySQL
 Target Server Version : 80044
 File Encoding         : 65001

 Date: 18/05/2026 18:05:52
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for gen_table
-- ----------------------------
DROP TABLE IF EXISTS `gen_table`;
CREATE TABLE `gen_table` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '编号',
  `table_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '表名',
  `table_comment` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '解释',
  `class_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '类名',
  `module_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '模块名称',
  `author` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '作者',
  `router` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '路由',
  `menu_id` bigint DEFAULT '0' COMMENT '上级菜单编号',
  `permission` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '权限标识',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '备注',
  `created_at` timestamp NOT NULL COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=269 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='代码生成表';

-- ----------------------------
-- Records of gen_table
-- ----------------------------
BEGIN;
INSERT INTO `gen_table` VALUES (251, 'sys_conf', '参数配置表', 'Conf', 'sys', '', 'conf', 0, 'sys:conf', '', '2026-04-16 15:41:12', '2026-05-18 17:22:42');
INSERT INTO `gen_table` VALUES (252, 'sys_oauth_client', '终端应用', 'OauthClient', 'sys', '', 'oauth/client', 100000, 'sys:oauth:client', '', '2026-04-16 15:41:12', '2026-05-18 17:22:51');
INSERT INTO `gen_table` VALUES (253, 'sys_menu', '菜单权限表', 'Menu', 'sys', '', 'menu', 0, 'sys:menu', '', '2026-04-18 08:31:34', '2026-05-18 17:22:55');
INSERT INTO `gen_table` VALUES (254, 'sys_menu_tl', '菜单翻译', 'MenuTl', 'sys', '', 'menu/tl', 0, 'sys:menu:tl', '', '2026-04-18 08:31:34', '2026-04-18 09:23:19');
INSERT INTO `gen_table` VALUES (256, 'sys_role', '角色信息表', 'Role', 'sys', '', 'role', 0, 'sys:role', '', '2026-04-18 14:47:43', '2026-04-18 14:47:48');
INSERT INTO `gen_table` VALUES (260, 'sys_dept', '部门表', 'Dept', 'sys', '', 'dept', 0, 'sys:dept', '', '2026-04-18 20:11:12', '2026-04-18 20:11:20');
INSERT INTO `gen_table` VALUES (261, 'sys_post', '岗位信息表', 'Post', 'sys', '', 'post', 0, 'sys:post', '', '2026-04-18 20:11:12', '2026-04-18 20:11:28');
INSERT INTO `gen_table` VALUES (262, 'sys_user', '用户信息表', 'User', 'sys', '', 'user', 0, 'sys:user', '', '2026-04-19 09:07:00', '2026-04-19 09:07:11');
INSERT INTO `gen_table` VALUES (263, 'sys_notice', '通知公告', 'Notice', 'sys', '', 'notice', 100000, 'sys:notice', '', '2026-04-19 11:42:31', '2026-04-21 14:35:06');
INSERT INTO `gen_table` VALUES (264, 'sys_lang', '语言', 'Lang', 'sys', '', 'lang', 100221, 'sys:lang', '', '2026-04-19 12:44:00', '2026-04-19 12:56:30');
INSERT INTO `gen_table` VALUES (265, 'sys_lang_resource', '语言资源', 'LangResource', 'sys', '', 'lang/resource', 100221, 'sys:lang:resource', '', '2026-04-19 12:44:00', '2026-04-19 12:56:37');
INSERT INTO `gen_table` VALUES (266, 'sys_job', '定时任务配置表', 'Job', 'sys', '', 'job', 0, 'sys:job', '', '2026-04-19 15:02:14', '2026-04-19 15:02:58');
INSERT INTO `gen_table` VALUES (267, 'sys_job_log', '定时任务执行日志表', 'JobLog', 'sys', '', 'job/log', 0, 'sys:job:log', '', '2026-04-19 15:02:14', '2026-04-19 15:03:04');
INSERT INTO `gen_table` VALUES (268, 'sys_user_credential', '用户Passkey凭据', 'UserCredential', 'sys', '', 'user/credential', 0, 'sys:user:credential', '', '2026-04-25 20:33:32', '2026-04-25 20:34:01');
COMMIT;

-- ----------------------------
-- Table structure for gen_table_column
-- ----------------------------
DROP TABLE IF EXISTS `gen_table_column`;
CREATE TABLE `gen_table_column` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '编号',
  `table_id` bigint NOT NULL COMMENT '所属表编号',
  `column_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'DB字段名称',
  `column_comment` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'DB字段解释',
  `column_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'DB字段数据类型',
  `data_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '数据类型',
  `field_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '字段名称',
  `json` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'JSON',
  `is_primary_key` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否主键（0否，1是）',
  `is_auto_increment` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否自增（0否，1是）',
  `is_required` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否必填（0否，1是）',
  `is_insert` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否插入字段（0否，1是）',
  `is_edit` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否编辑字段（0否，1是）',
  `is_list` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否列表字段（0否，1是）',
  `sortable` tinyint(1) NOT NULL DEFAULT '0' COMMENT '列表排序字段（0否，1是）',
  `filterable` tinyint(1) NOT NULL DEFAULT '0' COMMENT '列表筛选字段（0否，1是）',
  `condition` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '查询条件',
  `dict_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '字典类型',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序',
  `component` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'UI组件',
  `default_value` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '默认值',
  `data_length` int DEFAULT NULL COMMENT '数据长款',
  `data_precision` tinyint(1) DEFAULT NULL COMMENT '数据精度',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2613 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='代码生成表字段';

-- ----------------------------
-- Records of gen_table_column
-- ----------------------------
BEGIN;
INSERT INTO `gen_table_column` VALUES (2398, 251, 'id', '编号', 'bigint', 'int64', 'Id', 'id,string', 1, 1, 1, 0, 1, 1, 0, 0, '', '', 1, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2399, 251, 'name', '名称', 'varchar(255)', 'string', 'Name', 'name', 0, 0, 1, 1, 1, 1, 0, 0, 'like', '', 2, '', '', 255, 0);
INSERT INTO `gen_table_column` VALUES (2400, 251, 'conf_key', '配置KEY', 'varchar(64)', 'string', 'ConfKey', 'confKey', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 3, '', '', 64, 0);
INSERT INTO `gen_table_column` VALUES (2401, 251, 'conf_value', '配置值', 'text', 'string', 'ConfValue', 'confValue', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 4, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2402, 251, 'is_sys', '是否系统内置（0否 1是）', 'tinyint(1)', 'int16', 'IsSys', 'isSys', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 5, '', '', 1, 0);
INSERT INTO `gen_table_column` VALUES (2403, 251, 'scope', '范围（0公开配置，1登录配置，2后台服务配置）', 'tinyint(1)', 'int16', 'Scope', 'scope', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 6, '', '', 1, 0);
INSERT INTO `gen_table_column` VALUES (2404, 251, 'is_secret', '是否加密/脱敏（0否 1是）', 'tinyint(1)', 'int16', 'IsSecret', 'isSecret', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 7, '', '', 1, 0);
INSERT INTO `gen_table_column` VALUES (2405, 251, 'remark', '备注', 'varchar(500)', 'string', 'Remark', 'remark', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 8, '', '', 500, 0);
INSERT INTO `gen_table_column` VALUES (2406, 251, 'created_by', '', 'bigint', 'int64', 'CreatedBy', 'createdBy,string', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 9, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2407, 251, 'created_at', '', 'timestamp', 'time.Time', 'CreatedAt', 'createdAt', 0, 0, 0, 1, 1, 1, 0, 0, 'between', '', 10, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2408, 251, 'updated_by', '', 'bigint', 'int64', 'UpdatedBy', 'updatedBy,string', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 11, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2409, 251, 'updated_at', '', 'timestamp', 'time.Time', 'UpdatedAt', 'updatedAt', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 12, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2410, 252, 'id', 'ID', 'bigint', 'int64', 'Id', 'id,string', 1, 1, 1, 0, 1, 1, 0, 0, '', '', 1, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2411, 252, 'client_id', '客户端ID', 'varchar(32)', 'string', 'ClientId', 'clientId', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', '', 2, 'input', '', 32, 0);
INSERT INTO `gen_table_column` VALUES (2412, 252, 'client_secret', '客户端秘钥', 'varchar(256)', 'string', 'ClientSecret', 'clientSecret', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 3, 'input', '', 256, 0);
INSERT INTO `gen_table_column` VALUES (2413, 252, 'client_name', '客户端名称', 'varchar(255)', 'string', 'ClientName', 'clientName', 0, 0, 1, 1, 1, 1, 0, 0, 'like', '', 4, 'input', '', 255, 0);
INSERT INTO `gen_table_column` VALUES (2414, 252, 'logo_uri', '客户端LOGO', 'varchar(255)', 'string', 'LogoUri', 'logoUri', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 5, 'input', '', 255, 0);
INSERT INTO `gen_table_column` VALUES (2415, 252, 'scope', '授权范围', 'varchar(256)', 'string', 'Scope', 'scope', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 6, 'select', '', 256, 0);
INSERT INTO `gen_table_column` VALUES (2416, 252, 'authorized_grant_types', '授权类型', 'varchar(256)', 'string', 'AuthorizedGrantTypes', 'authorizedGrantTypes', 0, 0, 1, 1, 1, 1, 0, 0, 'in', 'sys_authorized_grant_type', 7, 'select', '', 256, 0);
INSERT INTO `gen_table_column` VALUES (2417, 252, 'web_server_redirect_uri', '回调地址', 'varchar(256)', 'string', 'WebServerRedirectUri', 'webServerRedirectUri', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 8, 'input', '', 256, 0);
INSERT INTO `gen_table_column` VALUES (2418, 252, 'access_token_validity', '访问令牌有效期（秒）', 'int', 'int', 'AccessTokenValidity', 'accessTokenValidity', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 9, 'input-number', '43200', 0, 0);
INSERT INTO `gen_table_column` VALUES (2419, 252, 'refresh_token_validity', '刷新令牌有效期（秒）', 'int', 'int', 'RefreshTokenValidity', 'refreshTokenValidity', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 10, 'input', '2592000', 0, 0);
INSERT INTO `gen_table_column` VALUES (2420, 252, 'autoapprove', '自动授权（0否，1是）', 'tinyint(1)', 'int16', 'Autoapprove', 'autoapprove', 0, 0, 1, 1, 1, 1, 0, 0, '', 'sys_is', 11, 'switch', '0', 1, 0);
INSERT INTO `gen_table_column` VALUES (2421, 252, 'created_by', '创建人', 'bigint', 'int64', 'CreatedBy', 'createdBy,string', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 12, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2422, 252, 'created_at', '创建时间', 'timestamp', 'time.Time', 'CreatedAt', 'createdAt', 0, 0, 0, 0, 0, 1, 0, 0, 'between', '', 13, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2423, 252, 'updated_by', '修改人', 'bigint', 'int64', 'UpdatedBy', 'updatedBy,string', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 14, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2424, 252, 'updated_at', '更新时间', 'timestamp', 'time.Time', 'UpdatedAt', 'updatedAt', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 15, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2425, 253, 'id', '菜单ID', 'bigint', 'int64', 'Id', 'id,string', 1, 1, 1, 0, 1, 1, 0, 0, '', '', 1, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2426, 253, 'parent_id', '父级菜单Id', 'bigint', 'int64', 'ParentId', 'parentId,string', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', '', 2, 'tree-select', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2427, 253, 'title', '菜单名称', 'varchar(20)', 'string', 'Title', 'title', 0, 0, 1, 1, 1, 1, 0, 0, 'like', '', 3, 'input', '', 20, 0);
INSERT INTO `gen_table_column` VALUES (2428, 253, 'type', '类型（1目录 2菜单 3按钮）', 'tinyint(1)', 'int16', 'Type', 'type', 0, 0, 1, 1, 1, 1, 0, 0, 'in', '', 4, 'radio', '', 1, 0);
INSERT INTO `gen_table_column` VALUES (2429, 253, 'sort', '顺序', 'bigint', 'int64', 'Sort', 'sort,string', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 5, 'input-number', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2430, 253, 'path', '路由地址', 'varchar(100)', 'string', 'Path', 'path', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 6, 'input', '', 100, 0);
INSERT INTO `gen_table_column` VALUES (2431, 253, 'component', '组件路径', 'varchar(100)', 'string', 'Component', 'component', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 7, 'input', '', 100, 0);
INSERT INTO `gen_table_column` VALUES (2432, 253, 'query', '路由参数', 'varchar(255)', 'string', 'Query', 'query', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 8, 'input', '', 255, 0);
INSERT INTO `gen_table_column` VALUES (2433, 253, 'visible', '是否可见（0隐藏 1显示）', 'tinyint(1)', 'int16', 'Visible', 'visible', 0, 0, 1, 1, 1, 1, 0, 0, '', 'sys_is', 9, 'switch', '', 1, 0);
INSERT INTO `gen_table_column` VALUES (2434, 253, 'status', '状态（0正常 1停用）', 'tinyint(1)', 'int16', 'Status', 'status', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', 'sys_common_status', 10, 'radio', '', 1, 0);
INSERT INTO `gen_table_column` VALUES (2435, 253, 'is_frame', '是否为外链（0否 1是）', 'tinyint(1)', 'int16', 'IsFrame', 'isFrame', 0, 0, 1, 1, 1, 1, 0, 0, '', 'sys_is', 11, 'switch', '', 1, 0);
INSERT INTO `gen_table_column` VALUES (2436, 253, 'permission', '权限标识', 'varchar(32)', 'string', 'Permission', 'permission', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 12, 'input', '', 32, 0);
INSERT INTO `gen_table_column` VALUES (2437, 253, 'icon', '图标', 'varchar(64)', 'string', 'Icon', 'icon', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 13, 'input', '', 64, 0);
INSERT INTO `gen_table_column` VALUES (2438, 253, 'title_key', '菜单key', 'varchar(64)', 'string', 'TitleKey', 'titleKey', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 14, '', '', 64, 0);
INSERT INTO `gen_table_column` VALUES (2439, 253, 'active_id', '激活ID', 'bigint', 'int64', 'ActiveId', 'activeId,string', 0, 0, 0, 1, 1, 1, 0, 0, 'eq', '', 15, 'tree-select', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2440, 253, 'created_by', '', 'bigint', 'int64', 'CreatedBy', 'createdBy,string', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 16, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2441, 253, 'created_at', '', 'timestamp', 'time.Time', 'CreatedAt', 'createdAt', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 17, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2442, 253, 'updated_by', '', 'bigint', 'int64', 'UpdatedBy', 'updatedBy,string', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 18, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2443, 253, 'updated_at', '', 'timestamp', 'time.Time', 'UpdatedAt', 'updatedAt', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 19, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2444, 254, 'menu_id', '菜单Id', 'bigint', 'int64', 'MenuId', 'menuId,string', 1, 0, 1, 0, 1, 1, 0, 0, 'eq', '', 1, '', '0', 0, 0);
INSERT INTO `gen_table_column` VALUES (2445, 254, 'lang_code', '语言标识', 'varchar(10)', 'string', 'LangCode', 'langCode', 1, 0, 1, 0, 1, 1, 0, 0, '', '', 2, '', '', 10, 0);
INSERT INTO `gen_table_column` VALUES (2446, 254, 'title', '菜单名称', 'varchar(20)', 'string', 'Title', 'title', 0, 0, 1, 1, 1, 1, 0, 0, 'like', '', 3, '', '', 20, 0);
INSERT INTO `gen_table_column` VALUES (2458, 256, 'id', '角色ID', 'bigint', 'int64', 'Id', 'id,string', 1, 1, 1, 0, 1, 1, 0, 0, '', '', 1, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2459, 256, 'role_name', '角色名称', 'varchar(50)', 'string', 'RoleName', 'roleName', 0, 0, 1, 1, 1, 1, 0, 0, 'like', '', 2, '', '', 50, 0);
INSERT INTO `gen_table_column` VALUES (2460, 256, 'role_key', '角色权限字符串', 'varchar(32)', 'string', 'RoleKey', 'roleKey', 0, 0, 0, 1, 1, 1, 0, 0, 'like', '', 3, '', '', 32, 0);
INSERT INTO `gen_table_column` VALUES (2461, 256, 'sort', '显示顺序', 'bigint', 'int64', 'Sort', 'sort,string', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 4, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2462, 256, 'data_scope', '数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）', 'tinyint(1)', 'int16', 'DataScope', 'dataScope', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', '', 5, '', '', 1, 0);
INSERT INTO `gen_table_column` VALUES (2463, 256, 'status', '状态（0正常 1停用）', 'tinyint(1)', 'int16', 'Status', 'status', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', '', 6, '', '', 1, 0);
INSERT INTO `gen_table_column` VALUES (2464, 256, 'remark', '备注', 'varchar(500)', 'string', 'Remark', 'remark', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 7, '', '', 500, 0);
INSERT INTO `gen_table_column` VALUES (2465, 256, 'created_by', '', 'bigint', 'int64', 'CreatedBy', 'createdBy,string', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 8, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2466, 256, 'created_at', '', 'timestamp', 'time.Time', 'CreatedAt', 'createdAt', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 9, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2467, 256, 'updated_by', '', 'bigint', 'int64', 'UpdatedBy', 'updatedBy,string', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 10, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2468, 256, 'updated_at', '', 'timestamp', 'time.Time', 'UpdatedAt', 'updatedAt', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 11, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2503, 260, 'id', '部门编号', 'bigint', 'int64', 'Id', 'id,string', 1, 1, 1, 0, 1, 1, 0, 0, '', '', 1, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2504, 260, 'parent_id', '父级编号', 'bigint', 'int64', 'ParentId', 'parentId,string', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', '', 2, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2505, 260, 'ancestors', '祖级', 'varchar(255)', 'string', 'Ancestors', 'ancestors', 0, 0, 0, 0, 0, 0, 0, 0, 'like', '', 3, '', '', 255, 0);
INSERT INTO `gen_table_column` VALUES (2506, 260, 'dept_type', '类型（1机构，2部门）', 'tinyint', 'int16', 'DeptType', 'deptType', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', 'sys_dept_type', 4, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2507, 260, 'dept_name', '部门名称', 'varchar(30)', 'string', 'DeptName', 'deptName', 0, 0, 1, 1, 1, 1, 0, 0, 'like', '', 5, '', '', 30, 0);
INSERT INTO `gen_table_column` VALUES (2508, 260, 'sort', '顺序', 'int', 'int', 'Sort', 'sort', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 6, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2509, 260, 'status', '状态（0正常 1停用）', 'tinyint(1)', 'int16', 'Status', 'status', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', 'sys_common_status', 7, '', '', 1, 0);
INSERT INTO `gen_table_column` VALUES (2510, 260, 'remark', '备注', 'varchar(255)', 'string', 'Remark', 'remark', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 8, '', '', 255, 0);
INSERT INTO `gen_table_column` VALUES (2511, 260, 'created_by', '', 'bigint', 'int64', 'CreatedBy', 'createdBy,string', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 9, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2512, 260, 'created_at', '', 'timestamp', 'time.Time', 'CreatedAt', 'createdAt', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 10, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2513, 260, 'updated_by', '', 'bigint', 'int64', 'UpdatedBy', 'updatedBy,string', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 11, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2514, 260, 'updated_at', '', 'timestamp', 'time.Time', 'UpdatedAt', 'updatedAt', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 12, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2515, 261, 'id', '岗位编号', 'bigint', 'int64', 'Id', 'id,string', 1, 1, 1, 0, 1, 1, 0, 0, '', '', 1, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2516, 261, 'post_key', '岗位标识', 'varchar(32)', 'string', 'PostKey', 'postKey', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 2, '', '', 32, 0);
INSERT INTO `gen_table_column` VALUES (2517, 261, 'post_name', '岗位名称', 'varchar(50)', 'string', 'PostName', 'postName', 0, 0, 1, 1, 1, 1, 0, 0, 'like', '', 3, '', '', 50, 0);
INSERT INTO `gen_table_column` VALUES (2518, 261, 'sort', '顺序', 'int', 'int', 'Sort', 'sort', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 4, 'input-number', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2519, 261, 'dept_id', '归属部门', 'bigint', 'int64', 'DeptId', 'deptId,string', 0, 0, 0, 1, 1, 1, 0, 0, 'eq', '', 5, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2520, 261, 'status', '状态（0正常 1停用）', 'tinyint(1)', 'int16', 'Status', 'status', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', 'sys_common_status', 6, 'switch', '', 1, 0);
INSERT INTO `gen_table_column` VALUES (2521, 261, 'remark', '备注', 'varchar(500)', 'string', 'Remark', 'remark', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 7, 'textarea', '', 500, 0);
INSERT INTO `gen_table_column` VALUES (2522, 261, 'created_by', '', 'bigint', 'int64', 'CreatedBy', 'createdBy,string', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 8, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2523, 261, 'created_at', '', 'timestamp', 'time.Time', 'CreatedAt', 'createdAt', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 9, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2524, 261, 'updated_by', '', 'bigint', 'int64', 'UpdatedBy', 'updatedBy,string', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 10, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2525, 261, 'updated_at', '', 'timestamp', 'time.Time', 'UpdatedAt', 'updatedAt', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 11, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2526, 262, 'id', '用户编号', 'bigint', 'int64', 'Id', 'id,string', 1, 1, 1, 0, 1, 1, 0, 0, '', '', 1, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2527, 262, 'username', '用户名', 'varchar(32)', 'string', 'Username', 'username', 0, 0, 1, 1, 1, 1, 0, 0, 'like', '', 2, '', '', 32, 0);
INSERT INTO `gen_table_column` VALUES (2528, 262, 'nickname', '呢称', 'varchar(50)', 'string', 'Nickname', 'nickname', 0, 0, 1, 1, 1, 1, 0, 0, 'like', '', 3, '', '', 50, 0);
INSERT INTO `gen_table_column` VALUES (2529, 262, 'type', '用户类型', 'tinyint', 'int16', 'Type', 'type', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', 'sys_user_type', 4, 'radio', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2530, 262, 'email', '邮箱', 'varchar(64)', 'string', 'Email', 'email', 0, 0, 0, 1, 1, 1, 0, 0, 'like_l', '', 5, '', '', 64, 0);
INSERT INTO `gen_table_column` VALUES (2531, 262, 'mobile', '手机号', 'varchar(11)', 'string', 'Mobile', 'mobile', 0, 0, 0, 1, 1, 1, 0, 0, 'like_l', '', 6, '', '', 11, 0);
INSERT INTO `gen_table_column` VALUES (2532, 262, 'sex', '性别（0未知 1男 2女）', 'tinyint', 'int16', 'Sex', 'sex', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', 'sys_user_sex', 7, 'radio', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2533, 262, 'avatar', '头像', 'varchar(255)', 'string', 'Avatar', 'avatar', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 8, '', '', 255, 0);
INSERT INTO `gen_table_column` VALUES (2534, 262, 'autograph', '电子签名', 'varchar(255)', 'string', 'Autograph', 'autograph', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 9, '', '', 255, 0);
INSERT INTO `gen_table_column` VALUES (2535, 262, 'password', '密码', 'varchar(128)', 'string', 'Password', 'password', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 10, '', '', 128, 0);
INSERT INTO `gen_table_column` VALUES (2536, 262, 'status', '状态（0正常 1停用）', 'tinyint', 'int16', 'Status', 'status', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', 'sys_common_status', 11, 'radio', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2537, 262, 'created_by', '', 'bigint', 'int64', 'CreatedBy', 'createdBy,string', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 12, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2538, 262, 'created_at', '', 'timestamp', 'time.Time', 'CreatedAt', 'createdAt', 0, 0, 0, 0, 0, 1, 0, 0, '', '', 13, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2539, 262, 'updated_by', '', 'bigint', 'int64', 'UpdatedBy', 'updatedBy,string', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 14, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2540, 262, 'updated_at', '', 'timestamp', 'time.Time', 'UpdatedAt', 'updatedAt', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 15, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2541, 263, 'id', '公告ID', 'int', 'int', 'Id', 'id', 1, 1, 1, 0, 1, 1, 0, 0, '', '', 1, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2542, 263, 'title', '公告标题', 'varchar(50)', 'string', 'Title', 'title', 0, 0, 1, 1, 1, 1, 0, 0, 'like', '', 2, '', '', 50, 0);
INSERT INTO `gen_table_column` VALUES (2543, 263, 'type', '公告类型（1通知 2公告）', 'tinyint(1)', 'int16', 'Type', 'type', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', 'sys_notice_type', 3, 'radio', '0', 1, 0);
INSERT INTO `gen_table_column` VALUES (2544, 263, 'content', '公告内容', 'varchar(3000)', 'string', 'Content', 'content', 0, 0, 1, 1, 1, 1, 0, 0, 'like', '', 4, '', '', 3000, 0);
INSERT INTO `gen_table_column` VALUES (2545, 263, 'status', '状态（0正常 1关闭）', 'tinyint(1)', 'int16', 'Status', 'status', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', 'sys_common_status', 5, 'radio', '0', 1, 0);
INSERT INTO `gen_table_column` VALUES (2546, 263, 'remark', '备注', 'varchar(255)', 'string', 'Remark', 'remark', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 6, '', '', 255, 0);
INSERT INTO `gen_table_column` VALUES (2547, 263, 'created_by', '创建者', 'bigint', 'int64', 'CreatedBy', 'createdBy,string', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 7, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2548, 263, 'created_at', '创建时间', 'timestamp', 'time.Time', 'CreatedAt', 'createdAt', 0, 0, 0, 0, 0, 1, 0, 0, '', '', 8, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2549, 263, 'updated_by', '更新者', 'bigint', 'int64', 'UpdatedBy', 'updatedBy,string', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 9, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2550, 263, 'updated_at', '更新时间', 'timestamp', 'time.Time', 'UpdatedAt', 'updatedAt', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 10, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2551, 264, 'id', '编号', 'bigint', 'int64', 'Id', 'id,string', 1, 1, 1, 0, 1, 1, 0, 0, '', '', 1, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2552, 264, 'lang_code', '语言编码（如zh-CN、en-US）', 'varchar(10)', 'string', 'LangCode', 'langCode', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', '', 2, 'input', '', 10, 0);
INSERT INTO `gen_table_column` VALUES (2553, 264, 'lang_name', '语言名称（如简体中文、English）', 'varchar(50)', 'string', 'LangName', 'langName', 0, 0, 1, 1, 1, 1, 0, 0, 'like', '', 3, 'input', '', 50, 0);
INSERT INTO `gen_table_column` VALUES (2554, 264, 'is_default', '是否默认语言（0否，1是）', 'tinyint(1)', 'int16', 'IsDefault', 'isDefault', 0, 0, 1, 1, 1, 1, 0, 0, '', 'sys_is', 4, 'switch', '0', 1, 0);
INSERT INTO `gen_table_column` VALUES (2555, 264, 'sort', '排序值（升序）', 'int', 'int', 'Sort', 'sort', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 5, 'input-number', '0', 0, 0);
INSERT INTO `gen_table_column` VALUES (2556, 264, 'status', '状态（0正常,1停用）', 'tinyint(1)', 'int16', 'Status', 'status', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', 'sys_common_status', 6, 'radio', '0', 1, 0);
INSERT INTO `gen_table_column` VALUES (2557, 264, 'remark', '备注', 'varchar(200)', 'string', 'Remark', 'remark', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 7, '', '', 200, 0);
INSERT INTO `gen_table_column` VALUES (2558, 264, 'created_at', '创建时间', 'timestamp', 'time.Time', 'CreatedAt', 'createdAt', 0, 0, 0, 0, 0, 1, 0, 0, '', '', 8, '', 'CURRENT_TIMESTAMP', 0, 0);
INSERT INTO `gen_table_column` VALUES (2559, 264, 'updated_at', '更新时间', 'timestamp', 'time.Time', 'UpdatedAt', 'updatedAt', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 9, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2560, 265, 'id', '主键ID', 'bigint', 'int64', 'Id', 'id,string', 1, 1, 1, 0, 1, 1, 0, 0, '', '', 1, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2561, 265, 'resource_key', '资源标识', 'varchar(128)', 'string', 'ResourceKey', 'resourceKey', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', '', 2, 'input', '', 128, 0);
INSERT INTO `gen_table_column` VALUES (2562, 265, 'namespace', '命名空间（default）', 'varchar(32)', 'string', 'Namespace', 'namespace', 0, 0, 1, 1, 1, 1, 0, 0, 'like', '', 3, 'input', 'default', 32, 0);
INSERT INTO `gen_table_column` VALUES (2563, 265, 'lang_code', '语言编码（如zh）', 'varchar(10)', 'string', 'LangCode', 'langCode', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', '', 4, 'input', '', 10, 0);
INSERT INTO `gen_table_column` VALUES (2564, 265, 'content', '内容', 'varchar(500)', 'string', 'Content', 'content', 0, 0, 1, 1, 1, 1, 0, 0, 'like', '', 5, 'textarea', '', 500, 0);
INSERT INTO `gen_table_column` VALUES (2565, 265, 'status', '状态（0正常,1停用）', 'tinyint(1)', 'int16', 'Status', 'status', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', 'sys_common_status', 6, 'switch', '0', 1, 0);
INSERT INTO `gen_table_column` VALUES (2567, 265, 'created_at', '创建时间', 'timestamp', 'time.Time', 'CreatedAt', 'createdAt', 0, 0, 0, 0, 0, 1, 0, 0, '', '', 8, '', 'CURRENT_TIMESTAMP', 0, 0);
INSERT INTO `gen_table_column` VALUES (2568, 265, 'updated_at', '更新时间', 'timestamp', 'time.Time', 'UpdatedAt', 'updatedAt', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 9, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2569, 266, 'id', '主键ID', 'bigint', 'int64', 'Id', 'id,string', 1, 1, 1, 0, 1, 1, 0, 0, '', '', 1, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2570, 266, 'job_name', '任务名称', 'varchar(64)', 'string', 'JobName', 'jobName', 0, 0, 1, 1, 1, 1, 0, 0, 'like', '', 2, 'input', '', 64, 0);
INSERT INTO `gen_table_column` VALUES (2571, 266, 'job_group', '任务分组', 'varchar(64)', 'string', 'JobGroup', 'jobGroup', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 3, 'input', '', 64, 0);
INSERT INTO `gen_table_column` VALUES (2572, 266, 'job_key', '任务唯一标识(不可重复)', 'varchar(64)', 'string', 'JobKey', 'jobKey', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', '', 4, 'input', '', 64, 0);
INSERT INTO `gen_table_column` VALUES (2573, 266, 'cron_expr', 'cron执行表达式', 'varchar(32)', 'string', 'CronExpr', 'cronExpr', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 5, 'input', '', 32, 0);
INSERT INTO `gen_table_column` VALUES (2574, 266, 'invoke_target', '执行目标：服务.方法名', 'varchar(128)', 'string', 'InvokeTarget', 'invokeTarget', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 6, 'input', '', 128, 0);
INSERT INTO `gen_table_column` VALUES (2575, 266, 'status', '状态 0正常 1暂停', 'tinyint', 'int16', 'Status', 'status', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', 'sys_active', 7, 'switch', '0', 0, 0);
INSERT INTO `gen_table_column` VALUES (2576, 266, 'concurrent', '是否并发 0禁止 1允许', 'tinyint', 'int16', 'Concurrent', 'concurrent', 0, 0, 1, 1, 1, 1, 0, 0, '', 'sys_is', 8, 'switch', '0', 0, 0);
INSERT INTO `gen_table_column` VALUES (2577, 266, 'retry_count', '失败重试次数', 'int', 'int', 'RetryCount', 'retryCount', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 9, 'input-number', '0', 0, 0);
INSERT INTO `gen_table_column` VALUES (2578, 266, 'timeout', '执行超时时间(秒)', 'int', 'int', 'Timeout', 'timeout', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 10, 'input-number', '0', 0, 0);
INSERT INTO `gen_table_column` VALUES (2579, 266, 'remark', '备注说明', 'varchar(255)', 'string', 'Remark', 'remark', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 11, 'textarea', '', 255, 0);
INSERT INTO `gen_table_column` VALUES (2580, 266, 'created_by', '创建人ID', 'bigint', 'int64', 'CreatedBy', 'createdBy,string', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 12, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2581, 266, 'created_at', '创建时间', 'timestamp', 'time.Time', 'CreatedAt', 'createdAt', 0, 0, 0, 0, 0, 1, 0, 0, '', '', 13, '', 'CURRENT_TIMESTAMP', 0, 0);
INSERT INTO `gen_table_column` VALUES (2582, 266, 'updated_at', '更新时间', 'timestamp', 'time.Time', 'UpdatedAt', 'updatedAt', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 14, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2583, 266, 'updated_by', '更新人ID', 'bigint', 'int64', 'UpdatedBy', 'updatedBy,string', 0, 0, 0, 0, 0, 0, 0, 0, '', '', 15, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2584, 267, 'id', '日志ID', 'bigint', 'int64', 'Id', 'id,string', 1, 1, 1, 0, 0, 1, 0, 0, '', '', 1, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2585, 267, 'job_id', '任务ID', 'bigint', 'int64', 'JobId', 'jobId,string', 0, 0, 1, 1, 0, 1, 0, 0, 'eq', '', 2, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2586, 267, 'job_name', '任务名称（冗余）', 'varchar(64)', 'string', 'JobName', 'jobName', 0, 0, 0, 1, 0, 1, 0, 0, 'like', '', 3, '', '', 64, 0);
INSERT INTO `gen_table_column` VALUES (2587, 267, 'job_group', '任务分组（冗余）', 'varchar(64)', 'string', 'JobGroup', 'jobGroup', 0, 0, 0, 1, 0, 1, 0, 0, 'eq', '', 4, '', '', 64, 0);
INSERT INTO `gen_table_column` VALUES (2588, 267, 'status', '执行状态 0成功 1失败 2进行中', 'tinyint', 'int16', 'Status', 'status', 0, 0, 0, 1, 0, 1, 0, 0, 'eq', '', 5, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2589, 267, 'message', '执行日志/异常信息', 'text', 'string', 'Message', 'message', 0, 0, 0, 1, 0, 1, 0, 0, '', '', 6, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2590, 267, 'duration', '执行耗时(毫秒)', 'int', 'int', 'Duration', 'duration', 0, 0, 0, 1, 0, 1, 0, 0, '', '', 7, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2591, 267, 'trigger_type', '触发类型：自动/手动', 'varchar(16)', 'string', 'TriggerType', 'triggerType', 0, 0, 0, 1, 0, 1, 0, 0, 'eq', '', 8, '', '', 16, 0);
INSERT INTO `gen_table_column` VALUES (2592, 267, 'start_time', '开始时间', 'datetime', 'time.Time', 'StartTime', 'startTime', 0, 0, 0, 1, 0, 1, 0, 0, '', '', 9, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2593, 267, 'end_time', '结束时间', 'datetime', 'time.Time', 'EndTime', 'endTime', 0, 0, 0, 1, 0, 1, 0, 0, '', '', 10, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2594, 267, 'created_at', '日志生成时间', 'timestamp', 'time.Time', 'CreatedAt', 'createdAt', 0, 0, 0, 1, 0, 1, 0, 0, 'between', '', 11, '', 'CURRENT_TIMESTAMP', 0, 0);
INSERT INTO `gen_table_column` VALUES (2595, 268, 'id', '内部主键', 'bigint', 'int64', 'Id', 'id,string', 1, 1, 1, 0, 1, 1, 0, 0, '', '', 1, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2596, 268, 'user_id', '系统用户ID（关联用户表）', 'bigint', 'int64', 'UserId', 'userId,string', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', '', 2, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2597, 268, 'username', '冗余用户名（用于快速查询/审计）', 'varchar(64)', 'string', 'Username', 'username', 0, 0, 1, 1, 1, 1, 0, 0, 'like', '', 3, '', '', 64, 0);
INSERT INTO `gen_table_column` VALUES (2598, 268, 'credential_id', '凭证ID（二进制存储，WebAuthn原始字节）', 'varbinary(255)', 'string', 'CredentialId', 'credentialId', 0, 0, 1, 1, 1, 1, 0, 0, 'eq', '', 4, '', '', 255, 0);
INSERT INTO `gen_table_column` VALUES (2599, 268, 'public_key_cose', 'COSE格式公钥（二进制）', 'varbinary(512)', 'string', 'PublicKeyCose', 'publicKeyCose', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 5, '', '', 512, 0);
INSERT INTO `gen_table_column` VALUES (2600, 268, 'user_handle', 'WebAuthn用户标识（稳定且不可变）', 'varbinary(64)', 'string', 'UserHandle', 'userHandle', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 6, '', '', 64, 0);
INSERT INTO `gen_table_column` VALUES (2601, 268, 'signature_count', '签名计数器（防克隆攻击）', 'bigint unsigned', 'uint64', 'SignatureCount', 'signatureCount,string', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 7, '', '0', 0, 0);
INSERT INTO `gen_table_column` VALUES (2602, 268, 'aaguid', '认证器设备类型ID', 'char(36)', 'string', 'Aaguid', 'aaguid', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 8, '', '', 36, 0);
INSERT INTO `gen_table_column` VALUES (2603, 268, 'attestation_type', '证明类型（basic/self/none等）', 'varchar(32)', 'string', 'AttestationType', 'attestationType', 0, 0, 0, 1, 1, 1, 0, 0, 'in', '', 9, '', '', 32, 0);
INSERT INTO `gen_table_column` VALUES (2604, 268, 'attachment', '认证器类型（platform/cross-platform）', 'varchar(32)', 'string', 'Attachment', 'attachment', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 10, '', '', 32, 0);
INSERT INTO `gen_table_column` VALUES (2605, 268, 'transports', '认证器传输方式（usb,nfc,ble,internal）', 'varchar(255)', 'string', 'Transports', 'transports', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 11, '', '', 255, 0);
INSERT INTO `gen_table_column` VALUES (2606, 268, 'device_type', '设备类型（single-device / multi-device）', 'varchar(32)', 'string', 'DeviceType', 'deviceType', 0, 0, 0, 1, 1, 1, 0, 0, 'in', '', 12, '', '', 32, 0);
INSERT INTO `gen_table_column` VALUES (2607, 268, 'backup_state', '是否支持云同步（passkey关键属性）', 'tinyint(1)', 'int16', 'BackupState', 'backupState', 0, 0, 0, 1, 1, 1, 0, 0, 'in', '', 13, '', '', 1, 0);
INSERT INTO `gen_table_column` VALUES (2608, 268, 'device_name', '用户自定义设备名称', 'varchar(255)', 'string', 'DeviceName', 'deviceName', 0, 0, 0, 1, 1, 1, 0, 0, 'like', '', 14, '', '', 255, 0);
INSERT INTO `gen_table_column` VALUES (2609, 268, 'is_revoked', '是否已吊销（禁用凭证）', 'tinyint(1)', 'int16', 'IsRevoked', 'isRevoked', 0, 0, 1, 1, 1, 1, 0, 0, '', '', 15, '', '0', 1, 0);
INSERT INTO `gen_table_column` VALUES (2610, 268, 'last_used_at', '最后使用时间', 'timestamp', 'time.Time', 'LastUsedAt', 'lastUsedAt', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 16, '', '', 0, 0);
INSERT INTO `gen_table_column` VALUES (2611, 268, 'created_at', '创建时间', 'timestamp', 'time.Time', 'CreatedAt', 'createdAt', 0, 0, 1, 1, 1, 1, 0, 0, 'between', '', 17, '', 'CURRENT_TIMESTAMP', 0, 0);
INSERT INTO `gen_table_column` VALUES (2612, 268, 'updated_at', '更新时间', 'timestamp', 'time.Time', 'UpdatedAt', 'updatedAt', 0, 0, 0, 1, 1, 1, 0, 0, '', '', 18, '', '', 0, 0);
COMMIT;

-- ----------------------------
-- Table structure for sys_conf
-- ----------------------------
DROP TABLE IF EXISTS `sys_conf`;
CREATE TABLE `sys_conf` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '编号',
  `name` varchar(255) DEFAULT NULL COMMENT '名称',
  `conf_key` varchar(64) DEFAULT NULL COMMENT '配置KEY',
  `conf_value` text COMMENT '配置值',
  `is_sys` tinyint(1) DEFAULT NULL COMMENT '是否系统内置（0否 1是）',
  `scope` tinyint(1) DEFAULT NULL COMMENT '范围（0公开配置，1登录配置，2后台服务配置）',
  `is_secret` tinyint(1) DEFAULT NULL COMMENT '是否加密/脱敏（0否 1是）',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '备注',
  `created_by` bigint DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_by` bigint DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='参数配置表';

-- ----------------------------
-- Records of sys_conf
-- ----------------------------
BEGIN;
INSERT INTO `sys_conf` VALUES (1, '账号初始密码', 'sys.user.initPassword', 'kwpzXGcBZbG5R3X02qZtEw==', 1, 2, 1, '初始化密码 123456', 1, '2022-12-11 16:51:52', 1, '2026-05-18 18:05:19');
INSERT INTO `sys_conf` VALUES (2, '验证码开关', 'sys.captcha.enabled', 'true', 1, 0, 0, '验证码开关（true开启，false关闭）', 1, '2022-12-11 16:51:52', 1, '2025-04-16 09:18:39');
INSERT INTO `sys_conf` VALUES (3, '验证码类型', 'sys.captcha.type', 'SLIDER', 1, 0, 0, 'SLIDER/ROTATE/CONCAT/WORD_IMAGE_CLICK', 1, '2025-08-23 09:34:48', 0, '2025-08-23 09:34:55');
INSERT INTO `sys_conf` VALUES (5, '开启用户注册功能', 'sys.registerEnabled', 'false', 1, 0, 0, '开启用户注册功能（true开启，false关闭）', 1, '2022-12-11 16:51:52', 1, '2025-04-16 09:18:43');
INSERT INTO `sys_conf` VALUES (6, '资源访问前缀地址', 'sys.ossBaseUrl', 'http://127.0.0.1:9000/', 1, 0, 0, '资源访问前缀地址', 1, '2023-09-04 18:44:54', 1, '2025-04-16 09:18:45');
INSERT INTO `sys_conf` VALUES (8, '系统名称', 'sys.systemName', 'Aevons', 1, 0, 0, '系统名称', 1, '2023-09-06 18:36:50', 1, '2026-04-21 16:59:31');
INSERT INTO `sys_conf` VALUES (13, 'IP黑名单', 'sys.blacklist', '#', 1, 0, 0, '多个逗号分隔', 1, '2023-09-26 16:40:40', 1, '2025-04-16 09:18:51');
INSERT INTO `sys_conf` VALUES (15, 'Mail 驱动', 'sys.mail.driver', 'qq', 1, 2, 0, 'qq/google/163/', 1, '2025-08-09 20:50:37', 0, '2025-08-09 20:51:01');
INSERT INTO `sys_conf` VALUES (16, 'Mail 主机', 'sys.mail.host', 'smtp.qq.com', 1, 2, 0, 'QQ Email 主机', 1, '2025-08-09 20:53:09', 0, '2025-08-09 20:53:31');
INSERT INTO `sys_conf` VALUES (17, 'Mail 协议', 'sys.mail.protocol', 'smtp', 1, 2, 0, 'Mail 协议', 1, '2025-08-09 20:55:23', 0, '2025-08-09 20:55:36');
INSERT INTO `sys_conf` VALUES (18, '发件地址', 'sys.mail.email', 'bQ4xAlQQxGUSrTp4KtGtyg==', 1, 2, 1, '发件地址', 1, '2025-08-09 20:56:38', 0, '2026-05-18 18:05:31');
INSERT INTO `sys_conf` VALUES (19, '发件授权码', 'sys.mail.code', '', 1, 2, 0, '发件授权码', 1, '2025-08-09 21:00:11', 0, '2026-04-18 08:19:43');
INSERT INTO `sys_conf` VALUES (21, '百度翻译AppId', 'i18n.baidu.translator.app_id', '', 1, 2, 0, '', 1, '2026-02-02 10:45:45', 0, '2026-04-18 08:20:09');
INSERT INTO `sys_conf` VALUES (22, '百度翻译SecretKey', 'i18n.baidu.translator.secret_key', '', 1, 2, 0, '', 1, '2026-02-02 10:46:31', 0, '2026-04-18 08:20:04');
INSERT INTO `sys_conf` VALUES (23, '百度翻译AppKey', 'i18n.baidu.translator.api_key', '', 1, 2, 0, '', 1, '2026-02-02 13:34:16', 0, '2026-04-18 08:20:12');
INSERT INTO `sys_conf` VALUES (24, '多语言翻译强制更新开关', 'i18n.forced_update', '1', 1, 2, 0, '（1开启，0关闭）, 开启后翻译过的也会更新', 1, '2026-02-02 13:34:16', 0, '2026-02-02 10:47:40');
COMMIT;

-- ----------------------------
-- Table structure for sys_dept
-- ----------------------------
DROP TABLE IF EXISTS `sys_dept`;
CREATE TABLE `sys_dept` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '部门编号',
  `parent_id` bigint DEFAULT NULL COMMENT '父级编号',
  `ancestors` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '祖级',
  `dept_type` tinyint DEFAULT NULL COMMENT '类型（1机构，2部门）',
  `dept_name` varchar(30) DEFAULT NULL COMMENT '部门名称',
  `sort` int DEFAULT NULL COMMENT '顺序',
  `status` tinyint(1) DEFAULT NULL COMMENT '状态（0正常 1停用）',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '备注',
  `created_by` bigint DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_by` bigint DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=120 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='部门表';

-- ----------------------------
-- Records of sys_dept
-- ----------------------------
BEGIN;
INSERT INTO `sys_dept` VALUES (100, 0, '0', 1, '光阴', 1, 0, '', 1, '2022-12-11 16:51:52', 1, '2026-04-24 14:46:40');
INSERT INTO `sys_dept` VALUES (105, 100, '0,100', 2, '市场部', 3, 0, '', 1, '2022-12-11 16:51:52', 1, '2025-04-11 09:10:27');
INSERT INTO `sys_dept` VALUES (110, 100, '0,100', 2, '运营部', 4, 0, '', 1, '2023-11-03 14:11:59', 1, '2026-04-24 14:49:20');
INSERT INTO `sys_dept` VALUES (111, 100, '0,100', 2, '财务部', 5, 0, '', 1, '2023-11-03 14:12:52', 1, '2026-04-24 14:48:21');
INSERT INTO `sys_dept` VALUES (112, 100, '0,100', 2, '研发部', 6, 0, '', 1, '2023-11-03 14:13:35', 1, '2026-04-24 14:48:16');
INSERT INTO `sys_dept` VALUES (114, 100, '0,100', 2, '测试组', 1, 0, '', 1, '2024-01-23 16:54:33', 1, '2026-04-20 14:54:13');
COMMIT;

-- ----------------------------
-- Table structure for sys_dict
-- ----------------------------
DROP TABLE IF EXISTS `sys_dict`;
CREATE TABLE `sys_dict` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '字典主键',
  `dict_type` varchar(32) DEFAULT NULL COMMENT '字典类型',
  `dict_name` varchar(50) DEFAULT NULL COMMENT '字典名称',
  `status` tinyint(1) DEFAULT NULL COMMENT '状态（0正常 1停用）',
  `is_sys` tinyint(1) DEFAULT NULL COMMENT '系统内置（0否 1是）',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '备注',
  `created_by` bigint DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `updated_by` bigint DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=304 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='字典类型表';

-- ----------------------------
-- Records of sys_dict
-- ----------------------------
BEGIN;
INSERT INTO `sys_dict` VALUES (1, 'sys_common_status', '通用状态', 0, 1, '', 1, '2024-12-17 13:16:00', '2025-04-12 22:50:25', NULL);
INSERT INTO `sys_dict` VALUES (2, 'sys_is', '是否', 0, 1, '', 1, '2024-12-18 13:16:00', '2025-04-12 19:21:14', NULL);
INSERT INTO `sys_dict` VALUES (3, 'sys_user_type', '用户类型', 0, 1, NULL, 1, '2024-12-19 13:16:00', NULL, NULL);
INSERT INTO `sys_dict` VALUES (4, 'sys_user_sex', '用户性别', 0, 1, '用户性别列表', 1, '2024-12-20 13:16:00', NULL, NULL);
INSERT INTO `sys_dict` VALUES (5, 'sys_visible', '可见', 0, 1, '可见', 1, '2024-12-21 13:16:00', '2025-05-10 15:01:41', NULL);
INSERT INTO `sys_dict` VALUES (6, 'sys_conf_scope', '系统配置范围', 0, 1, '', 1, '2024-12-22 13:16:00', '2025-08-21 14:31:35', NULL);
INSERT INTO `sys_dict` VALUES (7, 'sys_oper_type', '操作类型', 0, 1, '操作类型列表', 1, '2024-12-23 13:16:00', NULL, NULL);
INSERT INTO `sys_dict` VALUES (8, 'sys_success_status', '成功状态', 0, 1, '操作状态', 1, '2024-12-24 13:16:00', '2026-04-20 15:39:44', NULL);
INSERT INTO `sys_dict` VALUES (9, 'sys_dict_tag_style', '字典标签风格', 0, 1, '字典标签风格', 1, '2024-12-27 13:16:00', '2026-04-18 09:58:16', NULL);
INSERT INTO `sys_dict` VALUES (10, 'sys_menu_type', '菜单类型', 0, 1, '菜单类型', 1, '2024-12-27 13:16:00', '2026-04-18 09:58:16', NULL);
INSERT INTO `sys_dict` VALUES (11, 'sys_client_type', '客户端类型', 0, 0, '客户端类型', 1, '2024-12-27 13:16:00', '2025-04-12 22:54:46', NULL);
INSERT INTO `sys_dict` VALUES (12, 'sys_dept_type', '部门类型', 0, 0, NULL, 1, '2024-12-28 13:16:00', NULL, NULL);
INSERT INTO `sys_dict` VALUES (13, 'sys_data_scope', '角色数据范围', 0, 1, '角色数据范围', 1, '2024-12-28 13:16:00', '2026-04-18 19:17:17', NULL);
INSERT INTO `sys_dict` VALUES (15, 'sys_active', '激活状态', 0, 0, NULL, 1, '2025-12-11 14:04:50', NULL, NULL);
INSERT INTO `sys_dict` VALUES (16, 'sys_condition', '条件类型', 0, 1, '', 0, '2026-04-16 09:19:07', '2026-04-16 12:54:14', 0);
INSERT INTO `sys_dict` VALUES (17, 'sys_authorized_grant_type', '授权类型', 0, 1, '', 0, '2026-04-16 16:36:58', '2026-04-16 16:37:02', 0);
INSERT INTO `sys_dict` VALUES (18, 'sys_authorized_scope', '授权范围', 0, 1, '', 0, '2026-04-16 17:41:54', '2026-04-16 17:41:54', 0);
INSERT INTO `sys_dict` VALUES (20, 'sys_device_type', '设备类型', 0, 0, '', 1, '2025-11-24 13:58:31', '2026-04-17 22:23:18', NULL);
INSERT INTO `sys_dict` VALUES (21, 'sys_locale_namespace', '语言资源命名空间', 0, 0, '语言资源命名空间', 1, '2026-02-03 11:51:47', '2026-04-17 22:23:08', NULL);
INSERT INTO `sys_dict` VALUES (30, 'sys_notice_type', '通知类型', 0, 0, '通知类型列表', 1, '2024-12-29 13:16:00', NULL, NULL);
INSERT INTO `sys_dict` VALUES (32, 'sys_template_level', '模板级别', 0, 0, '', 1, '2024-12-26 13:16:00', '2024-02-21 09:24:03', NULL);
INSERT INTO `sys_dict` VALUES (33, 'sys_template_type', '模板类型', 0, 0, '模板类型', 1, '2024-12-25 13:16:00', NULL, NULL);
INSERT INTO `sys_dict` VALUES (200, 'xcode_component', '生成代码-UI组件', 0, 0, '', 0, '2026-04-16 15:44:25', '2026-04-17 22:21:36', 0);
INSERT INTO `sys_dict` VALUES (201, 'xcode_table_type', '表类型', 0, 0, '', 1, '2025-01-03 13:16:00', '2025-04-17 10:49:21', NULL);
INSERT INTO `sys_dict` VALUES (202, 'xcode_sort_type', '排序类型', 0, 0, '', 1, '2025-01-04 13:16:00', '2025-04-12 20:45:29', NULL);
INSERT INTO `sys_dict` VALUES (203, 'xcode_data_type', '数据类型', 0, 0, '', 1, '2025-01-05 13:16:00', '2025-04-12 20:42:00', NULL);
INSERT INTO `sys_dict` VALUES (204, 'xcode_condition', '条件', 0, 0, '', 1, '2025-01-06 13:16:00', '2025-04-12 20:46:23', NULL);
INSERT INTO `sys_dict` VALUES (205, 'xcode_control', '控件', 0, 0, '', 1, '2025-01-07 13:16:00', '2025-04-18 09:22:16', NULL);
INSERT INTO `sys_dict` VALUES (206, 'xcode_meta_type', '字段元数据类型', 0, 0, '', 1, '2025-01-08 13:16:00', '2025-04-18 09:29:36', NULL);
INSERT INTO `sys_dict` VALUES (257, 'xcode_page_type', '页面类型', 0, 0, NULL, 1, '2025-09-18 14:11:51', NULL, NULL);
INSERT INTO `sys_dict` VALUES (259, 'job_status', '定时任务状态', 0, 0, NULL, 1, '2025-11-18 14:54:06', NULL, NULL);
INSERT INTO `sys_dict` VALUES (260, 'job_misfire_policy', '定时任务策略', 0, 0, NULL, 1, '2025-11-18 16:19:05', NULL, NULL);
INSERT INTO `sys_dict` VALUES (261, 'job_concurrent', '定时任务并发', 0, 0, '', 1, '2025-11-18 16:25:28', '2025-11-19 12:51:29', NULL);
INSERT INTO `sys_dict` VALUES (300, 'flow_category', '流程分类', 0, 0, '', 1, '2024-12-31 13:16:00', '2025-04-12 23:00:27', NULL);
INSERT INTO `sys_dict` VALUES (301, 'flow_approve_status', '流程审批状态', 0, 0, '', 1, '2025-01-01 13:16:00', '2025-04-12 22:58:27', NULL);
INSERT INTO `sys_dict` VALUES (302, 'flow_deploy_status', '流程部署状态', 0, 0, '流程部署状态', 1, '2025-01-02 13:16:00', '2025-04-12 22:57:19', NULL);
INSERT INTO `sys_dict` VALUES (303, 'flow_active_status', '流程激活状态', 0, 0, NULL, 1, '2025-01-09 13:16:00', NULL, NULL);
COMMIT;

-- ----------------------------
-- Table structure for sys_dict_data
-- ----------------------------
DROP TABLE IF EXISTS `sys_dict_data`;
CREATE TABLE `sys_dict_data` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '字典编号',
  `dict_type` varchar(32) DEFAULT NULL COMMENT '字典类型',
  `dict_value` varchar(32) DEFAULT NULL COMMENT '字典键值',
  `status` tinyint(1) DEFAULT NULL COMMENT '状态（0正常 1停用）',
  `sort` int DEFAULT NULL COMMENT '顺序',
  `tag_type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '标签风格',
  `tag_class` varchar(10) DEFAULT NULL COMMENT '样式类名',
  `created_by` bigint DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `updated_by` bigint DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_dict_type_value` (`dict_type`,`dict_value`)
) ENGINE=InnoDB AUTO_INCREMENT=242 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='字典数据表';

-- ----------------------------
-- Records of sys_dict_data
-- ----------------------------
BEGIN;
INSERT INTO `sys_dict_data` VALUES (1, 'sys_common_status', '1', 0, 1, 'default', '', NULL, '2022-12-18 22:01:08', '2026-04-20 16:02:08', NULL);
INSERT INTO `sys_dict_data` VALUES (2, 'sys_common_status', '0', 0, 0, 'success', '', NULL, '2022-12-18 22:01:21', '2026-04-25 17:14:09', NULL);
INSERT INTO `sys_dict_data` VALUES (3, 'sys_is', '1', 0, 1, 'info', '', NULL, '2023-09-06 17:52:14', '2026-04-20 16:00:54', NULL);
INSERT INTO `sys_dict_data` VALUES (4, 'sys_is', '0', 0, 2, 'warning', '', NULL, '2023-09-06 17:52:25', '2026-04-20 16:00:46', NULL);
INSERT INTO `sys_dict_data` VALUES (5, 'sys_user_type', '0', 0, 0, 'info', '', NULL, '2024-11-03 10:15:56', '2026-04-18 12:05:50', NULL);
INSERT INTO `sys_dict_data` VALUES (6, 'sys_user_type', '1', 0, 0, 'info', '', NULL, '2024-11-03 10:16:28', '2026-04-18 12:04:39', NULL);
INSERT INTO `sys_dict_data` VALUES (7, 'sys_user_sex', '1', 0, 1, '', '', NULL, '2022-12-11 16:51:52', '2026-04-18 12:06:16', NULL);
INSERT INTO `sys_dict_data` VALUES (8, 'sys_user_sex', '2', 0, 2, '', '', NULL, '2022-12-11 16:51:52', '2026-04-18 12:06:35', NULL);
INSERT INTO `sys_dict_data` VALUES (9, 'sys_user_sex', '0', 0, 3, '', '', NULL, '2022-12-11 16:51:52', '2026-04-18 12:06:50', NULL);
INSERT INTO `sys_dict_data` VALUES (10, 'sys_visible', '1', 0, 1, 'info', '', NULL, '2022-12-11 16:51:52', '2026-04-18 12:07:16', NULL);
INSERT INTO `sys_dict_data` VALUES (11, 'sys_visible', '0', 0, 2, 'danger', '', NULL, '2022-12-11 16:51:52', '2026-04-18 12:07:22', NULL);
INSERT INTO `sys_dict_data` VALUES (12, 'sys_conf_scope', '0', 0, 1, 'default', '', NULL, '2025-08-09 19:19:56', '2026-04-25 17:29:28', NULL);
INSERT INTO `sys_dict_data` VALUES (13, 'sys_conf_scope', '1', 0, 2, 'default', '', NULL, '2025-08-09 19:20:11', '2026-04-25 17:29:21', NULL);
INSERT INTO `sys_dict_data` VALUES (14, 'sys_conf_scope', '2', 0, 3, 'default', '', NULL, '2025-08-09 19:21:44', '2026-04-25 17:29:13', NULL);
INSERT INTO `sys_dict_data` VALUES (15, 'sys_oper_type', 'OTHER', 0, 99, 'info', '', NULL, '2022-12-11 16:51:52', '2026-04-18 12:11:49', NULL);
INSERT INTO `sys_dict_data` VALUES (16, 'sys_oper_type', 'INSERT', 0, 1, 'info', '', NULL, '2022-12-11 16:51:52', '2026-04-18 12:08:53', NULL);
INSERT INTO `sys_dict_data` VALUES (17, 'sys_oper_type', 'UPDATE', 0, 2, 'success', '', NULL, '2022-12-11 16:51:52', '2026-04-18 12:09:08', NULL);
INSERT INTO `sys_dict_data` VALUES (18, 'sys_oper_type', 'DELETE', 0, 3, 'danger', '', NULL, '2022-12-11 16:51:52', '2026-04-18 12:09:15', NULL);
INSERT INTO `sys_dict_data` VALUES (19, 'sys_oper_type', 'AUTH', 0, 4, 'default', '', NULL, '2022-12-11 16:51:52', '2026-04-18 12:09:34', NULL);
INSERT INTO `sys_dict_data` VALUES (20, 'sys_oper_type', 'EXPORT', 0, 5, 'default', '', NULL, '2022-12-11 16:51:52', '2026-04-18 12:09:47', NULL);
INSERT INTO `sys_dict_data` VALUES (21, 'sys_oper_type', 'IMPORT', 0, 6, 'default', '', NULL, '2022-12-11 16:51:52', '2026-04-18 12:09:55', NULL);
INSERT INTO `sys_dict_data` VALUES (22, 'sys_oper_type', 'KICKED', 0, 7, 'danger', '', NULL, '2022-12-11 16:51:52', '2026-04-18 12:10:42', NULL);
INSERT INTO `sys_dict_data` VALUES (23, 'sys_oper_type', 'CLEAN', 0, 9, 'danger', '', NULL, '2022-12-11 16:51:52', '2026-04-18 12:11:07', NULL);
INSERT INTO `sys_dict_data` VALUES (24, 'sys_oper_type', 'SETUP', 0, 11, 'default', '', NULL, '2025-06-02 12:03:06', '2026-04-18 12:11:20', NULL);
INSERT INTO `sys_dict_data` VALUES (25, 'sys_oper_type', 'SAVE', 0, 10, 'success', '', NULL, '2025-09-04 08:45:58', '2026-04-18 12:11:13', NULL);
INSERT INTO `sys_dict_data` VALUES (26, 'sys_oper_type', 'RELEASE', 0, 12, 'warning', '', NULL, '2025-09-04 08:46:34', '2026-04-18 12:11:29', NULL);
INSERT INTO `sys_dict_data` VALUES (27, 'sys_oper_type', 'COPY', 0, 13, 'info', '', NULL, '2025-09-04 08:47:30', '2026-04-18 12:11:35', NULL);
INSERT INTO `sys_dict_data` VALUES (29, 'sys_success_status', '0', 0, 2, 'warning', '', NULL, '2022-12-11 16:51:52', '2026-04-20 15:57:31', NULL);
INSERT INTO `sys_dict_data` VALUES (30, 'sys_template_type', '2', 0, 2, 'default', '', NULL, '2024-02-21 09:13:47', '2026-04-18 14:11:33', NULL);
INSERT INTO `sys_dict_data` VALUES (31, 'sys_template_type', '3', 0, 3, 'default', '', NULL, '2024-02-21 09:14:39', '2026-04-18 14:11:25', NULL);
INSERT INTO `sys_dict_data` VALUES (32, 'sys_template_type', '4', 0, 4, 'default', '', NULL, '2024-02-21 09:15:12', '2026-04-18 14:11:17', NULL);
INSERT INTO `sys_dict_data` VALUES (33, 'sys_template_type', '1', 0, 1, 'default', '', NULL, '2024-02-21 09:15:50', '2026-04-18 14:11:40', NULL);
INSERT INTO `sys_dict_data` VALUES (34, 'sys_template_level', '1', 0, 1, 'info', '', NULL, '2024-02-21 09:24:55', '2026-04-18 14:09:49', NULL);
INSERT INTO `sys_dict_data` VALUES (35, 'sys_template_level', '2', 0, 2, 'warning', '', NULL, '2024-02-21 09:26:32', '2026-04-18 14:09:59', NULL);
INSERT INTO `sys_dict_data` VALUES (36, 'sys_template_level', '3', 0, 3, 'danger', '', NULL, '2024-02-21 09:27:44', '2026-04-18 14:10:06', NULL);
INSERT INTO `sys_dict_data` VALUES (37, 'sys_client_type', 'Windows', 0, 3, 'default', '', NULL, '2023-09-25 17:31:03', '2026-04-18 14:05:06', NULL);
INSERT INTO `sys_dict_data` VALUES (38, 'sys_client_type', 'Android', 0, 1, 'default', '', NULL, '2023-09-25 17:31:39', '2026-04-18 14:05:14', NULL);
INSERT INTO `sys_dict_data` VALUES (39, 'sys_client_type', 'IOS', 0, 2, 'default', '', NULL, '2023-09-25 17:31:56', '2026-04-18 14:05:11', NULL);
INSERT INTO `sys_dict_data` VALUES (40, 'sys_dept_type', '1', 0, 1, 'info', '', NULL, '2024-11-02 20:05:05', '2026-04-18 14:05:50', NULL);
INSERT INTO `sys_dict_data` VALUES (41, 'sys_dept_type', '2', 0, 2, 'warning', '', NULL, '2024-11-02 20:05:13', '2026-04-18 14:06:02', NULL);
INSERT INTO `sys_dict_data` VALUES (42, 'sys_notice_type', '1', 0, 1, 'warning', '', NULL, '2022-12-11 16:51:52', '2026-04-18 14:08:14', NULL);
INSERT INTO `sys_dict_data` VALUES (43, 'sys_notice_type', '2', 0, 2, 'success', '', NULL, '2022-12-11 16:51:52', '2026-04-18 14:09:01', NULL);
INSERT INTO `sys_dict_data` VALUES (46, 'flow_category', '0', 0, 1, 'success', '', NULL, '2023-11-27 10:47:56', '2026-04-29 12:39:09', NULL);
INSERT INTO `sys_dict_data` VALUES (47, 'flow_category', '1', 0, 2, 'warning', '', NULL, '2023-12-08 14:29:53', '2026-04-29 12:39:20', NULL);
INSERT INTO `sys_dict_data` VALUES (48, 'flow_approve_status', '0', 0, 0, 'info', '', NULL, '2024-01-23 14:11:07', '2026-04-29 12:38:07', NULL);
INSERT INTO `sys_dict_data` VALUES (49, 'flow_approve_status', '1', 0, 1, 'info', '', NULL, '2024-01-23 14:11:20', '2026-04-29 12:38:21', NULL);
INSERT INTO `sys_dict_data` VALUES (50, 'flow_approve_status', '2', 0, 2, 'warning', '', NULL, '2024-01-23 14:11:37', '2026-04-29 12:38:33', NULL);
INSERT INTO `sys_dict_data` VALUES (51, 'flow_approve_status', '4', 0, 4, 'danger', '', NULL, '2024-01-23 14:12:03', '2026-04-29 12:38:44', NULL);
INSERT INTO `sys_dict_data` VALUES (52, 'flow_approve_status', '6', 0, 6, 'success', '', NULL, '2024-01-23 14:12:22', '2026-04-29 12:38:54', NULL);
INSERT INTO `sys_dict_data` VALUES (53, 'flow_deploy_status', '0', 0, 1, 'success', '', NULL, '2024-11-04 16:53:37', '2026-04-29 12:36:53', NULL);
INSERT INTO `sys_dict_data` VALUES (54, 'flow_deploy_status', '4', 0, 2, 'warning', '', NULL, '2024-11-04 16:53:47', '2026-04-29 12:37:20', NULL);
INSERT INTO `sys_dict_data` VALUES (55, 'xcode_table_type', 'COMMON', 0, 1, 'info', NULL, NULL, '2024-06-17 10:50:09', '2025-04-17 10:49:21', NULL);
INSERT INTO `sys_dict_data` VALUES (56, 'xcode_table_type', 'FLOW', 0, 3, 'warning', NULL, NULL, '2024-06-17 10:50:51', '2025-04-17 10:49:21', NULL);
INSERT INTO `sys_dict_data` VALUES (57, 'xcode_table_type', 'TREE', 0, 2, 'success', NULL, NULL, '2024-06-17 10:51:40', '2025-04-17 10:49:21', NULL);
INSERT INTO `sys_dict_data` VALUES (58, 'xcode_table_type', 'FLOW_NODE', 1, 4, 'danger', NULL, NULL, '2024-10-31 17:06:56', '2025-04-18 10:22:55', NULL);
INSERT INTO `sys_dict_data` VALUES (60, 'xcode_table_type', 'CUSTOM', 0, 6, 'info', NULL, NULL, '2025-04-18 10:22:11', '2025-04-18 10:23:39', NULL);
INSERT INTO `sys_dict_data` VALUES (61, 'xcode_sort_type', '1', 0, 0, 'default', NULL, NULL, '2024-07-30 10:17:06', '2025-04-12 20:45:29', NULL);
INSERT INTO `sys_dict_data` VALUES (62, 'xcode_sort_type', '2', 0, 1, 'default', NULL, NULL, '2024-07-30 10:17:19', '2025-04-12 20:45:29', NULL);
INSERT INTO `sys_dict_data` VALUES (63, 'xcode_data_type', 'bigint', 0, 3, 'default', NULL, NULL, '2024-11-30 16:41:58', '2025-04-12 20:42:00', NULL);
INSERT INTO `sys_dict_data` VALUES (64, 'xcode_data_type', 'bit', 0, 10, 'default', NULL, NULL, '2024-11-30 16:42:11', '2025-04-12 20:42:00', NULL);
INSERT INTO `sys_dict_data` VALUES (65, 'xcode_data_type', 'char', 0, 7, 'default', NULL, NULL, '2024-11-30 16:42:24', '2025-04-12 20:42:00', NULL);
INSERT INTO `sys_dict_data` VALUES (66, 'xcode_data_type', 'date', 0, 4, 'default', NULL, NULL, '2024-11-30 16:42:42', '2025-04-12 20:42:00', NULL);
INSERT INTO `sys_dict_data` VALUES (67, 'xcode_data_type', 'datetime', 0, 5, 'default', NULL, NULL, '2024-11-30 16:43:00', '2025-04-12 20:42:00', NULL);
INSERT INTO `sys_dict_data` VALUES (68, 'xcode_data_type', 'decimal', 0, 6, 'default', NULL, NULL, '2024-11-30 16:43:16', '2025-04-12 20:42:00', NULL);
INSERT INTO `sys_dict_data` VALUES (69, 'xcode_data_type', 'int', 0, 2, 'default', NULL, NULL, '2024-11-30 16:44:52', '2025-04-12 20:42:00', NULL);
INSERT INTO `sys_dict_data` VALUES (70, 'xcode_data_type', 'text', 0, 8, 'default', NULL, NULL, '2024-11-30 16:45:21', '2025-04-12 20:42:00', NULL);
INSERT INTO `sys_dict_data` VALUES (71, 'xcode_data_type', 'tinyint', 0, 9, 'default', NULL, NULL, '2024-11-30 16:45:36', '2025-04-12 20:42:00', NULL);
INSERT INTO `sys_dict_data` VALUES (72, 'xcode_data_type', 'varchar', 0, 1, 'default', NULL, NULL, '2024-11-30 16:45:57', '2025-04-12 20:42:00', NULL);
INSERT INTO `sys_dict_data` VALUES (73, 'xcode_condition', 'eq', 0, 1, 'default', NULL, NULL, '2024-11-30 17:06:32', '2025-04-17 18:59:24', NULL);
INSERT INTO `sys_dict_data` VALUES (74, 'xcode_condition', 'neq', 0, 2, 'default', NULL, NULL, '2024-11-30 17:06:47', '2025-04-17 18:59:13', NULL);
INSERT INTO `sys_dict_data` VALUES (75, 'xcode_condition', 'gt', 0, 3, 'default', NULL, NULL, '2024-11-30 17:07:01', '2025-04-17 18:58:31', NULL);
INSERT INTO `sys_dict_data` VALUES (76, 'xcode_condition', 'ge', 0, 4, 'default', NULL, NULL, '2024-11-30 17:07:12', '2025-04-17 18:58:35', NULL);
INSERT INTO `sys_dict_data` VALUES (77, 'xcode_condition', 'lt', 0, 5, 'default', NULL, NULL, '2024-11-30 17:07:26', '2025-04-17 18:58:40', NULL);
INSERT INTO `sys_dict_data` VALUES (78, 'xcode_condition', 'le', 0, 6, 'default', NULL, NULL, '2024-11-30 17:07:38', '2025-04-17 18:58:45', NULL);
INSERT INTO `sys_dict_data` VALUES (79, 'xcode_condition', 'like', 0, 7, 'default', NULL, NULL, '2024-11-30 17:07:53', '2025-04-17 18:59:01', NULL);
INSERT INTO `sys_dict_data` VALUES (80, 'xcode_condition', 'left_like', 0, 8, 'default', NULL, NULL, '2024-11-30 17:08:21', '2025-04-17 19:00:18', NULL);
INSERT INTO `sys_dict_data` VALUES (81, 'xcode_condition', 'right_like', 0, 9, 'default', NULL, NULL, '2024-11-30 17:10:14', '2025-04-17 19:00:37', NULL);
INSERT INTO `sys_dict_data` VALUES (82, 'xcode_condition', 'be', 0, 10, 'default', NULL, NULL, '2024-11-30 17:10:34', '2025-04-17 19:01:36', NULL);
INSERT INTO `sys_dict_data` VALUES (83, 'xcode_condition', 'not_null', 0, 11, 'default', NULL, NULL, '2024-11-30 17:11:16', '2025-04-17 18:57:26', NULL);
INSERT INTO `sys_dict_data` VALUES (84, 'xcode_condition', 'is_null', 0, 12, 'default', NULL, NULL, '2024-11-30 17:12:31', '2025-04-17 18:57:34', NULL);
INSERT INTO `sys_dict_data` VALUES (85, 'xcode_condition', 'not_empty', 0, 13, 'default', NULL, NULL, '2024-11-30 17:17:22', '2025-04-17 18:58:11', NULL);
INSERT INTO `sys_dict_data` VALUES (86, 'xcode_condition', 'is_empty', 0, 14, 'default', NULL, NULL, '2024-11-30 17:17:38', '2025-04-17 18:58:04', NULL);
INSERT INTO `sys_dict_data` VALUES (87, 'xcode_condition', 'in_list', 0, 15, 'default', NULL, NULL, '2024-11-30 17:18:34', '2025-04-17 19:01:00', NULL);
INSERT INTO `sys_dict_data` VALUES (88, 'xcode_condition', 'not_list', 0, 16, 'default', NULL, NULL, '2024-11-30 17:19:25', '2025-04-17 19:01:09', NULL);
INSERT INTO `sys_dict_data` VALUES (89, 'xcode_control', 'text', 0, 1, 'default', NULL, NULL, '2025-04-17 13:28:06', '2025-04-18 09:22:16', NULL);
INSERT INTO `sys_dict_data` VALUES (90, 'xcode_control', 'number', 0, 2, 'default', NULL, NULL, '2025-04-17 13:28:26', '2025-04-18 09:22:16', NULL);
INSERT INTO `sys_dict_data` VALUES (91, 'xcode_control', 'textarea', 0, 3, 'default', NULL, NULL, '2025-04-17 13:28:44', '2025-04-18 09:22:16', NULL);
INSERT INTO `sys_dict_data` VALUES (92, 'xcode_control', 'password', 0, 23, 'default', NULL, NULL, '2025-04-17 18:49:18', '2025-05-10 11:42:23', NULL);
INSERT INTO `sys_dict_data` VALUES (93, 'xcode_control', 'color', 0, 34, 'default', NULL, NULL, '2025-04-17 18:49:29', '2025-05-10 11:41:00', NULL);
INSERT INTO `sys_dict_data` VALUES (94, 'xcode_control', 'icon', 0, 33, 'default', NULL, NULL, '2025-04-17 18:49:54', '2025-05-10 11:40:55', NULL);
INSERT INTO `sys_dict_data` VALUES (95, 'xcode_control', 'editor', 0, 30, 'default', NULL, NULL, '2025-04-17 18:50:10', '2025-05-10 11:40:34', NULL);
INSERT INTO `sys_dict_data` VALUES (96, 'xcode_control', 'date', 0, 5, 'default', NULL, NULL, '2025-04-17 18:50:22', '2025-04-24 11:43:18', NULL);
INSERT INTO `sys_dict_data` VALUES (97, 'xcode_control', 'datetime', 0, 6, 'default', NULL, NULL, '2025-04-17 18:50:36', '2025-04-24 11:43:22', NULL);
INSERT INTO `sys_dict_data` VALUES (98, 'xcode_control', 'select', 0, 9, 'default', NULL, NULL, '2025-04-17 18:51:15', '2025-05-04 15:33:20', NULL);
INSERT INTO `sys_dict_data` VALUES (99, 'xcode_control', 'select_multiple', 1, 20, 'default', NULL, NULL, '2025-04-17 18:51:31', '2025-05-10 11:42:08', NULL);
INSERT INTO `sys_dict_data` VALUES (100, 'xcode_control', 'radio', 0, 7, 'default', NULL, NULL, '2025-04-17 18:51:50', '2025-04-18 09:22:16', NULL);
INSERT INTO `sys_dict_data` VALUES (101, 'xcode_control', 'checkbox', 0, 8, 'default', NULL, NULL, '2025-04-17 18:52:02', '2025-04-18 09:22:16', NULL);
INSERT INTO `sys_dict_data` VALUES (102, 'xcode_control', 'file', 0, 21, 'default', NULL, NULL, '2025-04-17 18:52:43', '2025-05-10 11:42:15', NULL);
INSERT INTO `sys_dict_data` VALUES (103, 'xcode_control', 'datetime_slot', 0, 32, 'default', NULL, NULL, '2025-04-24 11:42:33', '2025-05-10 11:40:41', NULL);
INSERT INTO `sys_dict_data` VALUES (104, 'xcode_control', 'date_slot', 0, 31, 'default', NULL, NULL, '2025-04-24 11:42:56', '2025-05-10 11:40:38', NULL);
INSERT INTO `sys_dict_data` VALUES (105, 'xcode_control', 'bigfile', 0, 22, 'default', NULL, NULL, '2025-05-04 15:26:56', '2025-05-10 11:42:18', NULL);
INSERT INTO `sys_dict_data` VALUES (106, 'xcode_control', 'radio_button', 0, 11, 'default', NULL, NULL, '2025-05-09 11:16:08', '2025-05-10 11:42:56', NULL);
INSERT INTO `sys_dict_data` VALUES (107, 'xcode_meta_type', 'expression', 0, 3, 'default', NULL, NULL, '2025-04-18 09:24:22', '2025-04-18 12:18:41', NULL);
INSERT INTO `sys_dict_data` VALUES (108, 'xcode_meta_type', 'enum', 0, 1, 'default', NULL, NULL, '2025-04-18 09:24:40', '2025-04-18 12:13:21', NULL);
INSERT INTO `sys_dict_data` VALUES (109, 'xcode_meta_type', 'sort', 0, 2, 'default', NULL, NULL, '2025-04-18 09:25:31', '2025-04-18 12:18:38', NULL);
INSERT INTO `sys_dict_data` VALUES (110, 'xcode_meta_type', 'one_one', 0, 4, 'default', NULL, NULL, '2025-04-18 12:11:27', '2025-04-18 12:18:48', NULL);
INSERT INTO `sys_dict_data` VALUES (111, 'xcode_meta_type', 'one_many', 0, 5, 'default', NULL, NULL, '2025-04-18 12:12:03', '2025-04-18 12:18:58', NULL);
INSERT INTO `sys_dict_data` VALUES (112, 'xcode_meta_type', 'many_one', 0, 6, 'default', NULL, NULL, '2025-04-18 12:12:24', '2025-04-18 12:19:06', NULL);
INSERT INTO `sys_dict_data` VALUES (113, 'xcode_meta_type', 'many_many', 0, 7, 'default', NULL, NULL, '2025-04-18 12:12:52', '2025-04-18 12:19:10', NULL);
INSERT INTO `sys_dict_data` VALUES (128, 'flow_active_status', '1', 0, 1, 'success', '', NULL, '2025-09-18 12:00:55', '2026-04-29 12:32:40', NULL);
INSERT INTO `sys_dict_data` VALUES (129, 'flow_active_status', '2', 0, 2, 'info', NULL, NULL, '2025-09-18 12:01:32', '2026-04-18 08:07:56', NULL);
INSERT INTO `sys_dict_data` VALUES (130, 'xcode_page_type', '1', 0, 1, 'default', NULL, NULL, '2025-09-18 14:12:33', '2025-09-18 14:13:48', NULL);
INSERT INTO `sys_dict_data` VALUES (131, 'xcode_page_type', '2', 0, 2, 'success', NULL, NULL, '2025-09-18 14:13:45', '2025-09-18 14:32:29', NULL);
INSERT INTO `sys_dict_data` VALUES (132, 'xcode_page_type', '3', 0, 3, 'info', NULL, NULL, '2025-09-18 14:14:10', '2025-09-18 14:32:35', NULL);
INSERT INTO `sys_dict_data` VALUES (133, 'xcode_page_type', '4', 0, 4, 'warning', NULL, NULL, '2025-09-18 14:15:03', '2025-09-18 14:32:48', NULL);
INSERT INTO `sys_dict_data` VALUES (134, 'xcode_page_type', '5', 0, 5, 'danger', NULL, NULL, '2025-09-18 14:15:16', '2025-09-18 14:32:54', NULL);
INSERT INTO `sys_dict_data` VALUES (138, 'xcode_control', 'treeselect', 0, 10, 'default', NULL, NULL, '2025-11-12 16:47:20', '2025-11-12 16:47:53', NULL);
INSERT INTO `sys_dict_data` VALUES (139, 'job_status', '0', 0, 1, 'danger', NULL, NULL, '2025-11-18 14:57:36', NULL, NULL);
INSERT INTO `sys_dict_data` VALUES (140, 'job_status', '1', 0, 2, 'info', NULL, NULL, '2025-11-18 14:59:40', NULL, NULL);
INSERT INTO `sys_dict_data` VALUES (141, 'job_misfire_policy', '1', 0, 1, 'info', NULL, NULL, '2025-11-18 16:19:38', NULL, NULL);
INSERT INTO `sys_dict_data` VALUES (142, 'job_misfire_policy', '2', 0, 2, 'default', NULL, NULL, '2025-11-18 16:19:57', NULL, NULL);
INSERT INTO `sys_dict_data` VALUES (143, 'job_misfire_policy', '3', 0, 3, 'warning', NULL, NULL, '2025-11-18 16:20:23', NULL, NULL);
INSERT INTO `sys_dict_data` VALUES (144, 'job_concurrent', '0', 0, 1, 'default', NULL, NULL, '2025-11-18 16:26:07', '2025-11-19 12:51:29', NULL);
INSERT INTO `sys_dict_data` VALUES (145, 'job_concurrent', '1', 0, 2, 'warning', NULL, NULL, '2025-11-18 16:26:39', '2025-11-19 12:51:29', NULL);
INSERT INTO `sys_dict_data` VALUES (148, 'sys_device_type', 'COMPUTER', 0, 1, '', '', NULL, '2025-11-24 13:58:49', '2026-04-18 14:07:18', NULL);
INSERT INTO `sys_dict_data` VALUES (149, 'sys_device_type', 'MOBILE', 0, 2, '', '', NULL, '2025-11-24 13:59:00', '2026-04-18 14:07:26', NULL);
INSERT INTO `sys_dict_data` VALUES (150, 'sys_device_type', 'TABLET', 0, 3, '', '', NULL, '2025-11-24 13:59:14', '2026-04-18 14:07:34', NULL);
INSERT INTO `sys_dict_data` VALUES (151, 'sys_device_type', 'UNKNOWN', 0, 4, '', '', NULL, '2025-11-24 13:59:25', '2026-04-18 14:07:46', NULL);
INSERT INTO `sys_dict_data` VALUES (152, 'sys_active', '0', 0, 2, 'warning', '', NULL, '2025-12-11 14:05:27', '2026-04-18 14:03:40', NULL);
INSERT INTO `sys_dict_data` VALUES (153, 'sys_active', '1', 0, 1, 'success', '', NULL, '2025-12-11 14:05:45', '2026-04-18 14:03:22', NULL);
INSERT INTO `sys_dict_data` VALUES (154, 'sys_locale_namespace', 'default', 0, 1, 'default', '', NULL, '2026-02-03 11:52:25', '2026-04-18 14:06:56', NULL);
INSERT INTO `sys_dict_data` VALUES (174, 'sys_condition', 'eq', 0, 1, 'default', '', 0, '2026-04-16 09:19:46', '2026-04-16 09:19:46', 0);
INSERT INTO `sys_dict_data` VALUES (175, 'sys_condition', 'ne', 0, 2, '', '', 0, '2026-04-16 09:20:14', '2026-04-16 09:20:23', 0);
INSERT INTO `sys_dict_data` VALUES (176, 'sys_condition', 'gt', 0, 3, '', '', 0, '2026-04-16 09:20:41', '2026-04-16 09:20:41', 0);
INSERT INTO `sys_dict_data` VALUES (177, 'sys_condition', 'gte', 0, 4, '', '', 0, '2026-04-16 09:20:57', '2026-04-16 09:20:57', 0);
INSERT INTO `sys_dict_data` VALUES (178, 'sys_condition', 'lt', 0, 5, '', '', 0, '2026-04-16 09:21:11', '2026-04-16 09:21:11', 0);
INSERT INTO `sys_dict_data` VALUES (179, 'sys_condition', 'lte', 0, 6, 'default', '', 0, '2026-04-16 09:21:28', '2026-04-16 09:21:28', 0);
INSERT INTO `sys_dict_data` VALUES (180, 'sys_condition', 'like', 0, 7, 'default', '', 0, '2026-04-16 09:21:52', '2026-04-16 09:21:52', 0);
INSERT INTO `sys_dict_data` VALUES (181, 'sys_condition', 'like_l', 0, 8, '', '', 0, '2026-04-16 09:22:20', '2026-04-16 09:22:20', 0);
INSERT INTO `sys_dict_data` VALUES (182, 'sys_condition', 'like_r', 0, 9, '', '', 0, '2026-04-16 09:23:01', '2026-04-16 09:23:01', 0);
INSERT INTO `sys_dict_data` VALUES (183, 'sys_condition', 'in', 0, 10, '', '', 0, '2026-04-16 09:23:17', '2026-04-16 09:24:50', 0);
INSERT INTO `sys_dict_data` VALUES (184, 'sys_condition', 'not_in', 0, 11, '', '', 0, '2026-04-16 09:23:36', '2026-04-16 09:23:36', 0);
INSERT INTO `sys_dict_data` VALUES (185, 'sys_condition', 'between', 0, 12, '', '', 0, '2026-04-16 09:23:54', '2026-04-16 09:23:54', 0);
INSERT INTO `sys_dict_data` VALUES (186, 'sys_condition', 'is_null', 0, 13, '', '', 0, '2026-04-16 09:24:13', '2026-04-16 09:24:13', 0);
INSERT INTO `sys_dict_data` VALUES (187, 'sys_condition', 'not_null', 0, 14, '', '', 0, '2026-04-16 09:24:30', '2026-04-16 09:24:30', 0);
INSERT INTO `sys_dict_data` VALUES (197, 'xcode_component', 'input', 0, 1, 'default', '', 0, '2026-04-16 16:01:08', '2026-04-18 14:12:00', 0);
INSERT INTO `sys_dict_data` VALUES (198, 'xcode_component', 'input-number', 0, 2, 'default', '', 0, '2026-04-16 16:01:28', '2026-04-18 14:12:28', 0);
INSERT INTO `sys_dict_data` VALUES (199, 'xcode_component', 'select', 0, 3, 'default', '', 0, '2026-04-16 16:01:41', '2026-04-18 14:12:19', 0);
INSERT INTO `sys_dict_data` VALUES (200, 'xcode_component', 'checkbox', 0, 4, 'default', '', 0, '2026-04-16 16:01:52', '2026-04-18 14:12:25', 0);
INSERT INTO `sys_dict_data` VALUES (201, 'xcode_component', 'radio', 0, 5, 'default', '', 0, '2026-04-16 16:02:04', '2026-04-18 14:12:38', 0);
INSERT INTO `sys_dict_data` VALUES (202, 'xcode_component', 'switch', 0, 6, 'default', '', 0, '2026-04-16 16:02:14', '2026-04-18 14:12:47', 0);
INSERT INTO `sys_dict_data` VALUES (203, 'xcode_component', 'textarea', 0, 7, 'default', '', 0, '2026-04-16 16:02:27', '2026-04-18 14:12:57', 0);
INSERT INTO `sys_dict_data` VALUES (204, 'xcode_component', 'multiple-select', 0, 8, 'default', '', 0, '2026-04-16 16:02:40', '2026-04-18 14:13:06', 0);
INSERT INTO `sys_dict_data` VALUES (205, 'xcode_component', 'tree-select', 0, 9, 'default', '', 0, '2026-04-16 16:02:54', '2026-04-18 14:13:25', 0);
INSERT INTO `sys_dict_data` VALUES (206, 'xcode_component', 'range-picker', 0, 10, 'default', '', 0, '2026-04-16 16:03:08', '2026-04-18 14:13:36', 0);
INSERT INTO `sys_dict_data` VALUES (207, 'xcode_component', 'date-picker', 0, 11, 'default', '', 0, '2026-04-16 16:03:20', '2026-04-18 14:13:44', 0);
INSERT INTO `sys_dict_data` VALUES (208, 'xcode_component', 'time-picker', 0, 12, '', '', 0, '2026-04-16 16:03:31', '2026-04-18 14:13:56', 0);
INSERT INTO `sys_dict_data` VALUES (209, 'xcode_component', 'input-tag', 0, 13, '', '', 0, '2026-04-16 16:03:44', '2026-04-18 14:14:11', 0);
INSERT INTO `sys_dict_data` VALUES (210, 'xcode_component', 'upload', 0, 14, '', '', 0, '2026-04-16 16:04:00', '2026-04-18 14:14:25', 0);
INSERT INTO `sys_dict_data` VALUES (211, 'xcode_component', 'slider', 0, 15, '', '', 0, '2026-04-16 16:04:13', '2026-04-18 14:14:33', 0);
INSERT INTO `sys_dict_data` VALUES (212, 'sys_authorized_grant_type', 'password', 0, 1, 'default', '', 0, '2026-04-16 16:51:29', '2026-04-16 16:51:51', 0);
INSERT INTO `sys_dict_data` VALUES (213, 'sys_authorized_grant_type', 'refresh_token', 0, 2, 'default', '', 0, '2026-04-16 16:51:45', '2026-04-16 16:51:45', 0);
INSERT INTO `sys_dict_data` VALUES (214, 'sys_authorized_grant_type', 'authorization_code', 0, 3, 'default', '', 0, '2026-04-16 16:52:10', '2026-04-16 16:52:10', 0);
INSERT INTO `sys_dict_data` VALUES (215, 'sys_authorized_grant_type', 'client_credentials', 0, 4, 'default', '', 0, '2026-04-16 16:52:31', '2026-04-16 16:52:31', 0);
INSERT INTO `sys_dict_data` VALUES (216, 'sys_authorized_grant_type', 'mobile', 0, 5, 'default', '', 0, '2026-04-16 16:52:53', '2026-04-16 16:52:53', 0);
INSERT INTO `sys_dict_data` VALUES (217, 'sys_authorized_grant_type', 'passkey', 0, 6, 'default', '', 0, '2026-04-16 16:53:08', '2026-04-16 16:53:08', 0);
INSERT INTO `sys_dict_data` VALUES (218, 'sys_authorized_grant_type', 'email', 0, 7, 'default', '', 0, '2026-04-16 16:53:32', '2026-04-16 16:58:36', 0);
INSERT INTO `sys_dict_data` VALUES (219, 'sys_authorized_scope', 'openid', 0, 1, 'default', '', 0, '2026-04-16 17:42:35', '2026-04-16 17:42:35', 0);
INSERT INTO `sys_dict_data` VALUES (220, 'sys_authorized_scope', 'profile', 0, 2, 'default', '', 0, '2026-04-16 17:42:48', '2026-04-16 17:42:48', 0);
INSERT INTO `sys_dict_data` VALUES (221, 'sys_menu_type', '1', 0, 1, 'default', '', 0, '2026-04-18 10:00:10', '2026-04-18 10:00:10', 0);
INSERT INTO `sys_dict_data` VALUES (222, 'sys_menu_type', '2', 0, 2, 'success', '', 0, '2026-04-18 10:00:40', '2026-04-18 10:00:40', 0);
INSERT INTO `sys_dict_data` VALUES (223, 'sys_menu_type', '3', 0, 3, 'purple', '', 0, '2026-04-18 10:01:20', '2026-04-18 10:03:27', 0);
INSERT INTO `sys_dict_data` VALUES (224, 'sys_dict_tag_style', 'default', 0, 1, 'default', '', 0, '2026-04-18 13:40:07', '2026-04-18 13:40:07', 0);
INSERT INTO `sys_dict_data` VALUES (225, 'sys_dict_tag_style', 'info', 0, 2, 'info', '', 0, '2026-04-18 13:40:34', '2026-04-18 13:40:34', 0);
INSERT INTO `sys_dict_data` VALUES (226, 'sys_dict_tag_style', 'success', 0, 3, 'success', '', 0, '2026-04-18 13:40:46', '2026-04-18 13:41:09', 0);
INSERT INTO `sys_dict_data` VALUES (227, 'sys_dict_tag_style', 'warning', 0, 4, 'warning', '', 0, '2026-04-18 13:41:04', '2026-04-18 13:41:04', 0);
INSERT INTO `sys_dict_data` VALUES (228, 'sys_dict_tag_style', 'danger', 0, 5, 'danger', '', 0, '2026-04-18 13:41:27', '2026-04-18 13:41:27', 0);
INSERT INTO `sys_dict_data` VALUES (233, 'sys_dict_tag_style', 'purple', 0, 10, 'purple', '', 0, '2026-04-18 13:43:07', '2026-04-18 13:43:57', 0);
INSERT INTO `sys_dict_data` VALUES (234, 'sys_dict_tag_style', 'cyan', 0, 11, 'cyan', '', 0, '2026-04-18 13:43:22', '2026-04-18 13:43:22', 0);
INSERT INTO `sys_dict_data` VALUES (235, 'sys_data_scope', '0', 0, 1, 'default', '', 0, '2026-04-18 19:19:09', '2026-04-18 19:30:44', 0);
INSERT INTO `sys_dict_data` VALUES (236, 'sys_data_scope', '1', 0, 4, 'purple', '', 0, '2026-04-18 19:23:55', '2026-04-18 19:30:44', 0);
INSERT INTO `sys_dict_data` VALUES (237, 'sys_data_scope', '2', 0, 2, 'success', '', 0, '2026-04-18 19:24:40', '2026-04-18 19:30:44', 0);
INSERT INTO `sys_dict_data` VALUES (238, 'sys_data_scope', '3', 0, 3, 'cyan', '', 0, '2026-04-18 19:25:13', '2026-04-18 19:30:44', 0);
INSERT INTO `sys_dict_data` VALUES (239, 'sys_success_status', '1', 0, 1, 'success', '', 0, '2026-04-18 22:51:44', '2026-04-20 15:41:38', 0);
INSERT INTO `sys_dict_data` VALUES (240, 'sys_locale_namespace', 'template', 0, 4, '', '', 0, '2026-04-20 14:46:06', '2026-04-20 14:46:06', 0);
INSERT INTO `sys_dict_data` VALUES (241, 'sys_success_status', '4', 0, 3, 'danger', '', 0, '2026-04-20 15:58:00', '2026-04-20 15:58:00', 0);
COMMIT;

-- ----------------------------
-- Table structure for sys_dict_data_tl
-- ----------------------------
DROP TABLE IF EXISTS `sys_dict_data_tl`;
CREATE TABLE `sys_dict_data_tl` (
  `dict_data_id` bigint NOT NULL COMMENT '字典数据ID',
  `lang_code` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '语言标识',
  `label` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '标签翻译',
  `tip` varchar(200) DEFAULT NULL COMMENT '提示翻译',
  PRIMARY KEY (`dict_data_id`,`lang_code`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='字典数据翻译';

-- ----------------------------
-- Records of sys_dict_data_tl
-- ----------------------------
BEGIN;
INSERT INTO `sys_dict_data_tl` VALUES (1, 'en-US', 'Disable', '');
INSERT INTO `sys_dict_data_tl` VALUES (1, 'zh-CN', '停用', '');
INSERT INTO `sys_dict_data_tl` VALUES (2, 'en-US', 'Normal', '');
INSERT INTO `sys_dict_data_tl` VALUES (2, 'zh-CN', '正常', '');
INSERT INTO `sys_dict_data_tl` VALUES (3, 'en-US', 'Yes', '');
INSERT INTO `sys_dict_data_tl` VALUES (3, 'zh-CN', '是', '');
INSERT INTO `sys_dict_data_tl` VALUES (4, 'en-US', 'No', '');
INSERT INTO `sys_dict_data_tl` VALUES (4, 'zh-CN', '否', '');
INSERT INTO `sys_dict_data_tl` VALUES (5, 'en-US', 'Regular user', '');
INSERT INTO `sys_dict_data_tl` VALUES (5, 'zh-CN', '普通用户', '');
INSERT INTO `sys_dict_data_tl` VALUES (6, 'en-US', 'Administrator', '');
INSERT INTO `sys_dict_data_tl` VALUES (6, 'zh-CN', '管理员', '');
INSERT INTO `sys_dict_data_tl` VALUES (7, 'en-US', 'Man', '');
INSERT INTO `sys_dict_data_tl` VALUES (7, 'zh-CN', '男', '');
INSERT INTO `sys_dict_data_tl` VALUES (8, 'en-US', 'Woman', '');
INSERT INTO `sys_dict_data_tl` VALUES (8, 'zh-CN', '女', '');
INSERT INTO `sys_dict_data_tl` VALUES (9, 'en-US', 'Unknown', '');
INSERT INTO `sys_dict_data_tl` VALUES (9, 'zh-CN', '未知', '');
INSERT INTO `sys_dict_data_tl` VALUES (10, 'en-US', 'Show', '');
INSERT INTO `sys_dict_data_tl` VALUES (10, 'zh-CN', '显示', '');
INSERT INTO `sys_dict_data_tl` VALUES (11, 'en-US', 'Hide', '');
INSERT INTO `sys_dict_data_tl` VALUES (11, 'zh-CN', '隐藏', '');
INSERT INTO `sys_dict_data_tl` VALUES (12, 'en-US', 'Public', '');
INSERT INTO `sys_dict_data_tl` VALUES (12, 'zh-CN', '公开配置', '');
INSERT INTO `sys_dict_data_tl` VALUES (13, 'en-US', 'Login', '');
INSERT INTO `sys_dict_data_tl` VALUES (13, 'zh-CN', '登录配置', '');
INSERT INTO `sys_dict_data_tl` VALUES (14, 'en-US', 'Backend', '');
INSERT INTO `sys_dict_data_tl` VALUES (14, 'zh-CN', '后台服务配置', '');
INSERT INTO `sys_dict_data_tl` VALUES (15, 'en-US', 'Other', '');
INSERT INTO `sys_dict_data_tl` VALUES (15, 'zh-CN', '其他', '');
INSERT INTO `sys_dict_data_tl` VALUES (16, 'en-US', 'Add', '');
INSERT INTO `sys_dict_data_tl` VALUES (16, 'zh-CN', '新增', '');
INSERT INTO `sys_dict_data_tl` VALUES (17, 'en-US', 'Modify', '');
INSERT INTO `sys_dict_data_tl` VALUES (17, 'zh-CN', '修改', '');
INSERT INTO `sys_dict_data_tl` VALUES (18, 'en-US', 'Delete', '');
INSERT INTO `sys_dict_data_tl` VALUES (18, 'zh-CN', '删除', '');
INSERT INTO `sys_dict_data_tl` VALUES (19, 'en-US', 'Authorization', '');
INSERT INTO `sys_dict_data_tl` VALUES (19, 'zh-CN', '授权', '');
INSERT INTO `sys_dict_data_tl` VALUES (20, 'en-US', 'Export', '');
INSERT INTO `sys_dict_data_tl` VALUES (20, 'zh-CN', '导出', '');
INSERT INTO `sys_dict_data_tl` VALUES (21, 'en-US', 'Import', '');
INSERT INTO `sys_dict_data_tl` VALUES (21, 'zh-CN', '导入', '');
INSERT INTO `sys_dict_data_tl` VALUES (22, 'en-US', 'Force quit', '');
INSERT INTO `sys_dict_data_tl` VALUES (22, 'zh-CN', '强退', '');
INSERT INTO `sys_dict_data_tl` VALUES (23, 'en-US', 'Clear', '');
INSERT INTO `sys_dict_data_tl` VALUES (23, 'zh-CN', '清空', '');
INSERT INTO `sys_dict_data_tl` VALUES (24, 'en-US', 'Setup', '');
INSERT INTO `sys_dict_data_tl` VALUES (24, 'zh-CN', '设置', '');
INSERT INTO `sys_dict_data_tl` VALUES (25, 'en-US', 'Save', '');
INSERT INTO `sys_dict_data_tl` VALUES (25, 'zh-CN', '保存', '');
INSERT INTO `sys_dict_data_tl` VALUES (26, 'en-US', 'Release', '');
INSERT INTO `sys_dict_data_tl` VALUES (26, 'zh-CN', '发布', '');
INSERT INTO `sys_dict_data_tl` VALUES (27, 'en-US', 'Copy', '');
INSERT INTO `sys_dict_data_tl` VALUES (27, 'zh-CN', '复制', '');
INSERT INTO `sys_dict_data_tl` VALUES (29, 'en-US', 'Failure', '');
INSERT INTO `sys_dict_data_tl` VALUES (29, 'zh-CN', '失败', '');
INSERT INTO `sys_dict_data_tl` VALUES (30, 'en-US', 'SMS Template', '');
INSERT INTO `sys_dict_data_tl` VALUES (30, 'zh-CN', '短信模板', '');
INSERT INTO `sys_dict_data_tl` VALUES (31, 'en-US', 'App Template', '');
INSERT INTO `sys_dict_data_tl` VALUES (31, 'zh-CN', 'APP模板', '');
INSERT INTO `sys_dict_data_tl` VALUES (32, 'en-US', 'WeChat Template', '');
INSERT INTO `sys_dict_data_tl` VALUES (32, 'zh-CN', '微信模板', '');
INSERT INTO `sys_dict_data_tl` VALUES (33, 'en-US', 'PC Template', '');
INSERT INTO `sys_dict_data_tl` VALUES (33, 'zh-CN', 'PC模板', '');
INSERT INTO `sys_dict_data_tl` VALUES (34, 'en-US', 'Notice', '');
INSERT INTO `sys_dict_data_tl` VALUES (34, 'zh-CN', '通知', '');
INSERT INTO `sys_dict_data_tl` VALUES (35, 'en-US', 'Warning', '');
INSERT INTO `sys_dict_data_tl` VALUES (35, 'zh-CN', '警告', '');
INSERT INTO `sys_dict_data_tl` VALUES (36, 'en-US', 'Error', '');
INSERT INTO `sys_dict_data_tl` VALUES (36, 'zh-CN', '异常', '');
INSERT INTO `sys_dict_data_tl` VALUES (37, 'en-US', 'Windows', '');
INSERT INTO `sys_dict_data_tl` VALUES (37, 'zh-CN', 'Windows', '');
INSERT INTO `sys_dict_data_tl` VALUES (38, 'en-US', 'Android', '');
INSERT INTO `sys_dict_data_tl` VALUES (38, 'zh-CN', 'Android', '');
INSERT INTO `sys_dict_data_tl` VALUES (39, 'en-US', 'IOS', '');
INSERT INTO `sys_dict_data_tl` VALUES (39, 'zh-CN', 'IOS', '');
INSERT INTO `sys_dict_data_tl` VALUES (40, 'en-US', 'Organization', '');
INSERT INTO `sys_dict_data_tl` VALUES (40, 'zh-CN', '机构', '');
INSERT INTO `sys_dict_data_tl` VALUES (41, 'en-US', 'Department', '');
INSERT INTO `sys_dict_data_tl` VALUES (41, 'zh-CN', '部门', '');
INSERT INTO `sys_dict_data_tl` VALUES (42, 'en-US', 'Notice', '');
INSERT INTO `sys_dict_data_tl` VALUES (42, 'zh-CN', '通知', '');
INSERT INTO `sys_dict_data_tl` VALUES (43, 'en-US', 'Announcement', '');
INSERT INTO `sys_dict_data_tl` VALUES (43, 'zh-CN', '公告', '');
INSERT INTO `sys_dict_data_tl` VALUES (46, 'en-US', 'Approval process', '');
INSERT INTO `sys_dict_data_tl` VALUES (46, 'zh-CN', '审批流程', '');
INSERT INTO `sys_dict_data_tl` VALUES (47, 'en-US', 'Business Process', '');
INSERT INTO `sys_dict_data_tl` VALUES (47, 'zh-CN', '业务流程', '');
INSERT INTO `sys_dict_data_tl` VALUES (48, 'en-US', 'Not submitted', '');
INSERT INTO `sys_dict_data_tl` VALUES (48, 'zh-CN', '未提交', '');
INSERT INTO `sys_dict_data_tl` VALUES (49, 'en-US', 'Under Approval', '');
INSERT INTO `sys_dict_data_tl` VALUES (49, 'zh-CN', '审批中', '');
INSERT INTO `sys_dict_data_tl` VALUES (50, 'en-US', 'Revoked', '');
INSERT INTO `sys_dict_data_tl` VALUES (50, 'zh-CN', '已撤销', '');
INSERT INTO `sys_dict_data_tl` VALUES (51, 'en-US', 'Rejected', '');
INSERT INTO `sys_dict_data_tl` VALUES (51, 'zh-CN', '已驳回', '');
INSERT INTO `sys_dict_data_tl` VALUES (52, 'en-US', 'Approved', '');
INSERT INTO `sys_dict_data_tl` VALUES (52, 'zh-CN', '已通过', '');
INSERT INTO `sys_dict_data_tl` VALUES (53, 'en-US', 'Normal', '');
INSERT INTO `sys_dict_data_tl` VALUES (53, 'zh-CN', '正常', '');
INSERT INTO `sys_dict_data_tl` VALUES (54, 'en-US', 'Outdated', '');
INSERT INTO `sys_dict_data_tl` VALUES (54, 'zh-CN', '过时', '');
INSERT INTO `sys_dict_data_tl` VALUES (55, 'zh-CN', '普通数据表', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (56, 'zh-CN', '工作流数据表', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (57, 'zh-CN', '树形结构表', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (58, 'zh-CN', '工作流节点数据表', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (60, 'zh-CN', '自定义', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (61, 'zh-CN', '同级排序', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (62, 'zh-CN', '全局排序', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (63, 'zh-CN', 'bigint', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (64, 'zh-CN', 'bit', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (65, 'zh-CN', 'char', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (66, 'zh-CN', 'date', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (67, 'zh-CN', 'datetime', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (68, 'zh-CN', 'decimal', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (69, 'zh-CN', 'int', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (70, 'zh-CN', 'text', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (71, 'zh-CN', 'tinyint', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (72, 'zh-CN', 'varchar', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (73, 'zh-CN', '等于', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (74, 'zh-CN', '不等于', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (75, 'zh-CN', '大于', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (76, 'zh-CN', '大于等于', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (77, 'zh-CN', '小于', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (78, 'zh-CN', '小于等于', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (79, 'zh-CN', '包含', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (80, 'zh-CN', '左包含', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (81, 'zh-CN', '右包含', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (82, 'zh-CN', '区间', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (83, 'zh-CN', '不是NULL', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (84, 'zh-CN', '是NULL', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (85, 'zh-CN', '不是空', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (86, 'zh-CN', '是空', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (87, 'zh-CN', '在列表中', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (88, 'zh-CN', '不在列表中', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (89, 'zh-CN', '文本', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (90, 'zh-CN', '数字', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (91, 'zh-CN', '大文本', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (92, 'zh-CN', '密码', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (93, 'zh-CN', '颜色', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (94, 'zh-CN', '图标', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (95, 'zh-CN', '富文本', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (96, 'zh-CN', '日期', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (97, 'zh-CN', '时间', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (98, 'zh-CN', '下拉选择', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (99, 'zh-CN', '下拉选择（多选）', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (100, 'zh-CN', '单选框', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (101, 'zh-CN', '多选框', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (102, 'zh-CN', '文件', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (103, 'zh-CN', '时间段', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (104, 'zh-CN', '日期段', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (105, 'zh-CN', '大文件', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (106, 'zh-CN', '单选按钮', '');
INSERT INTO `sys_dict_data_tl` VALUES (107, 'zh-CN', '表达式', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (108, 'zh-CN', '枚举', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (109, 'zh-CN', '顺序', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (110, 'zh-CN', '一对一关系', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (111, 'zh-CN', '一对多关系', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (112, 'zh-CN', '多对一关系', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (113, 'zh-CN', '多对多关系', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (128, 'en-US', 'Active', '');
INSERT INTO `sys_dict_data_tl` VALUES (128, 'zh-CN', '激活', '');
INSERT INTO `sys_dict_data_tl` VALUES (129, 'zh-CN', '挂起', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (130, 'zh-CN', '表格', '');
INSERT INTO `sys_dict_data_tl` VALUES (131, 'zh-CN', '列表', '');
INSERT INTO `sys_dict_data_tl` VALUES (132, 'zh-CN', '卡片', '');
INSERT INTO `sys_dict_data_tl` VALUES (133, 'zh-CN', '表单', '');
INSERT INTO `sys_dict_data_tl` VALUES (134, 'zh-CN', '详情', '');
INSERT INTO `sys_dict_data_tl` VALUES (138, 'zh-CN', '下拉树', '');
INSERT INTO `sys_dict_data_tl` VALUES (139, 'zh-CN', '暂停', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (140, 'zh-CN', '正常', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (141, 'zh-CN', '触发续跑', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (142, 'zh-CN', '自然续跑', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (143, 'zh-CN', '全量回溯', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (144, 'zh-CN', '禁止', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (145, 'zh-CN', '允许', NULL);
INSERT INTO `sys_dict_data_tl` VALUES (148, 'en-US', 'Computer', '');
INSERT INTO `sys_dict_data_tl` VALUES (148, 'zh-CN', '电脑', '');
INSERT INTO `sys_dict_data_tl` VALUES (149, 'en-US', 'Mobile', '');
INSERT INTO `sys_dict_data_tl` VALUES (149, 'zh-CN', '手机', '');
INSERT INTO `sys_dict_data_tl` VALUES (150, 'en-US', 'Tablet', '');
INSERT INTO `sys_dict_data_tl` VALUES (150, 'zh-CN', '平板', '');
INSERT INTO `sys_dict_data_tl` VALUES (151, 'en-US', 'Unknown', '');
INSERT INTO `sys_dict_data_tl` VALUES (151, 'zh-CN', '未知', '');
INSERT INTO `sys_dict_data_tl` VALUES (152, 'en-US', 'Disable', '');
INSERT INTO `sys_dict_data_tl` VALUES (152, 'zh-CN', '停用', '');
INSERT INTO `sys_dict_data_tl` VALUES (153, 'en-US', 'Active', '');
INSERT INTO `sys_dict_data_tl` VALUES (153, 'zh-CN', '激活', '');
INSERT INTO `sys_dict_data_tl` VALUES (154, 'en-US', 'Default', '');
INSERT INTO `sys_dict_data_tl` VALUES (154, 'zh-CN', '默认', '');
INSERT INTO `sys_dict_data_tl` VALUES (174, 'en-US', 'EQ', '');
INSERT INTO `sys_dict_data_tl` VALUES (174, 'ru-RU', 'EQ', '');
INSERT INTO `sys_dict_data_tl` VALUES (174, 'zh-CN', '等于', '');
INSERT INTO `sys_dict_data_tl` VALUES (175, 'en-US', 'NE', '');
INSERT INTO `sys_dict_data_tl` VALUES (175, 'ru-RU', 'NE', '');
INSERT INTO `sys_dict_data_tl` VALUES (175, 'zh-CN', '不等于', '');
INSERT INTO `sys_dict_data_tl` VALUES (176, 'en-US', 'GT', '');
INSERT INTO `sys_dict_data_tl` VALUES (176, 'ru-RU', 'GT', '');
INSERT INTO `sys_dict_data_tl` VALUES (176, 'zh-CN', '大于', '');
INSERT INTO `sys_dict_data_tl` VALUES (177, 'en-US', 'GTE', '');
INSERT INTO `sys_dict_data_tl` VALUES (177, 'ru-RU', 'GTE', '');
INSERT INTO `sys_dict_data_tl` VALUES (177, 'zh-CN', '大于等于', '');
INSERT INTO `sys_dict_data_tl` VALUES (178, 'en-US', 'LT', '');
INSERT INTO `sys_dict_data_tl` VALUES (178, 'ru-RU', 'LT', '');
INSERT INTO `sys_dict_data_tl` VALUES (178, 'zh-CN', '小于', '');
INSERT INTO `sys_dict_data_tl` VALUES (179, 'en-US', 'LTE', '');
INSERT INTO `sys_dict_data_tl` VALUES (179, 'ru-RU', 'LTE', '');
INSERT INTO `sys_dict_data_tl` VALUES (179, 'zh-CN', '小于等于', '');
INSERT INTO `sys_dict_data_tl` VALUES (180, 'en-US', 'Like', '');
INSERT INTO `sys_dict_data_tl` VALUES (180, 'ru-RU', 'Like', '');
INSERT INTO `sys_dict_data_tl` VALUES (180, 'zh-CN', '模糊', '');
INSERT INTO `sys_dict_data_tl` VALUES (181, 'en-US', 'Like left', '');
INSERT INTO `sys_dict_data_tl` VALUES (181, 'ru-RU', 'Like left', '');
INSERT INTO `sys_dict_data_tl` VALUES (181, 'zh-CN', '左模糊', '');
INSERT INTO `sys_dict_data_tl` VALUES (182, 'en-US', 'Like right', '');
INSERT INTO `sys_dict_data_tl` VALUES (182, 'ru-RU', 'Like right', '');
INSERT INTO `sys_dict_data_tl` VALUES (182, 'zh-CN', '右模糊', '');
INSERT INTO `sys_dict_data_tl` VALUES (183, 'en-US', 'In', '');
INSERT INTO `sys_dict_data_tl` VALUES (183, 'ru-RU', 'In', '');
INSERT INTO `sys_dict_data_tl` VALUES (183, 'zh-CN', '包含', '');
INSERT INTO `sys_dict_data_tl` VALUES (184, 'en-US', 'Not in', '');
INSERT INTO `sys_dict_data_tl` VALUES (184, 'ru-RU', 'Not in', '');
INSERT INTO `sys_dict_data_tl` VALUES (184, 'zh-CN', '不包含', '');
INSERT INTO `sys_dict_data_tl` VALUES (185, 'en-US', 'Between', '');
INSERT INTO `sys_dict_data_tl` VALUES (185, 'ru-RU', 'Between', '');
INSERT INTO `sys_dict_data_tl` VALUES (185, 'zh-CN', '区间', '');
INSERT INTO `sys_dict_data_tl` VALUES (186, 'en-US', 'Is null', '');
INSERT INTO `sys_dict_data_tl` VALUES (186, 'ru-RU', 'Is null', '');
INSERT INTO `sys_dict_data_tl` VALUES (186, 'zh-CN', '为NILL', '');
INSERT INTO `sys_dict_data_tl` VALUES (187, 'en-US', 'Not null', '');
INSERT INTO `sys_dict_data_tl` VALUES (187, 'ru-RU', 'Not null', '');
INSERT INTO `sys_dict_data_tl` VALUES (187, 'zh-CN', '不为NULL', '');
INSERT INTO `sys_dict_data_tl` VALUES (197, 'en-US', 'Input', '');
INSERT INTO `sys_dict_data_tl` VALUES (197, 'zh-CN', '输入框', '');
INSERT INTO `sys_dict_data_tl` VALUES (198, 'en-US', 'Input number', '');
INSERT INTO `sys_dict_data_tl` VALUES (198, 'zh-CN', '数字输入框', '');
INSERT INTO `sys_dict_data_tl` VALUES (199, 'en-US', 'Select', '');
INSERT INTO `sys_dict_data_tl` VALUES (199, 'zh-CN', '选择器', '');
INSERT INTO `sys_dict_data_tl` VALUES (200, 'en-US', 'Checkbox', '');
INSERT INTO `sys_dict_data_tl` VALUES (200, 'zh-CN', '多选框', '');
INSERT INTO `sys_dict_data_tl` VALUES (201, 'en-US', 'Radio', '');
INSERT INTO `sys_dict_data_tl` VALUES (201, 'zh-CN', '单选框', '');
INSERT INTO `sys_dict_data_tl` VALUES (202, 'en-US', 'Switch', '');
INSERT INTO `sys_dict_data_tl` VALUES (202, 'zh-CN', '开关', '');
INSERT INTO `sys_dict_data_tl` VALUES (203, 'en-US', 'Textarea', '');
INSERT INTO `sys_dict_data_tl` VALUES (203, 'zh-CN', '文本域', '');
INSERT INTO `sys_dict_data_tl` VALUES (204, 'en-US', 'Multiple select', '');
INSERT INTO `sys_dict_data_tl` VALUES (204, 'zh-CN', '多选选择器 multiple-select', '');
INSERT INTO `sys_dict_data_tl` VALUES (205, 'en-US', 'Tree select', '');
INSERT INTO `sys_dict_data_tl` VALUES (205, 'zh-CN', '树型选择器', '');
INSERT INTO `sys_dict_data_tl` VALUES (206, 'en-US', 'Range picker', '');
INSERT INTO `sys_dict_data_tl` VALUES (206, 'zh-CN', '日期段选择器', '');
INSERT INTO `sys_dict_data_tl` VALUES (207, 'en-US', 'Date picker', '');
INSERT INTO `sys_dict_data_tl` VALUES (207, 'zh-CN', '日期选择器', '');
INSERT INTO `sys_dict_data_tl` VALUES (208, 'en-US', 'Time picker', '');
INSERT INTO `sys_dict_data_tl` VALUES (208, 'zh-CN', '时间选择器', '');
INSERT INTO `sys_dict_data_tl` VALUES (209, 'en-US', 'Input tag', '');
INSERT INTO `sys_dict_data_tl` VALUES (209, 'zh-CN', '标签输入框', '');
INSERT INTO `sys_dict_data_tl` VALUES (210, 'en-US', 'Upload file', '');
INSERT INTO `sys_dict_data_tl` VALUES (210, 'zh-CN', '上传文件', '');
INSERT INTO `sys_dict_data_tl` VALUES (211, 'en-US', 'Slider', '');
INSERT INTO `sys_dict_data_tl` VALUES (211, 'zh-CN', '滑块', '');
INSERT INTO `sys_dict_data_tl` VALUES (212, 'en-US', 'Password', '');
INSERT INTO `sys_dict_data_tl` VALUES (212, 'zh-CN', '密码模式', '');
INSERT INTO `sys_dict_data_tl` VALUES (213, 'en-US', 'Refresh token', '');
INSERT INTO `sys_dict_data_tl` VALUES (213, 'zh-CN', '刷新令牌', '');
INSERT INTO `sys_dict_data_tl` VALUES (214, 'en-US', 'Authorization code', '');
INSERT INTO `sys_dict_data_tl` VALUES (214, 'zh-CN', '授权码模式', '');
INSERT INTO `sys_dict_data_tl` VALUES (215, 'en-US', 'Client credentials', '');
INSERT INTO `sys_dict_data_tl` VALUES (215, 'zh-CN', '终端模式', '');
INSERT INTO `sys_dict_data_tl` VALUES (216, 'en-US', 'Mobile', '');
INSERT INTO `sys_dict_data_tl` VALUES (216, 'zh-CN', '手机验证模式', '');
INSERT INTO `sys_dict_data_tl` VALUES (217, 'en-US', 'Passkey', '');
INSERT INTO `sys_dict_data_tl` VALUES (217, 'zh-CN', 'Passkey模式', '');
INSERT INTO `sys_dict_data_tl` VALUES (218, 'en-US', 'Email', '');
INSERT INTO `sys_dict_data_tl` VALUES (218, 'zh-CN', '邮箱验证模式', '');
INSERT INTO `sys_dict_data_tl` VALUES (219, 'en-US', 'Open ID', '');
INSERT INTO `sys_dict_data_tl` VALUES (219, 'zh-CN', 'Open ID', '');
INSERT INTO `sys_dict_data_tl` VALUES (220, 'en-US', 'Profile', '');
INSERT INTO `sys_dict_data_tl` VALUES (220, 'zh-CN', 'Profile', '');
INSERT INTO `sys_dict_data_tl` VALUES (221, 'en-US', 'Directory', '');
INSERT INTO `sys_dict_data_tl` VALUES (221, 'zh-CN', '目录', '');
INSERT INTO `sys_dict_data_tl` VALUES (222, 'en-US', 'Menu', '');
INSERT INTO `sys_dict_data_tl` VALUES (222, 'zh-CN', '菜单', '');
INSERT INTO `sys_dict_data_tl` VALUES (223, 'en-US', 'Function', '');
INSERT INTO `sys_dict_data_tl` VALUES (223, 'zh-CN', '功能', '');
INSERT INTO `sys_dict_data_tl` VALUES (224, 'en-US', 'Default', '');
INSERT INTO `sys_dict_data_tl` VALUES (224, 'zh-CN', '默认', '');
INSERT INTO `sys_dict_data_tl` VALUES (225, 'en-US', 'Info', '');
INSERT INTO `sys_dict_data_tl` VALUES (225, 'zh-CN', '信息', '');
INSERT INTO `sys_dict_data_tl` VALUES (226, 'en-US', 'Success', '');
INSERT INTO `sys_dict_data_tl` VALUES (226, 'zh-CN', '成功', '');
INSERT INTO `sys_dict_data_tl` VALUES (227, 'en-US', 'Warning', '');
INSERT INTO `sys_dict_data_tl` VALUES (227, 'zh-CN', '警告', '');
INSERT INTO `sys_dict_data_tl` VALUES (228, 'en-US', 'Danger', '');
INSERT INTO `sys_dict_data_tl` VALUES (228, 'zh-CN', '危险', '');
INSERT INTO `sys_dict_data_tl` VALUES (233, 'en-US', 'Purple', '');
INSERT INTO `sys_dict_data_tl` VALUES (233, 'zh-CN', '紫色', '');
INSERT INTO `sys_dict_data_tl` VALUES (234, 'en-US', 'Cyan', '');
INSERT INTO `sys_dict_data_tl` VALUES (234, 'zh-CN', '青色', '');
INSERT INTO `sys_dict_data_tl` VALUES (235, 'en-US', 'All', '');
INSERT INTO `sys_dict_data_tl` VALUES (235, 'zh-CN', '全部', '');
INSERT INTO `sys_dict_data_tl` VALUES (236, 'en-US', 'Customize', '');
INSERT INTO `sys_dict_data_tl` VALUES (236, 'zh-CN', '自定义', '');
INSERT INTO `sys_dict_data_tl` VALUES (237, 'en-US', 'This department', '');
INSERT INTO `sys_dict_data_tl` VALUES (237, 'zh-CN', '本部门', '');
INSERT INTO `sys_dict_data_tl` VALUES (238, 'en-US', 'This department and below', '');
INSERT INTO `sys_dict_data_tl` VALUES (238, 'zh-CN', '本部门及以下', '');
INSERT INTO `sys_dict_data_tl` VALUES (239, 'en-US', 'Success', '');
INSERT INTO `sys_dict_data_tl` VALUES (239, 'zh-CN', '成功', '');
INSERT INTO `sys_dict_data_tl` VALUES (240, 'en-US', 'Template', '');
INSERT INTO `sys_dict_data_tl` VALUES (240, 'zh-CN', '模板', '');
INSERT INTO `sys_dict_data_tl` VALUES (241, 'en-US', 'Error', '');
INSERT INTO `sys_dict_data_tl` VALUES (241, 'zh-CN', '异常', '');
COMMIT;

-- ----------------------------
-- Table structure for sys_job
-- ----------------------------
DROP TABLE IF EXISTS `sys_job`;
CREATE TABLE `sys_job` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `job_name` varchar(64) DEFAULT NULL COMMENT '任务名称',
  `job_group` varchar(64) DEFAULT NULL COMMENT '任务分组',
  `job_key` varchar(64) NOT NULL COMMENT '任务唯一标识(不可重复)',
  `cron_expr` varchar(32) DEFAULT NULL COMMENT 'cron执行表达式',
  `invoke_target` varchar(128) DEFAULT NULL COMMENT '执行目标：服务.方法名',
  `status` tinyint DEFAULT '0' COMMENT '状态 0正常 1暂停',
  `concurrent` tinyint DEFAULT '0' COMMENT '是否并发 0禁止 1允许',
  `retry_count` int DEFAULT '0' COMMENT '失败重试次数',
  `timeout` int DEFAULT '0' COMMENT '执行超时时间(秒)',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注说明',
  `created_by` bigint DEFAULT NULL COMMENT '创建人ID',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `updated_by` bigint DEFAULT NULL COMMENT '更新人ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `job_key` (`job_key`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='定时任务配置表';

-- ----------------------------
-- Records of sys_job
-- ----------------------------
BEGIN;
INSERT INTO `sys_job` VALUES (1, 'sys.test', '', 'sys.test', '0 */1 * * * *', 'sys.test', 0, 0, 0, 30, '', 0, '2026-04-19 16:13:53', '2026-05-18 18:03:12', 0);
COMMIT;

-- ----------------------------
-- Table structure for sys_job_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_job_log`;
CREATE TABLE `sys_job_log` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `job_id` bigint NOT NULL COMMENT '任务ID',
  `job_name` varchar(64) DEFAULT NULL COMMENT '任务名称（冗余）',
  `job_group` varchar(64) DEFAULT NULL COMMENT '任务分组（冗余）',
  `status` tinyint DEFAULT NULL COMMENT '执行状态 0成功 1失败',
  `message` text COMMENT '执行日志/异常信息',
  `duration` int DEFAULT NULL COMMENT '执行耗时(毫秒)',
  `trigger_type` varchar(16) DEFAULT NULL COMMENT '触发类型：自动/手动',
  `start_time` datetime DEFAULT NULL COMMENT '开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '结束时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '日志生成时间',
  PRIMARY KEY (`id`),
  KEY `idx_job_id` (`job_id`)
) ENGINE=InnoDB AUTO_INCREMENT=208 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='定时任务执行日志表';

-- ----------------------------
-- Records of sys_job_log
-- ----------------------------
BEGIN;
INSERT INTO `sys_job_log` VALUES (205, 1, 'sys.test', '', 0, '[job] sys.test fired at 2026-05-18 18:03:16', 0, 'manual', '2026-05-18 18:03:17', '2026-05-18 18:03:17', '2026-05-18 18:03:17');
INSERT INTO `sys_job_log` VALUES (206, 1, 'sys.test', '', 0, '[job] sys.test fired at 2026-05-18 18:04:00', 0, 'auto', '2026-05-18 18:04:00', '2026-05-18 18:04:00', '2026-05-18 18:04:00');
INSERT INTO `sys_job_log` VALUES (207, 1, 'sys.test', '', 0, '[job] sys.test fired at 2026-05-18 18:05:00', 0, 'auto', '2026-05-18 18:05:00', '2026-05-18 18:05:00', '2026-05-18 18:05:00');
COMMIT;

-- ----------------------------
-- Table structure for sys_lang
-- ----------------------------
DROP TABLE IF EXISTS `sys_lang`;
CREATE TABLE `sys_lang` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '编号',
  `lang_code` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '语言编码（如zh-CN、en-US）',
  `lang_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '语言名称（如简体中文、English）',
  `is_default` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否默认语言（0否，1是）',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序值（升序）',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态（0正常,1停用）',
  `remark` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_lang_code` (`lang_code`) USING BTREE COMMENT '语言编码唯一，避免重复'
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='语言';

-- ----------------------------
-- Records of sys_lang
-- ----------------------------
BEGIN;
INSERT INTO `sys_lang` VALUES (1, 'zh-CN', '🇨🇳 简体中文', 1, 1, 0, '中文', '2026-01-29 17:39:01', '2026-04-16 09:38:44');
INSERT INTO `sys_lang` VALUES (2, 'en-US', '🇺🇸 English', 0, 2, 0, '英文', '2026-01-29 17:39:43', '2026-04-16 09:38:48');
INSERT INTO `sys_lang` VALUES (3, 'ru-Ru', '🇷🇺 Русский', 0, 3, 1, '俄语', '2026-02-01 15:57:25', '2026-04-25 20:20:50');
COMMIT;

-- ----------------------------
-- Table structure for sys_lang_resource
-- ----------------------------
DROP TABLE IF EXISTS `sys_lang_resource`;
CREATE TABLE `sys_lang_resource` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `resource_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '资源标识',
  `namespace` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'default' COMMENT '命名空间（default）',
  `lang_code` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '语言编码（如zh）',
  `content` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '内容',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态（0正常,1停用）',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uk_namespace_key_lang` (`namespace`,`resource_key`,`lang_code`) USING BTREE COMMENT '命名空间+KEY+语言+未删除 唯一',
  KEY `idx_namespace_lang` (`namespace`,`lang_code`) USING BTREE,
  KEY `idx_resource_key` (`resource_key`) USING BTREE,
  KEY `fk_i18n_resource_lang` (`lang_code`) USING BTREE,
  CONSTRAINT `sys_lang_resource_ibfk_1` FOREIGN KEY (`lang_code`) REFERENCES `sys_lang` (`lang_code`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2363 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='语言资源';

-- ----------------------------
-- Records of sys_lang_resource
-- ----------------------------
BEGIN;
INSERT INTO `sys_lang_resource` VALUES (1, 'common.button.submit', 'default', 'zh-CN', '提交', 0, '2026-01-30 10:27:46', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (2, 'common.button.confirm', 'default', 'zh-CN', '确认', 0, '2026-01-30 11:49:00', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (3, 'common.button.cancel', 'default', 'zh-CN', '取消', 0, '2026-01-30 11:49:00', '2026-02-03 12:27:52');
INSERT INTO `sys_lang_resource` VALUES (4, 'common.button.add', 'default', 'zh-CN', '新增', 0, '2026-01-30 11:49:00', '2026-02-03 12:27:52');
INSERT INTO `sys_lang_resource` VALUES (5, 'common.button.save', 'default', 'zh-CN', '保存', 0, '2026-01-30 11:49:00', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (6, 'common.button.setup', 'default', 'zh-CN', '设置', 0, '2026-01-30 11:49:00', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (7, 'common.button.delete', 'default', 'zh-CN', '删除', 0, '2026-01-30 11:49:00', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (8, 'common.button.modify', 'default', 'zh-CN', '修改', 0, '2026-01-30 11:49:00', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (9, 'common.button.export', 'default', 'zh-CN', '导出', 0, '2026-01-30 11:49:00', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (10, 'common.button.import', 'default', 'zh-CN', '导入', 0, '2026-01-30 11:49:00', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (11, 'common.button.batch_delete', 'default', 'zh-CN', '批量删除', 0, '2026-01-30 11:49:00', '2026-02-03 12:27:52');
INSERT INTO `sys_lang_resource` VALUES (12, 'common.button.reset', 'default', 'zh-CN', '重置', 0, '2026-01-30 11:52:22', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (13, 'common.button.clear', 'default', 'zh-CN', '清空', 0, '2026-01-30 11:52:22', '2026-02-03 12:27:52');
INSERT INTO `sys_lang_resource` VALUES (14, 'common.button.refresh', 'default', 'zh-CN', '刷新', 0, '2026-01-30 11:52:22', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (15, 'common.button.confirm', 'default', 'en-US', 'confirm', 0, '2026-01-30 11:49:00', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (16, 'common.button.cancel', 'default', 'en-US', 'cancel', 0, '2026-01-30 11:49:00', '2026-02-03 12:27:52');
INSERT INTO `sys_lang_resource` VALUES (17, 'common.button.add', 'default', 'en-US', 'add', 0, '2026-01-30 11:49:00', '2026-02-03 12:27:52');
INSERT INTO `sys_lang_resource` VALUES (18, 'common.button.save', 'default', 'en-US', 'save', 0, '2026-01-30 11:49:00', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (19, 'common.button.setup', 'default', 'en-US', 'Settings', 0, '2026-01-30 11:49:00', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (20, 'common.button.delete', 'default', 'en-US', 'delete', 0, '2026-01-30 11:49:00', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (21, 'common.button.modify', 'default', 'en-US', 'modify', 0, '2026-01-30 11:49:00', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (22, 'common.button.export', 'default', 'en-US', 'export', 0, '2026-01-30 11:49:00', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (23, 'common.button.import', 'default', 'en-US', 'import', 0, '2026-01-30 11:49:00', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (24, 'common.button.batch_delete', 'default', 'en-US', 'Batch Delete', 0, '2026-01-30 11:49:00', '2026-02-03 12:27:52');
INSERT INTO `sys_lang_resource` VALUES (25, 'common.button.reset', 'default', 'en-US', 'reset', 0, '2026-01-30 11:52:22', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (26, 'common.button.clear', 'default', 'en-US', 'clear', 0, '2026-01-30 11:52:22', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (27, 'common.button.refresh', 'default', 'en-US', 'refresh', 0, '2026-01-30 11:52:22', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (30, 'common.button.copy', 'default', 'zh-CN', '复制', 0, '2026-01-30 11:52:22', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (31, 'common.button.move', 'default', 'zh-CN', '移动', 0, '2026-01-30 11:52:22', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (32, 'common.button.search', 'default', 'zh-CN', '搜索', 0, '2026-01-30 11:52:22', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (33, 'common.button.update', 'default', 'zh-CN', '上传', 0, '2026-01-30 11:52:22', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (34, 'common.button.download', 'default', 'zh-CN', '下载', 0, '2026-01-30 11:52:22', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (35, 'common.button.copy', 'default', 'en-US', 'copy', 0, '2026-01-30 11:52:22', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (36, 'common.button.move', 'default', 'en-US', 'move', 0, '2026-01-30 11:52:22', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (37, 'common.button.search', 'default', 'en-US', 'search', 0, '2026-01-30 11:52:22', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (38, 'common.button.update', 'default', 'en-US', 'Upload', 0, '2026-01-30 11:52:22', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (39, 'common.button.download', 'default', 'en-US', 'download', 0, '2026-01-30 11:52:22', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (1832, 'common.button.add', 'default', 'ru-Ru', 'добавленный', 0, '2026-02-01 17:22:01', '2026-02-03 12:27:52');
INSERT INTO `sys_lang_resource` VALUES (1833, 'common.button.batch_delete', 'default', 'ru-Ru', 'Массовое удаление', 0, '2026-02-01 17:22:01', '2026-02-03 12:27:52');
INSERT INTO `sys_lang_resource` VALUES (1834, 'common.button.cancel', 'default', 'ru-Ru', 'отменить', 0, '2026-02-01 17:22:01', '2026-02-03 12:27:52');
INSERT INTO `sys_lang_resource` VALUES (1835, 'common.button.clear', 'default', 'ru-Ru', 'Очистить', 0, '2026-02-01 17:22:01', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (1836, 'common.button.confirm', 'default', 'ru-Ru', 'подтверждение', 0, '2026-02-01 17:22:01', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (1837, 'common.button.copy', 'default', 'ru-Ru', 'копировать', 0, '2026-02-01 17:22:01', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (1838, 'common.button.delete', 'default', 'ru-Ru', 'удалить', 0, '2026-02-01 17:22:01', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (1839, 'common.button.download', 'default', 'ru-Ru', 'скачать', 0, '2026-02-01 17:22:01', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (1840, 'common.button.export', 'default', 'ru-Ru', 'экспорт', 0, '2026-02-01 17:22:01', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (1841, 'common.button.import', 'default', 'ru-Ru', 'импорт', 0, '2026-02-01 17:22:01', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (1842, 'common.button.modify', 'default', 'ru-Ru', 'изменить', 0, '2026-02-01 17:22:01', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (1843, 'common.button.move', 'default', 'ru-Ru', 'двигаться', 0, '2026-02-01 17:22:01', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (1844, 'common.button.refresh', 'default', 'ru-Ru', 'обновить', 0, '2026-02-01 17:22:01', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (1845, 'common.button.reset', 'default', 'ru-Ru', 'сброс', 0, '2026-02-01 17:22:01', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (1846, 'common.button.save', 'default', 'ru-Ru', 'сохранить', 0, '2026-02-01 17:22:01', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (1847, 'common.button.search', 'default', 'ru-Ru', 'поиск', 0, '2026-02-01 17:22:01', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (1848, 'common.button.setup', 'default', 'ru-Ru', 'настройка', 0, '2026-02-01 17:22:01', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (1849, 'common.button.submit', 'default', 'ru-Ru', 'отправить', 0, '2026-02-01 17:22:01', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (1850, 'common.button.update', 'default', 'ru-Ru', 'загрузить', 0, '2026-02-01 17:22:01', '2026-02-01 18:00:48');
INSERT INTO `sys_lang_resource` VALUES (2177, 'common.button.remove', 'default', 'zh-CN', '移除', 0, '2026-02-05 11:55:06', NULL);
INSERT INTO `sys_lang_resource` VALUES (2185, 'common.button.remove', 'default', 'en-US', '移除', 0, '2026-02-05 11:55:07', '2026-02-05 12:05:25');
INSERT INTO `sys_lang_resource` VALUES (2188, 'common.button.remove', 'default', 'ru-Ru', 'Удалить', 0, '2026-02-05 11:55:07', '2026-02-05 11:55:08');
INSERT INTO `sys_lang_resource` VALUES (2256, 'common.button.submit', 'default', 'en-US', 'Submit', 0, '2026-02-05 17:24:13', '2026-02-05 17:33:03');
INSERT INTO `sys_lang_resource` VALUES (2263, 'common.button.choose', 'default', 'zh-CN', '选择', 0, '2026-02-05 17:37:01', '2026-02-05 17:37:42');
INSERT INTO `sys_lang_resource` VALUES (2265, 'common.button.choose', 'default', 'en-US', 'choice', 0, '2026-02-05 17:37:01', '2026-02-05 17:37:43');
INSERT INTO `sys_lang_resource` VALUES (2266, 'common.button.choose', 'default', 'ru-Ru', 'выбор', 0, '2026-02-05 17:37:03', '2026-02-05 17:37:45');
INSERT INTO `sys_lang_resource` VALUES (2270, 'common.alert.success.update', 'default', 'zh-CN', '修改成功', 0, '2026-02-05 17:42:27', '2026-02-05 17:43:16');
INSERT INTO `sys_lang_resource` VALUES (2271, 'common.alert.success.update', 'default', 'en-US', 'Modified successfully', 0, '2026-02-05 17:43:59', NULL);
INSERT INTO `sys_lang_resource` VALUES (2272, 'common.alert.success.update', 'default', 'ru-Ru', 'Изменение прошло успешно', 0, '2026-02-05 17:44:00', NULL);
INSERT INTO `sys_lang_resource` VALUES (2273, 'common.alert.success.save', 'default', 'zh-CN', '保存成功', 0, '2026-02-05 17:54:24', '2026-02-05 17:57:10');
INSERT INTO `sys_lang_resource` VALUES (2274, 'common.alert.success.delete', 'default', 'zh-CN', '删除成功', 0, '2026-02-05 17:55:25', '2026-02-05 17:57:27');
INSERT INTO `sys_lang_resource` VALUES (2275, 'common.alert.success.add', 'default', 'zh-CN', '新增成功', 0, '2026-02-05 17:55:42', '2026-02-05 17:57:35');
INSERT INTO `sys_lang_resource` VALUES (2276, 'common.alert.success.operation', 'default', 'zh-CN', '操作成功', 0, '2026-02-05 17:56:30', '2026-02-05 17:57:18');
INSERT INTO `sys_lang_resource` VALUES (2277, 'common.alert.success.save', 'default', 'en-US', 'Saved successfully', 0, '2026-02-05 17:57:13', NULL);
INSERT INTO `sys_lang_resource` VALUES (2278, 'common.alert.success.save', 'default', 'ru-Ru', 'Сохранено успешно', 0, '2026-02-05 17:57:14', NULL);
INSERT INTO `sys_lang_resource` VALUES (2279, 'common.alert.success.operation', 'default', 'en-US', 'Operation successful', 0, '2026-02-05 17:57:20', NULL);
INSERT INTO `sys_lang_resource` VALUES (2280, 'common.alert.success.operation', 'default', 'ru-Ru', 'Операция успешна', 0, '2026-02-05 17:57:22', NULL);
INSERT INTO `sys_lang_resource` VALUES (2281, 'common.alert.success.delete', 'default', 'en-US', 'Deleted successfully', 0, '2026-02-05 17:57:29', NULL);
INSERT INTO `sys_lang_resource` VALUES (2282, 'common.alert.success.delete', 'default', 'ru-Ru', 'Удаление успешно', 0, '2026-02-05 17:57:31', NULL);
INSERT INTO `sys_lang_resource` VALUES (2283, 'common.alert.success.add', 'default', 'en-US', 'Added successfully', 0, '2026-02-05 17:57:38', NULL);
INSERT INTO `sys_lang_resource` VALUES (2284, 'common.alert.success.add', 'default', 'ru-Ru', 'Добавить успех', 0, '2026-02-05 17:57:39', NULL);
INSERT INTO `sys_lang_resource` VALUES (2285, 'common.alert.fail.update', 'default', 'zh-CN', '修改失败', 0, '2026-02-05 17:42:27', '2026-02-05 17:43:16');
INSERT INTO `sys_lang_resource` VALUES (2286, 'common.alert.fail.save', 'default', 'zh-CN', '保存失败', 0, '2026-02-05 17:54:24', '2026-02-05 17:57:10');
INSERT INTO `sys_lang_resource` VALUES (2287, 'common.alert.fail.delete', 'default', 'zh-CN', '删除失败', 0, '2026-02-05 17:55:25', '2026-02-05 17:57:27');
INSERT INTO `sys_lang_resource` VALUES (2288, 'common.alert.fail.add', 'default', 'zh-CN', '新增失败', 0, '2026-02-05 17:55:42', '2026-02-05 17:57:35');
INSERT INTO `sys_lang_resource` VALUES (2289, 'common.alert.fail.operation', 'default', 'zh-CN', '操作失败', 0, '2026-02-05 17:56:30', '2026-02-05 17:57:18');
INSERT INTO `sys_lang_resource` VALUES (2290, 'common.alert.fail.update', 'default', 'en-US', 'Modification failed', 0, '2026-02-05 18:03:07', NULL);
INSERT INTO `sys_lang_resource` VALUES (2291, 'common.alert.fail.update', 'default', 'ru-Ru', 'Ошибка изменения', 0, '2026-02-05 18:03:09', NULL);
INSERT INTO `sys_lang_resource` VALUES (2292, 'common.alert.fail.save', 'default', 'en-US', 'Save failed', 0, '2026-02-05 18:03:49', NULL);
INSERT INTO `sys_lang_resource` VALUES (2293, 'common.alert.fail.save', 'default', 'ru-Ru', 'Не удалось сохранить', 0, '2026-02-05 18:03:50', NULL);
INSERT INTO `sys_lang_resource` VALUES (2294, 'common.alert.fail.operation', 'default', 'en-US', 'Operation failed', 0, '2026-02-05 18:04:04', NULL);
INSERT INTO `sys_lang_resource` VALUES (2295, 'common.alert.fail.operation', 'default', 'ru-Ru', 'Ошибка операции', 0, '2026-02-05 18:04:05', NULL);
INSERT INTO `sys_lang_resource` VALUES (2296, 'common.alert.fail.delete', 'default', 'en-US', 'Deletion failed', 0, '2026-02-05 18:04:12', NULL);
INSERT INTO `sys_lang_resource` VALUES (2297, 'common.alert.fail.delete', 'default', 'ru-Ru', 'Ошибка удаления', 0, '2026-02-05 18:04:13', NULL);
INSERT INTO `sys_lang_resource` VALUES (2298, 'common.alert.fail.add', 'default', 'en-US', 'Failed to add', 0, '2026-02-05 18:04:19', NULL);
INSERT INTO `sys_lang_resource` VALUES (2299, 'common.alert.fail.add', 'default', 'ru-Ru', 'Новые неудачи', 0, '2026-02-05 18:04:21', NULL);
INSERT INTO `sys_lang_resource` VALUES (2300, 'authentication.unauthorized', 'default', 'zh-CN', '认证失败', 0, '2026-02-06 10:37:38', '2026-02-06 10:38:23');
INSERT INTO `sys_lang_resource` VALUES (2301, 'authentication.username_not_found', 'default', 'zh-CN', '用户名或密码不正确', 0, '2026-02-06 10:37:38', '2026-02-06 10:38:12');
INSERT INTO `sys_lang_resource` VALUES (2302, 'authentication.bad_credentials', 'default', 'zh-CN', '用户名或密码不正确', 0, '2026-02-06 10:37:38', '2026-02-06 10:39:20');
INSERT INTO `sys_lang_resource` VALUES (2303, 'authentication.locked', 'default', 'zh-CN', '用户已被锁定，请联系管理员', 0, '2026-02-06 10:37:38', '2026-02-06 10:38:47');
INSERT INTO `sys_lang_resource` VALUES (2304, 'authentication.disabled', 'default', 'zh-CN', '用户已失效', 0, '2026-02-06 10:37:38', '2026-02-06 10:38:56');
INSERT INTO `sys_lang_resource` VALUES (2305, 'authentication.account_expired', 'default', 'zh-CN', '用户已过期', 0, '2026-02-06 10:37:38', '2026-02-06 10:39:27');
INSERT INTO `sys_lang_resource` VALUES (2306, 'authentication.credentials_expired', 'default', 'zh-CN', '凭证已过期', 0, '2026-02-06 10:37:38', '2026-02-06 10:39:12');
INSERT INTO `sys_lang_resource` VALUES (2307, 'authentication.denied_permission', 'default', 'zh-CN', '没有权限:[{0}]', 0, '2026-02-06 10:37:38', '2026-02-06 10:39:04');
INSERT INTO `sys_lang_resource` VALUES (2308, 'authentication.token_invalid', 'default', 'zh-CN', 'Token 无效', 0, '2026-02-06 10:37:38', '2026-02-06 10:38:39');
INSERT INTO `sys_lang_resource` VALUES (2309, 'authentication.token_valid', 'default', 'zh-CN', 'Token 有效', 0, '2026-02-06 10:37:38', '2026-02-06 10:38:31');
INSERT INTO `sys_lang_resource` VALUES (2310, 'exception.access_denied', 'default', 'zh-CN', '没有权限，请联系管理员授权', 0, '2026-02-06 10:37:38', '2026-02-06 10:40:13');
INSERT INTO `sys_lang_resource` VALUES (2311, 'exception.request_method_not_supported', 'default', 'zh-CN', '请求地址{0}，不支持{1}请求', 0, '2026-02-06 10:37:38', '2026-02-06 10:39:44');
INSERT INTO `sys_lang_resource` VALUES (2312, 'exception.missing_path_variable', 'default', 'zh-CN', '请求路径中缺少必需的路径变量{0}', 0, '2026-02-06 10:37:38', '2026-02-06 10:39:52');
INSERT INTO `sys_lang_resource` VALUES (2313, 'exception.method_argument_type_mismatch', 'default', 'zh-CN', '请求参数 [{0}] 类型不匹配', 0, '2026-02-06 10:37:38', '2026-02-06 10:39:59');
INSERT INTO `sys_lang_resource` VALUES (2314, 'exception.demo_mode', 'default', 'zh-CN', '演示模式，不能执行此操作', 0, '2026-02-06 10:37:38', '2026-02-06 10:40:06');
INSERT INTO `sys_lang_resource` VALUES (2315, 'authentication.username_not_found', 'default', 'en-US', 'Username or password incorrect ', 0, '2026-02-06 10:38:14', NULL);
INSERT INTO `sys_lang_resource` VALUES (2316, 'authentication.username_not_found', 'default', 'ru-Ru', 'Неверное имя пользователя или пароль', 0, '2026-02-06 10:38:16', NULL);
INSERT INTO `sys_lang_resource` VALUES (2317, 'authentication.unauthorized', 'default', 'en-US', 'Authentication failed', 0, '2026-02-06 10:38:26', NULL);
INSERT INTO `sys_lang_resource` VALUES (2318, 'authentication.unauthorized', 'default', 'ru-Ru', 'Ошибка аутентификации', 0, '2026-02-06 10:38:27', NULL);
INSERT INTO `sys_lang_resource` VALUES (2319, 'authentication.token_valid', 'default', 'en-US', 'Token is valid', 0, '2026-02-06 10:38:33', NULL);
INSERT INTO `sys_lang_resource` VALUES (2320, 'authentication.token_valid', 'default', 'ru-Ru', 'Токен работает', 0, '2026-02-06 10:38:35', NULL);
INSERT INTO `sys_lang_resource` VALUES (2321, 'authentication.token_invalid', 'default', 'en-US', 'Invalid token', 0, '2026-02-06 10:38:41', NULL);
INSERT INTO `sys_lang_resource` VALUES (2322, 'authentication.token_invalid', 'default', 'ru-Ru', 'Токен не работает', 0, '2026-02-06 10:38:43', NULL);
INSERT INTO `sys_lang_resource` VALUES (2323, 'authentication.locked', 'default', 'en-US', 'The user has been locked, please contact the administrator', 0, '2026-02-06 10:38:50', NULL);
INSERT INTO `sys_lang_resource` VALUES (2324, 'authentication.locked', 'default', 'ru-Ru', 'Пользователь заблокирован, свяжитесь с администратором', 0, '2026-02-06 10:38:51', NULL);
INSERT INTO `sys_lang_resource` VALUES (2325, 'authentication.disabled', 'default', 'en-US', 'User has expired', 0, '2026-02-06 10:38:58', NULL);
INSERT INTO `sys_lang_resource` VALUES (2326, 'authentication.disabled', 'default', 'ru-Ru', 'Пользователь утратил силу', 0, '2026-02-06 10:39:00', NULL);
INSERT INTO `sys_lang_resource` VALUES (2327, 'authentication.denied_permission', 'default', 'en-US', 'No permission: [{0}]', 0, '2026-02-06 10:39:06', NULL);
INSERT INTO `sys_lang_resource` VALUES (2328, 'authentication.denied_permission', 'default', 'ru-Ru', 'Нет прав: [0]', 0, '2026-02-06 10:39:08', NULL);
INSERT INTO `sys_lang_resource` VALUES (2329, 'authentication.credentials_expired', 'default', 'en-US', 'The voucher has expired', 0, '2026-02-06 10:39:15', NULL);
INSERT INTO `sys_lang_resource` VALUES (2330, 'authentication.credentials_expired', 'default', 'ru-Ru', 'Срок действия документов истек', 0, '2026-02-06 10:39:16', NULL);
INSERT INTO `sys_lang_resource` VALUES (2331, 'authentication.bad_credentials', 'default', 'en-US', 'Username or password incorrect ', 0, '2026-02-06 10:39:22', NULL);
INSERT INTO `sys_lang_resource` VALUES (2332, 'authentication.bad_credentials', 'default', 'ru-Ru', 'Неверное имя пользователя или пароль', 0, '2026-02-06 10:39:24', NULL);
INSERT INTO `sys_lang_resource` VALUES (2333, 'authentication.account_expired', 'default', 'en-US', 'User has expired', 0, '2026-02-06 10:39:29', NULL);
INSERT INTO `sys_lang_resource` VALUES (2334, 'authentication.account_expired', 'default', 'ru-Ru', 'Пользователь просрочен.', 0, '2026-02-06 10:39:31', NULL);
INSERT INTO `sys_lang_resource` VALUES (2335, 'exception.request_method_not_supported', 'default', 'en-US', 'Request address {0}, does not support {1} requests', 0, '2026-02-06 10:39:46', NULL);
INSERT INTO `sys_lang_resource` VALUES (2336, 'exception.request_method_not_supported', 'default', 'ru-Ru', 'Адрес запроса {0}, запрос {1} не поддерживается', 0, '2026-02-06 10:39:48', NULL);
INSERT INTO `sys_lang_resource` VALUES (2337, 'exception.missing_path_variable', 'default', 'en-US', 'The required path variable {0} is missing from the request path', 0, '2026-02-06 10:39:54', NULL);
INSERT INTO `sys_lang_resource` VALUES (2338, 'exception.missing_path_variable', 'default', 'ru-Ru', 'Отсутствует необходимая переменная пути {0}', 0, '2026-02-06 10:39:56', NULL);
INSERT INTO `sys_lang_resource` VALUES (2339, 'exception.method_argument_type_mismatch', 'default', 'en-US', 'Request parameter [{0}] type mismatch', 0, '2026-02-06 10:40:02', NULL);
INSERT INTO `sys_lang_resource` VALUES (2340, 'exception.method_argument_type_mismatch', 'default', 'ru-Ru', 'Параметры запроса [0] тип не совпадает', 0, '2026-02-06 10:40:03', NULL);
INSERT INTO `sys_lang_resource` VALUES (2341, 'exception.demo_mode', 'default', 'en-US', 'Demo mode, this operation cannot be performed', 0, '2026-02-06 10:40:09', NULL);
INSERT INTO `sys_lang_resource` VALUES (2342, 'exception.demo_mode', 'default', 'ru-Ru', 'Демонстрация режима, эта операция не может быть выполнена', 0, '2026-02-06 10:40:10', NULL);
INSERT INTO `sys_lang_resource` VALUES (2343, 'exception.access_denied', 'default', 'en-US', 'No permission, please contact the administrator for authorization', 0, '2026-02-06 10:40:16', NULL);
INSERT INTO `sys_lang_resource` VALUES (2344, 'exception.access_denied', 'default', 'ru-Ru', 'Без прав, свяжитесь с администратором для авторизации', 0, '2026-02-06 10:40:17', NULL);
INSERT INTO `sys_lang_resource` VALUES (2345, 'exception.bad_sql_grammar', 'default', 'zh-CN', '数据表已被更新，部分字段或表不存在', 0, '2026-02-06 10:42:16', '2026-02-06 10:42:26');
INSERT INTO `sys_lang_resource` VALUES (2346, 'exception.bad_sql_grammar', 'default', 'en-US', 'The data table has been updated, and some fields or tables do not exist', 0, '2026-02-06 10:42:28', NULL);
INSERT INTO `sys_lang_resource` VALUES (2347, 'exception.bad_sql_grammar', 'default', 'ru-Ru', 'Таблица данных обновлена, некоторые поля или таблицы не существуют', 0, '2026-02-06 10:42:30', NULL);
INSERT INTO `sys_lang_resource` VALUES (2360, 'template.test', 'template', 'zh-CN', '测试', 0, '2026-04-21 16:57:48', '2026-04-21 16:57:48');
INSERT INTO `sys_lang_resource` VALUES (2361, 'template.test', 'template', 'en-US', 'Test', 0, '2026-04-21 16:57:48', '2026-04-21 16:57:48');
INSERT INTO `sys_lang_resource` VALUES (2362, 'template.test', 'template', 'ru-Ru', '', 0, '2026-04-21 16:57:48', '2026-04-21 16:57:48');
COMMIT;

-- ----------------------------
-- Table structure for sys_login_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_login_log`;
CREATE TABLE `sys_login_log` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '编号',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '登录用户名',
  `client_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '客户端ID',
  `grant_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '授权类型',
  `os` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '系统',
  `browser` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '浏览器',
  `ip` varchar(128) DEFAULT NULL COMMENT '登录IP',
  `location` varchar(255) DEFAULT NULL COMMENT '登录地点',
  `status` tinyint(1) DEFAULT NULL COMMENT '状态（0失败 1成功）',
  `msg` varchar(50) DEFAULT NULL COMMENT '模块标题',
  `login_at` datetime DEFAULT NULL COMMENT '登录时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=259 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='操作日志记录';

-- ----------------------------
-- Records of sys_login_log
-- ----------------------------
BEGIN;
INSERT INTO `sys_login_log` VALUES (215, 'develop', 'auth-demo', 'authorization_code', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-22 12:35:05');
INSERT INTO `sys_login_log` VALUES (216, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-22 14:48:31');
INSERT INTO `sys_login_log` VALUES (217, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 0, 'auth.invalid_credentials', '2026-04-22 16:51:21');
INSERT INTO `sys_login_log` VALUES (218, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-22 16:51:25');
INSERT INTO `sys_login_log` VALUES (219, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-23 10:33:09');
INSERT INTO `sys_login_log` VALUES (220, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-23 10:46:21');
INSERT INTO `sys_login_log` VALUES (221, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-23 16:25:10');
INSERT INTO `sys_login_log` VALUES (222, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-24 14:42:19');
INSERT INTO `sys_login_log` VALUES (223, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-25 14:03:15');
INSERT INTO `sys_login_log` VALUES (224, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-25 16:46:22');
INSERT INTO `sys_login_log` VALUES (225, 'develop', 'auth-demo', 'authorization_code', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-25 16:57:07');
INSERT INTO `sys_login_log` VALUES (226, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-25 16:57:41');
INSERT INTO `sys_login_log` VALUES (227, 'develop', 'auth-demo', 'authorization_code', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-25 16:57:47');
INSERT INTO `sys_login_log` VALUES (228, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-25 17:09:46');
INSERT INTO `sys_login_log` VALUES (229, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-25 20:55:17');
INSERT INTO `sys_login_log` VALUES (230, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-25 20:56:30');
INSERT INTO `sys_login_log` VALUES (231, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-25 20:59:47');
INSERT INTO `sys_login_log` VALUES (232, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-25 21:19:03');
INSERT INTO `sys_login_log` VALUES (233, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-26 15:27:36');
INSERT INTO `sys_login_log` VALUES (234, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-26 15:47:55');
INSERT INTO `sys_login_log` VALUES (235, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-26 15:53:09');
INSERT INTO `sys_login_log` VALUES (236, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-26 16:01:18');
INSERT INTO `sys_login_log` VALUES (237, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-26 16:08:02');
INSERT INTO `sys_login_log` VALUES (238, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-26 16:13:35');
INSERT INTO `sys_login_log` VALUES (239, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-26 16:18:29');
INSERT INTO `sys_login_log` VALUES (240, 'develop', 'passkey', 'passkey', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Passkey login successful', '2026-04-26 17:10:00');
INSERT INTO `sys_login_log` VALUES (241, 'develop', 'passkey', 'passkey', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Passkey login successful', '2026-04-26 17:12:06');
INSERT INTO `sys_login_log` VALUES (242, 'develop', 'auth-demo', 'authorization_code', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-26 17:12:08');
INSERT INTO `sys_login_log` VALUES (243, 'develop', 'passkey', 'passkey', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Passkey login successful', '2026-04-27 11:35:12');
INSERT INTO `sys_login_log` VALUES (244, 'admin', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-27 14:24:52');
INSERT INTO `sys_login_log` VALUES (245, 'develop', 'passkey', 'passkey', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Passkey login successful', '2026-04-27 14:26:04');
INSERT INTO `sys_login_log` VALUES (246, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-28 10:30:11');
INSERT INTO `sys_login_log` VALUES (247, 'admin', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-28 11:00:10');
INSERT INTO `sys_login_log` VALUES (248, 'develop', 'passkey', 'passkey', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Passkey login successful', '2026-04-28 18:17:21');
INSERT INTO `sys_login_log` VALUES (249, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-28 18:17:37');
INSERT INTO `sys_login_log` VALUES (250, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-29 12:27:20');
INSERT INTO `sys_login_log` VALUES (251, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-04-30 10:22:45');
INSERT INTO `sys_login_log` VALUES (252, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-05-01 10:46:11');
INSERT INTO `sys_login_log` VALUES (253, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-05-01 10:57:17');
INSERT INTO `sys_login_log` VALUES (254, 'develop', '100000', 'password', 'Linux', 'Chrome', '127.0.0.1', '', 1, 'Login successful', '2026-05-03 12:45:57');
INSERT INTO `sys_login_log` VALUES (255, 'develop', '100000', 'passkey', 'Linux', 'Chrome', '192.168.0.102', '', 1, 'Passkey login successful', '2026-05-18 17:37:52');
INSERT INTO `sys_login_log` VALUES (256, 'develop', '100000', 'password', 'Linux', 'Chrome', '192.168.0.102', '', 1, 'Login successful', '2026-05-18 17:38:13');
INSERT INTO `sys_login_log` VALUES (257, '123123', '100000', 'password', 'Linux', 'Chrome', '192.168.0.102', '', 0, 'auth.invalid_credentials', '2026-05-18 18:01:39');
INSERT INTO `sys_login_log` VALUES (258, 'develop', '100000', 'passkey', 'Linux', 'Chrome', '192.168.0.102', '', 1, 'Passkey login successful', '2026-05-18 18:02:29');
COMMIT;

-- ----------------------------
-- Table structure for sys_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_menu`;
CREATE TABLE `sys_menu` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `parent_id` bigint DEFAULT NULL COMMENT '父级菜单Id',
  `type` tinyint(1) DEFAULT NULL COMMENT '类型（1目录 2菜单 3按钮）',
  `sort` int DEFAULT NULL COMMENT '顺序',
  `path` varchar(100) DEFAULT NULL COMMENT '路由地址',
  `component` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '组件路径',
  `query` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '路由参数',
  `visible` tinyint(1) DEFAULT NULL COMMENT '是否可见（0隐藏 1显示）',
  `status` tinyint(1) DEFAULT NULL COMMENT '状态（0正常 1停用）',
  `is_frame` tinyint(1) DEFAULT NULL COMMENT '是否为外链（0否 1是）',
  `permission` varchar(32) DEFAULT NULL COMMENT '权限标识',
  `icon` varchar(64) DEFAULT NULL COMMENT '图标',
  `active_id` bigint DEFAULT NULL COMMENT '激活ID',
  `created_by` bigint DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_by` bigint DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=100298 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='菜单权限表';

-- ----------------------------
-- Records of sys_menu
-- ----------------------------
BEGIN;
INSERT INTO `sys_menu` VALUES (100000, 0, 1, 99, '', '', '', 1, 0, 0, '', 'module', 0, 1, '2025-01-25 16:05:39', NULL, '2026-04-19 16:05:50');
INSERT INTO `sys_menu` VALUES (100002, 100294, 2, 1, '/system/user', 'system/user/index', '', 1, 0, 0, 'sys:user$list', 'user', 0, 1, '2025-01-25 16:05:39', NULL, '2026-04-23 16:02:19');
INSERT INTO `sys_menu` VALUES (100003, 100002, 3, 1, '', '', '', 1, 0, 0, 'sys:user$query', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-19 09:26:49');
INSERT INTO `sys_menu` VALUES (100004, 100002, 3, 2, '', '', '', 1, 0, 0, 'sys:user$add', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-19 09:26:58');
INSERT INTO `sys_menu` VALUES (100005, 100002, 3, 3, '', '', '', 1, 0, 0, 'sys:user$edit', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-19 09:27:05');
INSERT INTO `sys_menu` VALUES (100006, 100002, 3, 4, '', '', '', 1, 0, 0, 'sys:user$delete', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-19 09:27:16');
INSERT INTO `sys_menu` VALUES (100007, 100002, 3, 5, '', '', '', 1, 0, 0, 'sys:user$export', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-19 09:27:26');
INSERT INTO `sys_menu` VALUES (100008, 100002, 3, 6, '', '', '', 1, 0, 0, 'sys:user$import', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-19 09:27:36');
INSERT INTO `sys_menu` VALUES (100009, 100002, 3, 7, '', '', '', 1, 0, 0, 'sys:user$resetPwd', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-19 09:27:55');
INSERT INTO `sys_menu` VALUES (100010, 100002, 3, 8, '', '', '', 1, 0, 0, 'sys:user$updatePwd', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-19 09:28:26');
INSERT INTO `sys_menu` VALUES (100011, 100294, 2, 2, '/system/role', 'system/role/index', '', 1, 0, 0, 'sys:role$list', 'role', 0, 1, '2025-01-25 16:05:39', NULL, '2026-04-24 15:13:27');
INSERT INTO `sys_menu` VALUES (100012, 100011, 3, 1, '', '', '', 1, 0, 0, 'sys:role$query', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:25:00');
INSERT INTO `sys_menu` VALUES (100013, 100011, 3, 2, '', '', '', 1, 0, 0, 'sys:role$add', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:25:06');
INSERT INTO `sys_menu` VALUES (100014, 100011, 3, 3, '', '', '', 1, 0, 0, 'sys:role$edit', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:25:14');
INSERT INTO `sys_menu` VALUES (100015, 100011, 3, 4, '', '', '', 1, 0, 0, 'sys:role$delete', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:25:54');
INSERT INTO `sys_menu` VALUES (100016, 100011, 3, 5, '', '', '', 1, 0, 0, 'sys:role$export', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:25:44');
INSERT INTO `sys_menu` VALUES (100017, 100294, 2, 4, '/system/dept', 'system/dept/index', '', 1, 0, 0, 'sys:dept$list', 'tree', 0, 1, '2025-01-25 16:05:39', NULL, '2026-04-23 16:02:39');
INSERT INTO `sys_menu` VALUES (100018, 100017, 3, 1, '', '', '', 1, 0, 0, 'sys:dept$query', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 20:21:50');
INSERT INTO `sys_menu` VALUES (100019, 100017, 3, 2, '', '', '', 1, 0, 0, 'sys:dept$add', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 20:21:41');
INSERT INTO `sys_menu` VALUES (100020, 100017, 3, 3, '', '', '', 1, 0, 0, 'sys:dept$edit', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 20:21:31');
INSERT INTO `sys_menu` VALUES (100021, 100017, 3, 4, '', '', '', 1, 0, 0, 'sys:dept$delete', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 20:21:22');
INSERT INTO `sys_menu` VALUES (100022, 100294, 2, 5, '/system/post', 'system/post/index', '', 1, 0, 0, 'sys:post$list', 'post', 0, 1, '2025-01-25 16:05:39', NULL, '2026-04-23 16:06:18');
INSERT INTO `sys_menu` VALUES (100023, 100022, 3, 1, '', '', '', 1, 0, 0, 'sys:post$query', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 20:20:58');
INSERT INTO `sys_menu` VALUES (100024, 100022, 3, 2, '', '', '', 1, 0, 0, 'sys:post$add', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 20:20:45');
INSERT INTO `sys_menu` VALUES (100025, 100022, 3, 3, '', '', '', 1, 0, 0, 'sys:post$edit', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 20:20:39');
INSERT INTO `sys_menu` VALUES (100026, 100022, 3, 4, '', '', '', 1, 0, 0, 'sys:post$delete', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 20:20:21');
INSERT INTO `sys_menu` VALUES (100027, 100022, 3, 5, '', '', '', 1, 0, 0, 'sys:post$export', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 20:20:31');
INSERT INTO `sys_menu` VALUES (100053, 100295, 2, 3, '/system/menu', 'system/menu/index', '', 1, 0, 0, 'sys:menu$list', 'tree-table', 0, 1, '2025-01-25 16:05:39', NULL, '2026-04-23 16:03:08');
INSERT INTO `sys_menu` VALUES (100054, 100053, 3, 1, '', '', '', 1, 0, 0, 'sys:menu$query', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:23:58');
INSERT INTO `sys_menu` VALUES (100055, 100053, 3, 2, '', '', '', 1, 0, 0, 'sys:menu$add', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:24:06');
INSERT INTO `sys_menu` VALUES (100056, 100053, 3, 3, '', '', '', 1, 0, 0, 'sys:menu$edit', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:24:15');
INSERT INTO `sys_menu` VALUES (100057, 100053, 3, 4, '', '', '', 1, 0, 0, 'sys:menu$delete', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:24:23');
INSERT INTO `sys_menu` VALUES (100058, 100295, 2, 6, '/system/dict', 'system/dict/index', '', 1, 0, 0, 'sys:dict$list', 'dict', 0, 1, '2025-01-25 16:05:39', NULL, '2026-04-23 16:03:16');
INSERT INTO `sys_menu` VALUES (100059, 100058, 3, 1, '', '', '', 1, 0, 0, 'sys:dict$query', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 19:47:44');
INSERT INTO `sys_menu` VALUES (100060, 100058, 3, 3, '', '', '', 1, 0, 0, 'sys:dict$edit', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 19:48:00');
INSERT INTO `sys_menu` VALUES (100061, 100058, 3, 4, '', '', '', 1, 0, 0, 'sys:dict$delete', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 19:48:07');
INSERT INTO `sys_menu` VALUES (100062, 100058, 3, 5, '', '', '', 1, 0, 0, 'sys:dict$export', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 19:48:14');
INSERT INTO `sys_menu` VALUES (100063, 100295, 2, 7, '/system/config', 'system/config/index', '', 1, 0, 0, 'sys:conf$list', 'set-up', 0, 1, '2025-01-25 16:05:39', NULL, '2026-04-23 16:03:23');
INSERT INTO `sys_menu` VALUES (100064, 100063, 3, 1, '', '', '', 1, 0, 0, 'sys:conf$query', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:26:36');
INSERT INTO `sys_menu` VALUES (100065, 100063, 3, 2, '', '', '', 1, 0, 0, 'sys:conf$add', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:26:42');
INSERT INTO `sys_menu` VALUES (100066, 100063, 3, 3, '', '', '', 1, 0, 0, 'sys:conf$edit', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:26:50');
INSERT INTO `sys_menu` VALUES (100067, 100063, 3, 4, '', '', '', 1, 0, 0, 'sys:conf$delete', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:27:03');
INSERT INTO `sys_menu` VALUES (100068, 100063, 3, 5, '', '', '', 1, 0, 0, 'sys:conf$export', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:27:18');
INSERT INTO `sys_menu` VALUES (100094, 100000, 1, 98, '/system/monitor/online', '', '', 1, 0, 0, '', 'monitor', 0, 1, '2025-01-25 16:05:39', NULL, '2026-04-24 16:32:06');
INSERT INTO `sys_menu` VALUES (100095, 100094, 2, 1, '/system/monitor/online', 'monitor/online/index', '', 1, 0, 0, 'monitor:online$list', 'online', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 21:37:03');
INSERT INTO `sys_menu` VALUES (100097, 100095, 3, 2, '', '', '', 1, 0, 0, 'monitor:online$forceLogout', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 22:47:39');
INSERT INTO `sys_menu` VALUES (100099, 100094, 2, 4, '/system/monitor/server', 'monitor/server/index', '', 1, 0, 0, 'monitor:server$query', 'server', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 21:36:44');
INSERT INTO `sys_menu` VALUES (100100, 100094, 2, 7, '/system/monitor/oper_log', 'monitor/oper_log/index', '', 1, 0, 0, 'monitor:oper:log$query', 'form', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 22:43:10');
INSERT INTO `sys_menu` VALUES (100102, 100100, 3, 2, '', '', '', 1, 0, 0, 'monitor:oper:log$delete', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 22:44:34');
INSERT INTO `sys_menu` VALUES (100104, 100094, 2, 8, '/system/monitor/login_log', 'monitor/login_log/index', '', 1, 0, 0, 'monitor:login:log$query', 'logininfor', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 22:46:32');
INSERT INTO `sys_menu` VALUES (100105, 100104, 3, 1, '', '', '', 1, 0, 0, 'monitor:login:log$clear', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 22:47:28');
INSERT INTO `sys_menu` VALUES (100108, 100000, 1, 99, '/system/tool', '', '', 1, 0, 0, '', 'tool', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 12:00:34');
INSERT INTO `sys_menu` VALUES (100109, 100108, 2, 1, '/system/tool/job', 'tool/job/index', '', 1, 0, 0, 'job$list', 'job', 0, 1, '2025-01-25 16:05:39', 0, '2026-04-19 16:33:49');
INSERT INTO `sys_menu` VALUES (100110, 100109, 3, 1, '', '', '', 1, 0, 0, 'job$query', '', 0, 1, '2025-01-25 16:05:39', NULL, '2026-04-19 12:37:28');
INSERT INTO `sys_menu` VALUES (100111, 100109, 3, 2, '', '', '', 1, 0, 0, 'job$add', '', 0, 1, '2025-01-25 16:05:39', NULL, '2026-04-19 12:37:40');
INSERT INTO `sys_menu` VALUES (100112, 100109, 3, 3, '', '', '', 1, 0, 0, 'job$edit', '', 0, 1, '2025-01-25 16:05:39', NULL, '2026-04-19 12:37:49');
INSERT INTO `sys_menu` VALUES (100113, 100109, 3, 4, '', '', '', 1, 0, 0, 'job$delete', '', 0, 1, '2025-01-25 16:05:39', NULL, '2026-04-19 12:37:59');
INSERT INTO `sys_menu` VALUES (100119, 100108, 2, 3, '/system/tool/cache', 'tool/cache/index', '', 1, 0, 0, 'cache$list', 'redis', 0, 1, '2025-01-25 16:05:39', 0, '2026-04-22 12:50:33');
INSERT INTO `sys_menu` VALUES (100120, 100108, 2, 4, '/system/tool/gen', 'tool/gen/index', '', 1, 0, 0, 'gen:table$list', 'gencode', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 15:54:28');
INSERT INTO `sys_menu` VALUES (100121, 100120, 3, 1, '', '', '', 1, 0, 0, 'gen:table$query', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:13:07');
INSERT INTO `sys_menu` VALUES (100122, 100120, 3, 2, '', '', '', 1, 0, 0, 'gen:table$edit', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:11:44');
INSERT INTO `sys_menu` VALUES (100123, 100120, 3, 3, '', '', '', 1, 0, 0, 'gen:table$delete', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:11:59');
INSERT INTO `sys_menu` VALUES (100124, 100120, 3, 4, '', '', '', 1, 0, 0, 'gen:table$import', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:12:22');
INSERT INTO `sys_menu` VALUES (100125, 100120, 3, 5, '', '', '', 1, 0, 0, 'gen:table$preview', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:13:16');
INSERT INTO `sys_menu` VALUES (100126, 100120, 3, 6, '', '', '', 1, 0, 0, 'gen:table$download', '', NULL, 1, '2025-01-25 16:05:39', NULL, '2026-04-18 16:12:58');
INSERT INTO `sys_menu` VALUES (100183, 100109, 3, 7, '', '', '', 1, 0, 0, 'job:log$delete', '', 0, 1, '2025-11-19 11:12:46', NULL, '2026-04-20 15:01:10');
INSERT INTO `sys_menu` VALUES (100184, 100109, 3, 6, '', '', '', 1, 0, 0, 'job:log$list', '', 0, 1, '2025-11-19 11:13:04', NULL, '2026-04-20 15:01:05');
INSERT INTO `sys_menu` VALUES (100185, 100119, 3, 1, '', '', '', 1, 0, 0, 'cache$delete', '', 0, 1, '2025-11-19 11:20:08', NULL, '2026-04-19 12:37:15');
INSERT INTO `sys_menu` VALUES (100221, 100000, 1, 30, '/system/locale', '', '', 1, 0, 0, '', 'language', NULL, 1, '2026-01-28 09:04:20', 1, '2026-04-18 13:18:11');
INSERT INTO `sys_menu` VALUES (100249, 100294, 2, 10, '/system/oauth_client', 'system/oauth_client/index', '', 1, 0, 0, 'sys:oauth:client$list', 'app', 0, 1, '2026-04-16 17:04:02', 1, '2026-04-23 16:15:08');
INSERT INTO `sys_menu` VALUES (100250, 100249, 3, 1, '', '', '', 1, 0, 0, 'sys:oauth:client$query', '', NULL, 1, '2026-04-16 17:04:02', 1, '2026-04-18 20:59:28');
INSERT INTO `sys_menu` VALUES (100251, 100249, 3, 2, '', '', '', 1, 0, 0, 'sys:oauth:client$add', '', NULL, 1, '2026-04-16 17:04:02', 1, '2026-04-18 19:55:02');
INSERT INTO `sys_menu` VALUES (100252, 100249, 3, 3, '', '', '', 1, 0, 0, 'sys:oauth:client$edit', '', NULL, 1, '2026-04-16 17:04:02', 1, '2026-04-18 19:55:09');
INSERT INTO `sys_menu` VALUES (100253, 100249, 3, 4, '', '', '', 1, 0, 0, 'sys:oauth:client$delete', '', NULL, 1, '2026-04-16 17:04:02', 1, '2026-04-18 19:55:17');
INSERT INTO `sys_menu` VALUES (100254, 100249, 3, 5, '', '', '', 1, 0, 0, 'sys:oauth:client$export', '', NULL, 1, '2026-04-16 17:04:02', 1, '2026-04-18 19:55:24');
INSERT INTO `sys_menu` VALUES (100261, 100120, 3, 1, '', '', '', 1, 0, 0, 'gen:table$design', '', NULL, 0, '2026-04-18 16:09:10', 0, '2026-04-18 16:09:42');
INSERT INTO `sys_menu` VALUES (100263, 100058, 3, 7, '', '', '', 1, 0, 0, 'sys:dict$design', '', NULL, 1, '2026-04-18 16:34:09', 1, '2026-04-18 19:48:31');
INSERT INTO `sys_menu` VALUES (100267, 100058, 3, 2, '', '', '', 1, 0, 0, 'sys:dict$add', '', NULL, 0, '2026-04-18 19:34:09', 0, '2026-04-18 19:47:53');
INSERT INTO `sys_menu` VALUES (100268, 100058, 3, 6, '', '', '', 1, 0, 0, 'sys:dict$refresh', '', NULL, 0, '2026-04-18 19:34:51', 0, '2026-04-18 19:48:23');
INSERT INTO `sys_menu` VALUES (100269, 100017, 3, 5, '', '', '', 1, 0, 0, 'sys:dept$import', '', NULL, 0, '2026-04-18 20:22:33', 0, '2026-04-18 20:22:44');
INSERT INTO `sys_menu` VALUES (100270, 100017, 3, 6, '', '', '', 1, 0, 0, 'sys:dept$export', '', NULL, 0, '2026-04-18 20:23:21', 0, '2026-04-18 20:23:21');
INSERT INTO `sys_menu` VALUES (100271, 100100, 3, 1, '', '', '', 1, 0, 0, 'monitor:oper:log$clear', '', NULL, 0, '2026-04-18 22:44:28', 0, '2026-04-18 22:44:28');
INSERT INTO `sys_menu` VALUES (100272, 100104, 3, 1, '', '', '', 1, 0, 0, 'monitor:login:log$delete', '', NULL, 0, '2026-04-19 08:45:00', 0, '2026-04-19 08:45:00');
INSERT INTO `sys_menu` VALUES (100278, 100221, 2, 1, '/system/locale', 'system/language/index', '', 1, 0, 0, 'sys:lang$list', 'language', 0, 1, '2026-04-19 12:58:24', 1, '2026-04-19 13:01:58');
INSERT INTO `sys_menu` VALUES (100279, 100278, 3, 1, '', '', '', 1, 0, 0, 'sys:lang$query', '', NULL, 1, '2026-04-19 12:58:24', 1, '2026-04-19 12:58:24');
INSERT INTO `sys_menu` VALUES (100280, 100278, 3, 2, '', '', '', 1, 0, 0, 'sys:lang$add', '', NULL, 1, '2026-04-19 12:58:24', 1, '2026-04-19 12:58:24');
INSERT INTO `sys_menu` VALUES (100281, 100278, 3, 3, '', '', '', 1, 0, 0, 'sys:lang$edit', '', NULL, 1, '2026-04-19 12:58:24', 1, '2026-04-19 12:58:24');
INSERT INTO `sys_menu` VALUES (100282, 100278, 3, 4, '', '', '', 1, 0, 0, 'sys:lang$delete', '', NULL, 1, '2026-04-19 12:58:24', 1, '2026-04-19 12:58:24');
INSERT INTO `sys_menu` VALUES (100283, 100221, 2, 2, '/system/locale/resource', 'system/language/resource/index', '', 1, 0, 0, 'sys:lang:resource$list', 'resource', 0, 1, '2026-04-19 12:58:47', 1, '2026-04-24 15:42:18');
INSERT INTO `sys_menu` VALUES (100284, 100283, 3, 1, '', '', '', 1, 0, 0, 'sys:lang:resource$query', '', NULL, 1, '2026-04-19 12:58:47', 1, '2026-04-19 12:58:47');
INSERT INTO `sys_menu` VALUES (100285, 100283, 3, 2, '', '', '', 1, 0, 0, 'sys:lang:resource$add', '', NULL, 1, '2026-04-19 12:58:47', 1, '2026-04-19 12:58:47');
INSERT INTO `sys_menu` VALUES (100286, 100283, 3, 3, '', '', '', 1, 0, 0, 'sys:lang:resource$edit', '', NULL, 1, '2026-04-19 12:58:47', 1, '2026-04-19 12:58:47');
INSERT INTO `sys_menu` VALUES (100287, 100283, 3, 4, '', '', '', 1, 0, 0, 'sys:lang:resource$delete', '', NULL, 1, '2026-04-19 12:58:47', 1, '2026-04-19 12:58:47');
INSERT INTO `sys_menu` VALUES (100288, 100109, 3, 5, '', '', '', 1, 0, 0, 'job$trigger', '', 0, 0, '2026-04-19 15:57:07', 0, '2026-04-20 15:00:37');
INSERT INTO `sys_menu` VALUES (100289, 100296, 2, 9, '/system/notice', 'system/notice/index', '', 1, 0, 0, 'sys:notice$list', 'notice', 0, 1, '2026-04-21 14:36:13', 1, '2026-04-23 16:14:59');
INSERT INTO `sys_menu` VALUES (100290, 100289, 3, 1, '', '', '', 1, 0, 0, 'sys:notice$query', '', NULL, 1, '2026-04-21 14:36:13', 1, '2026-04-21 14:36:13');
INSERT INTO `sys_menu` VALUES (100291, 100289, 3, 2, '', '', '', 1, 0, 0, 'sys:notice$add', '', NULL, 1, '2026-04-21 14:36:13', 1, '2026-04-21 14:36:13');
INSERT INTO `sys_menu` VALUES (100292, 100289, 3, 3, '', '', '', 1, 0, 0, 'sys:notice$edit', '', NULL, 1, '2026-04-21 14:36:13', 1, '2026-04-21 14:36:13');
INSERT INTO `sys_menu` VALUES (100293, 100289, 3, 4, '', '', '', 1, 0, 0, 'sys:notice$delete', '', NULL, 1, '2026-04-21 14:36:13', 1, '2026-04-21 14:36:13');
INSERT INTO `sys_menu` VALUES (100294, 100000, 1, 1, '/system/user', '', '', 1, 0, 0, '', 'user-permissions', 0, 0, '2026-04-23 16:00:28', 0, '2026-04-24 16:32:32');
INSERT INTO `sys_menu` VALUES (100295, 100000, 1, 3, '', '', '', 1, 0, 0, '', 'setting', 0, 0, '2026-04-23 16:00:57', 0, '2026-04-24 15:15:56');
INSERT INTO `sys_menu` VALUES (100296, 100000, 1, 2, '', '', '', 1, 1, 0, '', 'operation', 0, 0, '2026-04-23 16:14:31', 0, '2026-05-01 10:46:41');
COMMIT;

-- ----------------------------
-- Table structure for sys_menu_tl
-- ----------------------------
DROP TABLE IF EXISTS `sys_menu_tl`;
CREATE TABLE `sys_menu_tl` (
  `menu_id` bigint NOT NULL DEFAULT '0' COMMENT '菜单Id',
  `lang_code` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '语言标识',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '菜单名称',
  PRIMARY KEY (`lang_code`,`menu_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='菜单翻译';

-- ----------------------------
-- Records of sys_menu_tl
-- ----------------------------
BEGIN;
INSERT INTO `sys_menu_tl` VALUES (100000, 'en-US', 'System');
INSERT INTO `sys_menu_tl` VALUES (100002, 'en-US', 'User');
INSERT INTO `sys_menu_tl` VALUES (100003, 'en-US', 'Query');
INSERT INTO `sys_menu_tl` VALUES (100004, 'en-US', 'Add');
INSERT INTO `sys_menu_tl` VALUES (100005, 'en-US', 'Edit');
INSERT INTO `sys_menu_tl` VALUES (100006, 'en-US', 'Delete');
INSERT INTO `sys_menu_tl` VALUES (100007, 'en-US', 'Export');
INSERT INTO `sys_menu_tl` VALUES (100008, 'en-US', 'Import');
INSERT INTO `sys_menu_tl` VALUES (100009, 'en-US', 'Reset password');
INSERT INTO `sys_menu_tl` VALUES (100010, 'en-US', 'Change password');
INSERT INTO `sys_menu_tl` VALUES (100011, 'en-US', 'Role');
INSERT INTO `sys_menu_tl` VALUES (100012, 'en-US', 'Query');
INSERT INTO `sys_menu_tl` VALUES (100013, 'en-US', 'Add');
INSERT INTO `sys_menu_tl` VALUES (100014, 'en-US', 'Edit');
INSERT INTO `sys_menu_tl` VALUES (100015, 'en-US', 'Delete');
INSERT INTO `sys_menu_tl` VALUES (100016, 'en-US', 'Export');
INSERT INTO `sys_menu_tl` VALUES (100017, 'en-US', 'Department');
INSERT INTO `sys_menu_tl` VALUES (100018, 'en-US', 'Query');
INSERT INTO `sys_menu_tl` VALUES (100019, 'en-US', 'Add');
INSERT INTO `sys_menu_tl` VALUES (100020, 'en-US', 'Edit');
INSERT INTO `sys_menu_tl` VALUES (100021, 'en-US', 'Delete');
INSERT INTO `sys_menu_tl` VALUES (100022, 'en-US', 'Position');
INSERT INTO `sys_menu_tl` VALUES (100023, 'en-US', 'Query');
INSERT INTO `sys_menu_tl` VALUES (100024, 'en-US', 'Add');
INSERT INTO `sys_menu_tl` VALUES (100025, 'en-US', 'Edit');
INSERT INTO `sys_menu_tl` VALUES (100026, 'en-US', 'Delete');
INSERT INTO `sys_menu_tl` VALUES (100027, 'en-US', 'Export');
INSERT INTO `sys_menu_tl` VALUES (100053, 'en-US', 'Menu');
INSERT INTO `sys_menu_tl` VALUES (100054, 'en-US', 'Query');
INSERT INTO `sys_menu_tl` VALUES (100055, 'en-US', 'Add');
INSERT INTO `sys_menu_tl` VALUES (100056, 'en-US', 'Edit');
INSERT INTO `sys_menu_tl` VALUES (100057, 'en-US', 'Delete');
INSERT INTO `sys_menu_tl` VALUES (100058, 'en-US', 'Dictionary');
INSERT INTO `sys_menu_tl` VALUES (100059, 'en-US', 'Query');
INSERT INTO `sys_menu_tl` VALUES (100060, 'en-US', 'Edit');
INSERT INTO `sys_menu_tl` VALUES (100061, 'en-US', 'Delete');
INSERT INTO `sys_menu_tl` VALUES (100062, 'en-US', 'Export');
INSERT INTO `sys_menu_tl` VALUES (100063, 'en-US', 'Settings');
INSERT INTO `sys_menu_tl` VALUES (100064, 'en-US', 'Query');
INSERT INTO `sys_menu_tl` VALUES (100065, 'en-US', 'Add');
INSERT INTO `sys_menu_tl` VALUES (100066, 'en-US', 'Edit');
INSERT INTO `sys_menu_tl` VALUES (100067, 'en-US', 'Delete');
INSERT INTO `sys_menu_tl` VALUES (100068, 'en-US', 'Export');
INSERT INTO `sys_menu_tl` VALUES (100094, 'en-US', 'Monitoring');
INSERT INTO `sys_menu_tl` VALUES (100095, 'en-US', 'Online users');
INSERT INTO `sys_menu_tl` VALUES (100097, 'en-US', 'Force logout');
INSERT INTO `sys_menu_tl` VALUES (100099, 'en-US', 'Service Monitoring');
INSERT INTO `sys_menu_tl` VALUES (100100, 'en-US', 'Operation log');
INSERT INTO `sys_menu_tl` VALUES (100102, 'en-US', 'Delete');
INSERT INTO `sys_menu_tl` VALUES (100104, 'en-US', 'Login log');
INSERT INTO `sys_menu_tl` VALUES (100105, 'en-US', 'Clear');
INSERT INTO `sys_menu_tl` VALUES (100106, 'en-US', 'Delete');
INSERT INTO `sys_menu_tl` VALUES (100108, 'en-US', 'Tools');
INSERT INTO `sys_menu_tl` VALUES (100109, 'en-US', 'Scheduler');
INSERT INTO `sys_menu_tl` VALUES (100110, 'en-US', 'Query');
INSERT INTO `sys_menu_tl` VALUES (100111, 'en-US', 'Add');
INSERT INTO `sys_menu_tl` VALUES (100112, 'en-US', 'Edit');
INSERT INTO `sys_menu_tl` VALUES (100113, 'en-US', 'Delete');
INSERT INTO `sys_menu_tl` VALUES (100119, 'en-US', 'Cache tool');
INSERT INTO `sys_menu_tl` VALUES (100120, 'en-US', 'Code generation');
INSERT INTO `sys_menu_tl` VALUES (100121, 'en-US', 'Query');
INSERT INTO `sys_menu_tl` VALUES (100122, 'en-US', 'Edit');
INSERT INTO `sys_menu_tl` VALUES (100123, 'en-US', 'Delete');
INSERT INTO `sys_menu_tl` VALUES (100124, 'en-US', 'Import');
INSERT INTO `sys_menu_tl` VALUES (100125, 'en-US', 'Preview');
INSERT INTO `sys_menu_tl` VALUES (100126, 'en-US', 'Download');
INSERT INTO `sys_menu_tl` VALUES (100183, 'en-US', 'Delete log');
INSERT INTO `sys_menu_tl` VALUES (100184, 'en-US', 'Log');
INSERT INTO `sys_menu_tl` VALUES (100185, 'en-US', 'Delete');
INSERT INTO `sys_menu_tl` VALUES (100221, 'en-US', 'Globalization');
INSERT INTO `sys_menu_tl` VALUES (100249, 'en-US', 'Application');
INSERT INTO `sys_menu_tl` VALUES (100250, 'en-US', 'Query');
INSERT INTO `sys_menu_tl` VALUES (100251, 'en-US', 'Add');
INSERT INTO `sys_menu_tl` VALUES (100252, 'en-US', 'Edit');
INSERT INTO `sys_menu_tl` VALUES (100253, 'en-US', 'Delete');
INSERT INTO `sys_menu_tl` VALUES (100254, 'en-US', 'Export');
INSERT INTO `sys_menu_tl` VALUES (100261, 'en-US', 'Design');
INSERT INTO `sys_menu_tl` VALUES (100263, 'en-US', 'Design');
INSERT INTO `sys_menu_tl` VALUES (100267, 'en-US', 'Add');
INSERT INTO `sys_menu_tl` VALUES (100268, 'en-US', 'Refresh');
INSERT INTO `sys_menu_tl` VALUES (100269, 'en-US', 'Import');
INSERT INTO `sys_menu_tl` VALUES (100270, 'en-US', 'Export');
INSERT INTO `sys_menu_tl` VALUES (100271, 'en-US', 'Clear');
INSERT INTO `sys_menu_tl` VALUES (100272, 'en-US', 'Delete');
INSERT INTO `sys_menu_tl` VALUES (100278, 'en-US', 'Language');
INSERT INTO `sys_menu_tl` VALUES (100279, 'en-US', 'Query');
INSERT INTO `sys_menu_tl` VALUES (100280, 'en-US', 'Add');
INSERT INTO `sys_menu_tl` VALUES (100281, 'en-US', 'Edit');
INSERT INTO `sys_menu_tl` VALUES (100282, 'en-US', 'Delete');
INSERT INTO `sys_menu_tl` VALUES (100283, 'en-US', 'Language resource');
INSERT INTO `sys_menu_tl` VALUES (100284, 'en-US', 'Query');
INSERT INTO `sys_menu_tl` VALUES (100285, 'en-US', 'Add');
INSERT INTO `sys_menu_tl` VALUES (100286, 'en-US', 'Edit');
INSERT INTO `sys_menu_tl` VALUES (100287, 'en-US', 'Delete');
INSERT INTO `sys_menu_tl` VALUES (100288, 'en-US', 'Trigger');
INSERT INTO `sys_menu_tl` VALUES (100289, 'en-US', 'Notice');
INSERT INTO `sys_menu_tl` VALUES (100290, 'en-US', 'Query');
INSERT INTO `sys_menu_tl` VALUES (100291, 'en-US', 'Add');
INSERT INTO `sys_menu_tl` VALUES (100292, 'en-US', 'Edit');
INSERT INTO `sys_menu_tl` VALUES (100293, 'en-US', 'Delete');
INSERT INTO `sys_menu_tl` VALUES (100294, 'en-US', 'Permission');
INSERT INTO `sys_menu_tl` VALUES (100295, 'en-US', 'Configuration');
INSERT INTO `sys_menu_tl` VALUES (100296, 'en-US', 'Operations');
INSERT INTO `sys_menu_tl` VALUES (100000, 'zh-CN', '系统');
INSERT INTO `sys_menu_tl` VALUES (100002, 'zh-CN', '用户管理');
INSERT INTO `sys_menu_tl` VALUES (100003, 'zh-CN', '查询');
INSERT INTO `sys_menu_tl` VALUES (100004, 'zh-CN', '新增');
INSERT INTO `sys_menu_tl` VALUES (100005, 'zh-CN', '用户修改');
INSERT INTO `sys_menu_tl` VALUES (100006, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` VALUES (100007, 'zh-CN', '导出');
INSERT INTO `sys_menu_tl` VALUES (100008, 'zh-CN', '导入');
INSERT INTO `sys_menu_tl` VALUES (100009, 'zh-CN', '重置密码');
INSERT INTO `sys_menu_tl` VALUES (100010, 'zh-CN', '修改密码');
INSERT INTO `sys_menu_tl` VALUES (100011, 'zh-CN', '角色管理');
INSERT INTO `sys_menu_tl` VALUES (100012, 'zh-CN', '查询');
INSERT INTO `sys_menu_tl` VALUES (100013, 'zh-CN', '新增');
INSERT INTO `sys_menu_tl` VALUES (100014, 'zh-CN', '修改');
INSERT INTO `sys_menu_tl` VALUES (100015, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` VALUES (100016, 'zh-CN', '导出');
INSERT INTO `sys_menu_tl` VALUES (100017, 'zh-CN', '部门管理');
INSERT INTO `sys_menu_tl` VALUES (100018, 'zh-CN', '查询');
INSERT INTO `sys_menu_tl` VALUES (100019, 'zh-CN', '新增');
INSERT INTO `sys_menu_tl` VALUES (100020, 'zh-CN', '修改');
INSERT INTO `sys_menu_tl` VALUES (100021, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` VALUES (100022, 'zh-CN', '岗位管理');
INSERT INTO `sys_menu_tl` VALUES (100023, 'zh-CN', '查询');
INSERT INTO `sys_menu_tl` VALUES (100024, 'zh-CN', '新增');
INSERT INTO `sys_menu_tl` VALUES (100025, 'zh-CN', '修改');
INSERT INTO `sys_menu_tl` VALUES (100026, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` VALUES (100027, 'zh-CN', '导出');
INSERT INTO `sys_menu_tl` VALUES (100053, 'zh-CN', '菜单管理');
INSERT INTO `sys_menu_tl` VALUES (100054, 'zh-CN', '查询');
INSERT INTO `sys_menu_tl` VALUES (100055, 'zh-CN', '新增');
INSERT INTO `sys_menu_tl` VALUES (100056, 'zh-CN', '修改');
INSERT INTO `sys_menu_tl` VALUES (100057, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` VALUES (100058, 'zh-CN', '字典管理');
INSERT INTO `sys_menu_tl` VALUES (100059, 'zh-CN', '查询');
INSERT INTO `sys_menu_tl` VALUES (100060, 'zh-CN', '修改');
INSERT INTO `sys_menu_tl` VALUES (100061, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` VALUES (100062, 'zh-CN', '导出');
INSERT INTO `sys_menu_tl` VALUES (100063, 'zh-CN', '系统配置');
INSERT INTO `sys_menu_tl` VALUES (100064, 'zh-CN', '查询');
INSERT INTO `sys_menu_tl` VALUES (100065, 'zh-CN', '新增');
INSERT INTO `sys_menu_tl` VALUES (100066, 'zh-CN', '修改');
INSERT INTO `sys_menu_tl` VALUES (100067, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` VALUES (100068, 'zh-CN', '导出');
INSERT INTO `sys_menu_tl` VALUES (100094, 'zh-CN', '系统监控');
INSERT INTO `sys_menu_tl` VALUES (100095, 'zh-CN', '在线用户');
INSERT INTO `sys_menu_tl` VALUES (100097, 'zh-CN', '强制下线');
INSERT INTO `sys_menu_tl` VALUES (100099, 'zh-CN', '服务监控');
INSERT INTO `sys_menu_tl` VALUES (100100, 'zh-CN', '操作日志');
INSERT INTO `sys_menu_tl` VALUES (100102, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` VALUES (100104, 'zh-CN', '登录日志');
INSERT INTO `sys_menu_tl` VALUES (100105, 'zh-CN', '清空');
INSERT INTO `sys_menu_tl` VALUES (100106, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` VALUES (100108, 'zh-CN', '工具');
INSERT INTO `sys_menu_tl` VALUES (100109, 'zh-CN', '定时任务');
INSERT INTO `sys_menu_tl` VALUES (100110, 'zh-CN', '查询');
INSERT INTO `sys_menu_tl` VALUES (100111, 'zh-CN', '新增');
INSERT INTO `sys_menu_tl` VALUES (100112, 'zh-CN', '修改');
INSERT INTO `sys_menu_tl` VALUES (100113, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` VALUES (100119, 'zh-CN', '缓存工具');
INSERT INTO `sys_menu_tl` VALUES (100120, 'zh-CN', '代码生成');
INSERT INTO `sys_menu_tl` VALUES (100121, 'zh-CN', '查询');
INSERT INTO `sys_menu_tl` VALUES (100122, 'zh-CN', '修改');
INSERT INTO `sys_menu_tl` VALUES (100123, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` VALUES (100124, 'zh-CN', '导入');
INSERT INTO `sys_menu_tl` VALUES (100125, 'zh-CN', '预览');
INSERT INTO `sys_menu_tl` VALUES (100126, 'zh-CN', '下载');
INSERT INTO `sys_menu_tl` VALUES (100183, 'zh-CN', '删除日志');
INSERT INTO `sys_menu_tl` VALUES (100184, 'zh-CN', '日志');
INSERT INTO `sys_menu_tl` VALUES (100185, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` VALUES (100221, 'zh-CN', '国际化');
INSERT INTO `sys_menu_tl` VALUES (100225, 'zh-CN', '用户列表');
INSERT INTO `sys_menu_tl` VALUES (100226, 'zh-CN', '隐藏用户');
INSERT INTO `sys_menu_tl` VALUES (100227, 'zh-CN', '用户分群');
INSERT INTO `sys_menu_tl` VALUES (100228, 'zh-CN', '高价值用户');
INSERT INTO `sys_menu_tl` VALUES (100229, 'zh-CN', '召回用户');
INSERT INTO `sys_menu_tl` VALUES (100231, 'zh-CN', '订单概览');
INSERT INTO `sys_menu_tl` VALUES (100232, 'zh-CN', '退货订单');
INSERT INTO `sys_menu_tl` VALUES (100249, 'zh-CN', '终端应用');
INSERT INTO `sys_menu_tl` VALUES (100250, 'zh-CN', '查询');
INSERT INTO `sys_menu_tl` VALUES (100251, 'zh-CN', '新增');
INSERT INTO `sys_menu_tl` VALUES (100252, 'zh-CN', '修改');
INSERT INTO `sys_menu_tl` VALUES (100253, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` VALUES (100254, 'zh-CN', '导出');
INSERT INTO `sys_menu_tl` VALUES (100261, 'zh-CN', '设计');
INSERT INTO `sys_menu_tl` VALUES (100263, 'zh-CN', '设计');
INSERT INTO `sys_menu_tl` VALUES (100267, 'zh-CN', '新增');
INSERT INTO `sys_menu_tl` VALUES (100268, 'zh-CN', '刷新');
INSERT INTO `sys_menu_tl` VALUES (100269, 'zh-CN', '导入');
INSERT INTO `sys_menu_tl` VALUES (100270, 'zh-CN', '导出');
INSERT INTO `sys_menu_tl` VALUES (100271, 'zh-CN', '清空');
INSERT INTO `sys_menu_tl` VALUES (100272, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` VALUES (100278, 'zh-CN', '语言');
INSERT INTO `sys_menu_tl` VALUES (100279, 'zh-CN', '查询');
INSERT INTO `sys_menu_tl` VALUES (100280, 'zh-CN', '新增');
INSERT INTO `sys_menu_tl` VALUES (100281, 'zh-CN', '修改');
INSERT INTO `sys_menu_tl` VALUES (100282, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` VALUES (100283, 'zh-CN', '语言资源');
INSERT INTO `sys_menu_tl` VALUES (100284, 'zh-CN', '查询');
INSERT INTO `sys_menu_tl` VALUES (100285, 'zh-CN', '新增');
INSERT INTO `sys_menu_tl` VALUES (100286, 'zh-CN', '修改');
INSERT INTO `sys_menu_tl` VALUES (100287, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` VALUES (100288, 'zh-CN', '触发');
INSERT INTO `sys_menu_tl` VALUES (100289, 'zh-CN', '通知公告');
INSERT INTO `sys_menu_tl` VALUES (100290, 'zh-CN', '查询');
INSERT INTO `sys_menu_tl` VALUES (100291, 'zh-CN', '新增');
INSERT INTO `sys_menu_tl` VALUES (100292, 'zh-CN', '修改');
INSERT INTO `sys_menu_tl` VALUES (100293, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` VALUES (100294, 'zh-CN', '权限管理');
INSERT INTO `sys_menu_tl` VALUES (100295, 'zh-CN', '配置管理');
INSERT INTO `sys_menu_tl` VALUES (100296, 'zh-CN', '运营');
COMMIT;

-- ----------------------------
-- Table structure for sys_notice
-- ----------------------------
DROP TABLE IF EXISTS `sys_notice`;
CREATE TABLE `sys_notice` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '公告ID',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '公告标题',
  `type` tinyint(1) NOT NULL COMMENT '公告类型（1通知 2公告）',
  `content` varchar(3000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '公告内容',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态（0正常 1关闭）',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '备注',
  `created_by` bigint DEFAULT NULL COMMENT '创建者',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
  `updated_by` bigint DEFAULT NULL COMMENT '更新者',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=52 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='通知公告表';

-- ----------------------------
-- Records of sys_notice
-- ----------------------------
BEGIN;
INSERT INTO `sys_notice` VALUES (2, '系统维护通知', 1, '为提升系统稳定性，将于今晚22:00-23:00进行服务器维护，期间系统暂停服务，请提前做好准备。', 0, '常规维护通知', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (3, '账号安全提醒', 1, '请定期修改登录密码，避免使用简单密码，保护个人账号信息安全。', 0, '安全提醒', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (4, '功能更新通知', 1, '系统新增数据导出功能，支持Excel、CSV格式导出，欢迎使用。', 0, '功能更新', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (5, '节假日放假通知', 1, '根据国家法定节假日安排，本周末正常休息，工作日调整至下周。', 0, '节假日通知', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (6, '数据备份通知', 1, '系统每日凌晨2点自动备份数据，保障数据安全，请勿在此时段操作重要业务。', 1, '备份通知', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (7, '权限调整通知', 1, '部分用户权限已完成优化，如需调整权限请联系管理员。', 0, '权限优化', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (8, '登录异常提醒', 1, '检测到您的账号在陌生设备登录，如非本人操作请及时修改密码。', 0, '安全预警', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (9, '系统升级通知', 1, '系统版本已升级至V2.1，修复已知bug，提升运行速度。', 0, '版本升级', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (10, '文件上传规范', 1, '上传文件请控制大小在500MB以内，格式支持jpg、png、pdf、docx。', 1, '上传规范', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (11, '用户培训通知', 1, '本周三下午15点开展系统使用培训会议，请相关人员准时参加。', 0, '培训通知', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (12, '网络调试通知', 1, '明日10:00-11:00进行网络调试，可能出现短暂卡顿，敬请谅解。', 0, '网络调试', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (13, '密码重置提醒', 1, '您的密码已超过90天未修改，请尽快前往个人中心重置密码。', 0, '密码提醒', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (14, '接口更新通知', 1, '第三方数据接口已完成更新，同步功能恢复正常使用。', 1, '接口调整', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (15, '存储空间提醒', 1, '您的个人存储空间即将不足，建议清理无用文件释放空间。', 0, '存储提醒', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (16, '服务热线变更', 1, '客服服务热线已变更为400-123-4567，如有问题可致电咨询。', 0, '热线变更', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (17, '数据统计通知', 1, '本月数据统计已完成，可在报表中心查看详细统计结果。', 0, '统计通知', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (18, '账号锁定提醒', 1, '您的账号因多次输错密码已临时锁定，1小时后自动解锁。', 1, '账号锁定', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (19, '系统巡检通知', 1, '每周一进行系统全面巡检，保障系统稳定运行，无特殊情况不影响使用。', 0, '巡检通知', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (20, '模板更新通知', 1, '报表模板已优化更新，新增多项自定义配置项。', 0, '模板更新', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (21, '登录方式升级', 1, '系统新增人脸识别登录功能，提升登录安全性和便捷性。', 0, '登录升级', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (22, '数据清理通知', 1, '系统将清理超过1年的历史临时数据，如需保留请提前备份。', 1, '数据清理', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (23, '使用反馈通道', 1, '欢迎通过系统内反馈通道提交使用建议，我们将持续优化产品。', 0, '反馈通知', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (24, '服务器扩容通知', 1, '服务器已完成扩容，系统承载能力大幅提升，运行更流畅。', 0, '服务器扩容', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (25, '权限认证通知', 1, '首次登录需完成实名认证，认证通过后可使用全部功能。', 0, '实名认证', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (26, '系统暂停服务', 1, '因不可抗力因素，系统临时暂停服务，恢复时间另行通知。', 1, '紧急通知', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (27, '公司管理制度公告', 2, '新版公司管理制度已正式生效，全体员工请严格遵守执行。', 0, '制度公告', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (28, '项目进度公告', 2, '年度核心项目已完成70%进度，各部门按计划推进后续工作。', 0, '项目公告', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (29, '优秀员工表彰', 2, '表彰本月优秀员工张三、李四，感谢他们的辛勤付出与突出贡献。', 0, '表彰公告', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (30, '会议安排公告', 2, '下周一上午9点召开全体部门工作会议，地点：三楼会议室。', 1, '会议公告', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (31, '招聘信息公告', 2, '公司现招聘Java开发、前端工程师、产品经理若干，简历投递至hr@company.com。', 0, '招聘公告', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (32, '福利调整公告', 2, '员工福利待遇优化调整，新增交通补贴、餐饮补贴，本月起执行。', 0, '福利公告', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (33, '办公区域调整', 2, '市场部、运营部办公区域调整至二楼西区，今日起正式搬迁。', 0, '区域调整', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (34, '规章制度公示', 2, '新版考勤管理制度公示3天，无异议后将于下月正式执行。', 0, '制度公示', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (35, '项目中标公告', 2, '恭喜公司成功中标XX项目，感谢全体团队的努力与付出！', 1, '中标公告', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (36, '员工体检通知', 2, '年度员工体检安排在下周五，地点：市第一人民医院，请准时参加。', 0, '体检公告', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (37, '团建活动公告', 2, '本月组织户外团建活动，时间：周六全天，自愿报名参加。', 0, '团建公告', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (38, '办公用品申领', 2, '每月15号可申领办公用品，各部门统一统计后提交申请。', 0, '申领公告', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (39, '消防安全公告', 2, '办公区域消防设施已完成检查，全体员工需掌握基本消防知识。', 1, '消防公告', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (40, '年度总结公告', 2, '年度工作总结报告已发布，各部门需完成年度复盘工作。', 0, '总结公告', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (41, '合作签约公告', 2, '公司与XX企业达成战略合作，共同推进业务拓展与创新。', 0, '合作公告', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (42, '设备更新公告', 2, '办公电脑、打印机等设备已完成更新换代，提升办公效率。', 0, '设备更新', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (43, '加班补贴公告', 2, '节假日加班补贴标准已明确，按国家规定足额发放。', 1, '补贴公告', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (44, '企业文化宣传', 2, '秉承诚信、创新、协作、共赢的企业文化，携手共创未来。', 0, '文化宣传', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (45, '资质认证公告', 2, '公司成功通过ISO9001质量体系认证，企业资质再升级。', 0, '资质认证', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (46, '客户答谢公告', 2, '感谢所有客户的信任与支持，我们将持续提供优质服务。', 0, '客户答谢', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (47, '规章制度修订', 2, '根据实际运营情况，修订部分规章制度，即日起生效。', 1, '制度修订', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (48, '办公环境优化', 2, '办公区域绿植、照明、空调已全面优化，打造舒适办公环境。', 0, '环境优化', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (49, '技能竞赛公告', 2, '公司举办职业技能竞赛，设置多项奖励，欢迎员工踊跃报名。', 0, '竞赛公告', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (50, '供应商合作公告', 2, '新增3家优质供应商，保障物资供应稳定与质量。', 0, '供应商公告', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
INSERT INTO `sys_notice` VALUES (51, '年终盛典公告', 2, '年终盛典将于12月31日举办，全体员工共同参与欢度新年。', 1, '盛典公告', 1, '2026-04-21 15:05:43', 1, '2026-04-21 15:05:43');
COMMIT;

-- ----------------------------
-- Table structure for sys_oauth_client
-- ----------------------------
DROP TABLE IF EXISTS `sys_oauth_client`;
CREATE TABLE `sys_oauth_client` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `client_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '客户端ID',
  `client_secret` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '客户端秘钥',
  `client_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '客户端名称',
  `resources` varchar(255) DEFAULT NULL COMMENT '资源，按服务名称逗号分割(如 ALL,auth-service,sys-service)',
  `logo_uri` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '客户端LOGO',
  `scope` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '授权范围',
  `authorized_grant_types` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '授权类型',
  `web_server_redirect_uri` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '回调地址',
  `access_token_validity` int NOT NULL COMMENT '访问令牌有效期（秒）',
  `refresh_token_validity` int NOT NULL COMMENT '刷新令牌有效期（秒）',
  `autoapprove` tinyint(1) NOT NULL COMMENT '自动授权（0否，1是）',
  `created_by` bigint DEFAULT NULL COMMENT '创建人',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
  `updated_by` bigint DEFAULT NULL COMMENT '修改人',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='终端应用';

-- ----------------------------
-- Records of sys_oauth_client
-- ----------------------------
BEGIN;
INSERT INTO `sys_oauth_client` VALUES (1, '100000', '$2a$10$1UpM1wYjOm.LNpWeu/FouupAkoZoxugN30tiqMphbjsaoOuhy82ea', 'Master', 'ALL', '/assets/images/logo.png', 'openid,profile', 'password,refresh_token,mobile,passkey,authorization_code,email', 'http://192.168.31.82:5173/callback', 43200, 2592000, 1, 1, '2024-04-10 13:42:46', 1, '2026-04-20 16:22:22');
INSERT INTO `sys_oauth_client` VALUES (3, 'auth-demo', '$2a$10$1UpM1wYjOm.LNpWeu/FouupAkoZoxugN30tiqMphbjsaoOuhy82ea', '授权应用示例', 'ALL', '/assets/images/logo.png', 'openid,profile', 'password,refresh_token,authorization_code,client_credentials,mobile,email,passkey', 'http://192.168.31.82:5174/callback', 43200, 2592000, 0, 1, '2024-04-10 13:42:46', 1, '2026-04-21 16:54:52');
COMMIT;

-- ----------------------------
-- Table structure for sys_oper_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_oper_log`;
CREATE TABLE `sys_oper_log` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '编号',
  `module` varchar(50) DEFAULT NULL COMMENT '模块',
  `type` char(10) DEFAULT NULL COMMENT '类型',
  `description` varchar(128) DEFAULT NULL COMMENT '描述',
  `method` varchar(10) DEFAULT NULL COMMENT '请求方法',
  `url` varchar(255) DEFAULT NULL COMMENT '请求URL',
  `ip` varchar(128) DEFAULT NULL COMMENT '请求IP',
  `location` varchar(255) DEFAULT NULL COMMENT '请求地点',
  `payload` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '请求参数',
  `result` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '响应数据',
  `device` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '设备',
  `os` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '系统',
  `browser` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '浏览器',
  `status` tinyint(1) DEFAULT NULL COMMENT '状态（0失败 1成功）',
  `error` varchar(2000) DEFAULT NULL COMMENT '错误消息',
  `time` bigint DEFAULT NULL COMMENT '耗时（毫秒）',
  `user_id` bigint DEFAULT NULL COMMENT '用户编号',
  `username` varchar(50) DEFAULT NULL COMMENT '用户名',
  `oper_at` datetime DEFAULT NULL COMMENT '操作时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='操作日志记录';

-- ----------------------------
-- Records of sys_oper_log
-- ----------------------------
BEGIN;
INSERT INTO `sys_oper_log` VALUES (1, 'OperLog-[操作日志清空]', 'CLEAN', '', 'DELETE', '/api/v1/system/oper/log', '127.0.0.1', '', '', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 60, 2, 'develop', '2026-04-25 17:51:16');
INSERT INTO `sys_oper_log` VALUES (2, 'Lang-[语言]', 'UPDATE', '', 'PUT', '/api/v1/system/lang/3', '127.0.0.1', '', '{\"createdAt\":\"2026-02-01T15:57:25+08:00\",\"id\":\"3\",\"isDefault\":0,\"langCode\":\"ru-Ru\",\"langName\":\"🇷🇺 Русский\",\"remark\":\"俄语\",\"sort\":3,\"status\":0,\"updatedAt\":\"2026-04-19T13:45:14+08:00\"}', '{\"code\":0,\"message\":\"success\",\"data\":3}\n', 'PC', 'Linux', 'Chrome', 1, '', 3, 2, 'develop', '2026-04-25 20:20:27');
INSERT INTO `sys_oper_log` VALUES (3, 'Lang-[语言]', 'UPDATE', '', 'PUT', '/api/v1/system/lang/3', '127.0.0.1', '', '{\"createdAt\":\"2026-02-01T15:57:25+08:00\",\"id\":\"3\",\"isDefault\":0,\"langCode\":\"ru-Ru\",\"langName\":\"🇷🇺 Русский\",\"remark\":\"俄语\",\"sort\":3,\"status\":1,\"updatedAt\":\"2026-04-25T20:20:27+08:00\"}', '{\"code\":0,\"message\":\"success\",\"data\":3}\n', 'PC', 'Linux', 'Chrome', 1, '', 44, 2, 'develop', '2026-04-25 20:20:50');
INSERT INTO `sys_oper_log` VALUES (4, '导入数据库表', 'IMPORT', '', 'POST', '/api/v1/gen/table/import', '127.0.0.1', '', '[\"sys_user_credential\"]', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 7, 2, 'develop', '2026-04-25 20:33:32');
INSERT INTO `sys_oper_log` VALUES (5, '代码生成表', 'UPDATE', '', 'PUT', '/api/v1/gen/table/268', '127.0.0.1', '', '{\"author\":\"\",\"className\":\"UserCredential\",\"createdAt\":\"2026-04-25T20:33:32+08:00\",\"id\":\"268\",\"menuId\":\"0\",\"moduleName\":\"system\",\"permission\":\"sys:user:credential\",\"remark\":\"\",\"router\":\"user/credential\",\"tableComment\":\"用户Passkey凭据\",\"tableName\":\"sys_user_credential\",\"updatedAt\":\"2026-04-25T20:33:32+08:00\"}', '{\"code\":0,\"message\":\"success\",\"data\":268}\n', 'PC', 'Linux', 'Chrome', 1, '', 4, 2, 'develop', '2026-04-25 20:34:01');
INSERT INTO `sys_oper_log` VALUES (6, 'Menu-[菜单权限表]', 'UPDATE', '', 'PUT', '/api/v1/system/menu/100296', '127.0.0.1', '', '{\"activeId\":\"0\",\"component\":\"\",\"createdAt\":\"0001-01-01T00:00:00Z\",\"createdBy\":\"0\",\"icon\":\"operation\",\"id\":\"100296\",\"isFrame\":0,\"langCode\":\"\",\"parentId\":\"100000\",\"path\":\"\",\"permission\":\"\",\"query\":\"\",\"sort\":2,\"status\":0,\"title\":\"\",\"translations\":{\"en-US\":\"Operations\",\"zh-CN\":\"运营\"},\"type\":1,\"updatedAt\":\"0001-01-01T00:00:00Z\",\"updatedBy\":\"0\",\"visible\":1}', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 43, 2, 'develop', '2026-04-26 17:04:25');
INSERT INTO `sys_oper_log` VALUES (7, 'User-[用户管理]', 'UPDATE', '', 'PUT', '/api/v1/system/user/1/reset-password', '127.0.0.1', '', '{\"password\":\"123456\"}', '{\"code\":0,\"message\":\"success\",\"data\":1}\n', 'PC', 'Linux', 'Chrome', 1, '', 55, 2, 'develop', '2026-04-28 10:59:57');
INSERT INTO `sys_oper_log` VALUES (8, 'User-[用户管理]', 'UPDATE', '', 'PUT', '/api/v1/system/user/2', '127.0.0.1', '', '{\"autograph\":\"\",\"avatar\":\"\",\"deptPosts\":[{\"deptId\":\"112\",\"postId\":\"4\"}],\"email\":\"looses118@gmail.com\",\"mobile\":\"\",\"nickname\":\"Develop\",\"password\":\"\",\"roleIds\":[\"3\"],\"sex\":0,\"status\":0,\"type\":0,\"username\":\"develop\"}', '{\"code\":0,\"message\":\"success\",\"data\":2}\n', 'PC', 'Linux', 'Chrome', 1, '', 3, 2, 'develop', '2026-04-28 18:19:59');
INSERT INTO `sys_oper_log` VALUES (9, '字典数据', 'UPDATE', '', 'PUT', '/api/v1/system/dict/data/${id}', '127.0.0.1', '', '{\"dictType\":\"flow_active_status\",\"dictValue\":\"1\",\"id\":\"128\",\"label\":\"\",\"langCode\":\"\",\"sort\":\"1\",\"status\":0,\"tagClass\":\"\",\"tagType\":\"success\",\"tip\":\"\",\"translations\":{\"en-US\":{\"label\":\"Active\",\"tip\":\"\"},\"zh-CN\":{\"label\":\"激活\",\"tip\":\"\"}}}', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 3, 2, 'develop', '2026-04-29 12:32:40');
INSERT INTO `sys_oper_log` VALUES (10, '字典数据', 'UPDATE', '', 'PUT', '/api/v1/system/dict/data/${id}', '127.0.0.1', '', '{\"dictType\":\"flow_deploy_status\",\"dictValue\":\"0\",\"id\":\"53\",\"label\":\"\",\"langCode\":\"\",\"sort\":\"1\",\"status\":0,\"tagClass\":\"\",\"tagType\":\"success\",\"tip\":\"\",\"translations\":{\"en-US\":{\"label\":\"Normal\",\"tip\":\"\"},\"zh-CN\":{\"label\":\"正常\",\"tip\":\"\"}}}', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 4, 2, 'develop', '2026-04-29 12:36:53');
INSERT INTO `sys_oper_log` VALUES (11, '字典数据', 'UPDATE', '', 'PUT', '/api/v1/system/dict/data/${id}', '127.0.0.1', '', '{\"dictType\":\"flow_deploy_status\",\"dictValue\":\"4\",\"id\":\"54\",\"label\":\"\",\"langCode\":\"\",\"sort\":\"2\",\"status\":0,\"tagClass\":\"\",\"tagType\":\"warning\",\"tip\":\"\",\"translations\":{\"en-US\":{\"label\":\"Outdated\",\"tip\":\"\"},\"zh-CN\":{\"label\":\"过时\",\"tip\":\"\"}}}', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 3, 2, 'develop', '2026-04-29 12:37:20');
INSERT INTO `sys_oper_log` VALUES (12, '字典数据', 'UPDATE', '', 'PUT', '/api/v1/system/dict/data/${id}', '127.0.0.1', '', '{\"dictType\":\"flow_approve_status\",\"dictValue\":\"0\",\"id\":\"48\",\"label\":\"\",\"langCode\":\"\",\"sort\":\"0\",\"status\":0,\"tagClass\":\"\",\"tagType\":\"info\",\"tip\":\"\",\"translations\":{\"en-US\":{\"label\":\"Not submitted\",\"tip\":\"\"},\"zh-CN\":{\"label\":\"未提交\",\"tip\":\"\"}}}', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 3, 2, 'develop', '2026-04-29 12:38:07');
INSERT INTO `sys_oper_log` VALUES (13, '字典数据', 'UPDATE', '', 'PUT', '/api/v1/system/dict/data/${id}', '127.0.0.1', '', '{\"dictType\":\"flow_approve_status\",\"dictValue\":\"1\",\"id\":\"49\",\"label\":\"\",\"langCode\":\"\",\"sort\":\"1\",\"status\":0,\"tagClass\":\"\",\"tagType\":\"info\",\"tip\":\"\",\"translations\":{\"en-US\":{\"label\":\"Under Approval\",\"tip\":\"\"},\"zh-CN\":{\"label\":\"审批中\",\"tip\":\"\"}}}', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 3, 2, 'develop', '2026-04-29 12:38:21');
INSERT INTO `sys_oper_log` VALUES (14, '字典数据', 'UPDATE', '', 'PUT', '/api/v1/system/dict/data/${id}', '127.0.0.1', '', '{\"dictType\":\"flow_approve_status\",\"dictValue\":\"2\",\"id\":\"50\",\"label\":\"\",\"langCode\":\"\",\"sort\":\"2\",\"status\":0,\"tagClass\":\"\",\"tagType\":\"warning\",\"tip\":\"\",\"translations\":{\"en-US\":{\"label\":\"Revoked\",\"tip\":\"\"},\"zh-CN\":{\"label\":\"已撤销\",\"tip\":\"\"}}}', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 2, 2, 'develop', '2026-04-29 12:38:33');
INSERT INTO `sys_oper_log` VALUES (15, '字典数据', 'UPDATE', '', 'PUT', '/api/v1/system/dict/data/${id}', '127.0.0.1', '', '{\"dictType\":\"flow_approve_status\",\"dictValue\":\"4\",\"id\":\"51\",\"label\":\"\",\"langCode\":\"\",\"sort\":\"4\",\"status\":0,\"tagClass\":\"\",\"tagType\":\"danger\",\"tip\":\"\",\"translations\":{\"en-US\":{\"label\":\"Rejected\",\"tip\":\"\"},\"zh-CN\":{\"label\":\"已驳回\",\"tip\":\"\"}}}', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 3, 2, 'develop', '2026-04-29 12:38:44');
INSERT INTO `sys_oper_log` VALUES (16, '字典数据', 'UPDATE', '', 'PUT', '/api/v1/system/dict/data/${id}', '127.0.0.1', '', '{\"dictType\":\"flow_approve_status\",\"dictValue\":\"6\",\"id\":\"52\",\"label\":\"\",\"langCode\":\"\",\"sort\":\"6\",\"status\":0,\"tagClass\":\"\",\"tagType\":\"success\",\"tip\":\"\",\"translations\":{\"en-US\":{\"label\":\"Approved\",\"tip\":\"\"},\"zh-CN\":{\"label\":\"已通过\",\"tip\":\"\"}}}', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 3, 2, 'develop', '2026-04-29 12:38:54');
INSERT INTO `sys_oper_log` VALUES (17, '字典数据', 'UPDATE', '', 'PUT', '/api/v1/system/dict/data/${id}', '127.0.0.1', '', '{\"dictType\":\"flow_category\",\"dictValue\":\"0\",\"id\":\"46\",\"label\":\"\",\"langCode\":\"\",\"sort\":\"1\",\"status\":0,\"tagClass\":\"\",\"tagType\":\"success\",\"tip\":\"\",\"translations\":{\"en-US\":{\"label\":\"Approval process\",\"tip\":\"\"},\"zh-CN\":{\"label\":\"审批流程\",\"tip\":\"\"}}}', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 4, 2, 'develop', '2026-04-29 12:39:09');
INSERT INTO `sys_oper_log` VALUES (18, '字典数据', 'UPDATE', '', 'PUT', '/api/v1/system/dict/data/${id}', '127.0.0.1', '', '{\"dictType\":\"flow_category\",\"dictValue\":\"1\",\"id\":\"47\",\"label\":\"\",\"langCode\":\"\",\"sort\":\"2\",\"status\":0,\"tagClass\":\"\",\"tagType\":\"warning\",\"tip\":\"\",\"translations\":{\"en-US\":{\"label\":\"Business Process\",\"tip\":\"\"},\"zh-CN\":{\"label\":\"业务流程\",\"tip\":\"\"}}}', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 3, 2, 'develop', '2026-04-29 12:39:20');
INSERT INTO `sys_oper_log` VALUES (19, '字典数据', 'DELETE', '', 'DELETE', '/api/v1/system/dict/data/146,147', '127.0.0.1', '', '', '{\"code\":0,\"message\":\"success\",\"data\":[146,147]}\n', 'PC', 'Linux', 'Chrome', 1, '', 4, 2, 'develop', '2026-04-29 12:40:07');
INSERT INTO `sys_oper_log` VALUES (20, '字典类型', 'DELETE', '', 'DELETE', '/api/v1/system/dict/262', '127.0.0.1', '', '', '{\"code\":0,\"message\":\"success\",\"data\":262}\n', 'PC', 'Linux', 'Chrome', 1, '', 4, 2, 'develop', '2026-04-29 12:40:17');
INSERT INTO `sys_oper_log` VALUES (21, 'Menu-[菜单权限表]', 'UPDATE', '', 'PUT', '/api/v1/system/menu/100296', '127.0.0.1', '', '{\"activeId\":\"0\",\"component\":\"\",\"createdAt\":\"0001-01-01T00:00:00Z\",\"createdBy\":\"0\",\"icon\":\"operation\",\"id\":\"100296\",\"isFrame\":0,\"langCode\":\"\",\"parentId\":\"100000\",\"path\":\"\",\"permission\":\"\",\"query\":\"\",\"sort\":2,\"status\":1,\"title\":\"\",\"translations\":{\"en-US\":\"Operations\",\"zh-CN\":\"运营\"},\"type\":1,\"updatedAt\":\"0001-01-01T00:00:00Z\",\"updatedBy\":\"0\",\"visible\":1}', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 10, 2, 'develop', '2026-05-01 10:46:41');
INSERT INTO `sys_oper_log` VALUES (23, 'Online-[在线用户]', 'DELETE', 'Online-[在线用户]', 'DELETE', '/api/v1/sys/monitor/online/728cd7fb-9ce5-485e-9fa7-f512b3aef76d', '192.168.0.102', '', '', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 0, 2, 'develop', '2026-05-18 17:20:32');
INSERT INTO `sys_oper_log` VALUES (24, 'Online-[在线用户]', 'DELETE', 'Online-[在线用户]', 'DELETE', '/api/v1/sys/monitor/online/f34e102e-0204-4b36-a59f-ab8717f3d834', '192.168.0.102', '', '', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 0, 2, 'develop', '2026-05-18 17:20:34');
INSERT INTO `sys_oper_log` VALUES (25, 'Online-[在线用户]', 'DELETE', 'Online-[在线用户]', 'DELETE', '/api/v1/sys/monitor/online/7cd0eff2-504a-419b-8316-5de9ec76b381', '192.168.0.102', '', '', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 0, 2, 'develop', '2026-05-18 17:20:36');
INSERT INTO `sys_oper_log` VALUES (26, 'Conf-[参数配置]', 'UPDATE', 'Conf-[参数配置]', 'PUT', '/api/v1/sys/conf/1', '192.168.0.102', '', '{\"id\":\"1\",\"name\":\"账号初始密码\",\"confKey\":\"sys.user.initPassword\",\"confValue\":\"123456\",\"isSys\":1,\"scope\":2,\"isSecret\":0,\"remark\":\"初始化密码 123456\",\"createdBy\":\"1\",\"createdAt\":\"2022-12-11T16:51:52+08:00\",\"updatedAt\":\"2026-04-18T08:19:55+08:00\",\"updatedBy\":\"1\"}', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 5, 2, 'develop', '2026-05-18 18:05:16');
INSERT INTO `sys_oper_log` VALUES (27, 'Conf-[参数配置]', 'UPDATE', 'Conf-[参数配置]', 'PUT', '/api/v1/sys/conf/1', '192.168.0.102', '', '{\"id\":\"1\",\"name\":\"账号初始密码\",\"confKey\":\"sys.user.initPassword\",\"confValue\":\"123456\",\"isSys\":1,\"scope\":2,\"isSecret\":1,\"remark\":\"初始化密码 123456\",\"createdBy\":\"1\",\"createdAt\":\"2022-12-11T16:51:52+08:00\",\"updatedAt\":\"2026-05-18T18:05:16+08:00\",\"updatedBy\":\"1\"}', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 4, 2, 'develop', '2026-05-18 18:05:19');
INSERT INTO `sys_oper_log` VALUES (28, 'Conf-[参数配置]', 'UPDATE', 'Conf-[参数配置]', 'PUT', '/api/v1/sys/conf/18', '192.168.0.102', '', '{\"id\":\"18\",\"name\":\"发件地址\",\"confKey\":\"sys.mail.email\",\"confValue\":\"1073602@qq.com\",\"isSys\":1,\"scope\":2,\"isSecret\":1,\"remark\":\"发件地址\",\"createdBy\":\"1\",\"createdAt\":\"2025-08-09T20:56:38+08:00\",\"updatedAt\":\"2026-04-20T12:30:55+08:00\",\"updatedBy\":\"1\"}', '{\"code\":0,\"message\":\"success\",\"data\":null}\n', 'PC', 'Linux', 'Chrome', 1, '', 7, 2, 'develop', '2026-05-18 18:05:31');
COMMIT;

-- ----------------------------
-- Table structure for sys_post
-- ----------------------------
DROP TABLE IF EXISTS `sys_post`;
CREATE TABLE `sys_post` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '岗位编号',
  `post_key` varchar(32) DEFAULT NULL COMMENT '岗位标识',
  `post_name` varchar(50) DEFAULT NULL COMMENT '岗位名称',
  `sort` int DEFAULT NULL COMMENT '顺序',
  `status` tinyint(1) DEFAULT NULL COMMENT '状态（0正常 1停用）',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '备注',
  `created_by` bigint DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_by` bigint DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='岗位信息表';

-- ----------------------------
-- Records of sys_post
-- ----------------------------
BEGIN;
INSERT INTO `sys_post` VALUES (1, 'GM', '总经理', 98, 0, 'General manager', 1, '2022-12-11 16:51:52', 1, '2025-04-11 08:50:58');
INSERT INTO `sys_post` VALUES (2, 'DGM', '副总经理', 97, 0, 'Deputy general manager', 1, '2022-12-11 16:51:52', 1, '2025-08-21 12:41:34');
INSERT INTO `sys_post` VALUES (3, 'CFO', '财务总监', 87, 0, '', 1, '2022-12-11 16:51:52', 1, '2025-04-11 09:20:31');
INSERT INTO `sys_post` VALUES (4, 'LC', '分管领导', 6, 0, '', 1, '2022-12-11 16:51:52', 1, '2025-04-11 09:19:01');
INSERT INTO `sys_post` VALUES (5, 'Staff', '普遍员工', 1, 0, '', 1, '2023-09-14 09:22:34', 1, '2025-04-11 09:18:37');
INSERT INTO `sys_post` VALUES (6, 'Mgr', '部门经理', 5, 0, 'Director', 1, '2022-12-11 16:51:52', 1, '2025-04-11 09:20:53');
INSERT INTO `sys_post` VALUES (8, 'Supv', '主管', 4, 0, '', 1, '2023-11-03 14:16:12', 1, '2025-04-11 09:21:03');
INSERT INTO `sys_post` VALUES (9, 'TL', '组长', 3, 0, 'Team Leader', 1, '2023-11-03 14:17:12', 1, '2025-04-11 09:20:22');
INSERT INTO `sys_post` VALUES (10, 'CN', '出纳', 9, 0, '', 1, '2023-11-03 14:17:26', 1, '2025-04-04 12:29:07');
INSERT INTO `sys_post` VALUES (11, 'BoD', '董事长', 99, 0, '', 1, '2025-09-17 09:34:59', 0, '2025-09-17 09:35:06');
INSERT INTO `sys_post` VALUES (12, 'Eng', '工程师', 2, 0, '', 1, '2025-09-17 09:47:04', 0, '2025-09-17 09:47:08');
COMMIT;

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `role_name` varchar(50) DEFAULT NULL COMMENT '角色名称',
  `role_key` varchar(32) DEFAULT NULL COMMENT '角色权限字符串',
  `sort` int DEFAULT NULL COMMENT '显示顺序',
  `data_scope` tinyint DEFAULT NULL COMMENT '数据范围',
  `status` tinyint(1) DEFAULT NULL COMMENT '状态（0正常 1停用）',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '备注',
  `created_by` bigint DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_by` bigint DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='角色信息表';

-- ----------------------------
-- Records of sys_role
-- ----------------------------
BEGIN;
INSERT INTO `sys_role` VALUES (1, 'Administrator', 'admin', 1, 0, 0, '管理员', 1, '2022-12-11 16:51:52', 1, '2026-04-23 16:15:35');
INSERT INTO `sys_role` VALUES (2, 'Ordinary', 'ordinary', 2, 2, 0, '普通用户', 1, '2022-12-11 16:51:52', 1, '2026-04-23 16:15:34');
INSERT INTO `sys_role` VALUES (3, 'Develop', 'develop', 3, 0, 0, '开发用', 1, '2023-09-06 10:14:39', 1, '2026-04-23 16:15:32');
COMMIT;

-- ----------------------------
-- Table structure for sys_role_dept
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_dept`;
CREATE TABLE `sys_role_dept` (
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `dept_id` bigint NOT NULL COMMENT '部门ID',
  PRIMARY KEY (`role_id`,`dept_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='角色和部门关联表';

-- ----------------------------
-- Records of sys_role_dept
-- ----------------------------
BEGIN;
INSERT INTO `sys_role_dept` VALUES (2, 112);
INSERT INTO `sys_role_dept` VALUES (2, 113);
COMMIT;

-- ----------------------------
-- Table structure for sys_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_menu`;
CREATE TABLE `sys_role_menu` (
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `menu_id` bigint NOT NULL COMMENT '菜单ID',
  PRIMARY KEY (`role_id`,`menu_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='角色和菜单关联表';

-- ----------------------------
-- Records of sys_role_menu
-- ----------------------------
BEGIN;
INSERT INTO `sys_role_menu` VALUES (1, 100000);
INSERT INTO `sys_role_menu` VALUES (1, 100002);
INSERT INTO `sys_role_menu` VALUES (1, 100003);
INSERT INTO `sys_role_menu` VALUES (1, 100004);
INSERT INTO `sys_role_menu` VALUES (1, 100005);
INSERT INTO `sys_role_menu` VALUES (1, 100006);
INSERT INTO `sys_role_menu` VALUES (1, 100007);
INSERT INTO `sys_role_menu` VALUES (1, 100008);
INSERT INTO `sys_role_menu` VALUES (1, 100009);
INSERT INTO `sys_role_menu` VALUES (1, 100010);
INSERT INTO `sys_role_menu` VALUES (1, 100011);
INSERT INTO `sys_role_menu` VALUES (1, 100012);
INSERT INTO `sys_role_menu` VALUES (1, 100013);
INSERT INTO `sys_role_menu` VALUES (1, 100014);
INSERT INTO `sys_role_menu` VALUES (1, 100015);
INSERT INTO `sys_role_menu` VALUES (1, 100016);
INSERT INTO `sys_role_menu` VALUES (1, 100017);
INSERT INTO `sys_role_menu` VALUES (1, 100018);
INSERT INTO `sys_role_menu` VALUES (1, 100019);
INSERT INTO `sys_role_menu` VALUES (1, 100020);
INSERT INTO `sys_role_menu` VALUES (1, 100021);
INSERT INTO `sys_role_menu` VALUES (1, 100022);
INSERT INTO `sys_role_menu` VALUES (1, 100023);
INSERT INTO `sys_role_menu` VALUES (1, 100024);
INSERT INTO `sys_role_menu` VALUES (1, 100025);
INSERT INTO `sys_role_menu` VALUES (1, 100026);
INSERT INTO `sys_role_menu` VALUES (1, 100027);
INSERT INTO `sys_role_menu` VALUES (1, 100053);
INSERT INTO `sys_role_menu` VALUES (1, 100054);
INSERT INTO `sys_role_menu` VALUES (1, 100055);
INSERT INTO `sys_role_menu` VALUES (1, 100056);
INSERT INTO `sys_role_menu` VALUES (1, 100057);
INSERT INTO `sys_role_menu` VALUES (1, 100058);
INSERT INTO `sys_role_menu` VALUES (1, 100059);
INSERT INTO `sys_role_menu` VALUES (1, 100060);
INSERT INTO `sys_role_menu` VALUES (1, 100061);
INSERT INTO `sys_role_menu` VALUES (1, 100062);
INSERT INTO `sys_role_menu` VALUES (1, 100063);
INSERT INTO `sys_role_menu` VALUES (1, 100064);
INSERT INTO `sys_role_menu` VALUES (1, 100065);
INSERT INTO `sys_role_menu` VALUES (1, 100066);
INSERT INTO `sys_role_menu` VALUES (1, 100067);
INSERT INTO `sys_role_menu` VALUES (1, 100068);
INSERT INTO `sys_role_menu` VALUES (1, 100094);
INSERT INTO `sys_role_menu` VALUES (1, 100095);
INSERT INTO `sys_role_menu` VALUES (1, 100097);
INSERT INTO `sys_role_menu` VALUES (1, 100099);
INSERT INTO `sys_role_menu` VALUES (1, 100100);
INSERT INTO `sys_role_menu` VALUES (1, 100102);
INSERT INTO `sys_role_menu` VALUES (1, 100104);
INSERT INTO `sys_role_menu` VALUES (1, 100105);
INSERT INTO `sys_role_menu` VALUES (1, 100108);
INSERT INTO `sys_role_menu` VALUES (1, 100109);
INSERT INTO `sys_role_menu` VALUES (1, 100110);
INSERT INTO `sys_role_menu` VALUES (1, 100111);
INSERT INTO `sys_role_menu` VALUES (1, 100112);
INSERT INTO `sys_role_menu` VALUES (1, 100113);
INSERT INTO `sys_role_menu` VALUES (1, 100119);
INSERT INTO `sys_role_menu` VALUES (1, 100120);
INSERT INTO `sys_role_menu` VALUES (1, 100121);
INSERT INTO `sys_role_menu` VALUES (1, 100122);
INSERT INTO `sys_role_menu` VALUES (1, 100123);
INSERT INTO `sys_role_menu` VALUES (1, 100124);
INSERT INTO `sys_role_menu` VALUES (1, 100125);
INSERT INTO `sys_role_menu` VALUES (1, 100126);
INSERT INTO `sys_role_menu` VALUES (1, 100183);
INSERT INTO `sys_role_menu` VALUES (1, 100184);
INSERT INTO `sys_role_menu` VALUES (1, 100185);
INSERT INTO `sys_role_menu` VALUES (1, 100221);
INSERT INTO `sys_role_menu` VALUES (1, 100249);
INSERT INTO `sys_role_menu` VALUES (1, 100250);
INSERT INTO `sys_role_menu` VALUES (1, 100251);
INSERT INTO `sys_role_menu` VALUES (1, 100252);
INSERT INTO `sys_role_menu` VALUES (1, 100253);
INSERT INTO `sys_role_menu` VALUES (1, 100254);
INSERT INTO `sys_role_menu` VALUES (1, 100261);
INSERT INTO `sys_role_menu` VALUES (1, 100263);
INSERT INTO `sys_role_menu` VALUES (1, 100267);
INSERT INTO `sys_role_menu` VALUES (1, 100268);
INSERT INTO `sys_role_menu` VALUES (1, 100271);
INSERT INTO `sys_role_menu` VALUES (1, 100272);
INSERT INTO `sys_role_menu` VALUES (1, 100278);
INSERT INTO `sys_role_menu` VALUES (1, 100279);
INSERT INTO `sys_role_menu` VALUES (1, 100280);
INSERT INTO `sys_role_menu` VALUES (1, 100281);
INSERT INTO `sys_role_menu` VALUES (1, 100282);
INSERT INTO `sys_role_menu` VALUES (1, 100283);
INSERT INTO `sys_role_menu` VALUES (1, 100284);
INSERT INTO `sys_role_menu` VALUES (1, 100285);
INSERT INTO `sys_role_menu` VALUES (1, 100286);
INSERT INTO `sys_role_menu` VALUES (1, 100287);
INSERT INTO `sys_role_menu` VALUES (1, 100288);
INSERT INTO `sys_role_menu` VALUES (1, 100289);
INSERT INTO `sys_role_menu` VALUES (1, 100290);
INSERT INTO `sys_role_menu` VALUES (1, 100291);
INSERT INTO `sys_role_menu` VALUES (1, 100292);
INSERT INTO `sys_role_menu` VALUES (1, 100293);
INSERT INTO `sys_role_menu` VALUES (1, 100294);
INSERT INTO `sys_role_menu` VALUES (1, 100295);
INSERT INTO `sys_role_menu` VALUES (1, 100296);
INSERT INTO `sys_role_menu` VALUES (2, 100000);
INSERT INTO `sys_role_menu` VALUES (2, 100002);
INSERT INTO `sys_role_menu` VALUES (2, 100003);
INSERT INTO `sys_role_menu` VALUES (2, 100004);
INSERT INTO `sys_role_menu` VALUES (2, 100005);
INSERT INTO `sys_role_menu` VALUES (2, 100006);
INSERT INTO `sys_role_menu` VALUES (2, 100007);
INSERT INTO `sys_role_menu` VALUES (2, 100008);
INSERT INTO `sys_role_menu` VALUES (2, 100009);
INSERT INTO `sys_role_menu` VALUES (2, 100010);
INSERT INTO `sys_role_menu` VALUES (2, 100011);
INSERT INTO `sys_role_menu` VALUES (2, 100012);
INSERT INTO `sys_role_menu` VALUES (2, 100013);
INSERT INTO `sys_role_menu` VALUES (2, 100014);
INSERT INTO `sys_role_menu` VALUES (2, 100015);
INSERT INTO `sys_role_menu` VALUES (2, 100016);
INSERT INTO `sys_role_menu` VALUES (2, 100017);
INSERT INTO `sys_role_menu` VALUES (2, 100018);
INSERT INTO `sys_role_menu` VALUES (2, 100019);
INSERT INTO `sys_role_menu` VALUES (2, 100020);
INSERT INTO `sys_role_menu` VALUES (2, 100021);
INSERT INTO `sys_role_menu` VALUES (2, 100022);
INSERT INTO `sys_role_menu` VALUES (2, 100023);
INSERT INTO `sys_role_menu` VALUES (2, 100024);
INSERT INTO `sys_role_menu` VALUES (2, 100025);
INSERT INTO `sys_role_menu` VALUES (2, 100026);
INSERT INTO `sys_role_menu` VALUES (2, 100027);
INSERT INTO `sys_role_menu` VALUES (2, 100053);
INSERT INTO `sys_role_menu` VALUES (2, 100054);
INSERT INTO `sys_role_menu` VALUES (2, 100055);
INSERT INTO `sys_role_menu` VALUES (2, 100056);
INSERT INTO `sys_role_menu` VALUES (2, 100057);
INSERT INTO `sys_role_menu` VALUES (2, 100058);
INSERT INTO `sys_role_menu` VALUES (2, 100059);
INSERT INTO `sys_role_menu` VALUES (2, 100060);
INSERT INTO `sys_role_menu` VALUES (2, 100061);
INSERT INTO `sys_role_menu` VALUES (2, 100062);
INSERT INTO `sys_role_menu` VALUES (2, 100063);
INSERT INTO `sys_role_menu` VALUES (2, 100064);
INSERT INTO `sys_role_menu` VALUES (2, 100065);
INSERT INTO `sys_role_menu` VALUES (2, 100066);
INSERT INTO `sys_role_menu` VALUES (2, 100067);
INSERT INTO `sys_role_menu` VALUES (2, 100068);
INSERT INTO `sys_role_menu` VALUES (2, 100094);
INSERT INTO `sys_role_menu` VALUES (2, 100095);
INSERT INTO `sys_role_menu` VALUES (2, 100097);
INSERT INTO `sys_role_menu` VALUES (2, 100099);
INSERT INTO `sys_role_menu` VALUES (2, 100100);
INSERT INTO `sys_role_menu` VALUES (2, 100102);
INSERT INTO `sys_role_menu` VALUES (2, 100104);
INSERT INTO `sys_role_menu` VALUES (2, 100105);
INSERT INTO `sys_role_menu` VALUES (2, 100108);
INSERT INTO `sys_role_menu` VALUES (2, 100109);
INSERT INTO `sys_role_menu` VALUES (2, 100110);
INSERT INTO `sys_role_menu` VALUES (2, 100111);
INSERT INTO `sys_role_menu` VALUES (2, 100112);
INSERT INTO `sys_role_menu` VALUES (2, 100113);
INSERT INTO `sys_role_menu` VALUES (2, 100119);
INSERT INTO `sys_role_menu` VALUES (2, 100120);
INSERT INTO `sys_role_menu` VALUES (2, 100121);
INSERT INTO `sys_role_menu` VALUES (2, 100122);
INSERT INTO `sys_role_menu` VALUES (2, 100123);
INSERT INTO `sys_role_menu` VALUES (2, 100124);
INSERT INTO `sys_role_menu` VALUES (2, 100125);
INSERT INTO `sys_role_menu` VALUES (2, 100126);
INSERT INTO `sys_role_menu` VALUES (2, 100183);
INSERT INTO `sys_role_menu` VALUES (2, 100184);
INSERT INTO `sys_role_menu` VALUES (2, 100185);
INSERT INTO `sys_role_menu` VALUES (2, 100221);
INSERT INTO `sys_role_menu` VALUES (2, 100249);
INSERT INTO `sys_role_menu` VALUES (2, 100250);
INSERT INTO `sys_role_menu` VALUES (2, 100251);
INSERT INTO `sys_role_menu` VALUES (2, 100252);
INSERT INTO `sys_role_menu` VALUES (2, 100253);
INSERT INTO `sys_role_menu` VALUES (2, 100254);
INSERT INTO `sys_role_menu` VALUES (2, 100261);
INSERT INTO `sys_role_menu` VALUES (2, 100263);
INSERT INTO `sys_role_menu` VALUES (2, 100267);
INSERT INTO `sys_role_menu` VALUES (2, 100268);
INSERT INTO `sys_role_menu` VALUES (2, 100269);
INSERT INTO `sys_role_menu` VALUES (2, 100270);
INSERT INTO `sys_role_menu` VALUES (2, 100271);
INSERT INTO `sys_role_menu` VALUES (2, 100272);
INSERT INTO `sys_role_menu` VALUES (2, 100278);
INSERT INTO `sys_role_menu` VALUES (2, 100279);
INSERT INTO `sys_role_menu` VALUES (2, 100280);
INSERT INTO `sys_role_menu` VALUES (2, 100281);
INSERT INTO `sys_role_menu` VALUES (2, 100282);
INSERT INTO `sys_role_menu` VALUES (2, 100283);
INSERT INTO `sys_role_menu` VALUES (2, 100284);
INSERT INTO `sys_role_menu` VALUES (2, 100285);
INSERT INTO `sys_role_menu` VALUES (2, 100286);
INSERT INTO `sys_role_menu` VALUES (2, 100287);
INSERT INTO `sys_role_menu` VALUES (2, 100288);
INSERT INTO `sys_role_menu` VALUES (2, 100289);
INSERT INTO `sys_role_menu` VALUES (2, 100290);
INSERT INTO `sys_role_menu` VALUES (2, 100291);
INSERT INTO `sys_role_menu` VALUES (2, 100292);
INSERT INTO `sys_role_menu` VALUES (2, 100293);
INSERT INTO `sys_role_menu` VALUES (2, 100294);
INSERT INTO `sys_role_menu` VALUES (2, 100295);
INSERT INTO `sys_role_menu` VALUES (2, 100296);
INSERT INTO `sys_role_menu` VALUES (3, 100000);
INSERT INTO `sys_role_menu` VALUES (3, 100002);
INSERT INTO `sys_role_menu` VALUES (3, 100003);
INSERT INTO `sys_role_menu` VALUES (3, 100004);
INSERT INTO `sys_role_menu` VALUES (3, 100005);
INSERT INTO `sys_role_menu` VALUES (3, 100006);
INSERT INTO `sys_role_menu` VALUES (3, 100007);
INSERT INTO `sys_role_menu` VALUES (3, 100008);
INSERT INTO `sys_role_menu` VALUES (3, 100009);
INSERT INTO `sys_role_menu` VALUES (3, 100010);
INSERT INTO `sys_role_menu` VALUES (3, 100011);
INSERT INTO `sys_role_menu` VALUES (3, 100012);
INSERT INTO `sys_role_menu` VALUES (3, 100013);
INSERT INTO `sys_role_menu` VALUES (3, 100014);
INSERT INTO `sys_role_menu` VALUES (3, 100015);
INSERT INTO `sys_role_menu` VALUES (3, 100016);
INSERT INTO `sys_role_menu` VALUES (3, 100017);
INSERT INTO `sys_role_menu` VALUES (3, 100018);
INSERT INTO `sys_role_menu` VALUES (3, 100019);
INSERT INTO `sys_role_menu` VALUES (3, 100020);
INSERT INTO `sys_role_menu` VALUES (3, 100021);
INSERT INTO `sys_role_menu` VALUES (3, 100022);
INSERT INTO `sys_role_menu` VALUES (3, 100023);
INSERT INTO `sys_role_menu` VALUES (3, 100024);
INSERT INTO `sys_role_menu` VALUES (3, 100025);
INSERT INTO `sys_role_menu` VALUES (3, 100026);
INSERT INTO `sys_role_menu` VALUES (3, 100027);
INSERT INTO `sys_role_menu` VALUES (3, 100053);
INSERT INTO `sys_role_menu` VALUES (3, 100054);
INSERT INTO `sys_role_menu` VALUES (3, 100055);
INSERT INTO `sys_role_menu` VALUES (3, 100056);
INSERT INTO `sys_role_menu` VALUES (3, 100057);
INSERT INTO `sys_role_menu` VALUES (3, 100058);
INSERT INTO `sys_role_menu` VALUES (3, 100059);
INSERT INTO `sys_role_menu` VALUES (3, 100060);
INSERT INTO `sys_role_menu` VALUES (3, 100061);
INSERT INTO `sys_role_menu` VALUES (3, 100062);
INSERT INTO `sys_role_menu` VALUES (3, 100063);
INSERT INTO `sys_role_menu` VALUES (3, 100064);
INSERT INTO `sys_role_menu` VALUES (3, 100065);
INSERT INTO `sys_role_menu` VALUES (3, 100066);
INSERT INTO `sys_role_menu` VALUES (3, 100067);
INSERT INTO `sys_role_menu` VALUES (3, 100068);
INSERT INTO `sys_role_menu` VALUES (3, 100094);
INSERT INTO `sys_role_menu` VALUES (3, 100095);
INSERT INTO `sys_role_menu` VALUES (3, 100097);
INSERT INTO `sys_role_menu` VALUES (3, 100099);
INSERT INTO `sys_role_menu` VALUES (3, 100100);
INSERT INTO `sys_role_menu` VALUES (3, 100102);
INSERT INTO `sys_role_menu` VALUES (3, 100104);
INSERT INTO `sys_role_menu` VALUES (3, 100105);
INSERT INTO `sys_role_menu` VALUES (3, 100108);
INSERT INTO `sys_role_menu` VALUES (3, 100109);
INSERT INTO `sys_role_menu` VALUES (3, 100110);
INSERT INTO `sys_role_menu` VALUES (3, 100111);
INSERT INTO `sys_role_menu` VALUES (3, 100112);
INSERT INTO `sys_role_menu` VALUES (3, 100113);
INSERT INTO `sys_role_menu` VALUES (3, 100119);
INSERT INTO `sys_role_menu` VALUES (3, 100120);
INSERT INTO `sys_role_menu` VALUES (3, 100121);
INSERT INTO `sys_role_menu` VALUES (3, 100122);
INSERT INTO `sys_role_menu` VALUES (3, 100123);
INSERT INTO `sys_role_menu` VALUES (3, 100124);
INSERT INTO `sys_role_menu` VALUES (3, 100125);
INSERT INTO `sys_role_menu` VALUES (3, 100126);
INSERT INTO `sys_role_menu` VALUES (3, 100183);
INSERT INTO `sys_role_menu` VALUES (3, 100184);
INSERT INTO `sys_role_menu` VALUES (3, 100185);
INSERT INTO `sys_role_menu` VALUES (3, 100221);
INSERT INTO `sys_role_menu` VALUES (3, 100249);
INSERT INTO `sys_role_menu` VALUES (3, 100250);
INSERT INTO `sys_role_menu` VALUES (3, 100251);
INSERT INTO `sys_role_menu` VALUES (3, 100252);
INSERT INTO `sys_role_menu` VALUES (3, 100253);
INSERT INTO `sys_role_menu` VALUES (3, 100254);
INSERT INTO `sys_role_menu` VALUES (3, 100261);
INSERT INTO `sys_role_menu` VALUES (3, 100263);
INSERT INTO `sys_role_menu` VALUES (3, 100267);
INSERT INTO `sys_role_menu` VALUES (3, 100268);
INSERT INTO `sys_role_menu` VALUES (3, 100269);
INSERT INTO `sys_role_menu` VALUES (3, 100270);
INSERT INTO `sys_role_menu` VALUES (3, 100271);
INSERT INTO `sys_role_menu` VALUES (3, 100272);
INSERT INTO `sys_role_menu` VALUES (3, 100278);
INSERT INTO `sys_role_menu` VALUES (3, 100279);
INSERT INTO `sys_role_menu` VALUES (3, 100280);
INSERT INTO `sys_role_menu` VALUES (3, 100281);
INSERT INTO `sys_role_menu` VALUES (3, 100282);
INSERT INTO `sys_role_menu` VALUES (3, 100283);
INSERT INTO `sys_role_menu` VALUES (3, 100284);
INSERT INTO `sys_role_menu` VALUES (3, 100285);
INSERT INTO `sys_role_menu` VALUES (3, 100286);
INSERT INTO `sys_role_menu` VALUES (3, 100287);
INSERT INTO `sys_role_menu` VALUES (3, 100288);
INSERT INTO `sys_role_menu` VALUES (3, 100289);
INSERT INTO `sys_role_menu` VALUES (3, 100290);
INSERT INTO `sys_role_menu` VALUES (3, 100291);
INSERT INTO `sys_role_menu` VALUES (3, 100292);
INSERT INTO `sys_role_menu` VALUES (3, 100293);
INSERT INTO `sys_role_menu` VALUES (3, 100294);
INSERT INTO `sys_role_menu` VALUES (3, 100295);
INSERT INTO `sys_role_menu` VALUES (3, 100296);
INSERT INTO `sys_role_menu` VALUES (5, 100002);
INSERT INTO `sys_role_menu` VALUES (5, 100003);
INSERT INTO `sys_role_menu` VALUES (5, 100004);
INSERT INTO `sys_role_menu` VALUES (5, 100005);
INSERT INTO `sys_role_menu` VALUES (5, 100006);
INSERT INTO `sys_role_menu` VALUES (5, 100007);
INSERT INTO `sys_role_menu` VALUES (5, 100008);
INSERT INTO `sys_role_menu` VALUES (5, 100009);
INSERT INTO `sys_role_menu` VALUES (5, 100010);
INSERT INTO `sys_role_menu` VALUES (6, 100000);
INSERT INTO `sys_role_menu` VALUES (6, 100002);
INSERT INTO `sys_role_menu` VALUES (6, 100003);
INSERT INTO `sys_role_menu` VALUES (6, 100004);
INSERT INTO `sys_role_menu` VALUES (6, 100005);
INSERT INTO `sys_role_menu` VALUES (6, 100006);
COMMIT;

-- ----------------------------
-- Table structure for sys_template
-- ----------------------------
DROP TABLE IF EXISTS `sys_template`;
CREATE TABLE `sys_template` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '模板编号',
  `tpl_key` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '模板KEY',
  `lang_code` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '语言编码（如zh）',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '模板名称',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '模板内容',
  `type` tinyint(1) NOT NULL COMMENT '模板类型',
  `level` tinyint(1) NOT NULL COMMENT '级别',
  `route_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '路由地址',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态（0正常 1停用）',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '备注',
  `created_at` timestamp NOT NULL COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=302 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='模板信息';

-- ----------------------------
-- Records of sys_template
-- ----------------------------
BEGIN;
INSERT INTO `sys_template` VALUES (100, '', '', '提交申请', '您提交了一条{applyTypeName} \n 申请编号：{applyId} \n 申请人：{applyUserName} \n 事由：{reason}\n\n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:08:38', '2025-09-28 10:03:02');
INSERT INTO `sys_template` VALUES (101, '', '', '驳回', '您收到一条驳回的申请 \n 申请编号：{applyId} \n 申请人：{applyUserName} \n 事由：{reason}\n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:09:00', '2025-09-28 10:03:13');
INSERT INTO `sys_template` VALUES (102, '', '', '待审', '您有一条待审的申请 \n 申请编号：{applyId} \n 申请人：{applyUserName} \n 事由：{reason}\n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:09:00', '2025-09-28 10:03:13');
INSERT INTO `sys_template` VALUES (103, '', '', '催办', '您有一条待处理的申请，请尽快处理 \n 申请编号：{applyId} \n 申请人：{applyUserName} \n 事由：{reason}\n \n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:09:00', '2025-09-28 10:03:02');
INSERT INTO `sys_template` VALUES (104, '', '', '撤回', '您有一条待处理的申请，请尽快处理 \n 申请编号：{applyId} \n 申请人：{applyUserName} \n 事由：{reason}\n 时间：{time}\n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:09:00', '2025-09-28 10:03:02');
INSERT INTO `sys_template` VALUES (106, '', '', '已通过', '您的{applyTypeName}已通过 \n 申请编号：{applyId} \n 申请人：{applyUserName} \n 事由：{reason}\n\n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:09:00', '2025-09-28 10:03:02');
INSERT INTO `sys_template` VALUES (107, '', '', '抄送', '您有一条{applyTypeName}的抄送 \n 申请编号：{applyId} \n 申请人：{applyUserName} \n 事由：{reason}\n\n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:09:00', '2025-09-28 10:03:02');
INSERT INTO `sys_template` VALUES (108, '', '', '撤消', '您撤消了一条{applyTypeName} \n 申请编号：{applyId} \n 申请人：{applyUserName} \n 事由：{reason}\n\n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:09:00', '2025-09-28 10:03:02');
INSERT INTO `sys_template` VALUES (109, '', '', '变更', '您变更了一条{applyTypeName} \n 申请编号：{applyId} \n 申请人：{applyUserName} \n 事由：{reason}\n\n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:09:00', '2025-09-28 10:03:02');
INSERT INTO `sys_template` VALUES (110, '', '', '待付款', '您有一条待付款的申请\n 申请编号：{applyId} \n 申请人：{applyUserName} \n {reason}！\n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:09:00', '2025-09-28 10:03:02');
INSERT INTO `sys_template` VALUES (111, '', '', '已付款', '您有一条付款信息\n 申请编号：{applyId} \n 申请人：{applyUserName} \n {reason}的报销已付款，请注意查收！\n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:09:00', NULL);
INSERT INTO `sys_template` VALUES (201, '', '', '备案', '您有一条需备案的申请 {applyTypeName} \n 申请编号：{applyId} \n 申请人：{applyUserName} \n 事由：{reason}\n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:09:00', NULL);
INSERT INTO `sys_template` VALUES (202, '', '', '已备案通过', '您有一条申请已备案通过 {applyTypeName} \n 申请编号：{applyId} \n 申请人：{applyUserName} \n 事由：{reason}\n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:09:00', NULL);
INSERT INTO `sys_template` VALUES (203, '', '', '待确认', '您有一条待确认的申请\n 申请编号：{applyId} \n 申请人：{applyUserName} \n {reason}！\n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:09:00', NULL);
INSERT INTO `sys_template` VALUES (204, '', '', '借款已付款\n', '您的借款已付款\n 申请编号：{applyId} \n 申请人：{applyUserName} \n {reason}！\n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:09:00', NULL);
INSERT INTO `sys_template` VALUES (205, '', '', '报销申请已付款\n', '您的报销申请已付款\n 申请编号：{applyId} \n 申请人：{applyUserName} \n {reason}！\n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:09:00', NULL);
INSERT INTO `sys_template` VALUES (206, '', '', '报销申请已退款', '您的报销申请已退款\n 申请编号：{applyId} \n 申请人：{applyUserName} \n {reason}！\n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:09:00', NULL);
INSERT INTO `sys_template` VALUES (207, '', '', '待付款的申请\n', '您有一条待付款的申请\n 申请编号：{applyId} \n 申请人：{applyUserName} \n {reason}！\n <a href=\"https://secapi.polybigdata.cn/oa/get/corp/weixin/code/{wxSuiteId}?redirectUrl=pages/record/travelDetail&curType={curType}&isApprove={isApprove}&isReim={isReim}&redirectId={applyId}&formId={formId}&suiteId={wxSuiteId}\">点击前往详情</a>', 4, 1, NULL, 1, '', '2023-05-30 10:09:00', NULL);
INSERT INTO `sys_template` VALUES (300, '', '', '劳务合同到期提醒', '{}劳务合同还有{}天到期!', 4, 1, NULL, 1, '', '2023-07-28 10:09:00', NULL);
INSERT INTO `sys_template` VALUES (301, '', '', 'EasyAdmin重置密码', '您的验证码：${code}', 1, 1, NULL, 1, '', '2025-08-23 11:49:15', NULL);
COMMIT;

-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户编号',
  `username` varchar(32) DEFAULT NULL COMMENT '用户名',
  `nickname` varchar(50) DEFAULT NULL COMMENT '呢称',
  `type` tinyint DEFAULT NULL COMMENT '用户类型',
  `email` varchar(64) DEFAULT NULL COMMENT '邮箱',
  `mobile` varchar(11) DEFAULT NULL COMMENT '手机号',
  `sex` tinyint DEFAULT NULL COMMENT '性别（0未知 1男 2女）',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
  `autograph` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '电子签名',
  `password` varchar(128) DEFAULT NULL COMMENT '密码',
  `status` tinyint DEFAULT NULL COMMENT '状态（0正常 1停用）',
  `created_by` bigint DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_by` bigint DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户信息表';

-- ----------------------------
-- Records of sys_user
-- ----------------------------
BEGIN;
INSERT INTO `sys_user` VALUES (1, 'admin', 'Admin', 0, '1073602@qq.com', '', 0, '', '', '$2a$10$XqznnjX9E8V1GXRtDzjjOuq241FotnGPD3jLWOpLMK.Fw8GbdSE3y', 0, 1, '2025-04-13 16:51:52', 1, '2026-04-28 10:59:57');
INSERT INTO `sys_user` VALUES (2, 'develop', 'Develop', 0, 'looses118@gmail.com', '', 0, '', '', '$2a$10$eS0i25p/wMfdmdhOxWryMuu2ZjSgYn6QyvSQLsQxxmXMhLMrSv0aC', 0, 1, '2025-04-13 16:51:52', 1, '2026-04-28 18:19:59');
INSERT INTO `sys_user` VALUES (3, 'test', 'Test', 0, '', '', 0, '', '', '$2a$10$caTsiM5lj44j7X2wOvWUueOHgE.4B5Eq9nx6KAfwybk7i3caE.85S', 0, 1, '2025-04-13 16:51:52', 1, '2026-04-20 16:04:26');
INSERT INTO `sys_user` VALUES (5, 'agang', '阿刚', 0, '', '', 0, '', '', '$2a$10$caTsiM5lj44j7X2wOvWUueOHgE.4B5Eq9nx6KAfwybk7i3caE.85S', 0, 1, '2025-04-13 16:51:52', 1, '2026-04-20 16:04:27');
INSERT INTO `sys_user` VALUES (6, 'lisi', '李四', 0, '', '', 1, '', '', '$2a$10$caTsiM5lj44j7X2wOvWUueOHgE.4B5Eq9nx6KAfwybk7i3caE.85S', 0, 1, '2025-04-13 16:51:52', 1, '2026-04-20 16:04:28');
COMMIT;

-- ----------------------------
-- Table structure for sys_user_credential
-- ----------------------------
DROP TABLE IF EXISTS `sys_user_credential`;
CREATE TABLE `sys_user_credential` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '内部主键',
  `user_id` bigint NOT NULL COMMENT '系统用户ID（关联用户表）',
  `username` varchar(64) NOT NULL COMMENT '冗余用户名（用于快速查询/审计）',
  `credential_id` varbinary(255) NOT NULL COMMENT '凭证ID（二进制存储，WebAuthn原始字节）',
  `public_key_cose` varbinary(512) NOT NULL COMMENT 'COSE格式公钥（二进制）',
  `user_handle` varbinary(64) NOT NULL COMMENT 'WebAuthn用户标识（稳定且不可变）',
  `signature_count` bigint unsigned NOT NULL DEFAULT '0' COMMENT '签名计数器（防克隆攻击）',
  `aaguid` char(36) DEFAULT NULL COMMENT '认证器设备类型ID',
  `attestation_type` varchar(32) DEFAULT NULL COMMENT '证明类型（basic/self/none等）',
  `attachment` varchar(32) DEFAULT NULL COMMENT '认证器类型（platform/cross-platform）',
  `transports` varchar(255) DEFAULT NULL COMMENT '认证器传输方式（usb,nfc,ble,internal）',
  `device_type` varchar(32) DEFAULT NULL COMMENT '设备类型（single-device / multi-device）',
  `backup_state` tinyint(1) DEFAULT NULL COMMENT '是否支持云同步（passkey关键属性）',
  `device_name` varchar(255) DEFAULT NULL COMMENT '用户自定义设备名称',
  `is_revoked` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已吊销（禁用凭证）',
  `last_used_at` timestamp NULL DEFAULT NULL COMMENT '最后使用时间',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_credential_id` (`credential_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户Passkey凭据';

-- ----------------------------
-- Records of sys_user_credential
-- ----------------------------
BEGIN;
INSERT INTO `sys_user_credential` VALUES (7, 2, 'develop', 0x032E86B4C9301ED58D832F96DC1C8BBF, 0xA501020326200121582037DE0D8711752F9E8E65FC569FEF86C64B1DDCD950F6DC8CBB1714BED85F0C1E225820FE478D6FF055A0FACB668278FABEC63093F8E45A8BE2203CCDFC66208281570D, 0x32, 0, 'ea9b8d664d011d213ce4b6b48cb575d4', 'none', '', 'hybrid,internal', '', 1, '', 0, '2026-05-18 18:02:29', '2026-04-26 16:30:18', '2026-05-18 18:02:29');
COMMIT;

-- ----------------------------
-- Table structure for sys_user_dept
-- ----------------------------
DROP TABLE IF EXISTS `sys_user_dept`;
CREATE TABLE `sys_user_dept` (
  `user_id` bigint NOT NULL COMMENT '用户编号',
  `dept_id` bigint NOT NULL COMMENT '部门编号',
  `post_id` bigint NOT NULL COMMENT '岗位编号',
  PRIMARY KEY (`user_id`,`dept_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户部门';

-- ----------------------------
-- Records of sys_user_dept
-- ----------------------------
BEGIN;
INSERT INTO `sys_user_dept` VALUES (1, 100, 1);
INSERT INTO `sys_user_dept` VALUES (2, 112, 4);
INSERT INTO `sys_user_dept` VALUES (3, 105, 6);
COMMIT;

-- ----------------------------
-- Table structure for sys_user_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_user_role`;
CREATE TABLE `sys_user_role` (
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`user_id`,`role_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户和角色关联表';

-- ----------------------------
-- Records of sys_user_role
-- ----------------------------
BEGIN;
INSERT INTO `sys_user_role` VALUES (1, 1);
INSERT INTO `sys_user_role` VALUES (2, 3);
INSERT INTO `sys_user_role` VALUES (3, 2);
INSERT INTO `sys_user_role` VALUES (5, 2);
INSERT INTO `sys_user_role` VALUES (6, 3);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
