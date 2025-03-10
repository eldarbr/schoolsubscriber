package queries

// answerId, startTime -> slot

const (
	CalendarAddBookingToEventSlot TOperationName = `calendarAddBookingToEventSlot`

	calendarAddBookingToEventSlotQuery = `mutation calendarAddBookingToEventSlot($answerId: ID!, $startTime: DateTime!, $wasStaffSlotChosen: Boolean!, $isOnline: Boolean) {
  student {
    addBookingP2PToEventSlot(
      answerId: $answerId
      startTime: $startTime
      wasStaffSlotChosen: $wasStaffSlotChosen
      isOnline: $isOnline
    ) {
      id
    }
  }
}
`
)

type VarsCalendarAddBookingToEventSlot struct {
	AnswerID           string `json:"answerId"`
	StartTime          string `json:"startTime"`
	IsOnline           bool   `json:"isOnline"`
	WasStaffSlotChosen bool   `json:"wasStaffSlotChosen"`
}

type ResponseCalendarAddBookingToEventSlot struct {
	BaseResponse
	Data struct {
		Student struct {
			AddBookingP2PToEventSlot struct {
				ID string `json:"id"`
			} `json:"addBookingP2PToEventSlot"`
		} `json:"student"`
	} `json:"data"`
}
