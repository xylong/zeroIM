CREATE TABLE `users` (
                         `id` varchar(24) COLLATE utf8mb4_unicode_ci  NOT NULL ,
                         `avatar` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
                         `nickname` varchar(24) COLLATE utf8mb4_unicode_ci NOT NULL,
                         `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
                         `password` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                         `status` tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                         `sex` tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                         `created_at` timestamp NULL DEFAULT NULL,
                         `updated_at` timestamp NULL DEFAULT NULL,
                         PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `wuid` (
                        `h` int(10) NOT NULL AUTO_INCREMENT,
                        `x` tinyint(4) NOT NULL DEFAULT '0',
                        PRIMARY KEY (`x`),
                        UNIQUE KEY `h` (`h`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=latin1;

CREATE TABLE `friends` (
                           `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                           `user_id` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
                           `friend_uid` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
                           `remark` varchar(255) DEFAULT NULL,
                           `add_source`  tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                           `created_at` timestamp NULL DEFAULT NULL,
                           PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '好友表';

CREATE TABLE `friend_requests` (
                                   `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                                   `user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户id',
                                   `req_uid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '申请好友id',
                                   `req_msg` varchar(255) NOT NULL DEFAULT '' COMMENT '请求信息',
                                   `req_time` timestamp NOT NULL COMMENT '请求时间',
                                   `handle_result` tinyint(4) NOT NULL DEFAULT '1' COMMENT '处理结果：1-未处理 2-通过 3-拒绝 4-取消 ',
                                   `handle_msg` varchar(255) NOT NULL DEFAULT '' COMMENT '处理结果信息',
                                   `handled_at` timestamp NULL DEFAULT NULL COMMENT '处理时间',
                                   `created_at` timestamp NOT NULL COMMENT '创建时间',
                                   `updated_at` timestamp NOT NULL COMMENT '更新时间',
                                   `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                                   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='好友申请表'

CREATE TABLE `groups` (
                          `id` varchar(24) COLLATE utf8mb4_unicode_ci  NOT NULL ,
                          `name` varchar(255) COLLATE utf8mb4_unicode_ci  NOT NULL ,
                          `icon` varchar(255) COLLATE utf8mb4_unicode_ci  NOT NULL ,
                          `status`  tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                          `creator_uid` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
                          `group_type` int(11) NOT NULL ,
                          `is_verify` boolean NOT NULL ,
                          `notification` varchar(255) DEFAULT NULL,
                          `notification_uid` varchar(64) DEFAULT NULL,
                          `created_at` timestamp NULL DEFAULT NULL,
                          `updated_at` timestamp NULL DEFAULT NULL,
                          PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '聊天群表';

CREATE TABLE `group_members` (
                                 `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                                 `group_id` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
                                 `user_id` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
                                 `role_level`  tinyint COLLATE utf8mb4_unicode_ci NOT NULL ,
                                 `join_time` timestamp NULL DEFAULT NULL,
                                 `join_source`  tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                                 `inviter_uid` varchar(64) DEFAULT NULL,
                                 `operator_uid` varchar(64) DEFAULT NULL,
                                 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '群成员表';

CREATE TABLE `group_requests` (
                                  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                                  `req_id` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
                                  `group_id` varchar(64) COLLATE utf8mb4_unicode_ci  NOT NULL ,
                                  `req_msg` varchar(255) DEFAULT NULL,
                                  `req_time` timestamp NULL DEFAULT NULL,
                                  `join_source`  tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                                  `inviter_user_id` varchar(64) DEFAULT NULL,
                                  `handle_user_id` varchar(64) DEFAULT NULL,
                                  `handle_time` timestamp NULL DEFAULT NULL,
                                  `handle_result`  tinyint COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                                  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '入群申请表';