// place files you want to import through the `$lib` alias in this folder.

export { apiClient, API_BASE_URL } from './clients/apiClient';

// Export types and functions from each API module
export type {
	SocialMediaPost,
	PostFilters,
	DashboardData,
	Platform,
	ContentFormat,
	PostStatus
} from './api/socialmediaposts';
export {
	listSocialMediaPosts,
	getSocialMediaPost,
	createSocialMediaPost,
	updateSocialMediaPost,
	updateSocialMediaPostStatus,
	deleteSocialMediaPost,
	getDashboardData
} from './api/socialmediaposts';

export type { ContentPost, ContentPostWithSocialPosts, CreateContentPostInput, RepurposeRequest } from './api/content';
export {
	listContentPosts,
	getContentPost,
	createContentPost,
	updateContentPostPriority,
	repurposeToPlatforms,
	getContentBacklog
} from './api/content';

export type { ScheduledPost, SchedulePostRequest, UpdateScheduleRequest } from './api/scheduling';
export {
	listScheduledPosts,
	getUpcomingScheduledPosts,
	getScheduledPost,
	schedulePost,
	updateSchedule,
	cancelSchedule,
	autoSchedule
} from './api/scheduling';

export type {
	PostAnalytics,
	AnalyticsSummary,
	TopPost,
	RecordAnalyticsRequest
} from './api/analytics';
export {
	getPostAnalytics,
	getAnalyticsSummary,
	getTopPosts,
	recordAnalytics
} from './api/analytics';

export type {
	PlatformConfig,
	OptimalTimesResponse,
	UpdatePlatformConfigRequest
} from './api/platforms';
export {
	listPlatforms,
	getPlatformConfig,
	getPlatformConfigByName,
	getOptimalTimes,
	updatePlatformConfig
} from './api/platforms';

export type { ContentAsset, CreateAssetRequest, UpdateAssetRequest } from './api/assets';
export {
	listAssets,
	getAsset,
	getAssetsByContentPost,
	getAssetsBySocialPost,
	createAsset,
	updateAsset,
	deleteAsset
} from './api/assets';
