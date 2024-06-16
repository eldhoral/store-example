-- store.cart definition

CREATE TABLE IF NOT EXISTS  `cart` (
                        `id` int(11) NOT NULL AUTO_INCREMENT,
                        `member_id` int(11) NOT NULL,
                        `product_id` int(11) NOT NULL,
                        `quantity` int(11) NOT NULL,
                        `created_date` timestamp NOT NULL DEFAULT current_timestamp(),
                        `is_active` tinyint(1) NOT NULL,
                        PRIMARY KEY (`id`)
);