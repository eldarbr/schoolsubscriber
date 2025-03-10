package queries

const (
	CalendarGetNameLessStudentTimeslotsForReview TOperationName = `calendarGetNameLessStudentTimeslotsForReview`

	calendarGetNameLessStudentTimeslotsForReviewQuery TQuery = `query calendarGetNameLessStudentTimeslotsForReview($from: DateTime!, $taskId: ID!, $to: DateTime!) {
    student {
      getNameLessStudentTimeslotsForReview(from: $from, taskId: $taskId, to: $to) {
        checkDuration
        projectReviewsInfo {
          ...ProjectReviewsInfo
        }
        timeSlots {
          ...CalendarNameLessTimeslot
        }
      }
    }
  }
  
  fragment ProjectReviewsInfo on ProjectReviewsInfo {
    reviewByStudentCount
    relevantReviewByStudentsCount
    reviewByInspectionStaffCount
    relevantReviewByInspectionStaffCount
  }
  
  fragment CalendarNameLessTimeslot on CalendarNamelessTimeSlot {
    start
    end
    validStartTimes
    staffSlot
  }
`
)

type VarsCalendarGetNameLessStudentTimeslotsForReview struct {
	From   string `json:"from"`
	To     string `json:"to"`
	TaskID string `json:"taskId"`
}

type ResponseCalendarGetNameLessStudentTimeslotsForReview struct {
	BaseResponse
	Data struct {
		Student struct {
			GetNameLessStudentTimeslotsForReview struct {
				CheckDuration      int `json:"checkDuration"`
				ProjectReviewsInfo struct {
					ReviewByStudentCount                 int `json:"reviewByStudentCount"`
					RelevantReviewByStudentsCount        int `json:"relevantReviewByStudentsCount"`
					ReviewByInspectionStaffCount         int `json:"reviewByInspectionStaffCount"`
					RelevantReviewByInspectionStaffCount int `json:"relevantReviewByInspectionStaffCount"`
				} `json:"projectReviewsInfo"`
				TimeSlots []struct {
					Start           string   `json:"start"`
					End             string   `json:"end"`
					ValidStartTimes []string `json:"validStartTimes"`
					StaffSlot       bool     `json:"staffSlot"`
				} `json:"timeSlots"`
			} `json:"getNameLessStudentTimeslotsForReview"`
		} `json:"student"`
	} `json:"data"`
}
