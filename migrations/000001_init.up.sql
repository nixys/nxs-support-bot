CREATE TABLE `users` (
  `tlgrm_userid` INT(11) NOT NULL,
  `rdmn_userid` INT(11) NOT NULL,
  `lang` VARCHAR(10) NOT NULL DEFAULT 'en',
  UNIQUE KEY `id` (`tlgrm_userid`, `rdmn_userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `issues_banch` (
  `tg_chat_id` INT(11) NOT NULL,
  `tg_message_id` INT(11) NOT NULL,
  `rdmn_issue_id` INT(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
