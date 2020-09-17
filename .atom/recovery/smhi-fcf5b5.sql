USE forelcastro;

CREATE TABLE `SMHI` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `reference_time` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `stripe_id` (`stripe_id`)
);
