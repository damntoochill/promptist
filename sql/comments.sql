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
  PRIMARY KEY ("comment_id")
);