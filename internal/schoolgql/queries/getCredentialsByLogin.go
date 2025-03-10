package queries

const (
	GetCredentialsByLogin      TOperationName = `getCredentialsByLogin`
	getCredentialsByLoginQuery TQuery         = `query getCredentialsByLogin($login: String!) {
  school21 {
    getStudentByLogin(login: $login) {
      studentId
      userId
      schoolId
      isActive
      isGraduate
      __typename
    }
    __typename
  }
}
`
)

type VarsGetCredentialsByLogin struct {
	Login string `json:"login"`
}

type ResponseGetCredentialsByLogin struct {
	BaseResponse
	Data struct {
		School21 struct {
			GetStudentByLogin struct {
				IsActive   bool   `json:"isActive"`
				IsGraduate bool   `json:"isGraduate"`
				SchoolID   string `json:"schoolId"`
				StudentID  string `json:"studentId"`
				UserID     string `json:"userId"`
			} `json:"getStudentByLogin"`
		} `json:"school21"`
	} `json:"data"`
}
