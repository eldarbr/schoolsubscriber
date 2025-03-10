package queries

const (
	SendInvitation TOperationName = `sendInvitation`

	sendInvitationQuery TQuery = `mutation sendInvitation($teamId: UUID!, $userId: ID!) {
  student {
    sendInvitation(teamId: $teamId, userId: $userId) {
      ...StudentInvitationInfo
      __typename
    }
    __typename
  }
}

fragment StudentInvitationInfo on StudentInvitationInfo {
  student {
    ...AvailableStudentForTeam
    __typename
  }
  invitationStatus
  schoolShortName
  __typename
}

fragment AvailableStudentForTeam on Student {
  id
  user {
    id
    login
    avatarUrl
    userExperience {
    ...CurrentUserExperience
    __typename
    }
    __typename
  }
  __typename
}

fragment CurrentUserExperience on UserExperience {
  id
  cookiesCount
  codeReviewPoints
  coinsCount
  level {
    id
    range {
      id
      levelCode
      __typename
    }
    __typename
  }
  __typename
}
`
)

type VarsSendInvitation struct {
	TeamID string `json:"teamId"`
	UserID string `json:"userId"`
}
