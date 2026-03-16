CREATE TABLE `friend_requests` (
                                   `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                   `user_id` int(10) unsigned NOT NULL COMMENT '发起申请的用户id',
                                   `req_uid` int(10) unsigned NOT NULL COMMENT '被申请的好友uid（目标用户）',
                                   `req_msg` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '申请附言/验证消息',
                                   `handle_result` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '0=待处理 1=通过 2=拒绝 3=申请人取消',
                                   `handle_msg` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '处理时的回复/拒绝理由',
                                   `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '申请创建时间',
                                   `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
                                   `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
                                   PRIMARY KEY (`id`),
                                   UNIQUE KEY `uk_user_target` (`user_id`,`req_uid`),
                                   KEY `idx_user_created` (`user_id`,`created_at`),
                                   KEY `idx_target_status_time` (`req_uid`,`handle_result`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='好友申请记录表（单向申请）';

CREATE TABLE `friends` (
                           `id`          int(10) unsigned NOT NULL AUTO_INCREMENT,
                           `user_id`     int(10) unsigned NOT NULL COMMENT '用户id',
                           `friend_uid`  int(10) unsigned NOT NULL COMMENT '好友uid',
                           `remark`      varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注（对好友的备注）',
                           `add_source`  tinyint(1) unsigned NOT NULL DEFAULT 1 COMMENT '添加方式：1搜索 2名片 3群聊 ...',
                           `created_at`  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '成为好友时间',
                           `updated_at`  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
                           `deleted_at`  timestamp NULL DEFAULT NULL COMMENT '删除/拉黑时间（软删除）',

                           PRIMARY KEY (`id`),
                           UNIQUE KEY `uk_user_friend` (`user_id`, `friend_uid`),
                           KEY `idx_user_created`     (`user_id`, `created_at`),
                           KEY `idx_friend_created`   (`friend_uid`, `created_at`)

) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='好友关系表（单向记录，双向好友需两条记录）';

CREATE TABLE `groups` (
                          `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                          `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '群名',
                          `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '群图标',
                          `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1开启 0关闭',
                          `creator_uid` int(11) NOT NULL DEFAULT '0' COMMENT '创建人uid',
                          `group_type` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1=普通群 2=企业群 3=聊天室',
                          `is_verify` tinyint(1) NOT NULL COMMENT '入群验证：1开启 2关闭',
                          `notification` text CHARACTER SET utf8 NOT NULL COMMENT '群公告',
                          `notification_uid` int(11) NOT NULL DEFAULT '0' COMMENT '最后设置公告的人uid（可选）',
                          `member_count` mediumint(6) NOT NULL DEFAULT '1' COMMENT '群人数',
                          `created_at` timestamp NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '创建时间',
                          `updated_at` timestamp NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                          `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                          PRIMARY KEY (`id`) USING BTREE,
                          KEY `idx_creator_uid` (`creator_uid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天群表';

CREATE TABLE `group_requests` (
                                  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                  `req_id` int(10) unsigned NOT NULL COMMENT '申请人/被邀请人 uid',
                                  `group_id` int(10) unsigned NOT NULL COMMENT '群id',
                                  `req_msg` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '申请/邀请附言',
                                  `join_source` tinyint(1) unsigned NOT NULL DEFAULT '2' COMMENT '1=被邀请入群 2=主动申请',
                                  `inviter_user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '邀请人uid（join_source=1时有效）',
                                  `handle_user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '处理人uid（群管理员/群主）',
                                  `handle_result` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '1=待处理 2=通过 3=拒绝 4=取消',
                                  `handled_at` timestamp NULL DEFAULT NULL COMMENT '处理时间',
                                  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                  PRIMARY KEY (`id`),
                                  UNIQUE KEY `uk_group_req` (`group_id`,`req_id`),
                                  KEY `idx_group_pending` (`group_id`,`handle_result`,`created_at`),
                                  KEY `idx_requester_status` (`req_id`,`handle_result`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='入群申请/邀请记录表';

CREATE TABLE `group_members` (
                                 `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                                 `group_id` int(10) unsigned NOT NULL COMMENT '群id',
                                 `user_id` int(10) unsigned NOT NULL COMMENT '用户uid',
                                 `role_level` tinyint(1) unsigned NOT NULL DEFAULT '3' COMMENT '0普通成员，10管理员，20群主',
                                 `join_time` timestamp NULL DEFAULT NULL COMMENT '入群时间',
                                 `join_source` tinyint(1) unsigned NOT NULL DEFAULT '2' COMMENT '入群方式：1=被邀请 2=主动申请通过',
                                 `inviter_uid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '邀请人uid（join_source=1时有效）',
                                 `last_operator_uid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '最后操作人uid（添加/修改/移除）',
                                 `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1正常，2禁言，3被踢',
                                 `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
                                 `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
                                 `deleted_at` timestamp NULL DEFAULT NULL COMMENT '退出/被踢时间（软删除）',
                                 PRIMARY KEY (`id`),
                                 UNIQUE KEY `uk_group_user` (`group_id`,`user_id`,`deleted_at`) USING BTREE,
                                 KEY `idx_user_group` (`user_id`,`group_id`),
                                 KEY `idx_group_role` (`group_id`,`role_level`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='群成员关系表';