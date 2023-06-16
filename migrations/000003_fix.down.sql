ALTER TABLE `users` MODIFY `tlgrm_userid` INT(11) NOT NULL;
ALTER TABLE `users` MODIFY `rdmn_userid` INT(11) NOT NULL;

ALTER TABLE `issues_banch` MODIFY `tg_chat_id` INT(11) NOT NULL;
ALTER TABLE `issues_banch` MODIFY `tg_message_id` INT(11) NOT NULL;
ALTER TABLE `issues_banch` MODIFY `rdmn_issue_id` INT(11) NOT NULL;
ALTER TABLE `issues_banch` DROP CONSTRAINT `elt`;

ALTER TABLE `presale_issues` MODIFY `tlgrm_userid` INT(11) NOT NULL;
ALTER TABLE `presale_issues` MODIFY `rdmn_issue_id` INT(11) NOT NULL;
