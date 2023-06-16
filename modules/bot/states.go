package tgbot

import tg "github.com/nixys/nxs-go-telegram"

var (
	stateHello     = tg.SessState("hello")
	stateBye       = tg.SessState("bye")
	stateHelp      = tg.SessState("help")
	stateForbidden = tg.SessState("forbidden")

	stateInitLang                = tg.SessState("initLang")
	stateInitMode                = tg.SessState("initMode")
	stateInitAccount             = tg.SessState("initAccount")
	stateInitRdmnAPIKeyIncorrect = tg.SessState("initRdmnAPIKeyIncorrectState")
	stateInitEnd                 = tg.SessState("initEnd")

	stateFeedback = tg.SessState("feedback")

	stateIssueCreate            = tg.SessState("issueCreate")
	stateIssueCreateProject     = tg.SessState("issueCreateProject")
	stateIssueCreatePriority    = tg.SessState("issueCreatePriority")
	stateIssueCreateSubject     = tg.SessState("issueCreateSubject")
	stateIssueCreateDescription = tg.SessState("issueCreateDescription")
	stateIssueCreateConfirm     = tg.SessState("issueCreateConfirm")
	stateIssueCreateEnd         = tg.SessState("issueCreateEnd")

	stateSettings                    = tg.SessState("settings")
	stateSettingsRdmn                = tg.SessState("settRdmn")
	stateSettingsRdmnAcc             = tg.SessState("settRdmnAcc")
	stateSettingsRdmnAPIKeySet       = tg.SessState("settRdmnAPIKeySet")
	stateSettingsRdmnAPIKeyIncorrect = tg.SessState("settRdmnApiKeyIncorrect")
	stateSettingsLangSelect          = tg.SessState("settLangSelect")
)
