-- store.`member` definition

CREATE TABLE IF NOT EXISTS `member` (
                          `id` int(11) NOT NULL AUTO_INCREMENT,
                          `channel_id` varchar(100) NOT NULL,
                          `username` varchar(100) NOT NULL,
                          `credential` varchar(100) NOT NULL,
                          `salt` varchar(100) NOT NULL,
                          `created_date` timestamp NULL DEFAULT current_timestamp(),
                          PRIMARY KEY (`id`)
);