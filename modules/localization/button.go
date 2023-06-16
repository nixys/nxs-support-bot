package localization

type Button string

const (
	ButtonOk               Button = "buttonOk"
	ButtonBack             Button = "buttonBack"
	ButtonQuit             Button = "buttonQuit"
	ButtonCancel           Button = "buttonCancel"
	ButtonPrevPage         Button = "buttonPrevPage"
	ButtonNextPage         Button = "buttonNextPage"
	ButtonRdmn             Button = "buttonRedmine"
	ButtonAccount          Button = "buttonAccount"
	ButtonLink             Button = "buttonLink"
	ButtonFavoriteProjects Button = "buttonFavoriteProjects"
	ButtonLang             Button = "buttonLanguage"
	ButtonEN               Button = "buttonEN"
	ButtonRU               Button = "buttonRU"
	ButtonProject          Button = "buttonProject"
	ButtonPriority         Button = "buttonPriority"
	ButtonCreateIssue      Button = "buttonCreateIssue"
	ButtonLeaveEmpty       Button = "buttonLeaveEmpty"
	ButtonAuthorize        Button = "buttonAuthorize"
	ButtonContactUs        Button = "buttonContactUs"
)

var buttons = []Button{
	ButtonOk,
	ButtonBack,
	ButtonQuit,
	ButtonCancel,
	ButtonPrevPage,
	ButtonNextPage,
	ButtonRdmn,
	ButtonAccount,
	ButtonLink,
	ButtonFavoriteProjects,
	ButtonLang,
	ButtonEN,
	ButtonRU,
	ButtonProject,
	ButtonPriority,
	ButtonCreateIssue,
	ButtonLeaveEmpty,
	ButtonAuthorize,
	ButtonContactUs,
}

func (b Button) String() string {
	return string(b)
}
