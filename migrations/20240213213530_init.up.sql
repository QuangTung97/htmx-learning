CREATE TABLE `users`
(
    `id`            INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,

    `provider`      VARCHAR(50)      NOT NULL,
    `oauth_user_id` VARCHAR(255)     NOT NULL,
    `email`         VARCHAR(255)     NOT NULL,
    `image_url`     VARCHAR(1024)    NOT NULL,

    `status`        TINYINT UNSIGNED NOT NULL,

    `created_at`    TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`    TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY `uk_provider_user_id` (`provider`, `oauth_user_id`)
);


CREATE TABLE `user_sessions`
(
    `user_id`    INT UNSIGNED,
    `session_id` VARCHAR(512)     NOT NULL,

    `status`     TINYINT UNSIGNED NOT NULL,

    `created_at` TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (`user_id`, `session_id`)
);