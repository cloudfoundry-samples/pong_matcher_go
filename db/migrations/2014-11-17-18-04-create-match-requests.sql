-- +migrate Up
CREATE TABLE `match_requests` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `uuid` varchar(255) DEFAULT NULL,
    `requester_id` varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`)
);

-- +migrate Down
DROP TABLE match_requests;
