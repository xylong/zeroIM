CREATE TABLE `users` (
                         `id` varchar(24) COLLATE utf8mb4_unicode_ci NOT NULL,
                         `avatar` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '头像',
                         `nickname` varchar(24) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '昵称',
                         `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '手机号',
                         `password` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '密码',
                         `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：1正常 2禁用',
                         `sex` tinyint(4) NOT NULL DEFAULT '0' COMMENT '性别：1男 2女 3未知',
                         `created_at` timestamp NOT NULL COMMENT '创建时间',
                         `updated_at` timestamp NOT NULL COMMENT '更新时间',
                         `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                         PRIMARY KEY (`id`),
                         UNIQUE KEY `uni_phone_creeated_at` (`phone`,`created_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `wuid` (
                        `h` int(10) NOT NULL AUTO_INCREMENT,
                        `x` tinyint(4) NOT NULL DEFAULT '0',
                        PRIMARY KEY (`x`),
                        UNIQUE KEY `h` (`h`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=latin1;