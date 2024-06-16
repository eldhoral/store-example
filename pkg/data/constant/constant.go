package constant

import "os"

const (
	Alphanumeric                  = ".*[a-zA-Z].*"
	DefaultDatetimeLayout         = "2006-01-02 15:04:05"
	DefaultDatetimeTimezoneLayout = "2006-01-02T15:04:05Z"

	DefaultDateLayout = "2006-01-02"

	DefaultTimeLayout = "15:04:05"

	YYYYmmddHHmmss = "20060102150405"

	YYYYLayout = "2006"
	MMLayout   = "01"

	TaskListDateLayout = "02 Jan 2006"
	HHMMLayout         = "15:04"

	MaxFileSize = 10485760

	EnterpreneurJob = "PENGUSAHA_WIRASWASTA_PERSEORANGAN_KORPORASI"

	Initialized          = 1
	PersonalInformation  = 2
	AccountInformation   = 4
	JobInformation       = 8
	EmergencyInformation = 16
	OTPValidated         = 32

	StateNotValidated         = 31
	StatePersonalInformation  = 3
	StateAccountInformation   = 7
	StateJobInformation       = 15
	StateEmergencyInformation = 31

	UrlLoginSuperApps             = "/auth/login"
	UrlDropdownEducation          = "/mobile/kyc-options/kyc_pendidikan"
	UrlDropdownMaritalStatus      = "/mobile/kyc-options/kyc_status_perkawinan"
	UrlDropdownResidenceOwnership = "/mobile/kyc-options/kyc_status_rumah"
	UrlDropdownReligion           = "/mobile/kyc-options/kyc_agama"
	UrlDropdownProvince           = "/mobile/areas/provinces"
	UrlDropdownDistricts          = "/mobile/areas/cities"
	UrlDropdownUrban              = "/mobile/areas/districts"
	UrlDropdownJobs               = "/mobile/kyc-options/kyc_pekerjaan"
	UrlDropdownBusinessFields     = "/mobile/kyc-options/kyc_bidang_usaha_2"
	UrlSearchArea                 = "/kyc-options/areas?keyword="
	UrlEvaluateToken              = "/auth/validate"
	UrlInquiryByNikAndChannel     = "/customer/inquiry-by-nik-and-channel"
	UrlMemberInfo                 = "/info"
	UrlMemberUpgradedInfo         = "/upgrade/info"

	//OTP
	UrlRegistrationOTP = "/paylater/registration/otp"
	UrlValidateOTP     = "/paylater/validate/otp"
	UrlDroppedInEmail  = "/paylater/droppedin/email"

	CookiePath = "Path=/"

	XRequestNameLogin    = "Login"
	XRequestNameNotif    = "Login"
	GetMemberInfo        = "GetMemberInfo"
	GetMemberUpgradeInfo = "GetMemberUpgradeInfo"

	CacheKeySearchProvince  = "paylater:dropdown:search.province:"
	CacheKeySearchCity      = "paylater:dropdown:search.city:"
	CacheKeySearchDistricts = "paylater:dropdown:search.districts:"
	CacheKeySearchUrban     = "paylater:dropdown:search.urban:"

	CacheKeySearchJobs               = "paylater:dropdown:search.jobs:"
	CacheKeySearchBusinessFields     = "paylater:dropdown:search.business_fields:"
	CacheKeySearchResidenceOwnership = "paylater:dropdown:search.residence_ownership:"
	CacheKeySearchEducation          = "paylater:dropdown:search.education:"
	CacheKeySearchMaritalStatus      = "paylater:dropdown:search.marital_status:"

	CacheKeyListJobs               = "paylater:dropdown:list.jobs"
	CacheKeyListBusinessFields     = "paylater:dropdown:list.business_fields"
	CacheKeyListResidenceOwnership = "paylater:dropdown:list.residence_ownership"
	CacheKeyListEducation          = "paylater:dropdown:list.education"
	CacheKeyListMaritalStatus      = "paylater:dropdown:list.marital_status"

	CacheKeyListProvince  = "paylater:dropdown:list.province:"
	CacheKeyListCity      = "paylater:dropdown:list.city:"
	CacheKeyListDistricts = "paylater:dropdown:list.districts:"
	CacheKeyListUrban     = "paylater:dropdown:list.urban:"

	CacheOtpRegistrationRequest = "paylater:otp:registration.request:"
	CacheOtpRegistrationInvalid = "paylater:otp:registration.invalid:"

	NotifType       = "PUSHNOTIF"
	NotifMessageApp = "Pengajuan Nobu Paylater kamu sedang diverifikasi. Kami akan mengirimkan verifikasi persetujuan dalam 1 x 24 jam."
	NotifSubjectApp = "Pengajuan Nobu Paylater kamu sedang diverifikasi"

	DropdownOption                   = "dropdown_options"
	DropdownOptionBusinessFields     = "business_fields"
	DropdownOptionResidenceOwnership = "residence_ownership"
	DropdownOptionJobs               = "jobs"
	DropdownOptionEducation          = "education"
	DropdownOptionMaritalStatus      = "marital_status"
)

var (
	EmptyArray         = []int{}
	EmptyStringPointer = func() *string { i := ""; return &i }()
	AllowedFile        = []string{".pdf", ".jpg", ".png", ".xlsx", ".xls", ".jpeg", ".docx", ".doc", ".csv", ".txt", ".ppt", ".pptx"}
	InformalJob        = []string{"IBU_RUMAH_TANGGA", "PEKERJA_SENL_OLAH_RAGA_KEAGAMAAN_SEJENISNYA", "SELEBRITI"}
	JobType            = []string{"EMPLOYEE", "ENTREPRENEUR", "INFORMAL_JOB"}
)

// Based on IsMySQL, we can change SQL, and timezone to match with DB
func IsMySQL() bool {
	return os.Getenv("DB_TZ") == "UTC"
}
