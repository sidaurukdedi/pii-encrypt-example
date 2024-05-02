CREATE TABLE `user_encrypt` (
  `uuid` char(36),
  `__encrypted__data_nama_crypt` varbinary(255) NOT NULL,
  `__encrypted__data_nama_hash` varbinary(32) NOT NULL,
  `__encrypted__data_email_crypt` varbinary(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`uuid`)
);