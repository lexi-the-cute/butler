package profile

import (
	"github.com/go-errors/errors"
	"github.com/itchio/butler/buse"
	"github.com/itchio/butler/buse/messages"
	"github.com/itchio/butler/database/models"
	"github.com/itchio/go-itchio"
)

func Register(router *buse.Router) {
	messages.ProfileList.Register(router, List)
	messages.ProfileLoginWithPassword.Register(router, LoginWithPassword)
	messages.ProfileLoginWithAPIKey.Register(router, LoginWithAPIKey)
	messages.ProfileUseSavedLogin.Register(router, UseSavedLogin)
	messages.ProfileForget.Register(router, Forget)
	messages.ProfileDataPut.Register(router, DataPut)
	messages.ProfileDataGet.Register(router, DataGet)
}

func List(rc *buse.RequestContext, params *buse.ProfileListParams) (*buse.ProfileListResult, error) {
	var profiles []*models.Profile
	err := rc.DB().Preload("User").Order("last_connected desc").Find(&profiles).Error
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	var formattedProfiles []*buse.Profile
	for _, profile := range profiles {
		formattedProfiles = append(formattedProfiles, formatProfile(profile))
	}

	return &buse.ProfileListResult{
		Profiles: formattedProfiles,
	}, nil
}

func formatProfile(p *models.Profile) *buse.Profile {
	return &buse.Profile{
		ID:            p.ID,
		LastConnected: p.LastConnected,
		User:          p.User,
	}
}

func LoginWithPassword(rc *buse.RequestContext, params *buse.ProfileLoginWithPasswordParams) (*buse.ProfileLoginWithPasswordResult, error) {
	if params.Username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if params.Password == "" {
		return nil, errors.New("password cannot be empty")
	}

	rootClient, err := rc.RootClient()
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	var key *itchio.APIKey
	var cookie itchio.Cookie

	{
		loginRes, err := rootClient.LoginWithPassword(&itchio.LoginWithPasswordParams{
			Username: params.Username,
			Password: params.Password,
		})
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}

		if loginRes.RecaptchaNeeded {
			// Captcha flow
			recaptchaRes, err := messages.ProfileRequestCaptcha.Call(rc, &buse.ProfileRequestCaptchaParams{
				RecaptchaURL: loginRes.RecaptchaURL,
			})
			if err != nil {
				return nil, errors.Wrap(err, 0)
			}

			if recaptchaRes.RecaptchaResponse == "" {
				return nil, &buse.ErrAborted{}
			}

			loginRes, err = rootClient.LoginWithPassword(&itchio.LoginWithPasswordParams{
				Username:          params.Username,
				Password:          params.Password,
				RecaptchaResponse: recaptchaRes.RecaptchaResponse,
			})
			if err != nil {
				return nil, errors.Wrap(err, 0)
			}
		}

		if loginRes.Token != "" {
			// TOTP flow
			totpRes, err := messages.ProfileRequestTOTP.Call(rc, &buse.ProfileRequestTOTPParams{})
			if err != nil {
				return nil, errors.Wrap(err, 0)
			}

			if totpRes.Code == "" {
				return nil, &buse.ErrAborted{}
			}

			verifyRes, err := rootClient.TOTPVerify(&itchio.TOTPVerifyParams{
				Token: loginRes.Token,
				Code:  totpRes.Code,
			})
			if err != nil {
				return nil, errors.Wrap(err, 0)
			}

			key = verifyRes.Key
			cookie = verifyRes.Cookie
		} else {
			// One-factor flow
			key = loginRes.Key
			cookie = loginRes.Cookie
		}
	}

	client, err := rc.KeyClient(key.Key)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	meRes, err := client.GetMe()
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	profile := &models.Profile{
		ID:     meRes.User.ID,
		APIKey: key.Key,
	}
	profile.UpdateFromUser(meRes.User)

	err = rc.DB().Save(profile).Error
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	res := &buse.ProfileLoginWithPasswordResult{
		Cookie:  cookie,
		Profile: formatProfile(profile),
	}
	return res, nil
}

func LoginWithAPIKey(rc *buse.RequestContext, params *buse.ProfileLoginWithAPIKeyParams) (*buse.ProfileLoginWithAPIKeyResult, error) {
	if params.APIKey == "" {
		return nil, errors.New("apiKey cannot be empty")
	}

	client, err := rc.KeyClient(params.APIKey)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	meRes, err := client.GetMe()
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	profile := &models.Profile{
		ID:     meRes.User.ID,
		APIKey: params.APIKey,
	}
	profile.UpdateFromUser(meRes.User)

	err = rc.DB().Save(profile).Error
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	res := &buse.ProfileLoginWithAPIKeyResult{
		Profile: formatProfile(profile),
	}
	return res, nil
}

func UseSavedLogin(rc *buse.RequestContext, params *buse.ProfileUseSavedLoginParams) (*buse.ProfileUseSavedLoginResult, error) {
	consumer := rc.Consumer

	profile, client := rc.ProfileClient(params.ProfileID)

	consumer.Opf("Validating credentials...")

	meRes, err := client.GetMe()
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	profile.UpdateFromUser(meRes.User)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	err = rc.DB().Save(profile).Error
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	consumer.Opf("Logged in!")

	res := &buse.ProfileUseSavedLoginResult{
		Profile: formatProfile(profile),
	}
	return res, nil
}

func Forget(rc *buse.RequestContext, params *buse.ProfileForgetParams) (*buse.ProfileForgetResult, error) {
	if params.ProfileID == 0 {
		return nil, errors.New("profileId must be set")
	}

	success := false

	profile := models.ProfileByID(rc.DB(), params.ProfileID)
	if profile != nil {
		err := rc.DB().Delete(profile).Error
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}
		success = true
	}

	res := &buse.ProfileForgetResult{
		Success: success,
	}
	return res, nil
}
