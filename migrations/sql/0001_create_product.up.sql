-- store.product definition

CREATE TABLE IF NOT EXISTS `product` (
                           `id` int(11) NOT NULL AUTO_INCREMENT,
                           `name` varchar(100) NOT NULL,
                           `category` varchar(100) NOT NULL,
                           `price` float NOT NULL,
                           `stock` int(11) NOT NULL,
                           PRIMARY KEY (`id`)
);