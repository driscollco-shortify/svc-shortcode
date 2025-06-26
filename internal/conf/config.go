package conf

import "time"

var Config config

type config struct {
	Behaviours struct {
		ShortCodes struct {
			ExtendLifetime struct {
				RemainingClicks struct {
					TriggerPoint int
					TopUp        int
				}
			}
		}
		SafeUrlChecks struct {
			CheckInterval time.Duration
		}
	}
	CloudFlare struct {
		Credentials struct {
			ApiKey string
			ZoneId string
		}
	}
	GCP struct {
		Credentials []byte
		FireStore   struct {
			Paths struct {
				ShortCodes       string
				ShortCodesLegacy string
				Stats            string
			}
		}
		ProjectId string `required:"true"`
		SafeSite  struct {
			ApiKey string
		}
	}
	Router struct {
		Port int `default:"80"`
	}
	Service struct {
		Handlers struct {
			Create struct {
				Domains []string
				Clicks  struct {
					Max int
				}
				FlowClient struct {
					HardRateLimitId string
					RateLimitId     string
					TargetId        string
				}
				Jwt struct {
					PublicKey  string
					PrivateKey string
				}
				Lifetime struct {
					Expiry time.Duration
				}
				PubSub struct {
					Topics struct {
						Created string
					}
				}
				ShortCode struct {
					Length int
				}
			}
		}
	}
	Stats struct {
		Year struct {
			StartDate string
		}
	}
}
