ALTER TABLE `users` MODIFY `tlgrm_userid` BIGINT NOT NULL;
ALTER TABLE `users` MODIFY `rdmn_userid` BIGINT NOT NULL;

ALTER TABLE `issues_banch` MODIFY `tg_chat_id` BIGINT NOT NULL;
ALTER TABLE `issues_banch` MODIFY `tg_message_id` BIGINT NOT NULL;
ALTER TABLE `issues_banch` MODIFY `rdmn_issue_id` BIGINT NOT NULL;
ALTER TABLE `issues_banch` ADD CONSTRAINT `elt` UNIQUE (`tg_chat_id`,`tg_message_id`, `rdmn_issue_id`);

ALTER TABLE `presale_issues` MODIFY `tlgrm_userid` BIGINT NOT NULL;
ALTER TABLE `presale_issues` MODIFY `rdmn_issue_id` BIGINT NOT NULL;

