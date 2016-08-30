package server

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/intervention-engine/fhir/auth"
	"github.com/mitre/heart"
	"golang.org/x/oauth2"
)

// RegisterController registers the CRUD routes (and middleware) for a FHIR resource
func RegisterController(name string, e *gin.Engine, m []gin.HandlerFunc, dal DataAccessLayer, config Config) {
	rc := NewResourceController(name, dal)
	rcBase := e.Group("/" + name)

	if len(m) > 0 {
		rcBase.Use(m...)
	}

	switch config.Auth.Method {
	case auth.AuthTypeNone:
		// do nothing
	case auth.AuthTypeOIDC:
		rcBase.Use(auth.HEARTScopesHandler(name))
	case auth.AuthTypeHEART:
		rcBase.Use(auth.HEARTScopesHandler(name))
	}

	rcBase.GET("", rc.IndexHandler)
	rcBase.POST("", rc.CreateHandler)
	rcBase.PUT("", rc.ConditionalUpdateHandler)
	rcBase.DELETE("", rc.ConditionalDeleteHandler)

	rcItem := rcBase.Group("/:id")
	rcItem.GET("", rc.ShowHandler)
	rcItem.PUT("", rc.UpdateHandler)
	rcItem.DELETE("", rc.DeleteHandler)
}

