CREATE TABLE `presale_issues` (
  `tlgrm_userid` INT(11) NOT NULL,
  `rdmn_issue_id` INT(11) NOT NULL,
  UNIQUE KEY `tlgrm_userid` (`tlgrm_userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
