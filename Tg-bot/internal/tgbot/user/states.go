package user

const (
	ERROR = iota
	NONE
	WaitingUrl
	WaitingHashtag
	WaitingUrlForRemove
)

type StateAction struct {
	NextState     int
	Message       string
	FieldtoChange string
}

var (
	TrackStates = map[int]StateAction{
		NONE:           {NextState: WaitingUrl, Message: "Пришлите мне ссылку:", FieldtoChange: ""},
		WaitingUrl:     {NextState: WaitingHashtag, Message: "Назовите категорию ссылки:", FieldtoChange: "Link"},
		WaitingHashtag: {NextState: NONE, Message: "Ссылка успешно сохранена!", FieldtoChange: "Category"}}

	UnTrackStates = map[int]StateAction{
		NONE:                {NextState: WaitingUrlForRemove, Message: "Пришлите мне ссылку для удаления:", FieldtoChange: ""},
		WaitingUrlForRemove: {NextState: NONE, Message: "Ваша ссылка удалена!", FieldtoChange: "Link"}}
)
