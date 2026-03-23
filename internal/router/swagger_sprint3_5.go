package router

// @Summary Get prayer time settings
// @Tags Worship
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tenant/prayer-time-settings [get]
func swaggerTenantPrayerTimeSettingsGet() {}

// @Summary Upsert prayer time settings
// @Tags Worship
// @Security BearerAuth
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/prayer-time-settings [put]
func swaggerTenantPrayerTimeSettingsPut() {}

// @Summary List prayer times daily
// @Tags Worship
// @Security BearerAuth
// @Param from query string false "from"
// @Param to query string false "to"
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/prayer-times-daily [get]
func swaggerTenantPrayerTimesDailyList() {}

// @Summary Create prayer times daily
// @Tags Worship
// @Security BearerAuth
// @Param payload body object true "payload"
// @Success 201 {object} map[string]interface{}
// @Router /tenant/prayer-times-daily [post]
func swaggerTenantPrayerTimesDailyCreate() {}

// @Summary Get prayer times daily
// @Tags Worship
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/prayer-times-daily/{id} [get]
func swaggerTenantPrayerTimesDailyGet() {}

// @Summary Update prayer times daily
// @Tags Worship
// @Security BearerAuth
// @Param id path int true "id"
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/prayer-times-daily/{id} [put]
func swaggerTenantPrayerTimesDailyUpdate() {}

// @Summary Delete prayer times daily
// @Tags Worship
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/prayer-times-daily/{id} [delete]
func swaggerTenantPrayerTimesDailyDelete() {}

// @Summary List prayer duties
// @Tags Worship
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tenant/prayer-duties [get]
func swaggerTenantPrayerDutiesList() {}

// @Summary Create prayer duty
// @Tags Worship
// @Security BearerAuth
// @Param payload body object true "payload"
// @Success 201 {object} map[string]interface{}
// @Router /tenant/prayer-duties [post]
func swaggerTenantPrayerDutiesCreate() {}

// @Summary Get prayer duty
// @Tags Worship
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/prayer-duties/{id} [get]
func swaggerTenantPrayerDutiesGet() {}

// @Summary Update prayer duty
// @Tags Worship
// @Security BearerAuth
// @Param id path int true "id"
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/prayer-duties/{id} [put]
func swaggerTenantPrayerDutiesUpdate() {}

// @Summary Delete prayer duty
// @Tags Worship
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/prayer-duties/{id} [delete]
func swaggerTenantPrayerDutiesDelete() {}

// @Summary List special days
// @Tags Worship
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tenant/special-days [get]
func swaggerTenantSpecialDaysList() {}

// @Summary Create special day
// @Tags Worship
// @Security BearerAuth
// @Param payload body object true "payload"
// @Success 201 {object} map[string]interface{}
// @Router /tenant/special-days [post]
func swaggerTenantSpecialDaysCreate() {}

// @Summary Get special day
// @Tags Worship
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/special-days/{id} [get]
func swaggerTenantSpecialDaysGet() {}

// @Summary Update special day
// @Tags Worship
// @Security BearerAuth
// @Param id path int true "id"
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/special-days/{id} [put]
func swaggerTenantSpecialDaysUpdate() {}

// @Summary Delete special day
// @Tags Worship
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/special-days/{id} [delete]
func swaggerTenantSpecialDaysDelete() {}

// @Summary Get prayer calendar
// @Tags Worship
// @Security BearerAuth
// @Param from query string false "from"
// @Param to query string false "to"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/prayer-calendar [get]
func swaggerTenantPrayerCalendarGet() {}

// @Summary List events
// @Tags Events
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tenant/events [get]
func swaggerTenantEventsList() {}

// @Summary Create event
// @Tags Events
// @Security BearerAuth
// @Param payload body object true "payload"
// @Success 201 {object} map[string]interface{}
// @Router /tenant/events [post]
func swaggerTenantEventsCreate() {}

// @Summary Get event
// @Tags Events
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/events/{id} [get]
func swaggerTenantEventsGet() {}

// @Summary Update event
// @Tags Events
// @Security BearerAuth
// @Param id path int true "id"
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/events/{id} [put]
func swaggerTenantEventsUpdate() {}