// RegisterRoutes registers the routes for each of the FHIR resources
func RegisterRoutes(e *gin.Engine, config map[string][]gin.HandlerFunc, dal DataAccessLayer, serverConfig Config) {

	switch serverConfig.Auth.Method {
	case auth.AuthTypeNone:
		// do nothing
	case auth.AuthTypeOIDC:
		// Set up sessions so we can keep track of the logged in user
		store := sessions.NewCookieStore([]byte(serverConfig.Auth.SessionSecret))
		e.Use(sessions.Sessions("mysession", store))
		// The OIDCAuthenticationHandler is set up before the IndexHandler in the handler function
		// chain. It will check to see if the user is logged in based on their session. If they are not
		// the user will be redirected to the authentication endpoint at the OP.
		oauthConfig := oauth2.Config{ClientID: serverConfig.Auth.ClientID,
			ClientSecret: serverConfig.Auth.ClientSecret,
			Endpoint: oauth2.Endpoint{AuthURL: serverConfig.Auth.AuthorizationURL,
				TokenURL: serverConfig.Auth.TokenURL},
		}
		oidcHandler := auth.OIDCAuthenticationHandler(oauthConfig)
		oauthHandler := auth.OAuthIntrospectionHandler(serverConfig.Auth.ClientID,
			serverConfig.Auth.ClientSecret, serverConfig.Auth.IntrospectionURL)
		e.Use(func(c *gin.Context) {
			if c.Request.Header.Get("Authorization") != "" {
				oauthHandler(c)
			} else {
				oidcHandler(c)
			}
		})
		// This handler is to take the redirect from the OP when the user logs in. It will
		// then fetch information about the user by hitting the user info endpoint and put
		// that in the session. Lastly, this handler is set up to redirect the user back
		// to the root.
		e.GET("/redirect", auth.RedirectHandler(oauthConfig, serverConfig.ServerURL,
			serverConfig.Auth.UserInfoURL))
		e.GET("/logout", heart.LogoutHandler)

	case auth.AuthTypeHEART:
		heart.SetUpRoutes(serverConfig.Auth.JWKPath, serverConfig.Auth.ClientID, serverConfig.Auth.OPURL,
			serverConfig.ServerURL, serverConfig.Auth.SessionSecret, e)

	}

	// Batch Support
	batch := NewBatchController(dal)
	batchHandlers := make([]gin.HandlerFunc, len(config["Batch"]))
	copy(batchHandlers, config["Batch"])
	batchHandlers = append(batchHandlers, batch.Post)
	e.POST("/", batchHandlers...)

	// Resources

	RegisterController("Appointment", e, config["Appointment"], dal, serverConfig)
	RegisterController("ReferralRequest", e, config["ReferralRequest"], dal, serverConfig)
	RegisterController("Account", e, config["Account"], dal, serverConfig)
	RegisterController("Provenance", e, config["Provenance"], dal, serverConfig)
	RegisterController("Questionnaire", e, config["Questionnaire"], dal, serverConfig)
	RegisterController("ExplanationOfBenefit", e, config["ExplanationOfBenefit"], dal, serverConfig)
	RegisterController("DocumentManifest", e, config["DocumentManifest"], dal, serverConfig)
	RegisterController("Specimen", e, config["Specimen"], dal, serverConfig)
	RegisterController("AllergyIntolerance", e, config["AllergyIntolerance"], dal, serverConfig)
	RegisterController("CarePlan", e, config["CarePlan"], dal, serverConfig)
	RegisterController("Goal", e, config["Goal"], dal, serverConfig)
	RegisterController("StructureDefinition", e, config["StructureDefinition"], dal, serverConfig)
	RegisterController("EnrollmentRequest", e, config["EnrollmentRequest"], dal, serverConfig)
	RegisterController("EpisodeOfCare", e, config["EpisodeOfCare"], dal, serverConfig)
	RegisterController("OperationOutcome", e, config["OperationOutcome"], dal, serverConfig)
	RegisterController("Medication", e, config["Medication"], dal, serverConfig)
	RegisterController("Procedure", e, config["Procedure"], dal, serverConfig)
	RegisterController("List", e, config["List"], dal, serverConfig)
	RegisterController("ConceptMap", e, config["ConceptMap"], dal, serverConfig)
	RegisterController("Subscription", e, config["Subscription"], dal, serverConfig)
	RegisterController("ValueSet", e, config["ValueSet"], dal, serverConfig)
	RegisterController("OperationDefinition", e, config["OperationDefinition"], dal, serverConfig)
	RegisterController("DocumentReference", e, config["DocumentReference"], dal, serverConfig)
	RegisterController("Order", e, config["Order"], dal, serverConfig)
	RegisterController("Immunization", e, config["Immunization"], dal, serverConfig)
	RegisterController("Device", e, config["Device"], dal, serverConfig)
	RegisterController("VisionPrescription", e, config["VisionPrescription"], dal, serverConfig)
	RegisterController("Media", e, config["Media"], dal, serverConfig)
	RegisterController("Conformance", e, config["Conformance"], dal, serverConfig)
	RegisterController("ProcedureRequest", e, config["ProcedureRequest"], dal, serverConfig)
	RegisterController("EligibilityResponse", e, config["EligibilityResponse"], dal, serverConfig)
	RegisterController("DeviceUseRequest", e, config["DeviceUseRequest"], dal, serverConfig)
	RegisterController("DeviceMetric", e, config["DeviceMetric"], dal, serverConfig)
	RegisterController("Flag", e, config["Flag"], dal, serverConfig)
	RegisterController("RelatedPerson", e, config["RelatedPerson"], dal, serverConfig)
	RegisterController("SupplyRequest", e, config["SupplyRequest"], dal, serverConfig)
	RegisterController("Practitioner", e, config["Practitioner"], dal, serverConfig)
	RegisterController("AppointmentResponse", e, config["AppointmentResponse"], dal, serverConfig)
	RegisterController("Observation", e, config["Observation"], dal, serverConfig)
	RegisterController("MedicationAdministration", e, config["MedicationAdministration"], dal, serverConfig)
	RegisterController("Slot", e, config["Slot"], dal, serverConfig)
	RegisterController("EnrollmentResponse", e, config["EnrollmentResponse"], dal, serverConfig)
	RegisterController("Binary", e, config["Binary"], dal, serverConfig)
	RegisterController("MedicationStatement", e, config["MedicationStatement"], dal, serverConfig)
	RegisterController("Person", e, config["Person"], dal, serverConfig)
	RegisterController("Contract", e, config["Contract"], dal, serverConfig)
	RegisterController("CommunicationRequest", e, config["CommunicationRequest"], dal, serverConfig)
	RegisterController("RiskAssessment", e, config["RiskAssessment"], dal, serverConfig)
	RegisterController("TestScript", e, config["TestScript"], dal, serverConfig)
	RegisterController("Basic", e, config["Basic"], dal, serverConfig)
	RegisterController("Group", e, config["Group"], dal, serverConfig)
	RegisterController("PaymentNotice", e, config["PaymentNotice"], dal, serverConfig)
	RegisterController("Organization", e, config["Organization"], dal, serverConfig)
	RegisterController("ImplementationGuide", e, config["ImplementationGuide"], dal, serverConfig)
	RegisterController("ClaimResponse", e, config["ClaimResponse"], dal, serverConfig)
	RegisterController("EligibilityRequest", e, config["EligibilityRequest"], dal, serverConfig)
	RegisterController("ProcessRequest", e, config["ProcessRequest"], dal, serverConfig)
	RegisterController("MedicationDispense", e, config["MedicationDispense"], dal, serverConfig)
	RegisterController("DiagnosticReport", e, config["DiagnosticReport"], dal, serverConfig)
	RegisterController("ImagingStudy", e, config["ImagingStudy"], dal, serverConfig)
	RegisterController("ImagingObjectSelection", e, config["ImagingObjectSelection"], dal, serverConfig)
	RegisterController("HealthcareService", e, config["HealthcareService"], dal, serverConfig)
	RegisterController("DataElement", e, config["DataElement"], dal, serverConfig)
	RegisterController("DeviceComponent", e, config["DeviceComponent"], dal, serverConfig)
	RegisterController("FamilyMemberHistory", e, config["FamilyMemberHistory"], dal, serverConfig)
	RegisterController("NutritionOrder", e, config["NutritionOrder"], dal, serverConfig)
	RegisterController("Encounter", e, config["Encounter"], dal, serverConfig)
	RegisterController("Substance", e, config["Substance"], dal, serverConfig)
	RegisterController("AuditEvent", e, config["AuditEvent"], dal, serverConfig)
	RegisterController("MedicationOrder", e, config["MedicationOrder"], dal, serverConfig)
	RegisterController("SearchParameter", e, config["SearchParameter"], dal, serverConfig)
	RegisterController("PaymentReconciliation", e, config["PaymentReconciliation"], dal, serverConfig)
	RegisterController("Communication", e, config["Communication"], dal, serverConfig)
	RegisterController("Condition", e, config["Condition"], dal, serverConfig)
	RegisterController("Composition", e, config["Composition"], dal, serverConfig)
	RegisterController("DetectedIssue", e, config["DetectedIssue"], dal, serverConfig)
	RegisterController("Bundle", e, config["Bundle"], dal, serverConfig)
	RegisterController("DiagnosticOrder", e, config["DiagnosticOrder"], dal, serverConfig)
	RegisterController("Patient", e, config["Patient"], dal, serverConfig)
	RegisterController("OrderResponse", e, config["OrderResponse"], dal, serverConfig)
	RegisterController("Coverage", e, config["Coverage"], dal, serverConfig)
	RegisterController("QuestionnaireResponse", e, config["QuestionnaireResponse"], dal, serverConfig)
	RegisterController("DeviceUseStatement", e, config["DeviceUseStatement"], dal, serverConfig)
	RegisterController("ProcessResponse", e, config["ProcessResponse"], dal, serverConfig)
	RegisterController("NamingSystem", e, config["NamingSystem"], dal, serverConfig)
	RegisterController("Schedule", e, config["Schedule"], dal, serverConfig)
	RegisterController("SupplyDelivery", e, config["SupplyDelivery"], dal, serverConfig)
	RegisterController("ClinicalImpression", e, config["ClinicalImpression"], dal, serverConfig)
	RegisterController("MessageHeader", e, config["MessageHeader"], dal, serverConfig)
	RegisterController("Claim", e, config["Claim"], dal, serverConfig)
	RegisterController("ImmunizationRecommendation", e, config["ImmunizationRecommendation"], dal, serverConfig)
	RegisterController("Location", e, config["Location"], dal, serverConfig)
	RegisterController("BodySite", e, config["BodySite"], dal, serverConfig)
}
