
-- 设置【上级父菜单ID】（请修改为你实际的父目录ID）
SET @parent_menu_id = {{.MenuId}};



INSERT INTO `sys_menu` (`parent_id`, `type`, `sort`, `path`, `component`, `query`, `is_frame`, `visible`, `status`, `permission`, `icon`, `created_by`, `created_at`, `updated_by`, `updated_at`
) VALUES (@parent_menu_id,2,1,'/{{.ModuleName}}/{{.CleanName}}','{{.ModuleName}}/{{.CleanName}}/index','',0,1,0,'{{.Permission}}$list','file',1,NOW(),1,NOW());

SET @current_id = LAST_INSERT_ID();
INSERT INTO `sys_menu_tl` (`menu_id`, `lang_code`, `title`) VALUES (@current_id, 'zh-CN', '{{.Comment}}');
INSERT INTO `sys_menu_tl` (`menu_id`, `lang_code`, `title`) VALUES (@current_id, 'en-US', '{{.ClassName}}');

-- 查询按钮
INSERT INTO `sys_menu` (`parent_id`, `type`, `sort`, `path`, `component`, `query`, `is_frame`, `visible`, `status`, `permission`, `icon`, `created_by`, `created_at`, `updated_by`, `updated_at`)
VALUES (@current_id, 3, 1, '', '', '', 0, 1, 0, '{{.Permission}}$query', '', 1, NOW(), 1, NOW());
SET @query_id = LAST_INSERT_ID();
INSERT INTO `sys_menu_tl` (`menu_id`, `lang_code`, `title`) VALUES (@query_id, 'zh-CN', '查询');
INSERT INTO `sys_menu_tl` (`menu_id`, `lang_code`, `title`) VALUES (@query_id, 'en-US', 'Query');

-- 新增按钮
INSERT INTO `sys_menu` (`parent_id`,`type`, `sort`, `path`, `component`, `query`, `is_frame`, `visible`, `status`, `permission`, `icon`, `created_by`, `created_at`, `updated_by`, `updated_at`)
VALUES (@current_id, 3, 2, '', '', '', 0, 1, 0,'{{.Permission}}$add', '', 1, NOW(), 1, NOW());
SET @add_id = LAST_INSERT_ID();
INSERT INTO `sys_menu_tl` (`menu_id`, `lang_code`, `title`) VALUES (@add_id, 'zh-CN', '新增');
INSERT INTO `sys_menu_tl` (`menu_id`, `lang_code`, `title`) VALUES (@add_id, 'en-US', 'Add');

-- 修改按钮
INSERT INTO `sys_menu` (`parent_id`, `type`, `sort`, `path`, `component`, `query`, `is_frame`, `visible`, `status`, `permission`, `icon`, `created_by`, `created_at`, `updated_by`, `updated_at`)
VALUES (@current_id, 3, 3, '', '', '', 0, 1, 0,'{{.Permission}}$edit', '', 1, NOW(), 1, NOW());
SET @edit_id = LAST_INSERT_ID();
INSERT INTO `sys_menu_tl` (`menu_id`, `lang_code`, `title`) VALUES (@edit_id, 'zh-CN', '修改');
INSERT INTO `sys_menu_tl` (`menu_id`, `lang_code`, `title`) VALUES (@edit_id, 'en-US', 'Edit');

-- 删除按钮
INSERT INTO `sys_menu` (`parent_id`, `type`, `sort`, `path`, `component`, `query`, `is_frame`, `visible`, `status`, `permission`, `icon`, `created_by`, `created_at`, `updated_by`, `updated_at`)
VALUES (@current_id, 3, 4, '', '', '', 0, 1, 0, '{{.Permission}}$delete', '', 1, NOW(), 1, NOW());
SET @delete_id = LAST_INSERT_ID();
INSERT INTO `sys_menu_tl` (`menu_id`, `lang_code`, `title`) VALUES (@delete_id, 'zh-CN', '删除');
INSERT INTO `sys_menu_tl` (`menu_id`, `lang_code`, `title`) VALUES (@delete_id, 'en-US', 'Delete');
