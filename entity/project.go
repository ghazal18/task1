package entity

type Project struct {
	ID          int    `pg:"id"`
	OwnerID     int    `pg:"owner_id"`
	Name        string `pg:"name"`
	Company     string `pg:"company"`
	Description string `pg:"description"`
	SocialLinks map[string]string `pg:"social_links"`
}
