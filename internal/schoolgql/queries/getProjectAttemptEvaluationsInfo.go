package queries

// get goalId -> studentAnswerId from this

const (
	GetProjectAttemptEvaluationsInfoByStudent TOperationName = `getProjectAttemptEvaluationsInfoByStudent`

	getProjectAttemptEvaluationsInfoByStudentQuery = `query getProjectAttemptEvaluationsInfoByStudent($goalId: ID!, $studentId: UUID!) {
  school21 {
    getProjectAttemptEvaluationsInfo(goalId: $goalId, studentId: $studentId) {
      ...ProjectAttemptEvaluations
    }
  }
}

fragment OnlineReviewInfo on OnlineReview {
  isOnline
  videos {
    onlineVideoId
    link
    status
    statusDetails
    updateDateTime
    fileSize
  }
}

fragment EvaluationFeedback on ReviewFeedback {
  id
  comment
  filledChecklist {
    id
  }
  reviewFeedbackCategoryValues {
    feedbackCategory
    feedbackValue
    id
  }
}

fragment Checklist on FilledChecklist {
  id
  checklistId
  endTimeCheck
  startTimeCheck
  reviewer {
    avatarUrl
    login
    businessAdminRoles {
      id
      school {
        id
        organizationType
      }
    }
  }
  reviewFeedback {
    ...EvaluationFeedback
  }
  comment
  receivedPoint
  receivedPercentage
  quickAction
  checkType
  onlineReview {
    ...OnlineReviewInfo
  }
}

fragment AttemptTeamMember on User {
  id
  avatarUrl
  login
  userExperience {
    level {
      id
      range {
        levelCode
      }
    }
    cookiesCount
    codeReviewPoints
  }
}

fragment P2PEvaluation on P2PEvaluationInfo {
  status
  checklist {
    ...Checklist
  }
}

fragment AttemptTeamWithMembers on TeamWithMembers {
  team {
    id
    name
  }
  members {
    role
    user {
      ...AttemptTeamMember
    }
  }
}

fragment AtemptResult on StudentGoalAttempt {
  finalPointProject
  finalPercentageProject
  resultModuleCompletion
  resultDate
}

fragment ProjectAttemptEvaluations on ProjectAttemptEvaluationsInfo {
  studentAnswerId
  attemptResult {
    ...AtemptResult
  }
  team {
    ...AttemptTeamWithMembers
  }
  p2p {
    ...P2PEvaluation
  }
  auto {
    status
    receivedPercentage
    endTimeCheck
    resultInfo
  }
  codeReview {
    averageMark
    studentCodeReviews {
      user {
        avatarUrl
        login
      }
      finalMark
      markTime
      reviewerCommentsCount
    }
  }
}`
)

type VarsGetProjectAttemptEvaluationsInfoByStudent struct {
	GoalID    string `json:"goalId"`
	StudentID string `json:"studentId"`
}

type ResponseGetProjectAttemptEvaluationsInfoByStudent struct {
	BaseResponse
	Data struct {
		School21 struct {
			GetProjectAttemptEvaluationsInfo []struct {
				AttemptResult any `json:"attemptResult"`
				Auto          struct {
					EndTimeCheck       *string `json:"endTimeCheck"`
					ReceivedPercentage int     `json:"receivedPercentage"`
					ResultInfo         any     `json:"resultInfo"`
					Status             string  `json:"status"`
				} `json:"auto"`

				// CodeReview struct {
				// 	AverageMark        any   `json:"averageMark"`
				// 	StudentCodeReviews []interface{} `json:"studentCodeReviews"`
				// } `json:"codeReview"`

				P2P []struct {
					// Checklist any `json:"checklist"`
					Status string `json:"status"`
				} `json:"p2p"`

				StudentAnswerId string `json:"studentAnswerId"`
				// Team            any    `json:"team"`
			} `json:"getProjectAttemptEvaluationsInfo"`
		} `json:"school21"`
	} `json:"data"`
}
