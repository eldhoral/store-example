-- store.`transaction` definition

CREATE TABLE IF NOT EXISTS `transaction` (
                               `id` int(11) NOT NULL AUTO_INCREMENT,
                               `trx_code` varchar(100) DEFAULT NULL,
                               `channel_id` varchar(100) NOT NULL,
                               `channel_ref_no` varchar(100) NOT NULL,
                               `channel_time` varchar(100) NOT NULL,
                               `channel_date` varchar(100) NOT NULL,
                               `amount` float NOT NULL,
                               `amount_fee` float NOT NULL,
                               `member_id` int(11) NOT NULL,
                               `status` varchar(100) NOT NULL,
                               `product_id` int(11) NOT NULL,
                               `quantity` int(11) NOT NULL,
                               `created_date` timestamp NULL DEFAULT current_timestamp(),
                               `updated_date` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
                               PRIMARY KEY (`id`)
):