// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package core

import (
	"github.com/ameena3/Agents-for-golang/activity"
	"github.com/ameena3/Agents-for-golang/activity/card"
)

// HeroCard creates an Attachment containing the given HeroCard.
func HeroCard(c *card.HeroCard) *activity.Attachment {
	return &activity.Attachment{
		ContentType: card.ContentTypeHeroCard,
		Content:     c,
	}
}

// ThumbnailCard creates an Attachment containing the given ThumbnailCard.
func ThumbnailCard(c *card.ThumbnailCard) *activity.Attachment {
	return &activity.Attachment{
		ContentType: card.ContentTypeThumbnailCard,
		Content:     c,
	}
}

// OAuthCard creates an Attachment containing the given OAuthCard.
func OAuthCard(c *card.OAuthCard) *activity.Attachment {
	return &activity.Attachment{
		ContentType: card.ContentTypeOAuthCard,
		Content:     c,
	}
}

// SigninCard creates an Attachment containing the given SigninCard.
func SigninCard(c *card.SigninCard) *activity.Attachment {
	return &activity.Attachment{
		ContentType: card.ContentTypeSigninCard,
		Content:     c,
	}
}

// AdaptiveCard creates an Attachment for an Adaptive Card. The content parameter
// should be a map[string]interface{} or any JSON-serializable struct representing
// the Adaptive Card schema payload.
func AdaptiveCard(content interface{}) *activity.Attachment {
	return &activity.Attachment{
		ContentType: card.ContentTypeAdaptiveCard,
		Content:     content,
	}
}
