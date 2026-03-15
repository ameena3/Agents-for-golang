module github.com/microsoft/agents-sdk-go

go 1.24

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.16.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.8.0
	github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos v1.1.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.4.0
	github.com/AzureAD/microsoft-authentication-library-for-go v1.3.2
	github.com/golang-jwt/jwt/v5 v5.2.1
)

require (
	github.com/Azure/azure-sdk-for-go v68.0.0+incompatible // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)

// Replace the monolithic azure-sdk-for-go (required transitively by azcosmos)
// with a local empty stub so the module graph resolves in offline environments.
replace github.com/Azure/azure-sdk-for-go v68.0.0+incompatible => ./internal/azure-sdk-legacy
