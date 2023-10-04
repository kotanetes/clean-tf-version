locals { 
    software_releases = {
         2023.4 = { 
            config_versions = { 
                messaging-services = "1.3.37"
            }
        },
        2023.5 = {
            config_versions = { 
                messaging-services = "1.3.38"
            }
        }
    }
}