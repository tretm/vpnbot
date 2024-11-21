CREATE TABLE `ip_stories` 
(
  `id` int NOT NULL,
  `ip` varchar(200) DEFAULT NULL,
  `user_agent` varchar(200) DEFAULT NULL,
  `country` varchar(200) DEFAULT NULL,
  `city` varchar(200) DEFAULT NULL,
  
  `provider` varchar(200) DEFAULT NULL,
  `company` varchar(200) DEFAULT NULL,
  `link` varchar(20) DEFAULT NULL,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `ip_stories`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `ip_stories`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;