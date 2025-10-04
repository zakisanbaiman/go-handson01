create table `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ユーザーの識別子',
    `name` VARCHAR(20) NOT NULL COMMENT 'ユーザー名',
    `password` VARCHAR(80) NOT NULL COMMENT 'パスワード',
    `role` VARCHAR(80) NOT NULL COMMENT 'ロール',
    `created_at` DATETIME(6) NOT NULL COMMENT '作成日時',
    `modified_at` DATETIME(6) NOT NULL COMMENT '更新日時',
    PRIMARY KEY (`id`),
    UNIQUE KEY `name_unique` (`name`) USING BTREE
);

create table `tasks` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'タスクの識別子',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'ユーザーの識別子',
    `title` VARCHAR(128) NOT NULL COMMENT 'タスクのタイトル',
    `status` VARCHAR(20) NOT NULL COMMENT 'タスクのステータス',
    `created_at` DATETIME(6) NOT NULL COMMENT '作成日時',
    `modified_at` DATETIME(6) NOT NULL COMMENT '更新日時',
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_user_id`
        FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) 
        ON DELETE RESTRICT ON UPDATE RESTRICT
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='タスク';
