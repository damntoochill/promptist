CREATE TABLE `auth_users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `full_name` varchar(255) NOT NULL,
  `username` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `profile_photo` varchar(255) DEFAULT NULL,
  `forgot_password_token` varchar(255) DEFAULT NULL,
  `forgot_password_expiry` varchar(255) DEFAULT NULL,
  `forgot_password_requests` int NOT NULL DEFAULT '0',
  `email_verification_token` varchar(255) DEFAULT NULL,
  `email_verification_expiry` datetime DEFAULT NULL,
  `has_verified_email` tinyint(1) NOT NULL DEFAULT '0',
  `email_verification_requests` int NOT NULL DEFAULT '0',
  `is_staff` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
);


CREATE TABLE `follows` (
  `id` int NOT NULL AUTO_INCREMENT,
  `leader_id` int NOT NULL,
  `follower_id` int NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `leader_id` (`leader_id`),
  KEY `follower_id` (`follower_id`),
  CONSTRAINT `follows_ibfk_1` FOREIGN KEY (`leader_id`) REFERENCES `users` (`id`),
  CONSTRAINT `follows_ibfk_2` FOREIGN KEY (`follower_id`) REFERENCES `users` (`id`)
);