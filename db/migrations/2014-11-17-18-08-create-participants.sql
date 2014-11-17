-- +migrate Up
CREATE TABLE `participants` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `match_id` varchar(255) DEFAULT NULL,
    `match_request_uuid` varchar(255) DEFAULT NULL,
    `player_id` varchar(255) DEFAULT NULL,
    `opponent_id` varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `match_request_uuid` (`match_request_uuid`)
);

-- +migrate Down
DROP TABLE participants;