// @Summary Update event status
// @Tags Events
// @Security BearerAuth
// @Param id path int true "id"
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/events/{id}/status [patch]
func swaggerTenantEventsUpdateStatus() {}

// @Summary Delete event
// @Tags Events
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/events/{id} [delete]
func swaggerTenantEventsDelete() {}

// @Summary List gallery albums
// @Tags Gallery
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tenant/gallery/albums [get]
func swaggerTenantGalleryAlbumsList() {}

// @Summary Create gallery album
// @Tags Gallery
// @Security BearerAuth
// @Param payload body object true "payload"
// @Success 201 {object} map[string]interface{}
// @Router /tenant/gallery/albums [post]
func swaggerTenantGalleryAlbumsCreate() {}

// @Summary Get gallery album
// @Tags Gallery
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/gallery/albums/{id} [get]
func swaggerTenantGalleryAlbumsGet() {}

// @Summary Update gallery album
// @Tags Gallery
// @Security BearerAuth
// @Param id path int true "id"
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/gallery/albums/{id} [put]
func swaggerTenantGalleryAlbumsUpdate() {}

// @Summary Delete gallery album
// @Tags Gallery
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/gallery/albums/{id} [delete]
func swaggerTenantGalleryAlbumsDelete() {}

// @Summary List gallery items
// @Tags Gallery
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tenant/gallery/items [get]
func swaggerTenantGalleryItemsList() {}

// @Summary Create gallery item
// @Tags Gallery
// @Security BearerAuth
// @Param payload body object true "payload"
// @Success 201 {object} map[string]interface{}
// @Router /tenant/gallery/items [post]
func swaggerTenantGalleryItemsCreate() {}

// @Summary Get gallery item
// @Tags Gallery
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/gallery/items/{id} [get]
func swaggerTenantGalleryItemsGet() {}

// @Summary Update gallery item
// @Tags Gallery
// @Security BearerAuth
// @Param id path int true "id"
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/gallery/items/{id} [put]
func swaggerTenantGalleryItemsUpdate() {}

// @Summary Delete gallery item
// @Tags Gallery
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/gallery/items/{id} [delete]
func swaggerTenantGalleryItemsDelete() {}

// @Summary List management members
// @Tags ManagementMembers
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tenant/management-members [get]
func swaggerTenantManagementMembersList() {}

// @Summary Create management member
// @Tags ManagementMembers
// @Security BearerAuth
// @Param payload body object true "payload"
// @Success 201 {object} map[string]interface{}
// @Router /tenant/management-members [post]
func swaggerTenantManagementMembersCreate() {}

// @Summary Get management member
// @Tags ManagementMembers
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/management-members/{id} [get]
func swaggerTenantManagementMembersGet() {}

// @Summary Update management member
// @Tags ManagementMembers
// @Security BearerAuth
// @Param id path int true "id"
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/management-members/{id} [put]
func swaggerTenantManagementMembersUpdate() {}

// @Summary Delete management member
// @Tags ManagementMembers
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/management-members/{id} [delete]
func swaggerTenantManagementMembersDelete() {}

// @Summary List donation channels
// @Tags Engagement
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tenant/donation-channels [get]
func swaggerTenantDonationChannelsList() {}

// @Summary Create donation channel
// @Tags Engagement
// @Security BearerAuth
// @Param payload body object true "payload"
// @Success 201 {object} map[string]interface{}
// @Router /tenant/donation-channels [post]
func swaggerTenantDonationChannelsCreate() {}

// @Summary Get donation channel
// @Tags Engagement
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/donation-channels/{id} [get]
func swaggerTenantDonationChannelsGet() {}

// @Summary Update donation channel
// @Tags Engagement
// @Security BearerAuth
// @Param id path int true "id"
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/donation-channels/{id} [put]
func swaggerTenantDonationChannelsUpdate() {}

// @Summary Delete donation channel
// @Tags Engagement
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/donation-channels/{id} [delete]
func swaggerTenantDonationChannelsDelete() {}

