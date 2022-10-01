-- -------------------------------------------------------------
-- TablePlus 4.8.2(436)
--
-- https://tableplus.com/
--
-- Database: art
-- Generation Time: 2022-09-08 13:24:10.7710
-- -------------------------------------------------------------


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


CREATE TABLE "art_likes" (
  "id" int NOT NULL AUTO_INCREMENT,
  "piece_id" int NOT NULL,
  "user_id" int NOT NULL,
  "created_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);

CREATE TABLE "art_pieces" (
  "piece_id" int NOT NULL AUTO_INCREMENT,
  "user_id" int DEFAULT NULL,
  "created_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "name" varchar(255) DEFAULT NULL,
  "description" text,
  "slug" varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  "image_uuid" varchar(255) DEFAULT NULL,
  "is_draft" tinyint NOT NULL DEFAULT '1',
  "prompt" text,
  "likes" int NOT NULL DEFAULT '0',
  "views" int NOT NULL DEFAULT '0',
  "saves" int NOT NULL DEFAULT '0',
  "comments" int NOT NULL DEFAULT '0',
  "program_id" int NOT NULL DEFAULT '0',
  "full_name" varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL DEFAULT '',
  "username" varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL DEFAULT '',
  "profile_photo_uuid" varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  "tags_literal" text CHARACTER SET latin1 COLLATE latin1_swedish_ci,
  "program_slug" varchar(255) DEFAULT NULL,
  "program_name" varchar(255) DEFAULT NULL,
  "program_cover_image_uuid" varchar(255) DEFAULT NULL,
  PRIMARY KEY ("piece_id")
);

CREATE TABLE "art_pieces_tags" (
  "id" int NOT NULL AUTO_INCREMENT,
  "tag_name" varchar(255) NOT NULL,
  "piece_id" varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  PRIMARY KEY ("id")
);

CREATE TABLE "art_programs" (
  "program_id" int NOT NULL AUTO_INCREMENT,
  "name" varchar(255) NOT NULL,
  "description" text NOT NULL,
  "slug" varchar(255) NOT NULL,
  PRIMARY KEY ("program_id")
);

CREATE TABLE "art_tags" (
  "tag_id" int NOT NULL AUTO_INCREMENT,
  "name" varchar(255) NOT NULL,
  "description" varchar(255) DEFAULT NULL,
  "total" int NOT NULL DEFAULT '1',
  PRIMARY KEY ("tag_id")
);

CREATE TABLE "auth_users" (
  "id" int NOT NULL AUTO_INCREMENT,
  "email" varchar(255) NOT NULL,
  "password_hash" varchar(255) NOT NULL,
  "created_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "forgot_password_token" varchar(255) DEFAULT NULL,
  "forgot_password_expiry" varchar(255) DEFAULT NULL,
  "forgot_password_requests" int NOT NULL DEFAULT '0',
  "is_verified" tinyint(1) NOT NULL DEFAULT '0',
  "is_admin" tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY ("id")
);

CREATE TABLE "auth_verifications" (
  "id" int NOT NULL AUTO_INCREMENT,
  "user_id" int NOT NULL,
  "token" varchar(255) NOT NULL,
  "expiry" datetime NOT NULL,
  PRIMARY KEY ("id")
);

CREATE TABLE "collection_collections" (
  "collection_id" int NOT NULL AUTO_INCREMENT,
  "name" varchar(255) NOT NULL,
  "description" varchar(255) DEFAULT NULL,
  "created_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "user_id" int NOT NULL,
  "is_public" tinyint NOT NULL DEFAULT '0',
  "num_pieces" int NOT NULL DEFAULT '0',
  PRIMARY KEY ("collection_id")
);

CREATE TABLE "collection_pieces" (
  "id" int NOT NULL AUTO_INCREMENT,
  "collection_id" int NOT NULL,
  "piece_id" int NOT NULL,
  PRIMARY KEY ("id")
);

CREATE TABLE "comment_comments" (
  "comment_id" int NOT NULL AUTO_INCREMENT,
  "created_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "body" text NOT NULL,
  "num_replies" int NOT NULL DEFAULT '0',
  "num_likes" int DEFAULT '0',
  "user_id" int NOT NULL,
  "username" varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  "full_name" varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  "avatar" varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  "parent_id" int DEFAULT NULL,
  "piece_id" int NOT NULL,
  PRIMARY KEY ("comment_id")
);

CREATE TABLE "comment_pro_comments" (
  "pro_comment_id" int NOT NULL AUTO_INCREMENT,
  "created_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "body" text NOT NULL,
  "num_replies" int NOT NULL DEFAULT '0',
  "num_likes" int DEFAULT '0',
  "user_id" int NOT NULL,
  "username" varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  "full_name" varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  "avatar" varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  "parent_id" int DEFAULT NULL,
  "profile_id" int NOT NULL,
  PRIMARY KEY ("pro_comment_id")
);

CREATE TABLE "follow_relationships" (
  "relationship_id" int NOT NULL AUTO_INCREMENT,
  "leader_id" int NOT NULL,
  "follower_id" int NOT NULL,
  "created_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("relationship_id")
);

CREATE TABLE "post_forums" (
  "forum_id" int NOT NULL AUTO_INCREMENT,
  "name" varchar(255) NOT NULL,
  "num_posts" int NOT NULL DEFAULT '0',
  "num_views" int NOT NULL DEFAULT '0',
  "updated_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "display_order" tinyint NOT NULL DEFAULT '0',
  "about" varchar(255) NOT NULL DEFAULT 'about',
  PRIMARY KEY ("forum_id")
);

CREATE TABLE "post_posts" (
  "post_id" int NOT NULL AUTO_INCREMENT,
  "user_id" int NOT NULL,
  "created_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "body" text NOT NULL,
  "num_replies" int NOT NULL DEFAULT '0',
  "num_likes" int DEFAULT '0',
  "post_type" tinyint NOT NULL DEFAULT '0',
  "username" varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  "full_name" varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  "avatar" varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  "parent_id" int DEFAULT NULL,
  "forum_id" int NOT NULL DEFAULT '1',
  PRIMARY KEY ("post_id")
);

CREATE TABLE "post_replies" (
  "id" int NOT NULL AUTO_INCREMENT,
  "parent_id" int NOT NULL,
  "user_id" int NOT NULL,
  "created_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "body" text NOT NULL,
  "num_replies" int NOT NULL DEFAULT '0',
  "num_likes" int NOT NULL DEFAULT '0',
  PRIMARY KEY ("id")
);

CREATE TABLE "profile_profiles" (
  "id" int NOT NULL AUTO_INCREMENT,
  "user_id" int NOT NULL,
  "username" varchar(255) NOT NULL,
  "full_name" varchar(255) NOT NULL,
  "bio" text,
  "location" varchar(255) DEFAULT NULL,
  "photo_uuid" varchar(255) DEFAULT NULL,
  "num_following" int NOT NULL DEFAULT '0',
  "num_followers" int NOT NULL DEFAULT '0',
  "num_collections" int NOT NULL DEFAULT '0',
  "num_pieces" int NOT NULL DEFAULT '0',
  "num_likes" int NOT NULL DEFAULT '0',
  "num_views" int NOT NULL DEFAULT '0',
  "num_comments" int NOT NULL DEFAULT '0',
  PRIMARY KEY ("id")
);



/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;