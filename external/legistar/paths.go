package legistar

const (
	legistarBase       = "https://webapi.legistar.com/v1/%s" // %s is city/state/etc name
	matters            = legistarBase + "/matters"
	matter             = matters + "/%d"
	matterHistory      = matter + "/histories"
	matterTextVersions = matter + "/versions"
	matterText         = matter + "/texts/%s"
	person             = legistarBase + "/persons"
	personVote         = person + "/%d/votes"
	events             = legistarBase + "/Events"
	eventItems         = events + "/%d/EventItems/%d"
)
