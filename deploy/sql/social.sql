CREATE TABLE `friend_requests` (
                                   `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                                   `user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户id',
                                   `req_uid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '申请好友id',
                                   `req_msg` varchar(255) NOT NULL DEFAULT '' COMMENT '请求信息',
                                   `handle_result` tinyint(4) NOT NULL DEFAULT '1' COMMENT '处理结果：1-未处理 2-通过 3-拒绝 4-取消 ',
                                   `handle_msg` varchar(255) NOT NULL DEFAULT '' COMMENT '处理结果信息',
                                   `handled_at` timestamp NULL DEFAULT NULL COMMENT '处理时间',
                                   `created_at` timestamp NOT NULL COMMENT '创建时间',
                                   `updated_at` timestamp NOT NULL COMMENT '更新时间',
                                   `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                   PRIMARY KEY (`id`),
                                   KEY `idx_user_id_created_at` (`user_id`,`created_at`),
                                   KEY `idx_req_uid_created_at` (`req_uid`,`created_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='好友申请表';

CREATE TABLE `friends` (
                           `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                           `user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户?',
                           `friend_uid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '好友uid',
                           `remark` varchar(255) NOT NULL COMMENT '备注',
                           `add_source` tinyint(4) NOT NULL DEFAULT '1' COMMENT '添加方式',
                           `created_at` timestamp NOT NULL COMMENT '创建时间',
                           `updated_at` timestamp NOT NULL COMMENT '更新时间',
                           `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                           PRIMARY KEY (`id`),
                           UNIQUE KEY `uk_user_friend` (`user_id`,`friend_uid`),
                           KEY `idx_friend_uid` (`friend_uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='好友表';

CREATE TABLE `groups` (
                          `id` varchar(24) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
                          `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '群名',
                          `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '1' COMMENT '群图标',
                          `status` tinyint(4) NOT NULL COMMENT '状态：1正常 2解散',
                          `creator_uid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '创建人uid',
                          `group_type` int(11) NOT NULL COMMENT '群类型',
                          `is_verify` tinyint(1) NOT NULL COMMENT '入群验证：1开启 2关闭',
                          `notification` varchar(255) NOT NULL DEFAULT '' COMMENT '通知',
                          `notification_uid` varchar(64) NOT NULL DEFAULT '' COMMENT '通知人uid',
                          `created_at` timestamp NOT NULL COMMENT '创建时间',
                          `updated_at` timestamp NOT NULL COMMENT '更新时间',
                          `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                          PRIMARY KEY (`id`),
                          UNIQUE KEY `uk_creator_name` (`name`,`creator_uid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='聊天群表';

CREATE TABLE `group_requests` (
                                  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                                  `req_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '申请人uid',
                                  `group_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '群id',
                                  `req_msg` varchar(255) NOT NULL COMMENT '申请信息',
                                  `join_source` tinyint(4) NOT NULL DEFAULT '2' COMMENT '入群方式：1邀请 2申请',
                                  `inviter_user_id` varchar(64) NOT NULL COMMENT '邀请人uid',
                                  `handle_user_id` varchar(64) NOT NULL DEFAULT '' COMMENT '处理人uid',
                                  `handle_time` timestamp NULL DEFAULT NULL COMMENT '处理时间',
                                  `handle_result` tinyint(4) NOT NULL COMMENT '1未处理 2通过 3拒绝 4取消',
                                  `created_at` timestamp NOT NULL COMMENT '创建时间',
                                  `updated_at` timestamp NOT NULL COMMENT '更新时间',
                                  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                  PRIMARY KEY (`id`),
                                  KEY `idx_group_created` (`group_id`,`created_at`),
                                  KEY `idx_req_created` (`req_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='入群申请表';

CREATE TABLE `group_members` (
                                 `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                                 `group_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '群id',
                                 `user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户uid',
                                 `role_level` tinyint(4) NOT NULL DEFAULT '3' COMMENT '1.创建者 2.管理者 3.普通',
                                 `join_time` timestamp NULL DEFAULT NULL COMMENT '入群时间',
                                 `join_source` tinyint(4) DEFAULT '1' COMMENT '入群方式：1.邀请，2.申请\n',
                                 `inviter_uid` varchar(64) DEFAULT '' COMMENT '邀请人uid',
                                 `operator_uid` varchar(64) DEFAULT '' COMMENT '操作人uid',
                                 `created_at` timestamp NOT NULL COMMENT '创建时间',
                                 `updated_at` timestamp NOT NULL COMMENT '更新时间',
                                 `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                 PRIMARY KEY (`id`),
                                 UNIQUE KEY `uk_group_user` (`group_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='群成员表';