// @Summary List social links
// @Tags Engagement
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tenant/social-links [get]
func swaggerTenantSocialLinksList() {}

// @Summary Create social link
// @Tags Engagement
// @Security BearerAuth
// @Param payload body object true "payload"
// @Success 201 {object} map[string]interface{}
// @Router /tenant/social-links [post]
func swaggerTenantSocialLinksCreate() {}

// @Summary Get social link
// @Tags Engagement
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/social-links/{id} [get]
func swaggerTenantSocialLinksGet() {}

// @Summary Update social link
// @Tags Engagement
// @Security BearerAuth
// @Param id path int true "id"
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/social-links/{id} [put]
func swaggerTenantSocialLinksUpdate() {}

// @Summary Delete social link
// @Tags Engagement
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/social-links/{id} [delete]
func swaggerTenantSocialLinksDelete() {}

// @Summary List external links
// @Tags Engagement
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tenant/external-links [get]
func swaggerTenantExternalLinksList() {}

// @Summary Create external link
// @Tags Engagement
// @Security BearerAuth
// @Param payload body object true "payload"
// @Success 201 {object} map[string]interface{}
// @Router /tenant/external-links [post]
func swaggerTenantExternalLinksCreate() {}

// @Summary Get external link
// @Tags Engagement
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/external-links/{id} [get]
func swaggerTenantExternalLinksGet() {}

// @Summary Update external link
// @Tags Engagement
// @Security BearerAuth
// @Param id path int true "id"
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/external-links/{id} [put]
func swaggerTenantExternalLinksUpdate() {}

// @Summary Delete external link
// @Tags Engagement
// @Security BearerAuth
// @Param id path int true "id"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/external-links/{id} [delete]
func swaggerTenantExternalLinksDelete() {}

// @Summary List feature catalog
// @Tags Engagement
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tenant/feature-catalog [get]
func swaggerTenantFeatureCatalogList() {}

// @Summary List website features
// @Tags Engagement
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tenant/website-features [get]
func swaggerTenantWebsiteFeaturesList() {}

// @Summary Update website feature
// @Tags Engagement
// @Security BearerAuth
// @Param feature_id path int true "feature_id"
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/website-features/{feature_id} [put]
func swaggerTenantWebsiteFeaturesUpdate() {}

// @Summary Bulk update website features
// @Tags Engagement
// @Security BearerAuth
// @Param payload body object true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /tenant/website-features/bulk [patch]
func swaggerTenantWebsiteFeaturesBulkUpdate() {}

// @Summary Public events
// @Tags Public
// @Param hostname path string true "hostname"
// @Success 200 {object} map[string]interface{}
// @Router /public/{hostname}/events [get]
func swaggerPublicEventsList() {}

// @Summary Public gallery albums
// @Tags Public
// @Param hostname path string true "hostname"
// @Success 200 {object} map[string]interface{}
// @Router /public/{hostname}/gallery/albums [get]
func swaggerPublicGalleryAlbumsList() {}

// @Summary Public gallery items
// @Tags Public
// @Param hostname path string true "hostname"
// @Success 200 {object} map[string]interface{}
// @Router /public/{hostname}/gallery/items [get]
func swaggerPublicGalleryItemsList() {}

// @Summary Public management members
// @Tags Public
// @Param hostname path string true "hostname"
// @Success 200 {object} map[string]interface{}
// @Router /public/{hostname}/management-members [get]
func swaggerPublicManagementMembersList() {}

// @Summary Public donation channels
// @Tags Public
// @Param hostname path string true "hostname"
// @Success 200 {object} map[string]interface{}
// @Router /public/{hostname}/donation-channels [get]
func swaggerPublicDonationChannelsList() {}

// @Summary Public social links
// @Tags Public
// @Param hostname path string true "hostname"
// @Success 200 {object} map[string]interface{}
// @Router /public/{hostname}/social-links [get]
func swaggerPublicSocialLinksList() {}

// @Summary Public external links
// @Tags Public
// @Param hostname path string true "hostname"
// @Success 200 {object} map[string]interface{}
// @Router /public/{hostname}/external-links [get]
func swaggerPublicExternalLinksList() {}
