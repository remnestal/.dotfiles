use oauth;

drop table StripePaymentMethods;

create table StripePaymentMethods (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `stripe_payment_method_id` varchar(64) NOT NULL,
  `stripe_card_fingerprint` varchar(64) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`public_product_id`),
  UNIQUE KEY `stripe_product_id` (`stripe_product_id`)
) DEFAULT CHARSET=utf8;
