-- +migrate Up
CREATE TABLE `results` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `match_id` varchar(255) DEFAULT NULL,
    `winner` varchar(255) DEFAULT NULL,
    `loser` varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`)
);

-- +migrate Down
DROP TABLE `results`;
