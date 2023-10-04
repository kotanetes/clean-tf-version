package main

type TFModel struct {
	Locals []Local `hcl:"locals"`
}

type Local struct {
	SoftwareReleases any `hcl:"software_releases"`
}

type TFModelJson struct {
	SoftwareRelease     string            `json:"software_release"`
	ApplicationVersions map[string]string `json:"application_versions"`
	ConfigVersion       map[string]string `json:"config_version"`
}

type TenantResult struct {
	Region  string `json:"region"`
	Tenant  string `json:"tenant"`
	Version string `json:"version"`
}
