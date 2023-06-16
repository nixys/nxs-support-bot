package localization

type Message string

const (
	MsgHello                  Message = "msgHello"
	MsgBye                    Message = "msgBye"
	MsgHelp                   Message = "msgHelp"
	MsgAPIKeyIncorrect        Message = "msgAPIKeyIncorrect"
	MsgForbidden              Message = "msgForbidden"
	MsgOrphanedReply          Message = "msgOrphanedReply"
	MsgIssueCreate            Message = "msgIssueCreate"
	MsgIssueCreateProject     Message = "msgIssueCreateProject"
	MsgIssueCreatePriority    Message = "msgIssueCreatePriority"
	MsgIssueCreateSubject     Message = "msgIssueCreateSubject"
	MsgIssueCreateDescription Message = "msgIssueCreateDescription"
	MsgIssueCreateConfirm     Message = "msgIssueCreateConfirm"
	MsgIssueCreateEnd         Message = "msgIssueCreateEnd"
	MsgIssueCreated           Message = "msgIssueCreated"
	MsgIssueUpdated           Message = "msgIssueUpdated"
	MsgInitLang               Message = "msgInitLang"
	MsgInitMode               Message = "msgInitMode"
	MsgInitRdmnApiKeySet      Message = "msgInitRdmnApiKeySet"
	MsgInitEnd                Message = "msgInitEnd"
	MsgSettings               Message = "msgSettings"
	MsgSettingsLang           Message = "msgSettingsLang"
	MsgSettingsRdmn           Message = "msgSettingsRdmn"
	MsgSettingsRdmnAcc        Message = "msgSettingsRdmnAcc"
	MsgSettingsRdmnApiKeySet  Message = "msgSettingsRdmnApiKeySet"
	MsgFeedbackGreetings      Message = "msgFeedbackGreetings"
	MsgFeedbackIssueUpdated   Message = "msgFeedbackIssueUpdated"
)

func (m Message) String() string {
	return string(m)
}
