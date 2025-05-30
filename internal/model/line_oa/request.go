package line_oa

type CreateRequest struct {
	Name               string `json:"name" binding:"required,max=100"`
	ChannelID          string `json:"channel_id" binding:"required,max=100"`
	ChannelSecret      string `json:"channel_secret" binding:"required"`
	ChannelAccessToken string `json:"channel_access_token" binding:"required"`
}

type UpdateRequest struct {
	Name               string `json:"name" binding:"required,max=100"`
	ChannelID          string `json:"channel_id" binding:"required,max=100"`
	ChannelSecret      string `json:"channel_secret" binding:"required"`
	ChannelAccessToken string `json:"channel_access_token" binding:"required"`
}
