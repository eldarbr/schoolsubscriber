package queries

const (
	PublicProfileGetPersonalInfo TOperationName = `publicProfileGetPersonalInfo`

	publicProfileGetPersonalInfoQuery TQuery = `query publicProfileGetPersonalInfo($userId: UUID!, $studentId: UUID!, $login: String!, $schoolId: UUID!) {
  school21 {
    getAvatarByUserId(userId: $userId)
    getStageGroupS21PublicProfile(studentId: $studentId) {
      waveId
      waveName
      eduForm
      __typename
    }
    getExperiencePublicProfile(userId: $userId) {
      value
      level {
        levelCode
        range {
          leftBorder
          rightBorder
          __typename
        }
        __typename
      }
      cookiesCount
      coinsCount
      codeReviewPoints
      isReviewPointsConsistent
      __typename
    }
    getEmailbyUserId(userId: $userId)
    getClassRoomByLogin(login: $login) {
      id
      number
      floor
      __typename
    }
    __typename
  }
  student {
    getWorkstationByLogin(login: $login) {
      workstationId
      hostName
      row
      number
      __typename
    }
    getFeedbackStatisticsAverageScore(studentId: $studentId) {
      countFeedback
      feedbackAverageScore {
        categoryCode
        categoryName
        value
        __typename
      }
      __typename
    }
    __typename
  }
  user {
    getSchool(schoolId: $schoolId) {
      id
      fullName
      shortName
      address
      __typename
    }
    __typename
  }
}
`
)

type VarsPublicProfileGetPersonalInfo struct {
	Login     string `json:"login"`
	SchoolID  string `json:"schoolId"`
	StudentID string `json:"studentId"`
	UserID    string `json:"userId"`
}

type ResponsePublicProfileGetPersonalInfo struct {
	BaseResponse
	Data struct {
		School21 struct {
			GetAvatarByUserID          string `json:"getAvatarByUserId"`
			GetClassRoomByLogin        any    `json:"getClassRoomByLogin"`
			GetEmailbyUserID           string `json:"getEmailbyUserId"`
			GetExperiencePublicProfile struct {
				CodeReviewPoints         int  `json:"codeReviewPoints"`
				CoinsCount               int  `json:"coinsCount"`
				CookiesCount             int  `json:"cookiesCount"`
				IsReviewPointsConsistent bool `json:"isReviewPointsConsistent"`
				Level                    struct {
					LevelCode int `json:"levelCode"`
					Range     struct {
						LeftBorder  int `json:"leftBorder"`
						RightBorder int `json:"rightBorder"`
					} `json:"range"`
				} `json:"level"`
				Value int `json:"value"`
			} `json:"getExperiencePublicProfile"`
			GetStageGroupS21PublicProfile struct {
				EduForm  string `json:"eduForm"`
				WaveID   int    `json:"waveId"`
				WaveName string `json:"waveName"`
			} `json:"getStageGroupS21PublicProfile"`
		} `json:"school21"`
		Student struct {
			GetFeedbackStatisticsAverageScore struct {
				CountFeedback        int `json:"countFeedback"`
				FeedbackAverageScore []struct {
					CategoryCode string `json:"categoryCode"`
					CategoryName string `json:"categoryName"`
					Value        string `json:"value"`
				} `json:"feedbackAverageScore"`
			} `json:"getFeedbackStatisticsAverageScore"`
			GetWorkstationByLogin any `json:"getWorkstationByLogin"`
		} `json:"student"`
		User struct {
			GetSchool struct {
				Address   string `json:"address"`
				FullName  string `json:"fullName"`
				ID        string `json:"id"`
				ShortName string `json:"shortName"`
			} `json:"getSchool"`
		} `json:"user"`
	} `json:"data"`
}
